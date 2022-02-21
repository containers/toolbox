/*
 * Copyright © 2019 – 2021 Red Hat Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package utils

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/acobaugh/osrelease"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/docker/go-units"
	"github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

type ParseReleaseFunc func(string) (string, error)

type Distro struct {
	ContainerNamePrefix    string
	ImageBasename          string
	ParseRelease           ParseReleaseFunc
	ReleaseFormat          string
	Registry               string
	Repository             string
	RepositoryNeedsRelease bool
}

const (
	idTruncLength          = 12
	releaseDefaultFallback = "34"
)

const (
	// Based on the nameRegex value in:
	// https://github.com/containers/libpod/blob/master/libpod/options.go
	ContainerNameRegexp = "[a-zA-Z0-9][a-zA-Z0-9_.-]*"
)

var (
	ErrUnsupportedDistro = errors.New("linux distribution is not supported")
)

var (
	containerNamePrefixDefault = "fedora-toolbox"

	distroDefault = "fedora"

	preservedEnvironmentVariables = []string{
		"COLORTERM",
		"DBUS_SESSION_BUS_ADDRESS",
		"DBUS_SYSTEM_BUS_ADDRESS",
		"DESKTOP_SESSION",
		"DISPLAY",
		"LANG",
		"SHELL",
		"SSH_AUTH_SOCK",
		"TERM",
		"TOOLBOX_PATH",
		"USER",
		"VTE_VERSION",
		"WAYLAND_DISPLAY",
		"XAUTHORITY",
		"XDG_CURRENT_DESKTOP",
		"XDG_DATA_DIRS",
		"XDG_MENU_PREFIX",
		"XDG_RUNTIME_DIR",
		"XDG_SEAT",
		"XDG_SESSION_DESKTOP",
		"XDG_SESSION_ID",
		"XDG_SESSION_TYPE",
		"XDG_VTNR",
	}

	releaseDefault string

	supportedDistros = map[string]Distro{
		"fedora": {
			"fedora-toolbox",
			"fedora-toolbox",
			parseReleaseFedora,
			"<release>/f<release>",
			"registry.fedoraproject.org",
			"",
			false,
		},
		"rhel": {
			"rhel-toolbox",
			"toolbox",
			parseReleaseRHEL,
			"<major.minor>",
			"registry.access.redhat.com",
			"ubi8",
			false,
		},
	}
)

var (
	ContainerNameDefault string
)

func init() {
	releaseDefault = releaseDefaultFallback

	hostID, err := GetHostID()
	if err == nil {
		if distroObj, supportedDistro := supportedDistros[hostID]; supportedDistro {
			release, err := GetHostVersionID()
			if err == nil {
				containerNamePrefixDefault = distroObj.ContainerNamePrefix
				distroDefault = hostID
				releaseDefault = release
			}
		}
	}

	ContainerNameDefault = containerNamePrefixDefault + "-" + releaseDefault
}

func CallFlatpakSessionHelper() (string, error) {
	logrus.Debug("Calling org.freedesktop.Flatpak.SessionHelper.RequestSession")

	connection, err := dbus.SessionBus()
	if err != nil {
		return "", errors.New("failed to connect to the D-Bus session instance")
	}

	sessionHelper := connection.Object("org.freedesktop.Flatpak", "/org/freedesktop/Flatpak/SessionHelper")
	call := sessionHelper.Call("org.freedesktop.Flatpak.SessionHelper.RequestSession", 0)

	var result map[string]dbus.Variant
	err = call.Store(&result)
	if err != nil {
		logrus.Debugf("Call to org.freedesktop.Flatpak.SessionHelper.RequestSession failed: %s", err)
		return "", errors.New("failed to call org.freedesktop.Flatpak.SessionHelper.RequestSession")
	}

	pathVariant := result["path"]
	pathVariantSignature := pathVariant.Signature().String()
	if pathVariantSignature != "s" {
		return "", errors.New("unknown reply from org.freedesktop.Flatpak.SessionHelper.RequestSession")
	}

	pathValue := pathVariant.Value()
	path := pathValue.(string)
	return path, nil
}

func EnsureXdgRuntimeDirIsSet(uid int) {
	if _, ok := os.LookupEnv("XDG_RUNTIME_DIR"); !ok {
		logrus.Debug("XDG_RUNTIME_DIR is unset")

		xdgRuntimeDir := fmt.Sprintf("/run/user/%d", uid)
		os.Setenv("XDG_RUNTIME_DIR", xdgRuntimeDir)

		logrus.Debugf("XDG_RUNTIME_DIR set to %s", xdgRuntimeDir)
	}
}

func ForwardToHost() (int, error) {
	envOptions := GetEnvOptionsForPreservedVariables()
	toolboxPath := os.Getenv("TOOLBOX_PATH")
	commandLineArgs := os.Args[1:]

	var flatpakSpawnArgs []string

	flatpakSpawnArgs = append(flatpakSpawnArgs, envOptions...)

	flatpakSpawnArgs = append(flatpakSpawnArgs, []string{
		"--host",
		toolboxPath,
	}...)

	flatpakSpawnArgs = append(flatpakSpawnArgs, commandLineArgs...)

	logrus.Debug("Forwarding to host:")
	logrus.Debugf("%s", toolboxPath)
	for _, arg := range commandLineArgs {
		logrus.Debugf("%s", arg)
	}

	exitCode, err := shell.RunWithExitCode("flatpak-spawn", os.Stdin, os.Stdout, nil, flatpakSpawnArgs...)
	if err != nil {
		return exitCode, err
	}

	return exitCode, nil
}

// GetCgroupsVersion returns the cgroups version of the host
//
// Based on the IsCgroup2UnifiedMode function in:
// https://github.com/containers/libpod/tree/master/pkg/cgroups
func GetCgroupsVersion() (int, error) {
	var st syscall.Statfs_t

	if err := syscall.Statfs("/sys/fs/cgroup", &st); err != nil {
		return -1, err
	}

	version := 1
	if st.Type == unix.CGROUP2_SUPER_MAGIC {
		version = 2
	}

	return version, nil
}

func GetContainerNamePrefixForImage(image string) (string, error) {
	basename := ImageReferenceGetBasename(image)
	if basename == "" {
		return "", fmt.Errorf("failed to get the basename of image %s", image)
	}

	for _, distroObj := range supportedDistros {
		if distroObj.ImageBasename != basename {
			continue
		}

		return distroObj.ContainerNamePrefix, nil
	}

	return basename, nil
}

func getDefaultImageForDistro(distro, release string) string {
	if !IsDistroSupported(distro) {
		distro = distroDefault
	}

	distroObj, supportedDistro := supportedDistros[distro]
	if !supportedDistro {
		panicMsg := fmt.Sprintf("failed to find %s in the list of supported distributions", distro)
		panic(panicMsg)
	}

	image := distroObj.ImageBasename + ":" + release
	return image
}

func GetEnvOptionsForPreservedVariables() []string {
	logrus.Debug("Creating list of environment variables to forward")

	var envOptions []string

	for _, variable := range preservedEnvironmentVariables {
		value, found := os.LookupEnv(variable)
		if !found {
			logrus.Debugf("%s is unset", variable)
			continue
		}

		logrus.Debugf("%s=%s", variable, value)
		envOptions = append(envOptions, fmt.Sprintf("--env=%s=%s", variable, value))
	}

	return envOptions
}

func GetFullyQualifiedImageFromDistros(image, release string) (string, error) {
	logrus.Debugf("Resolving fully qualified name for image %s from known registries", image)

	if ImageReferenceHasDomain(image) {
		return image, nil
	}

	basename := ImageReferenceGetBasename(image)
	if basename == "" {
		return "", fmt.Errorf("failed to get the basename of image %s", image)
	}

	for _, distroObj := range supportedDistros {
		if distroObj.ImageBasename != basename {
			continue
		}

		var repository string

		if distroObj.RepositoryNeedsRelease {
			repository = fmt.Sprintf(distroObj.Repository, release)
		} else {
			repository = distroObj.Repository
		}

		imageFull := distroObj.Registry

		if repository != "" {
			imageFull = imageFull + "/" + repository
		}

		imageFull = imageFull + "/" + image

		logrus.Debugf("Resolved image %s to %s", image, imageFull)

		return imageFull, nil
	}

	return "", fmt.Errorf("failed to resolve image %s", image)
}

// GetGroupForSudo returns the name of the sudoers group.
//
// Some distros call it 'sudo' (eg. Ubuntu) and some call it 'wheel' (eg. Fedora).
func GetGroupForSudo() (string, error) {
	logrus.Debug("Looking up group for sudo")

	groups := []string{"sudo", "wheel"}

	for _, group := range groups {
		if _, err := user.LookupGroup(group); err == nil {
			logrus.Debugf("Group for sudo is %s", group)
			return group, nil
		}
	}

	return "", errors.New("group for sudo not found")
}

// GetHostID returns the ID from the os-release files
//
// Examples:
// - host is Fedora, returned string is 'fedora'
func GetHostID() (string, error) {
	osRelease, err := osrelease.Read()
	if err != nil {
		return "", err
	}

	return osRelease["ID"], nil
}

// GetHostVariantID returns the VARIANT_ID from the os-release files
//
// Examples:
// - host is Fedora Workstation, returned string is 'workstation'
func GetHostVariantID() (string, error) {
	osRelease, err := osrelease.Read()
	if err != nil {
		return "", err
	}

	return osRelease["VARIANT_ID"], nil
}

// GetHostVersionID returns the VERSION_ID from the os-release files
//
// Examples:
// - host is Fedora 32, returned string is '32'
func GetHostVersionID() (string, error) {
	osRelease, err := osrelease.Read()
	if err != nil {
		return "", err
	}

	return osRelease["VERSION_ID"], nil
}

// GetMountPoint returns the mount point of a target.
func GetMountPoint(target string) (string, error) {
	var stdout strings.Builder

	if err := shell.Run("df", nil, &stdout, nil, "--output=target", target); err != nil {
		return "", err
	}

	output := stdout.String()
	options := strings.Split(output, "\n")
	if len(options) != 3 {
		return "", errors.New("unexpected output from df(1)")
	}

	mountPoint := strings.TrimSpace(options[1])
	return mountPoint, nil
}

// GetMountOptions returns the mount options of a target.
func GetMountOptions(target string) (string, error) {
	var stdout strings.Builder
	findMntArgs := []string{"--noheadings", "--output", "OPTIONS", target}

	if err := shell.Run("findmnt", nil, &stdout, nil, findMntArgs...); err != nil {
		return "", err
	}

	output := stdout.String()
	options := strings.Split(output, "\n")

	mountOptions := strings.TrimSpace(options[0])
	return mountOptions, nil
}

// GetReleaseFormat returns the format string signifying supported release
// version formats.
//
// distro should be value found under ID in os-release.
//
// If distro is unsupported an empty string is returned.
func GetReleaseFormat(distro string) string {
	if !IsDistroSupported(distro) {
		return ""
	}

	return supportedDistros[distro].ReleaseFormat
}

func GetRuntimeDirectory(targetUser *user.User) (string, error) {
	gid, err := strconv.Atoi(targetUser.Gid)
	if err != nil {
		return "", fmt.Errorf("failed to convert group ID to integer: %w", err)
	}

	uid, err := strconv.Atoi(targetUser.Uid)
	if err != nil {
		return "", fmt.Errorf("failed to convert user ID to integer: %w", err)
	}

	var runtimeDirectory string

	if uid == 0 {
		runtimeDirectory = "/run"
	} else {
		runtimeDirectory = os.Getenv("XDG_RUNTIME_DIR")
	}

	toolboxRuntimeDirectory := path.Join(runtimeDirectory, "toolbox")
	logrus.Debugf("Creating runtime directory %s", toolboxRuntimeDirectory)

	if err := os.MkdirAll(toolboxRuntimeDirectory, 0700); err != nil {
		wrapped_err := fmt.Errorf("failed to create runtime directory %s: %w",
			toolboxRuntimeDirectory,
			err)
		return "", wrapped_err
	}

	if err := os.Chown(toolboxRuntimeDirectory, uid, gid); err != nil {
		wrapped_err := fmt.Errorf("failed to change ownership of the runtime directory %s: %w",
			toolboxRuntimeDirectory,
			err)
		return "", wrapped_err
	}

	return toolboxRuntimeDirectory, nil
}

// GetSupportedDistros returns a list of supported distributions
func GetSupportedDistros() []string {
	var distros []string
	for d := range supportedDistros {
		distros = append(distros, d)
	}
	return distros
}

// HumanDuration accepts a Unix time value and converts it into a human readable
// string.
//
// Examples: "5 minutes ago", "2 hours ago", "3 days ago"
func HumanDuration(duration int64) string {
	return units.HumanDuration(time.Since(time.Unix(duration, 0))) + " ago"
}

// IsDistroSupported signifies if a distribution has a toolbx image for it.
//
// distro should be value found under ID in os-release.
func IsDistroSupported(distro string) bool {
	_, ok := supportedDistros[distro]
	return ok
}

// ImageReferenceCanBeID checks if 'image' might be the ID of an image
func ImageReferenceCanBeID(image string) bool {
	matched, err := regexp.MatchString("^[a-f0-9]{6,64}$", image)
	if err != nil {
		panic("regular expression for ID reference matching is invalid")
	}
	return matched
}

func ImageReferenceGetBasename(image string) string {
	var i int

	if ImageReferenceHasDomain(image) {
		i = strings.IndexRune(image, '/')
	}

	remainder := image[i:]
	j := strings.IndexRune(remainder, ':')
	if j == -1 {
		j = len(remainder)
	}

	path := remainder[:j]
	basename := filepath.Base(path)
	return basename
}

func ImageReferenceGetDomain(image string) string {
	if !ImageReferenceHasDomain(image) {
		return ""
	}

	i := strings.IndexRune(image, '/')
	domain := image[:i]
	return domain
}

func ImageReferenceGetTag(image string) string {
	var i int

	if ImageReferenceHasDomain(image) {
		i = strings.IndexRune(image, '/')
	}

	remainder := image[i:]
	j := strings.IndexRune(remainder, ':')
	if j == -1 {
		return ""
	}

	tag := remainder[j+1:]
	return tag
}

// ImageReferenceHasDomain checks if the provided image has a domain definition in it.
func ImageReferenceHasDomain(image string) bool {
	i := strings.IndexRune(image, '/')
	if i == -1 {
		return false
	}

	prefix := image[:i]

	// A domain should contain a top level domain name. An exception is 'localhost'
	if !strings.ContainsAny(prefix, ".:") && prefix != "localhost" {
		return false
	}

	return true
}

func SetUpConfiguration() error {
	logrus.Debug("Setting up configuration")

	configFiles := []string{
		"/etc/containers/toolbox.conf",
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Debugf("Setting up configuration: failed to get the user config directory: %s", err)
		return errors.New("failed to get the user config directory")
	}

	userConfigPath := userConfigDir + "/containers/toolbox.conf"
	configFiles = append(configFiles, []string{
		userConfigPath,
	}...)

	viper.SetConfigType("toml")

	for _, configFile := range configFiles {
		viper.SetConfigFile(configFile)

		if err := viper.MergeInConfig(); err != nil {
			// Seems like Viper's errors can't be examined with
			// errors.As.

			// Seems like Viper doesn't actually throw
			// viper.ConfigFileNotFoundError if a configuration
			// file is not found. We still check for it for the
			// sake of completion or in case Viper uses it in a
			// different version.
			_, ok := err.(viper.ConfigFileNotFoundError)
			if ok || os.IsNotExist(err) {
				logrus.Debugf("Setting up configuration: file %s not found", configFile)
				continue
			}

			if _, ok := err.(viper.ConfigParseError); ok {
				logrus.Debugf("Setting up configuration: failed to parse file %s: %s", configFile, err)
				return fmt.Errorf("failed to parse file %s", configFile)
			}

			logrus.Debugf("Setting up configuration: failed to read file %s: %s", configFile, err)
			return fmt.Errorf("failed to read file %s", configFile)
		}
	}

	image, release, err := ResolveImageName("", "", "")
	if err != nil {
		logrus.Debugf("Setting up configuration: failed to resolve image name: %s", err)
		return errors.New("failed to resolve image name")
	}

	container, err := ResolveContainerName("", image, release)
	if err != nil {
		logrus.Debugf("Setting up configuration: failed to resolve container name: %s", err)
		return errors.New("failed to resolve container name")
	}

	ContainerNameDefault = container

	return nil
}

// ShortID shortens provided id to first 12 characters.
func ShortID(id string) string {
	if len(id) > idTruncLength {
		return id[:idTruncLength]
	}
	return id
}

// ParseRelease assesses if the requested version of a distribution is in
// the correct format.
//
// If distro is an empty string, a default value (value under the
// 'general.distro' key in a config file or 'fedora') is assumed.
func ParseRelease(distro, release string) (string, error) {
	if distro == "" {
		distro, _ = ResolveDistro(distro)
	}

	if !IsDistroSupported(distro) {
		distro = distroDefault
	}

	distroObj, supportedDistro := supportedDistros[distro]
	if !supportedDistro {
		panicMsg := fmt.Sprintf("failed to find %s in the list of supported distributions", distro)
		panic(panicMsg)
	}

	parseRelease := distroObj.ParseRelease
	release, err := parseRelease(release)
	return release, err
}

func parseReleaseFedora(release string) (string, error) {
	if strings.HasPrefix(release, "F") || strings.HasPrefix(release, "f") {
		release = release[1:]
	}

	releaseN, err := strconv.Atoi(release)
	if err != nil {
		return "", err
	}

	if releaseN <= 0 {
		return "", errors.New("release must be a positive integer")
	}

	return release, nil
}

func parseReleaseRHEL(release string) (string, error) {
	if i := strings.IndexRune(release, '.'); i == -1 {
		return "", errors.New("release must have a '.'")
	}

	releaseN, err := strconv.ParseFloat(release, 32)
	if err != nil {
		return "", err
	}

	if releaseN <= 0 {
		return "", errors.New("release must be a positive number")
	}

	return release, nil
}

// PathExists wraps around os.Stat providing a nice interface for checking an existence of a path.
func PathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

// IsContainerNameValid checks if the name of a container matches the right pattern
func IsContainerNameValid(containerName string) bool {
	pattern := "^" + ContainerNameRegexp + "$"
	matched, err := regexp.MatchString(pattern, containerName)
	if err != nil {
		panicMsg := fmt.Sprintf("failed to parse regular expression for container name: %v", err)
		panic(panicMsg)
	}

	return matched
}

func IsInsideContainer() bool {
	return PathExists("/run/.containerenv")
}

func IsInsideToolboxContainer() bool {
	return PathExists("/run/.toolboxenv")
}

// ResolveContainerName standardizes the name of a container
//
// If no container name is specified then the name of the image will be used.
func ResolveContainerName(container, image, release string) (string, error) {
	logrus.Debug("Resolving container name")
	logrus.Debugf("Container: '%s'", container)
	logrus.Debugf("Image: '%s'", image)
	logrus.Debugf("Release: '%s'", release)

	if container == "" {
		var err error
		container, err = GetContainerNamePrefixForImage(image)
		if err != nil {
			return "", err
		}

		tag := ImageReferenceGetTag(image)
		if tag != "" {
			container = container + "-" + tag
		}
	}

	logrus.Debug("Resolved container name")
	logrus.Debugf("Container: '%s'", container)

	return container, nil
}

// ResolveDistro assess if the requested distribution is supported and provides
// a default value if none is requested.
//
// If distro is empty, and the "general.distro" key in a config file is set,
// the value is read from the config file. If the key is not set, the default
// value ('fedora') is used instead.
func ResolveDistro(distro string) (string, error) {
	logrus.Debug("Resolving distribution")
	logrus.Debugf("Distribution: %s", distro)

	if distro == "" {
		distro = distroDefault
		if viper.IsSet("general.distro") {
			distro = viper.GetString("general.distro")
		}
	}

	if !IsDistroSupported(distro) {
		return "", ErrUnsupportedDistro
	}

	logrus.Debug("Resolved distribution")
	logrus.Debugf("Distribution: %s", distro)

	return distro, nil
}

// ResolveImageName standardizes the name of an image.
//
// If no image name is specified then the base image will reflect the platform of the host (even the version).
//
// If the host system is unknown then the base image will be 'fedora-toolbox' with a default version
func ResolveImageName(distroCLI, imageCLI, releaseCLI string) (string, string, error) {
	logrus.Debug("Resolving image name")
	logrus.Debugf("Distribution (CLI): '%s'", distroCLI)
	logrus.Debugf("Image (CLI): '%s'", imageCLI)
	logrus.Debugf("Release (CLI): '%s'", releaseCLI)

	distro, image, release := distroCLI, imageCLI, releaseCLI

	if distroCLI == "" {
		distro, _ = ResolveDistro(distroCLI)
	}

	if distro != distroDefault && releaseCLI == "" && !viper.IsSet("general.release") {
		return "", "", fmt.Errorf("release not found for non-default distribution %s", distro)
	}

	if releaseCLI == "" {
		release = releaseDefault
		if viper.IsSet("general.release") {
			release = viper.GetString("general.release")
		}
	}

	if imageCLI == "" {
		image = getDefaultImageForDistro(distro, release)

		if viper.IsSet("general.image") && distroCLI == "" && releaseCLI == "" {
			image = viper.GetString("general.image")

			release = ImageReferenceGetTag(image)
			if release == "" {
				release = releaseDefault
				if viper.IsSet("general.release") {
					release = viper.GetString("general.release")
				}
			}
		}
	} else {
		release = ImageReferenceGetTag(image)
		if release == "" {
			release = releaseDefault
			if viper.IsSet("general.release") {
				release = viper.GetString("general.release")
			}
		}
	}

	logrus.Debug("Resolved image name")
	logrus.Debugf("Image: '%s'", image)
	logrus.Debugf("Release: '%s'", release)

	return image, release, nil
}

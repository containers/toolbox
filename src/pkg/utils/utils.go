/*
 * Copyright © 2019 – 2025 Red Hat Inc.
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

type GetDefaultReleaseFunc func() (string, error)
type GetFullyQualifiedImageFunc func(string, string) string
type ParseReleaseFunc func(string) (string, error)

type Distro struct {
	ContainerNamePrefix    string
	ImageBasename          string
	ReleaseRequired        bool
	GetDefaultRelease      GetDefaultReleaseFunc
	GetFullyQualifiedImage GetFullyQualifiedImageFunc
	ParseRelease           ParseReleaseFunc
}

type OptionValueSource int

const (
	optionValueDefault OptionValueSource = iota
	optionValueConfigFile
	optionValueCLI
)

const (
	containerNamePrefixFallback = "fedora-toolbox"
	distroFallback              = "fedora"
	idTruncLength               = 12
	releaseFallback             = "40"
)

const (
	// Based on the nameRegex value in:
	// https://github.com/containers/libpod/blob/master/libpod/options.go
	ContainerNameRegexp = "[a-zA-Z0-9][a-zA-Z0-9_.-]*"
)

var (
	containerNamePrefixDefault string

	distroDefault string

	preservedEnvironmentVariables = []string{
		"COLORTERM",
		"CONTAINERS_STORAGE_CONF",
		"DBUS_SESSION_BUS_ADDRESS",
		"DBUS_SYSTEM_BUS_ADDRESS",
		"DESKTOP_SESSION",
		"DISPLAY",
		"HISTCONTROL",
		"HISTFILE",
		"HISTFILESIZE",
		"HISTIGNORE",
		"HISTSIZE",
		"HISTTIMEFORMAT",
		"KONSOLE_VERSION",
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
		"XDG_SESSION_CLASS",
		"XDG_SESSION_DESKTOP",
		"XDG_SESSION_ID",
		"XDG_SESSION_TYPE",
		"XDG_VTNR",
		"XTERM_VERSION",
	}

	releaseDefault string

	runtimeDirectories map[string]string

	supportedDistros = map[string]Distro{
		"arch": {
			"arch-toolbox",
			"arch-toolbox",
			false,
			getDefaultReleaseArch,
			getFullyQualifiedImageArch,
			parseReleaseArch,
		},
		"fedora": {
			"fedora-toolbox",
			"fedora-toolbox",
			true,
			getDefaultReleaseFedora,
			getFullyQualifiedImageFedora,
			parseReleaseFedora,
		},
		"rhel": {
			"rhel-toolbox",
			"toolbox",
			true,
			getDefaultReleaseRHEL,
			getFullyQualifiedImageRHEL,
			parseReleaseRHEL,
		},
		"ubuntu": {
			"ubuntu-toolbox",
			"ubuntu-toolbox",
			true,
			getDefaultReleaseUbuntu,
			getFullyQualifiedImageUbuntu,
			parseReleaseUbuntu,
		},
	}
)

var (
	ContainerNameDefault string

	ErrContainerNameFromImageInvalid = errors.New("container name generated from image is invalid")

	ErrContainerNameInvalid = errors.New("container name is invalid")

	ErrDistroUnsupported = errors.New("distribution is unsupported")

	ErrDistroWithoutRelease = errors.New("non-default distribution must specify release")

	ErrImageWithoutBasename = errors.New("image does not have a basename")
)

func init() {
	containerNamePrefixDefault = containerNamePrefixFallback
	distroDefault = distroFallback
	releaseDefault = releaseFallback

	hostID, err := GetHostID()
	if err == nil {
		if distroObj, supportedDistro := supportedDistros[hostID]; supportedDistro {
			release, err := getDefaultReleaseForDistro(hostID)
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

	exitCode, err := shell.RunWithExitCode("flatpak-spawn", os.Stdin, os.Stdout, os.Stderr, flatpakSpawnArgs...)
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

func getContainerNamePrefixForImage(image string) (string, error) {
	basename := ImageReferenceGetBasename(image)
	if basename == "" {
		return "", &ImageError{image, ErrImageWithoutBasename}
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
	if distro == "" {
		panic("distro not specified")
	}

	distroObj, supportedDistro := supportedDistros[distro]
	if !supportedDistro {
		panicMsg := fmt.Sprintf("failed to find %s in the list of supported distributions", distro)
		panic(panicMsg)
	}

	image := distroObj.ImageBasename + ":" + release
	return image
}

func getDefaultReleaseForDistro(distro string) (string, error) {
	if distro == "" {
		panic("distro not specified")
	}

	distroObj, supportedDistro := supportedDistros[distro]
	if !supportedDistro {
		panicMsg := fmt.Sprintf("failed to find %s in the list of supported distributions", distro)
		panic(panicMsg)
	}

	release, err := distroObj.GetDefaultRelease()
	if err != nil {
		return "", err
	}

	return release, nil
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

	if image == "" {
		panic("image not specified")
	}

	if release == "" {
		panic("release not specified")
	}

	if tag := ImageReferenceGetTag(image); tag != "" && release != tag {
		panicMsg := fmt.Sprintf("image %s does not match release %s", image, release)
		panic(panicMsg)
	}

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

		getFullyQualifiedImageImpl := distroObj.GetFullyQualifiedImage
		imageFull := getFullyQualifiedImageImpl(image, release)

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

// getHostVersionID returns the VERSION_ID from the os-release files
//
// Examples:
// - host is Fedora 32, returned string is '32'
func getHostVersionID() (string, error) {
	osRelease, err := osrelease.Read()
	if err != nil {
		return "", err
	}

	return osRelease["VERSION_ID"], nil
}

func GetInitializedStamp(entryPointPID int, targetUser *user.User) (string, error) {
	toolbxRuntimeDirectory, err := GetRuntimeDirectory(targetUser)
	if err != nil {
		return "", err
	}

	initializedStampBase := fmt.Sprintf("container-initialized-%d", entryPointPID)
	initializedStamp := filepath.Join(toolbxRuntimeDirectory, initializedStampBase)
	return initializedStamp, nil
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

func GetRuntimeDirectory(targetUser *user.User) (string, error) {
	if runtimeDirectories == nil {
		runtimeDirectories = make(map[string]string)
	}

	if toolboxRuntimeDirectory, ok := runtimeDirectories[targetUser.Uid]; ok {
		return toolboxRuntimeDirectory, nil
	}

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
		return "", fmt.Errorf("failed to create runtime directory %s: %w", toolboxRuntimeDirectory, err)
	}

	if err := os.Chown(toolboxRuntimeDirectory, uid, gid); err != nil {
		wrappedErr := fmt.Errorf("failed to change ownership of the runtime directory %s: %w",
			toolboxRuntimeDirectory,
			err)
		return "", wrappedErr
	}

	runtimeDirectories[targetUser.Uid] = toolboxRuntimeDirectory
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
			var errConfigFileNotFound viper.ConfigFileNotFoundError
			var errConfigParse viper.ConfigParseError

			if errors.As(err, &errConfigFileNotFound) || errors.Is(err, os.ErrNotExist) {
				logrus.Debugf("Setting up configuration: file %s not found", configFile)
				continue
			} else if errors.As(err, &errConfigParse) {
				logrus.Debugf("Setting up configuration: failed to parse file %s: %s", configFile, err)
				return fmt.Errorf("failed to parse file %s", configFile)
			} else {
				logrus.Debugf("Setting up configuration: failed to read file %s: %s", configFile, err)
				return fmt.Errorf("failed to read file %s", configFile)
			}
		}
	}

	container, _, _, err := ResolveContainerAndImageNames("", "", "", "")
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

func parseRelease(distro, release string) (string, error) {
	if distro == "" {
		panic("distro not specified")
	}

	distroObj, supportedDistro := supportedDistros[distro]
	if !supportedDistro {
		panicMsg := fmt.Sprintf("failed to find %s in the list of supported distributions", distro)
		panic(panicMsg)
	}

	parseReleaseImpl := distroObj.ParseRelease
	release, err := parseReleaseImpl(release)
	return release, err
}

// PathExists wraps around os.Stat providing a nice interface for checking an existence of a path.
func PathExists(path string) bool {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
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

// ResolveContainerAndImageNames takes care of standardizing names of containers and images.
//
// If no image name is specified then the base image will reflect the platform of the host (even the version).
// If no container name is specified then the name of the image will be used.
//
// If the host system is unknown then the base image will be 'fedora-toolbox' with a default version
func ResolveContainerAndImageNames(container, distroCLI, imageCLI, releaseCLI string) (string, string, string, error) {
	logrus.Debug("Resolving container and image names")
	logrus.Debugf("Container: '%s'", container)
	logrus.Debugf("Distribution (CLI): '%s'", distroCLI)
	logrus.Debugf("Image (CLI): '%s'", imageCLI)
	logrus.Debugf("Release (CLI): '%s'", releaseCLI)

	distro, distroSource := distroCLI, optionValueCLI
	image, release := imageCLI, releaseCLI

	if distroCLI == "" {
		distro, distroSource = distroDefault, optionValueDefault
		if viper.IsSet("general.distro") {
			distro, distroSource = viper.GetString("general.distro"), optionValueConfigFile
		}
	}

	distroObj, supportedDistro := supportedDistros[distro]
	if !supportedDistro {
		return "", "", "", &DistroError{distro, ErrDistroUnsupported}
	}

	if distro == distroDefault {
		if releaseCLI == "" {
			release = releaseDefault
			if viper.IsSet("general.release") {
				release = viper.GetString("general.release")
			}
		}
	} else {
		if distroObj.ReleaseRequired {
			if releaseCLI == "" && !viper.IsSet("general.release") {
				return "", "", "", &DistroError{distro, ErrDistroWithoutRelease}
			}

			if releaseCLI == "" {
				release = viper.GetString("general.release")
			}

			if release == "" {
				panicMsg := fmt.Sprintf("cannot find release for non-default distribution %s", distro)
				panic(panicMsg)
			}
		} else {
			switch distroSource {
			case optionValueCLI:
				// 'release' already set to 'releaseCLI'
			case optionValueConfigFile:
				if releaseCLI == "" {
					if viper.IsSet("general.release") {
						release = viper.GetString("general.release")
					}
				}
			case optionValueDefault:
				panic("distro must be non-default")
			default:
				panic("cannot handle new OptionValueSource")
			}
		}
	}

	release, err := parseRelease(distro, release)
	if err != nil {
		return "", "", "", err
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

	if container == "" {
		var err error
		container, err = getContainerNamePrefixForImage(image)
		if err != nil {
			return "", "", "", err
		}

		tag := ImageReferenceGetTag(image)
		if tag != "" {
			container = container + "-" + tag
		}

		if !IsContainerNameValid(container) {
			return "", "", "", &ContainerError{container, image, ErrContainerNameFromImageInvalid}
		}
	} else {
		if !IsContainerNameValid(container) {
			return "", "", "", &ContainerError{container, "", ErrContainerNameInvalid}
		}
	}

	logrus.Debug("Resolved container and image names")
	logrus.Debugf("Container: '%s'", container)
	logrus.Debugf("Image: '%s'", image)
	logrus.Debugf("Release: '%s'", release)

	return container, image, release, nil
}

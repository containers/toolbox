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
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/acobaugh/osrelease"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/docker/go-units"
	"github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type ParseReleaseFunc func(string) (string, error)

type Distro struct {
	ContainerNamePrefix    string
	ImageBasename          string
	ParseRelease           ParseReleaseFunc
	Registry               string
	Repository             string
	RepositoryNeedsRelease bool
}

const (
	idTruncLength          = 12
	releaseDefaultFallback = "33"
)

const (
	// Based on the nameRegex value in:
	// https://github.com/containers/libpod/blob/master/libpod/options.go
	ContainerNameRegexp = "[a-zA-Z0-9][a-zA-Z0-9_.-]*"
)

var (
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
			"registry.fedoraproject.org",
			"",
			false,
		},
		"rhel": {
			"rhel-toolbox",
			"ubi",
			parseReleaseRHEL,
			"registry.access.redhat.com",
			"ubi8",
			false,
		},
	}
)

var (
	ContainerNameDefault       string
	ContainerNamePrefixDefault = "fedora-toolbox"
)

func init() {
	releaseDefault = releaseDefaultFallback

	hostID, err := GetHostID()
	if err == nil {
		if distroObj, supportedDistro := supportedDistros[hostID]; supportedDistro {
			release, err := GetHostVersionID()
			if err == nil {
				ContainerNamePrefixDefault = distroObj.ContainerNamePrefix
				distroDefault = hostID
				releaseDefault = release
			}
		}
	}

	ContainerNameDefault = ContainerNamePrefixDefault + "-" + releaseDefault
}

func AskForConfirmation(prompt string) bool {
	var retVal bool

	for {
		fmt.Printf("%s ", prompt)

		var response string

		fmt.Scanf("%s", &response)
		if response == "" {
			response = "n"
		} else {
			response = strings.ToLower(response)
		}

		if response == "no" || response == "n" {
			break
		} else if response == "yes" || response == "y" {
			retVal = true
			break
		}
	}

	return retVal
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

func CreateErrorContainerNotFound(container, executableBase string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "container %s not found\n", container)
	fmt.Fprintf(&builder, "Use the 'create' command to create a toolbox.\n")
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func CreateErrorInvalidRelease(executableBase string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "invalid argument for '--release'\n")
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func EnsureXdgRuntimeDirIsSet(uid int) {
	if xdgRuntimeDir, ok := os.LookupEnv("XDG_RUNTIME_DIR"); !ok {
		logrus.Debug("XDG_RUNTIME_DIR is unset")

		xdgRuntimeDir = fmt.Sprintf("/run/user/%d", uid)
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

func GetDefaultImageForDistro(distro, release string) string {
	if _, supportedDistro := supportedDistros[distro]; !supportedDistro {
		distro = "fedora"
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

// HumanDuration accepts a Unix time value and converts it into a human readable
// string.
//
// Examples: "5 minutes ago", "2 hours ago", "3 days ago"
func HumanDuration(duration int64) string {
	return units.HumanDuration(time.Since(time.Unix(duration, 0))) + " ago"
}

// ImageReferenceCanBeID checks if 'image' might be the ID of an image
func ImageReferenceCanBeID(image string) (bool, error) {
	matched, err := regexp.MatchString("^[a-f0-9]\\{6,64\\}$", image)
	return matched, err
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

// ShortID shortens provided id to first 12 characters.
func ShortID(id string) string {
	if len(id) > idTruncLength {
		return id[:idTruncLength]
	}
	return id
}

func ParseRelease(distro, str string) (string, error) {
	if distro == "" {
		distro = distroDefault
	}

	if _, supportedDistro := supportedDistros[distro]; !supportedDistro {
		distro = "fedora"
	}

	distroObj, supportedDistro := supportedDistros[distro]
	if !supportedDistro {
		panicMsg := fmt.Sprintf("failed to find %s in the list of supported distributions", distro)
		panic(panicMsg)
	}

	parseRelease := distroObj.ParseRelease
	release, err := parseRelease(str)
	return release, err
}

func parseReleaseFedora(str string) (string, error) {
	var release string

	if strings.HasPrefix(str, "F") || strings.HasPrefix(str, "f") {
		release = str[1:]
	} else {
		release = str
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

func parseReleaseRHEL(str string) (string, error) {
	if i := strings.IndexRune(str, '.'); i == -1 {
		return "", errors.New("release must have a '.'")
	}

	releaseN, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return "", err
	}

	if releaseN <= 0 {
		return "", errors.New("release must be a positive number")
	}

	return str, nil
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
	if PathExists("/run/.containerenv") {
		return true
	}

	return false
}

func IsInsideToolboxContainer() bool {
	if PathExists("/run/.toolboxenv") {
		return true
	}

	return false
}

func JoinJSON(joinkey string, maps ...[]map[string]interface{}) []map[string]interface{} {
	var json []map[string]interface{}
	found := make(map[string]bool)

	// Iterate over every json provided and check if it is already in the final json
	// If it contains some invalid entry (equals nil), then skip that entry

	for _, m := range maps {
		for _, entry := range m {
			if entry["names"] == nil && entry["Names"] == nil {
				continue
			}
			key := entry[joinkey].(string)
			if _, ok := found[key]; !ok {
				found[key] = true
				json = append(json, entry)
			}
		}
	}
	return json
}

// ResolveContainerAndImageNames takes care of standardizing names of containers and images.
//
// If no image name is specified then the base image will reflect the platform of the host (even the version).
// If no container name is specified then the name of the image will be used.
//
// If the host system is unknown then the base image will be 'fedora-toolbox' with a default version
func ResolveContainerAndImageNames(container, distro, image, release string) (string, string, string, error) {
	logrus.Debug("Resolving container and image names")
	logrus.Debugf("Container: '%s'", container)
	logrus.Debugf("Distribution: '%s'", distro)
	logrus.Debugf("Image: '%s'", image)
	logrus.Debugf("Release: '%s'", release)

	if distro == "" {
		distro = distroDefault
	}

	if distro != distroDefault && release == "" {
		return "", "", "", fmt.Errorf("release not found for non-default distribution %s", distro)
	}

	if release == "" {
		release = releaseDefault
	}

	if image == "" {
		image = GetDefaultImageForDistro(distro, release)
	} else {
		release = ImageReferenceGetTag(image)
		if release == "" {
			release = releaseDefault
		}
	}

	if container == "" {
		var err error
		container, err = GetContainerNamePrefixForImage(image)
		if err != nil {
			return "", "", "", err
		}

		tag := ImageReferenceGetTag(image)
		if tag != "" {
			container = container + "-" + tag
		}
	}

	logrus.Debug("Resolved container and image names")
	logrus.Debugf("Container: '%s'", container)
	logrus.Debugf("Image: '%s'", image)
	logrus.Debugf("Release: '%s'", release)

	return container, image, release, nil
}

func ShowManual(manual string) error {
	manBinary, err := exec.LookPath("man")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("man(1) not found")
		}

		return errors.New("failed to lookup man(1)")
	}

	manualArgs := []string{"man", manual}
	env := os.Environ()

	stderrFd := os.Stderr.Fd()
	stderrFdInt := int(stderrFd)
	stdoutFd := os.Stdout.Fd()
	stdoutFdInt := int(stdoutFd)
	if err := syscall.Dup3(stdoutFdInt, stderrFdInt, 0); err != nil {
		return errors.New("failed to redirect standard error to standard output")
	}

	if err := syscall.Exec(manBinary, manualArgs, env); err != nil {
		return errors.New("failed to invoke man(1)")
	}

	return nil
}

func SortJSON(json []map[string]interface{}, key string, hasInterface bool) []map[string]interface{} {
	sort.Slice(json, func(i, j int) bool {
		if hasInterface {
			return json[i][key].([]interface{})[0].(string) < json[j][key].([]interface{})[0].(string)
		}
		return json[i][key].(string) < json[j][key].(string)
	})

	return json
}

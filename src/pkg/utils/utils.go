/*
 * Copyright © 2019 – 2020 Red Hat Inc.
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

const (
	idTruncLength          = 12
	releaseDefaultFallback = "30"
)

const (
	ContainerNamePrefixDefault = "fedora-toolbox"

	// Based on the nameRegex value in:
	// https://github.com/containers/libpod/blob/master/libpod/options.go
	ContainerNameRegexp = "[a-zA-Z0-9][a-zA-Z0-9_.-]*"
)

var (
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
		"VTE_VERSION",
		"WAYLAND_DISPLAY",
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
)

var (
	ContainerNameDefault string
)

func init() {
	releaseDefault = releaseDefaultFallback

	hostID, err := GetHostID()
	if err == nil {
		if hostID == "fedora" {
			release, err := GetHostVersionID()
			if err == nil {
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
		return "", fmt.Errorf("failed to connect to the D-Bus session instance: %w", err)
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
		return "", fmt.Errorf("unknown reply from org.freedesktop.Flatpak.SessionHelper.RequestSession: %w", err)
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
	if len(options) != 2 {
		return "", errors.New("unexpected output from findmnt(1)")
	}

	mountOptions := strings.TrimSpace(options[0])
	return mountOptions, nil
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

func ParseRelease(str string) (string, error) {
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

// PathExists wraps around os.Stat providing a nice interface for checking an existence of a path.
func PathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}

// IsContainerNameValid checks if the name of a container matches the right pattern
func IsContainerNameValid(containerName string) (bool, error) {
	pattern := "^" + ContainerNameRegexp + "$"
	matched, err := regexp.MatchString(pattern, containerName)
	return matched, err
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
func ResolveContainerAndImageNames(container, image, release string) (string, string, string, error) {
	logrus.Debug("Resolving container and image names")
	logrus.Debugf("Container: '%s'", container)
	logrus.Debugf("Image: '%s'", image)
	logrus.Debugf("Release: '%s'", release)

	if release == "" {
		release = releaseDefault
	}

	if image == "" {
		image = "fedora-toolbox:" + release
	} else {
		release = ImageReferenceGetTag(image)
		if release == "" {
			release = releaseDefault
		}
	}

	if container == "" {
		basename := ImageReferenceGetBasename(image)
		if basename == "" {
			return "", "", "", fmt.Errorf("failed to get the basename of image %s", image)
		}

		container = basename

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
		return fmt.Errorf("failed to redirect standard error to standard output: %w", err)
	}

	if err := syscall.Exec(manBinary, manualArgs, env); err != nil {
		return fmt.Errorf("failed to invoke man(1): %w", err)
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

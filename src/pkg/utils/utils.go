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
	"regexp"
	"sort"
	"syscall"

	"github.com/containers/toolbox/pkg/shell"
	"github.com/godbus/dbus/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

const (
	idTruncLength = 12
)

const (
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
)

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

// ShortID shortens provided id to first 12 characters.
func ShortID(id string) string {
	if len(id) > idTruncLength {
		return id[:idTruncLength]
	}
	return id
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
	if err := syscall.Dup2(stdoutFdInt, stderrFdInt); err != nil {
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

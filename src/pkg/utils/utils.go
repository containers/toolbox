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
	"syscall"

	"github.com/containers/toolbox/pkg/shell"
	"github.com/sirupsen/logrus"
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

// PathExists wraps around os.Stat providing a nice interface for checking an existence of a path.
func PathExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
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

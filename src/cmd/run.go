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

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	runFlags struct {
		container string
		distro    string
		release   string
	}
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command in an existing toolbox container",
	RunE:  run,
}

func init() {
	flags := runCmd.Flags()
	flags.SetInterspersed(false)

	flags.StringVarP(&runFlags.container,
		"container",
		"c",
		"",
		"Run command inside a toolbox container with the given name")

	flags.StringVarP(&runFlags.distro,
		"distro",
		"d",
		"",
		"Run command inside a toolbox container for a different operating system distribution than the host")

	flags.StringVarP(&runFlags.release,
		"release",
		"r",
		"",
		"Run command inside a toolbox container for a different operating system release than the host")

	runCmd.SetHelpFunc(runHelp)
	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a toolbox container")
		}

		if _, err := utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	var nonDefaultContainer bool

	if runFlags.container != "" {
		nonDefaultContainer = true

		if !utils.IsContainerNameValid(runFlags.container) {
			var builder strings.Builder
			fmt.Fprintf(&builder, "invalid argument for '--container'\n")
			fmt.Fprintf(&builder, "Container names must match '%s'\n", utils.ContainerNameRegexp)
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	}

	var release string
	if runFlags.release != "" {
		nonDefaultContainer = true

		var err error
		release, err = utils.ParseRelease(runFlags.distro, runFlags.release)
		if err != nil {
			err := utils.CreateErrorInvalidRelease(executableBase)
			return err
		}
	}

	if len(args) == 0 {
		var builder strings.Builder
		fmt.Fprintf(&builder, "missing argument for \"run\"\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	command := args

	container, image, release, err := utils.ResolveContainerAndImageNames(runFlags.container, runFlags.distro, "", release)
	if err != nil {
		return err
	}

	if err := runCommand(container,
		!nonDefaultContainer,
		image,
		release,
		command,
		false,
		false,
		true); err != nil {
		return err
	}

	return nil
}

func runCommand(container string,
	defaultContainer bool,
	image, release string,
	command []string,
	emitEscapeSequence, fallbackToBash, pedantic bool) error {
	if !pedantic {
		if image == "" {
			panic("image not specified")
		}

		if release == "" {
			panic("release not specified")
		}
	}

	logrus.Debugf("Checking if container %s exists", container)

	if _, err := podman.ContainerExists(container); err != nil {
		logrus.Debugf("Container %s not found", container)

		if pedantic {
			err := utils.CreateErrorContainerNotFound(container, executableBase)
			return err
		}

		containers, err := getContainers()
		if err != nil {
			err := utils.CreateErrorContainerNotFound(container, executableBase)
			return err
		}

		containersCount := len(containers)
		logrus.Debugf("Found %d containers", containersCount)

		if containersCount == 0 {
			var shouldCreateContainer bool
			promptForCreate := true

			if rootFlags.assumeYes {
				shouldCreateContainer = true
				promptForCreate = false
			}

			if promptForCreate {
				prompt := "No toolbox containers found. Create now? [y/N]"
				shouldCreateContainer = utils.AskForConfirmation(prompt)
			}

			if !shouldCreateContainer {
				fmt.Printf("A container can be created later with the 'create' command.\n")
				fmt.Printf("Run '%s --help' for usage.\n", executableBase)
				return nil
			}

			if err := createContainer(container, image, release, false); err != nil {
				return err
			}
		} else if containersCount == 1 && defaultContainer {
			fmt.Fprintf(os.Stderr, "Error: container %s not found\n", container)

			container = containers[0].Names[0]
			fmt.Fprintf(os.Stderr, "Entering container %s instead.\n", container)
			fmt.Fprintf(os.Stderr, "Use the 'create' command to create a different toolbox.\n")
			fmt.Fprintf(os.Stderr, "Run '%s --help' for usage.\n", executableBase)
		} else {
			var builder strings.Builder
			fmt.Fprintf(&builder, "container %s not found\n", container)
			fmt.Fprintf(&builder, "Use the '--container' option to select a toolbox.\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	}

	if err := callFlatpakSessionHelper(container); err != nil {
		return err
	}

	logrus.Debugf("Starting container %s", container)
	if err := startContainer(container); err != nil {
		return err
	}

	entryPoint, entryPointPID, err := getEntryPointAndPID(container)
	if err != nil {
		return err
	}

	if entryPoint != "toolbox" {
		var builder strings.Builder
		fmt.Fprintf(&builder, "container %s is too old and no longer supported \n", container)
		fmt.Fprintf(&builder, "Recreate it with Toolbox version 0.0.17 or newer.\n")

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if entryPointPID <= 0 {
		return fmt.Errorf("invalid entry point PID of container %s", container)
	}

	logrus.Debugf("Waiting for container %s to finish initializing", container)

	toolboxRuntimeDirectory, err := utils.GetRuntimeDirectory(currentUser)
	if err != nil {
		return err
	}

	initializedStamp := fmt.Sprintf("%s/container-initialized-%d", toolboxRuntimeDirectory, entryPointPID)

	logrus.Debugf("Checking if initialization stamp %s exists", initializedStamp)

	initializedTimeout := 25 // seconds
	for i := 0; !utils.PathExists(initializedStamp); i++ {
		if i == initializedTimeout {
			return fmt.Errorf("failed to initialize container %s", container)
		}

		time.Sleep(time.Second)
	}

	logrus.Debugf("Container %s is initialized", container)

	if _, err := isCommandPresent(container, command[0]); err != nil {
		if fallbackToBash {
			fmt.Fprintf(os.Stderr,
				"Error: command %s not found in container %s\n",
				command[0],
				container)
			fmt.Fprintf(os.Stderr, "Using /bin/bash instead.\n")

			command = []string{"/bin/bash", "-l"}
		} else {
			return fmt.Errorf("command %s not found in container %s", command[0], container)
		}
	}

	if pathPresent, _ := isPathPresent(container, workingDirectory); !pathPresent {
		fmt.Fprintf(os.Stderr, "Error: path %s not found in container %s\n",
			workingDirectory, container)
		fmt.Fprintf(os.Stderr, "Using %s instead.\n",
			currentUser.HomeDir)

		workingDirectory = currentUser.HomeDir
	}

	logrus.Debug("Checking if 'podman exec' supports disabling the detach keys")

	var detachKeys []string

	if podman.CheckVersion("1.8.1") {
		logrus.Debug("'podman exec' supports disabling the detach keys")
		detachKeys = []string{"--detach-keys", ""}
	}

	envOptions := utils.GetEnvOptionsForPreservedVariables()
	logLevelString := podman.LogLevel.String()

	execArgs := []string{
		"--log-level", logLevelString,
		"exec",
	}

	execArgs = append(execArgs, detachKeys...)

	execArgs = append(execArgs, []string{
		"--interactive",
		"--tty",
		"--user", currentUser.Username,
		"--workdir", workingDirectory,
	}...)

	execArgs = append(execArgs, envOptions...)

	execArgs = append(execArgs, []string{
		container,
		"capsh", "--caps=", "--", "-c", "exec \"$@\"", "/bin/sh",
	}...)

	execArgs = append(execArgs, command...)

	if emitEscapeSequence {
		fmt.Printf("\033]777;container;push;%s;toolbox;%s\033\\", container, currentUser.Uid)
	}

	logrus.Debugf("Running in container %s:", container)
	logrus.Debug("podman")
	for _, arg := range execArgs {
		logrus.Debugf("%s", arg)
	}

	exitCode, err := shell.RunWithExitCode("podman", os.Stdin, os.Stdout, nil, execArgs...)

	if emitEscapeSequence {
		fmt.Printf("\033]777;container;pop;;;%s\033\\", currentUser.Uid)
	}

	switch exitCode {
	case 0:
		if err != nil {
			panic("unexpected error: 'podman exec' finished successfully")
		}
	case 125:
		err = fmt.Errorf("failed to invoke 'podman exec' in container %s", container)
	case 126:
		err = fmt.Errorf("failed to invoke command %s in container %s", command[0], container)
	case 127:
		if pathPresent, _ := isPathPresent(container, workingDirectory); !pathPresent {
			err = fmt.Errorf("directory %s not found in container %s", workingDirectory, container)
		} else {
			err = fmt.Errorf("command %s not found in container %s", command[0], container)
		}
	default:
		err = nil
	}

	if err != nil {
		return err
	}

	return nil
}

func runHelp(cmd *cobra.Command, args []string) {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			fmt.Fprintf(os.Stderr, "Error: this is not a toolbox container\n")
			return
		}

		if _, err := utils.ForwardToHost(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return
		}

		return
	}

	if err := utils.ShowManual("toolbox-run"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func callFlatpakSessionHelper(container string) error {
	logrus.Debugf("Inspecting mounts of container %s", container)

	info, err := podman.Inspect("container", container)
	if err != nil {
		return fmt.Errorf("failed to inspect entry point of container %s", container)
	}

	var needsFlatpakSessionHelper bool

	mounts := info["Mounts"].([]interface{})
	for _, mount := range mounts {
		destination := mount.(map[string]interface{})["Destination"].(string)
		if destination == "/run/host/monitor" {
			logrus.Debug("Requires org.freedesktop.Flatpak.SessionHelper")
			needsFlatpakSessionHelper = true
			break
		}
	}

	if !needsFlatpakSessionHelper {
		return nil
	}

	if _, err := utils.CallFlatpakSessionHelper(); err != nil {
		return err
	}

	return nil
}

func getEntryPointAndPID(container string) (string, int, error) {
	logrus.Debugf("Inspecting entry point of container %s", container)

	info, err := podman.Inspect("container", container)
	if err != nil {
		return "", 0, fmt.Errorf("failed to inspect entry point of container %s", container)
	}

	config := info["Config"].(map[string]interface{})
	entryPoint := config["Cmd"].([]interface{})[0].(string)

	state := info["State"].(map[string]interface{})
	entryPointPID := state["Pid"]
	logrus.Debugf("Entry point PID is a %T", entryPointPID)

	var entryPointPIDInt int

	switch entryPointPID.(type) {
	case float64:
		entryPointPIDFloat := entryPointPID.(float64)
		entryPointPIDInt = int(entryPointPIDFloat)
	case int:
		entryPointPIDInt = entryPointPID.(int)
	default:
		return "", 0, fmt.Errorf("failed to inspect entry point PID of container %s", container)
	}

	logrus.Debugf("Entry point of container %s is %s (PID=%d)", container, entryPoint, entryPointPIDInt)

	return entryPoint, entryPointPIDInt, nil
}

func isCommandPresent(container, command string) (bool, error) {
	logrus.Debugf("Looking for command %s in container %s", command, container)

	logLevelString := podman.LogLevel.String()
	args := []string{
		"--log-level", logLevelString,
		"exec",
		"--user", currentUser.Username,
		container,
		"sh", "-c", "command -v \"$1\"", "sh", command,
	}

	if err := shell.Run("podman", nil, nil, nil, args...); err != nil {
		return false, err
	}

	return true, nil
}

func isPathPresent(container, path string) (bool, error) {
	logrus.Debugf("Looking for path %s in container %s", path, container)

	logLevelString := podman.LogLevel.String()
	args := []string{
		"--log-level", logLevelString,
		"exec",
		"--user", currentUser.Username,
		container,
		"sh", "-c", "test -d \"$1\"", "sh", path,
	}

	if err := shell.Run("podman", nil, nil, nil, args...); err != nil {
		return false, err
	}

	return true, nil
}

func startContainer(container string) error {
	var stderr strings.Builder
	if err := podman.Start(container, &stderr); err == nil {
		return nil
	}

	errString := stderr.String()
	if !strings.Contains(errString, "use system migrate to mitigate") {
		return fmt.Errorf("failed to start container %s", container)
	}

	logrus.Debug("Checking if 'podman system migrate' supports '--new-runtime'")

	if !podman.CheckVersion("1.6.2") {
		var builder strings.Builder

		fmt.Fprintf(&builder,
			"container %s doesn't support cgroups v%d\n",
			container,
			cgroupsVersion)

		fmt.Fprintf(&builder, "Update Podman to version 1.6.2 or newer.\n")

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	logrus.Debug("'podman system migrate' supports '--new-runtime'")

	ociRuntimeRequired := "runc"
	if cgroupsVersion == 2 {
		ociRuntimeRequired = "crun"
	}

	logrus.Debugf("Migrating containers to OCI runtime %s", ociRuntimeRequired)

	if err := podman.SystemMigrate(ociRuntimeRequired); err != nil {
		var builder strings.Builder
		fmt.Fprintf(&builder, "failed to migrate containers to OCI runtime %s\n", ociRuntimeRequired)
		fmt.Fprintf(&builder, "Factory reset with: podman system reset")

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if err := podman.Start(container, nil); err != nil {
		var builder strings.Builder
		fmt.Fprintf(&builder, "container %s doesn't support cgroups v%d\n", container, cgroupsVersion)
		fmt.Fprintf(&builder, "Factory reset with: podman system reset")

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	return nil
}

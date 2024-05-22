/*
 * Copyright © 2019 – 2024 Red Hat Inc.
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
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/term"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	runFlags struct {
		container   string
		distro      string
		preserveFDs uint
		release     string
	}

	runFallbackCommands = [][]string{{"/bin/bash", "-l"}}
	runFallbackWorkDirs = []string{"" /* $HOME */}
)

var runCmd = &cobra.Command{
	Use:               "run",
	Short:             "Run a command in an existing Toolbx container",
	RunE:              run,
	ValidArgsFunction: completionEmpty,
}

func init() {
	flags := runCmd.Flags()
	flags.SetInterspersed(false)

	flags.StringVarP(&runFlags.container,
		"container",
		"c",
		"",
		"Run command inside a Toolbx container with the given name")

	flags.StringVarP(&runFlags.distro,
		"distro",
		"d",
		"",
		"Run command inside a Toolbx container for a different operating system distribution than the host")

	flags.UintVar(&runFlags.preserveFDs,
		"preserve-fds",
		0,
		"Pass down to command N additional file descriptors (in addition to 0, 1, 2)")

	flags.StringVarP(&runFlags.release,
		"release",
		"r",
		"",
		"Run command inside a Toolbx container for a different operating system release than the host")

	runCmd.SetHelpFunc(runHelp)

	if err := runCmd.RegisterFlagCompletionFunc("container", completionContainerNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}
	if err := runCmd.RegisterFlagCompletionFunc("distro", completionDistroNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}

	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		if _, err := utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	var defaultContainer bool = true

	if runFlags.container != "" {
		defaultContainer = false
	}

	if runFlags.release != "" {
		defaultContainer = false
	}

	if len(args) == 0 {
		var builder strings.Builder
		fmt.Fprintf(&builder, "missing argument for \"run\"\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	command := args

	container, image, release, err := resolveContainerAndImageNames(runFlags.container,
		"--container",
		runFlags.distro,
		"",
		runFlags.release)

	if err != nil {
		return err
	}

	if err := runCommand(container,
		defaultContainer,
		image,
		release,
		runFlags.preserveFDs,
		command,
		false,
		false,
		true); err != nil {
		// runCommand returns exitError for the executed commands to properly
		// propagate return codes. Cobra prints all non-nil errors which in
		// that case is not desirable. In that scenario silence the errors and
		// leave the error handling to the root command.
		var errExit *exitError
		if errors.As(err, &errExit) {
			cmd.SilenceErrors = true
		}

		return err
	}

	return nil
}

func runCommand(container string,
	defaultContainer bool,
	image, release string,
	preserveFDs uint,
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
			err := createErrorContainerNotFound(container)
			return err
		}

		containers, err := getContainers()
		if err != nil {
			err := createErrorContainerNotFound(container)
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
				prompt := "No Toolbx containers found. Create now? [y/N]"
				shouldCreateContainer = askForConfirmation(prompt)
			}

			if !shouldCreateContainer {
				fmt.Printf("A container can be created later with the 'create' command.\n")
				fmt.Printf("Run '%s --help' for usage.\n", executableBase)
				return nil
			}

			if err := createContainer(container, image, release, "", false); err != nil {
				return err
			}
		} else if containersCount == 1 && defaultContainer {
			fmt.Fprintf(os.Stderr, "Error: container %s not found\n", container)

			container = containers[0].Name()
			fmt.Fprintf(os.Stderr, "Entering container %s instead.\n", container)
			fmt.Fprintf(os.Stderr, "Use the 'create' command to create a different Toolbx.\n")
			fmt.Fprintf(os.Stderr, "Run '%s --help' for usage.\n", executableBase)
		} else {
			var builder strings.Builder
			fmt.Fprintf(&builder, "container %s not found\n", container)
			fmt.Fprintf(&builder, "Use the '--container' option to select a Toolbx.\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	}

	logrus.Debugf("Inspecting container %s", container)
	containerObj, err := podman.InspectContainer(container)
	if err != nil {
		return fmt.Errorf("failed to inspect container %s", container)
	}

	entryPoint := containerObj.EntryPoint()
	entryPointPID := containerObj.EntryPointPID()
	logrus.Debugf("Entry point of container %s is %s (PID=%d)", container, entryPoint, entryPointPID)

	if entryPoint != "toolbox" {
		var builder strings.Builder
		fmt.Fprintf(&builder, "container %s is too old and no longer supported\n", container)
		fmt.Fprintf(&builder, "Recreate it with Toolbx version 0.0.17 or newer.")

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if err := callFlatpakSessionHelper(containerObj); err != nil {
		return err
	}

	if entryPointPID <= 0 {
		logrus.Debugf("Starting container %s", container)
		if err := startContainer(container); err != nil {
			return err
		}

		logrus.Debugf("Inspecting container %s", container)
		containerObj, err := podman.InspectContainer(container)
		if err != nil {
			return fmt.Errorf("failed to inspect container %s", container)
		}

		entryPointPID = containerObj.EntryPointPID()
		logrus.Debugf("Entry point of container %s is %s (PID=%d)", container, entryPoint, entryPointPID)

		if entryPointPID <= 0 {
			return fmt.Errorf("invalid entry point PID of container %s", container)
		}

		logrus.Debugf("Waiting for container %s to finish initializing", container)
	}

	toolboxRuntimeDirectory, err := utils.GetRuntimeDirectory(currentUser)
	if err != nil {
		return err
	}

	initializedStampBase := fmt.Sprintf("container-initialized-%d", entryPointPID)
	initializedStamp := filepath.Join(toolboxRuntimeDirectory, initializedStampBase)

	logrus.Debugf("Checking if initialization stamp %s exists", initializedStamp)

	initializedTimeout := 25 // seconds
	for i := 0; !utils.PathExists(initializedStamp); i++ {
		if i == initializedTimeout {
			return fmt.Errorf("failed to initialize container %s", container)
		}

		time.Sleep(time.Second)
	}

	logrus.Debugf("Container %s is initialized", container)

	if err := runCommandWithFallbacks(container,
		preserveFDs,
		command,
		emitEscapeSequence,
		fallbackToBash); err != nil {
		return err
	}

	return nil
}

func runCommandWithFallbacks(container string,
	preserveFDs uint,
	command []string,
	emitEscapeSequence, fallbackToBash bool) error {

	logrus.Debug("Checking if 'podman exec' supports disabling the detach keys")

	var detachKeysSupported bool

	if podman.CheckVersion("1.8.1") {
		logrus.Debug("'podman exec' supports disabling the detach keys")
		detachKeysSupported = true
	}

	envOptions := utils.GetEnvOptionsForPreservedVariables()
	preserveFDsString := fmt.Sprint(preserveFDs)

	var stderr io.Writer
	var ttyNeeded bool

	if term.IsTerminal(os.Stdin) && term.IsTerminal(os.Stdout) {
		ttyNeeded = true
		if logLevel := logrus.GetLevel(); logLevel >= logrus.DebugLevel {
			stderr = os.Stderr
		}
	} else {
		stderr = os.Stderr
	}

	runFallbackCommandsIndex := 0
	runFallbackWorkDirsIndex := 0
	workDir := workingDirectory

	for {
		execArgs := constructExecArgs(container,
			preserveFDsString,
			command,
			detachKeysSupported,
			envOptions,
			fallbackToBash,
			ttyNeeded,
			workDir)

		if emitEscapeSequence {
			fmt.Printf("\033]777;container;push;%s;toolbox;%s\033\\", container, currentUser.Uid)
		}

		logrus.Debugf("Running in container %s:", container)
		logrus.Debug("podman")
		for _, arg := range execArgs {
			logrus.Debugf("%s", arg)
		}

		exitCode, err := shell.RunWithExitCode("podman", os.Stdin, os.Stdout, stderr, execArgs...)

		if emitEscapeSequence {
			fmt.Printf("\033]777;container;pop;;;%s\033\\", currentUser.Uid)
		}

		switch exitCode {
		case 0:
			if err != nil {
				panic("unexpected error: 'podman exec' finished successfully")
			}
			return nil
		case 125:
			return &exitError{exitCode, fmt.Errorf("failed to invoke 'podman exec' in container %s", container)}
		case 126:
			return &exitError{exitCode, fmt.Errorf("failed to invoke command %s in container %s", command[0], container)}
		case 127:
			if pathPresent, _ := isPathPresent(container, workDir); !pathPresent {
				if runFallbackWorkDirsIndex < len(runFallbackWorkDirs) {
					fmt.Fprintf(os.Stderr,
						"Error: directory %s not found in container %s\n",
						workDir,
						container)

					workDir = runFallbackWorkDirs[runFallbackWorkDirsIndex]
					if workDir == "" {
						workDir = currentUser.HomeDir
					}

					fmt.Fprintf(os.Stderr, "Using %s instead.\n", workDir)
					runFallbackWorkDirsIndex++
				} else {
					return &exitError{exitCode, fmt.Errorf("directory %s not found in container %s", workDir, container)}
				}
			} else if _, err := isCommandPresent(container, command[0]); err != nil {
				if fallbackToBash && runFallbackCommandsIndex < len(runFallbackCommands) {
					fmt.Fprintf(os.Stderr,
						"Error: command %s not found in container %s\n",
						command[0],
						container)

					command = runFallbackCommands[runFallbackCommandsIndex]
					fmt.Fprintf(os.Stderr, "Using %s instead.\n", command[0])

					runFallbackCommandsIndex++
				} else {
					return &exitError{exitCode, fmt.Errorf("command %s not found in container %s", command[0], container)}
				}
			} else {
				return nil
			}
		default:
			return &exitError{exitCode, nil}
		}
	}

	// code should not be reached
}

func runHelp(cmd *cobra.Command, args []string) {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			fmt.Fprintf(os.Stderr, "Error: this is not a Toolbx container\n")
			return
		}

		if _, err := utils.ForwardToHost(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return
		}

		return
	}

	if err := showManual("toolbox-run"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func callFlatpakSessionHelper(container podman.Container) error {
	name := container.Name()
	logrus.Debugf("Inspecting mounts of container %s", name)

	var needsFlatpakSessionHelper bool

	mounts := container.Mounts()
	for _, mount := range mounts {
		if mount == "/run/host/monitor" {
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

func constructCapShArgs(command []string, useLoginShell bool) []string {
	capShArgs := []string{"capsh", "--caps=", "--"}

	if useLoginShell {
		capShArgs = append(capShArgs, []string{"--login"}...)
	}

	capShArgs = append(capShArgs, []string{"-c", "exec \"$@\"", "bash"}...)
	capShArgs = append(capShArgs, command...)

	return capShArgs
}

func constructExecArgs(container, preserveFDs string,
	command []string,
	detachKeysSupported bool,
	envOptions []string,
	fallbackToBash bool,
	ttyNeeded bool,
	workDir string) []string {

	logLevelString := podman.LogLevel.String()

	execArgs := []string{
		"--log-level", logLevelString,
		"exec",
	}

	if detachKeysSupported {
		execArgs = append(execArgs, []string{
			"--detach-keys", "",
		}...)
	}

	execArgs = append(execArgs, envOptions...)

	execArgs = append(execArgs, []string{
		"--interactive",
		"--preserve-fds", preserveFDs,
	}...)

	if ttyNeeded {
		execArgs = append(execArgs, []string{
			"--tty",
		}...)
	}

	execArgs = append(execArgs, []string{
		"--user", currentUser.Username,
		"--workdir", workDir,
	}...)

	execArgs = append(execArgs, []string{
		container,
	}...)

	capShArgs := constructCapShArgs(command, !fallbackToBash)
	execArgs = append(execArgs, capShArgs...)

	return execArgs
}

func isCommandPresent(container, command string) (bool, error) {
	logrus.Debugf("Looking up command %s in container %s", command, container)

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
	logrus.Debugf("Looking up path %s in container %s", path, container)

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

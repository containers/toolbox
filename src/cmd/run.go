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

package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/containers/toolbox/pkg/nvidia"
	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/term"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/go-logfmt/logfmt"
	"github.com/google/renameio/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"tags.cncf.io/container-device-interface/specs-go"
)

type collectEntryPointErrorFunc func(err error)

type entryPointError struct {
	msg string
}

var (
	runFlags struct {
		container   string
		distro      string
		preserveFDs uint
		release     string
		workDir     string
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

	flags.StringVarP(&runFlags.workDir,
		"workdir",
		"w",
		"",
		"Run command inside a Toolbx container within the given working directory.")

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

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
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
		runFlags.workDir,
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
	preserveFDs uint,
	workDir string,
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
		fmt.Fprintf(&builder, "Recreate it with Toolbx version 0.0.97 or newer.")

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if err := callFlatpakSessionHelper(containerObj); err != nil {
		return err
	}

	var cdiEnviron []string

	cdiSpecForNvidia, err := nvidia.GenerateCDISpec()
	if err != nil {
		if errors.Is(err, nvidia.ErrNVMLDriverLibraryVersionMismatch) {
			var builder strings.Builder
			builder.WriteString("the proprietary NVIDIA driver's kernel and user space don't match\n")
			builder.WriteString("Check the host operating system and systemd journal.")

			errMsg := builder.String()
			return errors.New(errMsg)
		} else if !errors.Is(err, nvidia.ErrPlatformUnsupported) {
			return err
		}
	} else {
		cdiEnviron = append(cdiEnviron, cdiSpecForNvidia.ContainerEdits.Env...)
	}

	p11KitServerEnviron, err := startP11KitServer()
	if err != nil {
		return err
	}

	startContainerTimestamp := time.Unix(-1, 0)

	if entryPointPID <= 0 {
		if cdiSpecForNvidia != nil {
			cdiFileForNvidia, err := getCDIFileForNvidia(currentUser)
			if err != nil {
				return err
			}

			if err := saveCDISpecTo(cdiSpecForNvidia, cdiFileForNvidia); err != nil {
				return err
			}
		}

		startContainerTimestamp = time.Now()

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
			if err := showEntryPointLogs(container, startContainerTimestamp); err != nil {
				var errEntryPoint *entryPointError
				if errors.As(err, &errEntryPoint) {
					return err
				}

				logrus.Debugf("Reading logs from container %s failed: %s", container, err)
			}

			return fmt.Errorf("invalid entry point PID of container %s", container)
		}

		logrus.Debugf("Waiting for container %s to finish initializing", container)
	}

	if err := ensureContainerIsInitialized(container, entryPointPID, startContainerTimestamp); err != nil {
		return err
	}

	logrus.Debugf("Container %s is initialized", container)

	environ := append(cdiEnviron, p11KitServerEnviron...)
	if err := runCommandWithFallbacks(container,
		preserveFDs,
		workDir,
		command,
		environ,
		emitEscapeSequence,
		fallbackToBash); err != nil {
		return err
	}

	return nil
}

func runCommandWithFallbacks(container string,
	preserveFDs uint,
	workDir string,
	command, environ []string,
	emitEscapeSequence, fallbackToBash bool) error {

	logrus.Debug("Checking if 'podman exec' supports disabling the detach keys")

	var detachKeysSupported bool

	if podman.CheckVersion("1.8.1") {
		logrus.Debug("'podman exec' supports disabling the detach keys")
		detachKeysSupported = true
	}

	envOptions := utils.GetEnvOptionsForPreservedVariables()
	for _, env := range environ {
		logrus.Debugf("%s", env)
		envOption := "--env=" + env
		envOptions = append(envOptions, envOption)
	}

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

	if workDir == "" {
		workDir = workingDirectory
	}

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
			errMsg := fmt.Sprintf("failed to invoke 'podman exec' in container %s", container)
			return &exitError{exitCode, errors.New(errMsg)}
		case 126:
			var err error
			if command[0] != "toolbox" {
				errMsg := fmt.Sprintf("failed to invoke command %s in container %s",
					command[0],
					container)
				err = errors.New(errMsg)
			}

			return &exitError{exitCode, err}
		case 127:
			if pathPresent, _ := isPathPresent(container, workDir); !pathPresent {
				if runFallbackWorkDirsIndex < len(runFallbackWorkDirs) {
					fmt.Fprintf(os.Stderr,
						"Error: directory %s not found in container %s\n",
						workDir,
						container)

					workDir = runFallbackWorkDirs[runFallbackWorkDirsIndex]
					if workDir == "" {
						workDir = getCurrentUserHomeDir()
					}

					fmt.Fprintf(os.Stderr, "Using %s instead.\n", workDir)
					runFallbackWorkDirsIndex++
				} else {
					errMsg := fmt.Sprintf("directory %s not found in container %s",
						workDir,
						container)
					return &exitError{exitCode, errors.New(errMsg)}
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
					errMsg := fmt.Sprintf("command %s not found in container %s",
						command[0],
						container)
					return &exitError{exitCode, errors.New(errMsg)}
				}
			} else if command[0] == "toolbox" {
				return &exitError{exitCode, nil}
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

	fmt.Fprintf(os.Stderr, "Warning: container %s uses deprecated features\n", name)
	fmt.Fprintf(os.Stderr, "Consider recreating it with Toolbx version 0.0.97 or newer.\n")

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

func ensureContainerIsInitialized(container string, entryPointPID int, timestamp time.Time) error {
	initializedStamp, err := utils.GetInitializedStamp(entryPointPID, currentUser)
	if err != nil {
		return err
	}

	logrus.Debugf("Checking if initialization stamp %s exists", initializedStamp)

	shouldUsePolling := isUsePollingSet()

	// Optional shortcut for containers that are already initialized
	if !shouldUsePolling && utils.PathExists(initializedStamp) {
		return nil
	}

	logrus.Debugf("Setting up initialization timeout for container %s", container)
	initializedTimeout := time.NewTimer(25 * time.Second)
	defer initializedTimeout.Stop()

	logrus.Debugf("Following logs for container %s", container)

	parentCtx := context.Background()
	logsCtx, logsCancel := context.WithCancelCause(parentCtx)
	defer logsCancel(errors.New("clean-up"))

	logsCh, logsErrCh := followEntryPointLogsAsync(logsCtx, container, timestamp)

	fallbackToPolling := shouldUsePolling
	var watcherForStamp *fsnotify.Watcher

	if !fallbackToPolling {
		logrus.Debugf("Setting up watches for file system events from container %s", container)

		var err error
		watcherForStamp, err = fsnotify.NewWatcher()
		if err != nil {
			if errors.Is(err, unix.EMFILE) ||
				errors.Is(err, unix.ENFILE) ||
				errors.Is(err, unix.ENOMEM) {
				logrus.Debugf("Setting up watches for file system events: failed to create Watcher: %s",
					err)
				logrus.Debug("Using polling instead")
				fallbackToPolling = true
			} else {
				return fmt.Errorf("failed to create Watcher: %w", err)
			}
		}
	}

	var watcherForStampErrors chan error
	var watcherForStampEvents chan fsnotify.Event

	if watcherForStamp != nil {
		defer watcherForStamp.Close()

		toolboxRuntimeDirectory, err := utils.GetRuntimeDirectory(currentUser)
		if err != nil {
			return err
		}

		if err := watcherForStamp.Add(toolboxRuntimeDirectory); err != nil {
			if errors.Is(err, unix.ENOMEM) || errors.Is(err, unix.ENOSPC) {
				logrus.Debugf("Setting up watches for file system events: failed to add path: %s",
					err)
				logrus.Debug("Using polling instead")
				fallbackToPolling = true
			} else {
				return fmt.Errorf("failed to add path to Watcher: %w", err)
			}
		}

		if !fallbackToPolling {
			watcherForStampErrors = watcherForStamp.Errors
			watcherForStampEvents = watcherForStamp.Events
		}
	}

	var tickerPolling *time.Ticker
	var tickerPollingCh <-chan time.Time

	if fallbackToPolling {
		logrus.Debugf("Setting up polling ticker for container %s", container)

		tickerPolling = time.NewTicker(time.Second)
		defer tickerPolling.Stop()

		tickerPollingCh = tickerPolling.C
	}

	// Initialization could have finished before the Watcher was set up
	if !shouldUsePolling && utils.PathExists(initializedStamp) {
		return nil
	}

	logrus.Debug("Listening to container, ticker and timeout events")

	var errReceivedFromEntryPoint error

	for {
		select {
		case <-initializedTimeout.C:
			logsCancel(context.DeadlineExceeded)
			if utils.PathExists(initializedStamp) {
				return nil
			} else {
				return fmt.Errorf("failed to initialize container %s", container)
			}
		case line, ok := <-logsCh:
			collectEntryPointErrorFn := func(err error) {
				if !errors.Is(errReceivedFromEntryPoint, err) {
					errReceivedFromEntryPoint = errors.Join(errReceivedFromEntryPoint, err)
				}
			}

			if err := handleEntryPointLog(logsCtx,
				container,
				!ok,
				line,
				collectEntryPointErrorFn); err != nil {
				if errors.Is(err, context.Canceled) {
					return nil
				} else if errReceivedFromEntryPoint != nil {
					return errReceivedFromEntryPoint
				} else {
					logsCh = nil
				}
			}
		case err := <-logsErrCh:
			logrus.Debugf("Received an error while following the logs: %s", err)
		case time := <-tickerPollingCh:
			if done := handlePollingTickForStamp(time, initializedStamp); done {
				cause := fmt.Errorf("%w: initialization stamp %s exists",
					context.Canceled,
					initializedStamp)
				logsCancel(cause)
			}
		case event := <-watcherForStampEvents:
			if done := handleFileSystemEventForStamp(event, initializedStamp); done {
				cause := fmt.Errorf("%w: initialization stamp %s exists",
					context.Canceled,
					initializedStamp)
				logsCancel(cause)
			}
		case err := <-watcherForStampErrors:
			logrus.Debugf("Received an error from the file system watcher: %s", err)
		}
	}

	// code should not be reached
}

func followEntryPointLogsAsync(ctx context.Context, container string, since time.Time) (
	<-chan string, <-chan error,
) {
	reader, writer := io.Pipe()
	retValCh := make(chan string)
	errCh := make(chan error)

	go func() {
		defer writer.Close()

		if err := podman.LogsContext(ctx, container, true, since, writer); err != nil {
			errCh <- err
			return
		}
	}()

	go func() {
		defer reader.Close()
		defer close(retValCh)

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			retValCh <- line
		}

		if err := scanner.Err(); err != nil {
			errCh <- err
		} else {
			errCh <- io.EOF
		}
	}()

	return retValCh, errCh
}

func handleEntryPointLog(ctx context.Context,
	container string,
	end bool,
	line string,
	collectEntryPointErrorFn collectEntryPointErrorFunc) error {
	if end {
		if cause := context.Cause(ctx); errors.Is(cause, context.Canceled) {
			return cause
		} else {
			logrus.Debugf("Reading logs from container %s failed: 'podman logs' finished unexpectedly",
				container)
			return errors.New("'podman logs' finished unexpectedly")
		}
	} else {
		if err := showEntryPointLog(line); err != nil {
			var errEntryPoint *entryPointError

			if errors.As(err, &errEntryPoint) {
				collectEntryPointErrorFn(err)
				return nil
			}

			logrus.Debugf("Parsing entry point log failed: %s:", err)
			logrus.Debugf("%s", line)
		}
	}

	return nil
}

func handleFileSystemEventForStamp(event fsnotify.Event, initializedStamp string) bool {
	eventOpString := event.Op.String()
	logrus.Debugf("Handling file system event: operation %s on %s", eventOpString, event.Name)

	if event.Name == initializedStamp && utils.PathExists(initializedStamp) {
		return true
	}

	return false
}

func handlePollingTickForStamp(time time.Time, initializedStamp string) bool {
	timeString := time.String()
	logrus.Debugf("Handling polling tick %s", timeString)

	if utils.PathExists(initializedStamp) {
		return true
	}

	return false
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

func isUsePollingSet() bool {
	valueString := os.Getenv("TOOLBX_RUN_USE_POLLING")
	if valueString == "" {
		return false
	}

	if valueBool, err := strconv.ParseBool(valueString); err == nil {
		return valueBool
	}

	return true
}

func saveCDISpecTo(spec *specs.Spec, path string) error {
	if path == "" {
		panic("path not specified")
	}

	if spec == nil {
		panic("spec not specified")
	}

	logrus.Debugf("Saving Container Device Interface to file %s", path)

	if extension := filepath.Ext(path); extension != ".json" {
		panicMsg := fmt.Sprintf("path has invalid extension %s", extension)
		panic(panicMsg)
	}

	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		logrus.Debugf("Saving Container Device Interface: failed to marshal JSON: %s", err)
		return errors.New("failed to marshal Container Device Interface to JSON")
	}

	if err := renameio.WriteFile(path, data, 0644); err != nil {
		logrus.Debugf("Saving Container Device Interface: failed to write file: %s", err)
		return errors.New("failed to write Container Device Interface to file")
	}

	return nil
}

func showEntryPointLog(line string) error {
	var logLevel logrus.Level
	var logLevelFound bool
	var logMsg string

	reader := strings.NewReader(line)
	decoder := logfmt.NewDecoder(reader)

	if decoder.ScanRecord() {
		for decoder.ScanKeyval() {
			value := decoder.Value()
			valueString := string(value)

			switch key := decoder.Key(); string(key) {
			case "level":
				logLevelFound = true

				var err error
				logLevel, err = logrus.ParseLevel(valueString)
				if err != nil {
					logrus.Debugf("Parsing entry point log-level %s failed: %s",
						valueString,
						err)
					logLevel = logrus.DebugLevel
				}
			case "msg":
				logMsg = valueString
			}
		}
	}

	if err := decoder.Err(); err != nil {
		return err
	}

	if !logLevelFound {
		errMsg, _ := strings.CutPrefix(line, "Error: ")
		return &entryPointError{errMsg}
	}

	logger := logrus.StandardLogger()
	logger.Logf(logLevel, "> %s", logMsg)
	return nil
}

func showEntryPointLogs(container string, since time.Time) error {
	var stderr strings.Builder
	if err := podman.Logs(container, since, &stderr); err != nil {
		return err
	}

	var errReceived error
	logs := stderr.String()
	reader := strings.NewReader(logs)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		line := scanner.Text()
		if err := showEntryPointLog(line); err != nil {
			var errEntryPoint *entryPointError

			if errors.As(err, &errEntryPoint) {
				if !errors.Is(errReceived, err) {
					errReceived = errors.Join(errReceived, err)
				}
			} else {
				logrus.Debugf("Parsing entry point log failed: %s:", err)
				logrus.Debugf("%s", line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		if errReceived == nil {
			return err
		}

		logrus.Debugf("Reading logs from container %s failed: %s", container, err)
	}

	return errReceived
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

func startP11KitServer() ([]string, error) {
	serverSocket, err := utils.GetP11KitServerSocket(currentUser)
	if err != nil {
		return nil, err
	}

	const logPrefix = "Starting 'p11-kit server'"
	logrus.Debugf("%s with socket %s", logPrefix, serverSocket)

	serverSocketLock, err := utils.GetP11KitServerSocketLock(currentUser)
	if err != nil {
		return nil, err
	}

	serverSocketLockFile, err := utils.Flock(serverSocketLock, syscall.LOCK_EX)
	if err != nil {
		logrus.Debugf("%s: %s", logPrefix, err)

		var errFlock *utils.FlockError

		if errors.As(err, &errFlock) {
			if errors.Is(err, utils.ErrFlockAcquire) {
				err = utils.ErrFlockAcquire
			} else if errors.Is(err, utils.ErrFlockCreate) {
				err = utils.ErrFlockCreate
			} else {
				panicMsg := fmt.Sprintf("unexpected %T: %s", err, err)
				panic(panicMsg)
			}
		}

		return nil, err
	}

	defer serverSocketLockFile.Close()

	serverSocketAddress := fmt.Sprintf("P11_KIT_SERVER_ADDRESS=unix:path=%s", serverSocket)
	serverEnviron := []string{
		serverSocketAddress,
	}

	if utils.PathExists(serverSocket) {
		logrus.Debugf("%s: socket %s already exists", logPrefix, serverSocket)
		logrus.Debugf("%s: skipping", logPrefix)
		return serverEnviron, nil
	}

	serverArgs := []string{
		"server",
		"--name", serverSocket,
		"--provider", "p11-kit-trust.so",
		"pkcs11:model=p11-kit-trust?write-protected=yes",
	}

	if err := shell.Run("p11-kit", nil, nil, nil, serverArgs...); err != nil {
		logrus.Debugf("%s failed: %s", logPrefix, err)
		return nil, nil
	}

	return serverEnviron, nil
}

func (err *entryPointError) Error() string {
	return err.msg
}

func (err *entryPointError) Is(target error) bool {
	targetErrMsg := target.Error()
	return err.msg == targetErrMsg
}

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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/containers/toolbox/pkg/nvidia"
	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/containers/toolbox/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	cgroupsVersion int

	currentUser *user.User

	executable string

	executableBase string

	rootCmd = &cobra.Command{
		Use:               "toolbox",
		Short:             "Tool for interactive command line environments on Linux",
		PersistentPreRunE: preRun,
		RunE:              rootRun,
		Version:           version.GetVersion(),
	}

	rootFlags struct {
		assumeYes bool
		logLevel  string
		logPodman bool
		verbose   int
	}

	workingDirectory string
)

type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string {
	if e.err != nil {
		return e.err.Error()
	} else {
		return ""
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if rootCmd.SilenceErrors {
			if errMsg := err.Error(); errMsg != "" {
				fmt.Fprintf(os.Stderr, "Error: %s\n", errMsg)
			}
		}

		var errExit *exitError
		if errors.As(err, &errExit) {
			os.Exit(errExit.code)
		}

		os.Exit(1)
	}

	os.Exit(0)
}

func init() {
	if err := setUpGlobals(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	persistentFlags := rootCmd.PersistentFlags()

	persistentFlags.BoolVarP(&rootFlags.assumeYes,
		"assumeyes",
		"y",
		false,
		"Automatically answer yes for all questions")

	persistentFlags.StringVar(&rootFlags.logLevel,
		"log-level",
		"error",
		"Log messages at the specified level: trace, debug, info, warn, error, fatal or panic")

	persistentFlags.BoolVar(&rootFlags.logPodman,
		"log-podman",
		false,
		"Show the log output of Podman. The log level is handled by the log-level option")

	persistentFlags.CountVarP(&rootFlags.verbose, "verbose", "v", "Set log-level to 'debug'")

	if err := rootCmd.RegisterFlagCompletionFunc("log-level", completionLogLevels); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}

	rootCmd.SetHelpFunc(rootHelp)

	usageTemplate := fmt.Sprintf("Run '%s --help' for usage.", executableBase)
	rootCmd.SetUsageTemplate(usageTemplate)
}

func preRun(cmd *cobra.Command, args []string) error {
	cmd.Root().SilenceErrors = true
	cmd.Root().SilenceUsage = true

	if err := setUpLoggers(); err != nil {
		return err
	}

	logrus.Debugf("Running as real user ID %s", currentUser.Uid)
	logrus.Debugf("Resolved absolute path to the executable as %s", executable)

	if !utils.IsInsideContainer() {
		logrus.Debugf("Running on a cgroups v%d host", cgroupsVersion)

		if _, err := validateSubIDRanges(cmd, args, currentUser); err != nil {
			return err
		}
	}

	toolbxDelayEntryPoint := os.Getenv("TOOLBX_DELAY_ENTRY_POINT")
	logrus.Debugf("TOOLBX_DELAY_ENTRY_POINT is %s", toolbxDelayEntryPoint)

	toolbxFailEntryPoint := os.Getenv("TOOLBX_FAIL_ENTRY_POINT")
	logrus.Debugf("TOOLBX_FAIL_ENTRY_POINT is %s", toolbxFailEntryPoint)

	toolboxPath := os.Getenv("TOOLBOX_PATH")

	if toolboxPath == "" {
		if utils.IsInsideContainer() {
			if err := preRunIsCoreOSBug(); err != nil {
				return err
			}

			return errors.New("TOOLBOX_PATH not set")
		}

		os.Setenv("TOOLBOX_PATH", executable)
		toolboxPath = os.Getenv("TOOLBOX_PATH")
	}

	logrus.Debugf("TOOLBOX_PATH is %s", toolboxPath)

	if err := migrate(cmd, args); err != nil {
		return err
	}

	if err := utils.SetUpConfiguration(); err != nil {
		return err
	}

	return nil
}

func rootHelp(cmd *cobra.Command, args []string) {
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

	manual := "toolbox"

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			manual = manual + "-" + arg
			break
		}
	}

	if err := showManual(manual); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func rootRun(cmd *cobra.Command, args []string) error {
	return rootRunImpl(cmd, args)
}

func migrate(cmd *cobra.Command, args []string) error {
	logrus.Debug("Migrating to newer Podman")

	if utils.IsInsideContainer() {
		logrus.Debug("Migration not needed: running inside a container")
		return nil
	}

	if cmdName, completionCmdName := cmd.Name(), completionCmd.Name(); cmdName == completionCmdName {
		logrus.Debugf("Migration not needed: command %s doesn't need it", cmdName)
		return nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		logrus.Debugf("Migrating to newer Podman: failed to get the user config directory: %s", err)
		return errors.New("failed to get the user config directory")
	}

	toolboxConfigDir := configDir + "/toolbox"
	stampPath := toolboxConfigDir + "/podman-system-migrate"
	logrus.Debugf("Toolbx config directory is %s", toolboxConfigDir)

	podmanVersion, err := podman.GetVersion()
	if err != nil {
		logrus.Debugf("Migrating to newer Podman: failed to get the Podman version: %s", err)
		return errors.New("failed to get the Podman version")
	}

	logrus.Debugf("Current Podman version is %s", podmanVersion)

	err = os.MkdirAll(toolboxConfigDir, 0775)
	if err != nil {
		logrus.Debugf("Migrating to newer Podman: failed to create configuration directory %s: %s",
			toolboxConfigDir,
			err)
		return errors.New("failed to create configuration directory")
	}

	toolboxRuntimeDirectory, err := utils.GetRuntimeDirectory(currentUser)
	if err != nil {
		return err
	}

	migrateLock := toolboxRuntimeDirectory + "/migrate.lock"

	migrateLockFile, err := utils.Flock(migrateLock, syscall.LOCK_EX)
	if err != nil {
		logrus.Debugf("Migrating to newer Podman: %s", err)

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

		return err
	}

	defer migrateLockFile.Close()

	stampBytes, err := ioutil.ReadFile(stampPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			logrus.Debugf("Migrating to newer Podman: failed to read migration stamp file %s: %s",
				stampPath,
				err)
			return errors.New("failed to read migration stamp file")
		}
	} else {
		stampString := string(stampBytes)
		podmanVersionOld := strings.TrimSpace(stampString)
		if podmanVersionOld != "" {
			logrus.Debugf("Old Podman version is %s", podmanVersionOld)

			if podmanVersion == podmanVersionOld {
				logrus.Debugf("Migration not needed: Podman version %s is unchanged", podmanVersion)
				return nil
			}

			if !podman.CheckVersion(podmanVersionOld) {
				logrus.Debugf("Migration not needed: Podman version %s is old", podmanVersion)
				return nil
			}
		}
	}

	if err = podman.SystemMigrate(""); err != nil {
		logrus.Debugf("Migrating to newer Podman: failed to migrate containers: %s", err)
		return errors.New("failed to migrate containers")
	}

	logrus.Debugf("Migration to Podman version %s was ok", podmanVersion)
	logrus.Debugf("Updating Podman version in %s", stampPath)

	podmanVersionBytes := []byte(podmanVersion + "\n")
	err = ioutil.WriteFile(stampPath, podmanVersionBytes, 0664)
	if err != nil {
		logrus.Debugf("Migrating to newer Podman: failed to update Podman version in migration stamp file %s: %s",
			stampPath,
			err)
		return errors.New("failed to update Podman version in migration stamp file")
	}

	return nil
}

func newSubIDError() error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Missing subgid and/or subuid ranges for user %s\n", currentUser.Username)
	fmt.Fprintf(&builder, "See the podman(1), subgid(5), subuid(5) and usermod(8) manuals for more\n")
	fmt.Fprintf(&builder, "information.")

	errMsg := builder.String()
	return errors.New(errMsg)
}

func setUpGlobals() error {
	var err error

	if !utils.IsInsideContainer() {
		cgroupsVersion, err = utils.GetCgroupsVersion()
		if err != nil {
			return fmt.Errorf("failed to get the cgroups version: %w", err)
		}
	}

	currentUser, err = user.Current()
	if err != nil {
		return fmt.Errorf("failed to get the current user: %w", err)
	}

	executable, err = os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get the path to the executable: %w", err)
	}

	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path to the executable: %w", err)
	}

	executableBase = filepath.Base(executable)

	workingDirectory, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get the working directory: %w", err)
	}

	return nil
}

func setUpLoggers() error {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	if rootFlags.verbose > 0 {
		rootFlags.logLevel = "debug"
	}

	logLevel, err := logrus.ParseLevel(rootFlags.logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log-level: %w", err)
	}

	logrus.SetLevel(logLevel)

	if rootFlags.verbose > 1 {
		nvidia.SetLogLevel(logLevel)
		rootFlags.logPodman = true
	}

	if rootFlags.logPodman {
		podman.SetLogLevel(logLevel)
	}

	return nil
}

func validateSubIDRanges(cmd *cobra.Command, args []string, user *user.User) (bool, error) {
	logrus.Debugf("Looking up sub-GID and sub-UID ranges for user %s", user.Username)

	if user.Uid == "0" {
		logrus.Debugf("Look-up not needed: user %s doesn't need them", user.Username)
		return true, nil
	}

	if utils.IsInsideContainer() {
		logrus.Debug("Look-up not needed: running inside a container")
		return true, nil
	}

	if cmdName, completionCmdName := cmd.Name(), completionCmd.Name(); cmdName == completionCmdName {
		logrus.Debugf("Look-up not needed: command %s doesn't need them", cmdName)
		return true, nil
	}

	if _, err := utils.ValidateSubIDRanges(user); err != nil {
		logrus.Debugf("Looking up sub-GID and sub-UID ranges failed: %s", err)
		return false, newSubIDError()
	}

	return true, nil
}

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

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

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
		Short:             "Unprivileged development environment",
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

func Execute() {
	if err := setUpGlobals(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func init() {
	persistentFlags := rootCmd.PersistentFlags()

	persistentFlags.BoolVarP(&rootFlags.assumeYes,
		"assumeyes",
		"y",
		false,
		"Automatically answer yes for all questions.")

	persistentFlags.StringVar(&rootFlags.logLevel,
		"log-level",
		"error",
		"Log messages at the specified level: trace, debug, info, warn, error, fatal or panic.")

	persistentFlags.BoolVar(&rootFlags.logPodman,
		"log-podman",
		false,
		"Show the log output of Podman. The log level is handled by the log-level option.")

	persistentFlags.CountVarP(&rootFlags.verbose, "verbose", "v", "Set log-level to 'debug'.")

	rootCmd.SetHelpFunc(rootHelp)
	rootCmd.SetUsageFunc(rootUsage)
}

func preRun(cmd *cobra.Command, args []string) error {
	cmd.Root().SilenceUsage = true

	if err := setUpLoggers(); err != nil {
		return err
	}

	logrus.Debugf("Running as real user ID %s", currentUser.Uid)
	logrus.Debugf("Resolved absolute path to the executable as %s", executable)

	if !utils.IsInsideContainer() {
		logrus.Debugf("Running on a cgroups v%d host", cgroupsVersion)
	}

	toolboxPath := os.Getenv("TOOLBOX_PATH")

	if utils.IsInsideContainer() {
		if toolboxPath == "" {
			return errors.New("TOOLBOX_PATH not set")
		}
	} else {
		if currentUser.Uid != "0" {
			logrus.Debugf("Checking if /etc/subgid and /etc/subuid have entries for user %s",
				currentUser.Username)

			if _, err := validateSubIDFile("/etc/subuid"); err != nil {
				return newSubIDFileError()
			}

			if _, err := validateSubIDFile("/etc/subgid"); err != nil {
				return newSubIDFileError()
			}
		}

		if toolboxPath == "" {
			os.Setenv("TOOLBOX_PATH", executable)
		}
	}

	toolboxPath = os.Getenv("TOOLBOX_PATH")
	logrus.Debugf("TOOLBOX_PATH is %s", toolboxPath)

	if cmd.Use != "reset" {
		if err := migrate(); err != nil {
			return err
		}
	}

	return nil
}

func rootHelp(cmd *cobra.Command, args []string) {
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

	manual := "toolbox"

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			manual = manual + "-" + arg
			break
		}
	}

	if err := utils.ShowManual(manual); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func rootRun(cmd *cobra.Command, args []string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "missing command\n")
	fmt.Fprintf(&builder, "\n")
	fmt.Fprintf(&builder, "create    Create a new toolbox container\n")
	fmt.Fprintf(&builder, "enter     Enter an existing toolbox container\n")
	fmt.Fprintf(&builder, "list      List all existing toolbox containers and images\n")
	fmt.Fprintf(&builder, "\n")
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func rootUsage(cmd *cobra.Command) error {
	err := fmt.Errorf("Run '%s --help' for usage.", executableBase)
	fmt.Fprintf(os.Stderr, "%s", err)
	return err
}

func migrate() error {
	if utils.IsInsideContainer() {
		return nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get the user config directory")
	}

	toolboxConfigDir := configDir + "/toolbox"
	stampPath := toolboxConfigDir + "/podman-system-migrate"
	logrus.Debugf("Toolbox config directory is %s", toolboxConfigDir)

	podmanVersion, err := podman.GetVersion()
	if err != nil {
		return fmt.Errorf("failed to get the Podman version: %w", err)
	}

	logrus.Debugf("Current Podman version is %s", podmanVersion)

	err = os.MkdirAll(toolboxConfigDir, 0775)
	if err != nil {
		return fmt.Errorf("failed to create configuration directory: %w", err)
	}

	runtimeDirectory := os.Getenv("XDG_RUNTIME_DIR")
	toolboxRuntimeDirectory := runtimeDirectory + "/toolbox"
	if err := os.MkdirAll(toolboxRuntimeDirectory, 0700); err != nil {
		return fmt.Errorf("failed to create runtime directory %s: %w", toolboxRuntimeDirectory, err)
	}

	lockFile := toolboxRuntimeDirectory + "/migrate.lock"

	lockFD, err := syscall.Open(lockFile,
		syscall.O_CREAT|syscall.O_WRONLY,
		syscall.S_IRUSR|syscall.S_IWUSR|syscall.S_IRGRP|syscall.S_IWGRP|syscall.S_IROTH)
	if err != nil {
		return fmt.Errorf("failed to open migration lock file: %w", err)
	}

	defer syscall.Close(lockFD)

	err = syscall.Flock(lockFD, syscall.LOCK_EX)
	if err != nil {
		return fmt.Errorf("failed to acquire migration lock: %w", err)
	}

	stampBytes, err := ioutil.ReadFile(stampPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to read migration stamp file: %w", err)
		}
	} else {
		stampString := string(stampBytes)
		podmanVersionOld := strings.TrimSpace(stampString)
		if podmanVersionOld != "" {
			logrus.Debugf("Old Podman version is %s", podmanVersionOld)

			if podmanVersion == podmanVersionOld {
				logrus.Debugf("Migration not needed: Podman version %s is unchanged",
					podmanVersion)

				return nil
			}

			if !podman.CheckVersion(podmanVersionOld) {
				logrus.Debugf("Migration not needed: Podman version %s is old", podmanVersion)
				return nil
			}
		}
	}

	if err = podman.SystemMigrate(""); err != nil {
		return fmt.Errorf("failed to migrate containers: %w", err)
	}

	logrus.Debugf("Migration to Podman version %s was ok", podmanVersion)
	logrus.Debugf("Updating Podman version in %s", stampPath)

	podmanVersionBytes := []byte(podmanVersion + "\n")
	err = ioutil.WriteFile(stampPath, podmanVersionBytes, 0664)
	if err != nil {
		return fmt.Errorf("failed to update Podman version in migration stamp file: %w", err)
	}

	return nil
}

func newSubIDFileError() error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "/etc/subgid and /etc/subuid don't have entries for user %s\n", currentUser.Username)
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
		rootFlags.logPodman = true
	}

	if rootFlags.logPodman {
		podman.SetLogLevel(logLevel)
	}

	return nil
}

func validateSubIDFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("failed to open %s: %w", path, err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	prefix := currentUser.Username + ":"

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			return true, nil
		}
	}

	return false, fmt.Errorf("failed to find an entry for user %s in %s", currentUser.Username, path)
}

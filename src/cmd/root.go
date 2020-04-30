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
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

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
	fmt.Fprintf(&builder, "Run 'toolbox --help' for usage.")

	errMsg := builder.String()
	return errors.New(errMsg)
}

func rootUsage(cmd *cobra.Command) error {
	err := errors.New("Run 'toolbox --help' for usage.")
	fmt.Fprintf(os.Stderr, "%s", err)
	return err
}

func setUpGlobals() error {
	var err error

	if !utils.IsInsideContainer() {
		cgroupsVersion, err = utils.GetCgroupsVersion()
		if err != nil {
			return errors.New("failed to get the cgroups version")
		}
	}

	currentUser, err = user.Current()
	if err != nil {
		return errors.New("failed to get the current user")
	}

	executable, err = os.Executable()
	if err != nil {
		return errors.New("failed to get the path to the executable")
	}

	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		return errors.New("failed to resolve absolute path to the executable")
	}

	executableBase = filepath.Base(executable)

	workingDirectory, err = os.Getwd()
	if err != nil {
		return errors.New("failed to get the working directory")
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
		return errors.New("failed to parse log-level")
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

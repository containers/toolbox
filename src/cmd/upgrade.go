/*
 * Copyright © 2025 - 2026 Hadi Chokr <hadichokr@icloud.com>
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

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/term"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	labelPackageManagerUpdate = "com.github.containers.toolbx.package-manager.update"

	labelPackageManagerUpdateAssumeYes = "com.github.containers.toolbx.package-manager.update.assumeyes"
)

type upgradeUnsupportedError struct {
	Container string
}

var (
	upgradeFlags struct {
		upgradeAll bool
	}
)

var upgradeCmd = &cobra.Command{
	Use:               "upgrade",
	Short:             "Upgrade the packages inside one or more Toolbx containers",
	RunE:              upgrade,
	ValidArgsFunction: completionContainerNamesFiltered,
}

func init() {
	flags := upgradeCmd.Flags()

	flags.BoolVarP(&upgradeFlags.upgradeAll, "all", "a", false, "Upgrade all Toolbx containers")

	upgradeCmd.SetHelpFunc(upgradeHelp)
	rootCmd.AddCommand(upgradeCmd)
}

func upgrade(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	if upgradeFlags.upgradeAll && len(args) != 0 {
		var builder strings.Builder
		fmt.Fprintf(&builder, "option --all and CONTAINER cannot be used together\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	var containers []string

	if upgradeFlags.upgradeAll {
		logrus.Debug("Getting all containers")

		toolbxContainers, err := podman.GetContainers()
		if err != nil {
			logrus.Debugf("Getting all containers failed: %s", err)
			return errors.New("failed to get containers")
		}

		for toolbxContainers.Next() {
			container := toolbxContainers.Get()
			name := container.Name()
			containers = append(containers, name)
		}
	} else {
		if len(args) == 0 {
			var builder strings.Builder
			fmt.Fprintf(&builder, "missing argument for \"upgrade\"\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}

		containers = args
	}

	var failed bool

	for _, container := range containers {
		err := upgradeContainer(container)
		if err == nil {
			continue
		}

		var errUnsupported *upgradeUnsupportedError
		if upgradeFlags.upgradeAll && errors.As(err, &errUnsupported) {
			fmt.Fprintf(os.Stderr,
				"Warning: skipping container %s: it doesn't support 'upgrade'\n",
				container)
			continue
		}

		if len(containers) == 1 {
			return err
		}

		if errMsg := err.Error(); errMsg != "" {
			fmt.Fprintf(os.Stderr, "Error: %s\n", errMsg)
		}

		failed = true
	}

	if failed {
		return &exitError{1, nil}
	}

	return nil
}

func upgradeContainer(container string) error {
	logrus.Debugf("Inspecting container %s", container)

	containerObj, err := podman.InspectContainer(container)
	if err != nil {
		logrus.Debugf("Inspecting container %s failed: %s", container, err)
		return fmt.Errorf("failed to inspect container %s", container)
	}

	if !containerObj.IsToolbx() {
		return fmt.Errorf("%s is not a Toolbx container", container)
	}

	labels := containerObj.Labels()
	commandInteractive := labels[labelPackageManagerUpdate]
	commandAssumeYes := labels[labelPackageManagerUpdateAssumeYes]

	logrus.Debugf("Upgrade command of container %s is %s", container, commandInteractive)
	logrus.Debugf("Unattended upgrade command of container %s is %s", container, commandAssumeYes)

	if commandInteractive == "" && commandAssumeYes == "" {
		return &upgradeUnsupportedError{container}
	}

	var command string

	if rootFlags.assumeYes {
		command = commandAssumeYes
		if command == "" {
			var builder strings.Builder
			fmt.Fprintf(&builder, "container %s doesn't support unattended upgrades\n", container)
			fmt.Fprintf(&builder, "Upgrade it without '--assumeyes', or recreate it with a newer image.")

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	} else {
		command = commandInteractive
		if command == "" {
			command = commandAssumeYes
		}

		if !term.IsTerminal(os.Stdin) || !term.IsTerminal(os.Stdout) {
			var builder strings.Builder
			fmt.Fprintf(&builder, "confirmation required to upgrade container %s\n", container)
			fmt.Fprintf(&builder, "Use option '--assumeyes' to upgrade without confirmation.\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	}

	fmt.Printf("Upgrading container %s with: %s\n", container, command)

	if !rootFlags.assumeYes {
		if !askForConfirmation("Continue? [y/N]:") {
			return nil
		}
	}

	if err := runCommand(container,
		false,
		"",
		"",
		0,
		[]string{"sh", "-c", command},
		false,
		false,
		true); err != nil {
		return err
	}

	return nil
}

func upgradeHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-upgrade"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func (err *upgradeUnsupportedError) Error() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "container %s does not support \"upgrade\"\n", err.Container)
	fmt.Fprintf(&builder, "Recreate it with a newer image or file a bug.")

	errMsg := builder.String()
	return errMsg
}

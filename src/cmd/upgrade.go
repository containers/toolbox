/*
 * Copyright © 2025 Hadi Chokr <hadichokr@icloud.com>
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

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	upgradeAll       bool
	upgradeContainer string
)

var upgradeCmd = &cobra.Command{
	Use:               "upgrade [container]",
	Short:             "Detect package manager and upgrade packages in toolbx containers",
	Args:              cobra.MaximumNArgs(1),
	RunE:              runUpgrade,
	ValidArgsFunction: completionContainerNamesFiltered,
}

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeAll, "all", false, "Upgrade all Toolbx containers")
	upgradeCmd.Flags().StringVar(&upgradeContainer, "container", "", "Name of the toolbox container to upgrade")

	// Register container flag completion
	if err := upgradeCmd.RegisterFlagCompletionFunc("container", completionContainerNames); err != nil {
		fmt.Fprintf(os.Stderr, "failed to register flag completion function: %v\n", err)
		os.Exit(1)
	}
	upgradeCmd.SetHelpFunc(upgradeHelp)
	rootCmd.AddCommand(upgradeCmd)
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	// Use positional argument as container name if --container not set
	if upgradeContainer == "" && len(args) == 1 {
		upgradeContainer = args[0]
	}

	if !upgradeAll && upgradeContainer == "" {
		return errors.New("must specify either --all or a container name")
	}

	if upgradeAll && upgradeContainer != "" {
		return errors.New("cannot specify both --all and a container name")
	}

	if upgradeAll {
		containers, err := getContainers()
		if err != nil {
			return err
		}

		if len(containers) == 0 {
			return errors.New("no Toolbx containers found")
		}

		for _, container := range containers {
			fmt.Printf("Upgrading container: %s\n", container.Name())
			if err := execUpgradeInContainer(container.Name()); err != nil {
				fmt.Fprintf(os.Stderr, "Error upgrading container %s: %v\n", container.Name(), err)
			}
		}
		return nil
	}

	return execUpgradeInContainer(upgradeContainer)
}

// execUpgradeInContainer runs detection and upgrade inside the specified container
func execUpgradeInContainer(container string) error {
	pkgs := []struct {
		detect  string
		upgrade string
		name    string
	}{
		{"command -v dnf", "sudo dnf -y upgrade", "dnf"},
		{"command -v microdnf", "sudo microdnf upgrade -y", "microdnf"},
		{"command -v yum", "sudo yum -y upgrade", "yum"},
		{"command -v apt", "sudo apt update && sudo apt upgrade -y", "apt"},
		{"command -v pacman", "sudo pacman -Syu --noconfirm", "pacman"},
		{"command -v xbps-install", "sudo xbps-install -Su -y", "xbps"},
		{"command -v zypper", "sudo zypper update -y", "zypper"},
		{"command -v apk", "sudo apk update && sudo apk upgrade", "apk"},
		{"command -v emerge", "sudo emerge --sync && sudo emerge -uDN @world", "emerge"},
		{"command -v slackpkg", "sudo slackpkg update && sudo slackpkg upgrade-all", "slackpkg"},
		{"command -v swupd", "sudo swupd update", "swupd"},
	}

	for _, pkg := range pkgs {
		// Use runCommand to check if package manager exists
		err := runCommand(container,
				  false,  // defaultContainer
		    "",     // image
		    "",     // release
		    0,      // preserveFDs
		    []string{"sh", "-c", pkg.detect},
		    false,  // emitEscapeSequence
		    false,  // fallbackToBash
		    true)   // pedantic

		if err == nil {
			fmt.Printf("Detected package manager: %s\n", pkg.name)

			// Use runCommand to execute the upgrade
			upgradeErr := runCommand(container,
						 false,
			    "",
			    "",
			    0,
			    []string{"sh", "-c", pkg.upgrade},
			    false,
			    false,
			    true)

			return upgradeErr
		}
	}
	return errors.New("no supported package manager found")
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

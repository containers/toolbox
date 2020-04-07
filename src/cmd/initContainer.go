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
	"fmt"
	"os"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	initContainerFlags struct {
		home        string
		homeLink    bool
		mediaLink   bool
		mntLink     bool
		monitorHost bool
		shell       string
		uid         int
		user        string
	}
)

var initContainerCmd = &cobra.Command{
	Use:    "init-container",
	Short:  "Initialize a running container",
	Hidden: true,
	RunE:   initContainer,
}

func init() {
	flags := initContainerCmd.Flags()

	flags.StringVar(&initContainerFlags.home,
		"home",
		"",
		"Create a user inside the toolbox container whose login directory is HOME.")
	initContainerCmd.MarkFlagRequired("home")

	flags.BoolVar(&initContainerFlags.homeLink,
		"home-link",
		false,
		"Make /home a symbolic link to /var/home.")

	flags.BoolVar(&initContainerFlags.mediaLink,
		"media-link",
		false,
		"Make /media a symbolic link to /run/media.")

	flags.BoolVar(&initContainerFlags.mntLink, "mnt-link", false, "Make /mnt a symbolic link to /var/mnt.")

	flags.BoolVar(&initContainerFlags.monitorHost,
		"monitor-host",
		false,
		"Ensure that certain configuration files inside the toolbox container are in sync with the host.")

	flags.StringVar(&initContainerFlags.shell,
		"shell",
		"",
		"Create a user inside the toolbox container whose login shell is SHELL.")
	initContainerCmd.MarkFlagRequired("shell")

	flags.IntVar(&initContainerFlags.uid,
		"uid",
		0,
		"Create a user inside the toolbox container whose numerical user ID is UID.")
	initContainerCmd.MarkFlagRequired("uid")

	flags.StringVar(&initContainerFlags.user,
		"user",
		"",
		"Create a user inside the toolbox container whose login name is USER.")
	initContainerCmd.MarkFlagRequired("user")

	initContainerCmd.SetHelpFunc(initContainerHelp)
	rootCmd.AddCommand(initContainerCmd)
}

func initContainer(cmd *cobra.Command, args []string) error {
	return nil
}

func initContainerHelp(cmd *cobra.Command, args []string) {
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

	if err := utils.ShowManual("toolbox-init-container"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

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
	"os"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	enterFlags struct {
		container string
		distro    string
		release   string
	}
)

var enterCmd = &cobra.Command{
	Use:               "enter",
	Short:             "Enter a Toolbx container for interactive use",
	RunE:              enter,
	ValidArgsFunction: completionContainerNamesFiltered,
}

func init() {
	flags := enterCmd.Flags()

	flags.StringVarP(&enterFlags.container,
		"container",
		"c",
		"",
		"Enter a Toolbx container with the given name")

	flags.StringVarP(&enterFlags.distro,
		"distro",
		"d",
		"",
		"Enter a Toolbx container for a different operating system distribution than the host")

	flags.StringVarP(&enterFlags.release,
		"release",
		"r",
		"",
		"Enter a Toolbx container for a different operating system release than the host")

	if err := enterCmd.RegisterFlagCompletionFunc("container", completionContainerNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}
	if err := enterCmd.RegisterFlagCompletionFunc("distro", completionDistroNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}

	enterCmd.SetHelpFunc(enterHelp)
	rootCmd.AddCommand(enterCmd)
}

func enter(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	var container string
	var containerArg string
	var defaultContainer bool = true

	if len(args) != 0 {
		container = args[0]
		containerArg = "CONTAINER"
	} else if enterFlags.container != "" {
		container = enterFlags.container
		containerArg = "--container"
	}

	if container != "" {
		defaultContainer = false
	}

	if enterFlags.release != "" {
		defaultContainer = false
	}

	container, image, release, err := resolveContainerAndImageNames(container,
		containerArg,
		enterFlags.distro,
		"",
		enterFlags.release)

	if err != nil {
		return err
	}

	command := []string{"toolbox", "sh", "--", "-l"}

	if err := runCommand(container, defaultContainer, image, release, 0, command, true, true, false); err != nil {
		return err
	}

	return nil
}

func enterHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-enter"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

/*
 * Copyright © 2021 Red Hat Inc.
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
	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename an existing toolbox container",
	RunE:  rename,
}

func init() {
	renameCmd.SetHelpFunc(renameHelp)
	rootCmd.AddCommand(renameCmd)
}

func rename(cmd *cobra.Command, args []string) error {
	var container, newName string
	var err error

	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a toolbox container")
		}

		if _, err = utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	if len(args) != 2 {
		var builder strings.Builder
		fmt.Fprintf(&builder, "The 'rename' command takes two arguments\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)
		errMsg := builder.String()

		return errors.New(errMsg)
	}

	container = args[0]
	newName = args[1]

	if !podman.CheckVersion("3.0.0") {
		var builder strings.Builder
		fmt.Fprintf(&builder, "The 'rename' command requires Podman v3.0.0\n")
		fmt.Fprintf(&builder, "Please, upgrade your version of Podman.")
		errMsg := builder.String()

		return errors.New(errMsg)
	}

	if ok := utils.IsContainerNameValid(newName); !ok {
		var builder strings.Builder
		fmt.Fprintf(&builder, "Invalid argument for NEWNAME\n")
		fmt.Fprintf(&builder, "Container names must match '%s'\n", utils.ContainerNameRegexp)
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)
		errMsg := builder.String()

		return errors.New(errMsg)
	}

	if ok, _ := podman.IsToolboxContainer(container); !ok {
		var builder strings.Builder
		fmt.Fprintf(&builder, "Invalid argument for CONTAINER\n")
		fmt.Fprintf(&builder, "Container %s is not a toolbox container\n", container)
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)
		errMsg := builder.String()

		return errors.New(errMsg)
	}

	err = podman.Rename(container, newName)
	if err != nil {
		return fmt.Errorf("failed to rename container %s: %w", container, err)
	}

	return nil
}

func renameHelp(cmd *cobra.Command, args []string) {
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

	if err := utils.ShowManual("toolbox-rename"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

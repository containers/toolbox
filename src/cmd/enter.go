/*
 * Copyright © 2019 – 2021 Red Hat Inc.
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
	Use:   "enter",
	Short: "Enter a toolbox container for interactive use",
	RunE:  enter,
}

func init() {
	flags := enterCmd.Flags()

	flags.StringVarP(&enterFlags.container,
		"container",
		"c",
		"",
		"Enter a toolbox container with the given name")

	flags.StringVarP(&enterFlags.distro,
		"distro",
		"d",
		"",
		"Enter a toolbox container for a different operating system distribution than the host")

	flags.StringVarP(&enterFlags.release,
		"release",
		"r",
		"",
		"Enter a toolbox container for a different operating system release than the host")

	enterCmd.SetHelpFunc(enterHelp)
	rootCmd.AddCommand(enterCmd)
}

func enter(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a toolbox container")
		}

		if _, err := utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	var container string
	var containerArg string
	var nonDefaultContainer bool

	if len(args) != 0 {
		container = args[0]
		containerArg = "CONTAINER"
	} else if enterFlags.container != "" {
		container = enterFlags.container
		containerArg = "--container"
	}

	if container != "" {
		nonDefaultContainer = true

		if !utils.IsContainerNameValid(container) {
			var builder strings.Builder
			fmt.Fprintf(&builder, "invalid argument for '%s'\n", containerArg)
			fmt.Fprintf(&builder, "Container names must match '%s'\n", utils.ContainerNameRegexp)
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}
	}

	var release string
	if enterFlags.release != "" {
		nonDefaultContainer = true

		var err error
		release, err = utils.ParseRelease(enterFlags.distro, enterFlags.release)
		if err != nil {
			err := utils.CreateErrorInvalidRelease(executableBase)
			return err
		}
	}

	container, image, release, err := utils.ResolveContainerAndImageNames(container, enterFlags.distro, "", release)
	if err != nil {
		return err
	}

	hostID, err := utils.GetHostID()
	if err != nil {
		return fmt.Errorf("failed to get the host ID: %w", err)
	}

	hostVariantID, err := utils.GetHostVariantID()
	if err != nil {
		return errors.New("failed to get the host VARIANT_ID")
	}

	var emitEscapeSequence bool

	if hostID == "fedora" && (hostVariantID == "silverblue" || hostVariantID == "workstation") {
		emitEscapeSequence = true
	}

	return runCommand(container,
		!nonDefaultContainer,
		image,
		release,
		nil,
		emitEscapeSequence,
		false)
}

func enterHelp(cmd *cobra.Command, args []string) {
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

	if err := utils.ShowManual("toolbox-enter"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

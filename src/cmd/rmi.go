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
	"strings"

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	rmiFlags struct {
		deleteAll   bool
		forceDelete bool
	}
)

var rmiCmd = &cobra.Command{
	Use:               "rmi",
	Short:             "Remove one or more Toolbx images",
	RunE:              rmi,
	ValidArgsFunction: completionImageNamesFiltered,
}

func init() {
	flags := rmiCmd.Flags()

	flags.BoolVarP(&rmiFlags.deleteAll, "all", "a", false, "Remove all Toolbx images")

	flags.BoolVarP(&rmiFlags.forceDelete,
		"force",
		"f",
		false,
		"Force the removal of Toolbx images that are used by Toolbx containers")

	rmiCmd.SetHelpFunc(rmiHelp)
	rootCmd.AddCommand(rmiCmd)
}

func rmi(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	if rmiFlags.deleteAll {
		toolboxImages, err := getImages(false)
		if err != nil {
			return err
		}

		for _, image := range toolboxImages {
			imageID := image.ID()
			if err := podman.RemoveImage(imageID, rmiFlags.forceDelete); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}
		}
	} else {
		if len(args) == 0 {
			var builder strings.Builder
			fmt.Fprintf(&builder, "missing argument for \"rmi\"\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}

		for _, image := range args {
			imageObj, err := podman.InspectImage(image)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to inspect image %s\n", image)
				continue
			}

			if !imageObj.IsToolbx() {
				fmt.Fprintf(os.Stderr, "Error: %s is not a Toolbx image\n", image)
				continue
			}

			if err := podman.RemoveImage(image, rmiFlags.forceDelete); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}
		}
	}

	return nil
}

func rmiHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-rmi"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

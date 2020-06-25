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
	"strings"

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rmiFlags struct {
		deleteAll   bool
		forceDelete bool
	}
)

var rmiCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Remove one or more toolbox images",
	RunE:  rmi,
}

func init() {
	flags := rmiCmd.Flags()

	flags.BoolVarP(&rmiFlags.deleteAll, "all", "a", false, "Remove all toolbox containers.")

	flags.BoolVarP(&rmiFlags.forceDelete,
		"force",
		"f",
		false,
		"Force the removal of running and paused toolbox containers.")

	rmiCmd.SetHelpFunc(rmiHelp)
	rootCmd.AddCommand(rmiCmd)
}

func rmi(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a toolbox container")
		}

		if _, err := utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	if rmiFlags.deleteAll {
		logrus.Debug("Fetching images with label=com.redhat.component=fedora-toolbox")
		args := []string{"--filter", "label=com.redhat.component=fedora-toolbox"}
		images_old, err := podman.GetImages(args...)
		if err != nil {
			return errors.New("failed to list images with com.redhat.component=fedora-toolbox")
		}

		logrus.Debug("Fetching images with label=com.github.debarshiray.toolbox=true")
		args = []string{"--filter", "label=com.github.debarshiray.toolbox=true"}
		images_new, err := podman.GetImages(args...)
		if err != nil {
			return errors.New("failed to list images with com.github.debarshiray.toolbox=true")
		}

		var idKey string
		if podman.CheckVersion("2.0.0") {
			idKey = "Id"
		} else if podman.CheckVersion("1.8.3") {
			idKey = "ID"
		} else {
			idKey = "id"
		}

		images := utils.JoinJSON(idKey, images_old, images_new)

		for _, image := range images {
			imageID := image[idKey].(string)
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
			if exists, err := podman.ImageExists(image); !exists {
				if errors.Is(err, podman.ErrImageNotExist) {
					fmt.Fprintf(os.Stderr, "Error: image %s does not exist\n", image)
				} else {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				}
				continue
			}

			if _, err := podman.IsToolboxImage(image); err != nil {
				if errors.Is(err, podman.ErrImageNotToolbox) {
					fmt.Fprintf(os.Stderr, "Error: image %s is not a toolbox image\n", image)
				} else {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				}
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
			fmt.Fprintf(os.Stderr, "Error: this is not a toolbox container\n")
			return
		}

		if _, err := utils.ForwardToHost(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return
		}

		return
	}

	if err := utils.ShowManual("toolbox-rmi"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

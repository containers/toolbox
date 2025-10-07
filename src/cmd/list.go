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
	"text/tabwriter"

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/term"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	listFlags struct {
		onlyContainers bool
		onlyImages     bool
	}
)

var listCmd = &cobra.Command{
	Use:               "list",
	Short:             "List existing Toolbx containers and images",
	RunE:              list,
	ValidArgsFunction: completionEmpty,
}

func init() {
	flags := listCmd.Flags()

	flags.BoolVarP(&listFlags.onlyContainers,
		"containers",
		"c",
		false,
		"List only Toolbx containers, not images")

	flags.BoolVarP(&listFlags.onlyImages,
		"images",
		"i",
		false,
		"List only Toolbx images, not containers")

	listCmd.SetHelpFunc(listHelp)
	rootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	lsContainers := true
	lsImages := true

	if !listFlags.onlyContainers && listFlags.onlyImages {
		lsContainers = false
	} else if listFlags.onlyContainers && !listFlags.onlyImages {
		lsImages = false
	}

	var images []podman.Image
	var containers []podman.Container
	var err error

	if lsImages {
		images, err = getImages(false)
		if err != nil {
			return err
		}
	}

	if lsContainers {
		containers, err = getContainers()
		if err != nil {
			return err
		}
	}

	listOutput(images, containers)
	return nil
}

func getContainers() ([]podman.Container, error) {
	logrus.Debug("Fetching all containers")
	args := []string{"--all", "--sort", "names"}
	containers, err := podman.GetContainers(args...)
	if err != nil {
		logrus.Debugf("Fetching all containers failed: %s", err)
		return nil, errors.New("failed to get containers")
	}

	var toolboxContainers []podman.Container

	for containers.Next() {
		if container := containers.Get(); container.IsToolbx() {
			toolboxContainers = append(toolboxContainers, container)
		}
	}

	return toolboxContainers, nil
}

func listHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-list"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func getImages(fillNameWithID bool) ([]podman.Image, error) {
	logrus.Debug("Fetching all images")
	var args []string
	images, err := podman.GetImages(fillNameWithID, true, args...)
	if err != nil {
		logrus.Debugf("Fetching all images failed: %s", err)
		return nil, errors.New("failed to get images")
	}

	var toolboxImages []podman.Image

	for images.Next() {
		if image := images.Get(); image.IsToolbx() {
			toolboxImages = append(toolboxImages, image)
		}
	}

	return toolboxImages, nil
}

func listOutput(images []podman.Image, containers []podman.Container) {
	if len(images) != 0 {
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(writer, "%s\t%s\t%s\n", "IMAGE ID", "IMAGE NAME", "CREATED")

		for _, image := range images {
			if len(image.Names()) != 1 {
				panic("cannot list unflattened Image")
			}

			fmt.Fprintf(writer, "%s\t%s\t%s\n",
				utils.ShortID(image.ID()),
				image.Names()[0],
				image.Created())
		}

		writer.Flush()
	}

	if len(images) != 0 && len(containers) != 0 {
		fmt.Println()
	}

	if len(containers) != 0 {
		const boldGreenColor = "\033[1;32m"
		const defaultColor = "\033[0;00m" // identical to resetColor, but same length as boldGreenColor
		const resetColor = "\033[0m"

		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		if term.IsTerminal(os.Stdout) {
			fmt.Fprintf(writer, "%s", defaultColor)
		}

		fmt.Fprintf(writer,
			"%s\t%s\t%s\t%s\t%s",
			"CONTAINER ID",
			"CONTAINER NAME",
			"CREATED",
			"STATUS",
			"IMAGE NAME")

		if term.IsTerminal(os.Stdout) {
			fmt.Fprintf(writer, "%s", resetColor)
		}

		fmt.Fprintf(writer, "\n")

		for _, container := range containers {
			isRunning := false
			if podman.CheckVersion("2.0.0") {
				status := container.Status()
				isRunning = status == "running"
			}

			if term.IsTerminal(os.Stdout) {
				var color string
				if isRunning {
					color = boldGreenColor
				} else {
					color = defaultColor
				}

				fmt.Fprintf(writer, "%s", color)
			}

			created := container.Created()
			id := container.ID()
			image := container.Image()
			name := container.Name()
			status := container.Status()
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s", utils.ShortID(id), name, created, status, image)

			if term.IsTerminal(os.Stdout) {
				fmt.Fprintf(writer, "%s", resetColor)
			}

			fmt.Fprintf(writer, "\n")
		}

		writer.Flush()
	}
}

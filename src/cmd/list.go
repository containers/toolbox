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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type toolboxImage struct {
	ID      string
	Names   []string
	Created string
	Labels  map[string]string
}

type toolboxContainer struct {
	ID      string
	Names   []string
	Status  string
	Created string
	Image   string
	Labels  map[string]string
}

var (
	listFlags struct {
		onlyContainers bool
		onlyImages     bool
	}

	// toolboxLabels holds labels used by containers/images that mark them as compatible with Toolbox
	toolboxLabels = map[string]string{
		"com.github.debarshiray.toolbox": "true",
		"com.github.containers.toolbox":  "true",
	}
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List existing toolbox containers and images",
	RunE:  list,
}

func init() {
	flags := listCmd.Flags()

	flags.BoolVarP(&listFlags.onlyContainers,
		"containers",
		"c",
		false,
		"List only toolbox containers, not images")

	flags.BoolVarP(&listFlags.onlyImages,
		"images",
		"i",
		false,
		"List only toolbox images, not containers")

	listCmd.SetHelpFunc(listHelp)
	rootCmd.AddCommand(listCmd)
}

func list(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a toolbox container")
		}

		if _, err := utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	lsContainers := true
	lsImages := true

	if !listFlags.onlyContainers && listFlags.onlyImages {
		lsContainers = false
	} else if listFlags.onlyContainers && !listFlags.onlyImages {
		lsImages = false
	}

	var images []toolboxImage
	var containers []toolboxContainer
	var err error

	if lsImages {
		images, err = getImages()
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

func getContainers() ([]toolboxContainer, error) {
	logrus.Debug("Fetching all containers")
	args := []string{"--all", "--sort", "names"}
	containers, err := podman.GetContainers(args...)
	if err != nil {
		logrus.Debugf("Fetching all containers failed: %s", err)
		return nil, errors.New("failed to get containers")
	}

	var toolboxContainers []toolboxContainer

	for _, container := range containers {
		var c toolboxContainer
		var isToolboxContainer bool = false

		containerJSON, err := json.Marshal(container)
		if err != nil {
			logrus.Errorf("failed to marshal container: %v", err)
			continue
		}

		err = c.UnmarshalJSON(containerJSON)
		if err != nil {
			logrus.Errorf("failed to unmarshal container: %v", err)
			continue
		}

		for label := range toolboxLabels {
			if _, ok := c.Labels[label]; ok {
				isToolboxContainer = true
				break
			}
		}

		if isToolboxContainer {
			toolboxContainers = append(toolboxContainers, c)
		}
	}

	return toolboxContainers, nil
}

func listHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-list"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func getImages() ([]toolboxImage, error) {
	logrus.Debug("Fetching all images")
	args := []string{"--sort", "repository"}
	images, err := podman.GetImages(args...)
	if err != nil {
		logrus.Debugf("Fetching all images failed: %s", err)
		return nil, errors.New("failed to get images")
	}

	var toolboxImages []toolboxImage

	for _, image := range images {
		var i toolboxImage
		var isToolboxImage bool = false

		imageJSON, err := json.Marshal(image)
		if err != nil {
			logrus.Errorf("failed to marshal toolbox image: %v", err)
			continue
		}

		err = i.UnmarshalJSON(imageJSON)
		if err != nil {
			logrus.Errorf("failed to unmarshal toolbox image: %v", err)
			continue
		}

		for label := range toolboxLabels {
			if _, ok := i.Labels[label]; ok {
				isToolboxImage = true
				break
			}
		}

		if isToolboxImage {
			toolboxImages = append(toolboxImages, i)
		}
	}

	return toolboxImages, nil
}

func listOutput(images []toolboxImage, containers []toolboxContainer) {
	if len(images) != 0 {
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(writer, "%s\t%s\t%s\n", "IMAGE ID", "IMAGE NAME", "CREATED")

		for _, image := range images {
			imageName := "<none>"
			if len(image.Names) != 0 {
				imageName = image.Names[0]
			}

			fmt.Fprintf(writer, "%s\t%s\t%s\n",
				utils.ShortID(image.ID),
				imageName,
				image.Created)
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

		stdoutFd := os.Stdout.Fd()
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		if isatty.IsTerminal(stdoutFd) {
			fmt.Fprintf(writer, "%s", defaultColor)
		}

		fmt.Fprintf(writer,
			"%s\t%s\t%s\t%s\t%s",
			"CONTAINER ID",
			"CONTAINER NAME",
			"CREATED",
			"STATUS",
			"IMAGE NAME")

		if isatty.IsTerminal(stdoutFd) {
			fmt.Fprintf(writer, "%s", resetColor)
		}

		fmt.Fprintf(writer, "\n")

		for _, container := range containers {
			isRunning := false
			if podman.CheckVersion("2.0.0") {
				isRunning = container.Status == "running"
			}

			if isatty.IsTerminal(stdoutFd) {
				var color string
				if isRunning {
					color = boldGreenColor
				} else {
					color = defaultColor
				}

				fmt.Fprintf(writer, "%s", color)
			}

			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s",
				utils.ShortID(container.ID),
				container.Names[0],
				container.Created,
				container.Status,
				container.Image)

			if isatty.IsTerminal(stdoutFd) {
				fmt.Fprintf(writer, "%s", resetColor)
			}

			fmt.Fprintf(writer, "\n")
		}

		writer.Flush()
	}
}

func (i *toolboxImage) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID      string
		Names   []string
		Created interface{}
		Labels  map[string]string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	i.ID = raw.ID
	i.Names = raw.Names
	// Until Podman 2.0.x the field 'Created' held a human-readable string in
	// format "5 minutes ago". Since Podman 2.1 the field holds an integer with
	// Unix time. Go interprets numbers in JSON as float64.
	switch value := raw.Created.(type) {
	case string:
		i.Created = value
	case float64:
		i.Created = utils.HumanDuration(int64(value))
	}

	i.Labels = raw.Labels

	return nil
}

func (c *toolboxContainer) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID      string
		Names   interface{}
		Status  string
		State   interface{}
		Created interface{}
		Image   string
		Labels  map[string]string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	c.ID = raw.ID
	// In Podman V1 the field 'Names' held a single string but since Podman V2 the
	// field holds an array of strings
	switch value := raw.Names.(type) {
	case string:
		c.Names = append(c.Names, value)
	case []interface{}:
		for _, v := range value {
			c.Names = append(c.Names, v.(string))
		}
	}

	// In Podman V1 the field holding a string about the container's state was
	// called 'Status' and field 'State' held a number representing the state. In
	// Podman V2 the string was moved to 'State' and field 'Status' was dropped.
	switch value := raw.State.(type) {
	case string:
		c.Status = value
	case float64:
		c.Status = raw.Status
	}

	// In Podman V1 the field 'Created' held a human-readable string in format
	// "5 minutes ago". Since Podman V2 the field holds an integer with Unix time.
	// After a discussion in https://github.com/containers/podman/issues/6594 the
	// previous value was moved to field 'CreatedAt'. Since we're already using
	// the 'github.com/docker/go-units' library, we'll stop using the provided
	// human-readable string and assemble it ourselves. Go interprets numbers in
	// JSON as float64.
	switch value := raw.Created.(type) {
	case string:
		c.Created = value
	case float64:
		c.Created = utils.HumanDuration(int64(value))
	}
	c.Image = raw.Image
	c.Labels = raw.Labels

	return nil
}

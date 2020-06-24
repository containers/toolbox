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
	"text/tabwriter"

	"github.com/containers/toolbox/pkg/podman"
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
		"List only toolbox containers, not images.")

	flags.BoolVarP(&listFlags.onlyImages,
		"images",
		"i",
		false,
		"List only toolbox images, not containers.")

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

	var images []map[string]interface{}
	var containers []map[string]interface{}
	var err error

	if lsImages {
		images, err = listImages()
		if err != nil {
			return err
		}
	}

	if lsContainers {
		containers, err = listContainers()
		if err != nil {
			return err
		}
	}

	listOutput(images, containers)
	return nil
}

func listContainers() ([]map[string]interface{}, error) {
	logrus.Debug("Fetching containers with label=com.redhat.component=fedora-toolbox")
	args := []string{"--all", "--filter", "label=com.redhat.component=fedora-toolbox"}
	containers_old, err := podman.GetContainers(args...)
	if err != nil {
		return nil, errors.New("failed to list containers with com.redhat.component=fedora-toolbox")
	}

	logrus.Debug("Fetching containers with label=com.github.debarshiray.toolbox=true")
	args = []string{"--all", "--filter", "label=com.github.debarshiray.toolbox=true"}
	containers_new, err := podman.GetContainers(args...)
	if err != nil {
		return nil, errors.New("failed to list containers with label=com.github.debarshiray.toolbox=true")
	}

	containers := utils.JoinJSON("ID", containers_old, containers_new)
	containers = utils.SortJSON(containers, "Names", false)
	return containers, nil
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

	if err := utils.ShowManual("toolbox-list"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func listImages() ([]map[string]interface{}, error) {
	logrus.Debug("Fetching images with label=com.redhat.component=fedora-toolbox")
	args := []string{"--filter", "label=com.redhat.component=fedora-toolbox"}
	images_old, err := podman.GetImages(args...)
	if err != nil {
		return nil, errors.New("failed to list images with com.redhat.component=fedora-toolbox")
	}

	logrus.Debug("Fetching images with label=com.github.debarshiray.toolbox=true")
	args = []string{"--filter", "label=com.github.debarshiray.toolbox=true"}
	images_new, err := podman.GetImages(args...)
	if err != nil {
		return nil, errors.New("failed to list images with com.github.debarshiray.toolbox=true")
	}

	var images []map[string]interface{}
	if podman.CheckVersion("2.0.0") {
		images = utils.JoinJSON("Id", images_old, images_new)
		images = utils.SortJSON(images, "Names", true)
	} else if podman.CheckVersion("1.8.3") {
		images = utils.JoinJSON("ID", images_old, images_new)
		images = utils.SortJSON(images, "Names", true)
	} else {
		images = utils.JoinJSON("id", images_old, images_new)
		images = utils.SortJSON(images, "names", true)
	}

	return images, nil
}

func listOutput(images, containers []map[string]interface{}) {
	if len(images) != 0 {
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(writer, "%s\t%s\t%s\n", "IMAGE ID", "IMAGE NAME", "CREATED")

		var idKey, nameKey, createdKey string
		if podman.CheckVersion("2.0.0") {
			idKey = "Id"
			nameKey = "Names"
			createdKey = "Created"
		} else if podman.CheckVersion("1.8.3") {
			idKey = "ID"
			nameKey = "Names"
			createdKey = "Created"
		} else {
			idKey = "id"
			nameKey = "names"
			createdKey = "created"
		}
		for _, image := range images {
			id := utils.ShortID(image[idKey].(string))
			name := image[nameKey].([]interface{})[0].(string)
			created := image[createdKey].(string)
			fmt.Fprintf(writer, "%s\t%s\t%s\n", id, name, created)
		}

		writer.Flush()
	}

	if len(images) != 0 && len(containers) != 0 {
		fmt.Println()
	}

	if len(containers) != 0 {
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(writer,
			"%s\t%s\t%s\t%s\t%s\n",
			"CONTAINER ID",
			"CONTAINER NAME",
			"CREATED",
			"STATUS",
			"IMAGE NAME")

		for _, container := range containers {
			id := utils.ShortID(container["ID"].(string))
			name := container["Names"].(string)
			created := container["Created"].(string)
			status := container["Status"].(string)
			imageName := container["Image"].(string)
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", id, name, created, status, imageName)
		}

		writer.Flush()
	}
}

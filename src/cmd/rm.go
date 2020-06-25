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
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rmFlags struct {
		deleteAll   bool
		forceDelete bool
	}
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove one or more toolbox containers",
	RunE:  rm,
}

func init() {
	flags := rmCmd.Flags()

	flags.BoolVarP(&rmFlags.deleteAll, "all", "a", false, "Remove all toolbox containers.")

	flags.BoolVarP(&rmFlags.forceDelete,
		"force",
		"f",
		false,
		"Force the removal of running and paused toolbox containers.")

	rmCmd.SetHelpFunc(rmHelp)
	rootCmd.AddCommand(rmCmd)
}

func rm(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a toolbox container")
		}

		if _, err := utils.ForwardToHost(); err != nil {
			return err
		}

		return nil
	}

	if rmFlags.deleteAll {
		logrus.Debug("Fetching containers with label=com.redhat.component=fedora-toolbox")
		args := []string{"--all", "--filter", "label=com.redhat.component=fedora-toolbox"}
		containers_old, err := podman.GetContainers(args...)
		if err != nil {
			return errors.New("failed to list containers with com.redhat.component=fedora-toolbox")
		}

		logrus.Debug("Fetching containers with label=com.github.debarshiray.toolbox=true")
		args = []string{"--all", "--filter", "label=com.github.debarshiray.toolbox=true"}
		containers_new, err := podman.GetContainers(args...)
		if err != nil {
			return errors.New("failed to list containers with com.github.debarshiray.toolbox=true")
		}

		var idKey string
		if podman.CheckVersion("2.0.0") {
			idKey = "Id"
		} else {
			idKey = "ID"
		}

		containers := utils.JoinJSON(idKey, containers_old, containers_new)

		for _, container := range containers {
			containerID := container[idKey].(string)
			if err := removeContainer(containerID); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}
		}
	} else {
		if len(args) == 0 {
			var builder strings.Builder
			fmt.Fprintf(&builder, "missing argument for \"rm\"\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}

		for _, container := range args {
			if _, err := podman.IsToolboxContainer(container); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}

			if err := removeContainer(container); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
				continue
			}
		}
	}

	return nil
}

func rmHelp(cmd *cobra.Command, args []string) {
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

	if err := utils.ShowManual("toolbox-rm"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func removeContainer(container string) error {
	logrus.Debugf("Removing container %s", container)

	logLevelString := podman.LogLevel.String()
	args := []string{"--log-level", logLevelString, "rm"}

	if rmFlags.forceDelete {
		args = append(args, "--force")
	}

	args = append(args, container)

	exitCode, err := shell.RunWithExitCode("podman", nil, nil, nil, args...)
	switch exitCode {
	case 0:
		if err != nil {
			panic("unexpected error: 'podman rm' finished successfully")
		}
	case 1:
		err = fmt.Errorf("container %s does not exist", container)
	case 2:
		err = fmt.Errorf("container %s is running", container)
	default:
		err = fmt.Errorf("failed to remove container %s", container)
	}

	if err != nil {
		return err
	}

	return nil
}

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

	"github.com/containers/toolbox/pkg/podman"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	runFlags struct {
		container string
		release   string
	}
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command in an existing toolbox container",
	RunE:  run,
}

func init() {
	flags := runCmd.Flags()
	flags.SetInterspersed(false)

	flags.StringVarP(&runFlags.container,
		"container",
		"c",
		"",
		"Run command inside a toolbox container with the given name.")

	flags.StringVarP(&runFlags.release,
		"release",
		"r",
		"",
		"Run command inside a toolbox container for a different operating system release than the host.")

	runCmd.SetHelpFunc(runHelp)
	rootCmd.AddCommand(runCmd)
}

func run(cmd *cobra.Command, args []string) error {
	return nil
}

func runHelp(cmd *cobra.Command, args []string) {
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

	if err := utils.ShowManual("toolbox-run"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func getEntryPointAndPID(container string) (string, int, error) {
	logrus.Debugf("Inspecting entry point of container %s", container)

	info, err := podman.Inspect("container", container)
	if err != nil {
		return "", 0, fmt.Errorf("failed to inspect entry point of container %s", container)
	}

	config := info["Config"].(map[string]interface{})
	entryPoint := config["Cmd"].([]interface{})[0].(string)

	state := info["State"].(map[string]interface{})
	entryPointPID := state["Pid"]
	logrus.Debugf("Entry point PID is a %T", entryPointPID)

	var entryPointPIDInt int

	switch entryPointPID.(type) {
	case float64:
		entryPointPIDFloat := entryPointPID.(float64)
		entryPointPIDInt = int(entryPointPIDFloat)
	case int:
		entryPointPIDInt = entryPointPID.(int)
	default:
		return "", 0, fmt.Errorf("failed to inspect entry point PID of container %s", container)
	}

	logrus.Debugf("Entry point of container %s is %s (PID=%d)", container, entryPoint, entryPointPIDInt)

	return entryPoint, entryPointPIDInt, nil
}

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

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	createFlags struct {
		container string
		image     string
		release   string
	}
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new toolbox container",
	RunE:  create,
}

func init() {
	flags := createCmd.Flags()

	flags.StringVarP(&createFlags.container,
		"container",
		"c",
		"",
		"Assign a different name to the toolbox container.")

	flags.StringVarP(&createFlags.image,
		"image",
		"i",
		"",
		"Change the name of the base image used to create the toolbox container.")

	flags.StringVarP(&createFlags.release,
		"release",
		"r",
		"",
		"Create a toolbox container for a different operating system release than the host.")

	createCmd.SetHelpFunc(createHelp)
	rootCmd.AddCommand(createCmd)
}

func create(cmd *cobra.Command, args []string) error {
	return nil
}

func createHelp(cmd *cobra.Command, args []string) {
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

	if err := utils.ShowManual("toolbox-create"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

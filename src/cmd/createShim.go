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

/*
	import (
		"errors"
		"fmt"
		"os"
		"strings"

		"github.com/containers/toolbox/pkg/shell"
		"github.com/containers/toolbox/pkg/utils"
		"github.com/sirupsen/logrus"
		"github.com/spf13/cobra"
	)

	var (
		forbiddenBinaryShims = []string{
			"toolbox",
		}
	)

	var createShimCmd = &cobra.Command{
		Use:     "create-shim",
		Short:   "Create a binary shim for a command to be run on the host",
		Hidden:  true,
		Example: "asdasdasd",
		Args:    cobra.ExactArgs(1),
		RunE:    createShim,
	}

	func init() {
		// flags := initContainerCmd.Flags()

		createShimCmd.SetHelpFunc(createShimHelp)
		rootCmd.AddCommand(createShimCmd)
	}

	func createShim(cmd *cobra.Command, args []string) error {
		if !utils.IsInsideContainer() {
			var builder strings.Builder
			fmt.Fprintf(&builder, "the 'create-shim' command can only be used inside containers\n")
			fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

			errMsg := builder.String()
			return errors.New(errMsg)
		}

		shimBinary := fmt.Sprintf("/usr/libexec/toolbox/%s", args[0])

		if !utils.PathExists(hostRunnerShim.containerPath) || !utils.PathExists(sudoShim.containerPath) {
			var builder strings.Builder
			fmt.Fprintf(&builder, "This toolbox container is not set up for creating binary shims\n")
			fmt.Fprintf(&builder, "You're possibly trying to use a newer Toolbox in a container created with an older version\n")
			fmt.Fprintf(&builder, "Try to update the Toolbox binary used by this container")

			errMsg := builder.String()
			return errors.New(errMsg)
		}

		if utils.PathExists(shimBinary) {
			fmt.Printf("The requested shim binary already exists.\n")
			return nil
		}

		err := redirectPath(shimBinary, hostRunnerShim.containerPath, false)
		if err != nil {
			return fmt.Errorf("Failed to create shim binary %s: %w", shimBinary, err)
		}

		return nil
	}

	func createShimHelp(cmd *cobra.Command, args []string) {
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

		if err := utils.ShowManual("toolbox-create-shim"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return
		}
	}

	func runOnHost(command string, commandLineArgs []string) (int, error) {
		envOptions := utils.GetEnvOptionsForPreservedVariables()

		var flatpakSpawnArgs []string

		flatpakSpawnArgs = append(flatpakSpawnArgs, envOptions...)

		flatpakSpawnArgs = append(flatpakSpawnArgs, []string{
			"--host",
			command,
		}...)

		flatpakSpawnArgs = append(flatpakSpawnArgs, commandLineArgs...)

		logrus.Debug("Forwarding to host:")
		logrus.Debugf("%s", command)
		for _, arg := range commandLineArgs {
			logrus.Debugf("%s", arg)
		}

		exitCode, err := shell.RunWithExitCode("flatpak-spawn", os.Stdin, os.Stdout, nil, flatpakSpawnArgs...)
		if err != nil {
			return exitCode, err
		}

		return exitCode, nil
	}
*/

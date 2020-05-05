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
	"os/exec"
	"strings"
	"syscall"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove all local podman (and toolbox) state",
	RunE:  reset,
}

func init() {
	resetCmd.SetHelpFunc(resetHelp)
	rootCmd.AddCommand(resetCmd)
}

func reset(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(os.Stderr, "'%s reset' is deprecated in favor of 'podman system reset'.\n", executableBase)

	if utils.IsInsideContainer() {
		var builder strings.Builder
		fmt.Fprintf(&builder, "the 'reset' command cannot be used inside containers\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	podmanBinary, err := exec.LookPath("podman")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("podman(1) not found")
		}

		return errors.New("failed to lookup podman(1)")
	}

	podmanArgs := []string{"podman", "system", "reset"}
	if rootFlags.assumeYes {
		podmanArgs = append(podmanArgs, []string{"--force"}...)
	}

	env := os.Environ()

	if err := syscall.Exec(podmanBinary, podmanArgs, env); err != nil {
		return errors.New("failed to invoke podman(1)")
	}

	return nil
}

func resetHelp(cmd *cobra.Command, args []string) {
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

	if err := utils.ShowManual("toolbox-reset"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

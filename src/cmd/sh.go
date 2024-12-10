/*
 * Copyright © 2022 – 2024 Yann Soubeyrand
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
	"os/user"
	"strings"
	"syscall"

	"github.com/containers/toolbox/pkg/utils"
	lcuser "github.com/opencontainers/runc/libcontainer/user"
	"github.com/spf13/cobra"
)

var shCmd = &cobra.Command{
	Use:    "sh",
	Short:  "Launch user configured shell inside container",
	Hidden: true,
	RunE:   sh,
}

func init() {
	shCmd.SetHelpFunc(shHelp)
	rootCmd.AddCommand(shCmd)
}

func sh(cmd *cobra.Command, args []string) error {
	if !utils.IsInsideContainer() {
		var builder strings.Builder
		fmt.Fprintf(&builder, "the 'sh' command can only be used inside containers\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %s", err)
	}

	lcu, err := lcuser.LookupUser(u.Username)
	if err != nil {
		return fmt.Errorf("failed to lookup current user shell: %s", err)
	}

	for i, shell := range []string{
		lcu.Shell,
		"/bin/bash",
		"/bin/sh",
	} {
		if i > 0 {
			fmt.Fprintf(os.Stderr, "Falling back to %q\n", shell)
		}

		err = syscall.Exec(shell, args, os.Environ())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute %q: %s\n", shell, err)
		}
	}

	return fmt.Errorf("failed to execute a shell")
}

func shHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-sh"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

/*
 * Copyright © 2026 Hadi Chokr <hadichokr@icloud.com>
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
	"path/filepath"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	unexportFlags struct {
		all       bool
		app       string
		bin       string
		container string
	}
)

var unexportCmd = &cobra.Command{
	Use:               "unexport",
	Short:             "Remove an application or a binary exported from a Toolbx container",
	RunE:              unexport,
	ValidArgsFunction: completionEmpty,
}

func init() {
	flags := unexportCmd.Flags()

	flags.BoolVar(&unexportFlags.all, "all", false, "Remove everything exported from the Toolbx container")

	flags.StringVar(&unexportFlags.app, "app", "", "Remove the application with the given desktop entry")

	flags.StringVar(&unexportFlags.bin, "bin", "", "Remove the binary with the given name")

	flags.StringVarP(&unexportFlags.container,
		"container",
		"c",
		"",
		"Remove what was exported from the Toolbx container with the given name")

	if err := unexportCmd.RegisterFlagCompletionFunc("container", completionContainerNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}

	unexportCmd.SetHelpFunc(unexportHelp)
	rootCmd.AddCommand(unexportCmd)
}

func unexport(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	if unexportFlags.all && (cmd.Flag("app").Changed || cmd.Flag("bin").Changed) {
		var builder strings.Builder
		fmt.Fprintf(&builder, "option --all cannot be used with --app or --bin\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if !unexportFlags.all && unexportFlags.app == "" && unexportFlags.bin == "" {
		var builder strings.Builder
		fmt.Fprintf(&builder, "option --all, --app or --bin is needed\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	container, _, _, err := resolveContainerAndImageNames(unexportFlags.container, "--container", "", "", "")
	if err != nil {
		return err
	}

	if unexportFlags.all {
		removed, err := unexportAll(container)
		if err != nil {
			return err
		}

		fmt.Printf("Removed %d files exported from container %s\n", removed, container)
		return nil
	}

	if unexportFlags.app != "" {
		return unexportApp(container, unexportFlags.app)
	}

	return unexportBin(container, unexportFlags.bin)
}

func unexportHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-unexport"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func unexportApp(container, app string) error {
	appID := strings.TrimSuffix(app, ".desktop")
	target := filepath.Join(exportApplicationsDir(), appID+".desktop")

	if err := exportRemove(target, container); err != nil {
		return err
	}

	icons := filepath.Join(exportDataHome(), "toolbx", container, "icons", appID)
	if err := os.RemoveAll(icons); err != nil {
		return fmt.Errorf("failed to remove %s: %w", icons, err)
	}

	fmt.Printf("Removed %s exported from container %s\n", appID+".desktop", container)
	return nil
}

func unexportBin(container, bin string) error {
	target := filepath.Join(exportBinHome(), filepath.Base(bin))

	if err := exportRemove(target, container); err != nil {
		return err
	}

	fmt.Printf("Removed %s exported from container %s\n", filepath.Base(target), container)
	return nil
}

func unexportAll(container string) (int, error) {
	var removed int

	for _, directory := range []string{exportApplicationsDir(), exportBinHome()} {
		entries, err := os.ReadDir(directory)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return removed, fmt.Errorf("failed to read %s: %w", directory, err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			path := filepath.Join(directory, entry.Name())
			owner, err := exportOwner(path)
			if err != nil || owner != container {
				continue
			}

			if err := os.Remove(path); err != nil {
				return removed, fmt.Errorf("failed to remove %s: %w", path, err)
			}

			removed++
		}
	}

	data := filepath.Join(exportDataHome(), "toolbx", container)
	if err := os.RemoveAll(data); err != nil {
		return removed, fmt.Errorf("failed to remove %s: %w", data, err)
	}

	return removed, nil
}

func exportRemove(path, container string) error {
	owner, err := exportOwner(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file %s not found", path)
		}

		return err
	}

	if owner != container {
		return fmt.Errorf("file %s was not exported from container %s", path, container)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}

	return nil
}

/*
 * Copyright © 2025 Hadi Chokr <hadichokr@icloud.com>
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
	unexportContainer string
	unexportBin       string
	unexportApp       string
	unexportAll       bool
)

var unexportCmd = &cobra.Command{
	Use:   "unexport",
	Short: "Remove exported binaries and applications for a specific toolbox container",
	RunE:  runUnexport,
}

func init() {
	unexportCmd.Flags().StringVar(&unexportContainer, "container", "", "Name of the toolbox container")
	unexportCmd.Flags().StringVar(&unexportBin, "bin", "", "Name of the exported binary to remove")
	unexportCmd.Flags().StringVar(&unexportApp, "app", "", "Name of the exported application to remove")
	unexportCmd.Flags().BoolVar(&unexportAll, "all", false, "Remove all exported binaries and applications for the container")
	unexportCmd.SetHelpFunc(unexportHelp)
	rootCmd.AddCommand(unexportCmd)
}

func runUnexport(cmd *cobra.Command, args []string) error {
	if unexportContainer == "" {
		return errors.New("must specify --container")
	}

	if !unexportAll && unexportBin == "" && unexportApp == "" {
		return errors.New("must specify --bin, --app, or --all")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	binDir := filepath.Join(homeDir, ".local", "bin")
	appsDir := filepath.Join(homeDir, ".local", "share", "applications")

	removedBins := []string{}
	removedApps := []string{}

	if unexportBin != "" {
		path := filepath.Join(binDir, unexportBin)
		if fileContainsContainer(path, unexportContainer) {
			if err := os.Remove(path); err == nil {
				removedBins = append(removedBins, path)
			}
		}
	}

	if unexportApp != "" {
		// Remove .desktop file that matches app name and container
		matches, _ := filepath.Glob(filepath.Join(appsDir, fmt.Sprintf("*%s-%s.desktop", unexportApp, unexportContainer)))
		for _, path := range matches {
			if err := os.Remove(path); err == nil {
				removedApps = append(removedApps, path)
			}
		}
	}

	if unexportAll {
		// Remove all binaries for this container in .local/bin
		binFiles, _ := os.ReadDir(binDir)
		for _, f := range binFiles {
			if f.IsDir() {
				continue
			}
			path := filepath.Join(binDir, f.Name())
			if fileContainsContainer(path, unexportContainer) {
				if err := os.Remove(path); err == nil {
					removedBins = append(removedBins, path)
				}
			}
		}

		// Remove all .desktop files for this container in .local/share/applications
		appFiles, _ := os.ReadDir(appsDir)
		for _, f := range appFiles {
			name := f.Name()
			if strings.HasSuffix(name, "-"+unexportContainer+".desktop") {
				path := filepath.Join(appsDir, name)
				if err := os.Remove(path); err == nil {
					removedApps = append(removedApps, path)
				}
			}
		}
	}

	fmt.Printf("Removed binaries:\n")
	for _, b := range removedBins {
		fmt.Printf("  %s\n", b)
	}
	fmt.Printf("Removed desktop files:\n")
	for _, a := range removedApps {
		fmt.Printf("  %s\n", a)
	}
	if len(removedBins) == 0 && len(removedApps) == 0 {
		fmt.Println("No exported binaries or desktop files found to remove for container", unexportContainer)
	}
	return nil
}

// fileContainsContainer returns true if the file exists and has a toolbox_binary comment with name: <container>
func fileContainsContainer(path, container string) bool {
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "# toolbox_binary") && strings.Contains(string(content), fmt.Sprintf("name: %s", container))
}

// Exported function: remove all exported binaries and desktop files for a container
func UnexportAll(container string) ([]string, []string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}
	binDir := filepath.Join(homeDir, ".local", "bin")
	appsDir := filepath.Join(homeDir, ".local", "share", "applications")

	removedBins := []string{}
	removedApps := []string{}

	binFiles, _ := os.ReadDir(binDir)
	for _, f := range binFiles {
		if f.IsDir() {
			continue
		}
		path := filepath.Join(binDir, f.Name())
		if fileContainsContainer(path, container) {
			if err := os.Remove(path); err == nil {
				removedBins = append(removedBins, path)
			}
		}
	}

	appFiles, _ := os.ReadDir(appsDir)
	for _, f := range appFiles {
		name := f.Name()
		if strings.HasSuffix(name, "-"+container+".desktop") {
			path := filepath.Join(appsDir, name)
			if err := os.Remove(path); err == nil {
				removedApps = append(removedApps, path)
			}
		}
	}

	return removedBins, removedApps, nil
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

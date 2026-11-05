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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/google/renameio/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const exportContainerKey = "X-Toolbx-Container"

type exportIcon struct {
	Path string `json:"path"`
	Data []byte `json:"data"`
}

type exportManifest struct {
	AppID   string       `json:"appId,omitempty"`
	Desktop string       `json:"desktop,omitempty"`
	Icons   []exportIcon `json:"icons,omitempty"`
	Binary  string       `json:"binary,omitempty"`
}

var (
	exportFlags struct {
		app       string
		bin       string
		container string
		force     bool
	}
)

var exportCmd = &cobra.Command{
	Use:               "export",
	Short:             "Export an application or a binary from a Toolbx container",
	RunE:              export,
	ValidArgsFunction: completionEmpty,
}

func init() {
	flags := exportCmd.Flags()

	flags.StringVar(&exportFlags.app, "app", "", "Export the application with the given desktop entry")

	flags.StringVar(&exportFlags.bin, "bin", "", "Export the binary with the given name")

	flags.StringVarP(&exportFlags.container,
		"container",
		"c",
		"",
		"Export from the Toolbx container with the given name")

	flags.BoolVar(&exportFlags.force,
		"force",
		false,
		"Overwrite files that weren't exported from the same container")

	if err := exportCmd.RegisterFlagCompletionFunc("container", completionContainerNames); err != nil {
		panicMsg := fmt.Sprintf("failed to register flag completion function: %v", err)
		panic(panicMsg)
	}

	exportCmd.SetHelpFunc(exportHelp)
	rootCmd.AddCommand(exportCmd)
}

func export(cmd *cobra.Command, args []string) error {
	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	if cmd.Flag("app").Changed && cmd.Flag("bin").Changed {
		var builder strings.Builder
		fmt.Fprintf(&builder, "options --app and --bin cannot be used together\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	if exportFlags.app == "" && exportFlags.bin == "" {
		var builder strings.Builder
		fmt.Fprintf(&builder, "option --app or --bin is needed\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	container, image, release, err := resolveContainerAndImageNames(exportFlags.container,
		"--container",
		"",
		"",
		"")

	if err != nil {
		return err
	}

	manifest, err := inspectExport(container, image, release)
	if err != nil {
		return err
	}

	if exportFlags.app != "" {
		return exportApp(container, manifest)
	}

	return exportBin(container, manifest)
}

func exportHelp(cmd *cobra.Command, args []string) {
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

	if err := showManual("toolbox-export"); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
}

func inspectExport(container, image, release string) (*exportManifest, error) {
	runtimeDirectory, err := utils.GetRuntimeDirectory(currentUser)
	if err != nil {
		return nil, err
	}

	manifestPath := filepath.Join(runtimeDirectory, fmt.Sprintf("export-%d.json", os.Getpid()))
	defer os.Remove(manifestPath)

	command := []string{"toolbox", "export-inspect", "--output", manifestPath}
	if exportFlags.app != "" {
		command = append(command, "--app", exportFlags.app)
	} else {
		command = append(command, "--bin", exportFlags.bin)
	}

	if err := runCommand(container, false, image, release, 0, command, false, false, true); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		logrus.Debugf("Reading %s failed: %s", manifestPath, err)
		return nil, errors.New("failed to read the export manifest")
	}

	var manifest exportManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		logrus.Debugf("Unmarshalling %s failed: %s", manifestPath, err)
		return nil, errors.New("failed to parse the export manifest")
	}

	return &manifest, nil
}

func exportApp(container string, manifest *exportManifest) error {
	applications := exportApplicationsDir()
	if err := os.MkdirAll(applications, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", applications, err)
	}

	target := filepath.Join(applications, manifest.AppID+".desktop")
	if err := exportCheckTarget(target, container); err != nil {
		return err
	}

	icons := filepath.Join(exportDataHome(), "toolbx", container, "icons", manifest.AppID)
	if err := os.RemoveAll(icons); err != nil {
		return fmt.Errorf("failed to remove %s: %w", icons, err)
	}

	icon, err := exportWriteIcons(icons, manifest.Icons)
	if err != nil {
		return err
	}

	entry := exportRewriteEntry(container, manifest.Desktop, icon)
	if err := renameio.WriteFile(target, []byte(entry), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", target, err)
	}

	fmt.Printf("Exported %s from container %s\n", manifest.AppID+".desktop", container)
	return nil
}

func exportBin(container string, manifest *exportManifest) error {
	directory := exportBinHome()
	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("failed to create %s: %w", directory, err)
	}

	target := filepath.Join(directory, filepath.Base(manifest.Binary))
	if err := exportCheckTarget(target, container); err != nil {
		return err
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "#!/bin/sh\n")
	fmt.Fprintf(&builder, "# %s: %s\n", exportContainerKey, container)
	fmt.Fprintf(&builder, "exec %s run --container %s %s \"$@\"\n", executable, container, manifest.Binary)

	if err := renameio.WriteFile(target, []byte(builder.String()), 0755); err != nil {
		return fmt.Errorf("failed to write %s: %w", target, err)
	}

	fmt.Printf("Exported %s from container %s\n", filepath.Base(target), container)

	if !exportIsInPath(directory) {
		fmt.Fprintf(os.Stderr, "Warning: %s is not in PATH\n", directory)
	}

	return nil
}

func exportCheckTarget(path, container string) error {
	owner, err := exportOwner(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	if exportFlags.force || owner == container {
		return nil
	}

	var builder strings.Builder
	if owner == "" {
		fmt.Fprintf(&builder, "file %s was not exported by Toolbx\n", path)
	} else {
		fmt.Fprintf(&builder, "file %s was exported from container %s\n", path, owner)
	}
	fmt.Fprintf(&builder, "Use option '--force' to overwrite it.\n")
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func exportOwner(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if value, ok := strings.CutPrefix(line, exportContainerKey+"="); ok {
			return strings.TrimSpace(value), nil
		}

		if value, ok := strings.CutPrefix(line, "# "+exportContainerKey+":"); ok {
			return strings.TrimSpace(value), nil
		}
	}

	return "", nil
}

func exportRewriteEntry(container, entry, icon string) string {
	prefix := fmt.Sprintf("%s run --container %s ", executable, container)

	var builder strings.Builder
	var group string

	for _, line := range strings.Split(strings.TrimRight(entry, "\n"), "\n") {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			if group == "Desktop Entry" {
				fmt.Fprintf(&builder, "%s=%s\n", exportContainerKey, container)
			}

			group = strings.Trim(trimmed, "[]")
			fmt.Fprintf(&builder, "%s\n", line)
			continue
		}

		key, value, found := strings.Cut(trimmed, "=")
		if !found {
			fmt.Fprintf(&builder, "%s\n", line)
			continue
		}

		switch strings.TrimSpace(key) {
		case "Exec":
			fmt.Fprintf(&builder, "Exec=%s%s\n", prefix, strings.TrimSpace(value))
		case "TryExec":
			fmt.Fprintf(&builder, "TryExec=%s\n", executable)
		case "DBusActivatable":
			fmt.Fprintf(&builder, "DBusActivatable=false\n")
		case "Path":
			fmt.Fprintf(&builder, "#Path=%s\n", strings.TrimSpace(value))
		case "Icon":
			if icon == "" {
				fmt.Fprintf(&builder, "%s\n", line)
			} else {
				fmt.Fprintf(&builder, "Icon=%s\n", icon)
			}
		case exportContainerKey:
		default:
			fmt.Fprintf(&builder, "%s\n", line)
		}
	}

	if group == "Desktop Entry" {
		fmt.Fprintf(&builder, "%s=%s\n", exportContainerKey, container)
	}

	return builder.String()
}

func exportWriteIcons(directory string, icons []exportIcon) (string, error) {
	var best string
	var bestScore int

	for _, icon := range icons {
		path := filepath.Join(directory, filepath.FromSlash(icon.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return "", fmt.Errorf("failed to create %s: %w", filepath.Dir(path), err)
		}

		if err := os.WriteFile(path, icon.Data, 0644); err != nil {
			return "", fmt.Errorf("failed to write %s: %w", path, err)
		}

		if score := exportIconScore(icon.Path); best == "" || score > bestScore {
			best = path
			bestScore = score
		}
	}

	return best, nil
}

func exportIconScore(path string) int {
	if filepath.Ext(path) == ".svg" {
		return 1 << 16
	}

	for _, component := range strings.Split(path, "/") {
		size, _, found := strings.Cut(component, "x")
		if !found {
			continue
		}

		if value, err := strconv.Atoi(size); err == nil {
			return value
		}
	}

	return 0
}

func exportDataHome() string {
	if dataHome := os.Getenv("XDG_DATA_HOME"); dataHome != "" {
		return dataHome
	}

	return filepath.Join(getCurrentUserHomeDir(), ".local", "share")
}

func exportApplicationsDir() string {
	return filepath.Join(exportDataHome(), "applications")
}

func exportBinHome() string {
	if binHome := os.Getenv("XDG_BIN_HOME"); binHome != "" {
		return binHome
	}

	return filepath.Join(getCurrentUserHomeDir(), ".local", "bin")
}

func exportIsInPath(directory string) bool {
	for _, component := range filepath.SplitList(os.Getenv("PATH")) {
		if component == directory {
			return true
		}
	}

	return false
}

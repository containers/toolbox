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
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/google/renameio/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var exportInspectIconExtensions = []string{".png", ".svg", ".xpm"}

var exportInspectFlags struct {
	app    string
	bin    string
	output string
}

var exportInspectCmd = &cobra.Command{
	Use:    "export-inspect",
	Short:  "Resolve an application or a binary inside a Toolbx container",
	Hidden: true,
	RunE:   exportInspect,
}

func init() {
	flags := exportInspectCmd.Flags()

	flags.StringVar(&exportInspectFlags.app, "app", "", "Resolve the application with the given desktop entry")

	flags.StringVar(&exportInspectFlags.bin, "bin", "", "Resolve the binary with the given name")

	flags.StringVar(&exportInspectFlags.output, "output", "", "Write the manifest to the given path")
	if err := exportInspectCmd.MarkFlagRequired("output"); err != nil {
		panic("Could not mark flag --output as required")
	}

	rootCmd.AddCommand(exportInspectCmd)
}

func exportInspect(cmd *cobra.Command, args []string) error {
	if !utils.IsInsideContainer() {
		var builder strings.Builder
		fmt.Fprintf(&builder, "the 'export-inspect' command can only be used inside containers\n")
		fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	var manifest exportManifest
	var err error

	if exportInspectFlags.app != "" {
		manifest, err = exportInspectApp(exportInspectFlags.app)
	} else {
		manifest, err = exportInspectBin(exportInspectFlags.bin)
	}

	if err != nil {
		return err
	}

	data, err := json.Marshal(&manifest)
	if err != nil {
		logrus.Debugf("Marshalling the manifest failed: %s", err)
		return errors.New("failed to marshal the export manifest")
	}

	if err := renameio.WriteFile(exportInspectFlags.output, data, 0600); err != nil {
		return fmt.Errorf("failed to write %s: %w", exportInspectFlags.output, err)
	}

	return nil
}

func exportInspectApp(app string) (exportManifest, error) {
	var manifest exportManifest

	path, err := exportInspectFindEntry(app)
	if err != nil {
		return manifest, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return manifest, fmt.Errorf("failed to read %s: %w", path, err)
	}

	manifest.AppID = strings.TrimSuffix(filepath.Base(path), ".desktop")
	manifest.Desktop = string(data)

	icon := exportInspectEntryKey(manifest.Desktop, "Icon")
	if icon != "" {
		manifest.Icons = exportInspectFindIcons(icon)
	}

	return manifest, nil
}

func exportInspectBin(bin string) (exportManifest, error) {
	var manifest exportManifest

	path, err := exec.LookPath(bin)
	if err != nil {
		logrus.Debugf("Looking up %s failed: %s", bin, err)
		return manifest, fmt.Errorf("binary %s not found in the container", bin)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		return manifest, fmt.Errorf("failed to resolve %s: %w", bin, err)
	}

	manifest.Binary = path
	return manifest, nil
}

func exportInspectIsExported(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, exportContainerKey+"=") {
			logrus.Debugf("Skipping %s: exported by Toolbx", path)
			return true
		}
	}

	return false
}

func exportInspectFindEntry(app string) (string, error) {
	name := strings.TrimSuffix(app, ".desktop")
	var fallback string

	for _, directory := range exportInspectDataDirs("applications") {
		path := filepath.Join(directory, name+".desktop")
		if utils.PathExists(path) && !exportInspectIsExported(path) {
			return path, nil
		}

		entries, err := os.ReadDir(directory)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			id, found := strings.CutSuffix(entry.Name(), ".desktop")
			if !found {
				continue
			}

			if fallback != "" || !strings.HasSuffix(strings.ToLower(id), "."+strings.ToLower(name)) {
				continue
			}

			candidate := filepath.Join(directory, entry.Name())
			if exportInspectIsExported(candidate) {
				continue
			}

			fallback = candidate
		}
	}

	if fallback != "" {
		return fallback, nil
	}

	return "", fmt.Errorf("application %s not found in the container", app)
}

func exportInspectFindIcons(icon string) []exportIcon {
	if filepath.IsAbs(icon) {
		data, err := os.ReadFile(icon)
		if err != nil {
			logrus.Debugf("Reading %s failed: %s", icon, err)
			return nil
		}

		return []exportIcon{{Path: filepath.Base(icon), Data: data}}
	}

	var icons []exportIcon

	for _, root := range append(exportInspectDataDirs("icons"), "/usr/share/pixmaps") {
		err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() {
				return nil
			}

			extension := filepath.Ext(entry.Name())
			if strings.TrimSuffix(entry.Name(), extension) != icon {
				return nil
			}

			var known bool
			for _, candidate := range exportInspectIconExtensions {
				if extension == candidate {
					known = true
					break
				}
			}

			if !known {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			relative, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}

			icons = append(icons, exportIcon{Path: filepath.ToSlash(relative), Data: data})
			return nil
		})

		if err != nil {
			logrus.Debugf("Walking %s failed: %s", root, err)
		}
	}

	return icons
}

func exportInspectEntryKey(entry, key string) string {
	var group string

	for _, line := range strings.Split(entry, "\n") {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			group = strings.Trim(line, "[]")
			continue
		}

		if group != "Desktop Entry" {
			continue
		}

		if value, ok := strings.CutPrefix(line, key+"="); ok {
			return strings.TrimSpace(value)
		}
	}

	return ""
}

func exportInspectDataDirs(subdirectory string) []string {
	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		dataHome = filepath.Join(getCurrentUserHomeDir(), ".local", "share")
	}

	dataDirs := os.Getenv("XDG_DATA_DIRS")
	if dataDirs == "" {
		dataDirs = "/usr/local/share:/usr/share"
	}

	directories := []string{filepath.Join(dataHome, subdirectory)}
	for _, directory := range filepath.SplitList(dataDirs) {
		directories = append(directories, filepath.Join(directory, subdirectory))
	}

	return directories
}

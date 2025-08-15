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
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	exportBin       string
	exportApp       string
	exportContainer string
)

var exportCmd = &cobra.Command{
	Use:               "export",
	Short:             "Export binaries or applications from a toolbox container",
	RunE:              runExport,
	ValidArgsFunction: completionContainerNamesFiltered,
}

func init() {
	exportCmd.Flags().StringVar(&exportBin, "bin", "", "Path or name of binary to export")
	exportCmd.Flags().StringVar(&exportApp, "app", "", "Path or name of application to export")
	exportCmd.Flags().StringVar(&exportContainer, "container", "", "Name of the toolbox container")

	if err := exportCmd.RegisterFlagCompletionFunc("container", completionContainerNames); err != nil {
		panic(fmt.Sprintf("failed to register flag completion function: %v", err))
	}

	exportCmd.SetHelpFunc(exportHelp)
	rootCmd.AddCommand(exportCmd)
}

func runExport(cmd *cobra.Command, args []string) error {
	if exportBin == "" && exportApp == "" {
		return errors.New("must specify either --bin or --app")
	}
	if exportContainer == "" {
		return errors.New("must specify --container")
	}

	if exportBin != "" {
		return exportBinary(exportBin, exportContainer)
	}
	return exportApplication(exportApp, exportContainer)
}

func exportBinary(binName, containerName string) error {
	// Run command with strict environment
	out, err := runCommandWithOutput(
		containerName,
		false, "", "", 0,
		[]string{"sh", "--noprofile", "--norc", "-c",
			fmt.Sprintf("command -v %q 2>/dev/null || which %q 2>/dev/null || type -P %q 2>/dev/null || true",
				    binName, binName, binName)},
				  false, false, true,
	)
	if err != nil {
		return fmt.Errorf("failed to run command inside container: %v", err)
	}

	// Nuclear option for cleaning output
	cleanOutput := func(s string) string {
		// 1. Remove all ANSI escape sequences
		ansiRegex := regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
		s = ansiRegex.ReplaceAllString(s, "")

		// 2. Remove non-printable characters except newlines
		s = strings.Map(func(r rune) rune {
			if r >= 32 && r <= 126 || r == '\n' {
				return r
			}
			return -1
		}, s)

		// 3. Remove hyperlinks and terminal OSC sequences
		s = regexp.MustCompile(`\x1B\]8;;.*?\x1B\\`).ReplaceAllString(s, "")

		// 4. Trim each line and remove empty lines
		lines := strings.Split(s, "\n")
		var cleanLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				cleanLines = append(cleanLines, line)
			}
		}
		return strings.Join(cleanLines, "\n")
	}

	cleaned := cleanOutput(out)

	// Extract first valid absolute path
	var binPath string
	lines := strings.Split(cleaned, "\n")
	for _, line := range lines {
		// Skip anything that looks like error messages
		if strings.Contains(line, ": ") || strings.Contains(line, "error") {
			continue
		}

		// Must be absolute path without spaces
		if strings.HasPrefix(line, "/") && !strings.ContainsAny(line, " \t\r") {
			binPath = filepath.Clean(line)
			break
		}
	}

	if binPath == "" {
		return fmt.Errorf("binary %q not found in container (searched in: %q)", binName, cleaned)
	}

	// Verify the binary exists and is executable
	if _, err := runCommandWithOutput(
		containerName, false, "", "", 0,
		[]string{"test", "-x", binPath}, false, false, true,
	); err != nil {
		return fmt.Errorf("found path %q but it's not executable: %v", binPath, err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	binDir := filepath.Join(homeDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %v", err)
	}

	exportedBinPath := filepath.Join(binDir, binName)
	script := fmt.Sprintf(`#!/bin/sh
	# toolbox_binary
	# name: %s
	BIN_PATH="%s"
	exec toolbox run -c %s "$BIN_PATH" "$@"
	`, containerName, binPath, containerName)

	if err := os.WriteFile(exportedBinPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to create wrapper: %v", err)
	}

	fmt.Printf("Successfully exported %s from container %s to %s\n", binName, containerName, exportedBinPath)
	return nil
}

func exportApplication(appName, containerName string) error {
	// Step 1: Run find inside container to locate candidate desktop files
	findCmd := []string{"sh", "-c", fmt.Sprintf("find /usr/share/applications -name '%s.desktop'", appName)}
	out, err := runCommandWithOutput(containerName, false, "", "", 0, findCmd, false, false, true)
	if err != nil || strings.TrimSpace(out) == "" {
		return fmt.Errorf("Error: application %s not found in container", appName)
	}

	// Step 2: Nuclear cleaning of find output
	cleanOutput := func(s string) string {
		ansiRegex := regexp.MustCompile(`\x1B(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)
		s = ansiRegex.ReplaceAllString(s, "")
		s = regexp.MustCompile(`\x1B\]8;;.*?\x1B\\`).ReplaceAllString(s, "")
		s = strings.Map(func(r rune) rune {
			if r >= 32 && r <= 126 || r == '\n' {
				return r
			}
			return -1
		}, s)
		lines := strings.Split(s, "\n")
		var cleanLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				cleanLines = append(cleanLines, line)
			}
		}
		return strings.Join(cleanLines, "\n")
	}

	cleaned := cleanOutput(out)

	// Step 3: Scan bottom-to-top for the first valid absolute path ending with .desktop
	var desktopFile string
	lines := strings.Split(cleaned, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.HasPrefix(line, "/") && strings.HasSuffix(line, ".desktop") && !strings.ContainsAny(line, " \t") {
			desktopFile = filepath.Clean(line)
			break
		}
	}

	if desktopFile == "" {
		return fmt.Errorf("invalid desktop file path in container output: %q", cleaned)
	}

	// Step 4: Read the desktop file content safely
	catCmd := []string{"cat", desktopFile}
	content, err := runCommandWithOutput(containerName, false, "", "", 0, catCmd, false, false, true)
	if err != nil {
		return fmt.Errorf("failed to read desktop file %q: %v", desktopFile, err)
	}

	lines = strings.Split(cleanOutput(content), "\n")
	var newLines []string
	started := false
	hasNameTranslations := false

	// Step 5: Rewrite desktop file fields
	for _, line := range lines {
		if !started {
			if strings.TrimSpace(line) == "[Desktop Entry]" {
				started = true
				newLines = append(newLines, line)
			}
			continue
		}

		if strings.HasPrefix(line, "Exec=") {
			execCmd := line[5:]
			line = fmt.Sprintf("Exec=toolbox run -c %s %s", containerName, execCmd)
		} else if strings.HasPrefix(line, "Name=") {
			line = fmt.Sprintf("Name=%s (on %s)", line[5:], containerName)
		} else if strings.HasPrefix(line, "Name[") {
			hasNameTranslations = true
		} else if strings.HasPrefix(line, "GenericName=") {
			line = fmt.Sprintf("GenericName=%s (on %s)", line[12:], containerName)
		} else if strings.HasPrefix(line, "TryExec=") || line == "DBusActivatable=true" {
			continue
		}
		newLines = append(newLines, line)
	}

	if hasNameTranslations {
		for i, line := range newLines {
			if strings.HasPrefix(line, "Name[") {
				lang := line[5:strings.Index(line, "]")]
				value := line[strings.Index(line, "=")+1:]
				newLines[i] = fmt.Sprintf("Name[%s]=%s (on %s)", lang, value, containerName)
			}
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	appsPath := filepath.Join(homeDir, ".local", "share", "applications")
	exportedPath := filepath.Join(appsPath, filepath.Base(desktopFile))
	exportedPath = strings.TrimSuffix(exportedPath, ".desktop") + "-" + containerName + ".desktop"

	if err := os.MkdirAll(appsPath, 0755); err != nil {
		return fmt.Errorf("failed to create applications directory: %v", err)
	}
	if err := os.WriteFile(exportedPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to create desktop file: %v", err)
	}

	exec.Command("update-desktop-database", appsPath).Run()

	fmt.Printf("Successfully exported %s from container %s to %s\n", appName, containerName, exportedPath)
	return nil
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

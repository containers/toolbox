package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	exportBin      string
	exportApp      string
	exportContainer string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export binaries or applications from a toolbox container",
	RunE:  runExport,
}

func init() {
	exportCmd.Flags().StringVar(&exportBin, "bin", "", "Path or name of binary to export")
	exportCmd.Flags().StringVar(&exportApp, "app", "", "Path or name of application to export")
	exportCmd.Flags().StringVar(&exportContainer, "container", "", "Name of the toolbox container")
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
	} else if exportApp != "" {
		return exportApplication(exportApp, exportContainer)
	}
	return nil
}

func exportBinary(binName, containerName string) error {
	// Find the binary's full path inside the container
	checkCmd := fmt.Sprintf("toolbox run -c %s which %s", containerName, binName)
	out, err := exec.Command("sh", "-c", checkCmd).Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return fmt.Errorf("binary %s not found in container %s", binName, containerName)
	}
	binPath := strings.TrimSpace(string(out))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	exportedBinPath := filepath.Join(homeDir, ".local", "bin", binName)

	script := fmt.Sprintf(`#!/bin/sh
	# toolbox_binary
	# name: %s
	exec toolbox run -c %s %s "$@"
	`, containerName, containerName, binPath)

	if err := os.WriteFile(exportedBinPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to create wrapper: %v", err)
	}

	fmt.Printf("Successfully exported %s from container %s to %s\n", binName, containerName, exportedBinPath)
	return nil
}

func exportApplication(appName, containerName string) error {
	// Find the desktop file inside the container
	findCmd := fmt.Sprintf("toolbox run -c %s sh -c 'find /usr/share/applications -name \"*%s*.desktop\" | head -1'", containerName, appName)
	out, err := exec.Command("sh", "-c", findCmd).Output()
	if err != nil || strings.TrimSpace(string(out)) == "" {
		return fmt.Errorf("application %s not found in container %s", appName, containerName)
	}
	desktopFile := strings.TrimSpace(string(out))

	// Read the desktop file content
	catCmd := fmt.Sprintf("toolbox run -c %s cat %s", containerName, desktopFile)
	content, err := exec.Command("sh", "-c", catCmd).Output()
	if err != nil {
		return fmt.Errorf("failed to read desktop file: %v", err)
	}
	lines := strings.Split(string(content), "\n")
	var newLines []string
	hasNameTranslations := false

	for _, line := range lines {
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

	// Update desktop database
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

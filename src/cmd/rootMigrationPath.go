//
// Copyright © 2021 – 2025 Red Hat Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

//go:build migration_path_for_coreos_toolbox
// +build migration_path_for_coreos_toolbox

package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

func preRunIsCoreOSBug() error {
	if containerType := os.Getenv("container"); containerType == "" {
		var builder strings.Builder
		builder.WriteString("/run/.containerenv found on what looks like the host\n")
		builder.WriteString("If this is the host, then remove /run/.containerenv and try again.\n")
		builder.WriteString("Otherwise, contact your system administrator or file a bug.")

		errMsg := builder.String()
		return errors.New(errMsg)
	}

	return nil
}

func rootRunImpl(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		panic("unexpected argument: commands known or unknown shouldn't reach here")
	}

	if utils.IsInsideContainer() {
		if !utils.IsInsideToolboxContainer() {
			return errors.New("this is not a Toolbx container")
		}

		exitCode, err := utils.ForwardToHost()
		return &exitError{exitCode, err}
	}

	container, image, release, err := resolveContainerAndImageNames("", "", "", "", "")
	if err != nil {
		return err
	}

	userShell := os.Getenv("SHELL")
	if userShell == "" {
		return errors.New("failed to get the current user's default shell")
	}

	command := []string{userShell, "-l"}

	if err := runCommand(container, true, image, release, 0, command, true, true, false); err != nil {
		return err
	}

	return nil
}

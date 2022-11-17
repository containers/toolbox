/*
 * Copyright © 2021 – 2022 Red Hat Inc.
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
	"os"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:                   "completion",
	Short:                 "Generate completion script",
	Hidden:                true,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "fish", "zsh"},
	Args:                  cobra.ExactValidArgs(1),
	RunE:                  completion,
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func completion(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "bash":
		err := cmd.Root().GenBashCompletionV2(os.Stdout, true)
		return err
	case "fish":
		err := cmd.Root().GenFishCompletion(os.Stdout, true)
		return err
	case "zsh":
		err := cmd.Root().GenZshCompletion(os.Stdout)
		return err
	}

	panic("code should not be reached")
}

func completionEmpty(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

func completionCommands(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	commandNames := []string{}
	commands := cmd.Root().Commands()
	for _, command := range commands {
		if strings.Contains(command.Name(), "complet") {
			continue
		}
		commandNames = append(commandNames, command.Name())
	}

	return commandNames, cobra.ShellCompDirectiveNoFileComp
}

func completionContainerNames(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	var containerNames []string
	if containers, err := getContainers(); err == nil {
		for _, container := range containers {
			containerNames = append(containerNames, container.Names[0])
		}
	}

	return containerNames, cobra.ShellCompDirectiveNoFileComp
}

func completionContainerNamesFiltered(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if cmd.Name() == "enter" && len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var containerNames []string
	if containers, err := getContainers(); err == nil {
		for _, container := range containers {
			skip := false
			for _, arg := range args {
				if container.Names[0] == arg {
					skip = true
					break
				}
			}

			if skip {
				continue
			}

			containerNames = append(containerNames, container.Names[0])
		}
	}

	return containerNames, cobra.ShellCompDirectiveNoFileComp

}

func completionDistroNames(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	imageFlag := cmd.Flag("image")
	if imageFlag != nil && imageFlag.Changed {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	supportedDistros := utils.GetSupportedDistros()

	return supportedDistros, cobra.ShellCompDirectiveNoFileComp
}

func completionImageNames(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	distroFlag := cmd.Flag("distro")
	if distroFlag != nil && distroFlag.Changed {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var imageNames []string
	if images, err := getImages(); err == nil {
		for _, image := range images {
			if len(image.Names) > 0 {
				imageNames = append(imageNames, image.Names[0])
			} else {
				imageNames = append(imageNames, image.ID)
			}
		}
	}

	return imageNames, cobra.ShellCompDirectiveNoFileComp
}

func completionImageNamesFiltered(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	var imageNames []string
	if images, err := getImages(); err == nil {
		for _, image := range images {
			skip := false
			var imageName string

			if len(image.Names) > 0 {
				imageName = image.Names[0]
			} else {
				imageName = image.ID
			}

			for _, arg := range args {
				if arg == imageName {
					skip = true
					break
				}
			}

			if skip {
				continue
			}

			imageNames = append(imageNames, imageName)
		}
	}

	return imageNames, cobra.ShellCompDirectiveNoFileComp
}

func completionLogLevels(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
}

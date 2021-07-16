package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(toolbox completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ toolbox completion bash > /etc/bash_completion.d/toolbox
  # macOS:
  $ toolbox completion bash > /usr/local/etc/bash_completion.d/toolbox

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ toolbox completion zsh > "${fpath[1]}/_toolbox"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ toolbox completion fish | source

  # To load completions for each session, execute once:
  $ toolbox completion fish > ~/.config/fish/completions/toolbox.fish

`,
	Hidden:                true,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			if err := cmd.Root().GenBashCompletion(os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err)
			}
		case "zsh":
			if err := cmd.Root().GenZshCompletion(os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err)
			}
		case "fish":
			if err := cmd.Root().GenFishCompletion(os.Stdout, true); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
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
	containerNames := []string{}
	if containers, err := getContainers(); err == nil {
		for _, container := range containers {
			containerNames = append(containerNames, container.Names[0])
		}
	}

	if len(containerNames) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return containerNames, cobra.ShellCompDirectiveNoFileComp
}

func completionContainerNamesFiltered(cmd *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if cmd.Name() == "enter" && len(args) >= 1 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	containerNames := []string{}
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

	if len(containerNames) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
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

	imageNames := []string{}
	if images, err := getImages(); err == nil {
		for _, image := range images {
			if len(image.Names) > 0 {
				imageNames = append(imageNames, image.Names[0])
			} else {
				imageNames = append(imageNames, image.ID)
			}
		}
	}

	if len(imageNames) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return imageNames, cobra.ShellCompDirectiveNoFileComp
}

func completionImageNamesFiltered(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	imageNames := []string{}
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

	if len(imageNames) == 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	return imageNames, cobra.ShellCompDirectiveNoFileComp
}

func completionLogLevels(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"trace", "debug", "info", "warn", "error", "fatal", "panic"}, cobra.ShellCompDirectiveNoFileComp
}

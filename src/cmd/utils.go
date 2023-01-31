/*
 * Copyright © 2020 – 2022 Red Hat Inc.
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
	"strings"
	"syscall"

	"github.com/containers/toolbox/pkg/utils"
)

// askForConfirmation prints prompt to stdout and waits for response from the
// user
//
// Expected answers are: "yes", "y", "no", "n"
//
// Answers are internally converted to lower case.
//
// The default answer is "no" ([y/N])
func askForConfirmation(prompt string) bool {
	var retVal bool

	for {
		fmt.Printf("%s ", prompt)

		var response string

		fmt.Scanf("%s", &response)
		if response == "" {
			response = "n"
		} else {
			response = strings.ToLower(response)
		}

		if response == "no" || response == "n" {
			break
		} else if response == "yes" || response == "y" {
			retVal = true
			break
		}
	}

	return retVal
}

func createErrorContainerNotFound(container string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "container %s not found\n", container)
	fmt.Fprintf(&builder, "Use the 'create' command to create a toolbox.\n")
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func createErrorDistroWithoutRelease(distro string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "option '--release' is needed\n")
	fmt.Fprintf(&builder, "Distribution %s doesn't match the host.\n", distro)
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func createErrorInvalidContainer(containerArg string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "invalid argument for '%s'\n", containerArg)
	fmt.Fprintf(&builder, "Container names must match '%s'.\n", utils.ContainerNameRegexp)
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func createErrorInvalidDistro(distro string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "invalid argument for '--distro'\n")
	fmt.Fprintf(&builder, "Distribution %s is unsupported.\n", distro)
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func createErrorInvalidImageForContainerName(container string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "invalid argument for '--image'\n")
	fmt.Fprintf(&builder, "Container name %s generated from image is invalid.\n", container)
	fmt.Fprintf(&builder, "Container names must match '%s'.\n", utils.ContainerNameRegexp)
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func createErrorInvalidImageWithoutBasename() error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "invalid argument for '--image'\n")
	fmt.Fprintf(&builder, "Images must have basenames.\n")
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func createErrorInvalidRelease(hint string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "invalid argument for '--release'\n")
	fmt.Fprintf(&builder, "%s\n", hint)
	fmt.Fprintf(&builder, "Run '%s --help' for usage.", executableBase)

	errMsg := builder.String()
	return errors.New(errMsg)
}

func getUsageForCommonCommands() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "create    Create a new toolbox container\n")
	fmt.Fprintf(&builder, "enter     Enter an existing toolbox container\n")
	fmt.Fprintf(&builder, "list      List all existing toolbox containers and images\n")

	usage := builder.String()
	return usage
}

func resolveContainerAndImageNames(container, containerArg, distroCLI, imageCLI, releaseCLI string) (
	string, string, string, error,
) {
	container, image, release, err := utils.ResolveContainerAndImageNames(container,
		distroCLI,
		imageCLI,
		releaseCLI)

	if err != nil {
		var errContainer *utils.ContainerError
		var errDistro *utils.DistroError
		var errImage *utils.ImageError
		var errParseRelease *utils.ParseReleaseError

		if errors.As(err, &errContainer) {
			if errors.Is(err, utils.ErrContainerNameInvalid) {
				if containerArg == "" {
					panicMsg := fmt.Sprintf("unexpected %T without containerArg: %s", err, err)
					panic(panicMsg)
				}

				err := createErrorInvalidContainer(containerArg)
				return "", "", "", err
			} else if errors.Is(err, utils.ErrContainerNameFromImageInvalid) {
				err := createErrorInvalidImageForContainerName(errContainer.Container)
				return "", "", "", err
			} else {
				panicMsg := fmt.Sprintf("unexpected %T: %s", err, err)
				panic(panicMsg)
			}
		} else if errors.As(err, &errDistro) {
			if errors.Is(err, utils.ErrDistroUnsupported) {
				err := createErrorInvalidDistro(errDistro.Distro)
				return "", "", "", err
			} else if errors.Is(err, utils.ErrDistroWithoutRelease) {
				err := createErrorDistroWithoutRelease(errDistro.Distro)
				return "", "", "", err
			} else {
				panicMsg := fmt.Sprintf("unexpected %T: %s", err, err)
				panic(panicMsg)
			}
		} else if errors.As(err, &errImage) {
			if errors.Is(err, utils.ErrImageWithoutBasename) {
				err := createErrorInvalidImageWithoutBasename()
				return "", "", "", err
			} else {
				panicMsg := fmt.Sprintf("unexpected %T: %s", err, err)
				panic(panicMsg)
			}
		} else if errors.As(err, &errParseRelease) {
			err := createErrorInvalidRelease(errParseRelease.Hint)
			return "", "", "", err
		} else {
			return "", "", "", err
		}
	}

	return container, image, release, nil
}

// showManual tries to open the specified manual page using man on stdout
func showManual(manual string) error {
	manBinary, err := exec.LookPath("man")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			fmt.Printf("toolbox - Tool for containerized command line environments on Linux\n")
			fmt.Printf("\n")
			fmt.Printf("Common commands are:\n")

			usage := getUsageForCommonCommands()
			fmt.Printf("%s", usage)

			fmt.Printf("\n")
			fmt.Printf("Go to https://github.com/containers/toolbox for further information.\n")
			return nil
		}

		return errors.New("failed to look up man(1)")
	}

	manualArgs := []string{"man", manual}
	env := os.Environ()

	stderrFd := os.Stderr.Fd()
	stderrFdInt := int(stderrFd)
	stdoutFd := os.Stdout.Fd()
	stdoutFdInt := int(stdoutFd)
	if err := syscall.Dup3(stdoutFdInt, stderrFdInt, 0); err != nil {
		return errors.New("failed to redirect standard error to standard output")
	}

	if err := syscall.Exec(manBinary, manualArgs, env); err != nil {
		return errors.New("failed to invoke man(1)")
	}

	return nil
}

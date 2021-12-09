/*
 * Copyright © 2020 – 2021 Red Hat Inc.
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

func createErrorInvalidRelease() error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "invalid argument for '--release'\n")
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

		return errors.New("failed to lookup man(1)")
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

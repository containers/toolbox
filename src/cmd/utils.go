/*
 * Copyright © 2020 – 2024 Red Hat Inc.
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
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type askForConfirmationPreFunc func() error
type pollFunc func(error, []unix.PollFd) error

var (
	errClosed = errors.New("closed")

	errContinue = errors.New("continue")

	errHUP = errors.New("HUP")
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

	ctx := context.Background()
	retValCh, errCh := askForConfirmationAsync(ctx, prompt, nil)

	select {
	case val := <-retValCh:
		retVal = val
	case err := <-errCh:
		logrus.Debugf("Failed to ask for confirmation: %s", err)
		retVal = false
	}

	return retVal
}

func askForConfirmationAsync(ctx context.Context,
	prompt string,
	askForConfirmationPreFn askForConfirmationPreFunc) (<-chan bool, <-chan error) {

	retValCh := make(chan bool, 1)
	errCh := make(chan error, 1)

	done := ctx.Done()
	eventFD := -1
	if done != nil {
		fd, err := unix.Eventfd(0, unix.EFD_CLOEXEC|unix.EFD_NONBLOCK)
		if err != nil {
			errCh <- fmt.Errorf("eventfd(2) failed: %w", err)
			return retValCh, errCh
		}

		eventFD = fd
	}

	go func() {
		for {
			fmt.Printf("%s ", prompt)
			if askForConfirmationPreFn != nil {
				if err := askForConfirmationPreFn(); err != nil {
					if errors.Is(err, errContinue) {
						continue
					}

					errCh <- err
					break
				}
			}

			var response string

			pollFn := func(errPoll error, pollFDs []unix.PollFd) error {
				if len(pollFDs) != 1 {
					panic("unexpected number of file descriptors")
				}

				if errPoll != nil {
					return errPoll
				}

				if pollFDs[0].Revents&unix.POLLIN != 0 {
					logrus.Debug("Returned from /dev/stdin: POLLIN")

					scanner := bufio.NewScanner(os.Stdin)
					if !scanner.Scan() {
						if err := scanner.Err(); err != nil {
							return err
						} else {
							return io.EOF
						}
					}

					response = scanner.Text()
					return nil
				}

				if pollFDs[0].Revents&unix.POLLHUP != 0 {
					logrus.Debug("Returned from /dev/stdin: POLLHUP")
					return errHUP
				}

				if pollFDs[0].Revents&unix.POLLNVAL != 0 {
					logrus.Debug("Returned from /dev/stdin: POLLNVAL")
					return errClosed
				}

				return errContinue
			}

			stdinFD := int32(os.Stdin.Fd())

			err := poll(pollFn, int32(eventFD), stdinFD)
			if err != nil {
				errCh <- err
				break
			}

			if response == "" {
				response = "n"
			} else {
				response = strings.ToLower(response)
			}

			if response == "no" || response == "n" {
				retValCh <- false
				break
			} else if response == "yes" || response == "y" {
				retValCh <- true
				break
			}
		}
	}()

	watchContextForEventFD(ctx, eventFD)
	return retValCh, errCh
}

func discardInputAsync(ctx context.Context) (<-chan int, <-chan error) {
	retValCh := make(chan int, 1)
	errCh := make(chan error, 1)

	done := ctx.Done()
	eventFD := -1
	if done != nil {
		fd, err := unix.Eventfd(0, unix.EFD_CLOEXEC|unix.EFD_NONBLOCK)
		if err != nil {
			errCh <- fmt.Errorf("eventfd(2) failed: %w", err)
			return retValCh, errCh
		}

		eventFD = fd
	}

	go func() {
		var total int

		for {
			pollFn := func(errPoll error, pollFDs []unix.PollFd) error {
				if len(pollFDs) != 1 {
					panic("unexpected number of file descriptors")
				}

				if errPoll != nil && !errors.Is(errPoll, context.Canceled) {
					return errPoll
				}

				if pollFDs[0].Revents&unix.POLLIN != 0 {
					logrus.Debug("Returned from /dev/stdin: POLLIN")

					buffer := make([]byte, bytes.MinRead)
					n, err := os.Stdin.Read(buffer)
					total += n

					if errPoll != nil {
						return errPoll
					} else if err != nil {
						return err
					}

					return nil
				}

				if pollFDs[0].Revents&unix.POLLHUP != 0 {
					logrus.Debug("Returned from /dev/stdin: POLLHUP")

					if errPoll != nil {
						return errPoll
					}

					return errHUP
				}

				if pollFDs[0].Revents&unix.POLLNVAL != 0 {
					logrus.Debug("Returned from /dev/stdin: POLLNVAL")

					if errPoll != nil {
						return errPoll
					}

					return errClosed
				}

				if errPoll != nil {
					return errPoll
				}

				return errContinue
			}

			stdinFD := int32(os.Stdin.Fd())

			err := poll(pollFn, int32(eventFD), stdinFD)
			if err != nil {
				retValCh <- total
				errCh <- err
				break
			}
		}
	}()

	watchContextForEventFD(ctx, eventFD)
	return retValCh, errCh
}

func createErrorContainerNotFound(container string) error {
	var builder strings.Builder
	fmt.Fprintf(&builder, "container %s not found\n", container)
	fmt.Fprintf(&builder, "Use the 'create' command to create a Toolbx.\n")
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

func getCDIFileForNvidia(targetUser *user.User) (string, error) {
	toolboxRuntimeDirectory, err := utils.GetRuntimeDirectory(targetUser)
	if err != nil {
		return "", err
	}

	cdiFile := filepath.Join(toolboxRuntimeDirectory, "cdi-nvidia.json")
	return cdiFile, nil
}

func getUsageForCommonCommands() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "create    Create a new Toolbx container\n")
	fmt.Fprintf(&builder, "enter     Enter an existing Toolbx container\n")
	fmt.Fprintf(&builder, "list      List all existing Toolbx containers and images\n")

	usage := builder.String()
	return usage
}

func poll(pollFn pollFunc, eventFD int32, fds ...int32) error {
	if len(fds) == 0 {
		panic("file descriptors not specified")
	}

	pollFDs := []unix.PollFd{
		{
			Fd:      eventFD,
			Events:  unix.POLLIN,
			Revents: 0,
		},
	}

	for _, fd := range fds {
		pollFD := unix.PollFd{Fd: fd, Events: unix.POLLIN, Revents: 0}
		pollFDs = append(pollFDs, pollFD)
	}

	for {
		if _, err := unix.Poll(pollFDs, -1); err != nil {
			if errors.Is(err, unix.EINTR) {
				logrus.Debugf("Failed to poll(2): %s: ignoring", err)
				continue
			}

			return fmt.Errorf("poll(2) failed: %w", err)
		}

		var err error

		if pollFDs[0].Revents&unix.POLLIN != 0 {
			logrus.Debug("Returned from eventfd: POLLIN")
			err = context.Canceled

			for {
				buffer := make([]byte, 8)
				if n, err := unix.Read(int(eventFD), buffer); n != len(buffer) || err != nil {
					break
				}
			}
		} else if pollFDs[0].Revents&unix.POLLNVAL != 0 {
			logrus.Debug("Returned from eventfd: POLLNVAL")
			err = context.Canceled
		}

		if err := pollFn(err, pollFDs[1:]); !errors.Is(err, errContinue) {
			return err
		}
	}
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
			fmt.Printf("toolbox - Tool for interactive command line environments on Linux\n")
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

func watchContextForEventFD(ctx context.Context, eventFD int) {
	done := ctx.Done()
	if done == nil {
		return
	}

	if eventFD < 0 {
		panic("invalid file descriptor for eventfd")
	}

	go func() {
		defer unix.Close(eventFD)

		select {
		case <-done:
			buffer := make([]byte, 8)
			binary.PutUvarint(buffer, 1)

			if _, err := unix.Write(eventFD, buffer); err != nil {
				panicMsg := fmt.Sprintf("write(2) to eventfd failed: %s", err)
				panic(panicMsg)
			}
		}
	}()
}

/*
 * Copyright © 2019 – 2023 Red Hat Inc.
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

package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Run(name string, stdin io.Reader, stdout, stderr io.Writer, arg ...string) error {
	ctx := context.Background()
	err := RunContext(ctx, name, stdin, stdout, stderr, arg...)
	return err
}

func RunContext(ctx context.Context, name string, stdin io.Reader, stdout, stderr io.Writer, arg ...string) error {
	exitCode, err := RunContextWithExitCode(ctx, name, stdin, stdout, stderr, arg...)
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("failed to invoke %s(1)", name)
	}
	return nil
}

func RunContextWithExitCode(ctx context.Context,
	name string,
	stdin io.Reader,
	stdout, stderr io.Writer,
	arg ...string) (int, error) {

	logLevel := logrus.GetLevel()
	if stderr == nil && logLevel >= logrus.DebugLevel {
		stderr = os.Stderr
	}

	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return 1, fmt.Errorf("%s(1) not found", name)
		}

		if ctxErr := ctx.Err(); ctxErr != nil {
			return 1, ctxErr
		}

		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode := exitErr.ExitCode()
			return exitCode, nil
		}

		return 1, fmt.Errorf("failed to invoke %s(1)", name)
	}

	return 0, nil
}

func RunWithExitCode(name string, stdin io.Reader, stdout, stderr io.Writer, arg ...string) (int, error) {
	ctx := context.Background()
	exitCode, err := RunContextWithExitCode(ctx, name, stdin, stdout, stderr, arg...)
	return exitCode, err
}

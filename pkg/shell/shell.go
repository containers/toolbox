/*
 * Copyright © 2019 – 2021 Red Hat Inc.
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
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Run(name string, stdin io.Reader, stdout, stderr io.Writer, arg ...string) error {
	exitCode, err := RunWithExitCode(name, stdin, stdout, stderr, arg...)
	if exitCode != 0 && err == nil {
		err = fmt.Errorf("failed to invoke %s(1)", name)
	}

	if err != nil {
		return err
	}

	return nil
}

func RunWithExitCode(name string, stdin io.Reader, stdout, stderr io.Writer, arg ...string) (int, error) {
	logLevel := logrus.GetLevel()
	if stderr == nil && logLevel >= logrus.DebugLevel {
		stderr = os.Stderr
	}

	cmd := exec.Command(name, arg...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return 1, fmt.Errorf("%s(1) not found", name)
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

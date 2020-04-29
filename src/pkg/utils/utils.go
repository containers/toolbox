/*
 * Copyright © 2019 – 2020 Red Hat Inc.
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

package utils

import (
	"os"
	"os/exec"
	"syscall"
)

func ShowManual(manual string) error {
	manBinary, err := exec.LookPath("man")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return errors.New("man(1) not found")
		}

		return errors.New("failed to lookup man(1)")
	}

	manualArgs := []string{"man", manual}
	env := os.Environ()

	stderrFd := os.Stderr.Fd()
	stderrFdInt := int(stderrFd)
	stdoutFd := os.Stdout.Fd()
	stdoutFdInt := int(stdoutFd)
	if err := syscall.Dup2(stdoutFdInt, stderrFdInt); err != nil {
		return errors.New("failed to redirect standard error to standard output")
	}

	if err := syscall.Exec(manBinary, manualArgs, env); err != nil {
		return errors.New("failed to invoke man(1)")
	}

	return nil
}

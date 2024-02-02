/*
 * Copyright © 2023 – 2024 Red Hat Inc.
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

package shell_test

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/containers/toolbox/pkg/shell"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestShellRun(t *testing.T) {
	type input struct {
		commandName string
		stdIn       io.Reader
		args        []string
		loglevel    logrus.Level
		useStdErr   bool
	}

	type expect struct {
		err    error
		stdout []byte
		stderr []byte
	}

	testCases := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "OK",
			input: input{
				commandName: "echo",
				stdIn:       os.Stdin,
				args:        []string{"Toolbx test"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    nil,
				stdout: []byte("Toolbx test\n"),
				stderr: nil,
			},
		},
		{
			name: "FAIL_NonExisting_Command",
			input: input{
				commandName: "no-exist-executable",
				stdIn:       os.Stdin,
				args:        []string{"Toolbx test"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    errors.New("no-exist-executable(1) not found"),
				stdout: nil,
				stderr: nil,
			},
		},
		{
			name: "FAIL_Unexpected_Command_Result",
			input: input{
				commandName: "cat",
				stdIn:       os.Stdin,
				args:        []string{"/bogus/file.foo"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   true,
			},
			expect: expect{
				err:    errors.New("failed to invoke cat(1)"),
				stdout: nil,
				stderr: []byte("cat: /bogus/file.foo: No such file or directory\n"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var actualStdOut outputMock
			var actualStdErr outputMock

			logrus.SetLevel(tc.input.loglevel)

			err := shell.Run(tc.input.commandName, tc.input.stdIn, &actualStdOut, &actualStdErr, tc.input.args...)

			assert.Equal(t, tc.expect.err, err)
			assert.Equal(t, tc.expect.stdout, actualStdOut.written)
			// If we expect and std err value, we need to cast with std out mock struct
			// to ensure the expected value, otherwise the std err should be empty value
			if tc.input.useStdErr {
				assert.Equal(t, tc.expect.stderr, actualStdErr.written)
			} else {
				assert.Empty(t, actualStdErr)
			}
		})
	}
}

func TestRunContextWithExitCode(t *testing.T) {
	testCases := []struct {
		cancel   bool
		command  []string
		err      error
		errMsg   string
		exitCode int
		stdout   []byte
		stderr   []byte
		timeout  time.Duration
	}{
		{
			command: []string{"true"},
		},
		{
			command:  []string{"false"},
			exitCode: 1,
		},
		{
			command: []string{"echo"},
			stdout:  []byte("\n"),
		},
		{
			command:  []string{"echo", "hello, world"},
			err:      nil,
			exitCode: 0,
			stdout:   []byte("hello, world\n"),
		},
		{
			command:  []string{"command-does-not-exist"},
			errMsg:   "command-does-not-exist(1) not found",
			exitCode: 1,
		},
		{
			command:  []string{"cat", "/file/does/not/exist"},
			exitCode: 1,
			stderr:   []byte("cat: /file/does/not/exist: No such file or directory\n"),
		},
		{
			cancel:   true,
			command:  []string{"sleep", "+Inf"},
			err:      context.Canceled,
			exitCode: 1,
		},
		{
			command:  []string{"sleep", "+Inf"},
			err:      context.DeadlineExceeded,
			exitCode: 1,
			timeout:  1 * time.Second,
		},
	}

	for _, tc := range testCases {
		name := strings.Join(tc.command, " ")
		if tc.cancel {
			name += " (cancel)"
		}
		if tc.timeout != 0 {
			name += " (timeout)"
		}

		t.Run(name, func(t *testing.T) {
			var cancel context.CancelFunc
			ctx := context.Background()
			if tc.cancel {
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()
			}
			if tc.timeout != 0 {
				ctx, cancel = context.WithTimeout(ctx, tc.timeout)
				defer cancel()
			}

			if tc.cancel {
				cancel()
			}

			var stdout outputMock
			var stderr outputMock

			exitCode, err := shell.RunContextWithExitCode(ctx,
				tc.command[0],
				os.Stdin,
				&stdout,
				&stderr,
				tc.command[1:]...)

			if tc.err == nil && tc.errMsg == "" {
				assert.NoError(t, err)
			}

			if tc.err != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.err)
			}

			if tc.errMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errMsg)
			}

			assert.Equal(t, tc.exitCode, exitCode)

			if tc.stdout == nil {
				assert.Empty(t, stdout.written)
			} else {
				assert.NotEmpty(t, stdout.written)
				assert.Equal(t, tc.stdout, stdout.written)
			}

			if tc.stderr == nil {
				assert.Empty(t, stderr.written)
			} else {
				assert.NotEmpty(t, stderr.written)
				assert.Equal(t, tc.stderr, stderr.written)
			}
		})
	}
}

func TestShellRunWithExitCode(t *testing.T) {
	type input struct {
		commandName string
		stdIn       io.Reader
		args        []string
		loglevel    logrus.Level
		useStdErr   bool
	}

	type expect struct {
		err    error
		code   int
		stdout []byte
		stderr []byte
	}

	testCases := []struct {
		name   string
		input  input
		expect expect
	}{
		{
			name: "OK_Without_stderr_and_info_log_level",
			input: input{
				commandName: "echo",
				stdIn:       os.Stdin,
				args:        []string{"Toolbx test"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    nil,
				code:   0,
				stdout: []byte("Toolbx test\n"),
				stderr: nil,
			},
		},
		{
			name: "OK_Without_stderr_and_debug_log_level",
			input: input{
				commandName: "echo",
				stdIn:       os.Stdin,
				args:        []string{"Toolbx test"},
				loglevel:    logrus.DebugLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    nil,
				code:   0,
				stdout: []byte("Toolbx test\n"),
				stderr: nil,
			},
		},
		{
			name: "OK_With_stderr_and_info_log_level",
			input: input{
				commandName: "echo",
				stdIn:       os.Stdin,
				args:        nil,
				loglevel:    logrus.InfoLevel,
				useStdErr:   true,
			},
			expect: expect{
				err:    nil,
				code:   0,
				stdout: []byte("\n"),
				stderr: nil,
			},
		},
		{
			name: "FAIL_NonExisting_Command",
			input: input{
				commandName: "no-exist-executable",
				stdIn:       os.Stdin,
				args:        []string{"Toolbx test"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    errors.New("no-exist-executable(1) not found"),
				code:   1,
				stdout: nil,
				stderr: nil,
			},
		},
		{
			name: "FAIL_Unexpected_Command_Result",
			input: input{
				commandName: "cat",
				stdIn:       os.Stdin,
				args:        []string{"/bogus/file.foo"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   true,
			},
			expect: expect{
				err:    nil,
				code:   1,
				stdout: nil,
				stderr: []byte("cat: /bogus/file.foo: No such file or directory\n"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var actualStdOut outputMock
			var actualStdErr outputMock

			logrus.SetLevel(tc.input.loglevel)

			code, err := shell.RunWithExitCode(tc.input.commandName, tc.input.stdIn, &actualStdOut, &actualStdErr, tc.input.args...)

			assert.Equal(t, tc.expect.err, err)
			assert.Equal(t, tc.expect.code, code)
			assert.Equal(t, tc.expect.stdout, actualStdOut.written)
			if tc.input.useStdErr {
				assert.Equal(t, tc.expect.stderr, actualStdErr.written)
			} else {
				assert.Empty(t, actualStdErr)
			}
		})
	}
}

// outputMock is a mock to ensure content written to stdout/stderr was correct
type outputMock struct {
	written []byte
}

func (mock *outputMock) Write(p []byte) (n int, err error) {
	mock.written = append(mock.written, p...)
	return len(p), nil
}

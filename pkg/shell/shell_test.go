package shell_test

import (
	"errors"
	"io"
	"os"
	"testing"

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
				args:        []string{"toolbox test"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    nil,
				stdout: []byte("toolbox test\n"),
				stderr: nil,
			},
		},
		{
			name: "FAIL_NonExisting_Command",
			input: input{
				commandName: "no-exist-executable",
				stdIn:       os.Stdin,
				args:        []string{"toolbox test"},
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
				args:        []string{"toolbox test"},
				loglevel:    logrus.InfoLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    nil,
				code:   0,
				stdout: []byte("toolbox test\n"),
				stderr: nil,
			},
		},
		{
			name: "OK_Without_stderr_and_debug_log_level",
			input: input{
				commandName: "echo",
				stdIn:       os.Stdin,
				args:        []string{"toolbox test"},
				loglevel:    logrus.DebugLevel,
				useStdErr:   false,
			},
			expect: expect{
				err:    nil,
				code:   0,
				stdout: []byte("toolbox test\n"),
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
				args:        []string{"toolbox test"},
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

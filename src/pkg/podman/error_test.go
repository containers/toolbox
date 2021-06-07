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

package podman

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalError(t *testing.T) {
	type expect struct {
		IsNil  bool
		Error  string
		Search string
		Wrap   []string
	}

	testCases := []struct {
		name   string
		input  string
		expect expect
	}{
		{
			name:  "Empty error message",
			input: "",
			expect: expect{
				IsNil: true,
				Error: "",
			},
		},
		{
			name:  "Text with no prolog before an error message",
			input: "There is no prolog before the error message.",
			expect: expect{
				Error: "There is no prolog before the error message.",
			},
		},
		{
			name:  "Text with prolog before an error message",
			input: "There is the prolog.Error: an error message",
			expect: expect{
				Error:  "an error message",
				Search: "an error message",
			},
		},
		{
			name:  "Text with prolog before an error message (separated by newline)",
			input: "There is the prolog.\nError: an error message",
			expect: expect{
				Error:  "an error message",
				Search: "an error message",
			},
		},
		{
			name:  "Error message with several wrapped errors (prepended with \"Error:\")",
			input: "Error: level 1: level 2: level 3: level 4",
			expect: expect{
				Error:  "level 1: level 2: level 3: level 4",
				Search: "level 4",
				Wrap:   []string{"level 1", "level 2", "level 3", "level 4"},
			},
		},
		{
			name:  "Error message with newline and with errors in \"bullet\" list",
			input: "Error: foobar:\n  err1\n  err2\n  err3",
			expect: expect{
				Error:  "foobar: err1: err2: err3",
				Search: "err2",
				Wrap:   []string{"foobar", "err1", "err2", "err3"},
			},
		},
		{
			name: "Error message from 'podman pull' - unknown registry",
			input: `Trying to pull foobar.com/foo:latest...
Error: Error initializing source docker://foobar.com/foo:latest: error pinging docker registry foobar.com: Get "https://foobar.com/v2/": x509: certificate has expired or is not yet valid: current time 2021-07-04T00:45:46+02:00 is after 2019-09-26T23:59:59Z`,
			expect: expect{
				Error: "Error initializing source docker://foobar.com/foo:latest: error pinging docker registry foobar.com: Get \"https://foobar.com/v2/\": x509: certificate has expired or is not yet valid: current time 2021-07-04T00:45:46+02:00 is after 2019-09-26T23:59:59Z",
			},
		},
		{
			name: "Error message from 'podman pull' - unauthorized (Docker Hub)",
			input: `Trying to pull docker.io/library/foobar:latest...
Error: Error initializing source docker://foobar:latest: Error reading manifest latest in docker.io/library/foobar: errors:
denied: requested access to the resource is denied
unauthorized: authentication required`,
			expect: expect{
				Error: "Error initializing source docker://foobar:latest: Error reading manifest latest in docker.io/library/foobar: denied: requested access to the resource is denied: unauthorized: authentication required",
			},
		},
		{
			name: "Error message from 'podman pull' - unauthorized (Red Hat Registry)",
			input: `Trying to pull registry.redhat.io/foobar:latest...
Error: Error initializing source docker://registry.redhat.io/foobar:latest: unable to retrieve auth token: invalid username/password: unauthorized: Please login to the Red Hat Registry using your Customer Portal credentials. Further instructions can be found here: https://access.redhat.com/RegistryAuthentication
`,
			expect: expect{
				Error: "Error initializing source docker://registry.redhat.io/foobar:latest: unable to retrieve auth token: invalid username/password: unauthorized: Please login to the Red Hat Registry using your Customer Portal credentials. Further instructions can be found here: https://access.redhat.com/RegistryAuthentication",
			},
		},
		{
			name: "Error message from 'podman pull' - unknwon image (Fedora Registry)",
			input: `Trying to pull registry.fedoraproject.org/foobar:latest...
Error: Error initializing source docker://registry.fedoraproject.org/foobar:latest: Error reading manifest latest in registry.fedoraproject.org/foobar: manifest unknown: manifest unknown`,
			expect: expect{
				Error: "Error initializing source docker://registry.fedoraproject.org/foobar:latest: Error reading manifest latest in registry.fedoraproject.org/foobar: manifest unknown: manifest unknown",
			},
		},
		{
			name: "Error message from 'podman pull' - unknown image (Docker Hub)",
			input: `Trying to pull docker.io/library/busybox:foobar...
Error: Error initializing source docker://busybox:foobar: Error reading manifest foobar in docker.io/library/busybox: manifest unknown: manifest unknown`,
			expect: expect{
				Error: "Error initializing source docker://busybox:foobar: Error reading manifest foobar in docker.io/library/busybox: manifest unknown: manifest unknown",
			},
		},
		{
			name: "Error message from 'podman pull' - unknown image (Red Hat Registry)",
			input: `Trying to pull registry.redhat.io/foobar:latest...
time="2021-07-04T01:08:22+02:00" level=warning msg="failed, retrying in 1s ... (1/3). Error: Error initializing source docker://registry.redhat.io/foobar:latest: Error reading manifest latest in registry.redhat.io/foobar: unknown: Not Found"
time="2021-07-04T01:08:25+02:00" level=warning msg="failed, retrying in 1s ... (2/3). Error: Error initializing source docker://registry.redhat.io/foobar:latest: Error reading manifest latest in registry.redhat.io/foobar: unknown: Not Found"
time="2021-07-04T01:08:27+02:00" level=warning msg="failed, retrying in 1s ... (3/3). Error: Error initializing source docker://registry.redhat.io/foobar:latest: Error reading manifest latest in registry.redhat.io/foobar: unknown: Not Found"
Error: Error initializing source docker://registry.redhat.io/foobar:latest: Error reading manifest latest in registry.redhat.io/foobar: unknown: Not Found
`,
			expect: expect{
				Error: "Error initializing source docker://registry.redhat.io/foobar:latest: Error reading manifest latest in registry.redhat.io/foobar: unknown: Not Found",
			},
		},
		{
			name:  "Error message from 'podman login' - no such host",
			input: `Error: authenticating creds for "foobar": error pinging docker registry foobar: Get "https://foobar/v2/": dial tcp: lookup foobar: no such host`,
			expect: expect{
				Error: `authenticating creds for "foobar": error pinging docker registry foobar: Get "https://foobar/v2/": dial tcp: lookup foobar: no such host`,
			},
		},
		{
			name:  "Error message from 'podman login' - invalid username/password",
			input: `Error: error logging into "registry.redhat.io": invalid username/password`,
			expect: expect{
				Error: `error logging into "registry.redhat.io": invalid username/password`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := parseErrorMsg(bytes.NewBufferString(tc.input))

			if tc.expect.IsNil {
				assert.Nil(t, err)
				return
			} else {
				assert.NotNil(t, err)
			}

			errInternal := err.(*internalError)
			assert.Equal(t, tc.expect.Error, errInternal.Error())

			if tc.expect.Search != "" {
				assert.True(t, errInternal.Is(errors.New(tc.expect.Search)))
			}

			if len(tc.expect.Wrap) != 0 {
				for {
					assert.Equal(t, len(tc.expect.Wrap), len(errInternal.errors))

					for i, part := range tc.expect.Wrap {
						assert.Equal(t, part, errInternal.errors[i])
					}

					err = errInternal.Unwrap()
					if err == nil {
						assert.Equal(t, len(tc.expect.Wrap), 1)
						break
					}
					errInternal = err.(*internalError)
					tc.expect.Wrap = tc.expect.Wrap[1:]
				}
			}
		})
	}
}

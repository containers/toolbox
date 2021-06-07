/*
 * Copyright © 2021 Red Hat Inc.
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
			name:  "Empty",
			input: "",
			expect: expect{
				IsNil: true,
				Error: "",
			},
		},
		{
			name:  "Text with only a prolog and no error message",
			input: "There is only a prolog and no error message",
			expect: expect{
				IsNil: true,
				Error: "",
			},
		},
		{
			name:  "Text with only a prolog and no error message",
			input: "There is only a prolog Error: not an error message",
			expect: expect{
				IsNil: true,
				Error: "",
			},
		},
		{
			name: "Text with a prolog before the error message",
			input: `There is a prolog
Error: an error message`,
			expect: expect{
				Error:  "an error message",
				Search: "an error message",
			},
		},
		{
			name:  "Error message with several wrapped errors",
			input: "Error: level 1: level 2: level 3: level 4",
			expect: expect{
				Error:  "level 1: level 2: level 3: level 4",
				Search: "level 4",
				Wrap:   []string{"level 1", "level 2", "level 3", "level 4"},
			},
		},
		{
			name: "Error message with a bullet list",
			input: `Error: an error message:
  err1
  err2
  err3`,
			expect: expect{
				Error:  "an error message: err1: err2: err3",
				Search: "err2",
				Wrap:   []string{"an error message", "err1", "err2", "err3"},
			},
		},
		{
			name: "Error message from 'podman pull' - unauthorized (Docker Hub)",
			input: `Trying to pull docker.io/library/foobar:latest...
Error: Error initializing source docker://foobar:latest: Error reading manifest latest in docker.io/library/foobar: errors:
denied: requested access to the resource is denied
unauthorized: authentication required`,
			expect: expect{
				Error: "Error initializing source docker://foobar:latest: Error reading manifest latest in docker.io/library/foobar: errors: denied: requested access to the resource is denied: unauthorized: authentication required",
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

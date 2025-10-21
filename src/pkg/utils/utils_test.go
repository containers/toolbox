/*
 * Copyright © 2021 – 2025 Red Hat Inc.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageReferenceCanBeID(t *testing.T) {
	testCases := []struct {
		name string
		ref  string
		ok   bool
	}{
		{
			name: "Valid ID (random 6 chars)",
			ref:  "34afbc",
			ok:   true,
		},
		{
			name: "Valid ID (random 64 chars)",
			ref:  "8215cb84fa588215cb84fa588215cb84fa588215cb84fa588215cb84fa58fbca",
			ok:   true,
		},
		{
			name: "Valid ID (Podman short)",
			ref:  "8215cb84fa58",
			ok:   true,
		},
		{
			name: "Valid ID (Podman long)",
			ref:  "8b9affd1dbc261a7f586ed06a8fd993d09449a5ac79ebc7e80e86efdf3c223f6",
			ok:   true,
		},
		{
			name: "Invalid ID (random <6 chars)",
			ref:  "acbdf",
			ok:   false,
		},
		{
			name: "Invalid ID (random >64 chars)",
			ref:  "8215cb84fa588215cb84fa588215cb84fa588215cb84fa588215cb84fa58fbcab",
			ok:   false,
		},
		{
			name: "Invalid ID (image URI; whole alphabet + slash)",
			ref:  "fedoraproject.org/fedora-toolbox:32",
			ok:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok := ImageReferenceCanBeID(tc.ref)
			assert.Equal(t, tc.ok, ok)
		})
	}
}

func TestParseRelease(t *testing.T) {
	testCases := []struct {
		inputDistro  string
		inputRelease string
		output       string
		errMsg       string
	}{
		{
			inputDistro:  "arch",
			inputRelease: "",
			output:       "latest",
		},
		{
			inputDistro:  "arch",
			inputRelease: "latest",
			output:       "latest",
		},
		{
			inputDistro:  "arch",
			inputRelease: "rolling",
			output:       "latest",
		},
		{
			inputDistro:  "arch",
			inputRelease: "foo",
			errMsg:       "The release must be 'latest'.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "f34",
			output:       "34",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "33",
			output:       "33",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "-3",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "-3.0",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "-3.1",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "0",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "0.0",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "0.1",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "3.0",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "3.1",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "foo",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "fedora",
			inputRelease: "3foo",
			errMsg:       "The release must be a positive integer.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.0",
			output:       "8.0",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.3",
			output:       "8.3",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.42",
			output:       "8.42",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.0.0",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.0.1",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.3.0",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.3.1",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "foo",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "8.2foo",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "-2.1",
			errMsg:       "The release must be a positive number.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "-2.-1",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "rhel",
			inputRelease: "2.-1",
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "4.10",
			output:       "4.10",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "5.04",
			output:       "5.04",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "20.04",
			output:       "20.04",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "20.10",
			output:       "20.10",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "20",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "20.04.0",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "20.04.1",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "foo",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "20foo",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "foo.bar",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "foo.bar.baz",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "3.10",
			errMsg:       "The release year must be 4 or more.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "202.4",
			errMsg:       "The release year cannot have more than two digits.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "202.04",
			errMsg:       "The release year cannot have more than two digits.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "2020.4",
			errMsg:       "The release year cannot have more than two digits.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "2020.04",
			errMsg:       "The release year cannot have more than two digits.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "04.10",
			errMsg:       "The release year cannot have a leading zero.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "4.bar",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "4.bar.baz",
			errMsg:       "The release must be in the 'YY.MM' format.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "4.0",
			errMsg:       "The release month must be between 01 and 12.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "4.00",
			errMsg:       "The release month must be between 01 and 12.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "4.13",
			errMsg:       "The release month must be between 01 and 12.",
		},
		{
			inputDistro:  "ubuntu",
			inputRelease: "20.4",
			errMsg:       "The release month must have two digits.",
		},
	}

	for _, tc := range testCases {
		name := tc.inputDistro + ", " + tc.inputRelease
		t.Run(name, func(t *testing.T) {
			release, err := parseRelease(tc.inputDistro, tc.inputRelease)

			if tc.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errMsg)

				var errParseRelease *ParseReleaseError
				assert.ErrorAs(t, err, &errParseRelease)
			}

			assert.Equal(t, tc.output, release)
		})
	}
}

func TestPathExistsDoesNotExist(t *testing.T) {
	exists := PathExists("/does/not/exist")
	assert.False(t, exists)
}

func TestPathExistsExecutable(t *testing.T) {
	path, err := os.Executable()
	require.NoError(t, err)
	exists := PathExists(path)
	assert.True(t, exists)
}

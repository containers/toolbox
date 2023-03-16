/*
 * Copyright © 2021 – 2022 Red Hat Inc.
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
	"testing"

	"github.com/stretchr/testify/assert"
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
		name         string
		inputDistro  string
		inputRelease string
		output       string
		ok           bool
		errMsg       string
	}{
		{
			name:         "Fedora; f34; valid",
			inputDistro:  "fedora",
			inputRelease: "f34",
			output:       "34",
			ok:           true,
		},
		{
			name:         "Fedora; 33; valid",
			inputDistro:  "fedora",
			inputRelease: "33",
			output:       "33",
			ok:           true,
		},
		{
			name:         "Fedora; -3; invalid; less than 0",
			inputDistro:  "fedora",
			inputRelease: "-3",
			ok:           false,
			errMsg:       "The release must be a positive integer.",
		},
		{
			name:         "Fedora; foo; invalid; non-numeric",
			inputDistro:  "fedora",
			inputRelease: "foo",
			ok:           false,
			errMsg:       "The release must be a positive integer.",
		},
		{
			name:         "RHEL; 8.3; valid",
			inputDistro:  "rhel",
			inputRelease: "8.3",
			output:       "8.3",
			ok:           true,
		},
		{
			name:         "RHEL; 8.42; valid",
			inputDistro:  "rhel",
			inputRelease: "8.42",
			output:       "8.42",
			ok:           true,
		},
		{
			name:         "RHEL; 8; invalid; missing point release",
			inputDistro:  "rhel",
			inputRelease: "8",
			ok:           false,
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			name:         "RHEL; 8.2foo; invalid; non-float",
			inputDistro:  "rhel",
			inputRelease: "8.2foo",
			ok:           false,
			errMsg:       "The release must be in the '<major>.<minor>' format.",
		},
		{
			name:         "RHEL; -2.1; invalid; less than 0",
			inputDistro:  "rhel",
			inputRelease: "-2.1",
			ok:           false,
			errMsg:       "The release must be a positive number.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			release, err := parseRelease(tc.inputDistro, tc.inputRelease)

			if tc.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				if tc.errMsg != "" {
					assert.EqualError(t, err, tc.errMsg)
				}
			}

			assert.Equal(t, tc.output, release)
		})
	}
}

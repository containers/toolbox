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

package utils

import (
	"io/ioutil"
	"os"
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

func TestPathExists(t *testing.T) {
	testCases := []struct {
		name                string
		targetExists        bool
		targetIsSymlink     bool
		targetIsDir         bool
		symlinkTargetExists bool
		ok                  bool
	}{
		{
			name:         "Target does not exist",
			targetExists: false,
			ok:           false,
		},
		{
			name:         "Target exists and is not a symlink",
			targetExists: true,
			ok:           true,
		},
		{
			name:                "Target exists and is a symlink; symlink target does exist",
			targetExists:        true,
			targetIsSymlink:     true,
			symlinkTargetExists: true,
			ok:                  true,
		},
		{
			name:                "Target exists and is a symlink; symlink target does not exist",
			targetExists:        true,
			targetIsSymlink:     true,
			symlinkTargetExists: false,
			ok:                  true,
		},
		{
			name:         "Target exists and is a directory",
			targetExists: true,
			targetIsDir:  true,
			ok:           true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var targetName string

			targetName = "nonexistent-name"

			if tc.targetIsDir && tc.targetIsSymlink {
				t.Fatal("Target can't be a symlink and a directory at the same time")
			}

			if tc.targetExists {
				if !tc.targetIsSymlink && !tc.targetIsDir {
					f, err := ioutil.TempFile("", "")
					assert.NoError(t, err)
					targetName = f.Name()
					defer f.Close()
					defer os.Remove(targetName)
				}

				if tc.targetIsSymlink || tc.targetIsDir {
					dir, err := ioutil.TempDir("", "")
					assert.NoError(t, err)
					targetName = dir
					defer os.Remove(targetName)

					if tc.targetIsSymlink {
						target, err := ioutil.TempFile(dir, "")
						assert.NoError(t, err)
						defer target.Close()
						defer os.Remove(target.Name())

						targetName += "/symlink"
						err = os.Symlink(target.Name(), targetName)
						assert.NoError(t, err)

						if !tc.symlinkTargetExists {
							target.Close()
							os.Remove(target.Name())
						}
					}
				}
			}

			assert.Equal(t, tc.ok, PathExists(targetName))
		})
	}
}

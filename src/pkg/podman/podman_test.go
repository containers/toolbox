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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckVersion(t *testing.T) {
	testCases := []struct {
		name           string
		setVersion     string
		checkedVersion string
		ok             bool
	}{
		{
			name:           "Podman 2.1; Check for 2.0",
			setVersion:     "2.1",
			checkedVersion: "2.0",
			ok:             true,
		},
		{
			name:           "Podman 2.1-dev; Check for 2.1",
			setVersion:     "2.1-dev",
			checkedVersion: "2.1",
			ok:             true,
		},
		{
			name:           "Podman 2.3.1; Check for 2.2.8",
			setVersion:     "2.3.1",
			checkedVersion: "2.2.8",
			ok:             true,
		},
		{
			name:           "Podman 3.0-rc3; Check for 3.0-rc2",
			setVersion:     "3.0-rc3",
			checkedVersion: "3.0-rc2",
			ok:             true,
		},
		{
			name:           "Podman 2.1-rc3; Check for 2.1-rc4",
			setVersion:     "2.1-rc3",
			checkedVersion: "2.1-rc4",
			ok:             false,
		},
		{
			name:           "Podman 2.2.3; Check for 2.2.4",
			setVersion:     "2.2.3",
			checkedVersion: "2.2.4",
			ok:             false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			podmanVersion = tc.setVersion
			assert.Equal(t, tc.ok, CheckVersion(tc.checkedVersion))
		})
	}
}

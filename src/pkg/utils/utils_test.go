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
	"strconv"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetReleaseFormat(t *testing.T) {
	testCases := []struct {
		name     string
		distro   string
		expected string
	}{
		{
			"Unknown distro",
			"foobar",
			"",
		},
		{
			"Known distro (fedora)",
			"fedora",
			supportedDistros["fedora"].ReleaseFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := GetReleaseFormat(tc.distro)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestGetSupportedDistros(t *testing.T) {
	refDistros := []string{"fedora", "rhel"}

	distros := GetSupportedDistros()
	for _, d := range distros {
		assert.Contains(t, refDistros, d)
	}
}

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

func TestIsDistroSupport(t *testing.T) {
	testCases := []struct {
		name   string
		distro string
		ok     bool
	}{
		{
			"Unsupported distro",
			"foobar",
			false,
		},
		{
			"Supported distro (fedora)",
			"fedora",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := IsDistroSupported(tc.distro)
			assert.Equal(t, tc.ok, res)
		})
	}
}

func TestResolveDistro(t *testing.T) {
	testCases := []struct {
		name        string
		distro      string
		expected    string
		configValue string
		err         bool
	}{
		{
			"Default - no distro provided; config unset",
			"",
			distroDefault,
			"",
			false,
		},
		{
			"Default - no distro provided; config set",
			"",
			"rhel",
			"rhel",
			false,
		},
		{
			"Fedora",
			"fedora",
			"fedora",
			"",
			false,
		},
		{
			"RHEL",
			"rhel",
			"rhel",
			"",
			false,
		},
		{
			"FooBar; wrong distro",
			"foobar",
			"",
			"",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.configValue != "" {
				viper.Set("general.distro", tc.configValue)
			}

			res, err := ResolveDistro(tc.distro)
			assert.Equal(t, tc.expected, res)
			if tc.err {
				assert.NotNil(t, err)
			}
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
		err          error
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
			errMsg:       "release must be a positive integer",
		},
		{
			name:         "Fedora; foo; invalid; non-numeric",
			inputDistro:  "fedora",
			inputRelease: "foo",
			ok:           false,
			err:          strconv.ErrSyntax,
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
			errMsg:       "release must have a '.'",
		},
		{
			name:         "RHEL; 8.2foo; invalid; non-float",
			inputDistro:  "rhel",
			inputRelease: "8.2foo",
			ok:           false,
			err:          strconv.ErrSyntax,
		},
		{
			name:         "RHEL; -2.1; invalid; less than 0",
			inputDistro:  "rhel",
			inputRelease: "-2.1",
			ok:           false,
			errMsg:       "release must be a positive number",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			release, err := ParseRelease(tc.inputDistro, tc.inputRelease)

			if tc.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)

				if tc.err != nil {
					assert.ErrorIs(t, err, tc.err)
				}

				if tc.errMsg != "" {
					assert.EqualError(t, err, tc.errMsg)
				}
			}

			assert.Equal(t, tc.output, release)
		})
	}
}

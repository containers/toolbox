/*
 * Copyright © 2019 – 2026 Red Hat Inc.
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

package architecture

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToEscapedString(t *testing.T) {
	testCases := []struct {
		name   string
		input  []byte
		expect string
	}{
		{
			name:   "ELF magic bytes",
			input:  []byte{0x7f, 0x45, 0x4c, 0x46},
			expect: `\x7f\x45\x4c\x46`,
		},
		{
			name:   "ELF magic bytes (char)",
			input:  []byte{0x7f, 'E', 'L', 'F'},
			expect: `\x7f\x45\x4c\x46`,
		},
		{
			name:   "single zero byte",
			input:  []byte{0x00},
			expect: `\x00`,
		},
		{
			name:   "single 0xff byte",
			input:  []byte{0xff},
			expect: `\xff`,
		},
		{
			name:   "empty input",
			input:  []byte{},
			expect: "",
		},
		{
			name:   "nil input",
			input:  nil,
			expect: "",
		},
		{
			name:   "all zeros",
			input:  []byte{0x00, 0x00, 0x00},
			expect: `\x00\x00\x00`,
		},
		{
			name:   "all 0xff",
			input:  []byte{0xff, 0xff, 0xff},
			expect: `\xff\xff\xff`,
		},
		{
			name: "aarch64 ELFMagic full",
			input: []byte{
				0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x02, 0x00, 0xb7, 0x00,
			},
			expect: `\x7f\x45\x4c\x46\x02\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\xb7\x00`,
		},
		{
			name: "aarch64 ELFMask full",
			input: []byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xfe, 0xff, 0xff, 0xff,
			},
			expect: `\xff\xff\xff\xff\xff\xff\xff\x00\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff\xff`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := bytesToEscapedString(tc.input)
			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestBuildRegistrationString(t *testing.T) {
	testCases := []struct {
		name   string
		reg    Registration
		expect string
	}{
		{
			name: "standard aarch64 registration",
			reg: Registration{
				Name:        "qemu-aarch64",
				MagicType:   "M",
				Offset:      "0",
				Magic:       []byte{0x7f, 0x45},
				Mask:        []byte{0xff, 0xff},
				Interpreter: "/run/host/usr/bin/qemu-aarch64-static",
				Flags:       "FC",
			},
			expect: ":qemu-aarch64:M:0:\\x7f\\x45:\\xff\\xff:/run/host/usr/bin/qemu-aarch64-static:FC",
		},
		{
			name: "registration with empty magic and mask",
			reg: Registration{
				Name:        "test",
				MagicType:   "M",
				Offset:      "0",
				Magic:       []byte{},
				Mask:        []byte{},
				Interpreter: "/usr/bin/test",
				Flags:       "F",
			},
			expect: ":test:M:0:::/usr/bin/test:F",
		},
		{
			name: "registration with nil magic and mask",
			reg: Registration{
				Name:        "test",
				MagicType:   "M",
				Offset:      "0",
				Magic:       nil,
				Mask:        nil,
				Interpreter: "/usr/bin/test",
				Flags:       "F",
			},
			expect: ":test:M:0:::/usr/bin/test:F",
		},
		{
			name: "registration with all empty strings",
			reg: Registration{
				Magic: []byte{},
				Mask:  []byte{},
			},
			expect: ":::::::",
		},
		{
			name: "registration with non-default offset",
			reg: Registration{
				Name:        "custom",
				MagicType:   "M",
				Offset:      "16",
				Magic:       []byte{0x7f, 0x45, 0x4c, 0x46},
				Mask:        []byte{0xfe, 0xff, 0xff, 0xff},
				Interpreter: "/run/host/usr/bin/qemu-aarch64",
				Flags:       "FC",
			},
			expect: ":custom:M:16:\\x7f\\x45\\x4c\\x46:\\xfe\\xff\\xff\\xff:/run/host/usr/bin/qemu-aarch64:FC",
		},
		{
			name: "registration with flags only F",
			reg: Registration{
				Name:        "qemu-x86_64",
				MagicType:   "M",
				Offset:      "0",
				Magic:       []byte{0x7f, 0x45},
				Mask:        []byte{0xff, 0xff},
				Interpreter: "/run/host/usr/bin/qemu-x86_64-static",
				Flags:       "F",
			},
			expect: ":qemu-x86_64:M:0:\\x7f\\x45:\\xff\\xff:/run/host/usr/bin/qemu-x86_64-static:F",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.reg.buildRegistrationString()
			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestGetDefaultRegistration(t *testing.T) {
	testCases := []struct {
		name            string
		archID          int
		interpreterPath string
		expectNil       bool
		expectName      string
		expectInterp    string
		expectFlags     string
		expectMagicType string
		expectOffset    string
	}{
		{
			name:            "aarch64 with absolute path",
			archID:          Aarch64,
			interpreterPath: "/usr/bin/qemu-aarch64-static",
			expectName:      "qemu-aarch64",
			expectInterp:    "/run/host/usr/bin/qemu-aarch64-static",
			expectFlags:     defaultFlags,
			expectMagicType: defaultMagicType,
			expectOffset:    defaultOffset,
		},
		{
			name:            "aarch64 with -run-host prefix already present",
			archID:          Aarch64,
			interpreterPath: "/run/host/usr/bin/qemu-aarch64-static",
			expectName:      "qemu-aarch64",
			expectInterp:    "/run/host/usr/bin/qemu-aarch64-static",
			expectFlags:     defaultFlags,
			expectMagicType: defaultMagicType,
			expectOffset:    defaultOffset,
		},
		{
			name:            "x86_64 with absolute path",
			archID:          X86_64,
			interpreterPath: "/usr/bin/qemu-x86_64-static",
			expectName:      "qemu-x86_64",
			expectInterp:    "/run/host/usr/bin/qemu-x86_64-static",
			expectFlags:     defaultFlags,
			expectMagicType: defaultMagicType,
			expectOffset:    defaultOffset,
		},
		{
			name:            "ppc64le with absolute path",
			archID:          Ppc64le,
			interpreterPath: "/usr/bin/qemu-ppc64le-static",
			expectName:      "qemu-ppc64le",
			expectInterp:    "/run/host/usr/bin/qemu-ppc64le-static",
			expectFlags:     defaultFlags,
			expectMagicType: defaultMagicType,
			expectOffset:    defaultOffset,
		},
		{
			name:            "invalid arch returns nil",
			archID:          999,
			interpreterPath: "/usr/bin/qemu-fake",
			expectNil:       true,
		},
		{
			name:            "NotSpecified returns nil",
			archID:          NotSpecified,
			interpreterPath: "/usr/bin/qemu-fake",
			expectNil:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reg := getDefaultRegistration(tc.archID, tc.interpreterPath)

			if tc.expectNil {
				assert.Nil(t, reg)
				return
			}

			assert.NotNil(t, reg)
			assert.Equal(t, tc.expectName, reg.Name)
			assert.Equal(t, tc.expectInterp, reg.Interpreter)
			assert.Equal(t, tc.expectFlags, reg.Flags)
			assert.Equal(t, tc.expectMagicType, reg.MagicType)
			assert.Equal(t, tc.expectOffset, reg.Offset)

			arch, _ := getArchitecture(tc.archID)
			assert.Equal(t, arch.ELFMagic, reg.Magic)
			assert.Equal(t, arch.ELFMask, reg.Mask)
		})
	}
}

func TestGetDefaultRegistrationInterpreterPathPrefixing(t *testing.T) {
	testCases := []struct {
		name            string
		interpreterPath string
		expectInterp    string
	}{
		{
			name:            "absolute path gets -run-host prefix",
			interpreterPath: "/usr/bin/qemu-aarch64-static",
			expectInterp:    "/run/host/usr/bin/qemu-aarch64-static",
		},
		{
			name:            "-run-host- prefix is not duplicated",
			interpreterPath: "/run/host/usr/bin/qemu-aarch64-static",
			expectInterp:    "/run/host/usr/bin/qemu-aarch64-static",
		},
		{
			name:            "nested path",
			interpreterPath: "/opt/custom/qemu/bin/qemu-aarch64-static",
			expectInterp:    "/run/host/opt/custom/qemu/bin/qemu-aarch64-static",
		},
		{
			name:            "-run-host without trailing slash gets prefixed",
			interpreterPath: "/run/host",
			expectInterp:    "/run/host/run/host",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reg := getDefaultRegistration(Aarch64, tc.interpreterPath)
			assert.NotNil(t, reg)
			assert.Equal(t, tc.expectInterp, reg.Interpreter)
		})
	}
}

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
	"bytes"
	"debug/elf"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseArgArchValue(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		expect int
		errMsg string
	}{
		{
			name:   "aarch64",
			input:  "aarch64",
			expect: Aarch64,
		},
		{
			name:   "arm64 alias for aarch64",
			input:  "arm64",
			expect: Aarch64,
		},
		{
			name:   "x86_64",
			input:  "x86_64",
			expect: X86_64,
		},
		{
			name:   "amd64 alias for x86_64",
			input:  "amd64",
			expect: X86_64,
		},
		{
			name:   "ppc64le",
			input:  "ppc64le",
			expect: Ppc64le,
		},
		{
			name:   "unsupported architecture",
			input:  "mips",
			errMsg: "architecture 'mips' is not supported by Toolbx",
		},
		{
			name:   "empty string",
			input:  "",
			errMsg: "architecture '' is not supported by Toolbx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseArgArchValue(tc.input)

			if tc.errMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errMsg)
				assert.Equal(t, NotSpecified, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, result)
			}
		})
	}
}

func TestGetArchNameOCI(t *testing.T) {
	testCases := []struct {
		name   string
		archID int
		expect string
	}{
		{
			name:   "aarch64 returns arm64",
			archID: Aarch64,
			expect: "arm64",
		},
		{
			name:   "x86_64 returns amd64",
			archID: X86_64,
			expect: "amd64",
		},
		{
			name:   "ppc64le returns ppc64le",
			archID: Ppc64le,
			expect: "ppc64le",
		},
		{
			name:   "NotSpecified returns fallback string",
			archID: NotSpecified,
			expect: "arch_not_specified",
		},
		{
			name:   "invalid arch ID returns empty",
			archID: 999,
			expect: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetArchNameOCI(tc.archID)
			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestGetArchNameBinfmt(t *testing.T) {
	testCases := []struct {
		name   string
		archID int
		expect string
	}{
		{
			name:   "aarch64",
			archID: Aarch64,
			expect: "aarch64",
		},
		{
			name:   "x86_64",
			archID: X86_64,
			expect: "x86_64",
		},
		{
			name:   "ppc64le",
			archID: Ppc64le,
			expect: "ppc64le",
		},
		{
			name:   "NotSpecified returns fallback string",
			archID: NotSpecified,
			expect: "arch_not_specified",
		},
		{
			name:   "invalid arch ID returns empty",
			archID: 999,
			expect: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getArchNameBinfmt(tc.archID)
			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestGetArchitecture(t *testing.T) {
	testCases := []struct {
		name   string
		archID int
		exists bool
	}{
		{
			name:   "aarch64 exists",
			archID: Aarch64,
			exists: true,
		},
		{
			name:   "x86_64 exists",
			archID: X86_64,
			exists: true,
		},
		{
			name:   "ppc64le exists",
			archID: Ppc64le,
			exists: true,
		},
		{
			name:   "NotSpecified does not exist",
			archID: NotSpecified,
			exists: false,
		},
		{
			name:   "invalid arch ID does not exist",
			archID: 999,
			exists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			arch, exists := getArchitecture(tc.archID)
			assert.Equal(t, tc.exists, exists)

			if tc.exists {
				assert.Equal(t, tc.archID, arch.ID)
				assert.NotEmpty(t, arch.NameBinfmt)
				assert.NotEmpty(t, arch.NameOCI)
				assert.NotEmpty(t, arch.Aliases)
				assert.NotEmpty(t, arch.ELFMagic)
				assert.NotEmpty(t, arch.ELFMask)
			}
		})
	}
}

func TestImageReferenceGetArchFromTag(t *testing.T) {
	testCases := []struct {
		name   string
		image  string
		expect int
	}{
		{
			name:   "image with aarch64 binfmt arch suffix",
			image:  "registry.fedoraproject.org/fedora-toolbox:41-aarch64",
			expect: Aarch64,
		},
		{
			name:   "image with arm64 OCI arch suffix",
			image:  "registry.fedoraproject.org/fedora-toolbox:41-arm64",
			expect: Aarch64,
		},
		{
			name:   "image with x86_64 binfmt arch suffix",
			image:  "registry.fedoraproject.org/fedora-toolbox:41-x86_64",
			expect: X86_64,
		},
		{
			name:   "image with amd64 OCI arch suffix",
			image:  "registry.fedoraproject.org/fedora-toolbox:41-amd64",
			expect: X86_64,
		},
		{
			name:   "image with ppc64le arch suffix",
			image:  "registry.fedoraproject.org/fedora-toolbox:41-ppc64le",
			expect: Ppc64le,
		},
		{
			name:   "image without arch suffix",
			image:  "registry.fedoraproject.org/fedora-toolbox:41",
			expect: NotSpecified,
		},
		{
			name:   "image without tag",
			image:  "registry.fedoraproject.org/fedora-toolbox",
			expect: NotSpecified,
		},
		{
			name:   "image with unknown arch suffix",
			image:  "registry.fedoraproject.org/fedora-toolbox:41-mips",
			expect: NotSpecified,
		},
		{
			name:   "empty image reference",
			image:  "",
			expect: NotSpecified,
		},
		{
			name:   "ubuntu image with amd64 OCI arch suffix",
			image:  "quay.io/toolbx-images/ubuntu-toolbox:22.04-amd64",
			expect: X86_64,
		},
		{
			name:   "ubuntu image with x86_64 binfmt arch suffix",
			image:  "quay.io/toolbx-images/ubuntu-toolbox:22.04-x86_64",
			expect: X86_64,
		},
		{
			name:   "ubuntu image with arm64 OCI arch suffix",
			image:  "quay.io/toolbx-images/ubuntu-toolbox:22.04-arm64",
			expect: Aarch64,
		},
		{
			name:   "ubuntu image with aarch64 binfmt arch suffix",
			image:  "quay.io/toolbx-images/ubuntu-toolbox:22.04-aarch64",
			expect: Aarch64,
		},
		{
			name:   "ubuntu image with unsupported arch suffix",
			image:  "quay.io/toolbx-images/ubuntu-toolbox:22.04-mips",
			expect: NotSpecified,
		},
		{
			name:   "ubuntu image without arch suffix",
			image:  "quay.io/toolbx-images/ubuntu-toolbox:22.04",
			expect: NotSpecified,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ImageReferenceGetArchFromTag(tc.image)
			assert.Equal(t, tc.expect, result)
		})
	}
}

// buildMinimalELF64 constructs a valid ELF64 binary with the given program header types.
// A binary without PT_INTERP is statically linked; one with PT_INTERP is dynamically linked.
func buildMinimalELF64(progTypes []elf.ProgType) []byte {
	var buf bytes.Buffer

	phnum := uint16(len(progTypes))
	phoff := uint64(64)

	hdr := elf.Header64{
		Ident:     [16]byte{0x7f, 'E', 'L', 'F', 2, 1, 1},
		Type:      uint16(elf.ET_EXEC),
		Machine:   uint16(elf.EM_X86_64),
		Version:   1,
		Phoff:     phoff,
		Ehsize:    64,
		Phentsize: 56,
		Phnum:     phnum,
	}

	binary.Write(&buf, binary.LittleEndian, &hdr)

	for _, pt := range progTypes {
		prog := elf.Prog64{
			Type: uint32(pt),
		}
		binary.Write(&buf, binary.LittleEndian, &prog)
	}

	return buf.Bytes()
}

func TestIsStaticallyLinkedELF(t *testing.T) {
	testCases := []struct {
		name    string
		content []byte
		expect  bool
	}{
		{
			name:    "statically linked ELF (PT_LOAD only)",
			content: buildMinimalELF64([]elf.ProgType{elf.PT_LOAD}),
			expect:  true,
		},
		{
			name:    "dynamically linked ELF (has PT_INTERP)",
			content: buildMinimalELF64([]elf.ProgType{elf.PT_INTERP, elf.PT_LOAD}),
			expect:  false,
		},
		{
			name:    "not an ELF file",
			content: []byte("this is not an ELF binary"),
			expect:  false,
		},
		{
			name:    "empty file",
			content: []byte{},
			expect:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			file := filepath.Join(dir, "test-binary")
			err := os.WriteFile(file, tc.content, 0755)
			require.NoError(t, err)

			result := isStaticallyLinkedELF(file)
			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestIsStaticallyLinkedELFFileDoesNotExist(t *testing.T) {
	result := isStaticallyLinkedELF("/does/not/exist")
	assert.False(t, result)
}

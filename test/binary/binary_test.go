/*
 * Copyright © 2022 Ondřej Míchal
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

package binary

import (
	"debug/elf"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	binaryPath         string
	buildOverridesPath string

	runpath       string
	dynamicLinker string

	toolbxFile *os.File
	toolbxElf  *elf.File
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [BINARY] [BUILD OVERRIDES FILE]\n", os.Args[0])
}

func TestMain(m *testing.M) {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "%s: wrong number of arguments\n", os.Args[0])
		usage()
		os.Exit(1)
	}

	binaryPath = filepath.Clean(args[0])
	buildOverridesPath = filepath.Clean(args[1])
	content, err := ioutil.ReadFile(buildOverridesPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to read file %s: %s\n", os.Args[0], buildOverridesPath, err)
	}

	buildOverrides := strings.Split(strings.TrimSpace(string(content)), " ")
	if len(buildOverrides) != 2 {
		fmt.Fprintf(os.Stderr, "%s: overrides file needs to have 2 space-separated values written inside", os.Args[0])
		os.Exit(1)
	}
	dynamicLinker = buildOverrides[0]
	runpath = buildOverrides[1]

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s: file %s does not exist\n", os.Args[0], binaryPath)
		os.Exit(1)
	}

	toolbxFile, err = os.Open(binaryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to open tested binary: %s", os.Args[0], err)
		os.Exit(1)
	}

	toolbxElf, err = elf.NewFile(toolbxFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: failed to prepare ELF file: %s", os.Args[0], err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestRunPath(t *testing.T) {
	testedRunpathSlice, err := toolbxElf.DynString(elf.DT_RUNPATH)
	assert.NoError(t, err)
	assert.Len(t, testedRunpathSlice, 1)

	testedRunpath := testedRunpathSlice[0]
	assert.Equal(t, testedRunpath, runpath)
}

func TestDynamicLinker(t *testing.T) {
	var pt_interp *elf.Prog = nil
	for _, prog := range toolbxElf.Progs {
		if prog.Type == elf.PT_INTERP {
			pt_interp = prog
			break
		}
	}
	assert.NotNil(t, pt_interp)

	buf := make([]byte, pt_interp.Memsz-pt_interp.Align)
	i, err := toolbxFile.ReadAt(buf, int64(pt_interp.Off))
	assert.NoError(t, err)
	assert.LessOrEqual(t, i, len(buf))

	testedDynamicLinker := string(buf)
	assert.Equal(t, dynamicLinker, testedDynamicLinker)
}

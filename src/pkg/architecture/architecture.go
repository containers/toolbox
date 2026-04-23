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
	"debug/elf"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Architecture struct {
	ID         int
	NameBinfmt string
	NameOCI    string
	Aliases    []string
	ELFMagic   []byte
	ELFMask    []byte

	BinfmtFlags     string
	BinfmtName      string
	BinfmtMagicType string
	BinfmtOffset    string
}

type Config struct {
	ID               int
	QemuEmulatorPath string
}

const (
	NotSpecified = iota
	Aarch64
	Ppc64le
	X86_64
)

var supportedArchitectures = map[int]Architecture{
	Aarch64: {
		ID:         Aarch64,
		NameBinfmt: "aarch64",
		NameOCI:    "arm64",
		Aliases:    []string{"aarch64", "arm64"},
		ELFMagic:   []byte{0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0xb7, 0x00},
		ELFMask:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff},
	},
	Ppc64le: {
		ID:         Ppc64le,
		NameBinfmt: "ppc64le",
		NameOCI:    "ppc64le",
		Aliases:    []string{"ppc64le"},
		ELFMagic:   []byte{0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x15, 0x00},
		ELFMask:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0x00},
	},
	X86_64: {
		ID:         X86_64,
		NameBinfmt: "x86_64",
		NameOCI:    "amd64",
		Aliases:    []string{"x86_64", "amd64"},
		ELFMagic:   []byte{0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x3e, 0x00},
		ELFMask:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0xfe, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff},
	},
}

var (
	HostArchID             int
	supportedArgArchValues map[string]int
)

func init() {
	supportedArgArchValues = make(map[string]int)
	for archID, arch := range supportedArchitectures {
		for _, alias := range arch.Aliases {
			supportedArgArchValues[alias] = archID
		}
	}

	HostArchID, _ = ParseArgArchValue(runtime.GOARCH)
}

func GetArchConfigDefault() Config {
	return Config{
		ID:               HostArchID,
		QemuEmulatorPath: "",
	}
}

func getArchitecture(archID int) (Architecture, bool) {
	arch, exists := supportedArchitectures[archID]
	return arch, exists
}

func getArchNameBinfmt(arch int) string {
	if arch == NotSpecified {
		logrus.Warnf("Getting arch name for not specified architecture")
		return "arch_not_specified"
	}
	if archObj, exists := supportedArchitectures[arch]; exists {
		return archObj.NameBinfmt
	}
	return ""
}

func GetArchNameOCI(arch int) string {
	if arch == NotSpecified {
		logrus.Warnf("Getting arch name for not specified architecture")
		return "arch_not_specified"
	}
	if archObj, exists := supportedArchitectures[arch]; exists {
		return archObj.NameOCI
	}
	return ""
}

func HasContainerNativeArch(archID int) bool {
	return archID == HostArchID
}

func ImageReferenceGetArchFromTag(image string) int {
	tag := utils.ImageReferenceGetTag(image)

	if tag == "" {
		return NotSpecified
	}

	i := strings.LastIndexByte(tag, '-')
	if i == -1 {
		return NotSpecified
	}

	archInTag := tag[i+1:]

	for archID, arch := range supportedArchitectures {
		if arch.NameBinfmt == archInTag || arch.NameOCI == archInTag {
			return archID
		}
	}

	return NotSpecified
}

func IsArchSupportedOnCreation(archID int) (string, error) {
	archName := getArchNameBinfmt(archID)
	archNameDebug := GetArchNameOCI(archID)
	logrus.Debugf("Checking QEMU emulation support for architecture %s", archNameDebug)

	qemuBinaryPossibleNames := []string{
		fmt.Sprintf("qemu-%s-static", archName),
		fmt.Sprintf("qemu-%s", archName),
	}

	foundQemuBinaryPath := ""
	for _, qemuName := range qemuBinaryPossibleNames {
		qemuBinaryPath, err := exec.LookPath(qemuName)

		if err != nil {
			if errors.Is(err, exec.ErrNotFound) {
				continue
			}

			return "", fmt.Errorf("failed to look up binary '%s': %w", qemuName, err)
		}

		if isStaticallyLinkedELF(qemuBinaryPath) {
			foundQemuBinaryPath = qemuBinaryPath
			break
		}
	}

	if foundQemuBinaryPath == "" {
		err := fmt.Errorf("The host system does not have the required support: No %s statically linked QEMU emulator binary found", archNameDebug)
		return "", err
	}

	if !validateBinfmtRegistration(archID, false) {
		err := fmt.Errorf("The host system does not have the required support: No %s binfmt_misc registration found", archNameDebug)
		return "", err
	}

	return foundQemuBinaryPath, nil
}

func IsArchSupportedOnInitialization(archID int, interpreterPath string) error {
	archName := getArchNameBinfmt(archID)
	archNameDebug := GetArchNameOCI(archID)
	logrus.Debugf("Checking QEMU emulation support for architecture %s", archNameDebug)

	if isStaticallyLinkedELF(interpreterPath) {
		if !validateBinfmtRegistration(archID, true) {
			return fmt.Errorf("The host system does not have the required support: No %s binfmt_misc registration found", archNameDebug)
		}
		return nil
	}

	// Fallback: check standard locations on the host
	logrus.Debugf("Interpreter at %s not found or not statically linked, checking fallback locations in '/run/host/usr/bin/'", interpreterPath)
	fmt.Fprintf(os.Stderr, "Warning: QEMU emulator not found at expected path '%s', using fallback at '/run/host/usr/bin/'\n", interpreterPath)

	qemuBinaryPossiblePaths := []string{
		fmt.Sprintf("/run/host/usr/bin/qemu-%s-static", archName),
		fmt.Sprintf("/run/host/usr/bin/qemu-%s", archName),
	}

	for _, qemuPath := range qemuBinaryPossiblePaths {
		if isStaticallyLinkedELF(qemuPath) {
			logrus.Debugf("Found valid QEMU binary at %s", qemuPath)

			if !validateBinfmtRegistration(archID, true) {
				return fmt.Errorf("The host system does not have the required support: No %s binfmt_misc registration found", archNameDebug)
			}
			return nil
		}
	}

	return fmt.Errorf("The host system does not have the required support: No %s statically linked QEMU emulator binary found", archNameDebug)
}

func isStaticallyLinkedELF(filePath string) bool {
	if !utils.PathExists(filePath) {
		logrus.Debugf("File '%s' does not exist\n", filePath)
		return false
	}

	f, err := elf.Open(filePath)
	if err != nil {
		logrus.Debugf("File '%s' is not an ELF file\n", filePath)
		return false
	}
	defer f.Close()

	// Check for PT_INTERP program header
	for _, prog := range f.Progs {
		if prog.Type == elf.PT_INTERP {
			// Dynamically linked
			logrus.Debugf("File '%s' is dynamically linked\n", filePath)
			return false
		}
	}

	// Statically linked
	return true
}

func ParseArgArchValue(value string) (int, error) {
	archID, exists := supportedArgArchValues[value]
	if !exists {
		return NotSpecified, fmt.Errorf("architecture '%s' is not supported by Toolbx", value)
	}

	return archID, nil
}

func validateBinfmtRegistration(archID int, withinContainer bool) bool {
	archName := getArchNameBinfmt(archID)
	inContainerPathPrefix := ""

	if withinContainer {
		inContainerPathPrefix = "/run/host"
	}

	qemuBinfmtPossiblePaths := []string{
		fmt.Sprintf("%s/proc/sys/fs/binfmt_misc/qemu-%s", inContainerPathPrefix, archName),
		fmt.Sprintf("%s/proc/sys/fs/binfmt_misc/qemu-%s-static", inContainerPathPrefix, archName),
	}

	for _, binfmtPath := range qemuBinfmtPossiblePaths {
		if utils.PathExists(binfmtPath) {
			logrus.Debugf("Architecture %s is supported", archName)
			return true
		}
	}
	return false
}

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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/toolbox/pkg/shell"
	"github.com/sirupsen/logrus"
)

type Registration struct {
	Name        string
	MagicType   string
	Offset      string
	Magic       []byte
	Mask        []byte
	Interpreter string
	Flags       string
}

const (
	defaultMagicType = "M"
	defaultFlags     = "FC"
	defaultOffset    = "0"
	binfmtMiscPath   = "/proc/sys/fs/binfmt_misc"
)

func (r *Registration) buildRegistrationString() string {
	return fmt.Sprintf(":%s:%s:%s:%s:%s:%s:%s",
		r.Name, r.MagicType, r.Offset,
		bytesToEscapedString(r.Magic),
		bytesToEscapedString(r.Mask),
		r.Interpreter, r.Flags)
}

func (r *Registration) register() error {
	logrus.Debugf("Registering binfmt_misc for %s", r.Name)

	regString := r.buildRegistrationString()
	logrus.Debugf("Registration string: %s", regString)

	if err := os.WriteFile(filepath.Join(binfmtMiscPath, "register"), []byte(regString), 0200); err != nil {
		return fmt.Errorf("failed to register binfmt_misc handler: %w", err)
	}
	return nil
}

func bytesToEscapedString(bytes []byte) string {
	var result strings.Builder
	for _, b := range bytes {
		result.WriteString(fmt.Sprintf("\\x%02x", b))
	}
	return result.String()
}

func getDefaultRegistration(archID int, interpreterPath string) *Registration {
	arch, exists := getArchitecture(archID)
	if !exists {
		return nil
	}

	var name string
	flags := defaultFlags
	magicType := defaultMagicType
	offset := defaultOffset

	if arch.BinfmtName != "" {
		name = arch.BinfmtName
	} else {
		name = "qemu-" + arch.NameBinfmt
	}

	if arch.BinfmtFlags != "" {
		flags = arch.BinfmtFlags
	}

	if arch.BinfmtMagicType != "" {
		magicType = arch.BinfmtMagicType
	}

	if arch.BinfmtOffset != "" {
		offset = arch.BinfmtOffset
	}

	interpreter := interpreterPath
	if !strings.HasPrefix(interpreterPath, "/run/host/") {
		interpreter = filepath.Join("/run/host", interpreter)
	}

	return &Registration{
		Name:        name,
		MagicType:   magicType,
		Offset:      offset,
		Magic:       arch.ELFMagic,
		Mask:        arch.ELFMask,
		Interpreter: interpreter,
		Flags:       flags,
	}
}

func MountBinfmtMisc() error {
	args := []string{
		"binfmt_misc",
		"-t",
		"binfmt_misc",
		binfmtMiscPath,
	}

	var stdout bytes.Buffer

	if err := shell.Run("mount", nil, &stdout, nil, args...); err != nil {
		return fmt.Errorf("failed to mount binfmt_misc: %w", err)
	}

	logrus.Debugf("Result of mount command: %s", stdout.String())

	return nil
}

func RegisterBinfmtMisc(archID int, interpreterPath string) error {
	reg := getDefaultRegistration(archID, interpreterPath)
	if reg == nil {
		logrus.Debugf("Unable to register binfmt_misc for architecture '%s'", GetArchNameOCI(archID))
		return fmt.Errorf("Toolbx does not support architecture '%s'", GetArchNameOCI(archID))
	}

	if err := reg.register(); err != nil {
		return err
	}

	return nil
}

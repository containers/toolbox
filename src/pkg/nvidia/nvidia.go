/*
 * Copyright © 2024 Red Hat Inc.
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

package nvidia

import (
	"errors"
	"io"

	"github.com/NVIDIA/go-nvlib/pkg/nvlib/info"
	"github.com/NVIDIA/nvidia-container-toolkit/pkg/nvcdi"
	nvspec "github.com/NVIDIA/nvidia-container-toolkit/pkg/nvcdi/spec"
	"github.com/sirupsen/logrus"
	"tags.cncf.io/container-device-interface/specs-go"
)

var (
	logLevel = logrus.ErrorLevel
)

var (
	ErrPlatformUnsupported = errors.New("platform is unsupported")
)

func createNullLogger() *logrus.Logger {
	null := logrus.New()
	null.SetLevel(logrus.PanicLevel)
	null.SetOutput(io.Discard)
	return null
}

func GenerateCDISpec() (*specs.Spec, error) {
	logrus.Debugf("Generating Container Device Interface for NVIDIA")

	var logger *logrus.Logger
	if logLevel < logrus.DebugLevel {
		logger = createNullLogger()
	} else {
		logger = logrus.StandardLogger()
	}

	info := info.New(info.WithLogger(logger))

	if ok, reason := info.HasDXCore(); ok {
		logrus.Debugf("Generating Container Device Interface for NVIDIA: Windows is unsupported: %s", reason)
		return nil, ErrPlatformUnsupported
	}

	hasNvml, reason := info.HasNvml()
	if !hasNvml {
		logrus.Debugf("Generating Container Device Interface for NVIDIA: NVML not found: %s", reason)
	}

	isTegra, reason := info.IsTegraSystem()
	if !isTegra {
		logrus.Debugf("Generating Container Device Interface for NVIDIA: not a Tegra system: %s", reason)
	}

	if !hasNvml && !isTegra {
		logrus.Debug("Generating Container Device Interface for NVIDIA: skipping")
		return nil, ErrPlatformUnsupported
	}

	cdi, err := nvcdi.New(nvcdi.WithInfoLib(info), nvcdi.WithLogger(logger))
	if err != nil {
		logrus.Debugf("Generating Container Device Interface for NVIDIA: failed to create library: %s", err)
		return nil, errors.New("failed to create Container Device Interface library for NVIDIA")
	}

	commonEdits, err := cdi.GetCommonEdits()
	if err != nil {
		logrus.Debugf("Generating Container Device Interface for NVIDIA: failed to get containerEdits: %s", err)
		return nil, errors.New("failed to get Container Device Interface containerEdits for NVIDIA")
	}

	spec, err := nvspec.New(nvspec.WithEdits(*commonEdits.ContainerEdits))
	if err != nil {
		logrus.Debugf("Generating Container Device Interface for NVIDIA: failed to generate: %s", err)
		return nil, errors.New("failed to generate Container Device Interface for NVIDIA")
	}

	specRaw := spec.Raw()
	logrus.Debugf("Generated Container Device Interface for NVIDIA with version %s", specRaw.Version)

	return specRaw, nil
}

func SetLogLevel(level logrus.Level) {
	logLevel = level
}

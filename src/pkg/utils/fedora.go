/*
 * Copyright © 2021 – 2025 Red Hat Inc.
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
	"strings"

	"github.com/sirupsen/logrus"
)

func getDefaultReleaseFedora() (string, error) {
	release, err := getHostVersionID()
	if err != nil {
		return "", err
	}

	return release, nil
}

func getFullyQualifiedImageFedora(image, release string) string {
	imageFull := "registry.fedoraproject.org/" + image
	return imageFull
}

func getP11KitClientPathsFedora() []string {
	paths := []string{"/usr/lib64/pkcs11/p11-kit-client.so"}
	return paths
}

func parseReleaseFedora(release string) (string, error) {
	if strings.HasPrefix(release, "F") || strings.HasPrefix(release, "f") {
		release = release[1:]
	}

	releaseN, err := strconv.Atoi(release)
	if err != nil {
		logrus.Debugf("Parsing release %s as an integer failed: %s", release, err)
		return "", &ParseReleaseError{"The release must be a positive integer."}
	}

	if releaseN <= 0 {
		return "", &ParseReleaseError{"The release must be a positive integer."}
	}

	return release, nil
}

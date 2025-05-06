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
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func getDefaultReleaseRHEL() (string, error) {
	release, err := getHostVersionID()
	if err != nil {
		return "", err
	}

	return release, nil
}

func getFullyQualifiedImageRHEL(image, release string) string {
	i := strings.IndexRune(release, '.')
	if i == -1 {
		panicMsg := fmt.Sprintf("release %s not in '<major>.<minor>' format", release)
		panic(panicMsg)
	}

	releaseMajor := release[:i]
	imageFull := "registry.access.redhat.com/ubi" + releaseMajor + "/" + image
	return imageFull
}

func getP11KitClientPathsRHEL() []string {
	paths := []string{"/usr/lib64/pkcs11/p11-kit-client.so"}
	return paths
}

func parseReleaseRHEL(release string) (string, error) {
	if i := strings.IndexRune(release, '.'); i == -1 {
		return "", &ParseReleaseError{"The release must be in the '<major>.<minor>' format."}
	}

	releaseN, err := strconv.ParseFloat(release, 32)
	if err != nil {
		logrus.Debugf("Parsing release %s as a float failed: %s", release, err)
		return "", &ParseReleaseError{"The release must be in the '<major>.<minor>' format."}
	}

	if releaseN <= 0 {
		return "", &ParseReleaseError{"The release must be a positive number."}
	}

	return release, nil
}

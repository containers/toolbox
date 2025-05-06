/*
 * Copyright © 2023 – 2025 Red Hat Inc.
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
	"unicode/utf8"

	"github.com/sirupsen/logrus"
)

func getDefaultReleaseUbuntu() (string, error) {
	release, err := getHostVersionID()
	if err != nil {
		return "", err
	}

	return release, nil
}

func getFullyQualifiedImageUbuntu(image, release string) string {
	imageFull := "quay.io/toolbx/" + image
	return imageFull
}

func parseReleaseUbuntu(release string) (string, error) {
	releaseParts := strings.Split(release, ".")
	if len(releaseParts) != 2 {
		return "", &ParseReleaseError{"The release must be in the 'YY.MM' format."}
	}

	releaseYear, err := strconv.Atoi(releaseParts[0])
	if err != nil {
		logrus.Debugf("Parsing release year %s as an integer failed: %s", releaseParts[0], err)
		return "", &ParseReleaseError{"The release must be in the 'YY.MM' format."}
	}

	if releaseYear < 4 {
		return "", &ParseReleaseError{"The release year must be 4 or more."}
	}

	releaseYearLen := utf8.RuneCountInString(releaseParts[0])
	if releaseYearLen > 2 {
		return "", &ParseReleaseError{"The release year cannot have more than two digits."}
	} else if releaseYear < 10 && releaseYearLen == 2 {
		return "", &ParseReleaseError{"The release year cannot have a leading zero."}
	}

	releaseMonth, err := strconv.Atoi(releaseParts[1])
	if err != nil {
		logrus.Debugf("Parsing release month %s as an integer failed: %s", releaseParts[1], err)
		return "", &ParseReleaseError{"The release must be in the 'YY.MM' format."}
	}

	if releaseMonth < 1 {
		return "", &ParseReleaseError{"The release month must be between 01 and 12."}
	} else if releaseMonth > 12 {
		return "", &ParseReleaseError{"The release month must be between 01 and 12."}
	}

	releaseMonthLen := utf8.RuneCountInString(releaseParts[1])
	if releaseMonthLen != 2 {
		return "", &ParseReleaseError{"The release month must have two digits."}
	}

	return release, nil
}

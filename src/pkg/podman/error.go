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

package podman

import (
	"bufio"
	"bytes"
	"strings"
)

// internalError serves for representing errors printed by Podman to stderr
type internalError struct {
	errors []string
}

func (e *internalError) Error() string {
	if e.errors == nil || len(e.errors) == 0 {
		return ""
	}

	var builder strings.Builder

	for i, part := range e.errors {
		if i != 0 {
			builder.WriteString(": ")
		}

		builder.WriteString(part)
	}

	return builder.String()
}

// Is lexically compares errors
//
// The comparison is done for every part in the error chain not across.
func (e *internalError) Is(target error) bool {
	if target == nil {
		return false
	}

	if e.errors == nil || len(e.errors) == 0 {
		return false
	}

	for _, part := range e.errors {
		if part == target.Error() {
			return true
		}
	}

	return false
}

func (e *internalError) Unwrap() error {
	if e.errors == nil || len(e.errors) <= 1 {
		return nil
	}

	return &internalError{e.errors[1:]}
}

// parseErrorMsg serves for converting error output of Podman into an error
// that can be further used in Go
func parseErrorMsg(stderr *bytes.Buffer) error {
	// Stderr is not used only for error messages but also for things like
	// progress bars. We're only interested in the error messages.

	var errMsgFound bool
	var errMsgParts []string

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "Error: ") {
			line = strings.TrimPrefix(line, "Error: ")
			errMsgFound = true
		}

		if errMsgFound {
			line = strings.TrimSpace(line)
			line = strings.Trim(line, ":")

			parts := strings.Split(line, ": ")
			errMsgParts = append(errMsgParts, parts...)
		}
	}

	if !errMsgFound {
		return nil
	}

	return &internalError{errMsgParts}
}

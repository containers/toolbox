/*
 * Copyright © 2019 – 2021 Red Hat Inc.
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

	for {
		if e.errors[0] == target.Error() {
			return true
		}

		err := e.Unwrap()
		if err == nil {
			return false
		}
		e = err.(*internalError)
	}
}

func (e internalError) Unwrap() error {
	if e.errors == nil || len(e.errors) <= 1 {
		return nil
	}

	return &internalError{e.errors[1:]}
}

// parseErrorMsg serves for converting error output of Podman into an error
// that can be further used in Go
func parseErrorMsg(stderr *bytes.Buffer) error {
	errMsg := stderr.String()
	errMsg = strings.TrimSpace(errMsg)
	if errMsg == "" {
		return nil
	}

	// Stderr is not used only for error messages but also for e.g. loading
	// bars. Also, several errors can be present (e.g., when 'podman pull' retries
	// to pull an image several times after erroring). Get all of them.
	errMsgSplit := strings.Split(errMsg, "Error: ")
	// We're only interested in the part with the last error
	errMsg = errMsgSplit[len(errMsgSplit)-1]
	// Sometimes an error contains a newline (e.g., responses from Docker
	// registry). Normalize them into further parseable error message.
	// See an example bellow.
	errMsg = strings.ReplaceAll(errMsg, "\n", ": ")
	// Wrapped error messages are usually separated by a colon followed by
	// a single space character
	errMsgPartsRaw := strings.Split(errMsg, ": ")
	// The parts of the err message still can have a whitespace or a colon at the
	// beginning or at the end. Trim them.
	var errMsgParts []string
	for _, part := range errMsgPartsRaw {
		part = strings.TrimSpace(part)
		part = strings.Trim(part, ":")
		// Podman error messages are usually prepended with the "Error:" string.
		// Sometimes the error contains errors in a bullet list. This list is
		// usually prepended with a message equal to "errors:".
		//
		// The colons at the end don't have to be checked since they've been
		// trimmed away.
		//
		// Example:
		//   Error: Error initializing source docker://foobar:latest: Error reading manifest latest in docker.io/library/foobar: errors:
		//   denied: requested access to the resource is denied
		//   unauthorized: authentication required
		if part == "Error" || part == "errors" {
			continue
		}

		errMsgParts = append(errMsgParts, part)
	}

	return &internalError{errMsgParts}
}

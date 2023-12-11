/*
 * Copyright © 2023 Red Hat Inc.
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

package term

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTerminalTemporaryFile(t *testing.T) {
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "TestIsTerminalTempFile")
	assert.NoError(t, err)
	fileName := file.Name()
	defer os.Remove(fileName)
	defer file.Close()

	ok := IsTerminal(file)
	assert.False(t, ok)
}

func TestIsTerminalTerminal(t *testing.T) {
	file, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	assert.NoError(t, err)
	defer file.Close()

	ok := IsTerminal(file)
	assert.True(t, ok)
}

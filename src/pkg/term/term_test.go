/*
 * Copyright © 2023 – 2024 Red Hat Inc.
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
	"golang.org/x/sys/unix"
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

func TestNewStateFromDifferent(t *testing.T) {
	file, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	assert.NoError(t, err)
	defer file.Close()

	oldState, err := GetState(file)
	assert.NoError(t, err)
	assert.Equal(t, uint32(unix.ECHO), oldState.Lflag&unix.ECHO)
	assert.Equal(t, uint32(unix.ICANON), oldState.Lflag&unix.ICANON)
	assert.NotEqual(t, uint8(13), oldState.Cc[unix.VMIN])
	assert.NotEqual(t, uint8(42), oldState.Cc[unix.VTIME])

	newState := NewStateFrom(oldState, WithVMIN(13), WithVTIME(42), WithoutECHO(), WithoutICANON())
	assert.Empty(t, newState.Lflag&unix.ECHO)
	assert.Empty(t, newState.Lflag&unix.ICANON)
	assert.Equal(t, uint8(13), newState.Cc[unix.VMIN])
	assert.Equal(t, uint8(42), newState.Cc[unix.VTIME])
}

func TestNewStateFromNOP(t *testing.T) {
	file, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	assert.NoError(t, err)
	defer file.Close()

	oldState, err := GetState(file)
	assert.NoError(t, err)

	newState := NewStateFrom(oldState)
	assert.Equal(t, oldState.Cc, newState.Cc)
	assert.Equal(t, oldState.Cflag, newState.Cflag)
	assert.Equal(t, oldState.Iflag, newState.Iflag)
	assert.Equal(t, oldState.Ispeed, newState.Ispeed)
	assert.Equal(t, oldState.Lflag, newState.Lflag)
	assert.Equal(t, oldState.Line, newState.Line)
	assert.Equal(t, oldState.Oflag, newState.Oflag)
	assert.Equal(t, oldState.Ospeed, newState.Ospeed)
}

func TestSetStateDifferent(t *testing.T) {
	file, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	assert.NoError(t, err)
	defer file.Close()

	oldState, err := GetState(file)
	assert.NoError(t, err)
	assert.Equal(t, uint32(unix.ECHO), oldState.Lflag&unix.ECHO)
	assert.Equal(t, uint32(unix.ICANON), oldState.Lflag&unix.ICANON)
	assert.NotEqual(t, uint8(13), oldState.Cc[unix.VMIN])
	assert.NotEqual(t, uint8(42), oldState.Cc[unix.VTIME])

	newState := NewStateFrom(oldState, WithVMIN(13), WithVTIME(42), WithoutECHO(), WithoutICANON())
	assert.Empty(t, newState.Lflag&unix.ECHO)
	assert.Empty(t, newState.Lflag&unix.ICANON)
	assert.Equal(t, uint8(13), newState.Cc[unix.VMIN])
	assert.Equal(t, uint8(42), newState.Cc[unix.VTIME])

	err = SetState(file, newState)
	assert.NoError(t, err)

	newState2, err := GetState(file)
	assert.NoError(t, err)
	assert.Equal(t, newState.Cc, newState2.Cc)
	assert.Equal(t, newState.Cflag, newState2.Cflag)
	assert.Equal(t, newState.Iflag, newState2.Iflag)
	assert.Equal(t, newState.Ispeed, newState2.Ispeed)
	assert.Equal(t, newState.Lflag, newState2.Lflag)
	assert.Equal(t, newState.Line, newState2.Line)
	assert.Equal(t, newState.Oflag, newState2.Oflag)
	assert.Equal(t, newState.Ospeed, newState2.Ospeed)
}

func TestSetStateNOP(t *testing.T) {
	file, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	assert.NoError(t, err)
	defer file.Close()

	oldState, err := GetState(file)
	assert.NoError(t, err)

	err = SetState(file, oldState)
	assert.NoError(t, err)

	newState, err := GetState(file)
	assert.NoError(t, err)
	assert.Equal(t, oldState.Cc, newState.Cc)
	assert.Equal(t, oldState.Cflag, newState.Cflag)
	assert.Equal(t, oldState.Iflag, newState.Iflag)
	assert.Equal(t, oldState.Ispeed, newState.Ispeed)
	assert.Equal(t, oldState.Lflag, newState.Lflag)
	assert.Equal(t, oldState.Line, newState.Line)
	assert.Equal(t, oldState.Oflag, newState.Oflag)
	assert.Equal(t, oldState.Ospeed, newState.Ospeed)
}

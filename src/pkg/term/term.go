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

	"golang.org/x/sys/unix"
)

type Option func(*unix.Termios)

func GetState(file *os.File) (*unix.Termios, error) {
	fileFD := file.Fd()
	fileFDInt := int(fileFD)
	state, err := unix.IoctlGetTermios(fileFDInt, unix.TCGETS)
	return state, err
}

func IsTerminal(file *os.File) bool {
	if _, err := GetState(file); err != nil {
		return false
	}

	return true
}

func NewStateFrom(oldState *unix.Termios, options ...Option) *unix.Termios {
	newState := *oldState
	for _, option := range options {
		option(&newState)
	}

	return &newState
}

func SetState(file *os.File, state *unix.Termios) error {
	fileFD := file.Fd()
	fileFDInt := int(fileFD)
	err := unix.IoctlSetTermios(fileFDInt, unix.TCSETS, state)
	return err
}

func WithVMIN(vmin uint8) Option {
	return func(state *unix.Termios) {
		state.Cc[unix.VMIN] = vmin
	}
}

func WithVTIME(vtime uint8) Option {
	return func(state *unix.Termios) {
		state.Cc[unix.VTIME] = vtime
	}
}

func WithoutECHO() Option {
	return func(state *unix.Termios) {
		state.Lflag &^= unix.ECHO
	}
}

func WithoutICANON() Option {
	return func(state *unix.Termios) {
		state.Lflag &^= unix.ICANON
	}
}

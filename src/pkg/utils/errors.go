/*
 * Copyright © 2022 Red Hat Inc.
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
)

type ConfigurationError struct {
	File string
	Err  error
}

type ContainerError struct {
	Container string
	Image     string
	Err       error
}

type DistroError struct {
	Distro string
	Err    error
}

type ImageError struct {
	Image string
	Err   error
}

type ParseReleaseError struct {
	Hint string
}

func (err *ConfigurationError) Error() string {
	errMsg := fmt.Sprintf("%s: %s", err.File, err.Err)
	return errMsg
}

func (err *ConfigurationError) Unwrap() error {
	return err.Err
}

func (err *ContainerError) Error() string {
	errMsg := fmt.Sprintf("%s: %s", err.Container, err.Err)
	return errMsg
}

func (err *ContainerError) Unwrap() error {
	return err.Err
}

func (err *DistroError) Error() string {
	errMsg := fmt.Sprintf("%s: %s", err.Distro, err.Err)
	return errMsg
}

func (err *DistroError) Unwrap() error {
	return err.Err
}

func (err *ImageError) Error() string {
	errMsg := fmt.Sprintf("%s: %s", err.Image, err.Err)
	return errMsg
}

func (err *ImageError) Unwrap() error {
	return err.Err
}

func (err *ParseReleaseError) Error() string {
	return err.Hint
}

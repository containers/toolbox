/*
 * Copyright © 2019 – 2024 Red Hat Inc.
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
	"encoding/json"
	"fmt"
	"io"

	"github.com/HarryMichal/go-version"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Container struct {
	ID      string
	Names   []string
	Status  string
	Created string
	Image   string
	Labels  map[string]string
}

type Image struct {
	ID      string
	Names   []string
	Created string
	Labels  map[string]string
}

type ImageSlice []Image

var (
	podmanVersion string
)

var (
	LogLevel = logrus.ErrorLevel
)

func (container *Container) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID      string
		Names   interface{}
		Status  string
		State   interface{}
		Created interface{}
		Image   string
		Labels  map[string]string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	container.ID = raw.ID

	// In Podman V1 the field 'Names' held a single string but since Podman V2 the
	// field holds an array of strings
	switch value := raw.Names.(type) {
	case string:
		container.Names = append(container.Names, value)
	case []interface{}:
		for _, v := range value {
			container.Names = append(container.Names, v.(string))
		}
	}

	// In Podman V1 the field holding a string about the container's state was
	// called 'Status' and field 'State' held a number representing the state. In
	// Podman V2 the string was moved to 'State' and field 'Status' was dropped.
	switch value := raw.State.(type) {
	case string:
		container.Status = value
	case float64:
		container.Status = raw.Status
	}

	// In Podman V1 the field 'Created' held a human-readable string in format
	// "5 minutes ago". Since Podman V2 the field holds an integer with Unix time.
	// After a discussion in https://github.com/containers/podman/issues/6594 the
	// previous value was moved to field 'CreatedAt'. Since we're already using
	// the 'github.com/docker/go-units' library, we'll stop using the provided
	// human-readable string and assemble it ourselves. Go interprets numbers in
	// JSON as float64.
	switch value := raw.Created.(type) {
	case string:
		container.Created = value
	case float64:
		container.Created = utils.HumanDuration(int64(value))
	}

	container.Image = raw.Image
	container.Labels = raw.Labels
	return nil
}

func (image *Image) FlattenNames(fillNameWithID bool) []Image {
	var ret []Image

	if len(image.Names) == 0 {
		flattenedImage := *image

		if fillNameWithID {
			shortID := utils.ShortID(image.ID)
			flattenedImage.Names = []string{shortID}
		} else {
			flattenedImage.Names = []string{"<none>"}
		}

		ret = []Image{flattenedImage}
		return ret
	}

	ret = make([]Image, 0, len(image.Names))

	for _, name := range image.Names {
		flattenedImage := *image
		flattenedImage.Names = []string{name}
		ret = append(ret, flattenedImage)
	}

	return ret
}

func (image *Image) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID      string
		Names   []string
		Created interface{}
		Labels  map[string]string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	image.ID = raw.ID
	image.Names = raw.Names

	// Until Podman 2.0.x the field 'Created' held a human-readable string in
	// format "5 minutes ago". Since Podman 2.1 the field holds an integer with
	// Unix time. Go interprets numbers in JSON as float64.
	switch value := raw.Created.(type) {
	case string:
		image.Created = value
	case float64:
		image.Created = utils.HumanDuration(int64(value))
	}

	image.Labels = raw.Labels
	return nil
}

func (images ImageSlice) Len() int {
	return len(images)
}

func (images ImageSlice) Less(i, j int) bool {
	if len(images[i].Names) != 1 {
		panic("cannot sort unflattened ImageSlice")
	}

	if len(images[j].Names) != 1 {
		panic("cannot sort unflattened ImageSlice")
	}

	return images[i].Names[0] < images[j].Names[0]
}

func (images ImageSlice) Swap(i, j int) {
	images[i], images[j] = images[j], images[i]
}

// CheckVersion compares provided version with the version of Podman.
//
// Takes in one string parameter that should be in the format that is used for versioning (eg. 1.0.0, 2.5.1-dev).
//
// Returns true if the current version is equal to or higher than the required version.
func CheckVersion(requiredVersion string) bool {
	currentVersion, _ := GetVersion()

	currentVersion = version.Normalize(currentVersion)
	requiredVersion = version.Normalize(requiredVersion)

	return version.CompareSimple(currentVersion, requiredVersion) >= 0
}

// ContainerExists checks using Podman if a container with given ID/name exists.
//
// Parameter container is a name or an id of a container.
func ContainerExists(container string) (bool, error) {
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "container", "exists", container}

	exitCode, err := shell.RunWithExitCode("podman", nil, nil, nil, args...)
	if exitCode != 0 && err == nil {
		err = fmt.Errorf("failed to find container %s", container)
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetContainers is a wrapper function around `podman ps --format json` command.
//
// Parameter args accepts an array of strings to be passed to the wrapped command (eg. ["-a", "--filter", "123"]).
//
// Returned value is a slice of Containers.
//
// If a problem happens during execution, first argument is nil and second argument holds the error message.
func GetContainers(args ...string) ([]Container, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args = append([]string{"--log-level", logLevelString, "ps", "--format", "json"}, args...)

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	data := stdout.Bytes()
	var containers []Container
	if err := json.Unmarshal(data, &containers); err != nil {
		return nil, err
	}

	return containers, nil
}

// GetImages is a wrapper function around `podman images --format json` command.
//
// Parameter args accepts an array of strings to be passed to the wrapped command (eg. ["-a", "--filter", "123"]).
//
// Returned value is a slice of Images.
//
// If a problem happens during execution, first argument is nil and second argument holds the error message.
func GetImages(args ...string) ([]Image, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args = append([]string{"--log-level", logLevelString, "images", "--format", "json"}, args...)
	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	data := stdout.Bytes()
	var images []Image
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}

	return images, nil
}

// GetVersion returns version of Podman in a string
func GetVersion() (string, error) {
	if podmanVersion != "" {
		return podmanVersion, nil
	}

	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "version", "--format", "json"}

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return "", err
	}

	output := stdout.Bytes()
	var jsonoutput map[string]interface{}
	if err := json.Unmarshal(output, &jsonoutput); err != nil {
		return "", err
	}

	podmanClientInfoInterface := jsonoutput["Client"]
	switch podmanClientInfo := podmanClientInfoInterface.(type) {
	case nil:
		podmanVersion = jsonoutput["Version"].(string)
	case map[string]interface{}:
		podmanVersion = podmanClientInfo["Version"].(string)
	}
	return podmanVersion, nil
}

// ImageExists checks using Podman if an image with given ID/name exists.
//
// Parameter image is a name or an id of an image.
func ImageExists(image string) (bool, error) {
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "image", "exists", image}

	exitCode, err := shell.RunWithExitCode("podman", nil, nil, nil, args...)
	if exitCode != 0 && err == nil {
		err = fmt.Errorf("failed to find image %s", image)
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Inspect is a wrapper around 'podman inspect' command
//
// Parameter 'typearg' takes in values 'container' or 'image' that is passed to the --type flag
func Inspect(typearg string, target string) (map[string]interface{}, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "inspect", "--format", "json", "--type", typearg, target}

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var info []map[string]interface{}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return info[0], nil
}

func IsToolboxContainer(container string) (bool, error) {
	info, err := Inspect("container", container)
	if err != nil {
		return false, fmt.Errorf("failed to inspect container %s", container)
	}

	labels, _ := info["Config"].(map[string]interface{})["Labels"].(map[string]interface{})
	if labels["com.github.containers.toolbox"] != "true" && labels["com.github.debarshiray.toolbox"] != "true" {
		return false, fmt.Errorf("%s is not a Toolbx container", container)
	}

	return true, nil
}

func IsToolboxImage(image string) (bool, error) {
	info, err := Inspect("image", image)
	if err != nil {
		return false, fmt.Errorf("failed to inspect image %s", image)
	}

	if info["Labels"] == nil {
		return false, fmt.Errorf("%s is not a Toolbx image", image)
	}

	labels := info["Labels"].(map[string]interface{})
	if labels["com.github.containers.toolbox"] != "true" && labels["com.github.debarshiray.toolbox"] != "true" {
		return false, fmt.Errorf("%s is not a Toolbx image", image)
	}

	return true, nil
}

// Pull pulls an image
//
// authfile is a path to a JSON authentication file and is internally used only
// if it is not an empty string.
func Pull(imageName string, authfile string) error {
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "pull"}

	if authfile != "" {
		args = append(args, []string{"--authfile", authfile}...)
	}

	args = append(args, imageName)

	if err := shell.Run("podman", nil, nil, nil, args...); err != nil {
		return err
	}

	return nil
}

func RemoveContainer(container string, forceDelete bool) error {
	logrus.Debugf("Removing container %s", container)

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "rm"}

	if forceDelete {
		args = append(args, "--force")
	}

	args = append(args, container)

	exitCode, err := shell.RunWithExitCode("podman", nil, nil, nil, args...)
	switch exitCode {
	case 0:
		if err != nil {
			panic("unexpected error: 'podman rm' finished successfully")
		}
	case 1:
		err = fmt.Errorf("container %s does not exist", container)
	case 2:
		err = fmt.Errorf("container %s is running", container)
	default:
		err = fmt.Errorf("failed to remove container %s", container)
	}

	if err != nil {
		return err
	}

	return nil
}

func RemoveImage(image string, forceDelete bool) error {
	logrus.Debugf("Removing image %s", image)

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "rmi"}

	if forceDelete {
		args = append(args, "--force")
	}

	args = append(args, image)

	exitCode, err := shell.RunWithExitCode("podman", nil, nil, nil, args...)
	switch exitCode {
	case 0:
		if err != nil {
			panic("unexpected error: 'podman rmi' finished successfully")
		}
	case 1:
		err = fmt.Errorf("image %s does not exist", image)
	case 2:
		err = fmt.Errorf("image %s has dependent children", image)
	default:
		err = fmt.Errorf("failed to remove image %s", image)
	}

	if err != nil {
		return err
	}

	return nil
}

func SetLogLevel(logLevel logrus.Level) {
	LogLevel = logLevel
}

func Start(container string, stderr io.Writer) error {
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "start", container}

	if err := shell.Run("podman", nil, nil, stderr, args...); err != nil {
		return err
	}

	return nil
}

func SystemMigrate(ociRuntimeRequired string) error {
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "system", "migrate"}
	if ociRuntimeRequired != "" {
		args = append(args, []string{"--new-runtime", ociRuntimeRequired}...)
	}

	if err := shell.Run("podman", nil, nil, nil, args...); err != nil {
		return err
	}

	return nil
}

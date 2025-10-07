/*
 * Copyright © 2019 – 2025 Red Hat Inc.
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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"time"

	"github.com/HarryMichal/go-version"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
)

var (
	podmanVersion string
)

var (
	ErrImageRepoTagsEmpty = errors.New("image has empty RepoTags")

	ErrImageRepoTagsMissing = errors.New("image has no RepoTags")

	LogLevel = logrus.ErrorLevel
)

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
func GetContainers(args ...string) (*Containers, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args = append([]string{"--log-level", logLevelString, "ps", "--format", "json"}, args...)

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	data := stdout.Bytes()
	var containers []containerPS
	if err := json.Unmarshal(data, &containers); err != nil {
		return nil, err
	}

	return &Containers{containers, 0}, nil
}

// GetImages is a wrapper function around `podman images --format json` command.
//
// Parameter fillNameWithID is a boolean that indicates if the image names should be filled with the ID, when there are no names.
// Parameter args accepts an array of strings to be passed to the wrapped command (eg. ["-a", "--filter", "123"]).
//
// Returned value is a slice of Images.
//
// If a problem happens during execution, first argument is nil and second argument holds the error message.
func GetImages(fillNameWithID bool, sortByName bool, args ...string) (*Images, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args = append([]string{"--log-level", logLevelString, "images", "--format", "json"}, args...)

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	data := stdout.Bytes()
	var images []imageImages
	if err := json.Unmarshal(data, &images); err != nil {
		return nil, err
	}

	// Images flattening
	processed := make(map[string]struct{})
	var retImages []imageImages

	for _, image := range images {
		if _, ok := processed[image.ID()]; ok {
			continue
		}

		processed[image.ID()] = struct{}{}

		flattenedImages := image.flattenNames(fillNameWithID)
		retImages = append(retImages, flattenedImages...)
	}

	ret := Images{retImages, 0}
	if sortByName {
		sort.Sort(ret)
	}

	return &ret, nil
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

func GetFullyQualifiedImageFromRepoTags(image string) (string, error) {
	logrus.Debugf("Resolving fully qualified name for image %s from RepoTags", image)

	var imageFull string

	if utils.ImageReferenceHasDomain(image) {
		imageFull = image
	} else {
		imageObj, err := InspectImage(image)
		if err != nil {
			return "", fmt.Errorf("failed to inspect image %s", image)
		}

		if imageObj.RepoTags() == nil {
			return "", &ImageError{image, ErrImageRepoTagsMissing}
		}

		repoTags := imageObj.RepoTags()
		if len(repoTags) == 0 {
			return "", &ImageError{image, ErrImageRepoTagsEmpty}
		}

		for _, repoTag := range repoTags {
			tag := utils.ImageReferenceGetTag(repoTag)
			if tag != "latest" {
				imageFull = repoTag
				break
			}
		}

		if imageFull == "" {
			imageFull = repoTags[0]
		}
	}

	logrus.Debugf("Resolved image %s to %s", image, imageFull)

	return imageFull, nil
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

// InspectContainer is a wrapper around 'podman inspect --type container' command
func InspectContainer(container string) (Container, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "inspect", "--format", "json", "--type", "container", container}

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var containers []containerInspect
	if err := json.Unmarshal(output, &containers); err != nil {
		return nil, err
	}

	return &containers[0], nil
}

// InspectImage is a wrapper around 'podman inspect --type image' command
func InspectImage(image string) (Image, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "inspect", "--format", "json", "--type", "image", image}

	if err := shell.Run("podman", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var images []imageInspect
	if err := json.Unmarshal(output, &images); err != nil {
		return nil, err
	}

	return &images[0], nil
}

func Logs(container string, since time.Time, stderr io.Writer) error {
	ctx := context.Background()
	err := LogsContext(ctx, container, false, since, stderr)
	return err
}

func LogsContext(ctx context.Context, container string, follow bool, since time.Time, stderr io.Writer) error {
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "logs"}

	if follow {
		args = append(args, "--follow")
	}

	if sinceUnix := since.Unix(); sinceUnix >= 0 {
		sinceUnixString := strconv.FormatInt(sinceUnix, 10)
		args = append(args, []string{"--since", sinceUnixString}...)
	}

	args = append(args, container)

	if err := shell.RunContext(ctx, "podman", nil, nil, stderr, args...); err != nil {
		return err
	}

	return nil
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
		err = fmt.Errorf("container %s not found", container)
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
		err = fmt.Errorf("image %s not found", image)
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

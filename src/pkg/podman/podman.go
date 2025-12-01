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
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/HarryMichal/go-version"
	"github.com/containers/toolbox/pkg/shell"
	"github.com/containers/toolbox/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Image struct {
	Created string
	ID      string
	Labels  map[string]string
	Names   []string
}

type BuildOptions struct {
	Context string
	Tag     string
}

type ImageSlice []Image

var (
	podmanVersion string
)

var (
	ErrImageRepoTagsEmpty = errors.New("image has empty RepoTags")

	ErrImageRepoTagsMissing = errors.New("image has no RepoTags")

	LogLevel = logrus.ErrorLevel
)

var (
	ErrBuildContextDoesNotExist = errors.New("build context does not exist")

	ErrBuildContextInvalid = errors.New("build context is not a directory with a Containerfile")
)

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
		Created interface{}
		ID      string
		Labels  map[string]string
		Names   []string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Until Podman 2.0.x the field 'Created' held a human-readable string in
	// format "5 minutes ago". Since Podman 2.1 the field holds an integer with
	// Unix time. Go interprets numbers in JSON as float64.
	switch value := raw.Created.(type) {
	case string:
		image.Created = value
	case float64:
		image.Created = utils.HumanDuration(int64(value))
	}

	image.ID = raw.ID
	image.Labels = raw.Labels
	image.Names = raw.Names
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

func BuildImage(build BuildOptions) (string, error) {
	if !utils.PathExists(build.Context) {
		return "", &utils.BuildError{BuildContext: build.Context, Err: ErrBuildContextDoesNotExist}
	}
	if stat, err := os.Stat(build.Context); err != nil {
		return "", err
	} else {
		if !stat.Mode().IsDir() {
			return "", &utils.BuildError{BuildContext: build.Context, Err: ErrBuildContextInvalid}
		}
	}
	if !utils.PathExists(build.Context+"/Containerfile") && !utils.PathExists(build.Context+"/Dockerfile") {
		return "", &utils.BuildError{BuildContext: build.Context, Err: ErrBuildContextInvalid}
	}
	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "build", build.Context}
	if build.Tag != "" {
		args = append(args, "--tag", build.Tag)
	}

	stdout := new(bytes.Buffer)
	if err := shell.Run("podman", nil, stdout, nil, args...); err != nil {
		return "", err
	}
	output := strings.TrimRight(stdout.String(), "\n")
	imageIdBegin := strings.LastIndex(output, "\n") + 1
	imageId := output[imageIdBegin:]

	var name string
	if build.Tag == "" {
		info, err := InspectImage(imageId)
		if err != nil {
			return "", err
		}
		name = info["Labels"].(map[string]interface{})["name"].(string)
		args = []string{"--log-level", logLevelString, "tag", imageId, name}
		if err := shell.Run("podman", nil, nil, nil, args...); err != nil {
			return "", err
		}
	} else {
		name = build.Tag
	}

	return name, nil
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

func GetFullyQualifiedImageFromRepoTags(image string) (string, error) {
	logrus.Debugf("Resolving fully qualified name for image %s from RepoTags", image)

	var imageFull string

	if utils.ImageReferenceHasDomain(image) {
		imageFull = image
	} else {
		info, err := InspectImage(image)
		if err != nil {
			return "", fmt.Errorf("failed to inspect image %s", image)
		}

		if info["RepoTags"] == nil {
			return "", &ImageError{image, ErrImageRepoTagsMissing}
		}

		repoTags := info["RepoTags"].([]interface{})
		if len(repoTags) == 0 {
			return "", &ImageError{image, ErrImageRepoTagsEmpty}
		}

		for _, repoTag := range repoTags {
			repoTagString := repoTag.(string)
			tag := utils.ImageReferenceGetTag(repoTagString)
			if tag != "latest" {
				imageFull = repoTagString
				break
			}
		}

		if imageFull == "" {
			imageFull = repoTags[0].(string)
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
func InspectImage(image string) (map[string]interface{}, error) {
	var stdout bytes.Buffer

	logLevelString := LogLevel.String()
	args := []string{"--log-level", logLevelString, "inspect", "--format", "json", "--type", "image", image}

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

func IsToolboxImage(image string) (bool, error) {
	info, err := InspectImage(image)
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

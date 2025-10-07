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
	"encoding/json"

	"github.com/containers/toolbox/pkg/utils"
)

type Image interface {
	Created() string
	ID() string
	IsToolbx() bool
	Labels() map[string]string
	Names() []string
	RepoTags() []string
}

type imageImages struct {
	created  string
	id       string
	labels   map[string]string
	names    []string
	repoTags []string
}

type imageInspect struct {
	created      string
	entrypoint   []string
	envVars      []string
	id           string
	labels       map[string]string
	namesHistory []string
	repoTags     []string
}

type Images struct {
	data []imageImages
	i    int
}

func (image *imageImages) Created() string {
	return image.created
}

func (image *imageImages) flattenNames(fillNameWithID bool) []imageImages {
	var ret []imageImages

	if len(image.Names()) == 0 {
		flattenedImage := *image

		if fillNameWithID {
			shortID := utils.ShortID(image.ID())
			flattenedImage.setNames([]string{shortID})
		} else {
			flattenedImage.setNames([]string{"<none>"})
		}

		ret = []imageImages{flattenedImage}
		return ret
	}

	ret = make([]imageImages, 0, len(image.Names()))

	for _, name := range image.Names() {
		flattenedImage := *image
		flattenedImage.setNames([]string{name})
		ret = append(ret, flattenedImage)
	}

	return ret
}

func (image *imageImages) ID() string {
	return image.id
}

func (image *imageImages) IsToolbx() bool {
	return isToolbx(image.labels)
}

func (image *imageImages) Labels() map[string]string {
	return image.labels
}

func (image *imageImages) Names() []string {
	return image.names
}

func (image *imageImages) setNames(names []string) {
	image.names = names
}

func (image *imageImages) RepoTags() []string {
	return image.repoTags
}

func (image *imageImages) UnmarshalJSON(data []byte) error {
	var raw struct {
		Created  interface{}
		ID       string
		Labels   map[string]string
		Names    []string
		RepoTags []string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Until Podman 2.0.x the field 'Created' held a human-readable string in
	// format "5 minutes ago". Since Podman 2.1 the field holds an integer with
	// Unix time. Go interprets numbers in JSON as float64.
	switch value := raw.Created.(type) {
	case string:
		image.created = value
	case float64:
		image.created = utils.HumanDuration(int64(value))
	}

	image.id = raw.ID
	image.labels = raw.Labels
	image.names = raw.Names
	image.repoTags = raw.RepoTags
	return nil
}

func (image *imageInspect) Created() string {
	return image.created
}

func (image *imageInspect) Entrypoint() []string {
	return image.entrypoint
}

func (image *imageInspect) EnvVars() []string {
	return image.envVars
}

func (image *imageInspect) ID() string {
	return image.id
}

func (image *imageInspect) IsToolbx() bool {
	return isToolbx(image.labels)
}

func (image *imageInspect) Labels() map[string]string {
	return image.labels
}

func (image *imageInspect) Names() []string {
	return image.namesHistory
}

func (image *imageInspect) RepoTags() []string {
	return image.repoTags
}

func (image *imageInspect) UnmarshalJSON(data []byte) error {
	var raw struct {
		Created interface{}
		ID      string
		Config  struct {
			Labels     map[string]string
			Env        []string
			Entrypoint []string
		}
		NamesHistory []string
		RepoTags     []string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch value := raw.Created.(type) {
	case string:
		image.created = value
	case float64:
		image.created = utils.HumanDuration(int64(value))
	}

	image.id = raw.ID
	image.labels = raw.Config.Labels
	image.envVars = raw.Config.Env
	image.namesHistory = raw.NamesHistory
	image.repoTags = raw.RepoTags
	image.entrypoint = raw.Config.Entrypoint
	return nil
}

func (images *Images) Get() Image {
	if images.i < 1 {
		panic("called Containers.Get() without calling Containers.Next()")
	}

	image := images.data[images.i-1]
	return &image
}

func (images *Images) Next() bool {
	available := images.i < len(images.data)
	if available {
		images.i++
	}

	return available
}

func (images *Images) Reset() {
	images.i = 0
}

func (images Images) Len() int {
	return len(images.data)
}

func (images Images) Less(i, j int) bool {
	if len(images.data[i].Names()) != 1 {
		panic("cannot sort unflattened Images")
	}

	if len(images.data[j].Names()) != 1 {
		panic("cannot sort unflattened Images")
	}

	return images.data[i].Names()[0] < images.data[j].Names()[0]
}

func (images Images) Swap(i, j int) {
	images.data[i], images.data[j] = images.data[j], images.data[i]
}

// ALREADY EXISTS IN CONTAINER.GO
//
// func isToolbx(labels map[string]string) bool {
// 	if labels["com.github.containers.toolbox"] == "true" || labels["com.github.debarshiray.toolbox"] == "true" {
// 		return true
// 	}

// 	return false
// }

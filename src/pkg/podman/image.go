/*
 * Copyright © 2025 – 2026 Red Hat Inc.
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
	Name() string
	Names() []string
}

type Images struct {
	data []imageImages
	i    int
}

type imageImages struct {
	created string
	id      string
	labels  map[string]string
	names   []string
}

func (images *Images) Get() Image {
	if images == nil {
		panic("called Images.Get() on a nil Images")
	}

	if images.i < 1 {
		panic("called Images.Get() without calling Images.Next()")
	}

	image := &images.data[images.i-1]
	return image
}

func (images *Images) Len() int {
	if images == nil {
		return 0
	}

	return len(images.data)
}

func (images *Images) Next() bool {
	if images == nil {
		return false
	}

	available := images.i < len(images.data)
	if available {
		images.i++
	}

	return available
}

func (images *Images) Reset() {
	if images == nil {
		return
	}

	images.i = 0
}

func (image *imageImages) Created() string {
	return image.created
}

func (image *imageImages) flattenNames(fillNameWithID bool) []imageImages {
	var ret []imageImages

	names := image.Names()
	namesCount := len(names)
	if namesCount == 0 {
		flattenedImage := *image

		if fillNameWithID {
			id := image.ID()
			shortID := utils.ShortID(id)
			flattenedImage.names = []string{shortID}
		} else {
			flattenedImage.names = []string{"<none>"}
		}

		ret = []imageImages{flattenedImage}
		return ret
	}

	ret = make([]imageImages, 0, namesCount)

	for _, name := range names {
		flattenedImage := *image
		flattenedImage.names = []string{name}
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
	if image.labels == nil {
		return nil
	}

	labelsCount := len(image.labels)
	ret := make(map[string]string, labelsCount)
	for label, value := range image.labels {
		ret[label] = value
	}

	return ret
}

func (image *imageImages) Name() string {
	if len(image.names) != 1 {
		panic("cannot get name from unflattened Image")
	}

	return image.names[0]
}

func (image *imageImages) Names() []string {
	if image.names == nil {
		return nil
	}

	namesCount := len(image.names)
	ret := make([]string, namesCount)
	copy(ret, image.names)
	return ret
}

func (image *imageImages) UnmarshalJSON(data []byte) error {
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
		image.created = value
	case float64:
		image.created = utils.HumanDuration(int64(value))
	}

	image.id = raw.ID
	image.labels = raw.Labels
	image.names = raw.Names
	return nil
}

/*
 * Copyright © 2024 Red Hat Inc.
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

type Container interface {
	Created() string
	ID() string
	Image() string
	Labels() map[string]string
	Name() string
	Names() []string
	Status() string
}

type containerPS struct {
	created string
	id      string
	image   string
	labels  map[string]string
	names   []string
	status  string
}

type Containers struct {
	data []containerPS
	i    int
}

func (container *containerPS) Created() string {
	return container.created
}

func (container *containerPS) ID() string {
	return container.id
}

func (container *containerPS) Image() string {
	return container.image
}

func (container *containerPS) Labels() map[string]string {
	return container.labels
}

func (container *containerPS) Name() string {
	return container.names[0]
}

func (container *containerPS) Names() []string {
	return container.names
}

func (container *containerPS) Status() string {
	return container.status
}

func (container *containerPS) UnmarshalJSON(data []byte) error {
	var raw struct {
		Created interface{}
		ID      string
		Image   string
		Labels  map[string]string
		Names   interface{}
		State   interface{}
		Status  string
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
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
		container.created = value
	case float64:
		container.created = utils.HumanDuration(int64(value))
	}

	container.id = raw.ID
	container.image = raw.Image
	container.labels = raw.Labels

	// In Podman V1 the field 'Names' held a single string but since Podman V2 the
	// field holds an array of strings
	switch value := raw.Names.(type) {
	case string:
		container.names = append(container.names, value)
	case []interface{}:
		for _, v := range value {
			container.names = append(container.names, v.(string))
		}
	}

	// In Podman V1 the field holding a string about the container's state was
	// called 'Status' and field 'State' held a number representing the state. In
	// Podman V2 the string was moved to 'State' and field 'Status' was dropped.
	switch value := raw.State.(type) {
	case string:
		container.status = value
	case float64:
		container.status = raw.Status
	}

	return nil
}

func (containers *Containers) Get() Container {
	if containers.i < 1 {
		panic("called Containers.Get() without calling Containers.Next()")
	}

	container := containers.data[containers.i-1]
	return &container
}

func (containers *Containers) Next() bool {
	available := containers.i < len(containers.data)
	if available {
		containers.i++
	}

	return available
}

func (containers *Containers) Reset() {
	containers.i = 0
}

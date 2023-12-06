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

package skopeo

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/containers/toolbox/pkg/shell"
)

type Layer struct {
	Size json.Number
}
type Image struct {
	LayersData []Layer
}

func Inspect(ctx context.Context, target string) (*Image, error) {
	var stdout bytes.Buffer

	targetWithTransport := "docker://" + target
	args := []string{"inspect", "--format", "json", targetWithTransport}

	if err := shell.RunContext(ctx, "skopeo", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var image Image
	if err := json.Unmarshal(output, &image); err != nil {
		return nil, err
	}

	return &image, nil
}

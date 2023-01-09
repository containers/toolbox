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
	"encoding/json"

	"github.com/containers/toolbox/pkg/shell"
)

func Inspect(target string) (interface{}, error) {

	var stdout bytes.Buffer
	args := []string{"inspect", target}

	if err := shell.Run("skopeo", nil, &stdout, nil, args...); err != nil {
		return nil, err
	}

	output := stdout.Bytes()
	var images interface{}
	if err := json.Unmarshal(output, &images); err != nil {
		return nil, err
	}

	m := images.(map[string]interface{})
	var totalSize float64 = 0
	for k, v := range m {

		if k == "LayersData" {
			vv := v.([]interface{})

			for _, u := range vv {

				uu := u.(map[string]interface{})
				for j, w := range uu {

					if j == "Size" {
						totalSize += w.(float64)
					}
				}
			}
		}
	}
	totalSizeIS := (totalSize / 1000000)
	return int(totalSizeIS), nil
}

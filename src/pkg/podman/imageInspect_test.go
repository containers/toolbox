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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageInspect(t *testing.T) {
	type expect struct {
		id           string
		isToolbx     bool
		labels       map[string]string
		namesHistory []string
		repoTags     []string
	}

	testCases := []struct {
		name    string
		data    string
		expects []expect
	}{
		{
			name: "podman 1.1.2, fedora-toolbox:29",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/sh\"," +
				"        \"-c\"," +
				"        \"/bin/sh\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f29container\"," +
				"        \"FGC=f29\"," +
				"        \"NAME=fedora-toolbox\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"VERSION=29\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"architecture\": \"x86_64\"," +
				"        \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"        \"build-date\": \"2019-11-06T11:24:10.495113\"," +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"distribution-scope\": \"public\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"f29/fedora-toolbox\"," +
				"        \"release\": \"9\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"29\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2019-11-06T11:36:10.586442Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2019-02-18T06:48:39Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2019-11-06T11:36:10.586442Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"4a6adf1f2a96adf5ea0c02b61f9fa574306f77fc522f39c2ce6bb164daead882\"," +
				"    \"Labels\": {" +
				"      \"architecture\": \"x86_64\"," +
				"      \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"      \"build-date\": \"2019-11-06T11:24:10.495113\"," +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.github.debarshiray.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"distribution-scope\": \"public\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"f29/fedora-toolbox\"," +
				"      \"release\": \"9\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"29\"" +
				"    }," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/f29/fedora-toolbox:29\"" +
				"    ]," +
				"    \"Size\": 437590766," +
				"    \"VirtualSize\": 437590766" +
				"  }," +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/sh\"," +
				"        \"-c\"," +
				"        \"/bin/sh\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f29container\"," +
				"        \"FGC=f29\"," +
				"        \"NAME=fedora-toolbox\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"VERSION=29\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"architecture\": \"x86_64\"," +
				"        \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"        \"build-date\": \"2019-11-06T11:24:10.495113\"," +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"distribution-scope\": \"public\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"f29/fedora-toolbox\"," +
				"        \"release\": \"9\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"29\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-10-06T09:10:26.394950134Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2019-02-18T06:48:39Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2019-11-06T11:36:10.586442Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2025-10-06T09:10:26.394950134Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4\"," +
				"    \"Labels\": {" +
				"      \"architecture\": \"x86_64\"," +
				"      \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"      \"build-date\": \"2019-11-06T11:24:10.495113\"," +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.github.debarshiray.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"distribution-scope\": \"public\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"f29/fedora-toolbox\"," +
				"      \"release\": \"9\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"29\"" +
				"    }," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"localhost/fedora-toolbox-user:29\"" +
				"    ]," +
				"    \"Size\": 437967580," +
				"    \"VirtualSize\": 437967580" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "4a6adf1f2a96adf5ea0c02b61f9fa574306f77fc522f39c2ce6bb164daead882",
					isToolbx: true,
					labels: map[string]string{
						"architecture":                   "x86_64",
						"authoritative-source-url":       "registry.fedoraproject.org",
						"build-date":                     "2019-11-06T11:24:10.495113",
						"com.github.containers.toolbox":  "true",
						"com.github.debarshiray.toolbox": "true",
						"com.redhat.component":           "fedora-toolbox",
						"distribution-scope":             "public",
						"license":                        "MIT",
						"name":                           "f29/fedora-toolbox",
						"release":                        "9",
						"vendor":                         "Fedora Project",
						"version":                        "29",
					},
					namesHistory: nil,
					repoTags: []string{
						"registry.fedoraproject.org/f29/fedora-toolbox:29",
					},
				},
				{
					id:       "8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4",
					isToolbx: true,
					labels: map[string]string{
						"architecture":                   "x86_64",
						"authoritative-source-url":       "registry.fedoraproject.org",
						"build-date":                     "2019-11-06T11:24:10.495113",
						"com.github.containers.toolbox":  "true",
						"com.github.debarshiray.toolbox": "true",
						"com.redhat.component":           "fedora-toolbox",
						"distribution-scope":             "public",
						"license":                        "MIT",
						"name":                           "f29/fedora-toolbox",
						"release":                        "9",
						"vendor":                         "Fedora Project",
						"version":                        "29",
					},
					namesHistory: nil,
					repoTags: []string{
						"localhost/fedora-toolbox-user:29",
					},
				},
			},
		},
		{
			name: "podman 1.8.0, fedora-toolbox:30",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/sh\"," +
				"        \"-c\"," +
				"        \"/bin/sh\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f30container\"," +
				"        \"FGC=f30\"," +
				"        \"NAME=fedora-toolbox\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"VERSION=30\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"architecture\": \"x86_64\"," +
				"        \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"        \"build-date\": \"2020-03-03T10:28:18.406675\"," +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"distribution-scope\": \"public\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"f30/fedora-toolbox\"," +
				"        \"release\": \"12\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"30\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2020-03-03T10:32:20.503064Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2020-02-24T07:48:26Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2020-03-03T10:32:20.503064Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8\"," +
				"    \"Labels\": {" +
				"      \"architecture\": \"x86_64\"," +
				"      \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"      \"build-date\": \"2020-03-03T10:28:18.406675\"," +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.github.debarshiray.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"distribution-scope\": \"public\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"f30/fedora-toolbox\"," +
				"      \"release\": \"12\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"30\"" +
				"    }," +
				"    \"NamesHistory\": []," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/f30/fedora-toolbox:30\"" +
				"    ]," +
				"    \"Size\": 404259600," +
				"    \"VirtualSize\": 404259600" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8",
					isToolbx: true,
					labels: map[string]string{
						"architecture":                   "x86_64",
						"authoritative-source-url":       "registry.fedoraproject.org",
						"build-date":                     "2020-03-03T10:28:18.406675",
						"com.github.containers.toolbox":  "true",
						"com.github.debarshiray.toolbox": "true",
						"com.redhat.component":           "fedora-toolbox",
						"distribution-scope":             "public",
						"license":                        "MIT",
						"name":                           "f30/fedora-toolbox",
						"release":                        "12",
						"vendor":                         "Fedora Project",
						"version":                        "30",
					},
					namesHistory: []string{},
					repoTags: []string{
						"registry.fedoraproject.org/f30/fedora-toolbox:30",
					},
				},
			},
		},
		{
			name: "podman 2.2.1, fedora-toolbox:32",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/sh\"," +
				"        \"-c\"," +
				"        \"/bin/sh\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f32container\"," +
				"        \"FGC=f32\"," +
				"        \"NAME=fedora-toolbox\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"VERSION=32\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"architecture\": \"x86_64\"," +
				"        \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"        \"build-date\": \"2021-04-25T12:01:08.917265\"," +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"distribution-scope\": \"public\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"release\": \"14\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"32\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2021-04-25T12:02:39.696984Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2021-01-06T06:48:53Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2021-04-25T12:02:39.696984Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6\"," +
				"    \"Labels\": {" +
				"      \"architecture\": \"x86_64\"," +
				"      \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"      \"build-date\": \"2021-04-25T12:01:08.917265\"," +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.github.debarshiray.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"distribution-scope\": \"public\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"release\": \"14\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"32\"" +
				"    }," +
				"    \"NamesHistory\": []," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:32\"" +
				"    ]," +
				"    \"Size\": 371369047," +
				"    \"VirtualSize\": 371369047" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6",
					isToolbx: true,
					labels: map[string]string{
						"architecture":                   "x86_64",
						"authoritative-source-url":       "registry.fedoraproject.org",
						"build-date":                     "2021-04-25T12:01:08.917265",
						"com.github.containers.toolbox":  "true",
						"com.github.debarshiray.toolbox": "true",
						"com.redhat.component":           "fedora-toolbox",
						"distribution-scope":             "public",
						"license":                        "MIT",
						"name":                           "fedora-toolbox",
						"release":                        "14",
						"vendor":                         "Fedora Project",
						"version":                        "32",
					},
					namesHistory: []string{},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:32",
					},
				},
			},
		},
		{
			name: "podman 3.4.7, fedora-toolbox:35",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f35container\"," +
				"        \"FGC=f35\"," +
				"        \"NAME=fedora-toolbox\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"VERSION=35\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"architecture\": \"x86_64\"," +
				"        \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"        \"build-date\": \"2022-11-09T13:15:18.185125\"," +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"distribution-scope\": \"public\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"release\": \"18\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"35\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2022-11-09T13:16:43.063221Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2022-09-23T06:49:12Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2022-11-09T13:16:43.063221Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227\"," +
				"    \"Labels\": {" +
				"      \"architecture\": \"x86_64\"," +
				"      \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"      \"build-date\": \"2022-11-09T13:15:18.185125\"," +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.github.debarshiray.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"distribution-scope\": \"public\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"release\": \"18\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"35\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"    ]," +
				"    \"Size\": 495703728," +
				"    \"VirtualSize\": 495703728" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227",
					isToolbx: true,
					labels: map[string]string{
						"architecture":                   "x86_64",
						"authoritative-source-url":       "registry.fedoraproject.org",
						"build-date":                     "2022-11-09T13:15:18.185125",
						"com.github.containers.toolbox":  "true",
						"com.github.debarshiray.toolbox": "true",
						"com.redhat.component":           "fedora-toolbox",
						"distribution-scope":             "public",
						"license":                        "MIT",
						"name":                           "fedora-toolbox",
						"release":                        "18",
						"vendor":                         "Fedora Project",
						"version":                        "35",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:35",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:35",
					},
				},
			},
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f38container\"," +
				"        \"FGC=f38\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"architecture\": \"x86_64\"," +
				"        \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"        \"build-date\": \"2024-02-01T19:07:19.783946\"," +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"distribution-scope\": \"public\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"release\": \"20\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-01T19:09:35.77181Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2024-01-25T06:49:26Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2024-02-01T19:09:35.77181Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce\"," +
				"    \"Labels\": {" +
				"      \"architecture\": \"x86_64\"," +
				"      \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"      \"build-date\": \"2024-02-01T19:07:19.783946\"," +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"distribution-scope\": \"public\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"release\": \"20\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"    ]," +
				"    \"Size\": 1748730844," +
				"    \"VirtualSize\": 1748730844" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce",
					isToolbx: true,
					labels: map[string]string{
						"architecture":                  "x86_64",
						"authoritative-source-url":      "registry.fedoraproject.org",
						"build-date":                    "2024-02-01T19:07:19.783946",
						"com.github.containers.toolbox": "true",
						"com.redhat.component":          "fedora-toolbox",
						"distribution-scope":            "public",
						"license":                       "MIT",
						"name":                          "fedora-toolbox",
						"release":                       "20",
						"vendor":                        "Fedora Project",
						"version":                       "38",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:38",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:38",
					},
				},
			},
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38 locally built",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f38container\"," +
				"        \"FGC=f38\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"io.buildah.version\": \"1.33.7\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2026-04-08T21:30:47Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2024-01-25T06:49:26Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2026-04-08T21:25:10Z\"," +
				"        \"created_by\": \"/bin/sh -c #(nop) ARG NAME\"," +
				"        \"empty_layer\": true" +
				"      }," +
				"      {" +
				"        \"created\": \"2026-04-08T21:25:11Z\"," +
				"        \"created_by\": \"/bin/sh -c #(nop) ARG NAME VERSION\"," +
				"        \"empty_layer\": true" +
				"      }," +
				"      {" +
				"        \"comment\": \"FROM registry.fedoraproject.org/fedora:38\"," +
				"        \"created\": \"2026-04-08T21:30:48Z\"," +
				"        \"created_by\": \"|2 NAME=fedora-toolbox VERSION=38 /bin/sh -c dnf clean all\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5abdeba6b7f0b1443294df0e64db58f6f7a3064d901e4c9cef488c95976d5b46\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"io.buildah.version\": \"1.33.7\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"localhost/fedora-toolbox:38\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"localhost/fedora-toolbox:38\"" +
				"    ]," +
				"    \"Size\": 1118076168," +
				"    \"VirtualSize\": 1118076168" +
				"  }," +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f38container\"," +
				"        \"FGC=f38\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-01-25T06:49:26Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2024-01-25T06:49:26Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"3eda970c1318e5ea564134338746c04d86ae9817c3a0defdc138b8ccdf787285\"," +
				"    \"Labels\": {" +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"Size\": 186903609," +
				"    \"VirtualSize\": 186903609" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "5abdeba6b7f0b1443294df0e64db58f6f7a3064d901e4c9cef488c95976d5b46",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox": "true",
						"com.redhat.component":          "fedora-toolbox",
						"io.buildah.version":            "1.33.7",
						"license":                       "MIT",
						"name":                          "fedora-toolbox",
						"vendor":                        "Fedora Project",
						"version":                       "38",
					},
					namesHistory: []string{
						"localhost/fedora-toolbox:38",
					},
					repoTags: []string{
						"localhost/fedora-toolbox:38",
					},
				},
				{
					id:       "3eda970c1318e5ea564134338746c04d86ae9817c3a0defdc138b8ccdf787285",
					isToolbx: false,
					labels: map[string]string{
						"license": "MIT",
						"name":    "fedora",
						"vendor":  "Fedora Project",
						"version": "38",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora:38",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora:38",
					},
				},
			},
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38 locally built without a name",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f38container\"," +
				"        \"FGC=f38\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"io.buildah.version\": \"1.33.7\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2026-04-08T21:30:47Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2024-01-25T06:49:26Z\"" +
				"      }," +
				"      {" +
				"        \"created\": \"2026-04-08T21:25:10Z\"," +
				"        \"created_by\": \"/bin/sh -c #(nop) ARG NAME\"," +
				"        \"empty_layer\": true" +
				"      }," +
				"      {" +
				"        \"created\": \"2026-04-08T21:25:11Z\"," +
				"        \"created_by\": \"/bin/sh -c #(nop) ARG NAME VERSION\"," +
				"        \"empty_layer\": true" +
				"      }," +
				"      {" +
				"        \"comment\": \"FROM registry.fedoraproject.org/fedora:38\"," +
				"        \"created\": \"2026-04-08T21:30:48Z\"," +
				"        \"created_by\": \"|2 NAME=fedora-toolbox VERSION=38 /bin/sh -c dnf clean all\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"77b838402716b59256e007a09f66a45b3bbeda584d6996dc451916d71ad7c865\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"io.buildah.version\": \"1.33.7\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"docker.io/library/0200639f72f073877d2db4a7f43d38ce51a4ef5c5add7fdae62e39da37cf605a-tmp:latest\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": []," +
				"    \"Size\": 1118078214," +
				"    \"VirtualSize\": 1118078214" +
				"  }," +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f38container\"," +
				"        \"FGC=f38\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-01-25T06:49:26Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2024-01-25T06:49:26Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"3eda970c1318e5ea564134338746c04d86ae9817c3a0defdc138b8ccdf787285\"," +
				"    \"Labels\": {" +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"Size\": 186903609," +
				"    \"VirtualSize\": 186903609" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "77b838402716b59256e007a09f66a45b3bbeda584d6996dc451916d71ad7c865",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox": "true",
						"com.redhat.component":          "fedora-toolbox",
						"io.buildah.version":            "1.33.7",
						"license":                       "MIT",
						"name":                          "fedora-toolbox",
						"vendor":                        "Fedora Project",
						"version":                       "38",
					},
					namesHistory: []string{
						"docker.io/library/0200639f72f073877d2db4a7f43d38ce51a4ef5c5add7fdae62e39da37cf605a-tmp:latest",
					},
					repoTags: []string{},
				},
				{
					id:       "3eda970c1318e5ea564134338746c04d86ae9817c3a0defdc138b8ccdf787285",
					isToolbx: false,
					labels: map[string]string{
						"license": "MIT",
						"name":    "fedora",
						"vendor":  "Fedora Project",
						"version": "38",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora:38",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora:38",
					},
				},
			},
		},
		{
			name: "podman 4.9.4, fedora-toolbox:39",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"DISTTAG=f39container\"," +
				"        \"FGC=f39\"," +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"39\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-11-26T07:51:25Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"Created by Image Factory\"," +
				"        \"created\": \"2024-11-26T07:51:25Z\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"20a55d9dc10ebf8483727242074d0c50086a82c987b408ec4fc281ea23d7558b\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"39\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:39\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:39\"" +
				"    ]," +
				"    \"Size\": 2105060434," +
				"    \"VirtualSize\": 2105060434" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "20a55d9dc10ebf8483727242074d0c50086a82c987b408ec4fc281ea23d7558b",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox": "true",
						"license":                       "MIT",
						"name":                          "fedora",
						"vendor":                        "Fedora Project",
						"version":                       "39",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:39",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:39",
					},
				},
			},
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Size\": 2204748042," +
				"    \"VirtualSize\": 2204748042" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
			},
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40 with its name removed",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": []," +
				"    \"Size\": 2204748042," +
				"    \"VirtualSize\": 2204748042" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
					repoTags: []string{},
				},
			},
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40 and its copy with different registries",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"localhost/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"localhost/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Size\": 2204748042," +
				"    \"VirtualSize\": 2204748042" +
				"  }," +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"localhost/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"localhost/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Size\": 2204748042," +
				"    \"VirtualSize\": 2204748042" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
						"localhost/fedora-toolbox:40",
					},
					repoTags: []string{
						"localhost/fedora-toolbox:40",
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
				{
					id:       "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
						"localhost/fedora-toolbox:40",
					},
					repoTags: []string{
						"localhost/fedora-toolbox:40",
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
			},
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40 and its copy with different tags",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-copy\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-copy\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Size\": 2204748042," +
				"    \"VirtualSize\": 2204748042" +
				"  }," +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-copy\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-copy\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Size\": 2204748042," +
				"    \"VirtualSize\": 2204748042" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
						"registry.fedoraproject.org/fedora-toolbox:40-copy",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:40-copy",
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
				{
					id:       "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
						"registry.fedoraproject.org/fedora-toolbox:40-copy",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:40-copy",
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
			},
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40, fedora-toolbox:40-aarch64, postgres:latest",
			data: "" +
				"[" +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"postgres\"" +
				"      ]," +
				"      \"Entrypoint\": [" +
				"        \"docker-entrypoint.sh\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"GOSU_VERSION=1.19\"," +
				"        \"LANG=en_US.utf8\"," +
				"        \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/postgresql/18/bin\"," +
				"        \"PGDATA=/var/lib/postgresql/18/docker\"," +
				"        \"PG_MAJOR=18\"," +
				"        \"PG_VERSION=18.0-1.pgdg13+3\"" +
				"      ]," +
				"      \"ExposedPorts\": {" +
				"        \"5432/tcp\": {}" +
				"      }," +
				"      \"StopSignal\": \"SIGINT\"," +
				"      \"Volumes\": {" +
				"        \"/var/lib/postgresql\": {}" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-09-25T18:22:35Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"comment\": \"debuerreotype 0.16\"," +
				"        \"created\": \"2025-09-25T18:22:35Z\"," +
				"        \"created_by\": \"# debian.sh --arch 'amd64' out/ 'trixie' '@1759104000'\"" +
				"      }," +
				"      {" +
				"        \"comment\": \"buildkit.dockerfile.v0\"," +
				"        \"created\": \"2025-09-25T18:22:35Z\"," +
				"        \"created_by\": \"ENV GOSU_VERSION=1.19\"," +
				"        \"empty_layer\": true" +
				"      }," +
				"      {" +
				"        \"comment\": \"buildkit.dockerfile.v0\"," +
				"        \"created\": \"2025-09-25T18:22:35Z\"," +
				"        \"created_by\": \"ENV LANG=en_US.utf8\"," +
				"        \"empty_layer\": true" +
				"      }," +
				"      {" +
				"        \"comment\": \"buildkit.dockerfile.v0\"," +
				"        \"created\": \"2025-09-25T18:22:35Z\"," +
				"        \"created_by\": \"CMD [\\\"postgres\\\"]\"," +
				"        \"empty_layer\": true" +
				"      }" +
				"    ]," +
				"    \"Id\": \"194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd\"," +
				"    \"Labels\": null," +
				"    \"NamesHistory\": [" +
				"      \"docker.io/library/postgres:latest\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"docker.io/library/postgres:latest\"" +
				"    ]," +
				"    \"Size\": 463150524," +
				"    \"VirtualSize\": 463150524" +
				"  }," +
				"  {" +
				"    \"Architecture\": \"amd64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Size\": 2204748042," +
				"    \"VirtualSize\": 2204748042" +
				"  }," +
				"  {" +
				"    \"Architecture\": \"arm64\"," +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"/bin/bash\"" +
				"      ]," +
				"      \"Env\": [" +
				"        \"container=oci\"" +
				"      ]," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"io.buildah.version\": \"1.39.2\"," +
				"        \"license\": \"MIT\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.license\": \"MIT\"," +
				"        \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"        \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"        \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"        \"org.opencontainers.image.version\": \"40\"," +
				"        \"vendor\": \"Fedora Project\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"    \"History\": [" +
				"      {" +
				"        \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"        \"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"        \"created_by\": \"KIWI 10.2.19\"" +
				"      }" +
				"    ]," +
				"    \"Id\": \"49008798958db0411330d8b6a996c03bdab1d1b10d139e7df2b75ed3298f84a7\"," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"io.buildah.version\": \"1.39.2\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.license\": \"MIT\"," +
				"      \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"      \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"      \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"      \"org.opencontainers.image.version\": \"40\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"40\"" +
				"    }," +
				"    \"NamesHistory\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-aarch64\"" +
				"    ]," +
				"    \"Os\": \"linux\"," +
				"    \"Parent\": \"\"," +
				"    \"RepoTags\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-aarch64\"" +
				"    ]," +
				"    \"Size\": 1253521285," +
				"    \"VirtualSize\": 1253521285" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd",
					isToolbx: false,
					labels:   nil,
					namesHistory: []string{
						"docker.io/library/postgres:latest",
					},
					repoTags: []string{
						"docker.io/library/postgres:latest",
					},
				},
				{
					id:       "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
				{
					id:       "49008798958db0411330d8b6a996c03bdab1d1b10d139e7df2b75ed3298f84a7",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.39.2",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "40",
						"vendor":                           "Fedora Project",
						"version":                          "40",
					},
					namesHistory: []string{
						"registry.fedoraproject.org/fedora-toolbox:40-aarch64",
					},
					repoTags: []string{
						"registry.fedoraproject.org/fedora-toolbox:40-aarch64",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.data)
			var images []imageInspect
			err := json.Unmarshal(data, &images)
			assert.NoError(t, err)
			expectsCount := len(tc.expects)
			require.Len(t, images, expectsCount)

			for i, expect := range tc.expects {
				image := images[i]
				assert.Equal(t, expect.id, image.ID())
				assert.Equal(t, expect.isToolbx, image.IsToolbx())

				labels := image.Labels()
				if labels != nil {
					labels["foo"] = "bar"
				}

				assert.Equal(t, expect.labels, image.Labels())

				names := image.Names()
				namesCount := len(names)
				if namesCount != 0 {
					names[0] = "foo/bar"
				}

				assert.Equal(t, expect.namesHistory, image.Names())

				switch namesCount {
				case 0:
					assert.Panics(t, func() { _ = image.Name() })
				default:
					assert.Equal(t, expect.namesHistory[0], image.Name())
				}

				repoTags := image.RepoTags()
				if len(repoTags) != 0 {
					repoTags[0] = "foo/bar"
				}

				assert.Equal(t, expect.repoTags, image.RepoTags())
			}
		})
	}
}

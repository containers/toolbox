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
)

func TestImageImages(t *testing.T) {
	type expect struct {
		id       string
		isToolbx bool
		labels   map[string]string
		names    []string
	}

	testCases := []struct {
		name        string
		data        string
		expects     []expect
		imagesCount int
	}{
		{
			name: "podman 1.1.2, fedora-toolbox:29",
			data: "" +
				"[" +
				"  {" +
				"    \"id\": \"4a6adf1f2a96adf5ea0c02b61f9fa574306f77fc522f39c2ce6bb164daead882\"," +
				"    \"names\": [" +
				"      \"registry.fedoraproject.org/f29/fedora-toolbox:29\"" +
				"    ]," +
				"    \"created\": \"2019-11-06T11:36:10.586442Z\"," +
				"    \"size\": 437590766" +
				"  }," +
				"  {" +
				"    \"id\": \"8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4\"," +
				"    \"names\": [" +
				"      \"localhost/fedora-toolbox-user:29\"" +
				"    ]," +
				"    \"created\": \"2025-10-06T09:10:26.394950134Z\"," +
				"    \"size\": 437967580" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "4a6adf1f2a96adf5ea0c02b61f9fa574306f77fc522f39c2ce6bb164daead882",
					isToolbx: false,
					labels:   nil,
					names: []string{
						"registry.fedoraproject.org/f29/fedora-toolbox:29",
					},
				},
				{
					id:       "8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4",
					isToolbx: false,
					labels:   nil,
					names: []string{
						"localhost/fedora-toolbox-user:29",
					},
				},
			},
			imagesCount: 2,
		},
		{
			name: "podman 1.8.0, fedora-toolbox:30",
			data: "" +
				"[" +
				"  {" +
				"    \"id\": \"c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8\"," +
				"    \"names\": [" +
				"      \"registry.fedoraproject.org/f30/fedora-toolbox:30\"" +
				"    ]," +
				"    \"created\": \"2020-03-03T10:32:20.503064Z\"," +
				"    \"size\": 404259600," +
				"    \"readonly\": false," +
				"    \"history\": []" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8",
					isToolbx: false,
					labels:   nil,
					names: []string{
						"registry.fedoraproject.org/f30/fedora-toolbox:30",
					},
				},
			},
			imagesCount: 1,
		},
		{
			name: "podman 2.2.1, fedora-toolbox:32",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6\"," +
				"    \"ParentId\": \"\"," +
				"    \"Size\": 371369047," +
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
				"    \"Containers\": 1," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:32\"" +
				"    ]," +
				"    \"Created\": 1619352159," +
				"    \"CreatedAt\": \"2021-04-25T12:02:39Z\"" +
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:32",
					},
				},
			},
			imagesCount: 1,
		},
		{
			name: "podman 3.4.7, fedora-toolbox:35",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 495703728," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 495703728," +
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
				"    \"Containers\": 1," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"    ]," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"    ]," +
				"    \"Created\": 1667999803," +
				"    \"CreatedAt\": \"2022-11-09T13:16:43Z\"" +
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:35",
					},
				},
			},
			imagesCount: 1,
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 1748730844," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 1748730844," +
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
				"    \"Containers\": 1," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"    ]," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"    ]," +
				"    \"Created\": 1706814575," +
				"    \"CreatedAt\": \"2024-02-01T19:09:35Z\"" +
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:38",
					},
				},
			},
			imagesCount: 1,
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38 locally built",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"5abdeba6b7f0b1443294df0e64db58f6f7a3064d901e4c9cef488c95976d5b46\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 1118076168," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 1118076168," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"io.buildah.version\": \"1.33.7\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"Containers\": 0," +
				"    \"Names\": [" +
				"      \"localhost/fedora-toolbox:38\"" +
				"    ]," +
				"    \"History\": [" +
				"      \"localhost/fedora-toolbox:38\"" +
				"    ]," +
				"    \"Created\": 1775683847," +
				"    \"CreatedAt\": \"2026-04-08T21:30:47Z\"" +
				"  }," +
				"  {" +
				"    \"Id\": \"3eda970c1318e5ea564134338746c04d86ae9817c3a0defdc138b8ccdf787285\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 186903609," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 186903609," +
				"    \"Labels\": {" +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"Containers\": 0," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"Created\": 1706165366," +
				"    \"CreatedAt\": \"2024-01-25T06:49:26Z\"" +
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
					names: []string{
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
					names: []string{
						"registry.fedoraproject.org/fedora:38",
					},
				},
			},
			imagesCount: 2,
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38 locally built without a name",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"77b838402716b59256e007a09f66a45b3bbeda584d6996dc451916d71ad7c865\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 1118078214," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 1118078214," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"com.redhat.component\": \"fedora-toolbox\"," +
				"      \"io.buildah.version\": \"1.33.7\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora-toolbox\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"Containers\": 0," +
				"    \"Dangling\": true," +
				"    \"History\": [" +
				"      \"docker.io/library/0200639f72f073877d2db4a7f43d38ce51a4ef5c5add7fdae62e39da37cf605a-tmp:latest\"" +
				"    ]," +
				"    \"Created\": 1775683847," +
				"    \"CreatedAt\": \"2026-04-08T21:30:47Z\"" +
				"  }," +
				"  {" +
				"    \"Id\": \"3eda970c1318e5ea564134338746c04d86ae9817c3a0defdc138b8ccdf787285\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 186903609," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 186903609," +
				"    \"Labels\": {" +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"38\"" +
				"    }," +
				"    \"Containers\": 0," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora:38\"" +
				"    ]," +
				"    \"Created\": 1706165366," +
				"    \"CreatedAt\": \"2024-01-25T06:49:26Z\"" +
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
					names: nil,
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
					names: []string{
						"registry.fedoraproject.org/fedora:38",
					},
				},
			},
			imagesCount: 2,
		},
		{
			name: "podman 4.9.4, fedora-toolbox:39",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"20a55d9dc10ebf8483727242074d0c50086a82c987b408ec4fc281ea23d7558b\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2105060434," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2105060434," +
				"    \"Labels\": {" +
				"      \"com.github.containers.toolbox\": \"true\"," +
				"      \"license\": \"MIT\"," +
				"      \"name\": \"fedora\"," +
				"      \"vendor\": \"Fedora Project\"," +
				"      \"version\": \"39\"" +
				"    }," +
				"    \"Containers\": 1," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:39\"" +
				"    ]," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:39\"" +
				"    ]," +
				"    \"Created\": 1732607485," +
				"    \"CreatedAt\": \"2024-11-26T07:51:25Z\"" +
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:39",
					},
				},
			},
			imagesCount: 1,
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2204748042," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2204748042," +
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
				"    \"Containers\": 1," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
			},
			imagesCount: 1,
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40 with its name removed",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2204748042," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2204748042," +
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
				"    \"Containers\": 1," +
				"    \"Dangling\": true," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
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
					names: nil,
				},
			},
			imagesCount: 1,
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40 and its copy with different registries",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2204748042," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2204748042," +
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
				"    \"Containers\": 1," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"localhost/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"localhost/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
				"  }," +
				"  {" +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2204748042," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2204748042," +
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
				"    \"Containers\": 1," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"localhost/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"localhost/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
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
					names: []string{
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
					names: []string{
						"localhost/fedora-toolbox:40",
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
			},
			imagesCount: 2,
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40 and its copy with different tags",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2204748042," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2204748042," +
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
				"    \"Containers\": 1," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-test\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-test\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
				"  }," +
				"  {" +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2204748042," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2204748042," +
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
				"    \"Containers\": 1," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-test\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-test\"," +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:40-test",
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:40-test",
						"registry.fedoraproject.org/fedora-toolbox:40",
					},
				},
			},
			imagesCount: 2,
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40, fedora-toolbox:40-aarch64, postgres:latest",
			data: "" +
				"[" +
				"  {" +
				"    \"Id\": \"194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 463150524," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 463150524," +
				"    \"Labels\": null," +
				"    \"Containers\": 0," +
				"    \"History\": [" +
				"      \"docker.io/library/postgres:latest\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"docker.io/library/postgres:latest\"" +
				"    ]," +
				"    \"Created\": 1758824555," +
				"    \"CreatedAt\": \"2025-09-25T18:22:35Z\"" +
				"  }," +
				"  {" +
				"    \"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 2204748042," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 2204748042," +
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
				"    \"Containers\": 1," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
				"  }," +
				"  {" +
				"    \"Id\": \"49008798958db0411330d8b6a996c03bdab1d1b10d139e7df2b75ed3298f84a7\"," +
				"    \"ParentId\": \"\"," +
				"    \"RepoTags\": null," +
				"    \"Size\": 1253521285," +
				"    \"SharedSize\": 0," +
				"    \"VirtualSize\": 1253521285," +
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
				"    \"Containers\": 0," +
				"    \"History\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-aarch64\"" +
				"    ]," +
				"    \"Names\": [" +
				"      \"registry.fedoraproject.org/fedora-toolbox:40-aarch64\"" +
				"    ]," +
				"    \"Created\": 1747122540," +
				"    \"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
				"  }" +
				"]",
			expects: []expect{
				{
					id:       "194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd",
					isToolbx: false,
					labels:   nil,
					names: []string{
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
					names: []string{
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
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:40-aarch64",
					},
				},
			},
			imagesCount: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.data)
			var images []imageImages
			err := json.Unmarshal(data, &images)
			assert.NoError(t, err)
			assert.Len(t, images, tc.imagesCount)

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

				assert.Equal(t, expect.names, image.Names())

				switch namesCount {
				case 0:
					assert.Panics(t, func() { _ = image.Name() })
				case 1:
					assert.Equal(t, expect.names[0], image.Name())
				default:
					assert.Panics(t, func() { _ = image.Name() })

					flattenedImages := image.flattenNames(false)
					for j, flattenedImage := range flattenedImages {
						assert.Equal(t, expect.names[j], flattenedImage.Name())
					}

				}
			}
		})
	}
}

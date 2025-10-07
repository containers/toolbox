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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageImages(t *testing.T) {
	type expect struct {
		id       string
		isToolbx bool
		labels   map[string]string
		names    []string
		repoTags []string
	}

	testCases := []struct {
		name     string
		data     string
		expects  []expect
		imageCnt int
	}{
		{
			name: "podman 1.1.2, fedora-toolbox:29",
			data: "" +
				"[" +
				"	{" +
				"		\"id\": \"4a6adf1f2a96adf5ea0c02b61f9fa574306f77fc522f39c2ce6bb164daead882\"," +
				"		\"names\": [" +
				"			\"registry.fedoraproject.org/f29/fedora-toolbox:29\"" +
				"		]," +
				"		\"digest\": \"sha256:f324a546a5e894af041eea47f8f2392bf2a9e2d88ee77199b25a129174b7a0e1\"," +
				"		\"created\": \"2019-11-06T11:36:10.586442Z\"," +
				"		\"size\": 437590766" +
				"	}," +
				"	{" +
				"		\"id\": \"8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4\"," +
				"		\"names\": [" +
				"			\"localhost/fedora-toolbox-dkricka:29\"" +
				"		]," +
				"		\"digest\": \"sha256:00d5cea66e5b3b043ced7ba922198bfea469aa5f7fd2e8f7993403c3e9cb30bf\"," +
				"		\"created\": \"2025-10-06T09:10:26.394950134Z\"," +
				"		\"size\": 437967580" +
				"	}" +
				"]",
			expects: []expect{
				{
					id:       "4a6adf1f2a96adf5ea0c02b61f9fa574306f77fc522f39c2ce6bb164daead882",
					isToolbx: false,
					labels:   nil,
					names: []string{
						"registry.fedoraproject.org/f29/fedora-toolbox:29",
					},
					repoTags: nil,
				},
				{
					id:       "8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4",
					isToolbx: false,
					labels:   nil,
					names: []string{
						"localhost/fedora-toolbox-dkricka:29",
					},
					repoTags: nil,
				},
			},
			imageCnt: 2,
		},
		{
			name: "podman 1.8.0, fedora-toolbox:30",
			data: "" +
				"[" +
				"	{" +
				"		\"id\": \"c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8\"," +
				"		\"names\": [" +
				"			\"registry.fedoraproject.org/f30/fedora-toolbox:30\"" +
				"		]," +
				"		\"digest\": \"sha256:9deba34b2c7fd7ccd2057c8abc03108494b2c57da5d4e1d8b3e41e6164dd805a\"," +
				"		\"digests\": [" +
				"			\"sha256:9deba34b2c7fd7ccd2057c8abc03108494b2c57da5d4e1d8b3e41e6164dd805a\"," +
				"			\"sha256:45afda5831c50fe1e6171d41dcf9c695e3fb04ffe6ce1cabb2c690ea5adb9cf4\"" +
				"		]," +
				"		\"created\": \"2020-03-03T10:32:20.503064Z\"," +
				"		\"size\": 404259600," +
				"		\"readonly\": false," +
				"		\"history\": []" +
				"	}" +
				"]",
			expects: []expect{
				{
					id:       "c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8",
					isToolbx: false,
					labels:   nil,
					names: []string{
						"registry.fedoraproject.org/f30/fedora-toolbox:30",
					},
					repoTags: nil,
				},
			},
			imageCnt: 1,
		},
		{
			name: "podman 2.2.1, fedora-toolbox:32",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6\"," +
				"		\"ParentId\": \"\"," +
				"		\"Size\": 371369047," +
				"		\"Labels\": {" +
				"			\"architecture\": \"x86_64\"," +
				"			\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"			\"build-date\": \"2021-04-25T12:01:08.917265\"," +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"com.github.debarshiray.toolbox\": \"true\"," +
				"			\"com.redhat.build-host\": \"osbs-node02.iad2.fedoraproject.org\"," +
				"			\"com.redhat.component\": \"fedora-toolbox\"," +
				"			\"distribution-scope\": \"public\"," +
				"			\"license\": \"MIT\"," +
				"			\"maintainer\": \"Debarshi Ray \u003crishi@fedoraproject.org\u003e\"," +
				"			\"name\": \"fedora-toolbox\"," +
				"			\"release\": \"14\"," +
				"			\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"			\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"			\"vcs-ref\": \"c81107b8acbbd44abafd31793d71141d160e70c1\"," +
				"			\"vcs-type\": \"git\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"32\"" +
				"		}," +
				"		\"Containers\": 1," +
				"		\"Names\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:32\"" +
				"		]," +
				"		\"Digest\": \"sha256:dd65063ba3ee25991de8e3e13b6a44b67e4d1d4f4875153591c55f2f0cb78515\"," +
				"		\"Digests\": [" +
				"			\"sha256:dd65063ba3ee25991de8e3e13b6a44b67e4d1d4f4875153591c55f2f0cb78515\"," +
				"			\"sha256:1d1077ce12b6d45990709ab81543884838d71a3b2014d023b82fd465a3ce4e23\"" +
				"		]," +
				"		\"Created\": 1619352159," +
				"		\"CreatedAt\": \"2021-04-25T12:02:39Z\"" +
				"	}" +
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
						"com.redhat.build-host":          "osbs-node02.iad2.fedoraproject.org",
						"com.redhat.component":           "fedora-toolbox",
						"distribution-scope":             "public",
						"license":                        "MIT",
						"maintainer":                     "Debarshi Ray \u003crishi@fedoraproject.org\u003e",
						"name":                           "fedora-toolbox",
						"release":                        "14",
						"summary":                        "Base image for creating Fedora toolbox containers",
						"usage":                          "This image is meant to be used with the toolbox command",
						"vcs-ref":                        "c81107b8acbbd44abafd31793d71141d160e70c1",
						"vcs-type":                       "git",
						"vendor":                         "Fedora Project",
						"version":                        "32",
					},
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:32",
					},
					repoTags: nil,
				},
			},
			imageCnt: 1,
		},
		{
			name: "podman 3.4.7, fedora-toolbox:35",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227\"," +
				"		\"ParentId\": \"\"," +
				"		\"RepoTags\": null," +
				"		\"RepoDigests\": [" +
				"			\"sha256:b5f419099423fae80421cda692b00e301894101575806f44558e6c9c911769e6\"," +
				"			\"sha256:bdab10512284f81480235e89939babff43b2a8ca5af870dbebdeb127c6533568\"" +
				"		]," +
				"		\"Size\": 495703728," +
				"		\"SharedSize\": 0," +
				"		\"VirtualSize\": 495703728," +
				"		\"Labels\": {" +
				"			\"architecture\": \"x86_64\"," +
				"			\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"			\"build-date\": \"2022-11-09T13:15:18.185125\"," +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"com.github.debarshiray.toolbox\": \"true\"," +
				"			\"com.redhat.build-host\": \"osbs-node02.iad2.fedoraproject.org\"," +
				"			\"com.redhat.component\": \"fedora-toolbox\"," +
				"			\"distribution-scope\": \"public\"," +
				"			\"license\": \"MIT\"," +
				"			\"maintainer\": \"Debarshi Ray \u003crishi@fedoraproject.org\u003e\"," +
				"			\"name\": \"fedora-toolbox\"," +
				"			\"release\": \"18\"," +
				"			\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"			\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"			\"vcs-ref\": \"274bdf9053aa6e56114c9f4fca26f09a89a72ad5\"," +
				"			\"vcs-type\": \"git\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"35\"" +
				"		}," +
				"		\"Containers\": 1," +
				"		\"Names\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"		]," +
				"		\"Digest\": \"sha256:b5f419099423fae80421cda692b00e301894101575806f44558e6c9c911769e6\"," +
				"		\"History\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"		]," +
				"		\"Created\": 1667999803," +
				"		\"CreatedAt\": \"2022-11-09T13:16:43Z\"" +
				"	}" +
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
						"com.redhat.build-host":          "osbs-node02.iad2.fedoraproject.org",
						"com.redhat.component":           "fedora-toolbox",
						"distribution-scope":             "public",
						"license":                        "MIT",
						"maintainer":                     "Debarshi Ray \u003crishi@fedoraproject.org\u003e",
						"name":                           "fedora-toolbox",
						"release":                        "18",
						"summary":                        "Base image for creating Fedora toolbox containers",
						"usage":                          "This image is meant to be used with the toolbox command",
						"vcs-ref":                        "274bdf9053aa6e56114c9f4fca26f09a89a72ad5",
						"vcs-type":                       "git",
						"vendor":                         "Fedora Project",
						"version":                        "35",
					},
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:35",
					},
					repoTags: nil,
				},
			},
			imageCnt: 1,
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce\"," +
				"		\"ParentId\": \"\"," +
				"		\"RepoTags\": null," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:8a52bb3a18508c15cd12e8770c0ce17c910c2bb96aabc2ffc4b0ceb669cef935\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:efde4efa9f18e619f62a14849959a436fc8062b2248aca794a0f4957b1d2eeca\"" +
				"		]," +
				"		\"Size\": 1748730844," +
				"		\"SharedSize\": 0," +
				"		\"VirtualSize\": 1748730844," +
				"		\"Labels\": {" +
				"			\"architecture\": \"x86_64\"," +
				"			\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"			\"build-date\": \"2024-02-01T19:07:19.783946\"," +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"com.redhat.build-host\": \"osbs-node01.iad2.fedoraproject.org\"," +
				"			\"com.redhat.component\": \"fedora-toolbox\"," +
				"			\"distribution-scope\": \"public\"," +
				"			\"license\": \"MIT\"," +
				"			\"maintainer\": \"Debarshi Ray \u003crishi@fedoraproject.org\u003e\"," +
				"			\"name\": \"fedora-toolbox\"," +
				"			\"release\": \"20\"," +
				"			\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"			\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"			\"vcs-ref\": \"c36ca58a44f947077b7c62e17d7b4fe4cd1b5797\"," +
				"			\"vcs-type\": \"git\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"38\"" +
				"		}," +
				"		\"Containers\": 1," +
				"		\"Names\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"		]," +
				"		\"Digest\": \"sha256:efde4efa9f18e619f62a14849959a436fc8062b2248aca794a0f4957b1d2eeca\"," +
				"		\"History\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"		]," +
				"		\"Created\": 1706814575," +
				"		\"CreatedAt\": \"2024-02-01T19:09:35Z\"" +
				"	}" +
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
						"com.redhat.build-host":         "osbs-node01.iad2.fedoraproject.org",
						"com.redhat.component":          "fedora-toolbox",
						"distribution-scope":            "public",
						"license":                       "MIT",
						"maintainer":                    "Debarshi Ray \u003crishi@fedoraproject.org\u003e",
						"name":                          "fedora-toolbox",
						"release":                       "20",
						"summary":                       "Base image for creating Fedora toolbox containers",
						"usage":                         "This image is meant to be used with the toolbox command",
						"vcs-ref":                       "c36ca58a44f947077b7c62e17d7b4fe4cd1b5797",
						"vcs-type":                      "git",
						"vendor":                        "Fedora Project",
						"version":                       "38",
					},
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:38",
					},
					repoTags: nil,
				},
			},
			imageCnt: 1,
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"		\"ParentId\": \"\"," +
				"		\"RepoTags\": null," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:3d8c46e803fd184e9318ee90d4e6f1ad0b591a1c49f4de683d4f8748c0b95c30\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:c3cc9836c9e55475f85496f9369e925164f1b7cd832b55e22c13b74840576c31\"" +
				"		]," +
				"		\"Size\": 2204748042," +
				"		\"SharedSize\": 0," +
				"		\"VirtualSize\": 2204748042," +
				"		\"Labels\": {" +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"io.buildah.version\": \"1.39.2\"," +
				"			\"license\": \"MIT\"," +
				"			\"name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.license\": \"MIT\"," +
				"			\"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"			\"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"			\"org.opencontainers.image.version\": \"40\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"40\"" +
				"		}," +
				"		\"Containers\": 1," +
				"		\"Digest\": \"sha256:c3cc9836c9e55475f85496f9369e925164f1b7cd832b55e22c13b74840576c31\"," +
				"		\"History\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"		]," +
				"		\"Names\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"		]," +
				"		\"Created\": 1747122540," +
				"		\"CreatedAt\": \"2025-05-13T07:49:00Z\"" +
				"	}" +
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
					repoTags: nil,
				},
			},
			imageCnt: 1,
		},
		{
			name: "podman 5.6.1, fedora-toolbox:41",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"b65ee158f921088db52e8b98a6f3267de75324c8a6a04afc9ff095338c40e59b\"," +
				"		\"ParentId\": \"\"," +
				"		\"RepoTags\": null," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:0a51adc6bab55d49ff00da8aaad81ca1f02315511ed23d55ee5bbbe1a976a663\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:8599f0c0d421c0dc01c4b7fb1c07b2780c0ab1931d0f22dd7b6da3b93ff6b77b\"" +
				"		]," +
				"		\"Size\": 2308518290," +
				"		\"SharedSize\": 0," +
				"		\"VirtualSize\": 2308518290," +
				"		\"Labels\": {" +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"io.buildah.version\": \"1.41.4\"," +
				"			\"license\": \"MIT\"," +
				"			\"name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.license\": \"MIT\"," +
				"			\"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"			\"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"			\"org.opencontainers.image.version\": \"41\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"41\"" +
				"		}," +
				"		\"Containers\": 1," +
				"		\"Digest\": \"sha256:0a51adc6bab55d49ff00da8aaad81ca1f02315511ed23d55ee5bbbe1a976a663\"," +
				"		\"History\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:41\"" +
				"		]," +
				"		\"Names\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:41\"" +
				"		]," +
				"		\"Created\": 1759729685," +
				"		\"CreatedAt\": \"2025-10-06T05:48:05Z\"" +
				"	}" +
				"]",
			expects: []expect{
				{
					id:       "b65ee158f921088db52e8b98a6f3267de75324c8a6a04afc9ff095338c40e59b",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.41.4",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "41",
						"vendor":                           "Fedora Project",
						"version":                          "41",
					},
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:41",
					},
					repoTags: nil,
				},
			},
			imageCnt: 1,
		},
		{
			name: "podman 5.6.2, fedora-toolbox:42, docker.io-postgres:18, fedora-toolbox:42-aarch64",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"341ada9b9af076546132ae6fe6e328eb777046581e23b1c232bfc71d844f8598\"," +
				"		\"ParentId\": \"\"," +
				"		\"RepoTags\": null," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:2fd640a4b02193f88972845e73e22f3943f9b69c69f13cd17c50cec098bd0715\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:83540b1b86984bc56e85c0607ec0cc2469c45f6716259bc98668cacd9bbce48a\"" +
				"		]," +
				"		\"Size\": 2136421808," +
				"		\"SharedSize\": 0," +
				"		\"VirtualSize\": 2136421808," +
				"		\"Labels\": {" +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"io.buildah.version\": \"1.41.5\"," +
				"			\"license\": \"MIT\"," +
				"			\"name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.license\": \"MIT\"," +
				"			\"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"			\"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"			\"org.opencontainers.image.version\": \"42\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"42\"" +
				"		}," +
				"		\"Containers\": 1," +
				"		\"Digest\": \"sha256:2fd640a4b02193f88972845e73e22f3943f9b69c69f13cd17c50cec098bd0715\"," +
				"		\"History\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:42\"" +
				"		]," +
				"		\"Names\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:42\"" +
				"		]," +
				"		\"Created\": 1759733390," +
				"		\"CreatedAt\": \"2025-10-06T06:49:50Z\"" +
				"	}," +
				"	{" +
				"		\"Id\": \"194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd\"," +
				"		\"ParentId\": \"\"," +
				"		\"RepoTags\": null," +
				"		\"RepoDigests\": [" +
				"			\"docker.io/library/postgres@sha256:073e7c8b84e2197f94c8083634640ab37105effe1bc853ca4d5fbece3219b0e8\"," +
				"			\"docker.io/library/postgres@sha256:28f01a051c819681a816dca282088111ade7c44f834dd83cfd044f0548d38c19\"" +
				"		]," +
				"		\"Size\": 463150524," +
				"		\"SharedSize\": 0," +
				"		\"VirtualSize\": 463150524," +
				"		\"Labels\": null," +
				"		\"Containers\": 0," +
				"		\"Digest\": \"sha256:073e7c8b84e2197f94c8083634640ab37105effe1bc853ca4d5fbece3219b0e8\"," +
				"		\"History\": [" +
				"			\"docker.io/library/postgres:latest\"" +
				"		]," +
				"		\"Names\": [" +
				"			\"docker.io/library/postgres:latest\"" +
				"		]," +
				"		\"Created\": 1758824555," +
				"		\"CreatedAt\": \"2025-09-25T18:22:35Z\"" +
				"	}," +
				"	{" +
				"		\"Id\": \"2718f04f884164eb696b9ea011e29bca3222e901e387db80915428f2ed5ca5d7\"," +
				"		\"ParentId\": \"\"," +
				"		\"RepoTags\": null," +
				"		\"RepoDigests\": [" +
				"			\"quay.io/fedora/fedora-toolbox@sha256:8d282ca2e63be6a19e0f82bb0dafe3aa58ea45335be5876d596bee0dc788e679\"" +
				"		]," +
				"		\"Size\": 2103319610," +
				"		\"SharedSize\": 0," +
				"		\"VirtualSize\": 2103319610," +
				"		\"Labels\": {" +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"io.buildah.version\": \"1.41.4\"," +
				"			\"license\": \"MIT\"," +
				"			\"name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.license\": \"MIT\"," +
				"			\"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"			\"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"			\"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"			\"org.opencontainers.image.version\": \"42\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"42\"" +
				"		}," +
				"		\"Containers\": 1," +
				"		\"Digest\": \"sha256:8d282ca2e63be6a19e0f82bb0dafe3aa58ea45335be5876d596bee0dc788e679\"," +
				"		\"History\": [" +
				"			\"quay.io/fedora/fedora-toolbox:42-aarch64\"" +
				"		]," +
				"		\"Names\": [" +
				"			\"quay.io/fedora/fedora-toolbox:42-aarch64\"" +
				"		]," +
				"		\"Created\": 1758524250," +
				"		\"CreatedAt\": \"2025-09-22T06:57:30Z\"" +
				"	}" +
				"]",
			expects: []expect{
				{
					id:       "341ada9b9af076546132ae6fe6e328eb777046581e23b1c232bfc71d844f8598",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.41.5",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "42",
						"vendor":                           "Fedora Project",
						"version":                          "42",
					},
					names: []string{
						"registry.fedoraproject.org/fedora-toolbox:42",
					},
					repoTags: nil,
				},
				{
					id:       "194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd",
					isToolbx: false,
					labels:   nil,
					names: []string{
						"docker.io/library/postgres:latest",
					},
					repoTags: nil,
				},
				{
					id:       "2718f04f884164eb696b9ea011e29bca3222e901e387db80915428f2ed5ca5d7",
					isToolbx: true,
					labels: map[string]string{
						"com.github.containers.toolbox":    "true",
						"io.buildah.version":               "1.41.4",
						"license":                          "MIT",
						"name":                             "fedora-toolbox",
						"org.opencontainers.image.license": "MIT",
						"org.opencontainers.image.name":    "fedora-toolbox",
						"org.opencontainers.image.url":     "https://fedoraproject.org/",
						"org.opencontainers.image.vendor":  "Fedora Project",
						"org.opencontainers.image.version": "42",
						"vendor":                           "Fedora Project",
						"version":                          "42",
					},
					names: []string{
						"quay.io/fedora/fedora-toolbox:42-aarch64",
					},
					repoTags: nil,
				},
			},
			imageCnt: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.data)
			var images []imageImages
			err := json.Unmarshal(data, &images)
			assert.NoError(t, err)
			assert.Len(t, images, tc.imageCnt)

			for i, expect := range tc.expects {
				image := images[i]
				assert.Equal(t, expect.id, image.ID())
				assert.Equal(t, expect.isToolbx, image.IsToolbx())
				assert.Equal(t, expect.labels, image.Labels())
				assert.Equal(t, expect.names, image.Names())
				assert.Equal(t, expect.repoTags, image.RepoTags())
			}
		})
	}
}

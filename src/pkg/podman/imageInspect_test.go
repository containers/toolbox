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

func TestImageInspect(t *testing.T) {
	type expect struct {
		id           string
		isToolbx     bool
		labels       map[string]string
		namesHistory []string
		repoTags     []string
		entrypoint   []string
		envVars      []string
	}

	testCases := []struct {
		name   string
		data   string
		expect expect
	}{
		{
			name: "podman 1.1.2, fedora-toolbox:29",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4\"," +
				"		\"Digest\": \"sha256:00d5cea66e5b3b043ced7ba922198bfea469aa5f7fd2e8f7993403c3e9cb30bf\"," +
				"		\"RepoTags\": [" +
				"			\"localhost/fedora-toolbox-dkricka:29\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"localhost/fedora-toolbox-dkricka@sha256:00d5cea66e5b3b043ced7ba922198bfea469aa5f7fd2e8f7993403c3e9cb30bf\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"Created by Image Factory\"," +
				"		\"Created\": \"2025-10-06T09:10:26.394950134Z\"," +
				"		\"Config\": {" +
				"			\"User\": \"dkricka\"," +
				"			\"Env\": [" +
				"				\"DISTTAG=f29container\"," +
				"				\"FGC=f29\"," +
				"				\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"				\"NAME=fedora-toolbox\"," +
				"				\"VERSION=29\"" +
				"			]," +
				"			\"Cmd\": [" +
				"				\"/bin/sh\"," +
				"				\"-c\"," +
				"				\"/bin/sh\"" +
				"			]," +
				"			\"WorkingDir\": \"/home/dkricka\"," +
				"			\"Labels\": {" +
				"				\"architecture\": \"x86_64\"," +
				"				\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"				\"build-date\": \"2019-11-06T11:24:10.495113\"," +
				"				\"com.github.containers.toolbox\": \"true\"," +
				"				\"com.github.debarshiray.toolbox\": \"true\"," +
				"				\"com.redhat.build-host\": \"osbs-node01.phx2.fedoraproject.org\"," +
				"				\"com.redhat.component\": \"fedora-toolbox\"," +
				"				\"distribution-scope\": \"public\"," +
				"				\"license\": \"MIT\"," +
				"				\"maintainer\": \"Debarshi Ray <rishi@fedoraproject.org>\"," +
				"				\"name\": \"f29/fedora-toolbox\"," +
				"				\"release\": \"9\"," +
				"				\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"				\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"				\"vcs-ref\": \"b2c70792f70eaea3ee62eee45b5e1e2c9ac220ae\"," +
				"				\"vcs-type\": \"git\"," +
				"				\"vendor\": \"Fedora Project\"," +
				"				\"version\": \"29\"" +
				"			}" +
				"		}," +
				"		\"Version\": \"\"," +
				"		\"Author\": \"\"," +
				"		\"Architecture\": \"amd64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 437967580," +
				"		\"VirtualSize\": 437967580," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"				\"LowerDir\": \"/home/dkricka/.local/share/containers/storage/overlay/040ef06274e5a1f9549501844a47261049df0bb7d0210ff7f6622363f5e09ea0/diff:/home/dkricka/.local/share/containers/storage/overlay/da22bd5bfb28630a6045ab8caa0770c672f2878f10a314a3cec19c66b17b7449/diff\"," +
				"				\"MergedDir\": \"/home/dkricka/.local/share/containers/storage/overlay/38e3e0aba22c2877c28126236df72c4b6a03e1d6c3a52816670a802b37d9ad80/merged\"," +
				"				\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/38e3e0aba22c2877c28126236df72c4b6a03e1d6c3a52816670a802b37d9ad80/diff\"," +
				"				\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/38e3e0aba22c2877c28126236df72c4b6a03e1d6c3a52816670a802b37d9ad80/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"				\"sha256:da22bd5bfb28630a6045ab8caa0770c672f2878f10a314a3cec19c66b17b7449\"," +
				"				\"sha256:bc33983a85c035fbc50d595a194e254e93bf3204a7b426d721ae43367e43acbe\"," +
				"				\"sha256:86bf168c691f547a2c9afcdcc39daa173021e935807b8f63ca0779419f7fa1e4\"" +
				"			]" +
				"		}," +
				"		\"Labels\": {" +
				"			\"architecture\": \"x86_64\"," +
				"			\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"			\"build-date\": \"2019-11-06T11:24:10.495113\"," +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"com.github.debarshiray.toolbox\": \"true\"," +
				"			\"com.redhat.build-host\": \"osbs-node01.phx2.fedoraproject.org\"," +
				"			\"com.redhat.component\": \"fedora-toolbox\"," +
				"			\"distribution-scope\": \"public\"," +
				"			\"license\": \"MIT\"," +
				"			\"maintainer\": \"Debarshi Ray <rishi@fedoraproject.org>\"," +
				"			\"name\": \"f29/fedora-toolbox\"," +
				"			\"release\": \"9\"," +
				"			\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"			\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"			\"vcs-ref\": \"b2c70792f70eaea3ee62eee45b5e1e2c9ac220ae\"," +
				"			\"vcs-type\": \"git\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"29\"" +
				"		}," +
				"		\"Annotations\": {}," +
				"		\"ManifestType\": \"application/vnd.oci.image.manifest.v1+json\"," +
				"		\"User\": \"dkricka\"," +
				"		\"History\": [" +
				"			{" +
				"				\"created\": \"2019-02-18T06:48:39Z\"," +
				"				\"comment\": \"Created by Image Factory\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2019-11-06T11:36:10.586442Z\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-10-06T09:10:26.394950134Z\"," +
				"				\"created_by\": \"/bin/sh\"" +
				"			}" +
				"		]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "8c5e0075ddaf4651b73e71dcf420bab922209f246e1165f306f223112a3061a4",
				isToolbx:   true,
				labels: map[string]string{
					"architecture":                   "x86_64",
					"authoritative-source-url":       "registry.fedoraproject.org",
					"build-date":                     "2019-11-06T11:24:10.495113",
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.build-host":          "osbs-node01.phx2.fedoraproject.org",
					"com.redhat.component":           "fedora-toolbox",
					"distribution-scope":             "public",
					"license":                        "MIT",
					"maintainer":                     "Debarshi Ray <rishi@fedoraproject.org>",
					"name":                           "f29/fedora-toolbox",
					"release":                        "9",
					"summary":                        "Base image for creating Fedora toolbox containers",
					"usage":                          "This image is meant to be used with the toolbox command",
					"vcs-ref":                        "b2c70792f70eaea3ee62eee45b5e1e2c9ac220ae",
					"vcs-type":                       "git",
					"vendor":                         "Fedora Project",
					"version":                        "29",
				},
				namesHistory: nil,
				repoTags: []string{
					"localhost/fedora-toolbox-dkricka:29",
				},
				envVars: []string{
					"DISTTAG=f29container",
					"FGC=f29",
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"NAME=fedora-toolbox",
					"VERSION=29",
				},
			},
		},
		{
			name: "podman 1.8.0, fedora-toolbox:30",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8\"," +
				"		\"Digest\": \"sha256:9deba34b2c7fd7ccd2057c8abc03108494b2c57da5d4e1d8b3e41e6164dd805a\"," +
				"		\"RepoTags\": [" +
				"			\"registry.fedoraproject.org/f30/fedora-toolbox:30\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/f30/fedora-toolbox@sha256:45afda5831c50fe1e6171d41dcf9c695e3fb04ffe6ce1cabb2c690ea5adb9cf4\"," +
				"			\"registry.fedoraproject.org/f30/fedora-toolbox@sha256:9deba34b2c7fd7ccd2057c8abc03108494b2c57da5d4e1d8b3e41e6164dd805a\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"\"," +
				"		\"Created\": \"2020-03-03T10:32:20.503064Z\"," +
				"		\"Config\": {" +
				"			\"Env\": [" +
				"				\"DISTTAG=f30container\"," +
				"				\"FGC=f30\"," +
				"				\"container=oci\"," +
				"				\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"				\"NAME=fedora-toolbox\"," +
				"				\"VERSION=30\"" +
				"			]," +
				"			\"Cmd\": [" +
				"				\"/bin/sh\"," +
				"				\"-c\"," +
				"				\"/bin/sh\"" +
				"			]," +
				"			\"Labels\": {" +
				"				\"architecture\": \"x86_64\"," +
				"				\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"				\"build-date\": \"2020-03-03T10:28:18.406675\"," +
				"				\"com.github.containers.toolbox\": \"true\"," +
				"				\"com.github.debarshiray.toolbox\": \"true\"," +
				"				\"com.redhat.build-host\": \"osbs-node02.phx2.fedoraproject.org\"," +
				"				\"com.redhat.component\": \"fedora-toolbox\"," +
				"				\"distribution-scope\": \"public\"," +
				"				\"license\": \"MIT\"," +
				"				\"maintainer\": \"Debarshi Ray <rishi@fedoraproject.org>\"," +
				"				\"name\": \"f30/fedora-toolbox\"," +
				"				\"release\": \"12\"," +
				"				\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"				\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"				\"vcs-ref\": \"41c23bef7659ce5e8798c5f6330e25ba1dc9d5ed\"," +
				"				\"vcs-type\": \"git\"," +
				"				\"vendor\": \"Fedora Project\"," +
				"				\"version\": \"30\"" +
				"			}" +
				"		}," +
				"		\"Version\": \"1.13.1\"," +
				"		\"Author\": \"\"," +
				"		\"Architecture\": \"amd64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 404259600," +
				"		\"VirtualSize\": 404259600," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"				\"LowerDir\": \"/home/dkricka/.local/share/containers/storage/overlay/90d9aeda29c691f01d9c2e29cecd47facbc77768696219a9c8d61f273581f3df/diff\"," +
				"				\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/02cc06abb36aceb54d97c7dd01146473b4d55d05758ef37a2203abaa803b39e2/diff\"," +
				"				\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/02cc06abb36aceb54d97c7dd01146473b4d55d05758ef37a2203abaa803b39e2/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"				\"sha256:90d9aeda29c691f01d9c2e29cecd47facbc77768696219a9c8d61f273581f3df\"," +
				"				\"sha256:43c2a7b14976c4e87d165ec4ba73dbc5e889f05e25d045e9788bc80f3a82d020\"" +
				"			]" +
				"		}," +
				"		\"Labels\": {" +
				"			\"architecture\": \"x86_64\"," +
				"			\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"			\"build-date\": \"2020-03-03T10:28:18.406675\"," +
				"			\"com.github.containers.toolbox\": \"true\"," +
				"			\"com.github.debarshiray.toolbox\": \"true\"," +
				"			\"com.redhat.build-host\": \"osbs-node02.phx2.fedoraproject.org\"," +
				"			\"com.redhat.component\": \"fedora-toolbox\"," +
				"			\"distribution-scope\": \"public\"," +
				"			\"license\": \"MIT\"," +
				"			\"maintainer\": \"Debarshi Ray <rishi@fedoraproject.org>\"," +
				"			\"name\": \"f30/fedora-toolbox\"," +
				"			\"release\": \"12\"," +
				"			\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"			\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"			\"vcs-ref\": \"41c23bef7659ce5e8798c5f6330e25ba1dc9d5ed\"," +
				"			\"vcs-type\": \"git\"," +
				"			\"vendor\": \"Fedora Project\"," +
				"			\"version\": \"30\"" +
				"		}," +
				"		\"Annotations\": {}," +
				"		\"ManifestType\": \"application/vnd.docker.distribution.manifest.v2+json\"," +
				"		\"User\": \"\"," +
				"		\"History\": [" +
				"			{" +
				"				\"created\": \"2020-02-24T07:48:26Z\"," +
				"				\"comment\": \"Created by Image Factory\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2020-03-03T10:32:20.503064Z\"" +
				"			}" +
				"		]," +
				"		\"NamesHistory\": []" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8",
				isToolbx:   true,
				labels: map[string]string{
					"architecture":                   "x86_64",
					"authoritative-source-url":       "registry.fedoraproject.org",
					"build-date":                     "2020-03-03T10:28:18.406675",
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.build-host":          "osbs-node02.phx2.fedoraproject.org",
					"com.redhat.component":           "fedora-toolbox",
					"distribution-scope":             "public",
					"license":                        "MIT",
					"maintainer":                     "Debarshi Ray <rishi@fedoraproject.org>",
					"name":                           "f30/fedora-toolbox",
					"release":                        "12",
					"summary":                        "Base image for creating Fedora toolbox containers",
					"usage":                          "This image is meant to be used with the toolbox command",
					"vcs-ref":                        "41c23bef7659ce5e8798c5f6330e25ba1dc9d5ed",
					"vcs-type":                       "git",
					"vendor":                         "Fedora Project",
					"version":                        "30",
				},
				namesHistory: []string{},
				repoTags: []string{
					"registry.fedoraproject.org/f30/fedora-toolbox:30",
				},
				envVars: []string{
					"DISTTAG=f30container",
					"FGC=f30",
					"container=oci",
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"NAME=fedora-toolbox",
					"VERSION=30",
				},
			},
		},
		{
			name: "podman 2.2.1, fedora-toolbox:32",
			data: "" +
				"[" +
				"  {" +
				"		\"Id\": \"6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6\"," +
				"		\"Digest\": \"sha256:dd65063ba3ee25991de8e3e13b6a44b67e4d1d4f4875153591c55f2f0cb78515\"," +
				"		\"RepoTags\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:32\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:1d1077ce12b6d45990709ab81543884838d71a3b2014d023b82fd465a3ce4e23\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:dd65063ba3ee25991de8e3e13b6a44b67e4d1d4f4875153591c55f2f0cb78515\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"\"," +
				"		\"Created\": \"2021-04-25T12:02:39.696984Z\"," +
				"		\"Config\": {" +
				"			\"Env\": [" +
				"				\"DISTTAG=f32container\"," +
				"				\"FGC=f32\"," +
				"				\"container=oci\"," +
				"				\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"				\"NAME=fedora-toolbox\"," +
				"				\"VERSION=32\"" +
				"			]," +
				"			\"Cmd\": [" +
				"				\"/bin/sh\"," +
				"				\"-c\"," +
				"				\"/bin/sh\"" +
				"			]," +
				"			\"Labels\": {" +
				"				\"architecture\": \"x86_64\"," +
				"				\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"				\"build-date\": \"2021-04-25T12:01:08.917265\"," +
				"				\"com.github.containers.toolbox\": \"true\"," +
				"				\"com.github.debarshiray.toolbox\": \"true\"," +
				"				\"com.redhat.build-host\": \"osbs-node02.iad2.fedoraproject.org\"," +
				"				\"com.redhat.component\": \"fedora-toolbox\"," +
				"				\"distribution-scope\": \"public\"," +
				"				\"license\": \"MIT\"," +
				"				\"maintainer\": \"Debarshi Ray \u003crishi@fedoraproject.org\u003e\"," +
				"				\"name\": \"fedora-toolbox\"," +
				"				\"release\": \"14\"," +
				"				\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"				\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"				\"vcs-ref\": \"c81107b8acbbd44abafd31793d71141d160e70c1\"," +
				"				\"vcs-type\": \"git\"," +
				"				\"vendor\": \"Fedora Project\"," +
				"				\"version\": \"32\"" +
				"			}" +
				"		}," +
				"		\"Version\": \"1.13.1\"," +
				"		\"Author\": \"\"," +
				"		\"Architecture\": \"amd64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 371369047," +
				"		\"VirtualSize\": 371369047," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"				\"LowerDir\": \"/home/dkricka/.local/share/containers/storage/overlay/48df8b28c24ba6e7356c15f0cbc05a48b6a23241bb93a7a6e72f742a2979111e/diff\"," +
				"				\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/218e89f69aeaeac820946361012122e2a857000328dbac9ecc8f8f2ec503cbb4/diff\"," +
				"				\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/218e89f69aeaeac820946361012122e2a857000328dbac9ecc8f8f2ec503cbb4/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"				\"sha256:48df8b28c24ba6e7356c15f0cbc05a48b6a23241bb93a7a6e72f742a2979111e\"," +
				"				\"sha256:14519312874b1edc602c8389ecc348bcb76018ac1b54d5ed4ddf18763c396eac\"" +
				"			]" +
				"		}," +
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
				"		\"Annotations\": {}," +
				"		\"ManifestType\": \"application/vnd.docker.distribution.manifest.v2+json\"," +
				"		\"User\": \"\"," +
				"		\"History\": [" +
				"			{" +
				"				\"created\": \"2021-01-06T06:48:53Z\"," +
				"				\"comment\": \"Created by Image Factory\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2021-04-25T12:02:39.696984Z\"" +
				"			}" +
				"		]," +
				"		\"NamesHistory\": []" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6",
				isToolbx:   true,
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
				namesHistory: []string{},
				repoTags: []string{
					"registry.fedoraproject.org/fedora-toolbox:32",
				},
				envVars: []string{
					"DISTTAG=f32container",
					"FGC=f32",
					"container=oci",
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"NAME=fedora-toolbox",
					"VERSION=32",
				},
			},
		},
		{
			name: "podman 3.4.7, fedora-toolbox:35",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227\"," +
				"		\"Digest\": \"sha256:b5f419099423fae80421cda692b00e301894101575806f44558e6c9c911769e6\"," +
				"		\"RepoTags\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:b5f419099423fae80421cda692b00e301894101575806f44558e6c9c911769e6\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:bdab10512284f81480235e89939babff43b2a8ca5af870dbebdeb127c6533568\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"\"," +
				"		\"Created\": \"2022-11-09T13:16:43.063221Z\"," +
				"		\"Config\": {" +
				"			\"Env\": [" +
				"				\"DISTTAG=f35container\"," +
				"				\"FGC=f35\"," +
				"				\"container=oci\"," +
				"				\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"				\"NAME=fedora-toolbox\"," +
				"				\"VERSION=35\"" +
				"			]," +
				"			\"Cmd\": [" +
				"				\"/bin/bash\"" +
				"			]," +
				"			\"Labels\": {" +
				"				\"architecture\": \"x86_64\"," +
				"				\"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"				\"build-date\": \"2022-11-09T13:15:18.185125\"," +
				"				\"com.github.containers.toolbox\": \"true\"," +
				"				\"com.github.debarshiray.toolbox\": \"true\"," +
				"				\"com.redhat.build-host\": \"osbs-node02.iad2.fedoraproject.org\"," +
				"				\"com.redhat.component\": \"fedora-toolbox\"," +
				"				\"distribution-scope\": \"public\"," +
				"				\"license\": \"MIT\"," +
				"				\"maintainer\": \"Debarshi Ray \u003crishi@fedoraproject.org\u003e\"," +
				"				\"name\": \"fedora-toolbox\"," +
				"				\"release\": \"18\"," +
				"				\"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"				\"usage\": \"This image is meant to be used with the toolbox command\"," +
				"				\"vcs-ref\": \"274bdf9053aa6e56114c9f4fca26f09a89a72ad5\"," +
				"				\"vcs-type\": \"git\"," +
				"				\"vendor\": \"Fedora Project\"," +
				"				\"version\": \"35\"" +
				"			}" +
				"		}," +
				"		\"Version\": \"1.13.1\"," +
				"		\"Author\": \"\"," +
				"		\"Architecture\": \"amd64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 495703728," +
				"		\"VirtualSize\": 495703728," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"				\"LowerDir\": \"/home/dkricka/.local/share/containers/storage/overlay/97d99750be28c24bd019af3255fc880119b5aa4d980b6ceda6b061d645913e6c/diff\"," +
				"				\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/3eef47e77a97efdd07baa9de60b0f7fba24f95ce89bde79fae2225f470dabb31/diff\"," +
				"				\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/3eef47e77a97efdd07baa9de60b0f7fba24f95ce89bde79fae2225f470dabb31/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"				\"sha256:97d99750be28c24bd019af3255fc880119b5aa4d980b6ceda6b061d645913e6c\"," +
				"				\"sha256:22c9ce93f33107b1bfd2f6ee0be474e946a314f2974c7789e767529275f7ae70\"" +
				"			]" +
				"		}," +
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
				"		\"Annotations\": {}," +
				"		\"ManifestType\": \"application/vnd.docker.distribution.manifest.v2+json\"," +
				"		\"User\": \"\"," +
				"		\"History\": [" +
				"			{" +
				"				\"created\": \"2022-09-23T06:49:12Z\"," +
				"				\"comment\": \"Created by Image Factory\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2022-11-09T13:16:43.063221Z\"" +
				"			}" +
				"		]," +
				"		\"NamesHistory\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:35\"" +
				"		]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227",
				isToolbx:   true,
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
				namesHistory: []string{
					"registry.fedoraproject.org/fedora-toolbox:35",
				},
				repoTags: []string{
					"registry.fedoraproject.org/fedora-toolbox:35",
				},
				envVars: []string{
					"DISTTAG=f35container",
					"FGC=f35",
					"container=oci",
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"NAME=fedora-toolbox",
					"VERSION=35",
				},
			},
		},
		{
			name: "podman 4.9.4, fedora-toolbox:38",
			data: "" +
				"[" +
				"	{" +
				"		 \"Id\": \"e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce\"," +
				"		 \"Digest\": \"sha256:efde4efa9f18e619f62a14849959a436fc8062b2248aca794a0f4957b1d2eeca\"," +
				"		 \"RepoTags\": [" +
				"			  \"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"		 ]," +
				"		 \"RepoDigests\": [" +
				"			  \"registry.fedoraproject.org/fedora-toolbox@sha256:8a52bb3a18508c15cd12e8770c0ce17c910c2bb96aabc2ffc4b0ceb669cef935\"," +
				"			  \"registry.fedoraproject.org/fedora-toolbox@sha256:efde4efa9f18e619f62a14849959a436fc8062b2248aca794a0f4957b1d2eeca\"" +
				"		 ]," +
				"		 \"Parent\": \"\"," +
				"		 \"Comment\": \"\"," +
				"		 \"Created\": \"2024-02-01T19:09:35.77181Z\"," +
				"		 \"Config\": {" +
				"			  \"Env\": [" +
				"				   \"DISTTAG=f38container\"," +
				"				   \"FGC=f38\"," +
				"				   \"container=oci\"," +
				"				   \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"" +
				"			  ]," +
				"			  \"Cmd\": [" +
				"				   \"/bin/bash\"" +
				"			  ]," +
				"			  \"Labels\": {" +
				"				   \"architecture\": \"x86_64\"," +
				"				   \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"				   \"build-date\": \"2024-02-01T19:07:19.783946\"," +
				"				   \"com.github.containers.toolbox\": \"true\"," +
				"				   \"com.redhat.build-host\": \"osbs-node01.iad2.fedoraproject.org\"," +
				"				   \"com.redhat.component\": \"fedora-toolbox\"," +
				"				   \"distribution-scope\": \"public\"," +
				"				   \"license\": \"MIT\"," +
				"				   \"maintainer\": \"Debarshi Ray <rishi@fedoraproject.org>\"," +
				"				   \"name\": \"fedora-toolbox\"," +
				"				   \"release\": \"20\"," +
				"				   \"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"				   \"usage\": \"This image is meant to be used with the toolbox command\"," +
				"				   \"vcs-ref\": \"c36ca58a44f947077b7c62e17d7b4fe4cd1b5797\"," +
				"				   \"vcs-type\": \"git\"," +
				"				   \"vendor\": \"Fedora Project\"," +
				"				   \"version\": \"38\"" +
				"			  }," +
				"			  \"ArgsEscaped\": true" +
				"		 }," +
				"		 \"Version\": \"1.13.1\"," +
				"		 \"Author\": \"\"," +
				"		 \"Architecture\": \"amd64\"," +
				"		 \"Os\": \"linux\"," +
				"		 \"Size\": 1748730844," +
				"		 \"VirtualSize\": 1748730844," +
				"		 \"GraphDriver\": {" +
				"			  \"Name\": \"overlay\"," +
				"			  \"Data\": {" +
				"				   \"LowerDir\": \"/home/dkricka/.local/share/containers/storage/overlay/c12fdfa4ea5b7d6efc14a116bd49660835bd9a960dd1e256962eb1d88359408e/diff\"," +
				"				   \"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/bfe5a8cd24a2f51b9e2e4743c21e270ccf2393093fb1650dbd0a8291000a98d6/diff\"," +
				"				   \"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/bfe5a8cd24a2f51b9e2e4743c21e270ccf2393093fb1650dbd0a8291000a98d6/work\"" +
				"			  }" +
				"		 }," +
				"		 \"RootFS\": {" +
				"			  \"Type\": \"layers\"," +
				"			  \"Layers\": [" +
				"				   \"sha256:c12fdfa4ea5b7d6efc14a116bd49660835bd9a960dd1e256962eb1d88359408e\"," +
				"				   \"sha256:1533eb61688ba8ff4d0e4a89e9e9d286bfa550773978717ccede5fc7d1f81ea3\"" +
				"			  ]" +
				"		 }," +
				"		 \"Labels\": {" +
				"			  \"architecture\": \"x86_64\"," +
				"			  \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"			  \"build-date\": \"2024-02-01T19:07:19.783946\"," +
				"			  \"com.github.containers.toolbox\": \"true\"," +
				"			  \"com.redhat.build-host\": \"osbs-node01.iad2.fedoraproject.org\"," +
				"			  \"com.redhat.component\": \"fedora-toolbox\"," +
				"			  \"distribution-scope\": \"public\"," +
				"			  \"license\": \"MIT\"," +
				"			  \"maintainer\": \"Debarshi Ray <rishi@fedoraproject.org>\"," +
				"			  \"name\": \"fedora-toolbox\"," +
				"			  \"release\": \"20\"," +
				"			  \"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"			  \"usage\": \"This image is meant to be used with the toolbox command\"," +
				"			  \"vcs-ref\": \"c36ca58a44f947077b7c62e17d7b4fe4cd1b5797\"," +
				"			  \"vcs-type\": \"git\"," +
				"			  \"vendor\": \"Fedora Project\"," +
				"			  \"version\": \"38\"" +
				"		 }," +
				"		 \"Annotations\": {}," +
				"		 \"ManifestType\": \"application/vnd.docker.distribution.manifest.v2+json\"," +
				"		 \"User\": \"\"," +
				"		 \"History\": [" +
				"			  {" +
				"				   \"created\": \"2024-01-25T06:49:26Z\"," +
				"				   \"comment\": \"Created by Image Factory\"" +
				"			  }," +
				"			  {" +
				"				   \"created\": \"2024-02-01T19:09:35.77181Z\"" +
				"			  }" +
				"		 ]," +
				"		 \"NamesHistory\": [" +
				"			  \"registry.fedoraproject.org/fedora-toolbox:38\"" +
				"		 ]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce",
				isToolbx:   true,
				labels: map[string]string{
					"architecture":                  "x86_64",
					"authoritative-source-url":      "registry.fedoraproject.org",
					"build-date":                    "2024-02-01T19:07:19.783946",
					"com.github.containers.toolbox": "true",
					"com.redhat.build-host":         "osbs-node01.iad2.fedoraproject.org",
					"com.redhat.component":          "fedora-toolbox",
					"distribution-scope":            "public",
					"license":                       "MIT",
					"maintainer":                    "Debarshi Ray <rishi@fedoraproject.org>",
					"name":                          "fedora-toolbox",
					"release":                       "20",
					"summary":                       "Base image for creating Fedora toolbox containers",
					"usage":                         "This image is meant to be used with the toolbox command",
					"vcs-ref":                       "c36ca58a44f947077b7c62e17d7b4fe4cd1b5797",
					"vcs-type":                      "git",
					"vendor":                        "Fedora Project",
					"version":                       "38",
				},
				namesHistory: []string{
					"registry.fedoraproject.org/fedora-toolbox:38",
				},
				repoTags: []string{
					"registry.fedoraproject.org/fedora-toolbox:38",
				},
				envVars: []string{
					"DISTTAG=f38container",
					"FGC=f38",
					"container=oci",
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
				},
			},
		},
		{
			name: "podman 5.4.2, fedora-toolbox:40",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628\"," +
				"		\"Digest\": \"sha256:3d8c46e803fd184e9318ee90d4e6f1ad0b591a1c49f4de683d4f8748c0b95c30\"," +
				"		\"RepoTags\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:3d8c46e803fd184e9318ee90d4e6f1ad0b591a1c49f4de683d4f8748c0b95c30\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:c3cc9836c9e55475f85496f9369e925164f1b7cd832b55e22c13b74840576c31\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"\"," +
				"		\"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"		\"Config\": {" +
				"			\"Env\": [" +
				"					\"container=oci\"" +
				"			]," +
				"			\"Cmd\": [" +
				"					\"/bin/bash\"" +
				"			]," +
				"			\"WorkingDir\": \"/\"," +
				"			\"Labels\": {" +
				"					\"com.github.containers.toolbox\": \"true\"," +
				"					\"io.buildah.version\": \"1.39.2\"," +
				"					\"license\": \"MIT\"," +
				"					\"name\": \"fedora-toolbox\"," +
				"					\"org.opencontainers.image.license\": \"MIT\"," +
				"					\"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"					\"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"					\"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"					\"org.opencontainers.image.version\": \"40\"," +
				"					\"vendor\": \"Fedora Project\"," +
				"					\"version\": \"40\"" +
				"			}" +
				"		}," +
				"		\"Version\": \"\"," +
				"		\"Author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"		\"Architecture\": \"amd64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 2204748042," +
				"		\"VirtualSize\": 2204748042," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"					\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/ff4e046e00ff0e39fbfef8473831eb608e0f0c12ab19118059d4162b5b8b25cc/diff\"," +
				"					\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/ff4e046e00ff0e39fbfef8473831eb608e0f0c12ab19118059d4162b5b8b25cc/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"					\"sha256:ff4e046e00ff0e39fbfef8473831eb608e0f0c12ab19118059d4162b5b8b25cc\"" +
				"			]" +
				"		}," +
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
				"		\"Annotations\": {" +
				"			\"org.opencontainers.image.base.digest\": \"\"," +
				"			\"org.opencontainers.image.base.name\": \"\"" +
				"		}," +
				"		\"ManifestType\": \"application/vnd.oci.image.manifest.v1+json\"," +
				"		\"User\": \"\"," +
				"		\"History\": [" +
				"			{" +
				"					\"created\": \"2025-05-13T07:49:06.775757207Z\"," +
				"					\"created_by\": \"KIWI 10.2.19\"," +
				"					\"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"" +
				"			}" +
				"		]," +
				"		\"NamesHistory\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:40\"" +
				"		]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "5c5b2e637806fccb644effa4affc4a5d08dc7e6140586ecb8c601c8739e12628",
				isToolbx:   true,
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
				envVars: []string{
					"container=oci",
				},
			},
		},
		{
			name: "podman 5.6.1, fedora-toolbox:41",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"b65ee158f921088db52e8b98a6f3267de75324c8a6a04afc9ff095338c40e59b\"," +
				"		\"Digest\": \"sha256:8599f0c0d421c0dc01c4b7fb1c07b2780c0ab1931d0f22dd7b6da3b93ff6b77b\"," +
				"		\"RepoTags\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:41\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:0a51adc6bab55d49ff00da8aaad81ca1f02315511ed23d55ee5bbbe1a976a663\"," +
				"			\"registry.fedoraproject.org/fedora-toolbox@sha256:8599f0c0d421c0dc01c4b7fb1c07b2780c0ab1931d0f22dd7b6da3b93ff6b77b\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"\"," +
				"		\"Created\": \"2025-10-06T05:48:05.571613915Z\"," +
				"		\"Config\": {" +
				"			\"Env\": [" +
				"					\"container=oci\"" +
				"			]," +
				"			\"Cmd\": [" +
				"					\"/bin/bash\"" +
				"			]," +
				"			\"WorkingDir\": \"/\"," +
				"			\"Labels\": {" +
				"					\"com.github.containers.toolbox\": \"true\"," +
				"					\"io.buildah.version\": \"1.41.4\"," +
				"					\"license\": \"MIT\"," +
				"					\"name\": \"fedora-toolbox\"," +
				"					\"org.opencontainers.image.license\": \"MIT\"," +
				"					\"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"					\"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"					\"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"					\"org.opencontainers.image.version\": \"41\"," +
				"					\"vendor\": \"Fedora Project\"," +
				"					\"version\": \"41\"" +
				"			}" +
				"		}," +
				"		\"Version\": \"\"," +
				"		\"Author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"		\"Architecture\": \"amd64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 2308518290," +
				"		\"VirtualSize\": 2308518290," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"					\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/5d97edd0878eb05ef62bfdccc49ad38009ba2eb15bbcce5ff9f4bacc44048427/diff\"," +
				"					\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/5d97edd0878eb05ef62bfdccc49ad38009ba2eb15bbcce5ff9f4bacc44048427/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"					\"sha256:5d97edd0878eb05ef62bfdccc49ad38009ba2eb15bbcce5ff9f4bacc44048427\"" +
				"			]" +
				"		}," +
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
				"		\"Annotations\": {" +
				"			\"org.opencontainers.image.base.digest\": \"\"," +
				"			\"org.opencontainers.image.base.name\": \"\"," +
				"			\"org.opencontainers.image.created\": \"2025-10-06T05:48:05.571613915Z\"" +
				"		}," +
				"		\"ManifestType\": \"application/vnd.oci.image.manifest.v1+json\"," +
				"		\"User\": \"\"," +
				"		\"History\": [" +
				"			{" +
				"					\"created\": \"2025-10-06T05:48:08.998261156Z\"," +
				"					\"created_by\": \"KIWI 10.2.33\"," +
				"					\"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"" +
				"			}" +
				"		]," +
				"		\"NamesHistory\": [" +
				"			\"registry.fedoraproject.org/fedora-toolbox:41\"" +
				"		]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "b65ee158f921088db52e8b98a6f3267de75324c8a6a04afc9ff095338c40e59b",
				isToolbx:   true,
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
				namesHistory: []string{
					"registry.fedoraproject.org/fedora-toolbox:41",
				},
				repoTags: []string{
					"registry.fedoraproject.org/fedora-toolbox:41",
				},
				envVars: []string{
					"container=oci",
				},
			},
		},
		{
			name: "podman 5.6.2, fedora-toolbox:42",
			data: "" +
				"[" +
				"	{" +
				"		 \"Id\": \"341ada9b9af076546132ae6fe6e328eb777046581e23b1c232bfc71d844f8598\"," +
				"		 \"Digest\": \"sha256:83540b1b86984bc56e85c0607ec0cc2469c45f6716259bc98668cacd9bbce48a\"," +
				"		 \"RepoTags\": [" +
				"			  \"registry.fedoraproject.org/fedora-toolbox:42\"" +
				"		 ]," +
				"		 \"RepoDigests\": [" +
				"			  \"registry.fedoraproject.org/fedora-toolbox@sha256:2fd640a4b02193f88972845e73e22f3943f9b69c69f13cd17c50cec098bd0715\"," +
				"			  \"registry.fedoraproject.org/fedora-toolbox@sha256:83540b1b86984bc56e85c0607ec0cc2469c45f6716259bc98668cacd9bbce48a\"" +
				"		 ]," +
				"		 \"Parent\": \"\"," +
				"		 \"Comment\": \"\"," +
				"		 \"Created\": \"2025-10-06T06:49:50.766447647Z\"," +
				"		 \"Config\": {" +
				"			  \"Env\": [" +
				"				   \"PATH=/usr/local/bin:/usr/bin\"," +
				"				   \"container=oci\"" +
				"			  ]," +
				"			  \"Cmd\": [" +
				"				   \"/bin/bash\"" +
				"			  ]," +
				"			  \"WorkingDir\": \"/\"," +
				"			  \"Labels\": {" +
				"				   \"com.github.containers.toolbox\": \"true\"," +
				"				   \"io.buildah.version\": \"1.41.5\"," +
				"				   \"license\": \"MIT\"," +
				"				   \"name\": \"fedora-toolbox\"," +
				"				   \"org.opencontainers.image.license\": \"MIT\"," +
				"				   \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"				   \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"				   \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"				   \"org.opencontainers.image.version\": \"42\"," +
				"				   \"vendor\": \"Fedora Project\"," +
				"				   \"version\": \"42\"" +
				"			  }" +
				"		 }," +
				"		 \"Version\": \"\"," +
				"		 \"Author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"		 \"Architecture\": \"amd64\"," +
				"		 \"Os\": \"linux\"," +
				"		 \"Size\": 2136421808," +
				"		 \"VirtualSize\": 2136421808," +
				"		 \"GraphDriver\": {" +
				"			  \"Name\": \"overlay\"," +
				"			  \"Data\": {" +
				"				   \"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/a9f697caafe52ff941ae3024b709517804a8a233fc4998d01867cc37425c0381/diff\"," +
				"				   \"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/a9f697caafe52ff941ae3024b709517804a8a233fc4998d01867cc37425c0381/work\"" +
				"			  }" +
				"		 }," +
				"		 \"RootFS\": {" +
				"			  \"Type\": \"layers\"," +
				"			  \"Layers\": [" +
				"				   \"sha256:a9f697caafe52ff941ae3024b709517804a8a233fc4998d01867cc37425c0381\"" +
				"			  ]" +
				"		 }," +
				"		 \"Labels\": {" +
				"			  \"com.github.containers.toolbox\": \"true\"," +
				"			  \"io.buildah.version\": \"1.41.5\"," +
				"			  \"license\": \"MIT\"," +
				"			  \"name\": \"fedora-toolbox\"," +
				"			  \"org.opencontainers.image.license\": \"MIT\"," +
				"			  \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"			  \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"			  \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"			  \"org.opencontainers.image.version\": \"42\"," +
				"			  \"vendor\": \"Fedora Project\"," +
				"			  \"version\": \"42\"" +
				"		 }," +
				"		 \"Annotations\": {" +
				"			  \"org.opencontainers.image.base.digest\": \"\"," +
				"			  \"org.opencontainers.image.base.name\": \"\"," +
				"			  \"org.opencontainers.image.created\": \"2025-10-06T06:49:50.766447647Z\"" +
				"		 }," +
				"		 \"ManifestType\": \"application/vnd.oci.image.manifest.v1+json\"," +
				"		 \"User\": \"\"," +
				"		 \"History\": [" +
				"			  {" +
				"				   \"created\": \"2025-10-06T06:49:54.42143611Z\"," +
				"				   \"created_by\": \"KIWI 10.2.33\"," +
				"				   \"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"" +
				"			  }" +
				"		 ]," +
				"		 \"NamesHistory\": [" +
				"			  \"registry.fedoraproject.org/fedora-toolbox:42\"" +
				"		 ]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "341ada9b9af076546132ae6fe6e328eb777046581e23b1c232bfc71d844f8598",
				isToolbx:   true,
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
				namesHistory: []string{
					"registry.fedoraproject.org/fedora-toolbox:42",
				},
				repoTags: []string{
					"registry.fedoraproject.org/fedora-toolbox:42",
				},
				envVars: []string{
					"PATH=/usr/local/bin:/usr/bin",
					"container=oci",
				},
			},
		},
		{
			name: "podman 5.6.2, fedora-toolbox:42-aarch64",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"2718f04f884164eb696b9ea011e29bca3222e901e387db80915428f2ed5ca5d7\"," +
				"		\"Digest\": \"sha256:8d282ca2e63be6a19e0f82bb0dafe3aa58ea45335be5876d596bee0dc788e679\"," +
				"		\"RepoTags\": [" +
				"			\"quay.io/fedora/fedora-toolbox:42-aarch64\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"quay.io/fedora/fedora-toolbox@sha256:8d282ca2e63be6a19e0f82bb0dafe3aa58ea45335be5876d596bee0dc788e679\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"\"," +
				"		\"Created\": \"2025-09-22T06:57:30.229081342Z\"," +
				"		\"Config\": {" +
				"			\"Env\": [" +
				"				\"PATH=/usr/local/bin:/usr/bin\"," +
				"				\"container=oci\"" +
				"			]," +
				"			\"Cmd\": [" +
				"				\"/bin/bash\"" +
				"			]," +
				"			\"WorkingDir\": \"/\"," +
				"			\"Labels\": {" +
				"				\"com.github.containers.toolbox\": \"true\"," +
				"				\"io.buildah.version\": \"1.41.4\"," +
				"				\"license\": \"MIT\"," +
				"				\"name\": \"fedora-toolbox\"," +
				"				\"org.opencontainers.image.license\": \"MIT\"," +
				"				\"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"				\"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"				\"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"				\"org.opencontainers.image.version\": \"42\"," +
				"				\"vendor\": \"Fedora Project\"," +
				"				\"version\": \"42\"" +
				"			}" +
				"		}," +
				"		\"Version\": \"\"," +
				"		\"Author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"," +
				"		\"Architecture\": \"arm64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 2103319610," +
				"		\"VirtualSize\": 2103319610," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"				\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/159efe71ff4d6e5cb4fc62edce161380e39eb0d07acef359489bda062ccd6695/diff\"," +
				"				\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/159efe71ff4d6e5cb4fc62edce161380e39eb0d07acef359489bda062ccd6695/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"				\"sha256:159efe71ff4d6e5cb4fc62edce161380e39eb0d07acef359489bda062ccd6695\"" +
				"			]" +
				"		}," +
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
				"		\"Annotations\": {" +
				"			\"org.opencontainers.image.base.digest\": \"\"," +
				"			\"org.opencontainers.image.base.name\": \"\"," +
				"			\"org.opencontainers.image.created\": \"2025-09-22T06:57:30.229081342Z\"" +
				"		}," +
				"		\"ManifestType\": \"application/vnd.oci.image.manifest.v1+json\"," +
				"		\"User\": \"\"," +
				"		\"History\": [" +
				"			{" +
				"				\"created\": \"2025-09-22T06:57:33.61369795Z\"," +
				"				\"created_by\": \"KIWI 10.2.33\"," +
				"				\"author\": \"Fedora Project Contributors <devel@lists.fedoraproject.org>\"" +
				"			}" +
				"		]," +
				"		\"NamesHistory\": [" +
				"			\"quay.io/fedora/fedora-toolbox:42-aarch64\"" +
				"		]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: nil,
				id:         "2718f04f884164eb696b9ea011e29bca3222e901e387db80915428f2ed5ca5d7",
				isToolbx:   true,
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
				namesHistory: []string{
					"quay.io/fedora/fedora-toolbox:42-aarch64",
				},
				repoTags: []string{
					"quay.io/fedora/fedora-toolbox:42-aarch64",
				},
				envVars: []string{
					"PATH=/usr/local/bin:/usr/bin",
					"container=oci",
				},
			},
		},
		{
			name: "podman 5.6.2, docker.io-postgres:18",
			data: "" +
				"[" +
				"	{" +
				"		\"Id\": \"194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd\"," +
				"		\"Digest\": \"sha256:28f01a051c819681a816dca282088111ade7c44f834dd83cfd044f0548d38c19\"," +
				"		\"RepoTags\": [" +
				"			\"docker.io/library/postgres:latest\"" +
				"		]," +
				"		\"RepoDigests\": [" +
				"			\"docker.io/library/postgres@sha256:073e7c8b84e2197f94c8083634640ab37105effe1bc853ca4d5fbece3219b0e8\"," +
				"			\"docker.io/library/postgres@sha256:28f01a051c819681a816dca282088111ade7c44f834dd83cfd044f0548d38c19\"" +
				"		]," +
				"		\"Parent\": \"\"," +
				"		\"Comment\": \"debuerreotype 0.16\"," +
				"		\"Created\": \"2025-09-25T18:22:35Z\"," +
				"		\"Config\": {" +
				"			\"ExposedPorts\": {" +
				"				\"5432/tcp\": {}" +
				"			}," +
				"			\"Env\": [" +
				"				\"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/postgresql/18/bin\"," +
				"				\"GOSU_VERSION=1.19\"," +
				"				\"LANG=en_US.utf8\"," +
				"				\"PG_MAJOR=18\"," +
				"				\"PG_VERSION=18.0-1.pgdg13+3\"," +
				"				\"PGDATA=/var/lib/postgresql/18/docker\"" +
				"			]," +
				"			\"Entrypoint\": [" +
				"				\"docker-entrypoint.sh\"" +
				"			]," +
				"			\"Cmd\": [" +
				"				\"postgres\"" +
				"			]," +
				"			\"Volumes\": {" +
				"				\"/var/lib/postgresql\": {}" +
				"			}," +
				"			\"StopSignal\": \"SIGINT\"" +
				"		}," +
				"		\"Version\": \"\"," +
				"		\"Author\": \"\"," +
				"		\"Architecture\": \"amd64\"," +
				"		\"Os\": \"linux\"," +
				"		\"Size\": 463150524," +
				"		\"VirtualSize\": 463150524," +
				"		\"GraphDriver\": {" +
				"			\"Name\": \"overlay\"," +
				"			\"Data\": {" +
				"				\"LowerDir\": \"/home/dkricka/.local/share/containers/storage/overlay/a1f10cfd8a7afaceb6e3a5e62efbd76473e2dfe63ef6d03710af50a30d1396c0/diff:/home/dkricka/.local/share/containers/storage/overlay/d7ea160218de0bcdc6429adf9967617b3cf6c9d278a742a6a129727081c48f74/diff:/home/dkricka/.local/share/containers/storage/overlay/5893d38a1da70682f5d60c96907dee1e8a7c12641ca58a4df59084f205206245/diff:/home/dkricka/.local/share/containers/storage/overlay/f1ce038dbbcc1ff5e5c542e074c73c7fd23d54b65666bfbff70b9d2ea20c73fa/diff:/home/dkricka/.local/share/containers/storage/overlay/ffc0cf8058edb65e2b6008914350b9b288a3c2184cc9f3156b548b10766f9a7c/diff:/home/dkricka/.local/share/containers/storage/overlay/1315540e4cc7ea57cb4176923c151764a0d3eaa84d1988aa91b6a1cbb182a55b/diff:/home/dkricka/.local/share/containers/storage/overlay/2401f98f3fa81f6f5ef7daf085d253b648ac8b344c8126241e831dcd92e71971/diff:/home/dkricka/.local/share/containers/storage/overlay/2327e93659ba0e088112cd7d85868e3a880c5cbb65ee170d226eaad1ff3f85c0/diff:/home/dkricka/.local/share/containers/storage/overlay/93a6331fd6d03678ed591ae0de4f46a7f2b7574210a736246f87c79337df3f6f/diff:/home/dkricka/.local/share/containers/storage/overlay/5a7249aa355de7c663047a9c91562385d69bae1f52dc0f8329bc8ecbc0df47af/diff:/home/dkricka/.local/share/containers/storage/overlay/99f8e117d4b6f889773125425ef8d75a8c5fabb3ea3564526e9cac035d60fbfd/diff:/home/dkricka/.local/share/containers/storage/overlay/7b38d5b07f61b7a96ea0f0de8e603e0949ed96125867696a58d174e91fbfa3e9/diff:/home/dkricka/.local/share/containers/storage/overlay/1d46119d249f7719e1820e24a311aa7c453f166f714969cffe89504678eaa447/diff\"," +
				"				\"UpperDir\": \"/home/dkricka/.local/share/containers/storage/overlay/7e452011f7ac750c95ecc68aebc97f6f39c08c5d001712d98d535b2e7b7d7199/diff\"," +
				"				\"WorkDir\": \"/home/dkricka/.local/share/containers/storage/overlay/7e452011f7ac750c95ecc68aebc97f6f39c08c5d001712d98d535b2e7b7d7199/work\"" +
				"			}" +
				"		}," +
				"		\"RootFS\": {" +
				"			\"Type\": \"layers\"," +
				"			\"Layers\": [" +
				"				\"sha256:1d46119d249f7719e1820e24a311aa7c453f166f714969cffe89504678eaa447\"," +
				"				\"sha256:a9b74bbbba249d0a370c711687340045a284abfa6e24f89c9d3c5a9be2de1aff\"," +
				"				\"sha256:af76e12aa831a3a89bd606289fdb9aa7b0f0e468acbe1216eee438865b4955eb\"," +
				"				\"sha256:01e8a18dc9b2e9d81cb2ec31a00d4e4a210382fde795d988367e462d7965d7c4\"," +
				"				\"sha256:60587f31ccd4bc683f744f53b44b0ea1321cd59e15fa6c09c7cf12dfdfcac8f6\"," +
				"				\"sha256:dad63314a339ae2200d0a62a6aa54f17f5dac2769bc82d482870e8fd50c99334\"," +
				"				\"sha256:04d52f0a5b32b0f627bbd4427a0374f0a8d2d409dbbfda0099d89b87c774df36\"," +
				"				\"sha256:ca7f2e5210d61f5c629aef27c6707a26ce04f1688b4d3e2fad36a50b2de346ff\"," +
				"				\"sha256:aae2dbb17823d9d8e51357c7b0939e85b574648cab0af3b13f09658774d8d160\"," +
				"				\"sha256:c9655bf591eeb52c441721c3367e47e074457612fd363530c3a4d7d00d8ba27e\"," +
				"				\"sha256:62f0a60f9afcb19c8cc0c8a40db2027e76556e87a406dcc047483fd132ed1823\"," +
				"				\"sha256:30d59a0e59badee5b8d4b170974204a86803c2993605d10ea41794eac8aab770\"," +
				"				\"sha256:e4a2b82444b16c92d347ce2d033be19e315292a04b8e341555f618c0ba5a993a\"," +
				"				\"sha256:88763bd2f91627360cd5de61bc91e3b096844e1646015e73f2b91d3486cd8e86\"" +
				"			]" +
				"		}," +
				"		\"Labels\": null," +
				"		\"Annotations\": {" +
				"			\"com.docker.official-images.bashbrew.arch\": \"amd64\"," +
				"			\"org.opencontainers.image.base.digest\": \"sha256:c99c73388e005d98f2f131b15fa9389f2a8eec2888a35dc30455e5936467803b\"," +
				"			\"org.opencontainers.image.base.name\": \"debian:trixie-slim\"," +
				"			\"org.opencontainers.image.created\": \"2025-09-25T18:22:35Z\"," +
				"			\"org.opencontainers.image.revision\": \"22ca5c8d8e4b37bece4d38dbce1a060583b5308a\"," +
				"			\"org.opencontainers.image.source\": \"https://github.com/docker-library/postgres.git#22ca5c8d8e4b37bece4d38dbce1a060583b5308a:18/trixie\"," +
				"			\"org.opencontainers.image.url\": \"https://hub.docker.com/_/postgres\"," +
				"			\"org.opencontainers.image.version\": \"18.0\"" +
				"		}," +
				"		\"ManifestType\": \"application/vnd.oci.image.manifest.v1+json\"," +
				"		\"User\": \"\"," +
				"		\"History\": [" +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"# debian.sh --arch 'amd64' out/ 'trixie' '@1759104000'\"," +
				"				\"comment\": \"debuerreotype 0.16\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -eux; \\tgroupadd -r postgres --gid=999; \\tuseradd -r -g postgres --uid=999 --home-dir=/var/lib/postgresql --shell=/bin/bash postgres; \\tinstall --verbose --directory --owner postgres --group postgres --mode 1777 /var/lib/postgresql # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -ex; \\tapt-get update; \\tapt-get install -y --no-install-recommends \\t\\tgnupg \\t\\tless \\t; \\trm -rf /var/lib/apt/lists/* # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"ENV GOSU_VERSION=1.19\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -eux; \\tsavedAptMark=\\\"$(apt-mark showmanual)\\\"; \\tapt-get update; \\tapt-get install -y --no-install-recommends ca-certificates wget; \\trm -rf /var/lib/apt/lists/*; \\tdpkgArch=\\\"$(dpkg --print-architecture | awk -F- '{ print $NF }')\\\"; \\twget -O /usr/local/bin/gosu \\\"https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch\\\"; \\twget -O /usr/local/bin/gosu.asc \\\"https://github.com/tianon/gosu/releases/download/$GOSU_VERSION/gosu-$dpkgArch.asc\\\"; \\texport GNUPGHOME=\\\"$(mktemp -d)\\\"; \\tgpg --batch --keyserver hkps://keys.openpgp.org --recv-keys B42F6819007F00F88E364FD4036A9C25BF357DD4; \\tgpg --batch --verify /usr/local/bin/gosu.asc /usr/local/bin/gosu; \\tgpgconf --kill all; \\trm -rf \\\"$GNUPGHOME\\\" /usr/local/bin/gosu.asc; \\tapt-mark auto '.*' > /dev/null; \\t[ -z \\\"$savedAptMark\\\" ] || apt-mark manual $savedAptMark > /dev/null; \\tapt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \\tchmod +x /usr/local/bin/gosu; \\tgosu --version; \\tgosu nobody true # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -eux; \\tif [ -f /etc/dpkg/dpkg.cfg.d/docker ]; then \\t\\tgrep -q '/usr/share/locale' /etc/dpkg/dpkg.cfg.d/docker; \\t\\tsed -ri '/\\/usr\\/share\\/locale/d' /etc/dpkg/dpkg.cfg.d/docker; \\t\\t! grep -q '/usr/share/locale' /etc/dpkg/dpkg.cfg.d/docker; \\tfi; \\tapt-get update; apt-get install -y --no-install-recommends locales; rm -rf /var/lib/apt/lists/*; \\techo 'en_US.UTF-8 UTF-8' >> /etc/locale.gen; \\tlocale-gen; \\tlocale -a | grep 'en_US.utf8' # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"ENV LANG=en_US.utf8\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -eux; \\tapt-get update; \\tapt-get install -y --no-install-recommends \\t\\tlibnss-wrapper \\t\\txz-utils \\t\\tzstd \\t; \\trm -rf /var/lib/apt/lists/* # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c mkdir /docker-entrypoint-initdb.d # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -ex; \\tkey='B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8'; \\texport GNUPGHOME=\\\"$(mktemp -d)\\\"; \\tmkdir -p /usr/local/share/keyrings/; \\tgpg --batch --keyserver keyserver.ubuntu.com --recv-keys \\\"$key\\\"; \\tgpg --batch --export --armor \\\"$key\\\" > /usr/local/share/keyrings/postgres.gpg.asc; \\tgpgconf --kill all; \\trm -rf \\\"$GNUPGHOME\\\" # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"ENV PG_MAJOR=18\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/postgresql/18/bin\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"ENV PG_VERSION=18.0-1.pgdg13+3\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -ex; \\t\\texport PYTHONDONTWRITEBYTECODE=1; \\t\\tdpkgArch=\\\"$(dpkg --print-architecture)\\\"; \\taptRepo=\\\"[ signed-by=/usr/local/share/keyrings/postgres.gpg.asc ] http://apt.postgresql.org/pub/repos/apt trixie-pgdg main $PG_MAJOR\\\"; \\tcase \\\"$dpkgArch\\\" in \\t\\tamd64 | arm64 | ppc64el) \\t\\t\\techo \\\"deb $aptRepo\\\" > /etc/apt/sources.list.d/pgdg.list; \\t\\t\\tapt-get update; \\t\\t\\t;; \\t\\t*) \\t\\t\\techo \\\"deb-src $aptRepo\\\" > /etc/apt/sources.list.d/pgdg.list; \\t\\t\\t\\t\\t\\tsavedAptMark=\\\"$(apt-mark showmanual)\\\"; \\t\\t\\t\\t\\t\\ttempDir=\\\"$(mktemp -d)\\\"; \\t\\t\\tcd \\\"$tempDir\\\"; \\t\\t\\t\\t\\t\\tapt-get update; \\t\\t\\tapt-get install -y --no-install-recommends dpkg-dev; \\t\\t\\techo \\\"deb [ trusted=yes ] file://$tempDir ./\\\" > /etc/apt/sources.list.d/temp.list; \\t\\t\\t_update_repo() { \\t\\t\\t\\tdpkg-scanpackages . > Packages; \\t\\t\\t\\tapt-get -o Acquire::GzipIndexes=false update; \\t\\t\\t}; \\t\\t\\t_update_repo; \\t\\t\\t\\t\\t\\tnproc=\\\"$(nproc)\\\"; \\t\\t\\texport DEB_BUILD_OPTIONS=\\\"nocheck parallel=$nproc\\\"; \\t\\t\\tapt-get build-dep -y postgresql-common-dev; \\t\\t\\tapt-get source --compile postgresql-common-dev; \\t\\t\\t_update_repo; \\t\\t\\tapt-get build-dep -y \\\"postgresql-$PG_MAJOR=$PG_VERSION\\\"; \\t\\t\\tapt-get source --compile \\\"postgresql-$PG_MAJOR=$PG_VERSION\\\"; \\t\\t\\t\\t\\t\\t\\t\\t\\tapt-mark showmanual | xargs apt-mark auto > /dev/null; \\t\\t\\tapt-mark manual $savedAptMark; \\t\\t\\t\\t\\t\\tls -lAFh; \\t\\t\\t_update_repo; \\t\\t\\tgrep '^Package: ' Packages; \\t\\t\\tcd /; \\t\\t\\t;; \\tesac; \\t\\tapt-get install -y --no-install-recommends postgresql-common; \\tsed -ri 's/#(create_main_cluster) .*$/\\\\1 = false/' /etc/postgresql-common/createcluster.conf; \\tapt-get install -y --no-install-recommends \\t\\t\\\"postgresql-$PG_MAJOR=$PG_VERSION\\\" \\t; \\tif apt-get install -s \\\"postgresql-$PG_MAJOR-jit\\\" > /dev/null 2>&1; then \\t\\tapt-get install -y --no-install-recommends \\\"postgresql-$PG_MAJOR-jit=$PG_VERSION\\\"; \\tfi; \\t\\trm -rf /var/lib/apt/lists/*; \\t\\tif [ -n \\\"$tempDir\\\" ]; then \\t\\tapt-get purge -y --auto-remove; \\t\\trm -rf \\\"$tempDir\\\" /etc/apt/sources.list.d/temp.list; \\tfi; \\t\\tfind /usr -name '*.pyc' -type f -exec bash -c 'for pyc; do dpkg -S \\\"$pyc\\\" &> /dev/null || rm -vf \\\"$pyc\\\"; done' -- '{}' +; \\t\\tpostgres --version # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c set -eux; \\tdpkg-divert --add --rename --divert \\\"/usr/share/postgresql/postgresql.conf.sample.dpkg\\\" \\\"/usr/share/postgresql/$PG_MAJOR/postgresql.conf.sample\\\"; \\tcp -v /usr/share/postgresql/postgresql.conf.sample.dpkg /usr/share/postgresql/postgresql.conf.sample; \\tln -sv ../postgresql.conf.sample \\\"/usr/share/postgresql/$PG_MAJOR/\\\"; \\tsed -ri \\\"s!^#?(listen_addresses)\\\\\\\\s*=\\\\\\\\s*\\\\\\\\S+.*!\\\\\\\\1 = '*'!\\\" /usr/share/postgresql/postgresql.conf.sample; \\tgrep -F \\\"listen_addresses = '*'\\\" /usr/share/postgresql/postgresql.conf.sample # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c install --verbose --directory --owner postgres --group postgres --mode 3777 /var/run/postgresql # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"ENV PGDATA=/var/lib/postgresql/18/docker\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c ln -svT . /var/lib/postgresql/data # https://github.com/docker-library/postgres/pull/1259#issuecomment-2215477494 # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"VOLUME [/var/lib/postgresql]\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"COPY docker-entrypoint.sh docker-ensure-initdb.sh /usr/local/bin/ # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"RUN /bin/sh -c ln -sT docker-ensure-initdb.sh /usr/local/bin/docker-enforce-initdb.sh # buildkit\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"ENTRYPOINT [\\\"docker-entrypoint.sh\\\"]\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"STOPSIGNAL SIGINT\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"EXPOSE map[5432/tcp:{}]\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}," +
				"			{" +
				"				\"created\": \"2025-09-25T18:22:35Z\"," +
				"				\"created_by\": \"CMD [\\\"postgres\\\"]\"," +
				"				\"comment\": \"buildkit.dockerfile.v0\"," +
				"				\"empty_layer\": true" +
				"			}" +
				"		]," +
				"		\"NamesHistory\": [" +
				"			\"docker.io/library/postgres:latest\"" +
				"		]" +
				"	}" +
				"]",
			expect: expect{
				entrypoint: []string{
					"docker-entrypoint.sh",
				},
				id:       "194f5f2a900a5775ecff2129be107adc7c1ce98b89ac00ca0bed141310b7e6cd",
				isToolbx: false,
				labels:   nil,
				namesHistory: []string{
					"docker.io/library/postgres:latest",
				},
				repoTags: []string{
					"docker.io/library/postgres:latest",
				},
				envVars: []string{
					"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/lib/postgresql/18/bin",
					"GOSU_VERSION=1.19",
					"LANG=en_US.utf8",
					"PG_MAJOR=18",
					"PG_VERSION=18.0-1.pgdg13+3",
					"PGDATA=/var/lib/postgresql/18/docker",
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
			assert.Len(t, images, 1)

			image := images[0]
			assert.Equal(t, tc.expect.id, image.ID())
			assert.Equal(t, tc.expect.isToolbx, image.IsToolbx())
			assert.Equal(t, tc.expect.labels, image.Labels())
			assert.Equal(t, tc.expect.namesHistory, image.Names())
			assert.Equal(t, tc.expect.repoTags, image.RepoTags())
			assert.Equal(t, tc.expect.entrypoint, image.Entrypoint())
			assert.Equal(t, tc.expect.envVars, image.EnvVars())
		})
	}
}

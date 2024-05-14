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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainerInspect(t *testing.T) {
	type expect struct {
		entryPoint    string
		entryPointPID int
		id            string
		image         string
		isToolbx      bool
		labels        map[string]string
		mounts        []string
		name          string
		status        string
	}

	testCases := []struct {
		name   string
		data   string
		expect expect
	}{
		{
			name: "podman 1.1.2, toolbx 0.0.9, configured",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"sleep\"," +
				"        \"+Inf\"" +
				"      ]," +
				"      \"Image\": \"localhost/fedora-toolbox-user:29\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"f29/fedora-toolbox\"," +
				"        \"version\": \"29\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"b23f7c69ddec697f803b8acc40e85d212198c1baea9ffe193e7b3e0d2a020a39\"," +
				"    \"Image\": \"c7fab9e10750847a20b6664f485d9a1430b326eaf2799b932d15eebea36d6f5f\"," +
				"    \"ImageName\": \"localhost/fedora-toolbox-user:29\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"options\": [" +
				"          \"rbind\"," +
				"          \"ro\"," +
				"          \"rprivate\"" +
				"        ]," +
				"        \"source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"destination\": \"/run/user/1000\"," +
				"        \"options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"," +
				"          \"rprivate\"," +
				"          \"rw\"" +
				"        ]," +
				"        \"source\": \"/run/user/1000\"," +
				"        \"type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"destination\": \"/usr/bin/toolbox\"," +
				"        \"options\": [" +
				"          \"rbind\"," +
				"          \"ro\"," +
				"          \"rprivate\"" +
				"        ]," +
				"        \"source\": \"/some/prefix/binary\"," +
				"        \"type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-user-29\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Pid\": 0," +
				"      \"Status\": \"configured\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "sleep",
				entryPointPID: 0,
				id:            "b23f7c69ddec697f803b8acc40e85d212198c1baea9ffe193e7b3e0d2a020a39",
				image:         "localhost/fedora-toolbox-user:29",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "f29/fedora-toolbox",
					"version":                        "29",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
				},
				name:   "fedora-toolbox-user-29",
				status: "configured",
			},
		},
		{
			name: "podman 1.1.2, toolbx 0.0.9, exited",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"sleep\"," +
				"        \"+Inf\"" +
				"      ]," +
				"      \"Image\": \"localhost/fedora-toolbox-user:29\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"f29/fedora-toolbox\"," +
				"        \"version\": \"29\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"b23f7c69ddec697f803b8acc40e85d212198c1baea9ffe193e7b3e0d2a020a39\"," +
				"    \"Image\": \"c7fab9e10750847a20b6664f485d9a1430b326eaf2799b932d15eebea36d6f5f\"," +
				"    \"ImageName\": \"localhost/fedora-toolbox-user:29\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"options\": [" +
				"          \"rbind\"," +
				"          \"ro\"," +
				"          \"rprivate\"" +
				"        ]," +
				"        \"source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"destination\": \"/run/user/1000\"," +
				"        \"options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"," +
				"          \"rprivate\"," +
				"          \"rw\"" +
				"        ]," +
				"        \"source\": \"/run/user/1000\"," +
				"        \"type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"destination\": \"/usr/bin/toolbox\"," +
				"        \"options\": [" +
				"          \"rbind\"," +
				"          \"ro\"," +
				"          \"rprivate\"" +
				"        ]," +
				"        \"source\": \"/some/prefix/binary\"," +
				"        \"type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-user-29\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 143," +
				"      \"Pid\": 5302," +
				"      \"Status\": \"exited\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "sleep",
				entryPointPID: 5302,
				id:            "b23f7c69ddec697f803b8acc40e85d212198c1baea9ffe193e7b3e0d2a020a39",
				image:         "localhost/fedora-toolbox-user:29",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "f29/fedora-toolbox",
					"version":                        "29",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
				},
				name:   "fedora-toolbox-user-29",
				status: "exited",
			},
		},
		{
			name: "podman 1.1.2, toolbx 0.0.9, running",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"sleep\"," +
				"        \"+Inf\"" +
				"      ]," +
				"      \"Image\": \"localhost/fedora-toolbox-user:29\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"f29/fedora-toolbox\"," +
				"        \"version\": \"29\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"b23f7c69ddec697f803b8acc40e85d212198c1baea9ffe193e7b3e0d2a020a39\"," +
				"    \"Image\": \"c7fab9e10750847a20b6664f485d9a1430b326eaf2799b932d15eebea36d6f5f\"," +
				"    \"ImageName\": \"localhost/fedora-toolbox-user:29\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"options\": [" +
				"          \"rbind\"," +
				"          \"ro\"," +
				"          \"rprivate\"" +
				"        ]," +
				"        \"source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"destination\": \"/run/user/1000\"," +
				"        \"options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"," +
				"          \"rprivate\"," +
				"          \"rw\"" +
				"        ]," +
				"        \"source\": \"/run/user/1000\"," +
				"        \"type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"destination\": \"/usr/bin/toolbox\"," +
				"        \"options\": [" +
				"          \"rbind\"," +
				"          \"ro\"," +
				"          \"rprivate\"" +
				"        ]," +
				"        \"source\": \"/some/prefix/binary\"," +
				"        \"type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-user-29\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Pid\": 5302," +
				"      \"Status\": \"running\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "sleep",
				entryPointPID: 5302,
				id:            "b23f7c69ddec697f803b8acc40e85d212198c1baea9ffe193e7b3e0d2a020a39",
				image:         "localhost/fedora-toolbox-user:29",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "f29/fedora-toolbox",
					"version":                        "29",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
				},
				name:   "fedora-toolbox-user-29",
				status: "running",
			},
		},
		{
			name: "podman 1.8.0, toolbx 0.0.18, configured",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--verbose\"," +
				"        \"init-container\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/f30/fedora-toolbox:30\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"f30/fedora-toolbox\"," +
				"        \"version\": \"30\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"4f8922191fc19f51fa120eda6b0bf0ca3c498469f30ee57a673e6c9ac2d0d4bb\"," +
				"    \"Image\": \"c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/f30/fedora-toolbox:30\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/host/monitor\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000/.flatpak-helper/monitor\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-30\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"configured\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "4f8922191fc19f51fa120eda6b0bf0ca3c498469f30ee57a673e6c9ac2d0d4bb",
				image:         "registry.fedoraproject.org/f30/fedora-toolbox:30",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "f30/fedora-toolbox",
					"version":                        "30",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/host/monitor",
					"/run/user/1000",
					"/usr/bin/toolbox",
				},
				name:   "fedora-toolbox-30",
				status: "configured",
			},
		},
		{
			name: "podman 1.8.0, toolbx 0.0.18, exited",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--verbose\"," +
				"        \"init-container\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/f30/fedora-toolbox:30\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"f30/fedora-toolbox\"," +
				"        \"version\": \"30\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"4f8922191fc19f51fa120eda6b0bf0ca3c498469f30ee57a673e6c9ac2d0d4bb\"," +
				"    \"Image\": \"c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/f30/fedora-toolbox:30\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/host/monitor\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000/.flatpak-helper/monitor\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-30\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 143," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"exited\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "4f8922191fc19f51fa120eda6b0bf0ca3c498469f30ee57a673e6c9ac2d0d4bb",
				image:         "registry.fedoraproject.org/f30/fedora-toolbox:30",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "f30/fedora-toolbox",
					"version":                        "30",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/host/monitor",
					"/run/user/1000",
					"/usr/bin/toolbox",
				},
				name:   "fedora-toolbox-30",
				status: "exited",
			},
		},
		{
			name: "podman 1.8.0, toolbx 0.0.18, running",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--verbose\"," +
				"        \"init-container\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/f30/fedora-toolbox:30\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"f30/fedora-toolbox\"," +
				"        \"version\": \"30\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"4f8922191fc19f51fa120eda6b0bf0ca3c498469f30ee57a673e6c9ac2d0d4bb\"," +
				"    \"Image\": \"c49513deb6160607504d2c9abf9523e81c02f69ea479fd07572a7a32b50beab8\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/f30/fedora-toolbox:30\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/host/monitor\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000/.flatpak-helper/monitor\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-30\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 11175," +
				"      \"Status\": \"running\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 11175,
				id:            "4f8922191fc19f51fa120eda6b0bf0ca3c498469f30ee57a673e6c9ac2d0d4bb",
				image:         "registry.fedoraproject.org/f30/fedora-toolbox:30",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "f30/fedora-toolbox",
					"version":                        "30",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/host/monitor",
					"/run/user/1000",
					"/usr/bin/toolbox",
				},
				name:   "fedora-toolbox-30",
				status: "running",
			},
		},
		{
			name: "podman 2.2.1, toolbx 0.0.99.1, configured",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--verbose\"," +
				"        \"init-container\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:32\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"32\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"7dfa257361547c0c67ed8678fe1c4de784b647c848deec2b0541cf040a1c64ee\"," +
				"    \"Image\": \"6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:32\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-32\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"configured\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "7dfa257361547c0c67ed8678fe1c4de784b647c848deec2b0541cf040a1c64ee",
				image:         "registry.fedoraproject.org/fedora-toolbox:32",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "fedora-toolbox",
					"version":                        "32",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-32",
				status: "configured",
			},
		},
		{
			name: "podman 2.2.1, toolbx 0.0.99.1, exited",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--verbose\"," +
				"        \"init-container\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:32\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"32\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"7dfa257361547c0c67ed8678fe1c4de784b647c848deec2b0541cf040a1c64ee\"," +
				"    \"Image\": \"6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:32\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-32\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 143," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"exited\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "7dfa257361547c0c67ed8678fe1c4de784b647c848deec2b0541cf040a1c64ee",
				image:         "registry.fedoraproject.org/fedora-toolbox:32",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "fedora-toolbox",
					"version":                        "32",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-32",
				status: "exited",
			},
		},
		{
			name: "podman 2.2.1, toolbx 0.0.99.1, running",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--verbose\"," +
				"        \"init-container\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:32\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"32\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"7dfa257361547c0c67ed8678fe1c4de784b647c848deec2b0541cf040a1c64ee\"," +
				"    \"Image\": \"6b2cbce8102fc0c0424b619ad199216c025efc374457dc7a61bb89d393e7eab6\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:32\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Name\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-32\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 4407," +
				"      \"Status\": \"running\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 4407,
				id:            "7dfa257361547c0c67ed8678fe1c4de784b647c848deec2b0541cf040a1c64ee",
				image:         "registry.fedoraproject.org/fedora-toolbox:32",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "fedora-toolbox",
					"version":                        "32",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-32",
				status: "running",
			},
		},
		{
			name: "podman 3.4.7, toolbx 0.0.99.3, configured",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:35\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"35\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"9effd41d07eea253926c08b7e61182d2cb6563abffc41c1ff7c1e57c42da1dab\"," +
				"    \"Image\": \"862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:35\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-35\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"configured\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "9effd41d07eea253926c08b7e61182d2cb6563abffc41c1ff7c1e57c42da1dab",
				image:         "registry.fedoraproject.org/fedora-toolbox:35",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "fedora-toolbox",
					"version":                        "35",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-35",
				status: "configured",
			},
		},
		{
			name: "podman 3.4.7, toolbx 0.0.99.3, exited",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:35\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"35\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"9effd41d07eea253926c08b7e61182d2cb6563abffc41c1ff7c1e57c42da1dab\"," +
				"    \"Image\": \"862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:35\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-35\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 143," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"exited\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "9effd41d07eea253926c08b7e61182d2cb6563abffc41c1ff7c1e57c42da1dab",
				image:         "registry.fedoraproject.org/fedora-toolbox:35",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "fedora-toolbox",
					"version":                        "35",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-35",
				status: "exited",
			},
		},
		{
			name: "podman 3.4.7, toolbx 0.0.99.3, running",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--monitor-host\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:35\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.github.debarshiray.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"35\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"9effd41d07eea253926c08b7e61182d2cb6563abffc41c1ff7c1e57c42da1dab\"," +
				"    \"Image\": \"862705390e8b1678bbac66beb30547e0ef59abd65b18e23ea533f059ba069227\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:35\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-35\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Healthcheck\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 8253," +
				"      \"Status\": \"running\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 8253,
				id:            "9effd41d07eea253926c08b7e61182d2cb6563abffc41c1ff7c1e57c42da1dab",
				image:         "registry.fedoraproject.org/fedora-toolbox:35",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox":  "true",
					"com.github.debarshiray.toolbox": "true",
					"com.redhat.component":           "fedora-toolbox",
					"name":                           "fedora-toolbox",
					"version":                        "35",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-35",
				status: "running",
			},
		},
		{
			name: "podman 4.9.4, toolbx 0.0.99.5, created",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:38\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"a4599f0effa73cb8051d0b5650e28be7f7f9cd6655a584c48c14e7075201b7d7\"," +
				"    \"Image\": \"e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:38\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-38\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Health\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"created\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "a4599f0effa73cb8051d0b5650e28be7f7f9cd6655a584c48c14e7075201b7d7",
				image:         "registry.fedoraproject.org/fedora-toolbox:38",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox": "true",
					"com.redhat.component":          "fedora-toolbox",
					"name":                          "fedora-toolbox",
					"version":                       "38",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-38",
				status: "created",
			},
		},
		{
			name: "podman 4.9.4, toolbx 0.0.99.5, exited",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:38\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"a4599f0effa73cb8051d0b5650e28be7f7f9cd6655a584c48c14e7075201b7d7\"," +
				"    \"Image\": \"e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:38\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-38\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 143," +
				"      \"Health\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 0," +
				"      \"Status\": \"exited\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "a4599f0effa73cb8051d0b5650e28be7f7f9cd6655a584c48c14e7075201b7d7",
				image:         "registry.fedoraproject.org/fedora-toolbox:38",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox": "true",
					"com.redhat.component":          "fedora-toolbox",
					"name":                          "fedora-toolbox",
					"version":                       "38",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-38",
				status: "exited",
			},
		},
		{
			name: "podman 4.9.4, toolbx 0.0.99.5, running",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:38\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"com.redhat.component\": \"fedora-toolbox\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"38\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"a4599f0effa73cb8051d0b5650e28be7f7f9cd6655a584c48c14e7075201b7d7\"," +
				"    \"Image\": \"e8c6a36c07b778f0efcf7adb0c317ea2405afed5a3547fe8272c54b2495955ce\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:38\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-38\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Health\": {" +
				"        \"FailingStreak\": 0," +
				"        \"Log\": null," +
				"        \"Status\": \"\"" +
				"      }," +
				"      \"Pid\": 11686," +
				"      \"Status\": \"running\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 11686,
				id:            "a4599f0effa73cb8051d0b5650e28be7f7f9cd6655a584c48c14e7075201b7d7",
				image:         "registry.fedoraproject.org/fedora-toolbox:38",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox": "true",
					"com.redhat.component":          "fedora-toolbox",
					"name":                          "fedora-toolbox",
					"version":                       "38",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-38",
				status: "running",
			},
		},
		{
			name: "podman 5.0.2, toolbx 0.0.99.5, created",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"6571a3f51998bccbee1608495c7bf28d42264b883c7cca9d03cfb6b5ef5f44f1\"," +
				"    \"Image\": \"27151f84995bacace815731ceccee16a902e6ed207a57ca4601ca9ad7b5c7a3c\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-40\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Pid\": 0," +
				"      \"Status\": \"created\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "6571a3f51998bccbee1608495c7bf28d42264b883c7cca9d03cfb6b5ef5f44f1",
				image:         "registry.fedoraproject.org/fedora-toolbox:40",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox": "true",
					"name":                          "fedora-toolbox",
					"version":                       "40",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-40",
				status: "created",
			},
		},
		{
			name: "podman 5.0.2, toolbx 0.0.99.5, exited",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"6571a3f51998bccbee1608495c7bf28d42264b883c7cca9d03cfb6b5ef5f44f1\"," +
				"    \"Image\": \"27151f84995bacace815731ceccee16a902e6ed207a57ca4601ca9ad7b5c7a3c\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-40\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 143," +
				"      \"Pid\": 0," +
				"      \"Status\": \"exited\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 0,
				id:            "6571a3f51998bccbee1608495c7bf28d42264b883c7cca9d03cfb6b5ef5f44f1",
				image:         "registry.fedoraproject.org/fedora-toolbox:40",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox": "true",
					"name":                          "fedora-toolbox",
					"version":                       "40",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-40",
				status: "exited",
			},
		},
		{
			name: "podman 5.0.2, toolbx 0.0.99.5, running",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"toolbox\"," +
				"        \"--log-level\"," +
				"        \"debug\"," +
				"        \"init-container\"," +
				"        \"--gid\"," +
				"        \"1000\"," +
				"        \"--home\"," +
				"        \"/home/user\"," +
				"        \"--shell\"," +
				"        \"/bin/bash\"," +
				"        \"--uid\"," +
				"        \"1000\"," +
				"        \"--user\"," +
				"        \"user\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"      \"Labels\": {" +
				"        \"com.github.containers.toolbox\": \"true\"," +
				"        \"name\": \"fedora-toolbox\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"6571a3f51998bccbee1608495c7bf28d42264b883c7cca9d03cfb6b5ef5f44f1\"," +
				"    \"Image\": \"27151f84995bacace815731ceccee16a902e6ed207a57ca4601ca9ad7b5c7a3c\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora-toolbox:40\"," +
				"    \"Mounts\": [" +
				"      {" +
				"        \"Destination\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/etc/profile.d/toolbox.sh\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/run/user/1000\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"nodev\"," +
				"          \"nosuid\"," +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": true," +
				"        \"Source\": \"/run/user/1000\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/usr/bin/toolbox\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": [" +
				"          \"rbind\"" +
				"        ]," +
				"        \"Propagation\": \"rprivate\"," +
				"        \"RW\": false," +
				"        \"Source\": \"/some/prefix/binary\"," +
				"        \"Type\": \"bind\"" +
				"      }," +
				"      {" +
				"        \"Destination\": \"/dev/pts\"," +
				"        \"Driver\": \"\"," +
				"        \"Mode\": \"\"," +
				"        \"Options\": []," +
				"        \"Propagation\": \"\"," +
				"        \"RW\": true," +
				"        \"Source\": \"devpts\"," +
				"        \"Type\": \"bind\"" +
				"      }" +
				"    ]," +
				"    \"Name\": \"fedora-toolbox-40\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Pid\": 3792," +
				"      \"Status\": \"running\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "toolbox",
				entryPointPID: 3792,
				id:            "6571a3f51998bccbee1608495c7bf28d42264b883c7cca9d03cfb6b5ef5f44f1",
				image:         "registry.fedoraproject.org/fedora-toolbox:40",
				isToolbx:      true,
				labels: map[string]string{
					"com.github.containers.toolbox": "true",
					"name":                          "fedora-toolbox",
					"version":                       "40",
				},
				mounts: []string{
					"/etc/profile.d/toolbox.sh",
					"/run/user/1000",
					"/usr/bin/toolbox",
					"/dev/pts",
				},
				name:   "fedora-toolbox-40",
				status: "running",
			},
		},
		{
			name: "podman 5.0.2, created",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"true\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora:40\"," +
				"      \"Labels\": {" +
				"        \"name\": \"fedora\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"f62203a35f867cbdeb0d340741455cea23bd5fccff19c33ef453aaa163152142\"," +
				"    \"Image\": \"4de8bd41536df94855e3fc830586d55477c3899c2484386d54013c0c1d0f1dd7\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora:40\"," +
				"    \"Mounts\": []," +
				"    \"Name\": \"zealous_chaum\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Pid\": 0," +
				"      \"Status\": \"created\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "true",
				entryPointPID: 0,
				id:            "f62203a35f867cbdeb0d340741455cea23bd5fccff19c33ef453aaa163152142",
				image:         "registry.fedoraproject.org/fedora:40",
				isToolbx:      false,
				labels: map[string]string{
					"name":    "fedora",
					"version": "40",
				},
				name:   "zealous_chaum",
				status: "created",
			},
		},
		{
			name: "podman 5.0.2, exited",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"true\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora:40\"," +
				"      \"Labels\": {" +
				"        \"name\": \"fedora\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"f62203a35f867cbdeb0d340741455cea23bd5fccff19c33ef453aaa163152142\"," +
				"    \"Image\": \"4de8bd41536df94855e3fc830586d55477c3899c2484386d54013c0c1d0f1dd7\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora:40\"," +
				"    \"Mounts\": []," +
				"    \"Name\": \"zealous_chaum\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Pid\": 0," +
				"      \"Status\": \"exited\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "true",
				entryPointPID: 0,
				id:            "f62203a35f867cbdeb0d340741455cea23bd5fccff19c33ef453aaa163152142",
				image:         "registry.fedoraproject.org/fedora:40",
				isToolbx:      false,
				labels: map[string]string{
					"name":    "fedora",
					"version": "40",
				},
				name:   "zealous_chaum",
				status: "exited",
			},
		},
		{
			name: "podman 5.0.2, running",
			data: "" +
				"[" +
				"  {" +
				"    \"Config\": {" +
				"      \"Cmd\": [" +
				"        \"true\"" +
				"      ]," +
				"      \"Image\": \"registry.fedoraproject.org/fedora:40\"," +
				"      \"Labels\": {" +
				"        \"name\": \"fedora\"," +
				"        \"version\": \"40\"" +
				"      }" +
				"    }," +
				"    \"Created\": \"2024-02-29T20:40:50.162797336+02:00\"," +
				"    \"Id\": \"f62203a35f867cbdeb0d340741455cea23bd5fccff19c33ef453aaa163152142\"," +
				"    \"Image\": \"4de8bd41536df94855e3fc830586d55477c3899c2484386d54013c0c1d0f1dd7\"," +
				"    \"ImageName\": \"registry.fedoraproject.org/fedora:40\"," +
				"    \"Mounts\": []," +
				"    \"Name\": \"zealous_chaum\"," +
				"    \"State\": {" +
				"      \"ExitCode\": 0," +
				"      \"Pid\": 4462," +
				"      \"Status\": \"running\"" +
				"    }" +
				"  }" +
				"]",
			expect: expect{
				entryPoint:    "true",
				entryPointPID: 4462,
				id:            "f62203a35f867cbdeb0d340741455cea23bd5fccff19c33ef453aaa163152142",
				image:         "registry.fedoraproject.org/fedora:40",
				isToolbx:      false,
				labels: map[string]string{
					"name":    "fedora",
					"version": "40",
				},
				name:   "zealous_chaum",
				status: "running",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.data)
			var containers []containerInspect
			err := json.Unmarshal(data, &containers)
			assert.NoError(t, err)
			assert.Len(t, containers, 1)

			container := containers[0]
			assert.Equal(t, tc.expect.entryPoint, container.EntryPoint())
			assert.Equal(t, tc.expect.entryPointPID, container.EntryPointPID())
			assert.Equal(t, tc.expect.id, container.ID())
			assert.Equal(t, tc.expect.image, container.Image())
			assert.Equal(t, tc.expect.isToolbx, container.IsToolbx())
			assert.Equal(t, tc.expect.labels, container.Labels())
			assert.Equal(t, tc.expect.mounts, container.Mounts())
			assert.Equal(t, tc.expect.name, container.Name())
			assert.Equal(t, []string{tc.expect.name}, container.Names())
			assert.Equal(t, tc.expect.status, container.Status())
		})
	}
}

/*
 * Copyright © 2019 – 2026 Red Hat Inc.
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
	"encoding/json"
	"testing"

	"github.com/containers/toolbox/pkg/architecture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageGetSize(t *testing.T) {
	testCases := []struct {
		name   string
		image  Image
		expect float64
		errMsg string
	}{
		{
			name: "single layer",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("73528142")},
				},
			},
			expect: 73528142,
		},
		{
			name: "multiple layers",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("73528142")},
					{Size: json.Number("1849805")},
					{Size: json.Number("512")},
				},
			},
			expect: 75378459,
		},
		{
			name: "zero-size layers",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("0")},
					{Size: json.Number("0")},
				},
			},
			expect: 0,
		},
		{
			name: "empty LayersData slice",
			image: Image{
				LayersData: []Layer{},
			},
			expect: 0,
		},
		{
			name: "nil LayersData",
			image: Image{
				LayersData: nil,
			},
			expect: -1,
			errMsg: "'skopeo inspect' did not have LayersData",
		},
		{
			name: "invalid size value",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("not-a-number")},
				},
			},
			expect: -1,
			errMsg: "strconv.ParseFloat: parsing \"not-a-number\": invalid syntax",
		},
		{
			name: "large layer sizes",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("2147483648")},
					{Size: json.Number("2147483648")},
				},
			},
			expect: 4294967296,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			size, err := tc.image.GetSize()

			if tc.errMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expect, size)
		})
	}
}

func TestImageGetSizeHuman(t *testing.T) {
	testCases := []struct {
		name   string
		image  Image
		expect string
		errMsg string
	}{
		{
			name: "typical image size",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("73528142")},
					{Size: json.Number("1849805")},
				},
			},
			expect: "75.38MB",
		},
		{
			name: "large image size over 1 GB",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("2204748042")},
				},
			},
			expect: "2.205GB",
		},
		{
			name: "zero size",
			image: Image{
				LayersData: []Layer{
					{Size: json.Number("0")},
				},
			},
			expect: "0B",
		},
		{
			name: "nil LayersData",
			image: Image{
				LayersData: nil,
			},
			errMsg: "'skopeo inspect' did not have LayersData",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sizeHuman, err := tc.image.GetSizeHuman()

			if tc.errMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errMsg)
				assert.Empty(t, sizeHuman)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, sizeHuman)
			}
		})
	}
}

func TestImageVerifyArchitectureMatch(t *testing.T) {
	testCases := []struct {
		name           string
		image          Image
		expectedArchID int
		errMsg         string
	}{
		{
			name: "amd64 image matches x86_64",
			image: Image{
				Architecture: "amd64",
				NameFull:     "registry.fedoraproject.org/fedora-toolbox:41",
			},
			expectedArchID: architecture.X86_64,
		},
		{
			name: "arm64 image matches aarch64",
			image: Image{
				Architecture: "arm64",
				NameFull:     "registry.fedoraproject.org/fedora-toolbox:41-aarch64",
			},
			expectedArchID: architecture.Aarch64,
		},
		{
			name: "ppc64le image matches ppc64le",
			image: Image{
				Architecture: "ppc64le",
				NameFull:     "registry.fedoraproject.org/fedora-toolbox:41-ppc64le",
			},
			expectedArchID: architecture.Ppc64le,
		},
		{
			name: "amd64 image does not match aarch64",
			image: Image{
				Architecture: "amd64",
				NameFull:     "registry.fedoraproject.org/fedora-toolbox:41",
			},
			expectedArchID: architecture.Aarch64,
			errMsg:         "image registry.fedoraproject.org/fedora-toolbox:41 is a single-architecture image for amd64, but arm64 was requested",
		},
		{
			name: "arm64 image does not match x86_64",
			image: Image{
				Architecture: "arm64",
				NameFull:     "registry.fedoraproject.org/fedora-toolbox:41-aarch64",
			},
			expectedArchID: architecture.X86_64,
			errMsg:         "image registry.fedoraproject.org/fedora-toolbox:41-aarch64 is a single-architecture image for arm64, but amd64 was requested",
		},
		{
			name: "unsupported architecture in image",
			image: Image{
				Architecture: "mips",
				NameFull:     "example.com/custom-image:latest",
			},
			expectedArchID: architecture.X86_64,
			errMsg:         "architecture 'mips' is not supported by Toolbx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.image.VerifyArchitectureMatch(tc.expectedArchID)

			if tc.errMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInspectJSONUnmarshal(t *testing.T) {
	type expect struct {
		architecture   string
		layersCount    int
		totalSize      float64
		matchArchID    int
		mismatchArchID int
	}

	testCases := []struct {
		name   string
		data   string
		expect expect
	}{
		{
			name: "fedora_38,skopeo_1.15.0",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:efde4efa9f18e619f62a14849959a436fc8062b2248aca794a0f4957b1d2eeca\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2024-02-01T19:09:35.77181Z\"," +
				"  \"DockerVersion\": \"1.13.1\"," +
				"  \"Labels\": {" +
				"    \"architecture\": \"x86_64\"," +
				"    \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"    \"build-date\": \"2024-02-01T19:07:19.783946\"," +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"com.redhat.build-host\": \"osbs-node01.iad2.fedoraproject.org\"," +
				"    \"com.redhat.component\": \"fedora-toolbox\"," +
				"    \"distribution-scope\": \"public\"," +
				"    \"license\": \"MIT\"," +
				"    \"maintainer\": \"Debarshi Ray \\u003crishi@fedoraproject.org\\u003e\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"release\": \"20\"," +
				"    \"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"    \"usage\": \"This image is meant to be used with the toolbox command\"," +
				"    \"vcs-ref\": \"c36ca58a44f947077b7c62e17d7b4fe4cd1b5797\"," +
				"    \"vcs-type\": \"git\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"38\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:0642c0ebdff310231260783e8d24b6b322dd220d43ac297538242fe87350fcb6\"," +
				"    \"sha256:badbfb3b457582ca3e9dcb7617fe76c45cedfb6bd4f3e6970703ac0161c0e591\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\"," +
				"      \"Digest\": \"sha256:0642c0ebdff310231260783e8d24b6b322dd220d43ac297538242fe87350fcb6\"," +
				"      \"Size\": 68602971," +
				"      \"Annotations\": null" +
				"    }," +
				"    {" +
				"      \"MIMEType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\"," +
				"      \"Digest\": \"sha256:badbfb3b457582ca3e9dcb7617fe76c45cedfb6bd4f3e6970703ac0161c0e591\"," +
				"      \"Size\": 248948495," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"DISTTAG=f38container\"," +
				"    \"FGC=f38\"," +
				"    \"container=oci\"," +
				"    \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    2,
				totalSize:      317551466,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
		{
			name: "fedora_38,skopeo_1.15.0,architecture_arm64",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:efde4efa9f18e619f62a14849959a436fc8062b2248aca794a0f4957b1d2eeca\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2024-02-01T19:10:17.64198Z\"," +
				"  \"DockerVersion\": \"19.03.13\"," +
				"  \"Labels\": {" +
				"    \"architecture\": \"arm64\"," +
				"    \"authoritative-source-url\": \"registry.fedoraproject.org\"," +
				"    \"build-date\": \"2024-02-01T19:07:18.853192\"," +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"com.redhat.build-host\": \"osbs-aarch64-node01.iad2.fedoraproject.org\"," +
				"    \"com.redhat.component\": \"fedora-toolbox\"," +
				"    \"distribution-scope\": \"public\"," +
				"    \"license\": \"MIT\"," +
				"    \"maintainer\": \"Debarshi Ray \\u003crishi@fedoraproject.org\\u003e\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"release\": \"20\"," +
				"    \"summary\": \"Base image for creating Fedora toolbox containers\"," +
				"    \"usage\": \"This image is meant to be used with the toolbox command\"," +
				"    \"vcs-ref\": \"c36ca58a44f947077b7c62e17d7b4fe4cd1b5797\"," +
				"    \"vcs-type\": \"git\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"38\"" +
				"  }," +
				"  \"Architecture\": \"arm64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:9104857c647cd9953bca093a735e9aeb8940d9852ae2ed4080def13a55506260\"," +
				"    \"sha256:ab482b93b9e3916c98db0482b961294e4624501a1f4e4920fe202d3d09e6086a\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\"," +
				"      \"Digest\": \"sha256:9104857c647cd9953bca093a735e9aeb8940d9852ae2ed4080def13a55506260\"," +
				"      \"Size\": 67268628," +
				"      \"Annotations\": null" +
				"    }," +
				"    {" +
				"      \"MIMEType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\"," +
				"      \"Digest\": \"sha256:ab482b93b9e3916c98db0482b961294e4624501a1f4e4920fe202d3d09e6086a\"," +
				"      \"Size\": 279142877," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"DISTTAG=f38container\"," +
				"    \"FGC=f38\"," +
				"    \"container=oci\"," +
				"    \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "arm64",
				layersCount:    2,
				totalSize:      346411505,
				matchArchID:    architecture.Aarch64,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_39,skopeo_1.16.1",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:11f3634d4a8f2d4a69c4ad8442133f69979be49fa6269eccc6ab0863c39d59d0\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2024-11-26T07:51:25Z\"," +
				"  \"DockerVersion\": \"1.10.1\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"39\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:40eb91ad8a1c0e03d84a95ef38af4fe3caac449d78e796539a6062056e1f5777\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\"," +
				"      \"Digest\": \"sha256:40eb91ad8a1c0e03d84a95ef38af4fe3caac449d78e796539a6062056e1f5777\"," +
				"      \"Size\": 362650126," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"DISTTAG=f39container\"," +
				"    \"FGC=f39\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    1,
				totalSize:      362650126,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
		{
			name: "fedora_39,skopeo_1.16.1,architecture_arm64",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:11f3634d4a8f2d4a69c4ad8442133f69979be49fa6269eccc6ab0863c39d59d0\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2024-11-26T07:50:54Z\"," +
				"  \"DockerVersion\": \"1.10.1\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"39\"" +
				"  }," +
				"  \"Architecture\": \"arm64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:56d6ea82c48514304c4dbf3a8296a159458440f4835520ae1783429fa3acd792\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\"," +
				"      \"Digest\": \"sha256:56d6ea82c48514304c4dbf3a8296a159458440f4835520ae1783429fa3acd792\"," +
				"      \"Size\": 336450730," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"DISTTAG=f39container\"," +
				"    \"FGC=f39\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "arm64",
				layersCount:    1,
				totalSize:      336450730,
				matchArchID:    architecture.Aarch64,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_39,skopeo_1.16.1,architecture_ppc64le",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:11f3634d4a8f2d4a69c4ad8442133f69979be49fa6269eccc6ab0863c39d59d0\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2024-11-26T07:58:00Z\"," +
				"  \"DockerVersion\": \"1.10.1\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"39\"" +
				"  }," +
				"  \"Architecture\": \"ppc64le\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:014b7e796167e4c90ac97df346776a9253b099f11ddff67656f100de92bedf5f\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.docker.image.rootfs.diff.tar.gzip\"," +
				"      \"Digest\": \"sha256:014b7e796167e4c90ac97df346776a9253b099f11ddff67656f100de92bedf5f\"," +
				"      \"Size\": 346423938," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"DISTTAG=f39container\"," +
				"    \"FGC=f39\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "ppc64le",
				layersCount:    1,
				totalSize:      346423938,
				matchArchID:    architecture.Ppc64le,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_40,skopeo_1.18.0",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:c3cc9836c9e55475f85496f9369e925164f1b7cd832b55e22c13b74840576c31\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2025-05-13T07:49:00.272423532Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.39.2\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"40\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"40\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:3fa9653b37941f2c2fcdbdaaf6ac8d97496e3e331cc4aeb27629bd65ee5b94a8\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:3fa9653b37941f2c2fcdbdaaf6ac8d97496e3e331cc4aeb27629bd65ee5b94a8\"," +
				"      \"Size\": 378214266," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    1,
				totalSize:      378214266,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
		{
			name: "fedora_40,skopeo_1.18.0,architecture_arm64",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:c3cc9836c9e55475f85496f9369e925164f1b7cd832b55e22c13b74840576c31\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2025-05-13T07:48:52.762718634Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.39.2\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"40\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"40\"" +
				"  }," +
				"  \"Architecture\": \"arm64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:8c8790edeeac9e5823cf1442623d8eac2d83c1b8a12465351e9c159a33a89565\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:8c8790edeeac9e5823cf1442623d8eac2d83c1b8a12465351e9c159a33a89565\"," +
				"      \"Size\": 357895766," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "arm64",
				layersCount:    1,
				totalSize:      357895766,
				matchArchID:    architecture.Aarch64,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_40,skopeo_1.18.0,architecture_ppc64le",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:c3cc9836c9e55475f85496f9369e925164f1b7cd832b55e22c13b74840576c31\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2025-05-13T07:50:04.841519173Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.39.2\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"40\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"40\"" +
				"  }," +
				"  \"Architecture\": \"ppc64le\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:6fa7808a885e4c0597fbe0da6832bfb114a1baf4cd0383a7f6ad14fd14c713fc\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:6fa7808a885e4c0597fbe0da6832bfb114a1baf4cd0383a7f6ad14fd14c713fc\"," +
				"      \"Size\": 362731624," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "ppc64le",
				layersCount:    1,
				totalSize:      362731624,
				matchArchID:    architecture.Ppc64le,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_41,skopeo_1.20.0",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:f5f4b8d1c904810404117d99325e7130452c71e94cd78077e0141d3f0141eda0\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2025-12-15T07:47:49.048519394Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.41.5\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"41\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"41\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:4050d18ab275ca85fcf4de4ea4f01355ecf605a1ac1c944fe540e3a1838e1559\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:4050d18ab275ca85fcf4de4ea4f01355ecf605a1ac1c944fe540e3a1838e1559\"," +
				"      \"Size\": 387608637," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    1,
				totalSize:      387608637,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
		{
			name: "fedora_41,skopeo_1.20.0,architecture_arm64",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:f5f4b8d1c904810404117d99325e7130452c71e94cd78077e0141d3f0141eda0\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2025-12-15T07:47:52.135380621Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.41.5\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"41\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"41\"" +
				"  }," +
				"  \"Architecture\": \"arm64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:3b9ad2fb7a433b49659fa781d9406a4faa86a92a284ebfcc9eeb2accae5fa8da\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:3b9ad2fb7a433b49659fa781d9406a4faa86a92a284ebfcc9eeb2accae5fa8da\"," +
				"      \"Size\": 358789413," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "arm64",
				layersCount:    1,
				totalSize:      358789413,
				matchArchID:    architecture.Aarch64,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_41,skopeo_1.20.0,architecture_ppc64le",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:f5f4b8d1c904810404117d99325e7130452c71e94cd78077e0141d3f0141eda0\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2025-12-15T07:48:34.495136412Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.41.5\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"41\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"41\"" +
				"  }," +
				"  \"Architecture\": \"ppc64le\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:9920118d3124e87eb834b4777f83ffb19a455d8cd3d7c64f95fbad6473acc43b\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:9920118d3124e87eb834b4777f83ffb19a455d8cd3d7c64f95fbad6473acc43b\"," +
				"      \"Size\": 364425315," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "ppc64le",
				layersCount:    1,
				totalSize:      364425315,
				matchArchID:    architecture.Ppc64le,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_42,skopeo_1.22.2",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:c4ebfe5df39c582a1d15c1a2251c6f7060bb07568f28913289e73b1b72b882ab\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2026-05-03T07:47:44.819587225Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.43.1\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.licenses\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.title\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"42\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"42\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:285a4586579509386fd5c5899923ced5b72c21c29682b32ed0ba0c68b909d04c\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:285a4586579509386fd5c5899923ced5b72c21c29682b32ed0ba0c68b909d04c\"," +
				"      \"Size\": 369354285," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/bin:/usr/bin\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    1,
				totalSize:      369354285,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
		{
			name: "fedora_42,skopeo_1.22.2,architecture_arm64",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:c4ebfe5df39c582a1d15c1a2251c6f7060bb07568f28913289e73b1b72b882ab\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2026-05-03T07:48:18.994213287Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.43.1\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.licenses\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.title\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"42\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"42\"" +
				"  }," +
				"  \"Architecture\": \"arm64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:fd87ecd39d657fd1bde5020ebc33460925b94ac43a4d09e8c5ad535504d746b9\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:fd87ecd39d657fd1bde5020ebc33460925b94ac43a4d09e8c5ad535504d746b9\"," +
				"      \"Size\": 349545018," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/bin:/usr/bin\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "arm64",
				layersCount:    1,
				totalSize:      349545018,
				matchArchID:    architecture.Aarch64,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_42,skopeo_1.22.2,architecture_ppc64le",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:c4ebfe5df39c582a1d15c1a2251c6f7060bb07568f28913289e73b1b72b882ab\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2026-05-03T07:49:12.61781795Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.43.1\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.licenses\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.title\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"42\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"42\"" +
				"  }," +
				"  \"Architecture\": \"ppc64le\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:db1e644afebaba886ca48a84fbd6f7aa3dd0f53e03af807f06b405a9a0fd443d\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:db1e644afebaba886ca48a84fbd6f7aa3dd0f53e03af807f06b405a9a0fd443d\"," +
				"      \"Size\": 358713438," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/bin:/usr/bin\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "ppc64le",
				layersCount:    1,
				totalSize:      358713438,
				matchArchID:    architecture.Ppc64le,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_43,skopeo_1.22.2",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:06c9a93965c14d87dcb75ab5e03df0232b5ed71ccc9e98b818bdf665d3d095ab\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2026-05-02T05:48:04.671009859Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.43.1\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.licenses\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.title\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"43\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"43\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:07e96dd49dae8b7557b1a483eb263f1f96def8f7a56ec76c70aa06802264f722\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:07e96dd49dae8b7557b1a483eb263f1f96def8f7a56ec76c70aa06802264f722\"," +
				"      \"Size\": 359026274," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/bin:/usr/bin\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    1,
				totalSize:      359026274,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
		{
			name: "fedora_43,skopeo_1.22.2,architecture_arm64",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:06c9a93965c14d87dcb75ab5e03df0232b5ed71ccc9e98b818bdf665d3d095ab\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2026-05-02T05:48:22.512575956Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.43.1\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.licenses\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.title\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"43\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"43\"" +
				"  }," +
				"  \"Architecture\": \"arm64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:467ccd74d1af1ef3dc3122490cf0c5083a8bd8c6dbe01be2ddb73921db7e8c64\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:467ccd74d1af1ef3dc3122490cf0c5083a8bd8c6dbe01be2ddb73921db7e8c64\"," +
				"      \"Size\": 339693369," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/bin:/usr/bin\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "arm64",
				layersCount:    1,
				totalSize:      339693369,
				matchArchID:    architecture.Aarch64,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "fedora_43,skopeo_1.22.2,architecture_ppc64le",
			data: "" +
				"{" +
				"  \"Name\": \"registry.fedoraproject.org/fedora-toolbox\"," +
				"  \"Digest\": \"sha256:06c9a93965c14d87dcb75ab5e03df0232b5ed71ccc9e98b818bdf665d3d095ab\"," +
				"  \"RepoTags\": [" +
				"    \"32\", \"32-12\", \"32-13\", \"32-14\", \"33\"," +
				"    \"33-12\", \"33-13\", \"33-14\", \"33-15\", \"33-16\"," +
				"    \"33-17\", \"33-18\", \"34\", \"34-12\", \"34-13\"," +
				"    \"34-14\", \"34-15\", \"34-16\", \"34-17\", \"34-18\"," +
				"    \"34-19\", \"34-20\", \"34-21\", \"34-22\", \"34-7\"," +
				"    \"34-9\", \"35\", \"35-1\", \"35-11\", \"35-12\"," +
				"    \"35-13\", \"35-14\", \"35-15\", \"35-16\", \"35-17\"," +
				"    \"35-18\", \"35-3\", \"35-4\", \"35-5\", \"35-6\"," +
				"    \"35-8\", \"35-9\", \"36\", \"36-10\", \"36-14\"," +
				"    \"36-15\", \"36-16\", \"36-2\", \"36-3\", \"36-4\"," +
				"    \"36-5\", \"36-6\", \"36-7\", \"36-8\", \"37\"," +
				"    \"37-1\", \"37-13\", \"37-14\", \"37-15\", \"37-16\"," +
				"    \"37-17\", \"37-2\", \"37-20\", \"37-3\", \"37-4\"," +
				"    \"37-5\", \"37-6\", \"37-8\", \"38\", \"38-1\"," +
				"    \"38-12\", \"38-13\", \"38-14\", \"38-15\", \"38-16\"," +
				"    \"38-17\", \"38-19\", \"38-2\", \"38-20\", \"38-3\"," +
				"    \"38-4\", \"39\", \"39-1\", \"39-2\", \"39-3\"," +
				"    \"39-4\", \"39-aarch64\", \"39-ppc64le\", \"39-s390x\", \"39-x86_64\"," +
				"    \"40\", \"40-aarch64\", \"40-ppc64le\", \"40-s390x\", \"40-x86_64\"," +
				"    \"41\", \"41-aarch64\", \"41-ppc64le\", \"41-s390x\", \"41-x86_64\"," +
				"    \"42\", \"42-aarch64\", \"42-ppc64le\", \"42-s390x\", \"42-x86_64\"," +
				"    \"43\", \"43-aarch64\", \"43-ppc64le\", \"43-s390x\", \"43-x86_64\"," +
				"    \"44\", \"44-aarch64\", \"44-ppc64le\", \"44-s390x\", \"44-x86_64\"," +
				"    \"45\", \"45-aarch64\", \"45-ppc64le\", \"45-s390x\", \"45-x86_64\"," +
				"    \"latest\", \"rawhide\", \"testing\"" +
				"  ]," +
				"  \"Created\": \"2026-05-02T05:49:14.078771354Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.43.1\"," +
				"    \"license\": \"MIT\"," +
				"    \"name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.license\": \"MIT\"," +
				"    \"org.opencontainers.image.licenses\": \"MIT\"," +
				"    \"org.opencontainers.image.name\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.title\": \"fedora-toolbox\"," +
				"    \"org.opencontainers.image.url\": \"https://fedoraproject.org/\"," +
				"    \"org.opencontainers.image.vendor\": \"Fedora Project\"," +
				"    \"org.opencontainers.image.version\": \"43\"," +
				"    \"vendor\": \"Fedora Project\"," +
				"    \"version\": \"43\"" +
				"  }," +
				"  \"Architecture\": \"ppc64le\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:ad3fe53c004ef47eb1fa876a2487d785732aeda88f41e969e02598dbf10ea1ca\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:ad3fe53c004ef47eb1fa876a2487d785732aeda88f41e969e02598dbf10ea1ca\"," +
				"      \"Size\": 343373524," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/bin:/usr/bin\"," +
				"    \"container=oci\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "ppc64le",
				layersCount:    1,
				totalSize:      343373524,
				matchArchID:    architecture.Ppc64le,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "ubuntu_24.04,skopeo_1.22.2",
			data: "" +
				"{" +
				"  \"Name\": \"quay.io/toolbx/ubuntu-toolbox\"," +
				"  \"Digest\": \"sha256:2989cc44ce8e136e5e33cdee1d0bb5fded5b782b4e863365762ac96d1f7fd7e1\"," +
				"  \"RepoTags\": [" +
				"    \"16.04\", \"18.04\", \"20.04\", \"22.04\", \"22.10\"," +
				"    \"23.04\", \"23.10\", \"24.04\", \"24.10\", \"25.04\"," +
				"    \"latest\"" +
				"  ]," +
				"  \"Created\": \"2026-04-27T00:59:54.842771346Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.33.7\"," +
				"    \"maintainer\": \"Ievgen Popovych \\u003cjmennius@gmail.com\\u003e\"," +
				"    \"name\": \"ubuntu-toolbox\"," +
				"    \"org.opencontainers.image.version\": \"24.04\"," +
				"    \"summary\": \"Base image for creating Ubuntu Toolbx containers\"," +
				"    \"usage\": \"This image is meant to be used with the toolbox command\"," +
				"    \"version\": \"24.04\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:9c6d3dc912b4c7f77d573326b2040a3d210efaf389734ccd1d335237e146b67a\"," +
				"    \"sha256:befa25c0402b5bd10c3b85e6f9ef7fbb8dec14b1f1fe2f7cba504f414ea96211\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:9c6d3dc912b4c7f77d573326b2040a3d210efaf389734ccd1d335237e146b67a\"," +
				"      \"Size\": 31654253," +
				"      \"Annotations\": null" +
				"    }," +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:befa25c0402b5bd10c3b85e6f9ef7fbb8dec14b1f1fe2f7cba504f414ea96211\"," +
				"      \"Size\": 175991436," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    2,
				totalSize:      207645689,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
		{
			name: "ubuntu_24.04,skopeo_1.22.2,architecture_arm64",
			data: "" +
				"{" +
				"  \"Name\": \"quay.io/toolbx/ubuntu-toolbox\"," +
				"  \"Digest\": \"sha256:2989cc44ce8e136e5e33cdee1d0bb5fded5b782b4e863365762ac96d1f7fd7e1\"," +
				"  \"RepoTags\": [" +
				"    \"16.04\", \"18.04\", \"20.04\", \"22.04\", \"22.10\"," +
				"    \"23.04\", \"23.10\", \"24.04\", \"24.10\", \"25.04\"," +
				"    \"latest\"" +
				"  ]," +
				"  \"Created\": \"2026-04-27T01:28:47.169971593Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.33.7\"," +
				"    \"maintainer\": \"Ievgen Popovych \\u003cjmennius@gmail.com\\u003e\"," +
				"    \"name\": \"ubuntu-toolbox\"," +
				"    \"org.opencontainers.image.version\": \"24.04\"," +
				"    \"summary\": \"Base image for creating Ubuntu Toolbx containers\"," +
				"    \"usage\": \"This image is meant to be used with the toolbox command\"," +
				"    \"version\": \"24.04\"" +
				"  }," +
				"  \"Architecture\": \"arm64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:77efc964b64f9c6338305d55a8c7042eeb441ee5d9b23f2e37ad6684bcc939cb\"," +
				"    \"sha256:b1e40f994cd2a81d2dfdcba9bccc3e8391253cf9eb94ae2219faab2b6a5948e3\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:77efc964b64f9c6338305d55a8c7042eeb441ee5d9b23f2e37ad6684bcc939cb\"," +
				"      \"Size\": 30753035," +
				"      \"Annotations\": null" +
				"    }," +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:b1e40f994cd2a81d2dfdcba9bccc3e8391253cf9eb94ae2219faab2b6a5948e3\"," +
				"      \"Size\": 174479376," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "arm64",
				layersCount:    2,
				totalSize:      205232411,
				matchArchID:    architecture.Aarch64,
				mismatchArchID: architecture.X86_64,
			},
		},
		{
			name: "arch,skopeo_1.22.2",
			data: "" +
				"{" +
				"  \"Name\": \"quay.io/toolbx/arch-toolbox\"," +
				"  \"Digest\": \"sha256:40a85480d3e8c72d9b2a56f8615b66accae502b9d233b3a18266121656467e5d\"," +
				"  \"RepoTags\": [" +
				"    \"latest\"" +
				"  ]," +
				"  \"Created\": \"2026-04-27T00:28:04.418383637Z\"," +
				"  \"DockerVersion\": \"\"," +
				"  \"Labels\": {" +
				"    \"com.github.containers.toolbox\": \"true\"," +
				"    \"io.buildah.version\": \"1.33.7\"," +
				"    \"maintainer\": \"Morten Linderud \\u003cfoxboron@archlinux.org\\u003e\"," +
				"    \"name\": \"arch-toolbox\"," +
				"    \"org.opencontainers.image.authors\": \"Santiago Torres-Arias \\u003csantiago@archlinux.org\\u003e (@SantiagoTorres), Christian Rebischke \\u003cChris.Rebischke@archlinux.org\\u003e (@shibumi), Justin Kromlinger \\u003chashworks@archlinux.org\\u003e (@hashworks)\"," +
				"    \"org.opencontainers.image.created\": \"2026-04-19T00:06:37+00:00\"," +
				"    \"org.opencontainers.image.description\": \"Official containerd image of Arch Linux, a simple, lightweight Linux distribution aimed for flexibility.\"," +
				"    \"org.opencontainers.image.documentation\": \"https://wiki.archlinux.org/title/Docker#Arch_Linux\"," +
				"    \"org.opencontainers.image.licenses\": \"GPL-3.0-or-later\"," +
				"    \"org.opencontainers.image.revision\": \"0d7c4c0017584f9bcb495105cc412d6575f04564\"," +
				"    \"org.opencontainers.image.source\": \"https://gitlab.archlinux.org/archlinux/archlinux-docker\"," +
				"    \"org.opencontainers.image.title\": \"Arch Linux base-devel Image\"," +
				"    \"org.opencontainers.image.url\": \"https://gitlab.archlinux.org/archlinux/archlinux-docker/-/blob/master/README.md\"," +
				"    \"org.opencontainers.image.version\": \"20260419.0.517065\"," +
				"    \"summary\": \"Base image for creating Arch Linux Toolbx containers\"," +
				"    \"usage\": \"This image is meant to be used with the toolbox command\"," +
				"    \"version\": \"base-devel\"" +
				"  }," +
				"  \"Architecture\": \"amd64\"," +
				"  \"Os\": \"linux\"," +
				"  \"Layers\": [" +
				"    \"sha256:9713e57a91bb6c4d42be627d18cb764a4eb6221fa6a6fec0f66312bbc4279714\"," +
				"    \"sha256:8e68f3509eb65780eb16bd28716312a4f2ef3613e694ff4cfadd156b779c99af\"," +
				"    \"sha256:9e59132d6ea60eb12fe1792013e57e192c2f51c133d07f977c97f1fc16a69462\"" +
				"  ]," +
				"  \"LayersData\": [" +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:9713e57a91bb6c4d42be627d18cb764a4eb6221fa6a6fec0f66312bbc4279714\"," +
				"      \"Size\": 256601029," +
				"      \"Annotations\": null" +
				"    }," +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:8e68f3509eb65780eb16bd28716312a4f2ef3613e694ff4cfadd156b779c99af\"," +
				"      \"Size\": 9659," +
				"      \"Annotations\": null" +
				"    }," +
				"    {" +
				"      \"MIMEType\": \"application/vnd.oci.image.layer.v1.tar+gzip\"," +
				"      \"Digest\": \"sha256:9e59132d6ea60eb12fe1792013e57e192c2f51c133d07f977c97f1fc16a69462\"," +
				"      \"Size\": 361785305," +
				"      \"Annotations\": null" +
				"    }" +
				"  ]," +
				"  \"Env\": [" +
				"    \"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\"," +
				"    \"LANG=C.UTF-8\"" +
				"  ]" +
				"}",
			expect: expect{
				architecture:   "amd64",
				layersCount:    3,
				totalSize:      618395993,
				matchArchID:    architecture.X86_64,
				mismatchArchID: architecture.Aarch64,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.data)
			var image Image
			err := json.Unmarshal(data, &image)
			require.NoError(t, err)

			assert.Equal(t, tc.expect.architecture, image.Architecture)
			assert.Len(t, image.LayersData, tc.expect.layersCount)

			size, err := image.GetSize()
			assert.NoError(t, err)
			assert.Equal(t, tc.expect.totalSize, size)

			image.NameFull = "test-image"

			err = image.VerifyArchitectureMatch(tc.expect.matchArchID)
			assert.NoError(t, err)

			err = image.VerifyArchitectureMatch(tc.expect.mismatchArchID)
			assert.Error(t, err)
		})
	}
}

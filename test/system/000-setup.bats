#!/usr/bin/env bats
#
# Copyright © 2021 – 2022 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

@test "test suite: Setup" {
    local os_release="$(find_os_release)"
    local system_id="$(get_system_id)"
    local system_version="$(get_system_version)"

    assert [ -n "$os_release" ]
    assert [ -n "$system_id" ]
    assert [ -n "$system_version" ]

    _setup_environment
    # Cache the default image for the system
    _pull_and_cache_distro_image "$system_id" "$system_version"

    # Cache all images that will be needed during the tests
    _pull_and_cache_distro_image fedora 32
    _pull_and_cache_distro_image busybox

    # If run on Fedora Rawhide, cache 2 extra images (previous Fedora versions)
    local rawhide_res="$(awk '/rawhide/' $os_release)"
    if [ "$system_id" = "fedora" ] && [ -n "$rawhide_res" ]; then
        _pull_and_cache_distro_image fedora "$((system_version-1))"
        _pull_and_cache_distro_image fedora "$((system_version-2))"
    fi

    _setup_docker_registry
}

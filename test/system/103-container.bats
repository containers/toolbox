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

setup() {
  _setup_environment
  cleanup_containers
}

teardown() {
  cleanup_containers
}


@test "container: Check container starts without issues" {
  readonly CONTAINER_NAME="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  run container_started $CONTAINER_NAME
  assert_success
}

@test "container(Fedora Rawhide): Containers with supported versions start without issues" {
  local os_release="$(find_os_release)"
  local system_id="$(get_system_id)"
  local system_version="$(get_system_version)"
  local rawhide_res="$(awk '/rawhide/' $os_release)"

  if [ "$system_id" != "fedora" ] || [ -z "$rawhide_res" ]; then
    skip "This test is only for Fedora Rawhide"
  fi

  create_distro_container "$system_id" "$system_version" latest
  run container_started latest
  assert_success

  create_distro_container "$system_id" "$((system_version-1))" second
  run container_started second
  assert_success

  create_distro_container "$system_id" "$((system_version-2))" third
  run container_started third
  assert_success
}

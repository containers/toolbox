# shellcheck shell=bats
#
# Copyright © 2021 – 2025 Red Hat, Inc.
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

# bats file_tags=commands-options

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.10.0
  cleanup_all
}

teardown() {
  cleanup_all
}


@test "container: Check container starts without issues" {
  CONTAINER_NAME="$(get_system_id)-toolbox-$(get_system_version)"
  readonly CONTAINER_NAME

  create_default_container

  run container_started "$CONTAINER_NAME"
  assert_success
}

@test "container: Start with an old forward incompatible runtime" {
  create_distro_container fedora 34 fedora-toolbox-34

  run container_started fedora-toolbox-34
  assert_success
}

@test "container(Fedora Rawhide): Containers with supported versions start without issues" {
  if ! is_fedora_rawhide; then
    skip "This test is only for Fedora Rawhide"
  fi

  local system_id
  system_id="$(get_system_id)"

  local system_version
  system_version="$(get_system_version)"

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

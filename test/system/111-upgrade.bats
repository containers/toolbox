# shellcheck shell=bats
#
# Copyright © 2025 Hadi Chokr <hadichokr@icloud.com>
# Licensed under the Apache License, Version 2.0 (the "License");
# You may not use this file except in compliance with the License.
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
  bats_require_minimum_version 1.8.0
  cleanup_all
}

teardown() {
  cleanup_all
}

@test "upgrade(Arch): Upgrade Arch container" {
  create_distro_container arch latest arch-toolbox-test
  run container_started arch-toolbox-test
  assert_success

  run "$TOOLBX" upgrade --container arch-toolbox-test
  assert_success
}

@test "upgrade(All): Upgrade all containers" {
  create_distro_container arch latest arch-toolbox-all-test
  create_distro_container ubuntu latest ubuntu-toolbox-all-test

  run container_started arch-toolbox-all-test
  assert_success
  run container_started ubuntu-toolbox-all-test
  assert_success

  run "$TOOLBX" upgrade --all
  assert_success
}

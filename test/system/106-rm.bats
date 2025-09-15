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


@test "rm: Try to remove a non-existent container" {
  container_name="nonexistentcontainer"
  run "$TOOLBX" rm "$container_name"

  #assert_failure  #BUG: it should return 1
  assert_output "Error: failed to inspect container $container_name"
}

@test "rm: Try to remove a running container" {
  skip "Bug: Fail in 'toolbox rm' does not return non-zero value"
  create_container running
  start_container running

  run "$TOOLBX" rm running

  #assert_failure  #BUG: it should return 1
  assert_output "Error: container running is running"
}

@test "rm: Remove a not running container" {
  create_container not-running

  run "$TOOLBX" rm not-running

  assert_success
  assert_output ""
}

@test "rm: Force remove a running container" {
  create_container running
  start_container running

  run "$TOOLBX" rm --force running

  assert_success
  assert_output ""
}

@test "rm: Force remove all containers (with 2 containers created and 1 running)" {
  num_of_containers="$(list_containers)"
  assert_equal "$num_of_containers" 0

  create_container running
  create_container not-running
  start_container running

  run "$TOOLBX" rm --force --all

  assert_success
  assert_output ""

  new_num_of_containers="$(list_containers)"

  assert_equal "$new_num_of_containers" "$num_of_containers"
}

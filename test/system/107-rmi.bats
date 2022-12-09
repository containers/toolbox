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
  cleanup_all
}

teardown() {
  cleanup_all
}


@test "rmi: Remove all images with the default image present" {
  num_of_images=$(list_images)
  assert_equal "$num_of_images" 0

  pull_default_image

  run --keep-empty-lines $TOOLBOX rmi --all

  assert_success
  assert_output ""

  new_num_of_images=$(list_images)

  assert_equal "$new_num_of_images" "$num_of_images"
}

@test "rmi: Try to remove all images with a container present and running" {
  skip "Bug: Fail in 'toolbox rmi' does not return non-zero value"
  num_of_images=$(list_images)
  assert_equal "$num_of_images" 0

  create_container foo
  start_container foo

  run --keep-empty-lines --separate-stderr $TOOLBOX rmi --all

  assert_failure
  lines=("${stderr_lines[@]}")
  assert_line --index 0 --regexp "Error: image .* has dependent children"

  new_num_of_images=$(list_images)

  assert_equal "$new_num_of_images" "$num_of_images"
}

@test "rmi: Force remove all images with a container present and running" {
  num_of_images=$(list_images)
  assert_equal "$num_of_images" 0

  create_container foo
  start_container foo

  run --keep-empty-lines $TOOLBOX rmi --all --force

  assert_success
  assert_output ""

  new_num_of_images=$(list_images)

  assert_equal "$new_num_of_images" "$num_of_images"
}

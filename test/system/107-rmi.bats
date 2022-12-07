#!/usr/bin/env bats
#
# Copyright © 2021 – 2023 Red Hat, Inc.
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
  bats_require_minimum_version 1.5.0
  _setup_environment
  cleanup_all
}

teardown() {
  cleanup_all
}


@test "rmi: --all without any images" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi --all

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: --all with the default image" {
  num_of_images=$(list_images)
  assert_equal "$num_of_images" 0

  pull_default_image

  run --keep-empty-lines $TOOLBOX rmi --all

  assert_success
  assert_output ""

  new_num_of_images=$(list_images)

  assert_equal "$new_num_of_images" "$num_of_images"
}

@test "rmi: An image by name" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  local default_image
  default_image="$(get_default_image)"

  pull_default_image

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$default_image"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: --all with an image without a name" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  build_image_without_name >/dev/null

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi --all

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: An image without a name" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  image="$(build_image_without_name)"

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$image"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: An image and its copy by name, separately" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  local default_image
  default_image="$(get_default_image)"

  pull_default_image_and_copy

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 2

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$default_image"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$default_image-copy"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: An image and its copy by name, separately (reverse order)" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  local default_image
  default_image="$(get_default_image)"

  pull_default_image_and_copy

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 2

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$default_image-copy"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$default_image"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: An image and its copy by name, together" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  local default_image
  default_image="$(get_default_image)"

  pull_default_image_and_copy

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 2

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$default_image" "$default_image-copy"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: An image and its copy by name, together (reverse order)" {
  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0

  local default_image
  default_image="$(get_default_image)"

  pull_default_image_and_copy

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 2

  run --keep-empty-lines --separate-stderr "$TOOLBOX" rmi "$default_image-copy" "$default_image"

  assert_success
  assert_output ""
  output="$stderr"
  assert_output ""
  if check_bats_version 1.7.0; then
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 0
}

@test "rmi: Try --all with a running container" {
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

@test "rmi: '--all --force' with a running container" {
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

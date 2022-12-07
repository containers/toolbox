#!/usr/bin/env bats
#
# Copyright © 2019 – 2022 Red Hat, Inc.
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


@test "list: Run 'list' with zero containers and zero images (the list should be empty)" {
  run --keep-empty-lines $TOOLBOX list

  assert_success
  assert_output ""
}

@test "list: Run 'list -c' with zero containers (the list should be empty)" {
  run --keep-empty-lines $TOOLBOX list -c

  assert_success
  assert_output ""
}

@test "list: Run 'list -i' with zero images (the list should be empty)" {
  run --keep-empty-lines $TOOLBOX list -i

  assert_success
  assert_output ""
}

@test "list: Run 'list' with zero toolbox's containers and images, but other image (the list should be empty)" {
  pull_distro_image busybox

  run podman images

  assert_output --partial "$BUSYBOX_IMAGE"

  run --keep-empty-lines $TOOLBOX list

  assert_success
  assert_output ""
}

@test "list: List an image without a name" {
  build_image_without_name >/dev/null

  run --keep-empty-lines --separate-stderr $TOOLBOX list

  assert_success
  assert_line --index 1 --partial "<none>"
  assert [ ${#lines[@]} -eq 3 ]
  if check_bats_version 1.7.0; then
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

@test "list: Image and its copy" {
  local default_image
  default_image="$(get_default_image)"

  pull_default_image_and_copy

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 2

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "$default_image"
  assert_line --index 2 --partial "$default_image-copy"
  assert [ ${#lines[@]} -eq 4 ]
  if check_bats_version 1.7.0; then
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

@test "list: Try to list images and containers (no flag) with 3 containers and 2 images (the list should have 3 images and 2 containers)" {
  # Pull the two images
  pull_default_image
  pull_distro_image fedora 34

  # Create three containers
  create_default_container
  create_container non-default-one
  create_container non-default-two

  # Check images
  run --keep-empty-lines --separate-stderr $TOOLBOX list --images

  assert_success
  assert_line --index 1 --partial "fedora-toolbox:34"
  assert_line --index 2 --partial "$(get_system_id)-toolbox:$(get_system_version)"
  assert [ ${#lines[@]} -eq 4 ]
  if check_bats_version 1.7.0; then
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  # Check containers
  run --keep-empty-lines --separate-stderr $TOOLBOX list --containers

  assert_success
  assert_line --index 1 --partial "$(get_system_id)-toolbox-$(get_system_version)"
  assert_line --index 2 --partial "non-default-one"
  assert_line --index 3 --partial "non-default-two"
  assert [ ${#lines[@]} -eq 5 ]
  if check_bats_version 1.7.0; then
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  # Check all together
  run --keep-empty-lines --separate-stderr $TOOLBOX list

  assert_success
  assert_line --index 1 --partial "fedora-toolbox:34"
  assert_line --index 2 --partial "$(get_system_id)-toolbox:$(get_system_version)"
  assert_line --index 5 --partial "$(get_system_id)-toolbox-$(get_system_version)"
  assert_line --index 6 --partial "non-default-one"
  assert_line --index 7 --partial "non-default-two"
  assert [ ${#lines[@]} -eq 9 ]
  if check_bats_version 1.7.0; then
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

@test "list: Images with and without names" {
  local default_image
  default_image="$(get_default_image)"

  pull_default_image
  pull_distro_image fedora 34
  build_image_without_name >/dev/null

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "<none>"
  assert_line --index 2 --partial "fedora-toolbox:34"
  assert_line --index 3 --partial "$default_image"
  assert [ ${#lines[@]} -eq 5 ]
  if check_bats_version 1.7.0; then
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

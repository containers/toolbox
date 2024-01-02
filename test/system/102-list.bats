# shellcheck shell=bats
#
# Copyright © 2019 – 2023 Red Hat, Inc.
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
  bats_require_minimum_version 1.7.0
  _setup_environment
  cleanup_all
}

teardown() {
  cleanup_all
}

@test "list: Smoke test" {
  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Smoke test (using --containers)" {
  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --containers

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Smoke test (using --images)" {
  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: With just one non-Toolbx image" {
  pull_distro_image busybox

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: With just one non-Toolbx image (using --images)" {
  pull_distro_image busybox

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Default image" {
  local default_image
  default_image="$(get_default_image)"

  pull_default_image

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "$default_image"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Default image (using --images)" {
  local default_image
  default_image="$(get_default_image)"

  pull_default_image

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "$default_image"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Arch Linux image" {
  pull_distro_image arch latest

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/arch-toolbox:latest"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Arch Linux image (using --images)" {
  pull_distro_image arch latest

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/arch-toolbox:latest"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Fedora 34 image" {
  pull_distro_image fedora 34

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "registry.fedoraproject.org/fedora-toolbox:34"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Fedora 34 image (using --images)" {
  pull_distro_image fedora 34

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "registry.fedoraproject.org/fedora-toolbox:34"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: RHEL 8.7 image" {
  pull_distro_image rhel 8.7

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "registry.access.redhat.com/ubi8/toolbox:8.7"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: RHEL 8.7 image (using --images)" {
  pull_distro_image rhel 8.7

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "registry.access.redhat.com/ubi8/toolbox:8.7"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Ubuntu 16.04 image" {
  pull_distro_image ubuntu 16.04

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/ubuntu-toolbox:16.04"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Ubuntu 16.04 image (using --images)" {
  pull_distro_image ubuntu 16.04

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/ubuntu-toolbox:16.04"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Ubuntu 18.04 image" {
  pull_distro_image ubuntu 18.04

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/ubuntu-toolbox:18.04"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Ubuntu 18.04 image (using --images)" {
  pull_distro_image ubuntu 18.04

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/ubuntu-toolbox:18.04"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Ubuntu 20.04 image" {
  pull_distro_image ubuntu 20.04

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/ubuntu-toolbox:20.04"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Ubuntu 20.04 image (using --images)" {
  pull_distro_image ubuntu 20.04

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "quay.io/toolbx/ubuntu-toolbox:20.04"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: An image without a name" {
  build_image_without_name >/dev/null

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "<none>"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: An image without a name (using --images)" {
  build_image_without_name >/dev/null

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "<none>"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 2 ]
  else
    assert [ ${#lines[@]} -eq 3 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
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

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 3 ]
  else
    assert [ ${#lines[@]} -eq 4 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Image and its copy (using --images)" {
  local default_image
  default_image="$(get_default_image)"

  pull_default_image_and_copy

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 2

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert_line --index 1 --partial "$default_image"
  assert_line --index 2 --partial "$default_image-copy"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 3 ]
  else
    assert [ ${#lines[@]} -eq 4 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Containers and images" {
  local default_image
  local fedora34image
  default_image="$(get_default_image)"
  fedora34image="registry.fedoraproject.org/fedora-toolbox:34"

  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  # Pull the two images
  pull_default_image
  pull_distro_image fedora 34

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 2

  # Create three containers
  create_default_container
  create_container non-default-one
  create_container non-default-two

  # Check images
  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  # Alphabetize images
  local image1
  local image2
  image1="$(echo -e "$fedora34image\n$default_image" | sort | head -n 1)"
  image2="$(echo -e "$fedora34image\n$default_image" | sort | tail -n 1)"

  assert_success
  assert_line --index 1 --partial "$image1"
  assert_line --index 2 --partial "$image2"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 3 ]
  else
    assert [ ${#lines[@]} -eq 4 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]

  # Check containers
  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --containers

  # Alphabetize containers
  local container1
  local container2
  local container3
  container1="$(echo -e "$default_container\nnon-default-one\nnon-default-two" | sort | head -n 1)"
  container2="$(echo -e "$default_container\nnon-default-one\nnon-default-two" | sort | head -n 2 | tail -n 1)"
  container3="$(echo -e "$default_container\nnon-default-one\nnon-default-two" | sort | tail -n 1)"

  assert_success
  assert_line --index 1 --partial "$container1"
  assert_line --index 2 --partial "$container2"
  assert_line --index 3 --partial "$container3"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 4 ]
  else
    assert [ ${#lines[@]} -eq 5 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]

  # Check all together
  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert_line --index 1 --partial "$image1"
  assert_line --index 2 --partial "$image2"
  assert_line --index 5 --partial "$container1"
  assert_line --index 6 --partial "$container2"
  assert_line --index 7 --partial "$container3"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 8 ]
  else
    assert [ ${#lines[@]} -eq 9 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Images with and without names" {
  local default_image
  local fedora34image
  default_image="$(get_default_image)"
  fedora34image="registry.fedoraproject.org/fedora-toolbox:34"

  pull_default_image
  pull_distro_image fedora 34
  build_image_without_name >/dev/null

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 3

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  # Alphabetize images
  local image1
  local image2
  image1="$(echo -e "$fedora34image\n$default_image" | sort | head -n 1)"
  image2="$(echo -e "$fedora34image\n$default_image" | sort | tail -n 1)"

  assert_success
  assert_line --index 1 --partial "<none>"
  assert_line --index 2 --partial "$image1"
  assert_line --index 3 --partial "$image2"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 4 ]
  else
    assert [ ${#lines[@]} -eq 5 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: Images with and without names (using --images)" {
  local default_image
  local fedora34image
  default_image="$(get_default_image)"
  fedora34image="registry.fedoraproject.org/fedora-toolbox:34"

  pull_default_image
  pull_distro_image fedora 34
  build_image_without_name >/dev/null

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  # Alphabetize images
  local image1
  local image2
  image1="$(echo -e "$fedora34image\n$default_image" | sort | head -n 1)"
  image2="$(echo -e "$fedora34image\n$default_image" | sort | tail -n 1)"

  assert_success
  assert_line --index 1 --partial "<none>"
  assert_line --index 2 --partial "$image1"
  assert_line --index 3 --partial "$image2"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 4 ]
  else
    assert [ ${#lines[@]} -eq 5 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: With just one non-Toolbx container and one non-Toolbx image" {
  local busybox_image
  busybox_image="$(get_busybox_image)"

  pull_distro_image busybox

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  $PODMAN create --name busybox-container "$busybox_image"

  local num_of_containers
  num_of_containers="$(list_containers)"
  assert_equal "$num_of_containers" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: With just one non-Toolbx container and one non-Toolbx image (using --containers)" {
  local busybox_image
  busybox_image="$(get_busybox_image)"

  pull_distro_image busybox

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  $PODMAN create --name busybox-container "$busybox_image"

  local num_of_containers
  num_of_containers="$(list_containers)"
  assert_equal "$num_of_containers" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --containers

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "list: With just one non-Toolbx container and one non-Toolbx image (using --images)" {
  local busybox_image
  busybox_image="$(get_busybox_image)"

  pull_distro_image busybox

  local num_of_images
  num_of_images="$(list_images)"
  assert_equal "$num_of_images" 1

  $PODMAN create --name busybox-container "$busybox_image"

  local num_of_containers
  num_of_containers="$(list_containers)"
  assert_equal "$num_of_containers" 1

  run --keep-empty-lines --separate-stderr "$TOOLBOX" list --images

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

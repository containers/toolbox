#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  cleanup_all
}

teardown() {
  cleanup_all
}


@test "list: Run 'list' with zero containers and zero images (the list should be empty)" {
  run $TOOLBOX list

  assert_success
  assert_output ""
}

@test "list: Run 'list -c' with zero containers (the list should be empty)" {
  run $TOOLBOX list -c

  assert_success
  assert_output ""
}

@test "list: Run 'list -i' with zero images (the list should be empty)" {
  run $TOOLBOX list -c

  assert_success
  assert_output ""
}

@test "list: Run 'list' with zero toolbox's containers and images, but other image (the list should be empty)" {
  pull_distro_image busybox

  run podman images

  assert_output --partial "$BUSYBOX_IMAGE"

  run $TOOLBOX list

  assert_success
  assert_output ""
}

@test "list: Try to list images and containers (no flag) with 3 containers and 2 images (the list should have 3 images and 2 containers)" {
  # Pull the two images
  pull_default_image
  pull_distro_image fedora 32

  # Create three containers
  create_default_container
  create_container non-default-one
  create_container non-default-two

  # Check images
  run $TOOLBOX list --images

  assert_success
  assert_output --partial "$(get_system_id)-toolbox:$(get_system_version)"
  assert_output --partial "fedora-toolbox:32"

  # Check containers
  run $TOOLBOX list --containers

  assert_success
  assert_output --partial "$(get_system_id)-toolbox-$(get_system_version)"
  assert_output --partial "non-default-one"
  assert_output --partial "non-default-two"

  # Check all together
  run $TOOLBOX list

  assert_success
  assert_output --partial "$(get_system_id)-toolbox:$(get_system_version)"
  assert_output --partial "fedora-toolbox:32"
  assert_output --partial "$(get_system_id)-toolbox-$(get_system_version)"
  assert_output --partial "non-default-one"
  assert_output --partial "non-default-two"
}

@test "list: Run 'list -i' with UBI image (8.4; public) present" {
  pull_distro_image rhel 8.4

  run toolbox list --images

  assert_success
  assert_output --partial "registry.access.redhat.com/ubi8/ubi:8.4"
}

@test "list: Run 'list' with UBI image (8.4; public), toolbox container and non-toolbox container" {
  local num_of_containers

  pull_distro_image rhel 8.4

  create_distro_container rhel 8.4 rhel-toolbox
  podman create --name podman-container ubi8/ubi:8.4 /bin/sh

  num_of_containers=$(list_containers)
  assert [ $num_of_containers -eq 2 ]

  run toolbox list

  assert_success
  assert_line --index 1 --partial "registry.access.redhat.com/ubi8/ubi:8.4"
  assert_line --index 3 --partial "rhel-toolbox"
  refute_output --partial "podman-container"
}

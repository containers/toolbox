#!/usr/bin/env bats

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
  run $TOOLBOX list -i

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
  assert_output --partial "$(toolbx_default_image_name)"
  assert_output --partial "fedora-toolbox:32"

  # Check containers
  run $TOOLBOX list --containers

  assert_success
  assert_output --partial "$(toolbx_default_container_name)"
  assert_output --partial "non-default-one"
  assert_output --partial "non-default-two"

  # Check all together
  run $TOOLBOX list

  assert_success
  assert_output --partial "$(toolbx_default_image_name)"
  assert_output --partial "fedora-toolbox:32"
  assert_output --partial "$(toolbx_default_container_name)"
  assert_output --partial "non-default-one"
  assert_output --partial "non-default-two"
}

@test "list: List an image without a name" {
    echo -e "FROM scratch\n\nLABEL com.github.containers.toolbox=\"true\"" > "$BATS_TMPDIR"/Containerfile

    run $PODMAN build "$BATS_TMPDIR"

    assert_success
    assert_line --index 0 --partial "FROM scratch"
    assert_line --index 1 --partial "LABEL com.github.containers.toolbox=\"true\""
    assert_line --index 2 --partial "COMMIT"
    assert_line --index 3 --regexp "^--> [a-z0-9]*$"

    run $TOOLBOX list

    assert_success
    assert_line --index 1 --partial "<none>"

    rm -f "$BATS_TMPDIR"/Containerfile
}

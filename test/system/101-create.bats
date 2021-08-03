#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  check_xdg_runtime_dir
  cleanup_containers
}

teardown() {
  cleanup_containers
}


@test "create: Create the default container" {
  pull_default_image

  run $TOOLBOX -y create

  assert_success
}

@test "create: Create a container with a valid custom name ('custom-containerName')" {
  pull_default_image

  run $TOOLBOX -y create -c "custom-containerName"

  assert_success
}

@test "create: Create a container with a custom image and name ('fedora32'; f32)" {
  pull_distro_image fedora 32

  run $TOOLBOX -y create -c "fedora32" -i fedora-toolbox:32

  assert_success
}

@test "create: Try to create a container with invalid custom name ('ßpeci@l.Nam€'; using positional argument)" {
  run $TOOLBOX -y create "ßpeci@l.Nam€"

  assert_failure
  assert_line --index 0 "Error: invalid argument for 'CONTAINER'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'"
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Try to create a container with invalid custom name ('ßpeci@l.Nam€'; using option --container)" {
  run $TOOLBOX -y create -c "ßpeci@l.Nam€"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--container'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'"
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Create a container with a distro and release options ('fedora'; f32)" {
  pull_distro_image fedora 32

  run $TOOLBOX -y create -d "fedora" -r f32

  assert_success
  assert_output --partial "Created container: fedora-toolbox-32"
  assert_output --partial "Enter with: toolbox enter fedora-toolbox-32"

  # Make sure the container has actually been created
  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+fedora-toolbox-32"
}

@test "create: Try to create a container based on non-existent image" {
  run $TOOLBOX -y create -i foo.org/bar

  assert_failure
  assert_line --index 0 "Error: failed to pull image foo.org/bar"
  assert_line --index 1 "If it was a private image, log in with: podman login foo.org"
  assert_line --index 2 "Use 'toolbox --verbose ...' for further details."
}

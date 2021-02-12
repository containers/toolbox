#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  cleanup_containers
}

teardown() {
  cleanup_containers
}


@test "create: Create the default container" {
  pull_default_image

  run toolbox -y create

  assert_success
}

@test "create: Create a container with a valid custom name ('custom-containerName')" {
  run toolbox -y create -c "custom-containerName"

  assert_success
}

@test "create: Create a container with a custom image and name ('fedora29'; f29)" {
  pull_image_old 29

  run toolbox -y create -c "fedora29" -i fedora-toolbox:29

  assert_success
}

@test "create: Try to create a container with invalid custom name ('ßpeci@l.Nam€'; using positional argument)" {
  run toolbox -y create "ßpeci@l.Nam€"

  assert_failure
  assert_line --index 0 "Error: invalid argument for 'CONTAINER'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'"
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Try to create a container with invalid custom name ('ßpeci@l.Nam€'; using option --container)" {
  run toolbox -y create -c "ßpeci@l.Nam€"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--container'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'"
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Create a container with a distro and release options ('fedora'; f29)" {
  pull_image 29

  run toolbox -y create -d "fedora" -r f29

  assert_success
  assert_output --partial "Created container: fedora-toolbox-29"
  assert_output --partial "Enter with: toolbox enter --release 29"

  # Make sure the container has actually been created
  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+fedora-toolbox-29"
}

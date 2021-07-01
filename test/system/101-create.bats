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
  assert_output --partial "Enter with: toolbox enter --release 32"

  # Make sure the container has actually been created
  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+fedora-toolbox-32"
}

@test "create: Try to create container based on non-existent image" {
  run $TOOLBOX -y create -i localhost:50000/unknownimage

  assert_failure
  assert_line --index 0 "Error: image localhost:50000/unknownimage does not exist"
  assert_line --index 1 "Make sure the image URI is correct."
}

@test "create: Try to create container based on unresolvable image" {
  run $TOOLBOX create -i foobar

  assert_failure
  assert_line --index 0 "Error: image foobar not found in local storage and known registries"
  assert_line --index 1 "Make sure the image URI is correct"
}

@test "create: Try to create container based on image from private registry" {
  run $TOOLBOX -y create -i localhost:50001/fedora-toolbox:32

  assert_failure
  assert_line --index 0 "Error: Could not pull image localhost:50001/fedora-toolbox:32"
  assert_line --index 1 "The registry requires logging in."
  assert_line --index 2 "See 'podman login --help' on how to login into a registry."
}

@test "create: Try to create container based on image from private registry (provide false credentials)" {
  run $TOOLBOX -y create -i localhost:50001/fedora-toolbox:32 --creds wrong:wrong

  assert_failure
  assert_line --index 0 "Could not pull image localhost:50001/fedora-toolbox:32"
  assert_line --index 1 "The registry requires logging in."
  assert_line --index 2 "Credentials were provided. Trying to log into localhost:50001"
  assert_line --index 3 --partial "Error: error logging into \"localhost:50001\""
}

@test "create: Create container based on image from private registry (provide correct credentials)" {
  run $TOOLBOX -y create my-fedora -i localhost:50001/fedora-toolbox:32 --creds user:user

  assert_success
  assert_line --index 0 "Could not pull image localhost:50001/fedora-toolbox:32"
  assert_line --index 1 "The registry requires logging in."
  assert_line --index 2 "Credentials were provided. Trying to log into localhost:50001"
  assert_line --index 3 "Login Succeeded!"
  assert_line --index 4 "Retrying to pull image localhost:50001/fedora-toolbox:32"
  assert_line --index 5 "Created container: my-fedora"
  assert_line --index 6 "Enter with: toolbox enter my-fedora"
}

@test "create: Create container based on image from private registry (log in beforehand)" {
  
  run $PODMAN login localhost:50001 --username user --password user

  assert_success

  run $TOOLBOX -y create my-fedora -i localhost:50001/fedora-toolbox:32

  assert_success
  assert_line --index 0 "Created container: my-fedora"
  assert_line --index 1 "Enter with: toolbox enter my-fedora"
}

#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  _setup_environment
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

@test "create: Try to create a container based on unsupported distribution" {
  local distro="foo"

  run $TOOLBOX -y create -d "$distro"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--distro'"
  # Distro names are in a hashtable and thus the order can change
  assert_line --index 1 --regexp "Supported values are: (.?(fedora|rhel))+"
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try to create a container based on Fedora but with wrong version" {
  run $TOOLBOX -y create -d fedora -r foobar

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "Supported values for distribution fedora are in format: <release>/f<release>"
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try to create a container based on non-default distribution without providing version" {
  local distro="fedora"
  local system_id="$(get_system_id)"

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run $TOOLBOX -y create -d "$distro"

  assert_failure
  assert_line --index 0 "Error: release not found for non-default distribution $distro"
  assert [ ${#lines[@]} -eq 1 ]
}

@test "create: Try to create a container and pass a non-existent file to the --authfile option" {
  local file="$BATS_RUN_TMPDIR/non-existent-file"

  run $TOOLBOX create --authfile "$file"

  assert_failure
  assert_output "Error: file $file not found"
}

@test "create: Create a container based on an image from locked registry using an authentication file" {
  local authfile="$BATS_RUN_TMPDIR/authfile"
  local image="fedora-toolbox:32"

  run $PODMAN login --authfile "$authfile" --username user --password user "$DOCKER_REG_URI"
  assert_success

  run $TOOLBOX --assumeyes create --image "$DOCKER_REG_URI/$image"

  assert_failure
  assert_line --index 0 "Error: failed to pull image $DOCKER_REG_URI/$image"
  assert_line --index 1 "If it was a private image, log in with: podman login $DOCKER_REG_URI"
  assert_line --index 2 "Use 'toolbox --verbose ...' for further details."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX --assumeyes create --authfile "$authfile" --image "$DOCKER_REG_URI/$image"

  rm "$authfile"

  assert_success
  assert_line --index 0 "Created container: fedora-toolbox-32"
  assert_line --index 1 "Enter with: toolbox enter fedora-toolbox-32"
  assert [ ${#lines[@]} -eq 2 ]
}

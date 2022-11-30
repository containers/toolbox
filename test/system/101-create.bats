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
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Try to create a container with invalid custom name ('ßpeci@l.Nam€'; using option --container)" {
  run $TOOLBOX -y create -c "ßpeci@l.Nam€"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--container'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Try to create a container with invalid custom image ('ßpeci@l.Nam€')" {
  local image="ßpeci@l.Nam€"

  run $TOOLBOX create --image "$image"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--image'"
  assert_line --index 1 "Container name $image generated from image is invalid."
  assert_line --index 2 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 3 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 4 ]
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

@test "create: Try to create a container based on unsupported distribution" {
  local distro="foo"

  run $TOOLBOX --assumeyes create --distro "$distro"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try to create a container based on non-existent image" {
  run $TOOLBOX -y create -i foo.org/bar

  assert_failure
  assert_line --index 0 "Error: failed to pull image foo.org/bar"
  assert_line --index 1 "If it was a private image, log in with: podman login foo.org"
  assert_line --index 2 "Use 'toolbox --verbose ...' for further details."
}

@test "create: Try to create a container based on Fedora but with wrong version" {
  run $TOOLBOX -y create -d fedora -r foobar

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX --assumeyes create --distro fedora --release -3

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try to create a container based on RHEL but with wrong version" {
  run $TOOLBOX --assumeyes create --distro rhel --release 8

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX --assumeyes create --distro rhel --release 8.2foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX --assumeyes create --distro rhel --release -2.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive number."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try to create a container based on non-default distribution without providing version" {
  local distro="fedora"
  local system_id="$(get_system_id)"

  assert [ -n "$system_id" ]

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run $TOOLBOX -y create -d "$distro"

  assert_failure
  assert_line --index 0 "Error: option '--release' is needed"
  assert_line --index 1 "Distribution $distro doesn't match the host."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try to create a container using both --distro and --image" {
  pull_distro_image fedora 32

  run $TOOLBOX --assumeyes create --distro "fedora" --image fedora-toolbox:32

  assert_failure
  assert_line --index 0 "Error: options --distro and --image cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 2 ]
}

@test "create: Try to create a container using both --image and --release" {
  pull_distro_image fedora 32

  run $TOOLBOX --assumeyes create --image fedora-toolbox:32 --release 32

  assert_failure
  assert_line --index 0 "Error: options --image and --release cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 2 ]
}

@test "create: Try to create a container and pass a non-existent file to the --authfile option" {
  local file="$BATS_RUN_TMPDIR/non-existent-file"

  run $TOOLBOX create --authfile "$file"

  assert_failure
  assert_line --index 0 "Error: file $file not found"
  assert_line --index 1 "'podman login' can be used to create the file."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
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

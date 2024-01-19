# shellcheck shell=bats
#
# Copyright © 2019 – 2024 Red Hat, Inc.
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
  cleanup_all
}

@test "create: Smoke test" {
  pull_default_image

  run "$TOOLBOX" --assumeyes create

  assert_success
}

@test "create: With a custom name (using option --container)" {
  pull_default_image

  run "$TOOLBOX" --assumeyes create --container "custom-containerName"

  assert_success
}

@test "create: With a custom image and name (using option --container)" {
  pull_distro_image fedora 34

  run "$TOOLBOX" --assumeyes create --container "fedora34" --image fedora-toolbox:34

  assert_success
}

@test "create: Try without --assumeyes" {
  run --keep-empty-lines --separate-stderr "$TOOLBOX" create

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: image required to create toolbox container."
  assert_line --index 1 "Use option '--assumeyes' to download the image."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try with an invalid custom name (using positional argument)" {
  run "$TOOLBOX" --assumeyes create "ßpeci@l.N@m€"

  assert_failure
  assert_line --index 0 "Error: invalid argument for 'CONTAINER'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Try with an invalid custom name (using option --container)" {
  run "$TOOLBOX" --assumeyes create --container "ßpeci@l.N@m€"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--container'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "create: Try with an invalid custom image" {
  local image="ßpeci@l.N@m€"

  run "$TOOLBOX" create --image "$image"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--image'"
  assert_line --index 1 "Container name $image generated from image is invalid."
  assert_line --index 2 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 3 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 4 ]
}

@test "create: Try with an invalid custom image (using --assumeyes)" {
  local image="ßpeci@l.N@m€"

  run --keep-empty-lines --separate-stderr "$TOOLBOX" --assumeyes create --image "$image"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--image'"
  assert_line --index 1 "Container name $image generated from image is invalid."
  assert_line --index 2 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 3 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 4 ]
}

@test "create: Arch Linux" {
  pull_distro_image arch latest

  run "$TOOLBOX" --assumeyes create --distro arch

  assert_success
  assert_output --partial "Created container: arch-toolbox-latest"
  assert_output --partial "Enter with: toolbox enter arch-toolbox-latest"

  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+arch-toolbox-latest"
}

@test "create: Arch Linux ('--release latest')" {
  pull_distro_image arch latest

  run "$TOOLBOX" --assumeyes create --distro arch --release latest

  assert_success
  assert_output --partial "Created container: arch-toolbox-latest"
  assert_output --partial "Enter with: toolbox enter arch-toolbox-latest"

  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+arch-toolbox-latest"
}

@test "create: Arch Linux ('--release rolling')" {
  pull_distro_image arch latest

  run "$TOOLBOX" --assumeyes create --distro arch --release rolling

  assert_success
  assert_output --partial "Created container: arch-toolbox-latest"
  assert_output --partial "Enter with: toolbox enter arch-toolbox-latest"

  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+arch-toolbox-latest"
}

@test "create: Fedora 34" {
  pull_distro_image fedora 34

  run "$TOOLBOX" --assumeyes create --distro fedora --release f34

  assert_success
  assert_output --partial "Created container: fedora-toolbox-34"
  assert_output --partial "Enter with: toolbox enter fedora-toolbox-34"

  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+fedora-toolbox-34"
}

@test "create: RHEL 8.9" {
  pull_distro_image rhel 8.9

  run "$TOOLBOX" --assumeyes create --distro rhel --release 8.9

  assert_success
  assert_output --partial "Created container: rhel-toolbox-8.9"
  assert_output --partial "Enter with: toolbox enter rhel-toolbox-8.9"

  run podman ps -a

  assert_output --regexp "Created[[:blank:]]+rhel-toolbox-8.9"
}

@test "create: Ubuntu 16.04" {
  pull_distro_image ubuntu 16.04

  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 16.04

  assert_success
  assert_output --partial "Created container: ubuntu-toolbox-16.04"
  assert_output --partial "Enter with: toolbox enter ubuntu-toolbox-16.04"

  run $PODMAN ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+ubuntu-toolbox-16.04"
}

@test "create: Ubuntu 18.04" {
  pull_distro_image ubuntu 18.04

  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 18.04

  assert_success
  assert_output --partial "Created container: ubuntu-toolbox-18.04"
  assert_output --partial "Enter with: toolbox enter ubuntu-toolbox-18.04"

  run $PODMAN ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+ubuntu-toolbox-18.04"
}

@test "create: Ubuntu 20.04" {
  pull_distro_image ubuntu 20.04

  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 20.04

  assert_success
  assert_output --partial "Created container: ubuntu-toolbox-20.04"
  assert_output --partial "Enter with: toolbox enter ubuntu-toolbox-20.04"

  run $PODMAN ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+ubuntu-toolbox-20.04"
}

@test "create: Try an unsupported distribution" {
  local distro="foo"

  run "$TOOLBOX" --assumeyes create --distro "$distro"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try a non-existent image" {
  run --keep-empty-lines --separate-stderr "$TOOLBOX" create --image foo.org/bar

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: image required to create toolbox container."
  assert_line --index 1 "Use option '--assumeyes' to download the image."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try a non-existent image (using --assumeyes)" {
  run "$TOOLBOX" --assumeyes create --image foo.org/bar

  assert_failure
  assert_line --index 0 "Error: failed to pull image foo.org/bar"
  assert_line --index 1 "If it was a private image, log in with: podman login foo.org"
  assert_line --index 2 "Use 'toolbox --verbose ...' for further details."
}

@test "create: Try Arch Linux with an invalid release ('--release foo')" {
  run "$TOOLBOX" --assumeyes create --distro arch --release foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be 'latest'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release -3')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release -3

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release -3.0')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release -3.0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release -3.1')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release -3.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 0')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release 0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 0.0')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release 0.0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 0.1')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release 0.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 3.0')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release 3.0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 3.1')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release 3.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release foo')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 3foo')" {
  run "$TOOLBOX" --assumeyes create --distro fedora --release 3foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release 8

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.0.0')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release 8.0.0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.0.1')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release 8.0.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.3.0')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release 8.3.0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.3.1')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release 8.3.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release foo')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.2foo')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release 8.2foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release -2.1')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release -2.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive number."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release -2.-1')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release -2.-1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 2.-1')" {
  run "$TOOLBOX" --assumeyes create --distro rhel --release 2.-1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 20

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20.04.0')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 20.04.0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20.04.1')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 20.04.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release foo')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20foo')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 20foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release foo.bar')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release foo.bar

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release foo.bar.baz')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release foo.bar.baz

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 3.10')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 3.10

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year must be 4 or more."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 202.4')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 202.4

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 202.04')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 202.04

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 2020.4')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 2020.4

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 2020.04')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 2020.04

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 04.10')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 04.10

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have a leading zero."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.bar')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 4.bar

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.bar.baz')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 4.bar.baz

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.0')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 4.0

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must be between 01 and 12."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.00')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 4.00

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must be between 01 and 12."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.13')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 4.13

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must be between 01 and 12."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20.4')" {
  run "$TOOLBOX" --assumeyes create --distro ubuntu --release 20.4

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must have two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try a non-default distro without a release" {
  local distro="fedora"

  local system_id
  system_id="$(get_system_id)"

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run "$TOOLBOX" --assumeyes create --distro "$distro"

  assert_failure
  assert_line --index 0 "Error: option '--release' is needed"
  assert_line --index 1 "Distribution $distro doesn't match the host."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try using both --distro and --image" {
  run --keep-empty-lines --separate-stderr "$TOOLBOX" create --distro fedora --image fedora-toolbox:34

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --distro and --image cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try using both --distro and --image (using --assumeyes)" {
  pull_distro_image fedora 34

  run "$TOOLBOX" --assumeyes create --distro fedora --image fedora-toolbox:34

  assert_failure
  assert_line --index 0 "Error: options --distro and --image cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 2 ]
}

@test "create: Try using both --image and --release" {
  run --keep-empty-lines --separate-stderr "$TOOLBOX" create --image fedora-toolbox:34 --release 34

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --image and --release cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 2 ]
}

@test "create: Try using both --image and --release (using --assumeyes)" {
  pull_distro_image fedora 34

  run "$TOOLBOX" --assumeyes create --image fedora-toolbox:34 --release 34

  assert_failure
  assert_line --index 0 "Error: options --image and --release cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 2 ]
}

@test "create: Try a non-existent authentication file" {
  local file="$BATS_RUN_TMPDIR/non-existent-file"

  run "$TOOLBOX" create --authfile "$file"

  assert_failure
  assert_line --index 0 "Error: file $file not found"
  assert_line --index 1 "'podman login' can be used to create the file."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try a non-existent authentication file (using --assumeyes)" {
  local file="$BATS_RUN_TMPDIR/non-existent-file"

  run --keep-empty-lines --separate-stderr "$TOOLBOX" --assumeyes create --authfile "$file"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: file $file not found"
  assert_line --index 1 "'podman login' can be used to create the file."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: With a custom image that needs an authentication file" {
  local authfile="$BATS_RUN_TMPDIR/authfile"
  local image="fedora-toolbox:34"

  run $PODMAN login --authfile "$authfile" --username user --password user "$DOCKER_REG_URI"
  assert_success

  run "$TOOLBOX" --assumeyes create --image "$DOCKER_REG_URI/$image"

  assert_failure
  assert_line --index 0 "Error: failed to pull image $DOCKER_REG_URI/$image"
  assert_line --index 1 "If it was a private image, log in with: podman login $DOCKER_REG_URI"
  assert_line --index 2 "Use 'toolbox --verbose ...' for further details."
  assert [ ${#lines[@]} -eq 3 ]

  run "$TOOLBOX" --assumeyes create --authfile "$authfile" --image "$DOCKER_REG_URI/$image"

  rm "$authfile"

  assert_success
  assert_line --index 0 "Created container: fedora-toolbox-34"
  assert_line --index 1 "Enter with: toolbox enter fedora-toolbox-34"
  assert [ ${#lines[@]} -eq 2 ]
}

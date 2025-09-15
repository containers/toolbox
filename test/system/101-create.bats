# shellcheck shell=bats
#
# Copyright © 2019 – 2025 Red Hat, Inc.
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

# bats file_tags=commands-options

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.10.0
  cleanup_all
}

teardown() {
  cleanup_all
}

@test "create: Smoke test" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  pull_default_image

  run --keep-empty-lines --separate-stderr "$TOOLBX" create

  assert_success
  assert_line --index 0 "Created container: $default_container"
  assert_line --index 1 "Enter with: toolbox enter"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$default_container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$default_container"

  assert_success
  assert_output "true"
}

@test "create: With a custom name (using option --container)" {
  pull_default_image

  local container="custom-containerName"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --container "$container"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: With a custom image and name (using option --container)" {
  pull_distro_image fedora 34

  local container="fedora34"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --container "$container" --image fedora-toolbox:34

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Try without --assumeyes" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" create

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: image required to create Toolbx container."
  assert_line --index 1 "Use option '--assumeyes' to download the image."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try with an invalid custom name (using positional argument)" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create "ßpeci@l.N@m€"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for 'CONTAINER'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try with an invalid custom name (using option --container)" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --container "ßpeci@l.N@m€"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--container'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try with an invalid custom image" {
  local image="ßpeci@l.N@m€"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --image "$image"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--image'"
  assert_line --index 1 "Container name $image generated from image is invalid."
  assert_line --index 2 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 3 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 4 ]
}

@test "create: Try with an invalid custom image (using --assumeyes)" {
  local image="ßpeci@l.N@m€"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --image "$image"

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
  local system_id
  system_id="$(get_system_id)"

  pull_distro_image arch latest

  local container="arch-toolbox-latest"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro arch

  assert_success
  assert_line --index 0 "Created container: $container"

  if [ "$system_id" = "arch" ]; then
    assert_line --index 1 "Enter with: toolbox enter"
  else
    assert_line --index 1 "Enter with: toolbox enter $container"
  fi

  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Arch Linux ('--release latest')" {
  local system_id
  system_id="$(get_system_id)"

  pull_distro_image arch latest

  local container="arch-toolbox-latest"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro arch --release latest

  assert_success
  assert_line --index 0 "Created container: $container"

  if [ "$system_id" = "arch" ]; then
    assert_line --index 1 "Enter with: toolbox enter"
  else
    assert_line --index 1 "Enter with: toolbox enter $container"
  fi

  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Arch Linux ('--release rolling')" {
  local system_id
  system_id="$(get_system_id)"

  pull_distro_image arch latest

  local container="arch-toolbox-latest"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro arch --release rolling

  assert_success
  assert_line --index 0 "Created container: $container"

  if [ "$system_id" = "arch" ]; then
    assert_line --index 1 "Enter with: toolbox enter"
  else
    assert_line --index 1 "Enter with: toolbox enter $container"
  fi

  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Fedora 34" {
  pull_distro_image fedora 34

  local container="fedora-toolbox-34"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro fedora --release f34

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: RHEL 8.10" {
  pull_distro_image rhel 8.10

  local container="rhel-toolbox-8.10"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro rhel --release 8.10

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Ubuntu 16.04" {
  pull_distro_image ubuntu 16.04

  local container="ubuntu-toolbox-16.04"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro ubuntu --release 16.04

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Ubuntu 18.04" {
  pull_distro_image ubuntu 18.04

  local container="ubuntu-toolbox-18.04"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro ubuntu --release 18.04

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Ubuntu 20.04" {
  pull_distro_image ubuntu 20.04

  local container="ubuntu-toolbox-20.04"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro ubuntu --release 20.04

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: With a custom image without a name" {
  image="$(build_image_without_name)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --image "$image"

  assert_success
  assert_line --index 0 "Created container: $image"
  assert_line --index 1 "Enter with: toolbox enter $image"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$image"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$image"

  assert_success
  assert_output "true"
}

@test "create: With a custom image without a name, and container name (using positional argument)" {
  image="$(build_image_without_name)"

  local container="non-default"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --image "$image" "$container"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: With a custom image without a name, and container name (using option --container)" {
  image="$(build_image_without_name)"

  local container="non-default"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --image "$image" --container "$container"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

@test "create: Try the same name again" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  pull_default_image

  run --keep-empty-lines --separate-stderr "$TOOLBX" create

  assert_success
  assert_line --index 0 "Created container: $default_container"
  assert_line --index 1 "Enter with: toolbox enter"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$default_container"

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" create

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container $default_container already exists"
  assert_line --index 1 "Enter with: toolbox enter"
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$default_container"
}

@test "create: Try an unsupported distribution" {
  local distro="foo"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro "$distro"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try a non-existent image" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" create --image foo.org/bar

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: image required to create Toolbx container."
  assert_line --index 1 "Use option '--assumeyes' to download the image."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try a non-existent image (using --assumeyes)" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --image foo.org/bar

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to pull image foo.org/bar"
  assert_line --index 1 "If it was a private image, log in with: podman login foo.org"
  assert_line --index 2 "Use 'toolbox --verbose ...' for further details."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Arch Linux with an invalid release ('--release foo')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro arch --release foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be 'latest'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release -3')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release -3

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release -3.0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release -3.0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release -3.1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release -3.1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release 0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 0.0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release 0.0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 0.1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release 0.1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 3.0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release 3.0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 3.1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release 3.1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release foo')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Fedora with an invalid release ('--release 3foo')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --release 3foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release 8

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.0.0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release 8.0.0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.0.1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release 8.0.1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.3.0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release 8.3.0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.3.1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release 8.3.1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release foo')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 8.2foo')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release 8.2foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release -2.1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release -2.1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive number."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release -2.-1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release -2.-1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try RHEL with an invalid release ('--release 2.-1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro rhel --release 2.-1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 20

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20.04.0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 20.04.0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20.04.1')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 20.04.1

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release foo')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20foo')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 20foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release foo.bar')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release foo.bar

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release foo.bar.baz')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release foo.bar.baz

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 3.10')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 3.10

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year must be 4 or more."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 202.4')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 202.4

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 202.04')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 202.04

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 2020.4')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 2020.4

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 2020.04')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 2020.04

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have more than two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 04.10')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 04.10

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release year cannot have a leading zero."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.bar')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 4.bar

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.bar.baz')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 4.bar.baz

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the 'YY.MM' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.0')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 4.0

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must be between 01 and 12."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.00')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 4.00

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must be between 01 and 12."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 4.13')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 4.13

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must be between 01 and 12."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try Ubuntu with an invalid release ('--release 20.4')" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro ubuntu --release 20.4

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release month must have two digits."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try a non-default distro without a release" {
  local distro="fedora"

  local system_id
  system_id="$(get_system_id)"

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro "$distro"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: option '--release' is needed"
  assert_line --index 1 "Distribution $distro doesn't match the host."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try using both --distro and --image" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" create --distro fedora --image fedora-toolbox:34

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --distro and --image cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try using both --distro and --image (using --assumeyes)" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --distro fedora --image fedora-toolbox:34

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --distro and --image cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try using both --image and --release" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" create --image fedora-toolbox:34 --release 34

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --image and --release cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try using both --image and --release (using --assumeyes)" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --image fedora-toolbox:34 --release 34

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --image and --release cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try a non-existent authentication file" {
  local file="$BATS_TEST_TMPDIR/non-existent-file"

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --authfile "$file"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: file $file not found"
  assert_line --index 1 "'podman login' can be used to create the file."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try a non-existent authentication file (using --assumeyes)" {
  local file="$BATS_TEST_TMPDIR/non-existent-file"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --authfile "$file"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: file $file not found"
  assert_line --index 1 "'podman login' can be used to create the file."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: With a custom image that needs an authentication file" {
  local authfile="$BATS_TEST_TMPDIR/authfile"
  local image="fedora-toolbox:34"

  run podman login --authfile "$authfile" --username user --password user "$DOCKER_REG_URI"
  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --image "$DOCKER_REG_URI/$image"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to pull image $DOCKER_REG_URI/$image"
  assert_line --index 1 "If it was a private image, log in with: podman login $DOCKER_REG_URI"
  assert_line --index 2 "Use 'toolbox --verbose ...' for further details."
  assert [ ${#stderr_lines[@]} -eq 3 ]

  local container="fedora-toolbox-34"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
                                             --authfile "$authfile" \
                                             --image "$DOCKER_REG_URI/$image"

  rm "$authfile"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "true"
}

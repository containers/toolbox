# shellcheck shell=bats
#
# Copyright © 2019 – 2026 Red Hat, Inc.
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

# bats file_tags=non-native

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup_file() {
  bats_require_minimum_version 1.10.0
  skip_if_no_cross_arch_support
}

setup() {
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

@test "create: Cross-arch Smoke test" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)-${cross_arch}"

  pull_default_image_cross_arch

  run --keep-empty-lines --separate-stderr "$TOOLBX" create --arch "${cross_arch}"

  assert_success
  assert_line --index 0 "Created container: $default_container"
  assert_line --index 1 "Enter with: toolbox enter $default_container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$default_container"

  run podman inspect \
        --format '{{index .Config.Labels "com.github.containers.toolbox"}} {{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$default_container"

  assert_success
  assert_output "true ${cross_arch}"
}

@test "create: Cross-arch with --arch and --distro/--release" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local container="fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: Cross-arch with --arch (binfmt alias)" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local container="fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${binfmt_arch}" \
    --distro fedora \
    --release 44

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: Cross-arch with --image (no --distro/--release)" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local container="fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --image "fedora-toolbox:44"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: Cross-arch inferred from image tag (no --arch)" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  local container="fedora-toolbox-44-${binfmt_arch}"

  # The image tag contains the architecture suffix — toolbox should infer
  # the architecture from it without needing --arch
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --image "fedora-toolbox:44-${binfmt_arch}"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: Cross-arch with --arch matching image tag suffix" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  local container="fedora-toolbox-44-${binfmt_arch}"

  # --arch matches the architecture in the image tag — no conflict
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --image "fedora-toolbox:44-${binfmt_arch}"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: Cross-arch with custom container name" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local container="my-cross-arch-container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --container "$container" \
    --distro fedora \
    --release 44

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: --arch matching host is treated as native" {
  local host_arch

  case "$(uname -m)" in
    x86_64)  host_arch="amd64" ;;
    aarch64) host_arch="arm64" ;;
    ppc64le) host_arch="ppc64le" ;;
  esac

  pull_default_image

  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${host_arch}"

  assert_success
  assert_line --index 0 "Created container: $default_container"
  assert_line --index 1 "Enter with: toolbox enter"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$default_container"

  assert_success
  assert_output "${host_arch}"
}

@test "create: --arch matching host with --image is native" {
  local host_arch

  case "$(uname -m)" in
    x86_64)  host_arch="amd64" ;;
    aarch64) host_arch="arm64" ;;
    ppc64le) host_arch="ppc64le" ;;
  esac

  pull_default_image

  local default_image
  default_image="$(get_default_image)"

  local container
  container="$(get_system_id)-toolbox-$(get_system_version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${host_arch}" \
    --image "$default_image"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${host_arch}"
}

@test "create: Image with arch-like suffix from unsupported distro is native" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local custom_image="localhost/custom-image:v1-${cross_arch}"

  pull_default_image_and_copy_to "$custom_image"

  local container="native-test-container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --image "$custom_image" \
    --container "$container"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  local host_arch
  case "$(uname -m)" in
    x86_64)  host_arch="amd64" ;;
    aarch64) host_arch="arm64" ;;
    ppc64le) host_arch="ppc64le" ;;
  esac

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${host_arch}"
}

@test "create: Native and cross-arch containers can coexist" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_default_image
  pull_default_image_cross_arch

  local native_container
  native_container="$(get_system_id)-toolbox-$(get_system_version)"

  local cross_container
  cross_container="$(get_system_id)-toolbox-$(get_system_version)-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create

  assert_success
  assert_line --index 0 "Created container: $native_container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create --arch "${cross_arch}"

  assert_success
  assert_line --index 0 "Created container: $cross_container"
  assert [ ${#lines[@]} -eq 2 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman ps --all --format "{{.Names}}"

  assert_success
  assert_line --index 0 "$native_container"
  assert_line --index 1 "$cross_container"
  assert [ ${#lines[@]} -eq 2 ]
}

@test "create: Cross-arch image stored with arch-suffixed tag" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local container="fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman image exists "${IMAGES[fedora]}:44-${cross_arch}"

  assert_success

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: Image tag with arch suffix is not doubled" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  local container="fedora-toolbox-44-${binfmt_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --image "fedora-toolbox:44-${binfmt_arch}"

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman image exists "${IMAGES[fedora]}:44-${binfmt_arch}-${cross_arch}"

  assert_failure
}

@test "create: Cross-arch entry point has --arch and --arch-emulator-path args" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local container="fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_success
  assert_line --index 0 "Created container: $container"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman inspect \
        --format '{{join .Config.Cmd " "}}' \
        --type container \
        "$container"

  assert_success
  assert_output --partial "--arch"
  assert_output --partial "--arch-emulator-path"

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "$container"

  assert_success
  assert_output "${cross_arch}"
}

@test "create: Try cross-arch with non-existing registry" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --image "nonexistent-registry.invalid/fake-toolbox:44"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to verify: image nonexistent-registry.invalid/fake-toolbox:44 does not support architecture ${cross_arch} or the image does not exists at all"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "create: Try cross-arch with single-arch Arch Linux image" {
  if [ "$(uname -m)" = "aarch64" ]; then
    skip "arm64 is the host architecture"
  fi

  local cross_arch
  cross_arch="arm64"

  local image_arch
  image_arch="${IMAGES[arch]}:latest"

  # Arch Linux Toolbox is a single-arch image (amd64 only).
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --image "${image_arch}"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: image ${image_arch} is a single-architecture image for amd64, but ${cross_arch} was requested"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "create: Try cross-arch with multi-arch Ubuntu image for ppc64le" {
  if [ "$(uname -m)" = "ppc64le" ]; then
    skip "ppc64le is the host architecture"
  fi

  local image_ubuntu
  image_ubuntu="${IMAGES[ubuntu]}:22.04"

  # Ubuntu Toolbox does not support ppc64le (amd64 and arm64 only).
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch ppc64le \
    --image "${image_ubuntu}"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to verify: image ${image_ubuntu} does not support architecture ppc64le or the image does not exists at all"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "create: Try with an unsupported architecture" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch sparc64 \
    --distro fedora \
    --release 44

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: architecture 'sparc64' is not supported by Toolbx"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "create: Try with an empty architecture value" {
  local host_arch
  case "$(uname -m)" in
    x86_64)  host_arch="amd64" ;;
    aarch64) host_arch="arm64" ;;
    ppc64le) host_arch="ppc64le" ;;
  esac

  pull_distro_image fedora 44

  local default_container
  default_container="fedora-toolbox-44-native"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "" \
    --distro fedora \
    --release 44 \
    --container "$default_container"

  assert_success
  assert_line --index 0 "Created container: ${default_container}"
  assert_line --index 1 "Enter with: toolbox enter ${default_container}"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run podman inspect \
        --format '{{index .Config.Labels "toolbox-arch"}}' \
        --type container \
        "${default_container}"

  assert_success
  assert_output "${host_arch}"
}

@test "create: Try with a case-sensitive architecture value" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch ARM64 \
    --distro fedora \
    --release 44

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: architecture 'ARM64' is not supported by Toolbx"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "create: Try with conflicting --arch and image tag" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local conflicting_arch
  if [ "${cross_arch}" = "arm64" ]; then
    conflicting_arch="ppc64le"
  else
    conflicting_arch="arm64"
  fi

  local tag_arch
  tag_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${conflicting_arch}" \
    --image "fedora-toolbox:44-${tag_arch}"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: conflicting architecture specifications"
  assert_line --index 1 "--arch=${conflicting_arch} but image tag specifies ${cross_arch}"
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try with --release containing architecture suffix" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --distro fedora \
    --release 44-aarch64

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try with --distro and --image together with --arch" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --image "fedora-toolbox:44"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --distro and --image cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try with --image and --release together with --arch" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --image "fedora-toolbox:44" \
    --release 44

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --image and --release cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try creating duplicate cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  local container="fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --image "fedora-toolbox:44-${cross_arch}"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container $container already exists"
  assert_line --index 1 "Enter with: toolbox enter $container"
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try cross-arch without QEMU in PATH" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  local restricted_path="$BATS_TEST_TMPDIR/no-qemu-path"
  build_restricted_path "$restricted_path" "qemu-*"

  if [ -e "$restricted_path/qemu-${binfmt_arch}-static" ] || [ -e "$restricted_path/qemu-${binfmt_arch}" ]; then
    fail "qemu binaries were not excluded from restricted PATH"
  fi

  run --keep-empty-lines --separate-stderr env PATH="$restricted_path" \
    "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: Cannot create container for architecture ${cross_arch}"
  assert_line --index 1 "The host system does not have the required support: No ${cross_arch} statically linked QEMU emulator binary found"
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try cross-arch with fake QEMU (shell script, not ELF)" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  local restricted_path="$BATS_TEST_TMPDIR/fake-qemu-path"
  build_restricted_path "$restricted_path" "qemu-*"

  cat > "$restricted_path/qemu-${binfmt_arch}-static" <<'FAKE'
#!/bin/sh
echo "I am not a real QEMU"
FAKE

  chmod +x "$restricted_path/qemu-${binfmt_arch}-static"

  run --keep-empty-lines --separate-stderr env PATH="$restricted_path" \
    "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: Cannot create container for architecture ${cross_arch}"
  assert_line --index 1 "The host system does not have the required support: No ${cross_arch} statically linked QEMU emulator binary found"
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "create: Try cross-arch without skopeo in PATH" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local restricted_path="$BATS_TEST_TMPDIR/no-skopeo-path"
  build_restricted_path "$restricted_path" "skopeo"

  if [ -e "$restricted_path/skopeo" ]; then
    fail "skopeo was not excluded from restricted PATH"
  fi

  run --keep-empty-lines --separate-stderr env PATH="$restricted_path" \
    "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: Cannot inspect image ${IMAGES[fedora]}:44 for architecture ${cross_arch}: skopeo is not installed."
  assert_line --index 1 "Skopeo is required for creating non-native architecture containers."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}


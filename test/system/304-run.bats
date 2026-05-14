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

@test "run: Smoke test with true(1) in cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Try with an unsupported --arch value" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch sparc64 \
    --distro fedora \
    --release 44 \
    true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: architecture 'sparc64' is not supported by Toolbx"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "run: Try with a missing command when --arch is specified" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: missing argument for \"run\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "run: Verify machine architecture with uname(1) in cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    uname -m

  assert_success
  assert_line --index 0 "${binfmt_arch}"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that /run/.toolboxenv exists in cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    test -f /run/.toolboxenv

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that a specific cross-arch container is used with --container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local container="cross-arch-run-container"
  create_distro_container_cross_arch fedora 44 "${cross_arch}" "$container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --container "$container" \
    true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that a cross-arch container is reused on second run" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that the native container is used when --arch is omitted" {
  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that binfmt_misc is mounted inside cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    test -d /proc/sys/fs/binfmt_misc

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that QEMU emulator is registered in binfmt_misc inside cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    sh -c "test -f /proc/sys/fs/binfmt_misc/qemu-${binfmt_arch} \
           || test -f /proc/sys/fs/binfmt_misc/qemu-${binfmt_arch}-static"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that QEMU emulator registration is enabled inside cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    sh -c "cat /proc/sys/fs/binfmt_misc/qemu-${binfmt_arch} 2>/dev/null \
           || cat /proc/sys/fs/binfmt_misc/qemu-${binfmt_arch}-static 2>/dev/null"

  assert_success
  assert_line --index 0 "enabled"

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Try a non-existent command inside a cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local cmd="nonexistent-command-xyz123"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run -127 --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    "$cmd"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 "bash: line 1: exec: $cmd: not found"
  assert_line --index 1 "Error: command $cmd not found in container fedora-toolbox-44-${cross_arch}"
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "run: Exit code propagation — false(1) in cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    false

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Exit code propagation — 'exit 2' in cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run -2 --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    /bin/sh -c 'exit 2'

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Short entry point error propagates from cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  # shellcheck disable=SC2030
  export TOOLBX_FAIL_ENTRY_POINT=1

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2031,SC2154
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: TOOLBX_FAIL_ENTRY_POINT is set"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "run: Multi-line entry point error propagates from cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  # shellcheck disable=SC2030
  export TOOLBX_FAIL_ENTRY_POINT=2

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2031,SC2154
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: TOOLBX_FAIL_ENTRY_POINT is set"
  assert_line --index 1 "This environment variable should only be set when testing."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "run: Smoke test with 5s entry point delay in cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  # shellcheck disable=SC2030
  export TOOLBX_DELAY_ENTRY_POINT=5

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2031,SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Stop and restart cross-arch container re-initializes correctly" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local container="fedora-toolbox-44-${cross_arch}"
  create_distro_container_cross_arch fedora 44 "${cross_arch}" "$container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_success

  podman stop "$container" >/dev/null

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: 'sudo id' inside cross-arch container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    sudo id

  assert_success
  assert_line --index 0 "uid=0(root) gid=0(root) groups=0(root)"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: binfmt arch name alias (aarch64) resolves same as OCI name (arm64)" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  if [ "${cross_arch}" != "arm64" ]; then
    skip "alias test only applies when cross-arch is arm64"
  fi

  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch aarch64 \
    --distro fedora \
    --release 44 \
    uname -m

  assert_success
  assert_line --index 0 "aarch64"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Native and cross-arch containers coexist without interference" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  local host_binfmt_arch
  host_binfmt_arch="$(uname -m)"

  create_default_container
  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run uname -m

  assert_success
  assert_line --index 0 "${host_binfmt_arch}"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    uname -m

  assert_success
  assert_line --index 0 "${binfmt_arch}"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Error when QEMU is missing from recorded path and from /usr/bin/ fallback" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  if [ -f "/usr/bin/qemu-${binfmt_arch}-static" ] || \
     [ -f "/usr/bin/qemu-${binfmt_arch}" ]; then
    skip "QEMU is in /usr/bin/ — fallback would succeed, see the Warning test"
  fi

  local real_qemu=""
  for candidate in "qemu-${binfmt_arch}-static" "qemu-${binfmt_arch}"; do
    if command -v "$candidate" >/dev/null 2>&1; then
      real_qemu="$(command -v "$candidate")"
      break
    fi
  done

  local restricted_path="$BATS_TEST_TMPDIR/fake-qemu-path"
  build_restricted_path "$restricted_path" "qemu-*"
  ln -s "$real_qemu" "$restricted_path/qemu-${binfmt_arch}-static"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  run --keep-empty-lines --separate-stderr env PATH="$restricted_path" \
    "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_success

  local expected_interp="/run/host${restricted_path}/qemu-${binfmt_arch}-static"

  rm "$restricted_path/qemu-${binfmt_arch}-static"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  lines=("${stderr_lines[@]}")
  assert_line --index 0 \
    "Warning: QEMU emulator not found at expected path '${expected_interp}', using fallback at '/run/host/usr/bin/'"
  assert_line --index 1 "Error: Cannot run container for architecture ${cross_arch}:"
  assert_line --index 2 "The host system does not have the required support: No ${cross_arch} statically linked QEMU emulator binary found"
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Warning when QEMU not at recorded path but found at /usr/bin/ fallback" {
  local cross_arch
  cross_arch="$(get_cross_arch)"

  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  if [ ! -f "/usr/bin/qemu-${binfmt_arch}-static" ] && \
     [ ! -f "/usr/bin/qemu-${binfmt_arch}" ]; then
    skip "QEMU not in /usr/bin/ — fallback would fail, see the Error test"
  fi

  local real_qemu=""
  for candidate in "qemu-${binfmt_arch}-static" "qemu-${binfmt_arch}"; do
    if command -v "$candidate" >/dev/null 2>&1; then
      real_qemu="$(command -v "$candidate")"
      break
    fi
  done

  local restricted_path="$BATS_TEST_TMPDIR/fake-qemu-path"
  build_restricted_path "$restricted_path" "qemu-*"
  ln -s "$real_qemu" "$restricted_path/qemu-${binfmt_arch}-static"

  pull_distro_image_cross_arch fedora 44 "${cross_arch}"

  # Create the container: the symlink's path is stored as --arch-emulator-path.
  run --keep-empty-lines --separate-stderr env PATH="$restricted_path" \
    "$TOOLBX" --assumeyes create \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44

  assert_success

  local expected_interp="/run/host${restricted_path}/qemu-${binfmt_arch}-static"

  rm "$restricted_path/qemu-${binfmt_arch}-static"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  lines=("${stderr_lines[@]}")
  assert_line --index 0 \
    "Warning: QEMU emulator not found at expected path '${expected_interp}', using fallback at '/run/host/usr/bin/'"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# shellcheck shell=bats
#
# Copyright © 2025 - 2026 Hadi Chokr <hadichokr@icloud.com>
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

readonly UPGRADE_IMAGE="localhost/toolbx-upgrade-test:latest"

setup() {
  bats_require_minimum_version 1.10.0
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

# Builds an image carrying the labels read by 'upgrade'.
#
# Parameters:
# ===========
# - assumeyes - 'true' to also set the unattended variant of the label
build_upgrade_image() {
  local assumeyes
  assumeyes="$1"

  local default_image
  default_image="$(get_default_image)"

  pull_default_image

  {
    echo "FROM $default_image"
    echo ""
    echo "LABEL com.github.containers.toolbx.package-manager.update=\"echo upgrading\""

    if [ "$assumeyes" = "true" ]; then
      echo "LABEL com.github.containers.toolbx.package-manager.update.assumeyes=\"echo upgrading --assumeyes\""
    fi
  } >"$BATS_TEST_TMPDIR"/Containerfile

  podman build --quiet --tag "$UPGRADE_IMAGE" "$BATS_TEST_TMPDIR" >/dev/null \
    || fail "Podman couldn't build the test image"

  rm --force "$BATS_TEST_TMPDIR"/Containerfile
}

create_upgrade_container() {
  local container
  container="$1"

  "$TOOLBX" --assumeyes create --image "$UPGRADE_IMAGE" --container "$container" >/dev/null \
    || fail "Toolbx couldn't create container '$container'"
}

@test "upgrade: Smoke test (using --assumeyes)" {
  build_upgrade_image true
  create_upgrade_container upgrade-test

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes upgrade upgrade-test

  assert_success
  assert_line --index 0 "Upgrading container upgrade-test with: echo upgrading --assumeyes"
  assert_line --index 1 "upgrading --assumeyes"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "upgrade: All containers, with one that doesn't support it (using --assumeyes)" {
  build_upgrade_image true
  create_upgrade_container upgrade-test
  create_distro_container arch latest arch-toolbox-latest

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes upgrade --all

  assert_success
  assert_line --index 0 "Upgrading container upgrade-test with: echo upgrading --assumeyes"
  assert_line --index 1 "upgrading --assumeyes"
  assert [ ${#lines[@]} -eq 2 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Warning: skipping container arch-toolbox-latest: it doesn't support 'upgrade'"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "upgrade: All containers, without any containers" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" upgrade --all

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "upgrade: Try without any arguments" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" upgrade

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: missing argument for \"upgrade\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "upgrade: Try using both --all and a container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" upgrade --all some-container

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: option --all and CONTAINER cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "upgrade: Try a non-existent container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" upgrade wrong-container

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to inspect container wrong-container"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "upgrade: Try a non-Toolbx container" {
  local busybox_image
  busybox_image="$(get_busybox_image)"

  pull_distro_image busybox
  podman create --name busybox-container "$busybox_image" >/dev/null

  run --keep-empty-lines --separate-stderr "$TOOLBX" upgrade busybox-container

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: busybox-container is not a Toolbx container"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "upgrade: Try a container without the package manager labels (Arch Linux)" {
  create_distro_container arch latest arch-toolbox-latest

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes upgrade arch-toolbox-latest

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container arch-toolbox-latest does not support \"upgrade\""
  assert_line --index 1 "Recreate it with a newer image or file a bug."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "upgrade: Try without --assumeyes when not connected to a terminal" {
  build_upgrade_image true
  create_upgrade_container upgrade-test

  run --keep-empty-lines --separate-stderr "$TOOLBX" upgrade upgrade-test

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: confirmation required to upgrade container upgrade-test"
  assert_line --index 1 "Use option '--assumeyes' to upgrade without confirmation."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "upgrade: Try --assumeyes with a container that only supports interactive upgrades" {
  build_upgrade_image false
  create_upgrade_container upgrade-test

  run --keep-empty-lines --separate-stderr "$TOOLBX" --assumeyes upgrade upgrade-test

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container upgrade-test doesn't support unattended upgrades"
  assert_line --index 1 "Upgrade it without '--assumeyes', or recreate it with a newer image."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

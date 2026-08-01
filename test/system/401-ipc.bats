# shellcheck shell=bats
#
# Copyright © 2023 – 2026 Red Hat, Inc.
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
  cleanup_all
  pushd "$HOME" || return 1

  local cross_arch
  cross_arch="$(get_cross_arch)"
  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"
}

teardown_file() {
  popd || return 1
  cleanup_all
}

@test "ipc: No namespace inside non-native container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"
  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "${cross_arch}")"

  local ns_host
  ns_host=$(readlink /proc/$$/ns/ipc)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    sh -c 'readlink /proc/$$/ns/ipc'

  assert_success
  assert_line --index 0 "$ns_host"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # -- Verify the container runs under the expected non-native architecture
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "${cross_arch}" \
    --distro fedora \
    --release 44 \
    sh -c 'lscpu'

  assert_success
  assert_line --index 0 --regexp "^Architecture:[[:space:]]+${binfmt_arch}$"
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

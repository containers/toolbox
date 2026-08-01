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

@test "user: Separate namespace inside non-native container" {
  local ns_host
  ns_host=$(readlink /proc/$$/ns/user)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 sh -c 'readlink /proc/$$/ns/user'

  assert_success
  assert_line --index 0 --regexp '^user:\[[[:digit:]]+\]$'
  refute_line --index 0 "$ns_host"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # -- Architecture check
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 sh -c 'lscpu'

  assert_success
  local binfmt_arch
  binfmt_arch="$(oci_arch_to_binfmt "$(get_cross_arch)")"
  assert_line --index 0 --regexp "^Architecture:[[:space:]]+${binfmt_arch}$"
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside non-native container" {
  container_root_file_system="$(podman unshare podman mount fedora-toolbox-44-$(get_cross_arch))"

  "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 true

  run --keep-empty-lines --separate-stderr podman unshare cat "$container_root_file_system/etc/shadow"
  podman unshare podman unmount fedora-toolbox-44-$(get_cross_arch)

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside non-native container" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside non-native container" {
  container_root_file_system="$(podman unshare podman mount fedora-toolbox-44-$(get_cross_arch))"

  "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 true

  run --keep-empty-lines --separate-stderr podman unshare cat "$container_root_file_system/etc/shadow"
  podman unshare podman unmount fedora-toolbox-44-$(get_cross_arch)

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside non-native container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 cat /etc/group

  assert_success
  assert_line --regexp "^$USER:x:[[:digit:]]+:$USER$"
  assert_line --regexp "^wheel:x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: id(1) for $USER inside non-native container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 id

  assert_success
  assert [ ${#lines[@]} -eq 1 ]

  local output_id="${lines[0]}"

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 id "$USER"

  assert_success
  assert_line --index 0 "$output_id"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}
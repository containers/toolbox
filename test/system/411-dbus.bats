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

@test "dbus: Session bus inside non-native container" {
  local expected_response
  expected_response="$(gdbus call \
                         --session \
                         --dest org.freedesktop.DBus \
                         --object-path /org/freedesktop/DBus \
                         --method org.freedesktop.DBus.Peer.Ping)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "$(get_cross_arch)" \
    --distro fedora \
    --release 44 \
    gdbus call \
      --session \
      --dest org.freedesktop.DBus \
      --object-path /org/freedesktop/DBus \
      --method org.freedesktop.DBus.Peer.Ping

  assert_success
  assert_line --index 0 "$expected_response"
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

@test "dbus: System bus inside non-native container" {
  local expected_response
  expected_response="$(gdbus call \
                         --system \
                         --dest org.freedesktop.systemd1 \
                         --object-path /org/freedesktop/systemd1 \
                         --method org.freedesktop.DBus.Properties.Get \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "$(get_cross_arch)" \
    --distro fedora \
    --release 44 \
    gdbus call \
      --system \
      --dest org.freedesktop.systemd1 \
      --object-path /org/freedesktop/systemd1 \
      --method org.freedesktop.DBus.Properties.Get \
      org.freedesktop.systemd1.Manager \
      Version

  assert_success
  assert_line --index 0 "$expected_response"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}
# shellcheck shell=bats
#
# Copyright © 2025 – 2026 Red Hat, Inc.
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

setup() {
  bats_require_minimum_version 1.10.0
  cleanup_all
  pushd "$HOME" || return 1
  skip_if_no_cross_arch_support
}

teardown() {
  popd || return 1
  cleanup_all
}

@test "rpm: %_netsharedpath inside non-native container" {
  local cross_arch
  cross_arch="$(get_cross_arch)"
  create_distro_container_cross_arch fedora 44 "${cross_arch}" "fedora-toolbox-44-${cross_arch}"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 rpm --eval %_netsharedpath

  assert_success
  assert_line --index 0 "/dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "$(get_cross_arch)" \
    --distro fedora \
    --release 44 \
    cat /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 ""
  assert_line --index 3 "%_netsharedpath /dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --arch "$(get_cross_arch)" \
    --distro fedora \
    --release 44 \
    stat \
      --format "%A %U:%G" \
      /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

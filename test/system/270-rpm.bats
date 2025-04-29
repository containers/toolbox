# shellcheck shell=bats
#
# Copyright © 2025 Red Hat, Inc.
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

# bats file_tags=runtime-environment

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.10.0
  _setup_environment
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

# bats test_tags=arch-fedora
@test "rpm: %_netsharedpath inside the default container" {
  local system_id
  system_id="$(get_system_id)"

  if [ "$system_id" != "fedora" ]; then
    skip "doesn't use RPM"
  fi

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run rpm --eval %_netsharedpath

  assert_success
  assert_line --index 0 "/dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run cat /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://github.com/containers/toolbox"
  assert_line --index 2 ""
  assert_line --index 3 "%_netsharedpath /dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run stat \
                                                           --format "%A %U:%G" \
                                                           /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "rpm: %_netsharedpath inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 rpm --eval %_netsharedpath

  assert_success
  assert_line --index 0 "/dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    cat /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://github.com/containers/toolbox"
  assert_line --index 2 ""
  assert_line --index 3 "%_netsharedpath /dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    stat \
      --format "%A %U:%G" \
      /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "rpm: %_netsharedpath inside RHEL 8.10" {
  create_distro_container rhel 8.10 rhel-toolbox-8.10

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 rpm --eval %_netsharedpath

  assert_success
  assert_line --index 0 "/dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    cat /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://github.com/containers/toolbox"
  assert_line --index 2 ""
  assert_line --index 3 "%_netsharedpath /dev:/media:/mnt:/proc:/sys:/tmp:/var/lib/flatpak:/var/lib/libvirt"
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    stat \
      --format "%A %U:%G" \
      /usr/lib/rpm/macros.d/macros.toolbox

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

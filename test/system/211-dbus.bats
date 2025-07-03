# shellcheck shell=bats
#
# Copyright © 2023 – 2025 Red Hat, Inc.
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

setup_file() {
  bats_require_minimum_version 1.10.0
  cleanup_all
  pushd "$HOME" || return 1

  if echo "$TOOLBX_TEST_SYSTEM_TAGS" | grep "arch" >/dev/null 2>/dev/null; then
    create_distro_container arch latest arch-toolbox-latest
  fi

  if echo "$TOOLBX_TEST_SYSTEM_TAGS" | grep "fedora" >/dev/null 2>/dev/null; then
    create_default_container
    create_distro_container fedora 34 fedora-toolbox-34
    create_distro_container rhel 8.10 rhel-toolbox-8.10
  fi

  if echo "$TOOLBX_TEST_SYSTEM_TAGS" | grep "ubuntu" >/dev/null 2>/dev/null; then
    create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04
    create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04
    create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04
  fi
}

teardown_file() {
  popd || return 1
  cleanup_all
}

# bats test_tags=arch-fedora
@test "dbus: session bus inside the default container" {
  local expected_response
  expected_response="$(gdbus call \
                         --session \
                         --dest org.freedesktop.DBus \
                         --object-path /org/freedesktop/DBus \
                         --method org.freedesktop.DBus.Peer.Ping)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run gdbus call \
                                                           --session \
                                                           --dest org.freedesktop.DBus \
                                                           --object-path /org/freedesktop/DBus \
                                                           --method org.freedesktop.DBus.Peer.Ping

  assert_success
  assert_line --index 0 "$expected_response"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "dbus: session bus inside Arch Linux" {
  local expected_response
  expected_response="$(gdbus call \
                         --session \
                         --dest org.freedesktop.DBus \
                         --object-path /org/freedesktop/DBus \
                         --method org.freedesktop.DBus.Peer.Ping)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro arch \
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
}

# bats test_tags=arch-fedora
@test "dbus: session bus inside Fedora 34" {
  local expected_response
  expected_response="$(gdbus call \
                         --session \
                         --dest org.freedesktop.DBus \
                         --object-path /org/freedesktop/DBus \
                         --method org.freedesktop.DBus.Peer.Ping)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
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
}

# bats test_tags=arch-fedora
@test "dbus: session bus inside RHEL 8.10" {
  local expected_response
  expected_response="$(gdbus call \
                         --session \
                         --dest org.freedesktop.DBus \
                         --object-path /org/freedesktop/DBus \
                         --method org.freedesktop.DBus.Peer.Ping)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
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
}

# bats test_tags=ubuntu
@test "dbus: session bus inside Ubuntu 16.04" {
  busctl --user call org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.Peer Ping

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 16.04 \
    busctl --user call \
      org.freedesktop.DBus \
      /org/freedesktop/DBus \
      org.freedesktop.DBus.Peer \
      Ping

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "dbus: session bus inside Ubuntu 18.04" {
  busctl --user call org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.Peer Ping

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 18.04 \
    busctl --user call \
      org.freedesktop.DBus \
      /org/freedesktop/DBus \
      org.freedesktop.DBus.Peer \
      Ping

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "dbus: session bus inside Ubuntu 20.04" {
  busctl --user call org.freedesktop.DBus /org/freedesktop/DBus org.freedesktop.DBus.Peer Ping

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 20.04 \
    busctl --user call \
      org.freedesktop.DBus \
      /org/freedesktop/DBus \
      org.freedesktop.DBus.Peer \
      Ping

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "dbus: system bus inside the default container" {
  local expected_response
  expected_response="$(gdbus call \
                         --system \
                         --dest org.freedesktop.systemd1 \
                         --object-path /org/freedesktop/systemd1 \
                         --method org.freedesktop.DBus.Properties.Get \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run gdbus call \
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

# bats test_tags=arch-fedora
@test "dbus: system bus inside Arch Linux" {
  local expected_response
  expected_response="$(gdbus call \
                         --system \
                         --dest org.freedesktop.systemd1 \
                         --object-path /org/freedesktop/systemd1 \
                         --method org.freedesktop.DBus.Properties.Get \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro arch \
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

# bats test_tags=arch-fedora
@test "dbus: system bus inside Fedora 34" {
  local expected_response
  expected_response="$(gdbus call \
                         --system \
                         --dest org.freedesktop.systemd1 \
                         --object-path /org/freedesktop/systemd1 \
                         --method org.freedesktop.DBus.Properties.Get \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
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

# bats test_tags=arch-fedora
@test "dbus: system bus inside RHEL 8.10" {
  local expected_response
  expected_response="$(gdbus call \
                         --system \
                         --dest org.freedesktop.systemd1 \
                         --object-path /org/freedesktop/systemd1 \
                         --method org.freedesktop.DBus.Properties.Get \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
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

# bats test_tags=ubuntu
@test "dbus: system bus inside Ubuntu 16.04" {
  local expected_response
  expected_response="$(busctl --system get-property \
                         org.freedesktop.systemd1 \
                         /org/freedesktop/systemd1 \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 16.04 \
    busctl --system get-property \
      org.freedesktop.systemd1 \
      /org/freedesktop/systemd1 \
      org.freedesktop.systemd1.Manager \
      Version

  assert_success
  assert_line --index 0 "$expected_response"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "dbus: system bus inside Ubuntu 18.04" {
  local expected_response
  expected_response="$(busctl --system get-property \
                         org.freedesktop.systemd1 \
                         /org/freedesktop/systemd1 \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 18.04 \
    busctl --system get-property \
      org.freedesktop.systemd1 \
      /org/freedesktop/systemd1 \
      org.freedesktop.systemd1.Manager \
      Version

  assert_success
  assert_line --index 0 "$expected_response"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "dbus: system bus inside Ubuntu 20.04" {
  local expected_response
  expected_response="$(busctl --system get-property \
                         org.freedesktop.systemd1 \
                         /org/freedesktop/systemd1 \
                         org.freedesktop.systemd1.Manager \
                         Version)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 20.04 \
    busctl --system get-property \
      org.freedesktop.systemd1 \
      /org/freedesktop/systemd1 \
      org.freedesktop.systemd1.Manager \
      Version

  assert_success
  assert_line --index 0 "$expected_response"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

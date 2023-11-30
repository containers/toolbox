# shellcheck shell=bats
#
# Copyright © 2023 Red Hat, Inc.
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

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.7.0
  _setup_environment
  cleanup_containers
}

teardown() {
  cleanup_containers
}

@test "dbus: session bus" {
  local expected_response
  expected_response="$(gdbus call \
                         --session \
                         --dest org.freedesktop.DBus \
                         --object-path /org/freedesktop/DBus \
                         --method org.freedesktop.DBus.Peer.Ping)"

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run gdbus call \
                                                            --session \
                                                            --dest org.freedesktop.DBus \
                                                            --object-path /org/freedesktop/DBus \
                                                            --method org.freedesktop.DBus.Peer.Ping

  assert_success
  assert_line --index 0 "$expected_response"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "dbus: system bus" {
  local expected_response
  expected_response="$(gdbus call \
                         --system \
                         --dest org.freedesktop.systemd1 \
                         --object-path /org/freedesktop/systemd1 \
                         --method org.freedesktop.DBus.Properties.Get org.freedesktop.systemd1.Manager Version)"

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run gdbus call \
                                                            --system \
                                                            --dest org.freedesktop.systemd1 \
                                                            --object-path /org/freedesktop/systemd1 \
                                                            --method org.freedesktop.DBus.Properties.Get \
                                                              org.freedesktop.systemd1.Manager Version

  assert_success
  assert_line --index 0 "$expected_response"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

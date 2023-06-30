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

@test "network: no namespace" {
  local ns_host
  ns_host=$(readlink /proc/$$/ns/net)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run sh -c 'readlink /proc/$$/ns/net'

  assert_success
  assert_line --index 0 "$ns_host"
  assert [ ${#lines[@]} -eq 2 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "network: ping(8) inside the default container" {
  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "network: ping(8) inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro arch ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "network: ping(8) inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro fedora --release 34 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "network: ping(8) inside RHEL 8.7" {
  create_distro_container rhel 8.7 rhel-toolbox-8.7

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro rhel --release 8.7 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "network: ping(8) inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  run -2 --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 16.04 ping -c 2 f.root-servers.net

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -gt 0 ]

  skip "doesn't use ICMP Echo sockets"
}

@test "network: ping(8) inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 18.04 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "network: ping(8) inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 20.04 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

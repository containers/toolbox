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

readonly RESOLVER_PYTHON3='\
import socket; \
import sys; \
family = {"A": socket.AddressFamily.AF_INET, "AAAA": socket.AddressFamily.AF_INET6}; \
addr = socket.getaddrinfo(sys.argv[2], None, family[sys.argv[1]], socket.SocketKind.SOCK_RAW)[0][4][0]; \
print(addr)'

# shellcheck disable=SC2016
readonly RESOLVER_SH='resolvectl --legend false --no-pager --type "$0" query "$1" \
                      | cut --delimiter " " --fields 4'

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

@test "network: No namespace inside non-native container" {
  local ns_host
  ns_host=$(readlink /proc/$$/ns/net)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 sh -c 'readlink /proc/$$/ns/net'

  assert_success
  assert_line --index 0 "$ns_host"
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

@test "network: /etc/resolv.conf inside non-native container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 readlink /etc/resolv.conf

  assert_success

  if [ "${lines[0]}" = "/run/host/run/systemd/resolve/stub-resolv.conf" ]; then
    skip "host has absolute symlink"
  else
    assert_line --index 0 "/run/host/etc/resolv.conf"
  fi

  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "network: DNS inside non-native container" {
  local ipv4_skip=false
  local ipv4_addr
  if ! ipv4_addr="$(python3 -c "$RESOLVER_PYTHON3" A k.root-servers.net)"; then
    ipv4_skip=true
  fi

  local ipv6_skip=false
  local ipv6_addr
  if ! ipv6_addr="$(python3 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net)"; then
    ipv6_skip=true
  fi

  if $ipv4_skip && $ipv6_skip; then
    skip "DNS not working on host"
  fi

  if ! $ipv4_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --arch "$(get_cross_arch)" \
      --distro fedora \
      --release 44 \
      python3 -c "$RESOLVER_PYTHON3" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --arch "$(get_cross_arch)" \
      --distro fedora \
      --release 44 \
      python3 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

@test "network: ping(8) inside non-native container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --arch "$(get_cross_arch)" --distro fedora --release 44 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

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

readonly RESOLVER_PYTHON3='\
import socket; \
import sys; \
family = socket.AddressFamily.AF_INET if sys.argv[1] == "A" else 0; \
family = socket.AddressFamily.AF_INET6 if sys.argv[1] == "AAAA" else 0; \
addr = socket.getaddrinfo(sys.argv[2], None, family, socket.SocketKind.SOCK_RAW)[0][4][0]; \
print(addr)'

# shellcheck disable=SC2016
readonly RESOLVER_SH='resolvectl --legend false --no-pager --type "$0" query "$1" \
                      | cut --delimiter " " --fields 4'

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
@test "network: No namespace" {
  local ns_host
  ns_host=$(readlink /proc/$$/ns/net)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run sh -c 'readlink /proc/$$/ns/net'

  assert_success
  assert_line --index 0 "$ns_host"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "network: /etc/resolv.conf inside the default container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /etc/resolv.conf

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

# bats test_tags=arch-fedora
@test "network: /etc/resolv.conf inside Arch Linux" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch readlink /etc/resolv.conf

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

# bats test_tags=arch-fedora
@test "network: /etc/resolv.conf inside Fedora 34" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 readlink /etc/resolv.conf

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

# bats test_tags=arch-fedora
@test "network: /etc/resolv.conf inside RHEL 8.10" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 readlink /etc/resolv.conf

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

# bats test_tags=ubuntu
@test "network: /etc/resolv.conf inside Ubuntu 16.04" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 readlink /etc/resolv.conf

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

# bats test_tags=ubuntu
@test "network: /etc/resolv.conf inside Ubuntu 18.04" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 readlink /etc/resolv.conf

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

# bats test_tags=ubuntu
@test "network: /etc/resolv.conf inside Ubuntu 20.04" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 readlink /etc/resolv.conf

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

# bats test_tags=arch-fedora
@test "network: DNS inside the default container" {
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
    run --keep-empty-lines --separate-stderr "$TOOLBX" run python3 -c "$RESOLVER_PYTHON3" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run python3 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=arch-fedora
@test "network: DNS inside Arch Linux" {
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
      --distro arch \
      sh -c "$RESOLVER_SH" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --distro arch \
      sh -c "$RESOLVER_SH" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=arch-fedora
@test "network: DNS inside Fedora 34" {
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
      --distro fedora \
      --release 34 \
      python3 -c "$RESOLVER_PYTHON3" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --distro fedora \
      --release 34 \
      python3 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=arch-fedora
@test "network: DNS inside RHEL 8.10" {
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
      --distro rhel \
      --release 8.10 \
      /usr/libexec/platform-python3.6 -c "$RESOLVER_PYTHON3" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --distro rhel \
      --release 8.10 \
      /usr/libexec/platform-python3.6 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=ubuntu
@test "network: DNS inside Ubuntu 16.04" {
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
      --distro ubuntu \
      --release 16.04 \
      python3 -c "$RESOLVER_PYTHON3" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --distro ubuntu \
      --release 16.04 \
      python3 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=ubuntu
@test "network: DNS inside Ubuntu 18.04" {
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
      --distro ubuntu \
      --release 18.04 \
      python3 -c "$RESOLVER_PYTHON3" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --distro ubuntu \
      --release 18.04 \
      python3 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=ubuntu
@test "network: DNS inside Ubuntu 20.04" {
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
      --distro ubuntu \
      --release 20.04 \
      python3 -c "$RESOLVER_PYTHON3" A k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv4_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  if ! $ipv6_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run \
      --distro ubuntu \
      --release 20.04 \
      python3 -c "$RESOLVER_PYTHON3" AAAA k.root-servers.net

    assert_success
    assert_line --index 0 "$ipv6_addr"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=arch-fedora
@test "network: ping(8) inside the default container" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "network: ping(8) inside Arch Linux" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "network: ping(8) inside Fedora 34" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "network: ping(8) inside RHEL 8.10" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "network: ping(8) inside Ubuntu 16.04" {
  run -2 --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 ping -c 2 f.root-servers.net

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -gt 0 ]

  skip "doesn't use ICMP Echo sockets"
}

# bats test_tags=ubuntu
@test "network: ping(8) inside Ubuntu 18.04" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "network: ping(8) inside Ubuntu 20.04" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 ping -c 2 f.root-servers.net

  if [ "$status" -eq 1 ]; then
    skip "lost packets"
  fi

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

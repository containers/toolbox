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

@test "user: Separate namespace" {
  local ns_host
  ns_host=$(readlink /proc/$$/ns/user)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run sh -c 'readlink /proc/$$/ns/user'

  assert_success
  assert_line --index 0 --regexp '^user:\[[[:digit:]]+\]$'
  refute_line --index 0 "$ns_host"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside the default container" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount "$default_container")"

  "$TOOLBOX" run true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount "$default_container"

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount arch-toolbox-latest)"

  "$TOOLBOX" run --distro arch true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount arch-toolbox-latest

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount fedora-toolbox-34)"

  "$TOOLBOX" run --distro fedora --release 34 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount fedora-toolbox-34

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside RHEL 8.7" {
  create_distro_container rhel 8.7 rhel-toolbox-8.7
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount rhel-toolbox-8.7)"

  "$TOOLBOX" run --distro rhel --release 8.7 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount rhel-toolbox-8.7

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount ubuntu-toolbox-16.04)"

  "$TOOLBOX" run --distro ubuntu --release 16.04 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount ubuntu-toolbox-16.04

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount ubuntu-toolbox-18.04)"

  "$TOOLBOX" run --distro ubuntu --release 18.04 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount ubuntu-toolbox-18.04

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: root in shadow(5) inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount ubuntu-toolbox-20.04)"

  "$TOOLBOX" run --distro ubuntu --release 20.04 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount ubuntu-toolbox-20.04

  assert_success
  assert_line --regexp '^root::.+$'
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside the default container" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside Arch Linux" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  create_distro_container arch latest arch-toolbox-latest

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro arch cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside Fedora 34" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  create_distro_container fedora 34 fedora-toolbox-34

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro fedora --release 34 cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside RHEL 8.7" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  create_distro_container rhel 8.7 rhel-toolbox-8.7

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro rhel --release 8.7 cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside Ubuntu 16.04" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 16.04 cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside Ubuntu 18.04" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 18.04 cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in passwd(5) inside Ubuntu 20.04" {
  local user_gecos
  user_gecos="$(getent passwd "$USER" | cut --delimiter : --fields 5)"

  local user_id_real
  user_id_real="$(id --real --user)"

  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 20.04 cat /etc/passwd

  assert_success
  assert_line --regexp "^$USER::$user_id_real:$user_id_real:$user_gecos:$HOME:$SHELL$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside the default container" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount "$default_container")"

  "$TOOLBOX" run true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount "$default_container"

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount arch-toolbox-latest)"

  "$TOOLBOX" run --distro arch true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount arch-toolbox-latest

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount fedora-toolbox-34)"

  "$TOOLBOX" run --distro fedora --release 34 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount fedora-toolbox-34

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside RHEL 8.7" {
  create_distro_container rhel 8.7 rhel-toolbox-8.7
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount rhel-toolbox-8.7)"

  "$TOOLBOX" run --distro rhel --release 8.7 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount rhel-toolbox-8.7

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount ubuntu-toolbox-16.04)"

  "$TOOLBOX" run --distro ubuntu --release 16.04 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount ubuntu-toolbox-16.04

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount ubuntu-toolbox-18.04)"

  "$TOOLBOX" run --distro ubuntu --release 18.04 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount ubuntu-toolbox-18.04

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in shadow(5) inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04
  container_root_file_system="$("$PODMAN" unshare "$PODMAN" mount ubuntu-toolbox-20.04)"

  "$TOOLBOX" run --distro ubuntu --release 20.04 true

  run --keep-empty-lines --separate-stderr "$PODMAN" unshare cat "$container_root_file_system/etc/shadow"
  "$PODMAN" unshare "$PODMAN" unmount ubuntu-toolbox-20.04

  assert_success
  refute_line --regexp "^$USER:.*$"
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside the default container" {
  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run cat /etc/group

  assert_success
  assert_line --regexp "^(sudo|wheel):x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro arch cat /etc/group

  assert_success
  assert_line --regexp "^wheel:x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro fedora --release 34 cat /etc/group

  assert_success
  assert_line --regexp "^wheel:x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside RHEL 8.7" {
  create_distro_container rhel 8.7 rhel-toolbox-8.7

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro rhel --release 8.7 cat /etc/group

  assert_success
  assert_line --regexp "^wheel:x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 16.04 cat /etc/group

  assert_success
  assert_line --regexp "^sudo:x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 18.04 cat /etc/group

  assert_success
  assert_line --regexp "^sudo:x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "user: $USER in group(5) inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  run --keep-empty-lines --separate-stderr "$TOOLBOX" run --distro ubuntu --release 20.04 cat /etc/group

  assert_success
  assert_line --regexp "^sudo:x:[[:digit:]]+:$USER$"
  assert [ ${#lines[@]} -gt 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

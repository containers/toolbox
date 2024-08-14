# shellcheck shell=bats
#
# Copyright © 2023 – 2024 Red Hat, Inc.
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
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

@test "environment variables: HISTFILESIZE inside the default container" {
  create_default_container

  if [ "$HISTFILESIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTFILESIZE=1001
  else
    ((HISTFILESIZE++))
  fi

  export HISTFILESIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$HISTFILESIZE"'

  assert_success
  assert_line --index 0 "$HISTFILESIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTFILESIZE inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  # shellcheck disable=SC2031
  if [ "$HISTFILESIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTFILESIZE=1001
  else
    ((HISTFILESIZE++))
  fi

  export HISTFILESIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$HISTFILESIZE"'

  assert_success
  assert_line --index 0 "$HISTFILESIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTFILESIZE inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  # shellcheck disable=SC2031
  if [ "$HISTFILESIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTFILESIZE=1001
  else
    ((HISTFILESIZE++))
  fi

  export HISTFILESIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$HISTFILESIZE"'

  assert_success
  assert_line --index 0 "$HISTFILESIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTFILESIZE inside RHEL 8.10" {
  create_distro_container rhel 8.10 rhel-toolbox-8.10

  # shellcheck disable=SC2031
  if [ "$HISTFILESIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTFILESIZE=1001
  else
    ((HISTFILESIZE++))
  fi

  export HISTFILESIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$HISTFILESIZE"'

  assert_success
  assert_line --index 0 "$HISTFILESIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTFILESIZE inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  # shellcheck disable=SC2031
  if [ "$HISTFILESIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTFILESIZE=1001
  else
    ((HISTFILESIZE++))
  fi

  export HISTFILESIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 \
                                             bash -c 'echo "$HISTFILESIZE"'

  assert_success
  assert_line --index 0 "$HISTFILESIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTFILESIZE inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  # shellcheck disable=SC2031
  if [ "$HISTFILESIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTFILESIZE=1001
  else
    ((HISTFILESIZE++))
  fi

  export HISTFILESIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 \
                                             bash -c 'echo "$HISTFILESIZE"'

  assert_success
  assert_line --index 0 "$HISTFILESIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTFILESIZE inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  # shellcheck disable=SC2031
  if [ "$HISTFILESIZE" = "" ]; then
    HISTFILESIZE=1001
  else
    ((HISTFILESIZE++))
  fi

  export HISTFILESIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 \
                                             bash -c 'echo "$HISTFILESIZE"'

  assert_success
  assert_line --index 0 "$HISTFILESIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTSIZE inside the default container" {
  skip "https://pagure.io/setup/pull-request/48"

  create_default_container

  if [ "$HISTSIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTSIZE=1001
  else
    ((HISTSIZE++))
  fi

  export HISTSIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$HISTSIZE"'

  assert_success
  assert_line --index 0 "$HISTSIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTSIZE inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  # shellcheck disable=SC2031
  if [ "$HISTSIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTSIZE=1001
  else
    ((HISTSIZE++))
  fi

  export HISTSIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$HISTSIZE"'

  assert_success
  assert_line --index 0 "$HISTSIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTSIZE inside Fedora 34" {
  skip "https://pagure.io/setup/pull-request/48"

  create_distro_container fedora 34 fedora-toolbox-34

  # shellcheck disable=SC2031
  if [ "$HISTSIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTSIZE=1001
  else
    ((HISTSIZE++))
  fi

  export HISTSIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$HISTSIZE"'

  assert_success
  assert_line --index 0 "$HISTSIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTSIZE inside RHEL 8.10" {
  skip "https://pagure.io/setup/pull-request/48"

  create_distro_container rhel 8.10 rhel-toolbox-8.10

  # shellcheck disable=SC2031
  if [ "$HISTSIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTSIZE=1001
  else
    ((HISTSIZE++))
  fi

  export HISTSIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$HISTSIZE"'

  assert_success
  assert_line --index 0 "$HISTSIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTSIZE inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  # shellcheck disable=SC2031
  if [ "$HISTSIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTSIZE=1001
  else
    ((HISTSIZE++))
  fi

  export HISTSIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 bash -c 'echo "$HISTSIZE"'

  assert_success
  assert_line --index 0 "$HISTSIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTSIZE inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  # shellcheck disable=SC2031
  if [ "$HISTSIZE" = "" ]; then
    # shellcheck disable=SC2030
    HISTSIZE=1001
  else
    ((HISTSIZE++))
  fi

  export HISTSIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 bash -c 'echo "$HISTSIZE"'

  assert_success
  assert_line --index 0 "$HISTSIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HISTSIZE inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  # shellcheck disable=SC2031
  if [ "$HISTSIZE" = "" ]; then
    HISTSIZE=1001
  else
    ((HISTSIZE++))
  fi

  export HISTSIZE

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 bash -c 'echo "$HISTSIZE"'

  assert_success
  assert_line --index 0 "$HISTSIZE"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HOSTNAME inside the default container" {
  create_default_container

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 --regexp "^(toolbx|$HOSTNAME)$"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HOSTNAME inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HOSTNAME inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "$HOSTNAME"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HOSTNAME inside RHEL 8.10" {
  create_distro_container rhel 8.10 rhel-toolbox-8.10

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HOSTNAME inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HOSTNAME inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: HOSTNAME inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: KONSOLE_VERSION inside the default container" {
  create_default_container

  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: KONSOLE_VERSION inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: KONSOLE_VERSION inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: KONSOLE_VERSION inside RHEL 8.10" {
  create_distro_container rhel 8.10 rhel-toolbox-8.10

  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: KONSOLE_VERSION inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: KONSOLE_VERSION inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: KONSOLE_VERSION inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: XTERM_VERSION inside the default container" {
  create_default_container

  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: XTERM_VERSION inside Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: XTERM_VERSION inside Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: XTERM_VERSION inside RHEL 8.10" {
  create_distro_container rhel 8.10 rhel-toolbox-8.10

  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: XTERM_VERSION inside Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: XTERM_VERSION inside Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "environment variables: XTERM_VERSION inside Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

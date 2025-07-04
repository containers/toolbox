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
@test "environment variables: HISTFILESIZE inside the default container" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HISTFILESIZE inside Arch Linux" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HISTFILESIZE inside Fedora 34" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HISTFILESIZE inside RHEL 8.10" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HISTFILESIZE inside Ubuntu 16.04" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HISTFILESIZE inside Ubuntu 18.04" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HISTFILESIZE inside Ubuntu 20.04" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HISTSIZE inside the default container" {
  skip "https://pagure.io/setup/pull-request/48"

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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HISTSIZE inside Arch Linux" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HISTSIZE inside Fedora 34" {
  skip "https://pagure.io/setup/pull-request/48"

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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HISTSIZE inside RHEL 8.10" {
  skip "https://pagure.io/setup/pull-request/48"

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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HISTSIZE inside Ubuntu 16.04" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HISTSIZE inside Ubuntu 18.04" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HISTSIZE inside Ubuntu 20.04" {
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
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HOSTNAME inside the default container" {
  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 --regexp "^(toolbx|$HOSTNAME)$"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HOSTNAME inside Arch Linux" {
  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HOSTNAME inside Fedora 34" {
  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "$HOSTNAME"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: HOSTNAME inside RHEL 8.10" {
  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HOSTNAME inside Ubuntu 16.04" {
  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HOSTNAME inside Ubuntu 18.04" {
  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: HOSTNAME inside Ubuntu 20.04" {
  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 bash -c 'echo "$HOSTNAME"'

  assert_success
  assert_line --index 0 "toolbx"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: KONSOLE_VERSION inside the default container" {
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: KONSOLE_VERSION inside Arch Linux" {
  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: KONSOLE_VERSION inside Fedora 34" {
  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: KONSOLE_VERSION inside RHEL 8.10" {
  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: KONSOLE_VERSION inside Ubuntu 16.04" {
  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: KONSOLE_VERSION inside Ubuntu 18.04" {
  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: KONSOLE_VERSION inside Ubuntu 20.04" {
  # shellcheck disable=SC2031
  if [ "$KONSOLE_VERSION" = "" ]; then
    export KONSOLE_VERSION=230804
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 bash -c 'echo "$KONSOLE_VERSION"'

  assert_success
  assert_line --index 0 "$KONSOLE_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: XTERM_VERSION inside the default container" {
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: XTERM_VERSION inside Arch Linux" {
  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: XTERM_VERSION inside Fedora 34" {
  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro fedora --release 34 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "environment variables: XTERM_VERSION inside RHEL 8.10" {
  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro rhel --release 8.10 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: XTERM_VERSION inside Ubuntu 16.04" {
  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: XTERM_VERSION inside Ubuntu 18.04" {
  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    # shellcheck disable=SC2030
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "environment variables: XTERM_VERSION inside Ubuntu 20.04" {
  # shellcheck disable=SC2031
  if [ "$XTERM_VERSION" = "" ]; then
    export XTERM_VERSION="XTerm(385)"
  fi

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 bash -c 'echo "$XTERM_VERSION"'

  assert_success
  assert_line --index 0 "$XTERM_VERSION"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

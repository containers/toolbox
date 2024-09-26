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
@test "ulimit: real-time non-blocking time (hard)" {
  local limit
  limit=$(ulimit -H -R)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -R

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: real-time non-blocking time (soft)" {
  local limit
  limit=$(ulimit -S -R)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -R

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: core file size (hard)" {
  local limit
  limit=$(ulimit -H -c)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -c

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: core file size (soft)" {
  local limit
  limit=$(ulimit -S -c)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -c

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: data segment size (hard)" {
  local limit
  limit=$(ulimit -H -d)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -d

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: data segment size (soft)" {
  local limit
  limit=$(ulimit -S -d)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -d

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: scheduling priority (hard)" {
  local limit
  limit=$(ulimit -H -e)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -e

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: scheduling priority (soft)" {
  local limit
  limit=$(ulimit -S -e)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -e

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: file size (hard)" {
  local limit
  limit=$(ulimit -H -f)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -f

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: file size (soft)" {
  local limit
  limit=$(ulimit -S -f)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -f

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of pending signals (hard)" {
  local limit
  limit=$(ulimit -H -i)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -i

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of pending signals (soft)" {
  local limit
  limit=$(ulimit -S -i)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -i

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: locked memory size (hard)" {
  local limit
  limit=$(ulimit -H -l)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -l

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: locked memory size (soft)" {
  local limit
  limit=$(ulimit -S -l)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -l

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: resident memory size (hard)" {
  local limit
  limit=$(ulimit -H -m)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -m

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: resident memory size (soft)" {
  local limit
  limit=$(ulimit -S -m)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -m

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of open files (hard)" {
  local limit
  limit=$(ulimit -H -n)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -n

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of open files (soft)" {
  local limit
  limit=$(ulimit -H -n)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -n

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: pipe size (hard)" {
  local limit
  limit=$(ulimit -H -p)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -p

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: pipe size (soft)" {
  local limit
  limit=$(ulimit -S -p)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -p

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: POSIX message queue size (hard)" {
  local limit
  limit=$(ulimit -H -q)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -q

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: POSIX message queue size (soft)" {
  local limit
  limit=$(ulimit -S -q)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -q

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: real-time scheduling priority (hard)" {
  local limit
  limit=$(ulimit -H -r)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -r

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: real-time scheduling priority (soft)" {
  local limit
  limit=$(ulimit -S -r)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -r

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: stack size (hard)" {
  local limit
  limit=$(ulimit -H -s)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -s

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: stack size (soft)" {
  local limit
  limit=$(ulimit -S -s)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -s

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: CPU time (hard)" {
  local limit
  limit=$(ulimit -H -t)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -t

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: CPU time (soft)" {
  local limit
  limit=$(ulimit -S -t)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -t

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of user processes (hard)" {
  local limit
  limit=$(ulimit -H -u)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -u

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of user processes (soft)" {
  local limit
  limit=$(ulimit -S -u)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -u

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: virtual memory size (hard)" {
  local limit
  limit=$(ulimit -H -v)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -v

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: virtual memory size (soft)" {
  local limit
  limit=$(ulimit -S -v)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -v

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of file locks (hard)" {
  local limit
  limit=$(ulimit -H -x)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -H -x

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: number of file locks (soft)" {
  local limit
  limit=$(ulimit -S -x)

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run ulimit -S -x

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

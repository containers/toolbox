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
  create_default_container
}

teardown_file() {
  popd || return 1
  cleanup_all
}

# bats test_tags=arch-fedora
@test "ulimit: Real-time non-blocking time (hard)" {
  local limit
  limit=$(ulimit -H -R)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -R'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Real-time non-blocking time (soft)" {
  local limit
  limit=$(ulimit -S -R)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -R'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Core file size (hard)" {
  local limit
  limit=$(ulimit -H -c)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -c'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Core file size (soft)" {
  local limit
  limit=$(ulimit -S -c)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -c'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Data segment size (hard)" {
  local limit
  limit=$(ulimit -H -d)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -d'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Data segment size (soft)" {
  local limit
  limit=$(ulimit -S -d)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -d'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Scheduling priority (hard)" {
  local limit
  limit=$(ulimit -H -e)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -e'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Scheduling priority (soft)" {
  local limit
  limit=$(ulimit -S -e)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -e'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: File size (hard)" {
  local limit
  limit=$(ulimit -H -f)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -f'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: File size (soft)" {
  local limit
  limit=$(ulimit -S -f)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -f'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of pending signals (hard)" {
  local limit
  limit=$(ulimit -H -i)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -i'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of pending signals (soft)" {
  local limit
  limit=$(ulimit -S -i)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -i'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Locked memory size (hard)" {
  local limit
  limit=$(ulimit -H -l)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -l'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Locked memory size (soft)" {
  local limit
  limit=$(ulimit -S -l)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -l'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Resident memory size (hard)" {
  local limit
  limit=$(ulimit -H -m)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -m'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Resident memory size (soft)" {
  local limit
  limit=$(ulimit -S -m)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -m'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of open files (hard)" {
  local limit
  limit=$(ulimit -H -n)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -n'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of open files (soft)" {
  local limit
  limit=$(ulimit -H -n)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -n'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Pipe size (hard)" {
  local limit
  limit=$(ulimit -H -p)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -p'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Pipe size (soft)" {
  local limit
  limit=$(ulimit -S -p)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -p'

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

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -q'

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

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -q'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Real-time scheduling priority (hard)" {
  local limit
  limit=$(ulimit -H -r)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -r'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Real-time scheduling priority (soft)" {
  local limit
  limit=$(ulimit -S -r)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -r'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Stack size (hard)" {
  local limit
  limit=$(ulimit -H -s)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -s'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Stack size (soft)" {
  local limit
  limit=$(ulimit -S -s)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -s'

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

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -t'

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

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -t'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of user processes (hard)" {
  local limit
  limit=$(ulimit -H -u)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -u'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of user processes (soft)" {
  local limit
  limit=$(ulimit -S -u)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -u'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Virtual memory size (hard)" {
  local limit
  limit=$(ulimit -H -v)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -v'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Virtual memory size (soft)" {
  local limit
  limit=$(ulimit -S -v)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -v'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of file locks (hard)" {
  local limit
  limit=$(ulimit -H -x)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -H -x'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "ulimit: Number of file locks (soft)" {
  local limit
  limit=$(ulimit -S -x)

  run --keep-empty-lines --separate-stderr "$TOOLBX" run bash -c 'ulimit -S -x'

  assert_success
  assert_line --index 0 "$limit"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

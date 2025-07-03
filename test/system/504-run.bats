# shellcheck shell=bats
#
# Copyright © 2024 Red Hat, Inc.
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

# bats file_tags=commands-options

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.8.0
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

@test "run: Smoke test with true(1) (forwarded to host)" {
  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox run true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with false(1) (forwarded to host)" {
  create_default_container

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox run false

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Try an unsupported distribution (forwarded to host)" {
  create_default_container

  local distro="foo"

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox run --distro "$distro" ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Smoke test with 'exit 2' (forwarded to host)" {
  create_default_container

  run -2 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox run /bin/sh -c 'exit 2'

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Try /etc as a command (forwarded to host)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  run -126 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox run /etc

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "bash: line 1: /etc: Is a directory"
  assert_line --index 1 "bash: line 1: exec: /etc: cannot execute: Is a directory"
  assert_line --index 2 "Error: failed to invoke command /etc in container $default_container"
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try a non-existent command (forwarded to host)" {
  local cmd="non-existent-command"

  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  run -127 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox run "$cmd"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "bash: line 1: exec: $cmd: not found"
  assert_line --index 1 "Error: command $cmd not found in container $default_container"
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

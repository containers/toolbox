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

# bats file_tags=commands-options

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.10.0
}

@test "completion: Smoke test with 'bash'" {
  run "$TOOLBX" completion bash

  assert_success
  assert [ ${#lines[@]} -gt 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "completion: Smoke test with 'fish'" {
  run "$TOOLBX" completion fish

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "completion: Smoke test with 'zsh'" {
  run "$TOOLBX" completion zsh

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "completion: Try without any arguments" {
  run --separate-stderr "$TOOLBX" completion

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: accepts 1 arg(s), received 0"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "completion: Try with invalid arguments" {
  run --separate-stderr "$TOOLBX" completion foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument \"foo\" for \"toolbox completion\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "completion: Try with unknown flag" {
  run --separate-stderr "$TOOLBX" completion --foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "completion: Try with unsupported shell" {
  run --separate-stderr "$TOOLBX" completion powershell

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument \"powershell\" for \"toolbox completion\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

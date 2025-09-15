# shellcheck shell=bats
#
# Copyright © 2024 – 2025 Red Hat, Inc.
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
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

@test "enter: Try a non-existent container (forwarded to host)" {
  create_default_container

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox enter wrong-container

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container wrong-container not found"
  assert_line --index 1 "Use the '--container' option to select a Toolbx."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "enter: Try an unsupported distribution (forwarded to host)" {
  create_default_container

  local distro="foo"

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox enter --distro "$distro"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

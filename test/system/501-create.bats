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
  _setup_environment
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

@test "create: Try with an invalid custom name (using positional argument, forwarded to host)" {
  create_default_container

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox create "ßpeci@l.N@m€"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for 'CONTAINER'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try with an invalid custom name (using option --container, forwarded to host)" {
  create_default_container

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox create --container "ßpeci@l.N@m€"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--container'"
  assert_line --index 1 "Container names must match '[a-zA-Z0-9][a-zA-Z0-9_.-]*'."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "create: Try the same name again (forwarded to host)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox create

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container $default_container already exists"
  assert_line --index 1 "Enter with: toolbox enter"
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "create: Try an unsupported distribution (forwarded to host)" {
  create_default_container

  local distro="foo"

  run -1 --keep-empty-lines --separate-stderr "$TOOLBX" run toolbox create --distro "$distro"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

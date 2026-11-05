# shellcheck shell=bats
#
# Copyright © 2025 – 2026 Hadi Chokr <hadichokr@icloud.com>
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

  XDG_DATA_DIRS="/usr/local/share:/usr/share"
  export XDG_DATA_DIRS

  XDG_BIN_HOME="$HOME/.local/bin"
  export XDG_BIN_HOME

  PATH="$XDG_BIN_HOME:$PATH"
  export PATH

  cleanup_all
  cleanup_exports
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_exports
  cleanup_all
}

@test "unexport: Try without --all, --app or --bin" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" unexport

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: option --all, --app or --bin is needed"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "unexport: Remove an exported application" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"

  run "$TOOLBX" export --app org.example.ToolbxTest --container "$container"
  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" unexport \
                                             --app org.example.ToolbxTest \
                                             --container "$container"

  assert_success
  assert_line --index 0 "Removed org.example.ToolbxTest.desktop exported from container $container"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  assert [ ! -e "$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop" ]
  assert [ ! -e "$XDG_DATA_HOME/toolbx/$container/icons/org.example.ToolbxTest" ]
}

@test "unexport: Remove an exported binary" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_bin "$container"

  run "$TOOLBX" export --bin toolbx-test --container "$container"
  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" unexport --bin toolbx-test --container "$container"

  assert_success
  assert_line --index 0 "Removed toolbx-test exported from container $container"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
  assert [ ! -e "$XDG_BIN_HOME/toolbx-test" ]
}

@test "unexport: Remove everything exported from a container" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"
  create_test_bin "$container"

  run "$TOOLBX" export --app org.example.ToolbxTest --container "$container"
  assert_success

  run "$TOOLBX" export --bin toolbx-test --container "$container"
  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" unexport --all --container "$container"

  assert_success
  assert_line --index 0 "Removed 2 files exported from container $container"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  assert [ ! -e "$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop" ]
  assert [ ! -e "$XDG_BIN_HOME/toolbx-test" ]
  assert [ ! -e "$XDG_DATA_HOME/toolbx/$container" ]
}

@test "unexport: rm removes exported files along with the container" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"
  create_test_bin "$container"

  run "$TOOLBX" export --app org.example.ToolbxTest --container "$container"
  assert_success

  run "$TOOLBX" export --bin toolbx-test --container "$container"
  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" rm --force "$container"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  assert [ ! -e "$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop" ]
  assert [ ! -e "$XDG_BIN_HOME/toolbx-test" ]
  assert [ ! -e "$XDG_DATA_HOME/toolbx/$container" ]
}

@test "unexport: rm by ID removes exported files along with the container" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_bin "$container"

  run "$TOOLBX" export --bin toolbx-test --container "$container"
  assert_success

  local id
  id="$(podman inspect --format "{{.Id}}" --type container "$container")"

  run --keep-empty-lines --separate-stderr "$TOOLBX" rm --force "$id"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
  assert [ ! -e "$XDG_BIN_HOME/toolbx-test" ]
}

@test "unexport: Try to remove a non-exported application" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" unexport \
                                             --app non-existent-app \
                                             --container my-container

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: file $XDG_DATA_HOME/applications/non-existent-app.desktop not found"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "unexport: Try to remove a file exported from another container" {
  local entry="$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop"

  mkdir --parents "$XDG_DATA_HOME/applications"
  printf "%s\n" "[Desktop Entry]" "Type=Application" "X-Toolbx-Container=other" >"$entry"

  run --keep-empty-lines --separate-stderr "$TOOLBX" unexport \
                                             --app org.example.ToolbxTest \
                                             --container my-container

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: file $entry was not exported from container my-container"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

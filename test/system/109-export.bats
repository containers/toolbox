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

  # Don't let the environment of the CI leak into the look-up in the container.
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

@test "export: Try without --app or --bin" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" export

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: option --app or --bin is needed"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "export: Try using both --app and --bin" {
  run --keep-empty-lines --separate-stderr "$TOOLBX" export --app foo --bin bar

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: options --app and --bin cannot be used together"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "export: Export an application" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" export \
                                             --app org.example.ToolbxTest \
                                             --container "$container"

  assert_success
  assert_line --index 0 "Exported org.example.ToolbxTest.desktop from container $container"
  assert [ ${#lines[@]} -eq 1 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  assert [ -f "$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop" ]
  assert [ -f "$XDG_DATA_HOME/toolbx/$container/icons/org.example.ToolbxTest/hicolor/48x48/apps/toolbx-test.png" ]
}

@test "export: Rewrite the keys of an exported application" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"

  run "$TOOLBX" export --app org.example.ToolbxTest --container "$container"
  assert_success

  local entry="$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop"

  run --keep-empty-lines --separate-stderr cat "$entry"

  assert_success
  assert_line --index 0 "[Desktop Entry]"
  assert_line --regexp "^Exec=/.+/toolbox run --container $container toolbx-test %U$"
  assert_line --regexp "^TryExec=/.+/toolbox$"
  assert_line "DBusActivatable=false"
  assert_line "#Path=/usr/share"
  assert_line "Icon=$XDG_DATA_HOME/toolbx/$container/icons/org.example.ToolbxTest/hicolor/48x48/apps/toolbx-test.png"
  assert_line "X-Toolbx-Container=$container"
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "export: Export an application twice" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"

  run "$TOOLBX" export --app org.example.ToolbxTest --container "$container"
  assert_success

  run --keep-empty-lines --separate-stderr "$TOOLBX" export \
                                             --app org.example.ToolbxTest \
                                             --container "$container"

  assert_success
  assert_line --index 0 "Exported org.example.ToolbxTest.desktop from container $container"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # The entry exported to the home directory is visible inside the container.
  # It must not be picked up and prefixed a second time.
  run --keep-empty-lines --separate-stderr grep "^Exec=" \
                                             "$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop"

  assert_success
  assert_line --index 0 --regexp "^Exec=/.+/toolbox run --container $container toolbx-test %U$"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "export: Export an application using a shortened name" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" export --app ToolbxTest --container "$container"

  assert_success
  assert_line --index 0 "Exported org.example.ToolbxTest.desktop from container $container"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
  assert [ -f "$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop" ]
}

@test "export: Export a binary" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_bin "$container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" export --bin toolbx-test --container "$container"

  assert_success
  assert_line --index 0 "Exported toolbx-test from container $container"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
  assert [ -x "$XDG_BIN_HOME/toolbx-test" ]

  run --keep-empty-lines --separate-stderr "$XDG_BIN_HOME/toolbx-test"

  assert_success
  assert_line --index 0 "toolbx-test"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "export: Try to overwrite a file exported from another container" {
  create_default_container

  local container
  container="$(get_latest_container_name)"
  create_test_app "$container"

  local entry="$XDG_DATA_HOME/applications/org.example.ToolbxTest.desktop"

  mkdir --parents "$XDG_DATA_HOME/applications"
  printf "%s\n" "[Desktop Entry]" "Type=Application" "X-Toolbx-Container=other" >"$entry"

  run --keep-empty-lines --separate-stderr "$TOOLBX" export \
                                             --app org.example.ToolbxTest \
                                             --container "$container"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: file $entry was exported from container other"
  assert_line --index 1 "Use option '--force' to overwrite it."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "export: Try to export a non-existent application" {
  create_default_container

  local container
  container="$(get_latest_container_name)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" export \
                                             --app non-existent-app \
                                             --container "$container"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: application non-existent-app not found in the container"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "export: Try to export a non-existent binary" {
  create_default_container

  local container
  container="$(get_latest_container_name)"

  run --keep-empty-lines --separate-stderr "$TOOLBX" export \
                                             --bin non-existent-bin \
                                             --container "$container"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: binary non-existent-bin not found in the container"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# shellcheck shell=bats
#
# Copyright © 2025 Hadi Chokr <hadichokr@icloud.com>
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

install_test_apps() {
  run "$TOOLBX" run sudo dnf -y install gimp vlc neovim
  assert_success
}

@test "unexport: Remove exported GIMP app from Fedora container" {
  create_default_container
  install_test_apps
  run "$TOOLBX" export --app gimp --container "$(get_latest_container_name)"
  assert_success
  assert [ -f "$HOME/.local/share/applications/gimp-$(get_latest_container_name).desktop" ]

  run "$TOOLBX" unexport --app gimp --container "$(get_latest_container_name)"
  assert_success
  assert [ ! -f "$HOME/.local/share/applications/gimp-$(get_latest_container_name).desktop" ]
}

@test "unexport: Remove all exported items from Fedora container" {
  create_default_container
  install_test_apps
  run "$TOOLBX" export --app gimp --container "$(get_latest_container_name)"
  run "$TOOLBX" export --bin nvim --container "$(get_latest_container_name)"

  assert_success
  assert [ -f "$HOME/.local/share/applications/gimp-$(get_latest_container_name).desktop" ]
  assert [ -f "$HOME/.local/bin/nvim" ]

  run "$TOOLBX" unexport --all --container "$(get_latest_container_name)"
  assert_success
  assert [ ! -f "$HOME/.local/share/applications/gimp-$(get_latest_container_name).desktop" ]
  assert [ ! -f "$HOME/.local/bin/nvim" ]
}

@test "unexport: Fail to remove non-exported app" {
  create_default_container
  install_test_apps

  run --separate-stderr "$TOOLBX" unexport --app fakeapp --container "$(get_latest_container_name)"
  assert_failure
  assert_output --partial "Error: application fakeapp not exported from container"
}


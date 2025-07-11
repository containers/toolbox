# shellcheck shell=bats
#
# Copyright © 2021 – 2024 Red Hat, Inc.
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
}

teardown() {
  cleanup_all
}

@test "enter: Try to enter the default container with no containers created" {
  run $TOOLBX enter <<< "n"

  assert_success
  assert_line --index 0 "No Toolbx containers found. Create now? [y/N] A container can be created later with the 'create' command."
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter the default container with more than 1 other containers present" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_container first
  create_container second

  run $TOOLBX enter

  assert_failure
  assert_line --index 0 "Error: container $default_container_name not found"
  assert_line --index 1 "Use the '--container' option to select a Toolbx."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter a specific container with no containers created " {
  run $TOOLBX enter wrong-container <<< "n"

  assert_success
  assert_line --index 0 "No Toolbx containers found. Create now? [y/N] A container can be created later with the 'create' command."
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter a specific non-existent container with other containers present" {
  create_container other-container

  run $TOOLBX enter wrong-container

  assert_failure
  assert_line --index 0 "Error: container wrong-container not found"
  assert_line --index 1 "Use the '--container' option to select a Toolbx."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "enter: Try to enter a container based on unsupported distribution" {
  local distro="foo"

  run $TOOLBX --assumeyes enter --distro "$distro"

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "enter: Try to enter a container based on Fedora but with wrong version" {
  run $TOOLBX enter -d fedora -r foobar

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBX enter --distro fedora --release -3

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "enter: Try to enter a container based on RHEL but with wrong version" {
  run $TOOLBX enter --distro rhel --release 8

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBX enter --distro rhel --release 8.2foo

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBX enter --distro rhel --release -2.1

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive number."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "enter: Try to enter a container based on non-default distro without providing a version" {
  local distro="fedora"

  local system_id
  system_id="$(get_system_id)"

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run $TOOLBX enter -d "$distro"

  assert_failure
  assert_line --index 0 "Error: option '--release' is needed"
  assert_line --index 1 "Distribution $distro doesn't match the host."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

# TODO: Write the test
@test "enter: Enter the default Toolbx" {
  skip "Testing of entering Toolbxes is not implemented"
}

# TODO: Write the test
@test "enter: Enter the default Toolbx when only 1 non-default Toolbx is present" {
  skip "Testing of entering Toolbxes is not implemented"
}

# TODO: Write the test
@test "enter: Enter a specific Toolbx" {
  skip "Testing of entering Toolbxes is not implemented"
}

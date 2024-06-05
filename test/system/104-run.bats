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

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.7.0
  _setup_environment
  cleanup_containers
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_containers
}

@test "run: Smoke test with true(1)" {
  create_default_container

  run "$TOOLBX" run true

  assert_success
  assert_output ""
}

@test "run: Smoke test with true(1) (using polling fallback)" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  export TOOLBX_RUN_USE_POLLING=1
  run --separate-stderr "$TOOLBX" run --verbose true

  assert_success
  assert_output ""

  # shellcheck disable=SC2154
  output="$stderr"

  assert_output --partial "Setting up polling ticker for container $default_container_name"
  refute_output --partial "Setting up watches for file system events from container $default_container_name"
  assert_output --partial "Handling polling tick"
}

@test "run: Smoke test with true(1) (using entry point with 5s delay)" {
  # shellcheck disable=SC2030
  export TOOLBX_DELAY_ENTRY_POINT=5

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with false(1)" {
  create_default_container

  run -1 "$TOOLBX" run false

  assert_failure
  assert_output ""
}

@test "run: Smoke test with Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  run --separate-stderr "$TOOLBX" run --distro arch true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with Arch Linux ('--release latest')" {
  create_distro_container arch latest arch-toolbox-latest

  run --separate-stderr "$TOOLBX" run --distro arch --release latest true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with Arch Linux ('--release rolling')" {
  create_distro_container arch latest arch-toolbox-latest

  run --separate-stderr "$TOOLBX" run --distro arch --release rolling true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  run --separate-stderr "$TOOLBX" run --distro fedora --release 34 true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with RHEL 8.9" {
  create_distro_container rhel 8.9 rhel-toolbox-8.9

  run --separate-stderr "$TOOLBX" run --distro rhel --release 8.9 true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  run --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  run --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Smoke test with Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  run --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 true

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that a login shell is used to invoke the command" {
  create_default_container

  cp "$HOME"/.bash_profile "$HOME"/.bash_profile.orig
  echo "echo \"$HOME/.bash_profile read\"" >>"$HOME"/.bash_profile

  run "$TOOLBX" run true

  mv "$HOME"/.bash_profile.orig "$HOME"/.bash_profile

  assert_success
  assert_line --index 0 "$HOME/.bash_profile read"
  assert [ ${#lines[@]} -eq 1 ]
}

@test "run: 'echo \"Hello World\"' inside the default container" {
  create_default_container

  run "$TOOLBX" --verbose run echo "Hello World"

  assert_success
  assert_line --index $((${#lines[@]}-1)) "Hello World"
}

@test "run: 'echo \"Hello World\"' inside a restarted container" {
  create_container running

  start_container running
  stop_container running

  run "$TOOLBX" --verbose run --container running echo "Hello World"

  assert_success
  assert_line --index $((${#lines[@]}-1)) "Hello World"
}

@test "run: 'sudo id' inside the default container" {
  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run sudo id

  assert_success
  assert_line --index 0 "uid=0(root) gid=0(root) groups=0(root)"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 1 ]
  else
    assert [ ${#lines[@]} -eq 2 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that /run/.containerenv exists" {
  create_default_container

  run --separate-stderr "$TOOLBX" run cat /run/.containerenv

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that /run/.toolboxenv exists" {
  create_default_container

  run --separate-stderr "$TOOLBX" run test -f /run/.toolboxenv

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Ensure that the default container is used" {
  test -z "${name+x}"

  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container
  create_container other-container

  run --separate-stderr "$TOOLBX" run cat /run/.containerenv

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC1090
  source <(echo "$output")
  run echo "$name"

  assert_success
  assert_output "$default_container_name"
}

@test "run: Ensure that a specific container is used" {
  run echo "$name"

  assert_success
  assert_output ""

  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container
  create_container other-container

  run --separate-stderr "$TOOLBX" run --container other-container cat /run/.containerenv

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC1090
  source <(echo "$output")
  run echo "$name"

  assert_success
  assert_output "other-container"
}

@test "run: Ensure that $HOME is used as a fallback working directory" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  local host_only_dir
  host_only_dir="$(mktemp --directory /var/tmp/toolbx-test-XXXXXXXXXX)"

  pushd "$host_only_dir"
  run --separate-stderr "$TOOLBX" run pwd
  popd

  rm --force --recursive "$host_only_dir"

  assert_success
  assert_line --index 0 "$HOME"
  assert [ ${#lines[@]} -eq 1 ]
  lines=("${stderr_lines[@]}")
  assert_line --index $((${#stderr_lines[@]}-2)) \
    "Error: directory $host_only_dir not found in container $default_container_name"
  assert_line --index $((${#stderr_lines[@]}-1)) "Using $HOME instead."
  assert [ ${#stderr_lines[@]} -gt 2 ]
}

@test "run: Pass down 1 additional file descriptor" {
  create_default_container

  # File descriptors 3 and 4 are reserved by Bats.
  run --separate-stderr "$TOOLBX" run --preserve-fds 3 readlink /proc/self/fd/5 5>/dev/null

  assert_success
  assert_line --index 0 "/dev/null"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "run: Try the non-existent default container with none other present" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  run --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container $default_container_name not found"
  assert_line --index 1 "Use the 'create' command to create a Toolbx."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try the non-existent default container with another present" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_container other-container

  run --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container $default_container_name not found"
  assert_line --index 1 "Use the 'create' command to create a Toolbx."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try a specific non-existent container with another present" {
  create_container other-container

  run --separate-stderr "$TOOLBX" run --container wrong-container true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container wrong-container not found"
  assert_line --index 1 "Use the 'create' command to create a Toolbx."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try an unsupported distribution" {
  local distro="foo"

  run --separate-stderr "$TOOLBX" --assumeyes run --distro "$distro" ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release -3')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release -3 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release -3.0')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release -3.0 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release -3.1')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release -3.1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release 0')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release 0 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release 0.0')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release 0.0 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release 0.1')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release 0.1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release 3.0')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release 3.0 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release 3.1')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release 3.1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release foo')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release foo ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try Fedora with an invalid release ('--release 3foo')" {
  run --separate-stderr "$TOOLBX" run --distro fedora --release 3foo ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release 8')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release 8 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release 8.0.0')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release 8.0.0 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release 8.0.1')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release 8.0.1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release 8.3.0')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release 8.3.0 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release 8.3.1')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release 8.3.1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release foo')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release foo ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release 8.2foo')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release 8.2foo ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release -2.1')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release -2.1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive number."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release -2.-1')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release -2.-1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try RHEL with an invalid release ('--release 2.-1')" {
  run --separate-stderr "$TOOLBX" run --distro rhel --release 2.-1 ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try a non-default distro without a release" {
  local distro="fedora"

  local system_id
  system_id="$(get_system_id)"

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run --separate-stderr "$TOOLBX" run --distro "$distro" ls

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: option '--release' is needed"
  assert_line --index 1 "Distribution $distro doesn't match the host."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try a failing entry point with a short error and no delay" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  # shellcheck disable=SC2030
  export TOOLBX_FAIL_ENTRY_POINT=1

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: TOOLBX_FAIL_ENTRY_POINT is set"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "run: Try a failing entry point with a short error and 5s delay" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_DELAY_ENTRY_POINT=5

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_FAIL_ENTRY_POINT=1

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: TOOLBX_FAIL_ENTRY_POINT is set"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "run: Try a failing entry point with a short error and 30s delay" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_DELAY_ENTRY_POINT=30

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_FAIL_ENTRY_POINT=1

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to initialize container $default_container_name"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "run: Try a failing entry point with a long error and no delay" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_FAIL_ENTRY_POINT=2

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: TOOLBX_FAIL_ENTRY_POINT is set"
  assert_line --index 1 "This environment variable should only be set when testing."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "run: Try a failing entry point with a long error and 5s delay" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_DELAY_ENTRY_POINT=5

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_FAIL_ENTRY_POINT=2

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: TOOLBX_FAIL_ENTRY_POINT is set"
  assert_line --index 1 "This environment variable should only be set when testing."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "run: Try a failing entry point with a long error and 30s delay" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  # shellcheck disable=SC2030,SC2031
  export TOOLBX_DELAY_ENTRY_POINT=30

  # shellcheck disable=SC2031
  export TOOLBX_FAIL_ENTRY_POINT=2

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to initialize container $default_container_name"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "run: Try a slow entry point that times out" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  # shellcheck disable=SC2031
  export TOOLBX_DELAY_ENTRY_POINT=30

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to initialize container $default_container_name"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "run: Smoke test with 'exit 2'" {
  create_default_container

  run -2 "$TOOLBX" run /bin/sh -c 'exit 2'
  assert_failure
  assert_output ""
}

@test "run: Pass down 1 invalid file descriptor" {
  local default_container_name
  default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  # File descriptors 3 and 4 are reserved by Bats.
  run -125 --separate-stderr "$TOOLBX" run --preserve-fds 3 readlink /proc/self/fd/5

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: file descriptor 5 is not available - the preserve-fds option requires that file descriptors must be passed"
  assert_line --index 1 "Error: failed to invoke 'podman exec' in container $default_container_name"
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "run: Try /etc as a command" {
  create_default_container

  run -126 --separate-stderr "$TOOLBX" run /etc

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "bash: line 1: /etc: Is a directory"
  assert_line --index 1 "bash: line 1: exec: /etc: cannot execute: Is a directory"
  assert_line --index 2 "Error: failed to invoke command /etc in container $(get_latest_container_name)"
  assert [ ${#stderr_lines[@]} -eq 3 ]
}

@test "run: Try a non-existent command" {
  local cmd="non-existent-command"

  create_default_container

  run -127 --separate-stderr "$TOOLBX" run "$cmd"

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "bash: line 1: exec: $cmd: not found"
  assert_line --index 1 "Error: command $cmd not found in container $(get_latest_container_name)"
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "run: Try an old unsupported container" {
  local default_image
  default_image="$(get_default_image)"

  pull_default_image

  local container="ancient"

  run "$PODMAN" create --name "$container" "$default_image" true

  assert_success

  run $PODMAN ps --all

  assert_success
  assert_output --regexp "Created[[:blank:]]+$container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --container "$container" true

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: container $container is too old and no longer supported"
  assert_line --index 1 "Recreate it with Toolbx version 0.0.17 or newer."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

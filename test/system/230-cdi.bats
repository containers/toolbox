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

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.7.0
  _setup_environment
  cleanup_all
  pushd "$HOME" || return 1
  rm --force "$XDG_RUNTIME_DIR/toolbox/cdi-nvidia.json" || return 1
}

teardown() {
  rm --force "$XDG_RUNTIME_DIR/toolbox/cdi-nvidia.json" || return 1
  popd || return 1
  cleanup_all
}

@test "cdi: Smoke test" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-empty.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-empty.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    test /etc/ld.so.cache -ot "$toolbx_runtime_directory/cdi-nvidia.json"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/ld.so.conf.d/toolbx-nvidia.conf

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "cdi: ldconfig(8) with no folder" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-00.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    test "$toolbx_runtime_directory/cdi-nvidia.json" -ot /etc/ld.so.cache

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-00.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/ld.so.conf.d/toolbx-nvidia.conf

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "cdi: ldconfig(8) with one folder" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-01.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    test "$toolbx_runtime_directory/cdi-nvidia.json" -ot /etc/ld.so.cache

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-01.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run cat /etc/ld.so.conf.d/toolbx-nvidia.conf

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 ""
  assert_line --index 3 "/usr/lib64"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 4 ]
  else
    assert [ ${#lines[@]} -eq 5 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "cdi: ldconfig(8) with two folders" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-02.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    test "$toolbx_runtime_directory/cdi-nvidia.json" -ot /etc/ld.so.cache

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-02.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run cat /etc/ld.so.conf.d/toolbx-nvidia.conf

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 ""
  assert_line --index 3 "/usr/lib"
  assert_line --index 4 "/usr/lib64"

  if check_bats_version 1.10.0; then
    assert [ ${#lines[@]} -eq 5 ]
  else
    assert [ ${#lines[@]} -eq 6 ]
  fi

  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "cdi: Try invalid JSON" {
  local invalid_json="This is not JSON"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  echo "$invalid_json" >"$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if grep --invert-match --quiet --no-messages "^$invalid_json$" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to load Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "cdi: Try an empty file" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  touch "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if [ -s "$toolbx_runtime_directory/cdi-nvidia.json" ]; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to load Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "cdi: Try hook with invalid path" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-10.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-10.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid hook in Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "cdi: Try hook with unknown path" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-11.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    test /etc/ld.so.cache -ot "$toolbx_runtime_directory/cdi-nvidia.json"

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-11.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/ld.so.conf.d/toolbx-nvidia.conf

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "cdi: Try hook with unknown args" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-12.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    test /etc/ld.so.cache -ot "$toolbx_runtime_directory/cdi-nvidia.json"

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-12.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/ld.so.conf.d/toolbx-nvidia.conf

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "cdi: Try hook with invalid name" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-14.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-14.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid hook in Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "cdi: Try hook with unknown name" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-hooks-15.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    test /etc/ld.so.cache -ot "$toolbx_runtime_directory/cdi-nvidia.json"

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-hooks-15.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/ld.so.conf.d/toolbx-nvidia.conf

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "cdi: Try mount with invalid container path" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-mounts-10.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-mounts-10.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid mount in Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "cdi: Try mount with invalid host path" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-mounts-11.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-mounts-11.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid mount in Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

@test "cdi: Try mount with non-existent paths" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$BATS_TEST_DIRNAME/data/cdi-mounts-12.json" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$BATS_TEST_DIRNAME/data/cdi-mounts-12.json" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /non/existent/path

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

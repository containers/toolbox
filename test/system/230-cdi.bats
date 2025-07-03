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

# bats file_tags=runtime-environment

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  bats_require_minimum_version 1.10.0
  cleanup_all
  pushd "$HOME" || return 1
  rm --force "$XDG_RUNTIME_DIR/toolbox/cdi-nvidia.json" || return 1
}

teardown() {
  rm --force "$XDG_RUNTIME_DIR/toolbox/cdi-nvidia.json" || return 1
  popd || return 1
  cleanup_all
}

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with no link" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-00.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (absolute target, different parent)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-01.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (absolute target, different parent, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-01.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (absolute target, missing parent)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-02.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (absolute target, missing parent, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-02.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (absolute target, same parent)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-03.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (absolute target, same parent, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-03.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with three links (absolute targets, mixed parents)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-04.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with three links (absolute targets, mixed parents, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-04.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "/usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (relative target, different parent)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-05.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (relative target, different parent, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-05.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (relative target, missing parent)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-06.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "../../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (relative target, missing parent, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-06.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "../../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "../../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (relative target, same parent)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-07.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with one link (relative target, same parent, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-07.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with three links (relative targets, mixed parents)" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-08.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "../../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: create-symlinks with three links (relative targets, mixed parents, restart)" {
  local default_container
  default_container="$(get_system_id)-toolbox-$(get_system_version)"

  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-08.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "../../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  podman stop "$default_container"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /run/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /run/toolbox.1

  assert_success
  assert_line --index 0 "../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /opt/bin/toolbox

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /opt/bin/toolbox

  assert_success
  assert_line --index 0 "../../usr/bin/toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run readlink /usr/bin/toolbox.1

  assert_success
  assert_line --index 0 "toolbox"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
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
  assert [ ${#lines[@]} -eq 5 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with invalid path" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-30.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid hook in Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with unknown path" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-31.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with missing --link argument" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-32.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to create symlinks for Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with relative link in --link argument" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-33.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to create symlinks for Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with wrongly formatted --link argument ('foo')" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-34.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to create symlinks for Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with wrongly formatted --link argument ('foo::bar::baz')" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-35.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: failed to create symlinks for Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with invalid name" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-36.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run true

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid hook in Container Device Interface for NVIDIA"
  assert [ ${#stderr_lines[@]} -eq 1 ]
}

# bats test_tags=arch-fedora
@test "cdi: Try create-symlinks with unknown name" {
  local test_cdi_file="$BATS_TEST_DIRNAME/data/cdi-hooks-create-symlinks-37.json"
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$toolbx_runtime_directory"

  cp "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"
  chmod 644 "$toolbx_runtime_directory/cdi-nvidia.json"

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /usr/bin/toolbox.1

  if ! cmp --silent "$test_cdi_file" "$toolbx_runtime_directory/cdi-nvidia.json"; then
    skip "found NVIDIA hardware"
  fi

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

# bats test_tags=arch-fedora
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

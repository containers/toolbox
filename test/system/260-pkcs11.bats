# shellcheck shell=bats
#
# Copyright © 2025 Red Hat, Inc.
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
  _setup_environment
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

# bats test_tags=arch-fedora
@test "pkcs11: Smoke test" {
  local pkcs11_directory_absent=false
  local pkcs11_skip=false

  local system_id
  system_id="$(get_system_id)"

  local system_version
  system_version="$(get_system_version)"

  if [ "$system_id" = "fedora" ]; then
    pkcs11_skip=true
  elif [ "$system_id" = "ubuntu" ] && [ "$system_version" = "22.04" ]; then
    pkcs11_directory_absent=true
    pkcs11_skip=true
  fi

  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_default_container

  run --keep-empty-lines --separate-stderr "$TOOLBX" run test -S "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]

  # shellcheck disable=SC2154
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run sh -c 'echo "$P11_KIT_SERVER_ADDRESS"'

  assert_success
  assert_line --index 0 "unix:path=$toolbx_runtime_directory/pkcs11"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  if $pkcs11_skip; then
    if $pkcs11_directory_absent; then
      run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/pkcs11/modules

      assert_failure
      assert [ ${#lines[@]} -eq 0 ]
      assert [ ${#stderr_lines[@]} -eq 0 ]
    else
      run --keep-empty-lines --separate-stderr "$TOOLBX" run test -d /etc/pkcs11/modules

      assert_success
      assert [ ${#lines[@]} -eq 0 ]
      assert [ ${#stderr_lines[@]} -eq 0 ]
    fi

    run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/pkcs11/modules/p11-kit-trust.module

    assert_failure
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  else
    run --keep-empty-lines --separate-stderr "$TOOLBX" run cat /etc/pkcs11/modules/p11-kit-trust.module

    assert_success
    assert_line --index 0 "# Written by Toolbx"
    assert_line --index 1 "# https://containertoolbx.org/"
    assert_line --index 2 ""
    assert_line --index 3 "module: p11-kit-client.so"
    assert [ ${#lines[@]} -eq 4 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]

    run --keep-empty-lines --separate-stderr "$TOOLBX" run stat \
                                                             --format "%A %U:%G" \
                                                             /etc/pkcs11/modules/p11-kit-trust.module

    assert_success
    assert_line --index 0 "-rw-r--r-- root:root"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi

  run --keep-empty-lines --separate-stderr lsof -Fcfp "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert_line --index 0 --regexp '^p[[:digit:]]+$'
  assert_line --index 1 "cp11-kit-server"

  local pid="${lines[0]#?}"
  kill "$pid"

  assert_line --index 2 --regexp '^f[[:digit:]]+$'
  assert [ ${#lines[@]} -eq 3 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "pkcs11: Smoke test with Arch Linux" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_distro_container arch latest arch-toolbox-latest

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro arch \
    test -S "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro arch \
    sh -c 'echo "$P11_KIT_SERVER_ADDRESS"'

  assert_success
  assert_line --index 0 "unix:path=$toolbx_runtime_directory/pkcs11"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro arch \
    cat /etc/pkcs11/modules/p11-kit-trust.module

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 ""
  assert_line --index 3 "module: p11-kit-client.so"
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro arch \
    stat \
      --format "%A %U:%G" \
      /etc/pkcs11/modules/p11-kit-trust.module

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr lsof -Fcfp "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert_line --index 0 --regexp '^p[[:digit:]]+$'
  assert_line --index 1 "cp11-kit-server"

  local pid="${lines[0]#?}"
  kill "$pid"

  assert_line --index 2 --regexp '^f[[:digit:]]+$'
  assert [ ${#lines[@]} -eq 3 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "pkcs11: Smoke test with Fedora 34" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_distro_container fedora 34 fedora-toolbox-34

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    test -S "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    sh -c 'echo "$P11_KIT_SERVER_ADDRESS"'

  assert_success
  assert_line --index 0 "unix:path=$toolbx_runtime_directory/pkcs11"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    test -d /etc/pkcs11/modules

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    test -e /etc/pkcs11/modules/p11-kit-trust.module

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr lsof -Fcfp "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert_line --index 0 --regexp '^p[[:digit:]]+$'
  assert_line --index 1 "cp11-kit-server"

  local pid="${lines[0]#?}"
  kill "$pid"

  assert_line --index 2 --regexp '^f[[:digit:]]+$'
  assert [ ${#lines[@]} -eq 3 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "pkcs11: Smoke test with RHEL 8.10" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_distro_container rhel 8.10 rhel-toolbox-8.10

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    test -S "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    sh -c 'echo "$P11_KIT_SERVER_ADDRESS"'

  assert_success
  assert_line --index 0 "unix:path=$toolbx_runtime_directory/pkcs11"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    test -d /etc/pkcs11/modules

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    test -e /etc/pkcs11/modules/p11-kit-trust.module

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr lsof -Fcfp "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert_line --index 0 --regexp '^p[[:digit:]]+$'
  assert_line --index 1 "cp11-kit-server"

  local pid="${lines[0]#?}"
  kill "$pid"

  assert_line --index 2 --regexp '^f[[:digit:]]+$'
  assert [ ${#lines[@]} -eq 3 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "pkcs11: Smoke test with Ubuntu 16.04" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 16.04 \
    test -S "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 16.04 \
    sh -c 'echo "$P11_KIT_SERVER_ADDRESS"'

  assert_success
  assert_line --index 0 "unix:path=$toolbx_runtime_directory/pkcs11"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 16.04 \
    test -e /etc/pkcs11/modules

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 16.04 \
    test -e /etc/pkcs11/modules/p11-kit-trust.module

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr lsof -Fcfp "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert_line --index 0 --regexp '^p[[:digit:]]+$'
  assert_line --index 1 "cp11-kit-server"

  local pid="${lines[0]#?}"
  kill "$pid"

  assert_line --index 2 --regexp '^f[[:digit:]]+$'
  assert [ ${#lines[@]} -eq 3 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "pkcs11: Smoke test with Ubuntu 18.04" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 18.04 \
    test -S "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 18.04 \
    sh -c 'echo "$P11_KIT_SERVER_ADDRESS"'

  assert_success
  assert_line --index 0 "unix:path=$toolbx_runtime_directory/pkcs11"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 18.04 \
    cat /etc/pkcs11/modules/p11-kit-trust.module

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 ""
  assert_line --index 3 "module: p11-kit-client.so"
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 18.04 \
    stat \
      --format "%A %U:%G" \
      /etc/pkcs11/modules/p11-kit-trust.module

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr lsof -Fcfp "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert_line --index 0 --regexp '^p[[:digit:]]+$'
  assert_line --index 1 "cp11-kit-server"

  local pid="${lines[0]#?}"
  kill "$pid"

  assert_line --index 2 --regexp '^f[[:digit:]]+$'
  assert [ ${#lines[@]} -eq 3 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "pkcs11: Smoke test with Ubuntu 20.04" {
  local toolbx_runtime_directory="$XDG_RUNTIME_DIR/toolbox"

  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 20.04 \
    test -S "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  # shellcheck disable=SC2016
  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 20.04 \
    sh -c 'echo "$P11_KIT_SERVER_ADDRESS"'

  assert_success
  assert_line --index 0 "unix:path=$toolbx_runtime_directory/pkcs11"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 20.04 \
    cat /etc/pkcs11/modules/p11-kit-trust.module

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 ""
  assert_line --index 3 "module: p11-kit-client.so"
  assert [ ${#lines[@]} -eq 4 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 20.04 \
    stat \
      --format "%A %U:%G" \
      /etc/pkcs11/modules/p11-kit-trust.module

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr lsof -Fcfp "$toolbx_runtime_directory/pkcs11"

  assert_success
  assert_line --index 0 --regexp '^p[[:digit:]]+$'
  assert_line --index 1 "cp11-kit-server"

  local pid="${lines[0]#?}"
  kill "$pid"

  assert_line --index 2 --regexp '^f[[:digit:]]+$'
  assert [ ${#lines[@]} -eq 3 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

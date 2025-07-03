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
  cleanup_all
  pushd "$HOME" || return 1
}

teardown() {
  popd || return 1
  cleanup_all
}

# bats test_tags=arch-fedora
@test "kerberos: Smoke test" {
  local kerberos_skip=false

  local system_id
  system_id="$(get_system_id)"

  if [ "$system_id" != "fedora" ]; then
    kerberos_skip=true
  fi

  create_default_container

  if $kerberos_skip; then
    run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/krb5.conf.d

    assert_failure
    assert [ ${#lines[@]} -eq 0 ]
    # shellcheck disable=SC2154
    assert [ ${#stderr_lines[@]} -eq 0 ]

    run --keep-empty-lines --separate-stderr "$TOOLBX" run test -e /etc/krb5.conf.d/kcm_default_ccache

    assert_failure
    assert [ ${#lines[@]} -eq 0 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  else
    run --keep-empty-lines --separate-stderr "$TOOLBX" run cat /etc/krb5.conf.d/kcm_default_ccache

    assert_success
    assert_line --index 0 "# Written by Toolbx"
    assert_line --index 1 "# https://containertoolbx.org/"
    assert_line --index 2 "#"
    assert_line --index 3 "# # To disable the KCM credential cache, comment out the following lines."
    assert_line --index 4 ""
    assert_line --index 5 "[libdefaults]"
    assert_line --index 6 "    default_ccache_name = KCM:"
    assert [ ${#lines[@]} -eq 7 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]

    run --keep-empty-lines --separate-stderr "$TOOLBX" run stat \
                                                             --format "%A %U:%G" \
                                                             /etc/krb5.conf.d/kcm_default_ccache

    assert_success
    assert_line --index 0 "-rw-r--r-- root:root"
    assert [ ${#lines[@]} -eq 1 ]
    assert [ ${#stderr_lines[@]} -eq 0 ]
  fi
}

# bats test_tags=arch-fedora
@test "kerberos: Smoke test with Arch Linux" {
  create_distro_container arch latest arch-toolbox-latest

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch test -e /etc/krb5.conf.d

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro arch test -e /etc/krb5.conf.d/kcm_default_ccache

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "kerberos: Smoke test with Fedora 34" {
  create_distro_container fedora 34 fedora-toolbox-34

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    cat /etc/krb5.conf.d/kcm_default_ccache

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 "#"
  assert_line --index 3 "# # To disable the KCM credential cache, comment out the following lines."
  assert_line --index 4 ""
  assert_line --index 5 "[libdefaults]"
  assert_line --index 6 "    default_ccache_name = KCM:"
  assert [ ${#lines[@]} -eq 7 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro fedora \
    --release 34 \
    stat \
      --format "%A %U:%G" \
      /etc/krb5.conf.d/kcm_default_ccache

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=arch-fedora
@test "kerberos: Smoke test with RHEL 8.10" {
  create_distro_container rhel 8.10 rhel-toolbox-8.10

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    cat /etc/krb5.conf.d/kcm_default_ccache

  assert_success
  assert_line --index 0 "# Written by Toolbx"
  assert_line --index 1 "# https://containertoolbx.org/"
  assert_line --index 2 "#"
  assert_line --index 3 "# # To disable the KCM credential cache, comment out the following lines."
  assert_line --index 4 ""
  assert_line --index 5 "[libdefaults]"
  assert_line --index 6 "    default_ccache_name = KCM:"
  assert [ ${#lines[@]} -eq 7 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro rhel \
    --release 8.10 \
    stat \
      --format "%A %U:%G" \
      /etc/krb5.conf.d/kcm_default_ccache

  assert_success
  assert_line --index 0 "-rw-r--r-- root:root"
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "kerberos: Smoke test with Ubuntu 16.04" {
  create_distro_container ubuntu 16.04 ubuntu-toolbox-16.04

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 16.04 test -e /etc/krb5.conf.d

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 16.04 \
    test -e /etc/krb5.conf.d/kcm_default_ccache

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "kerberos: Smoke test with Ubuntu 18.04" {
  create_distro_container ubuntu 18.04 ubuntu-toolbox-18.04

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 18.04 test -e /etc/krb5.conf.d

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 18.04 \
    test -e /etc/krb5.conf.d/kcm_default_ccache

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

# bats test_tags=ubuntu
@test "kerberos: Smoke test with Ubuntu 20.04" {
  create_distro_container ubuntu 20.04 ubuntu-toolbox-20.04

  run --keep-empty-lines --separate-stderr "$TOOLBX" run --distro ubuntu --release 20.04 test -e /etc/krb5.conf.d

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]

  run --keep-empty-lines --separate-stderr "$TOOLBX" run \
    --distro ubuntu \
    --release 20.04 \
    test -e /etc/krb5.conf.d/kcm_default_ccache

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

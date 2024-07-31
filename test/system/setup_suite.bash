# shellcheck shell=bash
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

missing_dependencies=false

if [ -f "$BATS_TEST_DIRNAME/libs/bats-assert/load.bash" ] && [ -f "$BATS_TEST_DIRNAME/libs/bats-support/load.bash" ]; then
  load 'libs/helpers'
else
  missing_dependencies=true
fi

setup_suite() {
  bats_require_minimum_version 1.7.0
  echo "# test suite: Set up" >&3

  if $missing_dependencies; then
    echo "# Missing dependencies" >&3
    echo "# Forgot to run 'git submodule init' and 'git submodule update' ?" >&3
    return 1
  fi

  local system_id
  system_id="$(get_system_id)"

  local system_version
  system_version="$(get_system_version)"

  _setup_environment
  # Cache the default image for the system
  _pull_and_cache_distro_image "$system_id" "$system_version" || false
  # Cache all images that will be needed during the tests
  _pull_and_cache_distro_image arch latest || false
  _pull_and_cache_distro_image fedora 34 || false
  _pull_and_cache_distro_image rhel 8.10 || false
  _pull_and_cache_distro_image ubuntu 16.04 || false
  _pull_and_cache_distro_image ubuntu 18.04 || false
  _pull_and_cache_distro_image ubuntu 20.04 || false
  _pull_and_cache_distro_image busybox || false
  # If run on Fedora Rawhide, cache 2 extra images (previous Fedora versions)
  if is_fedora_rawhide; then
    _pull_and_cache_distro_image fedora "$((system_version-1))" || false
    _pull_and_cache_distro_image fedora "$((system_version-2))" || false
  fi

  _setup_docker_registry
}

teardown_suite() {
  bats_require_minimum_version 1.7.0
  echo "# test suite: Tear down" >&3

  if $missing_dependencies; then
    return 0
  fi

  _setup_environment

  _clean_cached_images
  _clean_docker_registry
  _clean_temporary_storage
}

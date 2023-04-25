#!/usr/bin/env bats
#
# Copyright © 2019 – 2022 Red Hat, Inc.
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
load 'libs/helpers.bash'

setup() {
  _setup_environment
}

@test "version: Check version using option --version" {
  bats_require_minimum_version 1.5.0
  run --separate-stderr $TOOLBOX --version

  assert_success
  assert_output --regexp '^toolbox version [0-9]+\.[0-9]+\.[0-9]+(\.[0-9]+)?$'
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

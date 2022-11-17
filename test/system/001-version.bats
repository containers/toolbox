#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers.bash'

setup() {
  _setup_environment
}

@test "version: Check version using option --version" {
  run --separate-stderr $TOOLBOX --version

  assert_success
  assert_output --regexp '^toolbox version [0-9]+\.[0-9]+\.[0-9]+(\.[0-9]+)?$'
  assert [ ${#lines[@]} -eq 1 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

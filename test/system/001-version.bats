#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers.bash'

setup() {
  _setup_environment
}

@test "version: Check version using option --version" {
  run $TOOLBOX --version

  assert_output --regexp '^toolbox version [0-9]+\.[0-9]+\.[0-9]+(\.[0-9]+)?$'
}

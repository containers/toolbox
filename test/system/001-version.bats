#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'

@test "version: Check version using option --version" {
  run toolbox --version

  assert_output --regexp '^toolbox version [0-9]+\.[0-9]+\.[0-9]+$'
}

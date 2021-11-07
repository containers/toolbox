#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  check_xdg_runtime_dir
  cleanup_containers
}

teardown() {
  cleanup_containers
}

@test "root: Try to enter the default container with no containers created" {
  run $TOOLBOX <<< "n"

  assert_success
  assert_line --index 0 "No toolbox containers found. Create now? [y/N] A container can be created later with the 'create' command."
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

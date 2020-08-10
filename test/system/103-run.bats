#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

# It seems like 'toolbox run' (or 'enter') doesn't work fine when
# the workdir is outside the $HOME.
# This hack is to make the tests work from outside the $HOME.
readonly CURDIR=$PWD

setup() {
  cd "$HOME" || return 1
  cleanup_containers
}

teardown() {
  cleanup_containers
  cd "$CURDIR" || return 1
}


@test "run: Try to run echo 'Hello World' with no containers created" {
  run toolbox run echo "Hello World"

  assert_failure
  assert_line --index 0 --regexp 'Error: container .* not found'
  assert_output --partial "Run 'toolbox --help' for usage."
}

#TODO: This should work without --partial
# The issue here is that toolbox output add the CRLF character at the end
@test "run: Run echo 'Hello World' inside of the default container" {
  create_default_container

  run toolbox --verbose run echo "Hello World"

  assert_success
  assert_output --partial "Hello World"
}

@test "run: Run echo 'Hello World' inside a container after being stopped" {
  create_container running

  start_container running
  stop_container running

  run toolbox --verbose run --container running echo -n "Hello World"

  assert_success
  assert_output --partial "Hello World"
}

@test "run: Run sudo id inside of the default container" {
  create_default_container

  run toolbox --verbose run sudo id

  assert_success
  assert_output --partial "uid=0(root)"
}

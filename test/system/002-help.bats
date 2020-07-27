#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'

@test "help: Try to run toolbox with no command (shows usage screen)" {
  run toolbox

  assert_failure
  assert_line --index 0 "Error: missing command"
  assert_output --partial "Run 'toolbox --help' for usage."
}

@test "help: Run command 'help'" {
  run toolbox help

  assert_success
  assert_output --partial "toolbox - Unprivileged development environment"
}

@test "help: Use flag '--help' (it should show usage screen)" {
  run toolbox --help

  assert_success
  assert_output --partial "toolbox - Unprivileged development environment"
}

@test "help: Try to run toolbox with non-existent command (shows usage screen)" {
  run toolbox foo

  assert_failure
  assert_line --index 0 "Error: unknown command \"foo\" for \"toolbox\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers.bash'

setup() {
  check_xdg_runtime_dir
}

@test "help: Run command 'help'" {
  run $TOOLBOX help

  assert_success
  assert_output --partial "toolbox - Tool for containerized command line environments on Linux"
}

@test "help: Run command 'help' with no man present" {
  if hash man 2>/dev/null; then
    skip "Test works only if man is not in PATH"
  fi

  run $TOOLBOX help

  assert_success
  assert_line --index 0 "toolbox - Tool for containerized command line environments on Linux"
  assert_line --index 1 "Common commands are:"
  assert_line --index 2 "create    Create a new toolbox container"
  assert_line --index 3 "enter     Enter an existing toolbox container"
  assert_line --index 4 "list      List all existing toolbox containers and images"
  assert_line --index 5 "Go to https://github.com/containers/toolbox for further information."
}

@test "help: Use flag '--help' (it should show usage screen)" {
  run $TOOLBOX --help

  assert_success
  assert_output --partial "toolbox - Tool for containerized command line environments on Linux"
}

@test "help: Try to run toolbox with non-existent command (shows usage screen)" {
  run $TOOLBOX foo

  assert_failure
  assert_line --index 0 "Error: unknown command \"foo\" for \"toolbox\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run toolbox with non-existent flag (shows usage screen)" {
  run $TOOLBOX --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox create' with non-existent flag (shows usage screen)" {
  run $TOOLBOX create --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox enter' with non-existent flag (shows usage screen)" {
  run $TOOLBOX enter --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox help' with non-existent flag (shows usage screen)" {
  run $TOOLBOX help --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox init-container' with non-existent flag (shows usage screen)" {
  run $TOOLBOX init-container --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox list' with non-existent flag (shows usage screen)" {
  run $TOOLBOX list --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox rm' with non-existent flag (shows usage screen)" {
  run $TOOLBOX rm --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox rmi' with non-existent flag (shows usage screen)" {
  run $TOOLBOX rmi --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

@test "help: Try to run 'toolbox run' with non-existent flag (shows usage screen)" {
  run $TOOLBOX run --foo

  assert_failure
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
}

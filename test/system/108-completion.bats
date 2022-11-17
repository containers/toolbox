#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  _setup_environment
}

@test "completion: Smoke test with 'bash'" {
  run $TOOLBOX completion bash

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "completion: Smoke test with 'fish'" {
  run $TOOLBOX completion fish

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "completion: Smoke test with 'zsh'" {
  run $TOOLBOX completion zsh

  assert_success
  assert [ ${#lines[@]} -gt 0 ]
  assert [ ${#stderr_lines[@]} -eq 0 ]
}

@test "completion: Try without any arguments" {
  run --separate-stderr $TOOLBOX completion

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: accepts 1 arg(s), received 0"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "completion: Try with invalid arguments" {
  run --separate-stderr $TOOLBOX completion foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument \"foo\" for \"toolbox completion\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "completion: Try with unknown flag" {
  run --separate-stderr $TOOLBOX completion --foo

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: unknown flag: --foo"
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

@test "completion: Try with unsupported shell" {
  run --separate-stderr $TOOLBOX completion powershell

  assert_failure
  assert [ ${#lines[@]} -eq 0 ]
  lines=("${stderr_lines[@]}")
  assert_line --index 0 "Error: invalid argument \"powershell\" for \"toolbox completion\""
  assert_line --index 1 "Run 'toolbox --help' for usage."
  assert [ ${#stderr_lines[@]} -eq 2 ]
}

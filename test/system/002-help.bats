#!/usr/bin/env bats

load helpers

@test "Show usage screen when no command is given" {
  run_toolbox 1
  is "${lines[0]}" "Error: missing command" "Usage line 1"
}

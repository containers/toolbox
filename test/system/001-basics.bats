#!/usr/bin/env bats

load helpers

function setup() {
  :
}

function teardown() {
  :
}

@test "Output version number using full flag" {
  skip "Not implemented"
  run_toolbox --version
}

@test "Output version number using command" {
  skip "Not implemented"
  run_toolbox version
}

@test "Show usage screen when no command is given" {
  run_toolbox 1
  is "${lines[0]}" "toolbox: missing command" "Usage line 1"
}

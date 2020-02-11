#!/usr/bin/env bats

load helpers

function setup() {
  setup_with_one_container
}

@test "Echo 'Hello World' inside of an container" {
  run_toolbox run echo "Hello World"
  is "$output" "Hello World" "Should say 'Hello World'"
}

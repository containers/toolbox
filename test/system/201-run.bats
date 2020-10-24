#!/usr/bin/env bats

load helpers

@test "Start the 'running' container and check it started alright" {
  run_podman --log-level debug start running

  is_toolbox_ready running
}

@test "Echo 'Hello World' inside of the default container" {
  run_toolbox run echo "Hello World"
  # is "$output" "Hello World" "Should say 'Hello World'"
}

@test "Echo 'Hello World' inside of the 'running' container" {
  run_toolbox run -c running echo "Hello World"
  # is "$output" "Hello World" "Should say 'Hello World'"
}

@test "Stop the 'running' container using 'podman stop'" {
  run_podman stop running
  is "${#lines[@]}" "1" "Expected number of lines of the output is 1 (with the id of the container)"
}

@test "Echo 'hello World' again in the 'running' container after being stopped and exit" {
  run_toolbox run -c running echo "Hello World"
  # is "$output" "Hello World" "Should say 'Hello World'"
}

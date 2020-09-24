#!/usr/bin/env bats

load helpers

@test "Start the 'running' container" {
  run_podman --log-level debug start running
}

@test "Logs of container 'running' look alright" {
  run_podman logs running
  is "${lines[${#lines[@]} - 1]}" "level=debug msg=\"Going to sleep\"" "The last line of the logs should say the entry-point went to sleep"
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

#!/usr/bin/env bats

load helpers

@test "Try to remove a nonexistent container" {
  run_toolbox rm nonexistentcontainer
  is "$lines[0]" "Error: failed to inspect container nonexistentcontainer" "Toolbox should fail to remove a non-existent container"
}

@test "Try to remove the running container 'running'" {
  run_toolbox rm running
  is "$output" "Error: container running is running" "Toolbox should fail to remove a running container"
}

@test "Remove the not running container 'not-running'" {
  run_toolbox rm not-running
  is "$output" "" "The output should be empty"
}

@test "Force remove the running container 'running'" {
  run_toolbox rm --force running
  is "$output" "" "The output should be empty"
}

@test "Force remove all remaining containers (only 1 should be left)" {
  run_toolbox rm --force --all
  is "$output" "" "The output should be empty"
}

#!/usr/bin/env bats

load helpers

@test "Try to remove a nonexistent container" {
  run_toolbox 1 rm nonexistentcontainer
  is "$output" "toolbox: failed to inspect $todelete" "Toolbox should fail with: no such container"
}

@test "Try to remove the running container 'running'" {
  run_toolbox 1 rm running
  is "$output" "toolbox: failed to remove container running" "Toolbox should fail to remove the running container"
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

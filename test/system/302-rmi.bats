#!/usr/bin/env bats

load helpers

@test "Remove all images (2 should be present; --force should not be necessary)" {
  run_toolbox rmi --all
  is "$output" "" "The output should be empty"
}

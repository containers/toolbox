#!/usr/bin/env bats

load helpers

@test "Create the default container" {
  run_toolbox -y create
}

@test "Create a container with a valid custom name ('not-running'; name passed as an argument)" {
  run_toolbox -y create "not-running"
}

@test "Create a container with a custom image and name ('running';f29; name passed as an option)" {
  run_toolbox -y create -c "running" -i fedora-toolbox:29
}

@test "Try to create a container with invalid custom name" {
  run_toolbox 1 -y create "ßpeci@l.Nam€"
  is "${lines[0]}" "Error: failed to create container ßpeci@l.Nam€" "Toolbox should fail to create a container with such name"
}

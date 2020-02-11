#!/usr/bin/env bats

load helpers

@test "Create the default container." {
  run_toolbox -y create
}

@test "Create a container with a valid custom name (whole word)" {
  run_toolbox -y create -c "customname"
}

@test "Try to create a container with a bad custom name (with special characters)" {
  run_toolbox 1 -y create -c "ßpeci@l.Nam€"
  is "${lines[0]}" "toolbox: invalid argument for '--container'" "Toolbox reports invalid argument for --container"
}

@test "Create a container with a custom image (f29)" {
  run_toolbox -y create -i fedora-toolbox:29
}

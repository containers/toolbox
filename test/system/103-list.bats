#!/usr/bin/env bats

load helpers

@test "Run list with three containers and two images" {
  run_toolbox list
  is "${#lines[@]}" "7" "Expected number of lines of the output is 7 (Img: 3 + Cont: 4)"

  is "${lines[1]}" ".*registry.fedoraproject.org/.*" "The first of the two images"
  is "${lines[2]}" ".*registry.fedoraproject.org/.*" "The second of the two images"

  is "${lines[4]}" ".*fedora-toolbox-.*" "The default container should be first in the list"
  is "${lines[5]}" ".*not-running.*" "The container 'not-running' should be second"
  is "${lines[6]}" ".*running.*" "The container 'running' should be third (last)"
}

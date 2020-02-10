#!/usr/bin/env bats

load helpers

@test "Run list with zero containers and two images" {
  run_toolbox list
  is "${#lines[@]}" "3" "Expected number of lines of the output is 3 (Img: 3 + Spc: 0 + Cont: 0)"

  is "${lines[1]}" ".*registry.fedoraproject.org/.*" "First of the two images"
  is "${lines[2]}" ".*registry.fedoraproject.org/.*" "Second of the two images"
}

@test "Run list with zero containers (-c flag)" {
  run_toolbox list -c
  is "$output" "" "Output of list should be blank"
}

@test "Run list with zero images (-i flag)" {
  run_toolbox list -i
  is "${#lines[@]}" "3" "Expected number of lines of the output is 3"

  is "${lines[1]}" ".*registry.fedoraproject.org/.*" "First of the two images"
  is "${lines[2]}" ".*registry.fedoraproject.org/.*" "Second of the two images"
}

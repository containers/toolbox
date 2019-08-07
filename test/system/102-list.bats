#!/usr/bin/env bats

load helpers

@test "Run list with zero containers and zero images" {
  remove_all_images
  remove_all_containers
  run_toolbox list
  is "$output" "" "Output of list should be blank"
}

@test "Run list with zero containers (-c flag)" {
  remove_all_containers
  run_toolbox list -c
  is "$output" "" "Output of list should be blank"
}

@test "Run list with zero images (-i flag)" {
  remove_all_images
  run_toolbox list -i
  is "$output" "" "Output of list should be blank"
}

@test "Run list with 1 default container and 1 default image" {
  create_toolbox
  run_toolbox list
  is "${lines[1]}" ".*registry.fedoraproject.org/.*" "Default image"
  is "${lines[3]}" ".*fedora-toolbox-.*" "Default container"
  is "${#lines[@]}" "4" "Expected length of output is 4"
}

@test "Run list with 3 containers (-c flag)" {
  create_toolbox 3 fedora
  run_toolbox list -c
  for i in $(seq 1 3); do
    is "${lines[$i]}" ".*fedora-$((i)) \+" "One of the containers"
  done
}

@test "Run list with 3 images (-i flag)" {
  get_images 3
  run_toolbox list -i
  is "${#lines[@]}" "4" "Expected length of output is 4"
}

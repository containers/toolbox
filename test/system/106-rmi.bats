#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  cleanup_all
}

teardown() {
  cleanup_all
}


@test "rmi: Remove all images with the default image present" {
  num_of_images=$(list_images)
  assert_equal "$num_of_images" 0

  pull_default_image

  run $TOOLBOX rmi --all

  assert_success
  assert_output ""

  new_num_of_images=$(list_images)

  assert_equal "$new_num_of_images" "$num_of_images"
}

@test "rmi: Try to remove all images with a container present and running" {
  skip "Bug: Fail in 'toolbox rmi' does not return non-zero value"
  num_of_images=$(list_images)
  assert_equal "$num_of_images" 0

  create_container foo
  start_container foo

  run $TOOLBOX rmi --all

  assert_failure
  assert_output --regexp "Error: image .* has dependent children"

  new_num_of_images=$(list_images)

  assert_equal "$new_num_of_images" "$num_of_images"
}

@test "rmi: Force remove all images with a container present and running" {
  num_of_images=$(list_images)
  assert_equal "$num_of_images" 0

  create_container foo
  start_container foo

  run $TOOLBOX rmi --all --force

  assert_success
  assert_output ""

  new_num_of_images=$(list_images)

  assert_equal "$new_num_of_images" "$num_of_images"
}

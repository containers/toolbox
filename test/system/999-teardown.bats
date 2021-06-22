#!/usr/bin/env bats

load 'libs/helpers'

@test "test suite: Teardown" {
  _clean_cached_images
  _clean_temporary_storage
}

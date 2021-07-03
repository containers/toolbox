#!/usr/bin/env bats

load 'libs/helpers'

@test "test suite: Teardown" {
  _restore_toolbox_configs
  _clean_cached_images
}

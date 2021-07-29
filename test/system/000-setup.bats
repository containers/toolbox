#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/helpers'

@test "test suite: Setup" {
    _setup_runtime_dir || fail "Failed to setup runtime dir ${RUNTIME_DIR}"
    # Cache the default image for the system
    _pull_and_cache_distro_image $(get_system_id) $(get_system_version) || die
    # Cache all images that will be needed during the tests
    _pull_and_cache_distro_image fedora 32 || die
    _pull_and_cache_distro_image busybox || die
}

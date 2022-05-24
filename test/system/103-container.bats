#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  _setup_environment
  cleanup_containers
}

teardown() {
  cleanup_containers
}


@test "container: Check container starts without issues" {
  readonly CONTAINER_NAME="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  run container_started $CONTAINER_NAME
  assert_success
}

@test "container(Fedora Rawhide): Containers with supported versions start without issues" {
  local os_release="$(find_os_release)"
  local system_id="$(get_system_id)"
  local system_version="$(get_system_version)"
  local rawhide_res="$(awk '/rawhide/' $os_release)"

  assert [ -n "$os_release" ]

  if [ "$system_id" != "fedora" ] || [ -z "$rawhide_res" ]; then
    skip "This test is only for Fedora Rawhide"
  fi

  create_distro_container "$system_id" "$system_version" latest
  run container_started latest
  assert_success

  create_distro_container "$system_id" "$((system_version-1))" second
  run container_started second
  assert_success

  create_distro_container "$system_id" "$((system_version-2))" third
  run container_started third
  assert_success
}

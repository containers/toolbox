#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  cleanup_containers
}

teardown() {
  cleanup_containers
}


@test "container: Check container starts without issues" {
  readonly CONTAINER_NAME="$(get_system_id)-toolbox-$(get_system_version)"

  create_default_container

  run $PODMAN start $CONTAINER_NAME

  CONTAINER_INITIALIZED=0

  for TRIES in 1 2 3 4 5
  do
    run $PODMAN logs $CONTAINER_NAME
    CONTAINER_OUTPUT=$output
    run grep 'Listening to file system and ticker events' <<< $CONTAINER_OUTPUT
    if [[ "$status" -eq 0 ]]; then
      CONTAINER_INITIALIZED=1
      break
    fi
    sleep 1
  done

  echo $CONTAINER_OUTPUT
  assert [ "$CONTAINER_INITIALIZED" -eq 1 ]
}


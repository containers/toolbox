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

@test "run: Try to run a command in the default container with no containers created" {
  local default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  run $TOOLBOX run true

  assert_failure
  assert_line --index 0 "Error: container $default_container_name not found"
  assert_line --index 1 "Use the 'create' command to create a toolbox."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "run: Try to run a command in the default container when 1 non-default container is present" {
  local default_container_name="$(get_system_id)-toolbox-$(get_system_version)"

  create_container other-container

  run $TOOLBOX run true

  assert_failure
  assert_line --index 0 "Error: container $default_container_name not found"
  assert_line --index 1 "Use the 'create' command to create a toolbox."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "run: Try to run a command in a specific non-existent container" {
  create_container other-container

  run $TOOLBOX run -c wrong-container true

  assert_failure
  assert_line --index 0 "Error: container wrong-container not found"
  assert_line --index 1 "Use the 'create' command to create a toolbox."
  assert_line --index 2 "Run 'toolbox --help' for usage."
}

@test "run: Run echo 'Hello World' inside of the default container" {
  create_default_container

  run $TOOLBOX --verbose run echo -n "Hello World"

  assert_success
  assert_line --index $((${#lines[@]}-1)) "Hello World"
}

@test "run: Run echo 'Hello World' inside a container after being stopped" {
  create_container running

  start_container running
  stop_container running

  run $TOOLBOX --verbose run --container running echo -n "Hello World"

  assert_success
  assert_line --index $((${#lines[@]}-1)) "Hello World"
}

@test "run: Run sudo id inside of the default container" {
  create_default_container

  output="$($TOOLBOX --verbose run sudo id 2>$BATS_TMPDIR/stderr)"
  status="$?"

  echo "# stderr"
  cat $BATS_TMPDIR/stderr
  echo "# stdout"
  echo $output

  assert_success
  assert_output --partial "uid=0(root)"
}

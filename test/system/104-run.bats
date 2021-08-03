#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

setup() {
  check_xdg_runtime_dir
  cleanup_containers
}

teardown() {
  cleanup_containers
}


@test "run: Try to run echo 'Hello World' with no containers created" {
  run $TOOLBOX run echo "Hello World"

  assert_failure
  assert_line --index 0 --regexp 'Error: container .* not found'
  assert_output --partial "Run 'toolbox --help' for usage."
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

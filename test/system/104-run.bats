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

@test "run: Run command exiting with zero code in the default container" {
  create_default_container

  run $TOOLBOX run true

  assert_success
  assert_output ""
}

@test "run: Ensure that a login shell is used to invoke the command" {
  create_default_container

  cp "$HOME"/.bash_profile "$HOME"/.bash_profile.orig
  echo "echo \"~/.bash_profile read\"" >>"$HOME"/.bash_profile

  run $TOOLBOX run true

  mv "$HOME"/.bash_profile.orig "$HOME"/.bash_profile

  assert_success
  assert_line --index 0 "~/.bash_profile read"
  assert [ ${#lines[@]} -eq 1 ]
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

@test "run: Ensure that $HOME is used as a fallback working directory" {
  local default_container_name="$(get_system_id)-toolbox-$(get_system_version)"
  create_default_container

  pushd /etc/kernel
  run --separate-stderr $TOOLBOX run pwd
  popd

  assert_success
  assert_line --index 0 "$HOME"
  assert [ ${#lines[@]} -eq 1 ]
  lines=("${stderr_lines[@]}")
  assert_line --index $((${#stderr_lines[@]}-2)) "Error: directory /etc/kernel not found in container $default_container_name"
  assert_line --index $((${#stderr_lines[@]}-1)) "Using $HOME instead."
  assert [ ${#stderr_lines[@]} -gt 2 ]
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

@test "run: Try to run a command in a container based on unsupported distribution" {
  local distro="foo"

  run $TOOLBOX --assumeyes run --distro "$distro" ls

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--distro'"
  assert_line --index 1 "Distribution $distro is unsupported."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "run: Try to run a command in a container based on Fedora but with wrong version" {
  run $TOOLBOX run -d fedora -r foobar ls

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX run --distro fedora --release -3 ls

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive integer."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "run: Try to run a command in a container based on RHEL but with wrong version" {
  run $TOOLBOX run --distro rhel --release 8 ls

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX run --distro rhel --release 8.2foo ls

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be in the '<major>.<minor>' format."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]

  run $TOOLBOX run --distro rhel --release -2.1 ls

  assert_failure
  assert_line --index 0 "Error: invalid argument for '--release'"
  assert_line --index 1 "The release must be a positive number."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "run: Try to run a command in a container based on non-default distro without providing a version" {
  local distro="fedora"
  local system_id="$(get_system_id)"

  if [ "$system_id" = "fedora" ]; then
    distro="rhel"
  fi

  run $TOOLBOX run -d "$distro" ls

  assert_failure
  assert_line --index 0 "Error: option '--release' is needed"
  assert_line --index 1 "Distribution $distro doesn't match the host."
  assert_line --index 2 "Run 'toolbox --help' for usage."
  assert [ ${#lines[@]} -eq 3 ]
}

@test "run: Run command exiting with non-zero code in the default container" {
  create_default_container

  run $TOOLBOX run /bin/sh -c 'exit 2'
  assert_failure
  assert [ $status -eq 2 ]
  assert_output ""
}

@test "run: Try to run non-existent command in the default container" {
  local cmd="non-existent-command"

  create_default_container

  run -127 $TOOLBOX run $cmd

  assert_failure
  assert [ $status -eq 127 ]
  assert_line --index 0 "bash: line 1: exec: $cmd: not found"
  assert_line --index 1 "Error: command $cmd not found in container $(get_latest_container_name)"
  assert [ ${#lines[@]} -eq 2 ]
}

@test "run: Try to run /etc as a command in the deault container" {
  create_default_container

  run $TOOLBOX run /etc

  assert_failure
  assert [ $status -eq 126 ]
  assert_line --index 0 "bash: line 1: /etc: Is a directory"
  assert_line --index 1 "bash: line 1: exec: /etc: cannot execute: Is a directory"
  assert_line --index 2 "Error: failed to invoke command /etc in container $(get_latest_container_name)"
  assert [ ${#lines[@]} -eq 3 ]
}

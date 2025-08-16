#!/usr/bin/env bats

load 'libs/bats-assert/load'
load 'libs/helpers'


setup() {
  cleanup_containers
}

teardown() {
  cleanup_containers
}


@test "container: Check required packages are present in the container" {
  create_default_container

  readonly REQUIRED_PACKAGES_FILE=$PROJECT_DIR/test/system/data/required-packages-$(get_system_id)-$(get_system_version)

  assert [ -f $REQUIRED_PACKAGES_FILE ]

  REQUIRED_PACKAGES_NUM=$(cat $REQUIRED_PACKAGES_FILE | wc -l)

  # tr -d '\r' is required to remove CRLF from toolbox output lines 
  OUTPUT_LINES_NUM=$($TOOLBOX run rpm -qa --qf "%{NAME}\n" | tr -d '\r' | sort | join - $REQUIRED_PACKAGES_FILE | wc -l)

  assert [ "$OUTPUT_LINES_NUM" = "$REQUIRED_PACKAGES_NUM" ]
}

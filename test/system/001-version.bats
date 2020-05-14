#!/usr/bin/env bats

load helpers

@test "Output version number using full flag" {
  run_toolbox --version
}

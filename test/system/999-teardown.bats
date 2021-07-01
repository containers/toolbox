#!/usr/bin/env bats

load 'libs/helpers'

@test "test suite: Teardown" {
  _clean_cached_images

  # Remove containers & with registries
  $PODMAN --root "$PODMAN_REG_ROOT" rm -af
  $PODMAN --root "$PODMAN_REG_ROOT" rmi -af

  # Remove test cache dir
  rm -rf "$CACHE_DIR"

  # Remove certificates
  rm -rf ~/.config/containers/certs.d/localhost:50000
  rm -rf ~/.config/containers/certs.d/localhost:50001
}

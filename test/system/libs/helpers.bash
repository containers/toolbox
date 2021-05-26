#!/usr/bin/env bash

# Podman and Toolbox commands to run
readonly PODMAN=${PODMAN:-podman}
readonly TOOLBOX=${TOOLBOX:-toolbox}
readonly SKOPEO=$(command -v skopeo)
readonly PROJECT_DIR=${PWD}

# Helpful globals
current_os_version=$(awk -F= '/VERSION_ID/ {print $2}' /etc/os-release)
readonly DEFAULT_FEDORA_VERSION=${DEFAULT_FEDORA_VERSION:-${current_os_version}}
readonly REGISTRY_URL=${REGISTRY_URL:-"registry.fedoraproject.org"}
readonly BUSYBOX_IMAGE="docker.io/library/busybox"


function cleanup_all() {
  $PODMAN system reset --force >/dev/null
}


function cleanup_containers() {
  $PODMAN rm --all --force >/dev/null
}


function get_busybox_image() {
  $PODMAN pull "$BUSYBOX_IMAGE" >/dev/null \
    || fail "Podman couldn't pull the image."
}


function pull_image() {
  local version
  local image
  version="$1"
  image="${REGISTRY_URL}/fedora-toolbox:${version}"

  $SKOPEO copy "dir:${PROJECT_DIR}/fedora-toolbox-${version}" "containers-storage:${image}"
  $PODMAN images
}


function pull_default_image() {
  pull_image "${DEFAULT_FEDORA_VERSION}"
}


function create_container() {
  local container_name
  local version
  local image
  container_name="$1"
  version="$DEFAULT_FEDORA_VERSION"
  image="${REGISTRY_URL}/fedora-toolbox:${version}"

  pull_image "$version"

  $TOOLBOX --assumeyes create --container "$container_name" \
    --image "$image" >/dev/null \
    || fail "Toolbox couldn't create the container '$container_name'"
}


function create_default_container() {
  create_container "fedora-toolbox-${DEFAULT_FEDORA_VERSION}"
}


function start_container() {
  local container_name
  container_name="$1"

  $PODMAN start "$container_name" >/dev/null \
    || fail "Podman couldn't start the container '$container_name'"
}


function stop_container() {
  local container_name
  container_name="$1"

  # Make sure the container is running before trying to stop it
  $PODMAN start "$container_name" >/dev/null \
    || fail "Podman couldn't start the container '$container_name'"
  $PODMAN stop "$container_name" >/dev/null \
    || fail "Podman couldn't stop the container '$container_name'"
}


function list_images() {
  $PODMAN images --all --quiet | wc -l
}


function list_containers() {
  $PODMAN ps --all --quiet | wc -l
}

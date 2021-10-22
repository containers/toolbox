#!/usr/bin/env bash

load 'libs/bats-support/load'

# Helpful globals
readonly TEMP_BASE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/toolbox"
readonly TEMP_STORAGE_DIR="${TEMP_BASE_DIR}/system-test-storage"

readonly IMAGE_CACHE_DIR="${BATS_RUN_TMPDIR}/image-cache"
readonly ROOTLESS_PODMAN_STORE_DIR="${TEMP_STORAGE_DIR}/storage"
readonly PODMAN_STORE_CONFIG_FILE="${TEMP_STORAGE_DIR}/store.conf"

# Podman and Toolbox commands to run
readonly PODMAN=${PODMAN:-podman}
readonly TOOLBOX=${TOOLBOX:-toolbox}
readonly SKOPEO=$(command -v skopeo)

# Images
declare -Ag IMAGES=([busybox]="quay.io/toolbox_tests/busybox" \
                   [fedora]="registry.fedoraproject.org/fedora-toolbox" \
                   [rhel]="registry.access.redhat.com/ubi8")


function cleanup_all() {
  $PODMAN system reset --force >/dev/null
}


function cleanup_containers() {
  $PODMAN rm --all --force >/dev/null
}


function _setup_environment() {
  _setup_containers_store
  check_xdg_runtime_dir
}

function _setup_containers_store() {
  mkdir -p ${TEMP_STORAGE_DIR}
  # Setup a storage config file for PODMAN
  echo -e "[storage]\n  driver = \"overlay\"\n  rootless_storage_path = \"${ROOTLESS_PODMAN_STORE_DIR}\"\n" > ${PODMAN_STORE_CONFIG_FILE}
  export CONTAINERS_STORAGE_CONF=${PODMAN_STORE_CONFIG_FILE}
}


function _clean_temporary_storage() {
  rm -rf ${TEMP_STORAGE_DIR}
}


# Pulls an image using Podman and saves it to a image dir using Skopeo
#
# Parameters
# ==========
# - distro - os-release field ID (e.g., fedora, rhel)
# - version - os-release field VERSION_ID (e.g., 33, 34, 8.4)
#
# Only use during test suite setup for caching all images to be used throught
# tests.
function _pull_and_cache_distro_image() {
  local num_of_retries=5
  local timeout=10
  local pulled=false
  local distro
  local version
  local image
  local image_archive

  distro="$1"
  version="$2"

  if [ ! -v IMAGES[$distro] ]; then
    fail "Requested distro (${distro}) does not have a matching image"
  fi

  image="${IMAGES[$distro]}"
  image_archive="${distro}-toolbox"

  if [[ $# -eq 2 ]]; then
    image="${image}:${version}"
    image_archive="${image_archive}-${version}"
  fi

  for ((i = ${num_of_retries}; i > 0; i--)); do
    run $PODMAN pull ${image}

    if [ "$status" -eq 0 ]; then
      pulled=true
      break
    fi

    sleep $timeout
  done

  if ! $pulled; then
    echo "Failed to pull image ${image}"
    assert_success
  fi

  if [ ! -d ${IMAGE_CACHE_DIR} ]; then
    mkdir -p ${IMAGE_CACHE_DIR}
  fi

  run $SKOPEO copy --dest-compress containers-storage:${image} dir:${IMAGE_CACHE_DIR}/${image_archive}

  if [ "$status" -ne 0 ]; then
    echo "Failed to cache image ${image} to ${IMAGE_CACHE_DIR}/${image_archive}"
    assert_success
  fi

  cleanup_all
}


# Removes the folder with cached images
function _clean_cached_images() {
  rm -rf ${IMAGE_CACHE_DIR}
}


# Copies an image from local storage to Podman's image store
# 
# Call before creating any container. Network failures are not nice.
#
# An image has to be cached first. See _pull_and_cache_distro_image()
#
# Parameters:
# ===========
# - distro - os-release field ID (e.g., fedora, rhel)
# - version - os-release field VERSION_ID (e.g., 33, 34, 8.4)
function pull_distro_image() {
  local distro
  local version
  local image
  local image_archive

  distro="$1"
  version="$2"

  if [ ! -v IMAGES[$distro] ]; then
    fail "Requested distro (${distro}) does not have a matching image"
  fi

  image="${IMAGES[$distro]}"
  image_archive="${distro}-toolbox"

  if [[ -n $version ]]; then
    image="${image}:${version}"
    image_archive="${image_archive}-${version}"
  fi

  # No need to copy if the image is already available in Podman
  run $PODMAN image exists ${image}
  if [[ "$status" -eq 0 ]]; then
    return
  fi

  # https://github.com/containers/skopeo/issues/547 for the options for containers-storage
  run $SKOPEO copy "dir:${IMAGE_CACHE_DIR}/${image_archive}" "containers-storage:[overlay@$ROOTLESS_PODMAN_STORE_DIR+$ROOTLESS_PODMAN_STORE_DIR]${image}"
  if [ "$status" -ne 0 ]; then
    echo "Failed to load image ${image} from cache ${IMAGE_CACHE_DIR}/${image_archive}"
    assert_success
  fi

  $PODMAN images
}


# Copies the system's default image to Podman's image store
#
# See pull_default_image() for more info.
function pull_default_image() {
  pull_distro_image $(get_system_id) $(get_system_version)
}


# Creates a container with specific name, distro and version
#
# Pulling of an image is taken care of by the function
#
# Parameters:
# ===========
# - distro - os-release field ID (e.g., fedora, rhel)
# - version - os-release field VERSION_ID (e.g., 33, 34, 8.4)
# - container_name - name of the container
function create_distro_container() {
  local distro
  local version
  local container_name

  distro="$1"
  version="$2"
  container_name="$3"

  pull_distro_image ${distro} ${version}

  $TOOLBOX --assumeyes create --container "${container_name}" --distro "${distro}" --release "${version}" >/dev/null \
    || fail "Toolbox couldn't create container '$container_name'"
}


# Creates a container with specific name matching the system
#
# Parameters:
# ===========
# - container_name - name of the container
function create_container() {
  local container_name

  container_name="$1"

  create_distro_container $(get_system_id) $(get_system_version) $container_name
}


# Creates a default container
function create_default_container() {
  pull_default_image

  $TOOLBOX --assumeyes create >/dev/null \
    || fail "Toolbox couldn't create default container"
}


function start_container() {
  local container_name
  container_name="$1"

  $PODMAN start "$container_name" >/dev/null \
    || fail "Podman couldn't start the container '$container_name'"
}


# Checks if a toolbox container started
#
# Parameters:
# ===========
# - container_name - name of the container
function container_started() {
  local container_name
  container_name="$1"

  run $PODMAN start $container_name

  container_initialized=0

  for TRIES in 1 2 3 4 5
  do
    run $PODMAN logs $container_name
    container_output=$output
    # Look for last line of the container startup log
    run grep 'Listening to file system and ticker events' <<< $container_output
    if [[ "$status" -eq 0 ]]; then
      container_initialized=1
      break
    fi
    sleep 1
  done

  echo $container_output >2

  echo $container_initialized
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


# Returns the path to os-release
function find_os_release() {
  if [[ -f "/etc/os-release" ]]; then
    echo "/etc/os-release"
  elif [[ -f "/usr/lib/os-release" ]]; then
    echo "/usr/lib/os-release"
  else
    echo ""
  fi
}


# Returns the content of field ID in os-release
function get_system_id() {
    local os_release

    os_release="$(find_os_release)"

    if [[ -z "$os_release" ]]; then
        echo ""
        return
    fi

    echo $(awk -F= '/ID/ {print $2}' $os_release | head -n 1)
}


# Returns the content of field VERSION_ID in os-release
function get_system_version() {
    local os_release

    os_release="$(find_os_release)"

    if [[ -z "$os_release" ]]; then
        echo ""
        return
    fi

    echo $(awk -F= '/VERSION_ID/ {print $2}' $os_release | head -n 1)
}


# Setup the XDG_RUNTIME_DIR variable if not set
function check_xdg_runtime_dir() {
    if [[ -z "${XDG_RUNTIME_DIR}" ]]; then
        export XDG_RUNTIME_DIR="/run/user/${UID}"
    fi
}

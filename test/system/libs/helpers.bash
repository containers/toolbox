#!/usr/bin/env bash

load 'libs/bats-support/load'
load 'libs/bats-assert/load'

# Helpful globals
readonly TEMP_BASE_DIR="${XDG_CACHE_HOME:-$HOME/.cache}/toolbox"
readonly TEMP_STORAGE_DIR="${TEMP_BASE_DIR}/system-test-storage"

readonly IMAGE_CACHE_DIR="${BATS_RUN_TMPDIR}/image-cache"
readonly ROOTLESS_PODMAN_STORE_DIR="${TEMP_STORAGE_DIR}/storage"
readonly ROOTLESS_PODMAN_RUNROOT_DIR="${TEMP_STORAGE_DIR}/runroot"
readonly PODMAN_STORE_CONFIG_FILE="${TEMP_STORAGE_DIR}/storage.conf"
readonly DOCKER_REG_ROOT="${TEMP_STORAGE_DIR}/docker-registry-root"
readonly DOCKER_REG_CERTS_DIR="${BATS_RUN_TMPDIR}/certs"
readonly DOCKER_REG_AUTH_DIR="${BATS_RUN_TMPDIR}/auth"
readonly DOCKER_REG_URI="localhost:50000"
readonly DOCKER_REG_NAME="docker-registry"

# Podman and Toolbox commands to run
readonly PODMAN=${PODMAN:-$(command -v podman)}
readonly TOOLBOX=${TOOLBOX:-$(command -v toolbox)}
readonly SKOPEO=${SKOPEO:-$(command -v skopeo)}

# Images
declare -Ag IMAGES=([busybox]="quay.io/toolbox_tests/busybox" \
                   [docker-reg]="quay.io/toolbox_tests/registry" \
                   [fedora]="registry.fedoraproject.org/fedora-toolbox" \
                   [rhel]="registry.access.redhat.com/ubi8/toolbox")


function cleanup_all() {
  $PODMAN rm --all --force >/dev/null
  $PODMAN rmi --all --force >/dev/null
}


function cleanup_containers() {
  $PODMAN rm --all --force >/dev/null
}


function _setup_environment() {
  _setup_containers_storage
  check_xdg_runtime_dir
}

function _setup_containers_storage() {
  mkdir -p ${TEMP_STORAGE_DIR}
  # Setup a storage config file for PODMAN
  echo -e "[storage]\n  driver = \"overlay\"\n  rootless_storage_path = \"${ROOTLESS_PODMAN_STORE_DIR}\"\n  runroot = \"${ROOTLESS_PODMAN_RUNROOT_DIR}\"\n" > ${PODMAN_STORE_CONFIG_FILE}
  export CONTAINERS_STORAGE_CONF=${PODMAN_STORE_CONFIG_FILE}
}


function _clean_temporary_storage() {
  $PODMAN system reset -f

  rm -rf ${ROOTLESS_PODMAN_STORE_DIR}
  rm -rf ${ROOTLESS_PODMAN_RUNROOT_DIR}
  rm -rf ${PODMAN_STORE_CONFIG_FILE}
  rm -rf ${TEMP_STORAGE_DIR}
}


# Caches an image associated with a distribution at a specific release to
# an image dir using Skopeo
#
# Parameters
# ==========
# - distro - os-release field ID (e.g., fedora, rhel)
# - version - os-release field VERSION_ID (e.g., 33, 34, 8.4) (optional)
#
# Only use during test suite setup for caching all images to be used throught
# tests.
function _pull_and_cache_distro_image() {
  local distro
  local image
  local image_archive
  local version

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

  _pull_and_cache_image ${image} ${image_archive}
}


# Caches an image to an image dir using Skopeo
#
# Parameters
# ==========
# - image - OCI image available via docker:// transport
# - image_archive - name of the cache archive (optional)
#
# Only use during test suite setup for caching all images to be used throught
# tests.
function _pull_and_cache_image() {
  local num_of_retries=5
  local timeout=10
  local cached=false
  local image
  local image_archive

  image="$1"
  image_archive="${image//:/-}"

  if [[ $# -eq 2 ]]; then
    image_archive="$2"
  fi

  image_archive="${IMAGE_CACHE_DIR}/${image_archive}"

  if [[ -d ${image_archive} ]] ; then
    return 0
  fi

  if [ ! -d $(dirname ${image_archive}) ]; then
    run mkdir -p $(dirname ${image_archive})
    assert_success
  fi

  for ((i = ${num_of_retries}; i > 0; i--)); do
    run $SKOPEO copy --dest-compress docker://${image} dir:${image_archive}

    if [ "$status" -eq 0 ]; then
      cached=true
      break
    fi

    sleep $timeout
  done

  if ! $cached; then
    echo "Failed to cache image ${image} to ${image_archive}"
    assert_success
  fi

  cleanup_all
}


# Removes the folder with cached images
function _clean_cached_images() {
  rm -rf ${IMAGE_CACHE_DIR}
}


# Prepares a localy hosted image registry
#
# The registry is set up with Podman set to an alternative root. It won't
# affect other containers or images in the default root.
#
# Instructions taken from https://docs.docker.com/registry/deploying/
function _setup_docker_registry() {
  # Create certificates for HTTPS
  # This is needed so that Podman does not have to be configured to work with
  # HTTP-only registries
  run mkdir -p "${DOCKER_REG_CERTS_DIR}"
  assert_success
  run openssl req \
    -newkey rsa:4096 \
    -nodes -sha256 \
    -keyout "${DOCKER_REG_CERTS_DIR}"/domain.key \
    -addext "subjectAltName= DNS:localhost" \
    -x509 \
    -days 365 \
    -subj '/' \
    -out "${DOCKER_REG_CERTS_DIR}"/domain.crt
  assert_success

  # Add certificate to Podman's trusted certificates (rootless)
  run mkdir -p "$HOME"/.config/containers/certs.d/"${DOCKER_REG_URI}"
  assert_success
  run cp "${DOCKER_REG_CERTS_DIR}"/domain.crt "$HOME"/.config/containers/certs.d/"${DOCKER_REG_URI}"/domain.crt
  assert_success

  # Create a registry user
  # username: user; password: user
  run mkdir -p "${DOCKER_REG_AUTH_DIR}"
  assert_success
  run htpasswd -Bbc "${DOCKER_REG_AUTH_DIR}"/htpasswd user user
  assert_success

  # Create separate Podman root
  run mkdir -p "${DOCKER_REG_ROOT}"
  assert_success

  # Pull Docker registry image
  run $PODMAN --root "${DOCKER_REG_ROOT}" pull "${IMAGES[docker-reg]}"
  assert_success

  # Create a Docker registry
  run $PODMAN --root "${DOCKER_REG_ROOT}" run -d \
    --rm \
    --name "${DOCKER_REG_NAME}" \
    --privileged \
    -v "${DOCKER_REG_AUTH_DIR}":/auth \
    -e REGISTRY_AUTH=htpasswd \
    -e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
    -e REGISTRY_AUTH_HTPASSWD_PATH="/auth/htpasswd" \
    -v "${DOCKER_REG_CERTS_DIR}":/certs \
    -e REGISTRY_HTTP_ADDR=0.0.0.0:443 \
    -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
    -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
    -p 50000:443 \
    "${IMAGES[docker-reg]}"
  assert_success

  run $PODMAN login \
    --authfile ${TEMP_BASE_DIR}/authfile.json \
    --username user \
    --password user \
    "${DOCKER_REG_URI}"
  assert_success

  # Add fedora-toolbox:32 image to the registry
  run $SKOPEO copy --dest-authfile ${TEMP_BASE_DIR}/authfile.json \
    dir:"${IMAGE_CACHE_DIR}"/fedora-toolbox-32 \
    docker://"${DOCKER_REG_URI}"/fedora-toolbox:32
  assert_success

  run rm ${TEMP_BASE_DIR}/authfile.json
  assert_success
}


# Stop, removes and cleans after a locally hosted Docker registry
function _clean_docker_registry() {
  # Stop Docker registry container
  $PODMAN --root "${DOCKER_REG_ROOT}" stop --time 0 "${DOCKER_REG_NAME}"
  # Clean up Podman's registry root state
  $PODMAN --root "${DOCKER_REG_ROOT}" rm --all --force
  $PODMAN --root "${DOCKER_REG_ROOT}" rmi --all --force
  # Remove Docker registry dir
  rm -rf "${DOCKER_REG_ROOT}"
  # Remove dir with created registry certificates
  rm -rf "$HOME"/.config/containers/certs.d/"${DOCKER_REG_URI}"
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

  pull_image ${image} ${image_archive}
}


#
#
# Parameters:
# ===========
# - image
# - image_archive
function pull_image() {
  local image
  local image_archive

  # No need to copy if the image is already available in Podman
  run $PODMAN image exists ${image}
  if [[ "$status" -eq 0 ]]; then
    return
  fi

  image="$1"
  image_archive="${image//:/-}"

  if [[ $# -eq 2 ]]; then
    image_archive="$2"
  fi

  image_archive="${IMAGE_CACHE_DIR}/${image_archive}"

  # https://github.com/containers/skopeo/issues/547 for the options for containers-storage
  run $SKOPEO copy "dir:${image_archive}" "containers-storage:[overlay@$ROOTLESS_PODMAN_STORE_DIR+$ROOTLESS_PODMAN_STORE_DIR]${image}"
  if [ "$status" -ne 0 ]; then
    echo "Failed to load image ${image} from cache ${image_archive}"
    assert_success
  fi

  $PODMAN images
}

# Copies the system's default image to Podman's image store
#
# See pull_default_image() for more info.
function pull_default_image() {
	pull_image $(toolbx_default_image)
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

  pull_default_image

  $TOOLBOX --assumeyes create --container "${container_name}" >/dev/null \
    || fail "Toolbox couldn't create container '${container_name}'"
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
#
# Returns:
# ========
# - 0 - container has not started
# - 1 - container has started
function container_started() {
  local container_name
  container_name="$1"

  run $PODMAN start $container_name

  # Used as a return value
  container_initialized=1

  for TRIES in 1 2 3 4 5
  do
    run $PODMAN logs $container_name
    container_output=$output
    # Look for last line of the container startup log
    run grep 'Listening to file system and ticker events' <<< $container_output
    if [[ "$status" -eq 0 ]]; then
      container_initialized=0
      break
    fi
    sleep 1
  done

  return $container_initialized
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


# Returns the name of the latest created container
function get_latest_container_name() {
  $PODMAN ps -l --format "{{ .Names }}"
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


# Prints a value in Toolbx config
#
# If key does not exist, prints nothing
#
# Parameters:
# ===========
# - config-key - config key
function toolbx_config_key() {
  local config_key="$1"

  echo $(TOOLBOX __test --type config-key "$config_key")
}


# Prints the default Toolbx container name
function toolbx_default_container_name() {
  echo $($TOOLBOX __test --type default-container-name)
}


# Prints the default Toolbx OCI image
function toolbx_default_image() {
  echo $($TOOLBOX __test --type default-image)
}


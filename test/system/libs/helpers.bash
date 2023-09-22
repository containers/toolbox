# shellcheck shell=bash

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
readonly PODMAN="${PODMAN:-$(command -v podman)}"
readonly TOOLBOX="${TOOLBOX:-$(command -v toolbox)}"
readonly SKOPEO="${SKOPEO:-$(command -v skopeo)}"

# Images
declare -Ag IMAGES=([arch]="quay.io/toolbx/arch-toolbox" \
                   [busybox]="quay.io/toolbox_tests/busybox" \
                   [docker-reg]="quay.io/toolbox_tests/registry" \
                   [fedora]="registry.fedoraproject.org/fedora-toolbox" \
                   [rhel]="registry.access.redhat.com/ubi8/toolbox" \
                   [ubuntu]="quay.io/toolbx/ubuntu-toolbox")


function cleanup_all() {
  "$PODMAN" rm --all --force >/dev/null
  "$PODMAN" rmi --all --force >/dev/null
}


function cleanup_containers() {
  "$PODMAN" rm --all --force >/dev/null
}


function _setup_environment() {
  _setup_containers_storage
  check_xdg_runtime_dir
}

function _setup_containers_storage() {
  mkdir -p "${TEMP_STORAGE_DIR}"
  # Set up a storage config file for PODMAN
  echo -e "[storage]\n  driver = \"overlay\"\n  rootless_storage_path = \"${ROOTLESS_PODMAN_STORE_DIR}\"\n  runroot = \"${ROOTLESS_PODMAN_RUNROOT_DIR}\"\n" > "${PODMAN_STORE_CONFIG_FILE}"
  export CONTAINERS_STORAGE_CONF="${PODMAN_STORE_CONFIG_FILE}"
}


function _clean_temporary_storage() {
  "$PODMAN" system reset --force

  rm --force --recursive "${ROOTLESS_PODMAN_STORE_DIR}"
  rm --force --recursive "${ROOTLESS_PODMAN_RUNROOT_DIR}"
  rm --force --recursive "${PODMAN_STORE_CONFIG_FILE}"
  rm --force --recursive "${TEMP_STORAGE_DIR}"
}


# Pulls an image using Podman and saves it to a image dir using Skopeo
#
# Parameters
# ==========
# - distro - os-release field ID (e.g., fedora, rhel)
# - version - os-release field VERSION_ID (e.g., 33, 34, 8.4)
#
# Only use during test suite setup for caching all images to be used throughout
# tests.
function _pull_and_cache_distro_image() {
  local num_of_retries=5
  local timeout=10
  local cached=false
  local distro
  local version
  local image
  local image_archive

  distro="$1"
  version="$2"

  if [ -z "${IMAGES[$distro]+x}" ]; then
    fail "Requested distro (${distro}) does not have a matching image"
    return 1
  fi

  image="${IMAGES[$distro]}"
  image_archive="${distro}-toolbox"

  if [[ $# -eq 2 ]]; then
    image="${image}:${version}"
    image_archive="${image_archive}-${version}"
  fi

  if [[ -d "${IMAGE_CACHE_DIR}/${image_archive}" ]] ; then
    return 0
  fi

  if [ ! -d "${IMAGE_CACHE_DIR}" ]; then
    run mkdir -p "${IMAGE_CACHE_DIR}"
    assert_success
  fi

  local -i j

  for ((j = 0; j < num_of_retries; j++)); do
    run "$SKOPEO" copy --dest-compress "docker://${image}" "dir:${IMAGE_CACHE_DIR}/${image_archive}"

    if [ "$status" -eq 0 ]; then
      cached=true
      break
    fi

    sleep "$timeout"
  done

  if ! $cached; then
    echo "Failed to cache image ${image} to ${IMAGE_CACHE_DIR}/${image_archive}"
    assert_success
  fi

  cleanup_all
}


# Removes the folder with cached images
function _clean_cached_images() {
  rm --force --recursive "${IMAGE_CACHE_DIR}"
}


# Prepares a locally hosted image registry
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
  run "$PODMAN" --root "${DOCKER_REG_ROOT}" pull "${IMAGES[docker-reg]}"
  assert_success

  # Create a Docker registry
  run "$PODMAN" --root "${DOCKER_REG_ROOT}" run -d \
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

  run "$PODMAN" login \
    --authfile "${TEMP_BASE_DIR}/authfile.json" \
    --username user \
    --password user \
    "${DOCKER_REG_URI}"
  assert_success

  # Add fedora-toolbox:34 image to the registry
  run "$SKOPEO" copy --dest-authfile "${TEMP_BASE_DIR}/authfile.json" \
    dir:"${IMAGE_CACHE_DIR}"/fedora-toolbox-34 \
    docker://"${DOCKER_REG_URI}"/fedora-toolbox:34
  assert_success

  run rm "${TEMP_BASE_DIR}/authfile.json"
  assert_success
}


# Stop, removes and cleans after a locally hosted Docker registry
function _clean_docker_registry() {
  # Stop Docker registry container
  if "$PODMAN" --root "$DOCKER_REG_ROOT" container exists "$DOCKER_REG_NAME"; then
    "$PODMAN" --root "${DOCKER_REG_ROOT}" stop --time 0 "${DOCKER_REG_NAME}"
  fi

  # Clean up Podman's registry root state
  "$PODMAN" --root "${DOCKER_REG_ROOT}" rm --all --force
  "$PODMAN" --root "${DOCKER_REG_ROOT}" rmi --all --force
  # Remove Docker registry dir
  rm --force --recursive "${DOCKER_REG_ROOT}"
  # Remove dir with created registry certificates
  rm --force --recursive "$HOME"/.config/containers/certs.d/"${DOCKER_REG_URI}"
}


function build_image_without_name() {
  echo -e "FROM scratch\n\nLABEL com.github.containers.toolbox=\"true\"" > "$BATS_TEST_TMPDIR"/Containerfile

  run "$PODMAN" build "$BATS_TEST_TMPDIR"

  assert_success
  assert_line --index 0 --partial "FROM scratch"
  assert_line --index 1 --partial "LABEL com.github.containers.toolbox=\"true\""
  assert_line --index 2 --partial "COMMIT"
  assert_line --index 3 --regexp "^--> [a-f0-9]{6,64}$"

  # shellcheck disable=SC2154
  last=$((${#lines[@]}-1))

  assert_line --index "$last" --regexp "^[a-f0-9]{64}$"

  rm -f "$BATS_TEST_TMPDIR"/Containerfile

  echo "${lines[$last]}"
}


function check_bats_version() {
    local required_version
    required_version="$1"

    if ! old_version=$(printf "%s\n%s\n" "$BATS_VERSION" "$required_version" | sort --version-sort | head --lines 1); then
        return 1
    fi

    if [ "$required_version" = "$old_version" ]; then
        return 0
    fi

    return 1
}


function get_busybox_image() {
  local image
  image="${IMAGES[busybox]}"
  echo "$image"
  return 0
}


function get_default_image() {
  local distro
  local image
  local release

  distro="$(get_system_id)"
  release="$(get_system_version)"
  image="${IMAGES[$distro]}:$release"

  echo "$image"
  return 0
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

  if [ -z "${IMAGES[$distro]+x}" ]; then
    fail "Requested distro (${distro}) does not have a matching image"
    return 1
  fi

  image="${IMAGES[$distro]}"
  image_archive="${distro}-toolbox"

  if [[ -n $version ]]; then
    image="${image}:${version}"
    image_archive="${image_archive}-${version}"
  fi

  # No need to copy if the image is already available in Podman
  run "$PODMAN" image exists "${image}"
  if [[ "$status" -eq 0 ]]; then
    return
  fi

  # https://github.com/containers/skopeo/issues/547 for the options for containers-storage
  run "$SKOPEO" copy "dir:${IMAGE_CACHE_DIR}/${image_archive}" "containers-storage:[overlay@$ROOTLESS_PODMAN_STORE_DIR+$ROOTLESS_PODMAN_STORE_DIR]${image}"
  if [ "$status" -ne 0 ]; then
    echo "Failed to load image ${image} from cache ${IMAGE_CACHE_DIR}/${image_archive}"
    assert_success
  fi
}


# Copies the system's default image to Podman's image store
#
# See pull_default_image() for more info.
function pull_default_image() {
  pull_distro_image "$(get_system_id)" "$(get_system_version)"
}


function pull_default_image_and_copy() {
  pull_default_image

  local distro
  local version
  local image

  distro="$(get_system_id)"
  version="$(get_system_version)"
  image="${IMAGES[$distro]}:$version"

  # https://github.com/containers/skopeo/issues/547 for the options for containers-storage
  run "$SKOPEO" copy \
      "containers-storage:[overlay@$ROOTLESS_PODMAN_STORE_DIR+$ROOTLESS_PODMAN_STORE_DIR]$image" \
      "containers-storage:[overlay@$ROOTLESS_PODMAN_STORE_DIR+$ROOTLESS_PODMAN_STORE_DIR]$image-copy"

  if [ "$status" -ne 0 ]; then
    echo "Failed to copy image $image to $image-copy"
    assert_success
  fi
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

  pull_distro_image "${distro}" "${version}"

  "$TOOLBOX" --assumeyes create --container "${container_name}" --distro "${distro}" --release "${version}" >/dev/null \
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

  create_distro_container "$(get_system_id)" "$(get_system_version)" "$container_name"
}


# Creates a default container
function create_default_container() {
  pull_default_image

  "$TOOLBOX" --assumeyes create >/dev/null \
    || fail "Toolbox couldn't create default container"
}


function start_container() {
  local container_name
  container_name="$1"

  "$PODMAN" start "$container_name" >/dev/null \
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

  local -i ret_val=1

  start_container "$container_name"

  local -i j
  local num_of_retries=5

  for ((j = 0; j < num_of_retries; j++)); do
    run --separate-stderr "$PODMAN" logs "$container_name"

    # shellcheck disable=SC2154
    if [ "$status" -ne 0 ]; then
      fail "Failed to invoke '$PODMAN logs'"
      ret_val="$status"
      break
    fi

    # Look for last line of the container startup log
    # shellcheck disable=SC2154
    if echo "$output $stderr" | grep "Listening to file system and ticker events"; then
      ret_val=0
      break
    fi

    sleep 1
  done

  if [ "$ret_val" -ne 0 ]; then
    if [ "$j" -eq "$num_of_retries" ]; then
      fail "Failed to initialize container $container_name"
    fi

    [ "$output" != "" ] && echo "$output"
    [ "$stderr" != "" ] && echo "$stderr" >&2
  fi

  return "$ret_val"
}


function stop_container() {
  local container_name
  container_name="$1"

  # Make sure the container is running before trying to stop it
  "$PODMAN" start "$container_name" >/dev/null \
    || fail "Podman couldn't start the container '$container_name'"
  "$PODMAN" stop "$container_name" >/dev/null \
    || fail "Podman couldn't stop the container '$container_name'"
}


# Returns the name of the latest created container
function get_latest_container_name() {
  "$PODMAN" ps --latest --format "{{ .Names }}"
}


function list_images() {
  "$PODMAN" images --all --format "{{.ID}}" | wc --lines
}


function list_containers() {
  "$PODMAN" ps --all --quiet | wc --lines
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
function get_system_id() (
  local os_release

  os_release="$(find_os_release)"

  if [[ -z "$os_release" ]]; then
    echo ""
    return
  fi

  # shellcheck disable=SC1090
  . "$os_release"

  echo "$ID"
)


# Returns the content of field VERSION_ID in os-release
function get_system_version() (
  local os_release

  os_release="$(find_os_release)"

  if [[ -z "$os_release" ]]; then
    echo ""
    return
  fi

  # shellcheck disable=SC1090
  . "$os_release"

  echo "$VERSION_ID"
)


function is_fedora_rawhide() (
  local os_release
  os_release="$(find_os_release)"
  [ -z "$os_release" ] && return 1

  # shellcheck disable=SC1090
  . "$os_release"

  [ "$ID" != "fedora" ] && return 1
  [ "$REDHAT_BUGZILLA_PRODUCT_VERSION" != "rawhide" ] && return 1

  return 0
)


# Set up the XDG_RUNTIME_DIR variable if not set
function check_xdg_runtime_dir() {
  if [[ -z "${XDG_RUNTIME_DIR}" ]]; then
    export XDG_RUNTIME_DIR="/run/user/${UID}"
  fi
}

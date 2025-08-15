# shellcheck shell=bash

load 'libs/bats-support/load'
load 'libs/bats-assert/load'

readonly HOME="$BATS_SUITE_TMPDIR/home"
export HOME

readonly XDG_CACHE_HOME="$HOME/.cache"
export XDG_CACHE_HOME

readonly XDG_CONFIG_HOME="$HOME/.config"
export XDG_CONFIG_HOME

readonly XDG_DATA_HOME="$HOME/.local/share"
export XDG_DATA_HOME

readonly XDG_RUNTIME_DIR="${XDG_RUNTIME_DIR:-/run/user/$UID}"
export XDG_RUNTIME_DIR

readonly XDG_STATE_HOME="$HOME/.local/state"
export XDG_STATE_HOME

readonly CONTAINERS_STORAGE_CONF="$XDG_CONFIG_HOME/containers/storage.conf"
export CONTAINERS_STORAGE_CONF

# Helpful globals
readonly IMAGE_CACHE_DIR="$BATS_SUITE_TMPDIR/image-cache"
readonly TOOLBX_ROOTLESS_STORAGE_PATH="$XDG_DATA_HOME/containers/storage"
readonly ROOTLESS_PODMAN_RUNROOT_DIR="$BATS_SUITE_TMPDIR/runroot"
readonly DOCKER_REG_ROOT="$BATS_SUITE_TMPDIR/docker-registry-root"
readonly DOCKER_REG_CERTS_DIR="$BATS_SUITE_TMPDIR/certs"
readonly DOCKER_REG_AUTH_DIR="$BATS_SUITE_TMPDIR/auth"
readonly DOCKER_REG_URI="localhost:50000"
readonly DOCKER_REG_NAME="docker-registry"

# Podman and Toolbx commands to run
readonly TOOLBX="${TOOLBX:-$(command -v toolbox)}"
readonly TOOLBX_TEST_SYSTEM_TAGS_ALL="arch-fedora,commands-options,custom-image,runtime-environment,ubuntu"
readonly TOOLBX_TEST_SYSTEM_TAGS="${TOOLBX_TEST_SYSTEM_TAGS:-$TOOLBX_TEST_SYSTEM_TAGS_ALL}"

# Images
declare -Ag IMAGES=([arch]="quay.io/toolbx/arch-toolbox" \
                   [busybox]="quay.io/toolbox_tests/busybox" \
                   [docker-reg]="quay.io/toolbox_tests/registry" \
                   [fedora]="registry.fedoraproject.org/fedora-toolbox" \
                   [rhel]="registry.access.redhat.com/ubi8/toolbox" \
                   [ubuntu]="quay.io/toolbx/ubuntu-toolbox")


function cleanup_all() {
  podman rm --all --force >/dev/null
  podman rmi --all --force >/dev/null
}


function _setup_environment() {
  # shellcheck disable=SC2174
  mkdir --mode 700 --parents "$HOME"
  _setup_containers_storage
}

function _setup_containers_storage() {
  # Set up a storage config file for PODMAN
  configuration_directory="$(dirname "$CONTAINERS_STORAGE_CONF")"
  mkdir --parents "$configuration_directory"
  echo -e "[storage]\n  driver = \"overlay\"\n  runroot = \"$ROOTLESS_PODMAN_RUNROOT_DIR\"\n" > "$CONTAINERS_STORAGE_CONF"
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

  local error_message
  local -i j
  local -i ret_val

  for ((j = 0; j < num_of_retries; j++)); do
    error_message="$( (skopeo copy --dest-compress \
                          "docker://${image}" \
                          "dir:${IMAGE_CACHE_DIR}/${image_archive}" >/dev/null) 2>&1)"
    ret_val="$?"

    if [ "$ret_val" -eq 0 ]; then
      cached=true
      break
    fi

    sleep "$timeout"
  done

  if ! $cached; then
    echo "Failed to cache image ${image} to ${IMAGE_CACHE_DIR}/${image_archive}" >&2
    [ "$error_message" != "" ] && echo "$error_message" >&2
    return "$ret_val"
  fi

  cleanup_all
  ret_val="$?"

  return "$ret_val"
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
  run podman --root "${DOCKER_REG_ROOT}" pull "${IMAGES[docker-reg]}"
  assert_success

  # Create a Docker registry
  run podman --root "${DOCKER_REG_ROOT}" run \
    --detach \
    --env REGISTRY_AUTH=htpasswd \
    --env REGISTRY_AUTH_HTPASSWD_PATH="/auth/htpasswd" \
    --env REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
    --env REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
    --env REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
    --name "${DOCKER_REG_NAME}" \
    --privileged \
    --publish 50000:5000 \
    --rm \
    --volume "${DOCKER_REG_AUTH_DIR}":/auth \
    --volume "${DOCKER_REG_CERTS_DIR}":/certs \
    "${IMAGES[docker-reg]}"
  assert_success

  run podman login \
    --authfile "${BATS_SUITE_TMPDIR}/authfile.json" \
    --username user \
    --password user \
    "${DOCKER_REG_URI}"
  assert_success

  # Add fedora-toolbox:34 image to the registry
  run skopeo copy --dest-authfile "${BATS_SUITE_TMPDIR}/authfile.json" \
    dir:"${IMAGE_CACHE_DIR}"/fedora-toolbox-34 \
    docker://"${DOCKER_REG_URI}"/fedora-toolbox:34
  assert_success

  run rm "${BATS_SUITE_TMPDIR}/authfile.json"
  assert_success
}


# Stop, removes and cleans after a locally hosted Docker registry
function _clean_docker_registry() {
  # Stop Docker registry container
  if podman --root "$DOCKER_REG_ROOT" container exists "$DOCKER_REG_NAME"; then
    podman --root "${DOCKER_REG_ROOT}" stop --time 0 "${DOCKER_REG_NAME}"
  fi

  # Clean up Podman's registry root state
  podman --root "${DOCKER_REG_ROOT}" rm --all --force
  podman --root "${DOCKER_REG_ROOT}" rmi --all --force
  # Remove Docker registry dir
  rm --force --recursive "${DOCKER_REG_ROOT}"
  # Remove dir with created registry certificates
  rm --force --recursive "$HOME"/.config/containers/certs.d/"${DOCKER_REG_URI}"
}


function build_image_without_name() {
  echo -e "FROM scratch\n\nLABEL com.github.containers.toolbox=\"true\"" > "$BATS_TEST_TMPDIR"/Containerfile

  run podman build --quiet "$BATS_TEST_TMPDIR"

  assert_success
  assert_line --index 0 --regexp "^[a-f0-9]{64}$"

  # shellcheck disable=SC2154
  assert [ ${#lines[@]} -eq 1 ]

  rm -f "$BATS_TEST_TMPDIR"/Containerfile

  echo "${lines[0]}"
}


function build_non_toolbx_image() {
  local image_name="localhost/non-toolbx:test-$$"

  echo -e "FROM scratch\n\nLABEL test=\"non-toolbx\"" > "$BATS_TEST_TMPDIR"/Containerfile

  run podman build --quiet --tag "$image_name" "$BATS_TEST_TMPDIR"

  assert_success
  assert_line --index 0 --regexp "^[a-f0-9]{64}$"

  # shellcheck disable=SC2154
  assert [ ${#lines[@]} -eq 1 ]

  rm -f "$BATS_TEST_TMPDIR"/Containerfile

  echo "$image_name"
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
  if podman image exists "${image}"; then
    return 0
  fi

  # https://github.com/containers/skopeo/issues/547 for the options for containers-storage
  run skopeo copy "dir:${IMAGE_CACHE_DIR}/${image_archive}" "containers-storage:[overlay@$TOOLBX_ROOTLESS_STORAGE_PATH+$TOOLBX_ROOTLESS_STORAGE_PATH]${image}"

  # shellcheck disable=SC2154
  if [ "$status" -ne 0 ]; then
    echo "Failed to load image ${image} from cache ${IMAGE_CACHE_DIR}/${image_archive}"
    assert_success
  fi

  return 0
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
  run skopeo copy \
      "containers-storage:[overlay@$TOOLBX_ROOTLESS_STORAGE_PATH+$TOOLBX_ROOTLESS_STORAGE_PATH]$image" \
      "containers-storage:[overlay@$TOOLBX_ROOTLESS_STORAGE_PATH+$TOOLBX_ROOTLESS_STORAGE_PATH]$image-copy"

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

  "$TOOLBX" --assumeyes create --container "${container_name}" --distro "${distro}" --release "${version}" >/dev/null \
    || fail "Toolbx couldn't create container '$container_name'"
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

  "$TOOLBX" --assumeyes create >/dev/null \
    || fail "Toolbx couldn't create default container"
}


# Creates a container with specific name and image
#
# Parameters:
# ===========
# - image - name of the image
# - container_name - name of the container
function create_image_container() {
  local image
  local container_name

  image="$1"
  container_name="$2"

  "$TOOLBX" --assumeyes create --container "${container_name}" --image "${image}" >/dev/null \
    || fail "Toolbx couldn't create container '$container_name'"
}


function start_container() {
  local container_name
  container_name="$1"

  podman start "$container_name" >/dev/null \
    || fail "Podman couldn't start the container '$container_name'"
}


# Checks if a Toolbx container started
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
    run --separate-stderr podman logs "$container_name"

    # shellcheck disable=SC2154
    if [ "$status" -ne 0 ]; then
      fail "Failed to invoke 'podman logs'"
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
  podman start "$container_name" >/dev/null \
    || fail "Podman couldn't start the container '$container_name'"
  podman stop "$container_name" >/dev/null \
    || fail "Podman couldn't stop the container '$container_name'"
}


# Returns the name of the latest created container
function get_latest_container_name() {
  podman ps --latest --format "{{ .Names }}"
}


function list_images() {
  podman images --all --format "{{.ID}}" | wc --lines
}


function list_containers() {
  podman ps --all --quiet | wc --lines
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

  local system_version="$VERSION_ID"
  [ "$ID" = "arch" ] && system_version="latest"

  echo "$system_version"
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

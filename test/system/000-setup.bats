#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'
load 'libs/helpers'

@test "test suite: Setup" {
    # Cache the default image for the system
    _pull_and_cache_distro_image $(get_system_id) $(get_system_version) || die
    # Cache all images that will be needed during the tests
    _pull_and_cache_distro_image fedora 32 || die
    _pull_and_cache_distro_image rhel 8.4 || die
    _pull_and_cache_distro_image busybox || die

    # Prepare localy hosted image registries
    # The registries need to live in a separate instance of Podman to prevent
    # them from being removed when cleaning state between test cases.

    # Create certificates for HTTPS
    mkdir -p "$CERTS_DIR"
    openssl req \
        -newkey rsa:4096 \
        -nodes -sha256 \
        -keyout "$CERTS_DIR"/domain.key \
        -addext "subjectAltName = DNS:localhost" \
        -x509 \
        -days 365 \
        -subj '/' \
        -out "$CERTS_DIR"/domain.crt
    assert [ $? -eq 0 ]

    # Add certificate to Podman's trusted certificates (rootless)
    mkdir -p ~/.config/containers/certs.d/localhost:50000
    cp "$CERTS_DIR"/domain.crt ~/.config/containers/certs.d/localhost:50000
    mkdir -p ~/.config/containers/certs.d/localhost:50001
    cp "$CERTS_DIR"/domain.crt ~/.config/containers/certs.d/localhost:50001

    # Create a registry user
    # username: user
    # password: user
    mkdir -p "$AUTH_DIR"
    htpasswd -Bbn user user > "$AUTH_DIR"/htpasswd
    assert [ $? -eq 0 ]

    # Create separate Podman root
    mkdir -p "$PODMAN_REG_ROOT"

    # Create a Docker registry without authentication
    run $PODMAN --root "$PODMAN_REG_ROOT" run -d \
        --rm \
        --name docker-registry-noauth \
        --privileged \
        -v "$CERTS_DIR":/certs \
        -e REGISTRY_HTTP_ADDR=0.0.0.0:443 \
        -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
        -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
        -p 50000:443 \
        docker.io/library/registry:2
    assert_success
    echo "# Podman logs - docker-registry-noauth"
    $PODMAN --root "$PODMAN_REG_ROOT" logs docker-registry-noauth

    # Create a Docker registry with authentication
    run $PODMAN --root "$PODMAN_REG_ROOT" run -d \
        --rm \
        --name docker-registry-auth \
        --privileged \
        -v "$AUTH_DIR":/auth \
        -e REGISTRY_AUTH=htpasswd \
        -e REGISTRY_AUTH_HTPASSWD_REALM="Registry Realm" \
        -e REGISTRY_AUTH_HTPASSWD_PATH=/auth/htpasswd \
        -v "$CERTS_DIR":/certs \
        -e REGISTRY_HTTP_ADDR=0.0.0.0:443 \
        -e REGISTRY_HTTP_TLS_CERTIFICATE=/certs/domain.crt \
        -e REGISTRY_HTTP_TLS_KEY=/certs/domain.key \
        -p 50001:443 \
        docker.io/library/registry:2
    assert_success
    echo "# Podman logs - docker-registry-auth"
    $PODMAN --root "$PODMAN_REG_ROOT" logs docker-registry-auth

    # Add UBI8 to created registries
    run $SKOPEO copy "dir:${IMAGE_CACHE_DIR}/fedora-toolbox-32" "docker://localhost:50000/fedora-toolbox:32"
    assert_success

    run $SKOPEO copy \
        --dest-creds user:user \
        "dir:${IMAGE_CACHE_DIR}/fedora-toolbox-32" "docker://localhost:50001/fedora-toolbox:32"
    assert_success

}

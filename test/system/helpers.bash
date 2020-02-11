#!/usr/bin/env bash

# Podman and Toolbox commands to run
PODMAN=${PODMAN:-podman}
TOOLBOX=${TOOLBOX:-toolbox}

# Helpful globals
LATEST_FEDORA_VERSION=${LATEST_FEDORA_VERSION:-"32"}
DEFAULT_FEDORA_VERSION=${DEFAULT_FEDORA_VERSION:-"f31"}
REGISTRY_URL=${REGISTRY_URL:-"registry.fedoraproject.org"}
TOOLBOX_DEFAULT_IMAGE=${TOOLBOX_DEFAULT_IMAGE:-"registry.fedoraproject.org/f31/fedora-toolbox:31"}
TOOLBOX_TIMEOUT=${TOOLBOX_TIMEOUT:-100}
PODMAN_TIMEOUT=${PODMAN_TIMEOUT:-100}

# Colors
LGC='\033[1;32m' # Light Green Color
LBC='\033[1;34m' # Light Blue Color
NC='\033[0m' # No Color

# Basic setup
function basic_setup() {
    echo "# [basic_setup]" >&2
    # Make sure desired images are present
    if [ -z "$found_needed_image" ]; then
        run_podman pull "$TOOLBOX_DEFAULT_IMAGE"
    fi
}

function setup_with_one_container() {
    echo "# [setup_with_one_container]" >&2
    # Clean up all images except for the default one
    remove_all_images_but_default
    # Create a new (default) container if no other are present
    run_toolbox -y create
}

function basic_teardown() {
    echo "# [basic_teardown]" >&2
    # Clean up all containers
    remove_all_containers
    # Clean up all images except for the default one
    remove_all_images_but_default
}

# Set the default setup function
function setup() {
    basic_setup
}

function teardown() {
    basic_teardown
}


################
#  run_podman  #  Invoke $PODMAN, with timeout, using BATS 'run'
################
#
# This is the preferred mechanism for invoking podman: first, it
# invokes $PODMAN, which may be 'podman-remote' or '/some/path/podman'.
#
# Second, we log the command run and its output. This doesn't normally
# appear in BATS output, but it will if there's an error.
#
# Next, we check exit status. Since the normal desired code is 0,
# that's the default; but the first argument can override:
#
#     run_podman 125  nonexistent-subcommand
#     run_podman '?'  some-other-command       # let our caller check status
#
# Since we use the BATS 'run' mechanism, $output and $status will be
# defined for our caller.
#
function run_podman() {
    # Number as first argument = expected exit code; default 0
    expected_rc=0
    case "$1" in
        [0-9])           expected_rc=$1; shift;;
        [1-9][0-9])      expected_rc=$1; shift;;
        [12][0-9][0-9])  expected_rc=$1; shift;;
        '?')             expected_rc=  ; shift;;  # ignore exit code
    esac

    # stdout is only emitted upon error; this echo is to help a debugger
    echo "\$ $PODMAN $*"
    run timeout --foreground -v --kill=10 $PODMAN_TIMEOUT $PODMAN "$@" 3>/dev/null
    # without "quotes", multiple lines are glommed together into one
    if [ -n "$output" ]; then
        echo "$output"
    fi
    if [ "$status" -ne 0 ]; then
        echo -n "[ rc=$status ";
        if [ -n "$expected_rc" ]; then
            if [ "$status" -eq "$expected_rc" ]; then
                echo -n "(expected) ";
            else
                echo -n "(** EXPECTED $expected_rc **) ";
            fi
        fi
        echo "]"
    fi

    if [ -n "$expected_rc" ]; then
        if [ "$status" -ne "$expected_rc" ]; then
            die "exit code is $status; expected $expected_rc"
        fi
    fi
}

function run_toolbox() {
    # Number as first argument = expected exit code; default 0
    expected_rc=0
    case "$1" in
        [0-9])           expected_rc=$1; shift;;
        [1-9][0-9])      expected_rc=$1; shift;;
        [12][0-9][0-9])  expected_rc=$1; shift;;
        '?')             expected_rc=  ; shift;;  # ignore exit code
    esac

    # stdout is only emitted upon error; this echo is to help a debugger
    echo "\$ $TOOLBOX $*"
    run timeout --foreground -v --kill=10 $TOOLBOX_TIMEOUT $TOOLBOX "$@" 3>/dev/null
    # without "quotes", multiple lines are glommed together into one
    if [ -n "$output" ]; then
        echo "$output"
    fi
    if [ "$status" -ne 0 ]; then
        echo -n "[ rc=$status ";
        if [ -n "$expected_rc" ]; then
            if [ "$status" -eq "$expected_rc" ]; then
                echo -n "(expected) ";
            else
                echo -n "(** EXPECTED $expected_rc **) ";
            fi
        fi
        echo "]"
    fi

    if [ -n "$expected_rc" ]; then
        if [ "$status" -ne "$expected_rc" ]; then
            die "exit code is $status; expected $expected_rc"
        fi
    fi
}

# Functions to prepare environment


function create_toolbox() {
    echo "# [create_toolbox]"

    local numberof="$1"
    local naming="$2"
    local image="$3"

    if [ "$numberof" = "" ]; then
        numberof=${numberof:-1}
    fi

    if [ "$image" = "" ]; then
        image=$TOOLBOX_DEFAULT_IMAGE
    fi

    for i in $(seq "$numberof"); do
        if [ "$naming" = "" ]; then
            run_toolbox '?' -y create -i "$image"
        else
            run_toolbox '?' -y create -c "$naming-$i" -i "$image"
        fi
    done
}


function get_images() {
    echo "# [get_images]"

    local numberof="$1"
    local image=""

    if [ "$numberof" = "" ]; then
        numberof=${numberof:-1}
    fi

    for i in $(seq $numberof); do
        local version=$[$LATEST_FEDORA_VERSION-$i]
        image="$REGISTRY_URL/f$version/fedora-toolbox:$version"
        run_podman pull "$image" || echo "Podman couldn't pull the image."
    done
}


function remove_all_images() {
    echo "# [remove_all_images]"
    run_podman images --all --format '{{.Repository}}:{{.Tag}} {{.ID}}'
    for line in "${lines[@]}"; do
        set $line
        run_podman rmi --force "$1" >/dev/null 2>&1 || true
        run_podman rmi --force "$2" >/dev/null 2>&1 || true
    done
}

function remove_all_images_but_default() {
    echo "# [remove_all_images_but_default]"
    found_needed_image=1
    run_podman images --all --format '{{.Repository}}:{{.Tag}} {{.ID}}'
    for line in "${lines[@]}"; do
        set $line
        if [ "$1" == "$TOOLBOX_DEFAULT_IMAGE" ]; then
            found_needed_image=1
        else
            run_podman rmi --force "$1" >/dev/null 2>&1 || true
        fi
    done
}


function get_image_name() {
    echo "# [get_image_name]"
    local type="$1"
    local version="$2"

    if [ -z "$type" ]; then
        type=${type:-fedora}
    fi

    if [ -z "$version" ]; then
        version=${version:-$DEFAULT_FEDORA_VERSION}
    fi

    case "$type" in
        fedora)
            echo "$REGISTRY_URL/f$version/fedora-toolbox:$version"
            ;;
    esac
}

function remove_all_containers() {
    echo "# [remove_all_containers]"
    run_toolbox '?' rm --all --force
}


# BATS specific functions

#########
#  die  #  Abort with helpful message
#########
function die() {
    echo "#/vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv"  >&2
    echo "#| FAIL: $*"                                           >&2
    echo "#\\^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^" >&2
    false
}


########
#  is  #  Compare actual vs expected string; fail w/diagnostic if mismatch
########
#
# Compares given string against expectations, using 'expr' to allow patterns.
#
# Examples:failed to inspect
#
#   is "$actual" "$expected" "descriptive test name"
#   is "apple" "orange"  "name of a test that will fail in most universes"
#   is "apple" "[a-z]\+" "this time it should pass"
#
function is() {
    local actual="$1"
    local expect="$2"
    local testname="${3:-FIXME}"

    if [ -z "$expect" ]; then
        if [ -z "$actual" ]; then
            return
        fi
        expect='[no output]'
    elif expr "$actual" : "$expect" >/dev/null; then
        return
    fi

    # This is a multi-line message, which may in turn contain multi-line
    # output, so let's format it ourself, readably
    local -a actual_split
    readarray -t actual_split <<<"$actual"
    printf "#/vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv\n" >&2
    printf "#|     FAIL: $testname\n"                          >&2
    printf "#| expected: '%s'\n" "$expect"                     >&2
    printf "#|   actual: '%s'\n" "${actual_split[0]}"          >&2
    local line
    for line in "${actual_split[@]:1}"; do
        printf "#|         > '%s'\n" "$line"                   >&2
    done
    printf "#\\^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^\n" >&2
    false
}

#################
#  parse_table  #  Split a table on '|' delimiters; return space-separated
#################
#
# See sample .bats scripts for examples. The idea is to list a set of
# tests in a table, then use simple logic to iterate over each test.
# Columns are separated using '|' (pipe character) because sometimes
# we need spaces in our fields.
#
function parse_table() {
    while read line; do
        test -z "$line" && continue

        declare -a row=()
        while read col; do
            dprint "col=<<$col>>"
            row+=("$col")
        done <  <(echo "$line" | tr '|' '\012' | sed -e 's/^ *//' -e 's/\\/\\\\/g')

        printf "%q " "${row[@]}"
        printf "\n"
    done <<<"$1"
}

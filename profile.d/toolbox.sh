# shellcheck shell=sh

# shellcheck disable=SC2153
[ "${BASH_VERSION:-}" != "" ] || [ "${ZSH_VERSION:-}" != "" ] || return 0
[ "$PS1" != "" ] || return 0

toolbox_config="$HOME/.config/toolbox"
host_welcome_stub="$toolbox_config/host-welcome-shown"
toolbox_welcome_stub="$toolbox_config/toolbox-welcome-shown"

# shellcheck disable=SC1091
# shellcheck disable=SC2046
eval $(
          if [ -f /etc/os-release ]; then
              . /etc/os-release
          else
              . /usr/lib/os-release
          fi

          echo ID="$ID"
          echo PRETTY_NAME="\"$PRETTY_NAME\""
          echo VARIANT_ID="$VARIANT_ID"
      )

if [ -f /run/ostree-booted ] \
   && ! [ -f "$host_welcome_stub" ] \
   && [ "${ID}" = "fedora" ] \
   && { [ "${VARIANT_ID}" = "workstation" ] \
        || [ "${VARIANT_ID}" = "silverblue" ] \
        || [ "${VARIANT_ID}" = "kinoite" ] \
        || [ "${VARIANT_ID}" = "sericea" ]; }; then
    echo ""
    echo "Welcome to ${PRETTY_NAME:-Linux}."
    echo ""
    echo "This terminal is running on the host system. You may want to try"
    echo "out the Toolbx for a directly mutable environment that allows "
    echo "package installation with DNF."
    echo ""
    printf "For more information, see the "
    # shellcheck disable=SC1003
    printf '\033]8;;https://docs.fedoraproject.org/en-US/fedora-silverblue/toolbox/\033\\documentation\033]8;;\033\\'
    printf ".\n"
    echo ""

    mkdir -p "$toolbox_config"
    touch "$host_welcome_stub"
fi

if [ -f /run/.containerenv ] \
   && [ -f /run/.toolboxenv ]; then
    [ "${BASH_VERSION:-}" != "" ] && PS1=$(printf "\[\033[35m\]⬢ \[\033[0m\]%s" "[\u@\h \W]\\$ ")
    [ "${ZSH_VERSION:-}" != "" ] && PS1=$(printf "\033[35m⬢ \033[0m%s" "[%n@%m]%~%# ")

    if ! [ -f "$toolbox_welcome_stub" ]; then
        echo ""
        echo "Welcome to the Toolbx; a container where you can install and run"
        echo "all your tools."
        echo ""

        if [ "${ID}" = "fedora" ]; then
            echo " - Use DNF in the usual manner to install command line tools."
            echo " - To create a new tools container, run 'toolbox create'."
            echo ""
            printf "For more information, see the "
            # shellcheck disable=SC1003
            printf '\033]8;;https://docs.fedoraproject.org/en-US/fedora-silverblue/toolbox/\033\\documentation\033]8;;\033\\'
            printf ".\n"
        else
            echo " - To create a new tools container, run 'toolbox create'."
        fi

        echo ""

        mkdir -p "$toolbox_config"
        touch "$toolbox_welcome_stub"
    fi

    if ! [ -f /etc/profile.d/vte.sh ] && [ -z "$PROMPT_COMMAND" ] && [ "${VTE_VERSION:-0}" -ge 3405 ]; then
        case "$TERM" in
            xterm*|vte*)
                [ -n "${BASH_VERSION:-}" ] && PROMPT_COMMAND=" "
                ;;
        esac
    fi

    if [ "$TERM" != "" ]; then
        error_message="Error: terminfo entry not found for $TERM"
        term_without_first_character="${TERM#?}"
        term_just_first_character="${TERM%"$term_without_first_character"}"
        terminfo_sub_directory="$term_just_first_character/$TERM"

        if [ "$TERMINFO" = "" ]; then
          ! [ -e "/usr/share/terminfo/$terminfo_sub_directory" ] \
            && ! [ -e "/lib/terminfo/$terminfo_sub_directory" ] \
            && ! [ -e "$HOME/.terminfo/$terminfo_sub_directory" ] \
            && echo "$error_message" >&2
        else
          ! [ -e "$TERMINFO/$terminfo_sub_directory" ] \
            && echo "$error_message" >&2
        fi
    fi
fi

unset ID
unset PRETTY_NAME
unset VARIANT_ID
unset toolbox_config
unset host_welcome_stub
unset toolbox_welcome_stub

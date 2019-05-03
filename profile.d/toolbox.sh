[ "$BASH_VERSION" != "" ] || [ "$ZSH_VERSION" != "" ] || return 0
[ "$PS1" != "" ] || return 0

toolbox_config="$HOME/.config/toolbox"
host_welcome_stub="$toolbox_config/host-welcome-shown"
toolbox_welcome_stub="$toolbox_config/toolbox-welcome-shown"

if [ -f /run/ostree-booted ] \
   && ! [ -f "$host_welcome_stub" ]; then
    echo ""
    echo "Welcome to Fedora Silverblue. This terminal is running on the"
    echo "immutable host system. You may want to try out the Toolbox for a"
    echo "more traditional environment that allows package installation"
    echo "with DNF."
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
   && ! [ -f "$toolbox_welcome_stub" ] \
   && [ -f /run/.toolboxenv ]; then
    echo ""
    echo "Welcome to the Toolbox; a container where you can install and run"
    echo "all your tools."
    echo ""
    echo " - Use DNF in the usual manner to install command line tools."
    echo " - To create a new tools container, run 'toolbox create'."
    echo ""
    printf "For more information, see the "
    # shellcheck disable=SC1003
    printf '\033]8;;https://docs.fedoraproject.org/en-US/fedora-silverblue/toolbox/\033\\documentation\033]8;;\033\\'
    printf ".\n"
    echo ""

    mkdir -p "$toolbox_config"
    touch "$toolbox_welcome_stub"
fi

unset toolbox_config
unset host_welcome_stub
unset toolbox_welcome_stub

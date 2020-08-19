[ "$BASH_VERSION" != "" ] || [ "$ZSH_VERSION" != "" ] || return 0
[ "$PS1" != "" ] || return 0

# Add an indicator to the prompt that we're inside the toolbox
if [ -f /run/.containerenv ] \
   && [ -f /run/.toolboxenv ]; then
    PS1=$(printf "\[\033[35m\]⬢\[\033[0m\]%s" "[\u@\h \W]\\$ ")
fi

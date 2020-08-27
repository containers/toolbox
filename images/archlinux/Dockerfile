FROM docker.io/library/archlinux:latest

ENV NAME=archlinux-toolbox VERSION=latest
LABEL com.github.containers.toolbox="true" \
      com.github.debarshiray.toolbox="true" \
      name="$FCG/$NAME" \
      version="$VERSION" \
      usage="This image is meant to be used with the toolbox command" \
      summary="Base image for creating Arch Linux toolbox containers" \
      maintainer="Erazem Kokot <contact@erazem.eu>"

COPY extra-packages /
RUN pacman -Syu --noconfirm $(<extra-packages)
RUN rm /extra-packages

# Clean up all local caches
RUN pacman -Scc --noconfirm

# Allow sudo for users in group "sudo"
RUN groupadd -f sudo
RUN echo '%sudo ALL=(ALL) ALL' >/etc/sudoers.d/group-sudo

CMD /bin/sh
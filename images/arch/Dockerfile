FROM docker.io/archlinux/base:latest
MAINTAINER Erazem Kokot <contact at erazem dot eu>

ENV NAME=arch-toolbox VERSION=rolling
LABEL com.github.containers.toolbox="true" \
      com.github.debarshiray.toolbox="true" \
      name="$FCG/$NAME" \
      version="$VERSION" \
      usage="This image is meant to be used with the toolbox command" \
      summary="Base image for creating Arch toolbox containers" \
      maintainer="Erazem Kokot <contact at erazem dot eu>"

RUN pacman -Syu --noconfirm

# Install packages from the extra-packages file
COPY extra-packages /
RUN pacman -Sy --noconfirm $(<extra-packages)
RUN rm /extra-packages

# Clean up all local caches
RUN pacman -Scc --noconfirm

CMD /bin/sh

FROM debian:10

ENV NAME=debian-toolbox VERSION=10
LABEL com.github.containers.toolbox="true" \
      com.github.debarshiray.toolbox="true" \
      com.redhat.component="$NAME" \
      name="$NAME" \
      version="$VERSION" \
      usage="This image is meant to be used with the toolbox command" \
      summary="Base image for creating Fedora toolbox containers" \
      maintainer="Debarshi Ray <rishi@fedoraproject.org>"

COPY README.md /

WORKDIR /

RUN apt -y update

COPY missing-docs /
RUN apt -y reinstall $(cat missing-docs)
RUN rm /missing-docs

COPY extra-packages /
RUN apt -y install $(cat extra-packages)
RUN rm /extra-packages

RUN apt clean all

COPY exempt_group /etc/sudoers.d/

CMD /bin/sh

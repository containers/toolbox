FROM docker.io/library/debian:buster

ENV NAME=debian-toolbox VERSION=buster
LABEL com.github.containers.toolbox="true" \
      com.github.debarshiray.toolbox="true" \
      name="$NAME" \
      version="$VERSION" \
      usage="This image is meant to be used with the toolbox command" \
      summary="Base image for creating Debian buster (10) toolbox containers" \
      maintainer="Debarshi Ray <rishi@fedoraproject.org>"

COPY README.md /

RUN apt-get update

RUN apt-get -y upgrade

COPY extra-packages /
RUN apt-get -y install $(cat extra-packages | xargs)
RUN rm /extra-packages

RUN sed -i -e 's/ ALL$/ NOPASSWD:ALL/' /etc/sudoers

RUN apt-get clean

RUN echo VARIANT_ID=container >> /etc/os-release

CMD /bin/sh

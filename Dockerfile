FROM registry.fedoraproject.org/fedora:28

ENV NAME=fedora-toolbox VERSION=28
LABEL com.redhat.component="$NAME" \
      name="$FGC/$NAME" \
      version="$VERSION" \
      usage="This image is meant to be used with the fedora-toolbox command" \
      summary="Base image for creating Fedora toolbox containers" \
      maintainer="Debarshi Ray <rishi@fedoraproject.org>"

COPY README.md /

RUN sed -i '/tsflags=nodocs/d' /etc/dnf/dnf.conf
RUN dnf -y swap coreutils-single coreutils-full

COPY extra-packages /
RUN packages=; while read -r package; do packages="$packages $package"; done \
        <extra-packages; \
    dnf -y install $packages
RUN rm /extra-packages

CMD /bin/sh

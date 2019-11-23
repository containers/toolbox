FROM docker.io/library/centos:8

ENV NAME=centos-toolbox VERSION=8
LABEL com.github.containers.toolbox="true" \
      com.github.debarshiray.toolbox="true" \
      com.redhat.component="$NAME" \
      name="$FGC/$NAME" \
      version="$VERSION" \
      usage="This image is meant to be used with the toolbox command" \
      summary="Base image for creating CentOS 8 toolbox containers" \
      maintainer="Debarshi Ray <rishi@fedoraproject.org>"

COPY README.md /

RUN sed -i '/tsflags=nodocs/d' /etc/yum.conf
RUN yum -y swap coreutils-single coreutils

COPY missing-docs /
RUN yum -y reinstall $(<missing-docs)
RUN rm /missing-docs

COPY extra-packages /
RUN yum -y install $(<extra-packages)
RUN rm /extra-packages

RUN yum clean all

RUN echo VARIANT_ID=container >> /etc/os-release
RUN ln -s /etc/os-release /usr/lib/os-release

CMD /bin/sh

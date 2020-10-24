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
COPY missing-docs /
COPY extra-packages /

RUN sed -i '/tsflags=nodocs/d' /etc/yum.conf
RUN yum -y swap coreutils-single coreutils && \
    yum -y reinstall $(<missing-docs) && \
    yum -y install $(<extra-packages) && \
    yum clean all
RUN rm /extra-packages /missing-docs

RUN echo VARIANT_ID=container >> /etc/os-release

CMD /bin/sh

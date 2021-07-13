FROM docker.io/library/ubuntu:16.04

LABEL com.github.containers.toolbox="true" \
      com.github.debarshiray.toolbox="true" \
      name="ubuntu-toolbox" \
      version="16.04" \
      usage="This image is meant to be used with the toolbox command" \
      summary="Base image for creating Ubuntu toolbox containers" \
      maintainer="Ievgen Popovych <jmennius@gmail.com>"

COPY README.md /

# remove some container-optimized apt configuration;
# removing docker-gzip-indexes specifically helps with command-not-found
RUN rm /etc/apt/apt.conf.d/docker-gzip-indexes /etc/apt/apt.conf.d/docker-no-languages

# Install ubuntu-minimal & ubuntu-standard (recommendations are welcome)
# Ask resolvconf maintainer script not to touch resolve.conf file
# Install extra packages as well as libnss-myhostname
COPY extra-packages /
RUN apt-get update && \
    echo "resolvconf resolvconf/linkify-resolvconf boolean false" | debconf-set-selections && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
        ubuntu-minimal ubuntu-standard \
        libnss-myhostname \
        $(cat extra-packages | xargs) && \
    rm -rd /var/lib/apt/lists/*
RUN rm /extra-packages

# Enable myhostname plugin for clean hostname resolution w/o patching hosts (at least for sudo)
RUN sed -i '/^hosts:/s/$/ myhostname/' /etc/nsswitch.conf

# Allow authentication with empty password, promptless
RUN sed -i '/^auth.*pam_unix.so/s/nullok_secure/try_first_pass nullok/' /etc/pam.d/common-auth

# Fix empty bind-mount to clear selinuxfs (see #337)
RUN mkdir /usr/share/empty

CMD /bin/sh

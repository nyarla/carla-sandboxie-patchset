FROM ubuntu:20.04

ENV DEBIAN_FRONTEND noninteractive
ENV APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE 1

# Setup for stable wine
RUN   apt-get update \
  &&  apt-get install -y wget gnupg software-properties-common \
  &&  (wget -qO- https://dl.winehq.org/wine-builds/winehq.key | apt-key add -) \
  &&  apt-add-repository 'deb https://dl.winehq.org/wine-builds/ubuntu/ focal main' \
  &&  dpkg --add-architecture i386 \
  &&  apt-get update -qq \
  &&  apt-get install -y -o APT::Immediate-Configure=false libc6 libc6:i386 libgcc-s1:i386 \
  &&  apt-get install -y -f \
  && rm -rf /var/lib/apt/lists/*

# Install packages
RUN   apt-get update \
  &&  apt-get install -y \
        autoconf \
        binutils-mingw-w64-x86-64 \
        build-essential \
        cmake \
        curl \
        g++-mingw-w64-x86-64 \
        git \
        jq \
        libglib2.0-dev \
        llvm \
        locales \
        mingw-w64 \
        qttools5-dev-tools \
        winehq-stable \
        zip \
    && rm -rf /var/lib/apt/lists/*

# Install build dependences
RUN   sed -Ei 's/^# deb-src /deb-src /' /etc/apt/sources.list \
  &&  apt-get update \
  &&  apt-get install -y \
    gcc-multilib \
    libcrypt-dev \
    libffi-dev \
    libssl-dev \
    uuid-dev \
  && rm -rf /var/lib/apt/lists/*

# Setup build environment
RUN localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
RUN sed -i 's|define __MINGW_FORTIFY_LEVEL [0-9]|define __MINGW_FORTIFY_LEVEL 0|' /usr/share/mingw-w64/include/_mingw_mac.h

# Setup build account
RUN useradd -m builder
ENV LANG en_US.UTF-8

WORKDIR /home/builder
ENV HOME /home/builder
ENV TRAVIS_BUILD_DIR /home/builder/src
ENV BOOTSTRAP_VERSION 2
ENV TARGET=win64

ENTRYPOINT [ "/bin/bash" ]

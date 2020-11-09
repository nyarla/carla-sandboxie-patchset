FROM ubuntu:20.04

ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update \
  && apt-get install -y \
    build-essential autoconf mingw-w64 libglib2.0-dev cmake \
    zip xz-utils zstd \
    wget curl \
    locales \
    software-properties-common \
  && rm -rf /var/lib/apt/lists/*

RUN dpkg --add-architecture i386 \
  && (wget -O - https://dl.winehq.org/wine-builds/winehq.key | apt-key add -) \
  && apt-add-repository 'deb https://dl.winehq.org/wine-builds/ubuntu/ focal main' \
  && apt-get update \
  && apt-get install -y --install-recommends winehq-stable \
  && apt-get install -y winetricks \
  && rm -rf /var/lib/apt/lists/*

RUN localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
RUN sed -i 's|define __MINGW_FORTIFY_LEVEL [0-9]|define __MINGW_FORTIFY_LEVEL 0|' /usr/share/mingw-w64/include/_mingw_mac.h

RUN useradd -m builder
ENV LANG en_US.UTF-8

WORKDIR /home/builder

ENTRYPOINT [ "/bin/bash" ]

#!/usr/bin/env bash

set -e -u -o pipefail

cd ~/src/Carla

if test ! -e ~/builds/.done-patch ; then
  patch -p1 -i ../patches/fix-glib-patch.patch
  patch -p1 -i ../patches/fix-buildscript.patch
  patch -p1 -i ../patches/fix-pyqt-deps.patch
  patch -p1 -i ../patches/fix-pack-win.patch
  patch -p1 -i ../patches/sandboxie.patch
  patch -p1 -i ../patches/sandboxie-discovery.patch
  touch ~/builds/.done-patch
fi

test -e ~/builds/.done-deps || (bash data/windows/build-deps.sh && touch ~/builds/.done-deps)
test -e ~/builds/.done-pyqt || (bash data/windows/build-pyqt.sh && touch ~/builds/.done-pyqt)

if test ! -e ~/builds/.done-cx_Freeze ; then
  cd ~/builds/msys2-i686/mingw32/lib/python3.8/site-packages/cx_Freeze && patch -p2 -i ~/src/patches/cx_Freeze.patch
  cd ~/builds/msys2-x86_64/mingw64/lib/python3.8/site-packages/cx_Freeze && patch -p2 -i ~/src/patches/cx_Freeze.patch
  touch ~/builds/.done-cx_Freeze
  cd ~/src/Carla
fi

cp ~/builds/msys2-i686/mingw32/bin/libffi-{7,6}.dll
cp ~/builds/msys2-x86_64/mingw64/bin/libffi-{7,6}.dll

make clean
rm -rf data/windows/Carla*
bash data/windows/build-win.sh 64
env CARLA_DEV=1 bash data/windows/pack-win.sh 64

for app in Carla Carla.vst Carla.lv2 CarlaControl; do
  cp -R data/windows/$app ../dist/
done

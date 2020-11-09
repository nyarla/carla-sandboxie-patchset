#!/usr/bin/env bash

set -e -u -o pipefail

cd "$(dirname "${0}")/../dist"

GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o carla-bridge-win64.exe ../helpers/carla-bridge.go 
GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui" -o carla-bridge-win32.exe ../helpers/carla-bridge.go

GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o carla-discovery-win64.exe ../helpers/carla-discovery.go
GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui" -o carla-discovery-win32.exe ../helpers/carla-discovery.go

for app in Carla Carla.vst Carla.lv2 ; do
  cd $app
  if ! test -e _carla-discovery-win32.exe ; then
    mv {,_}carla-discovery-win32.exe
    cp ../carla-discovery-win32.exe .
  fi

  if ! test -e _carla-discovery-win64.exe ; then
    mv {,_}carla-discovery-win64.exe
    cp ../carla-discovery-win64.exe .
  fi

  if ! test -e _carla-bridge-win32.exe ; then
    mv {,_}carla-bridge-win32.exe
    cp ../carla-bridge-win32.exe .
  fi

  if ! test -e _carla-bridge-win64.exe ; then
    mv carla-bridge-native.exe carla-bridge-win64.exe
    mv {,_}carla-bridge-win64.exe
    cp ../carla-bridge-win64.exe .
  fi
  
  cd ../
done

rm *.exe

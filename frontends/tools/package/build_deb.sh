#!/bin/bash

DIR=fantasim-sdl_${FANTASIM_SDL_VERSION}_linux

if [ -d "$DIR" ]; then rm -rf "$DIR"; fi
cp -r "frontends/tools/package/fantasim-sdl.deb" $DIR

mkdir -p $DIR/usr/local/bin
cp frontends/sdl/fantasim-sdl $DIR/usr/local/bin

rpl e34f19fc-289d-4fb9-b134-ced07a29a273 $FANTASIM_SDL_VERSION $DIR/DEBIAN/control
rpl e34f19fc-299d-4fb9-b334-aed07b29a273 $TRAVIS_BUILD_NUMBER $DIR/DEBIAN/control

dpkg-deb --build $DIR
rm -rf $DIR

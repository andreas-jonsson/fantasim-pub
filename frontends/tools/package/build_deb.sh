#!/bin/bash

#DIR=fantasim-sdl_$FANTASIM_SDL_VERSION-$TRAVIS_BUILD_NUMBER
DIR=fantasim-sdl_$FANTASIM_SDL_VERSION

if [ -d "$DIR" ]; then rm -rf "$DIR"; fi
cp -r "frontends/tools/package/fantasim-$FANTASIM_SDL_VERSION-x" $DIR

mkdir -p $DIR/usr/local/bin
cp fantasim-sdl $DIR/usr/local/bin

rpl e34f19fc-289d-4fb9-b134-ced07a29a273 $FANTASIM_SDL_VERSION $DIR/DEBIAN/control
rpl e34f19fc-299d-4fb9-b334-aed07b29a273 $TRAVIS_BUILD_NUMBER $DIR/DEBIAN/control

dpkg-deb --build $DIR
rm -rf $DIR

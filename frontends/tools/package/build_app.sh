#!/bin/bash

if [ -d Fantasim-SDL.app ]; then rm -rf Fantasim-SDL.app; fi
cp -r "frontends/tools/package/Fantasim-SDL.app" Fantasim-SDL.app

rpl e34f49fc-229d-4fb9-b134-ced05a29a271 $FANTASIM_SDL_VERSION frontends/tools/package/Fantasim-SDL.app/Contents/Info.plist
rpl e34f49fc-129d-4fb9-b134-ced05a22a270 $FANTASIM_SDL_SHORT_VERSION frontends/tools/package/Fantasim-SDL.app/Contents/Info.plist

cp frontends/sdl/fantasim-sdl Fantasim-SDL.app/Contents/MacOS
dylibbundler -od -b -x Fantasim-SDL.app/Contents/MacOS/fantasim-sdl -d Fantasim-SDL.app/Contents/libs

zip -rq fantasim-sdl_{$FANTASIM_SDL_VERSION}_osx.zip Fantasim-SDL.app
rm -rf Fantasim-SDL.app

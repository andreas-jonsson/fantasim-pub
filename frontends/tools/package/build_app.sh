#!/bin/bash

if [ -d Fantasim-SDL.app ]; then rm -rf Fantasim-SDL.app; fi
cp -r "frontends/tools/package/Fantasim-SDL.app" Fantasim-SDL.app

cp fantasim-sdl Fantasim-SDL.app/Contents/MacOS
dylibbundler -od -b -x Fantasim-SDL.app/Contents/MacOS/fantasim-sdl -d Fantasim-SDL.app/Contents/libs

zip -rq fantasim-sdl_osx.zip Fantasim-SDL.app
rm -rf Fantasim-SDL.app

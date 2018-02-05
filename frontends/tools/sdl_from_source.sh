#!/bin/bash

wget -q http://www.libsdl.org/release/SDL2-2.0.4.tar.gz
tar xf SDL2-*.tar.gz
cd SDL2-* && ./configure --prefix=$SDL_PREFIX && make && sudo make install

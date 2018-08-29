#!/bin/bash

export DEST=fantasim-wasm_${FANTASIM_SDL_VERSION}_js
rm -rf $DEST
rm -rf ${DEST}.zip

mkdir $DEST
GOARCH=wasm GOOS=js go build -o ${DEST}/main.wasm frontends/wasm/main.go
cp frontends/wasm/index.html ${DEST}
cp frontends/wasm/wasm_exec.js ${DEST}

zip -rq ${DEST}.zip ${DEST}
rm -rf ${DEST}
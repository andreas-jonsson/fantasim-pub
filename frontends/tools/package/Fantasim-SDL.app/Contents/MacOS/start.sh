#!/bin/bash

osascript -e 'tell app "Terminal" to do script "\"'"$(dirname "$0")/fantasim-sdl" -url=$1'\""'

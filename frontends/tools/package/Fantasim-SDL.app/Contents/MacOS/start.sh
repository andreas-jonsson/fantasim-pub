#!/bin/bash

# echo $0 > /Users/andreasjonsson/Desktop/test.txt
# echo $1 >> /Users/andreasjonsson/Desktop/test.txt

# $(dirname "$0")/fantasim-sdl -url="$1"

osascript <<EOD
	on open location this_URL
		do shell script "/Users/andreasjonsson/Downloads/Fantasim-SDL/Contents/MacOS/tramp.sh '" & this_URL & "'"		
	end open location
EOD

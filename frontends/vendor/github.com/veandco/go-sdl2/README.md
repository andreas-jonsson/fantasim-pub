# SDL2 binding for Go [![Build Status](https://travis-ci.org/veandco/go-sdl2.svg?branch=master)](https://travis-ci.org/veandco/go-sdl2)
`go-sdl2` is SDL2 wrapped for Go users. It enables interoperability between Go and the SDL2 library which is written in C. That means the original SDL2 installation is required for this to work.


# Table of Contents
* [Documentation](#documentation)
* [Examples](#examples)
* [Requirements](#requirements)
* [Installation](#installation)
* [Cross-compiling](#cross-compiling)
* [FAQ](#faq)
* [License](#license)


# Documentation
* [GoDoc documentation for go-sdl2](https://godoc.org/github.com/veandco/go-sdl2)
* [Original SDL2 wiki](https://wiki.libsdl.org)


# Examples
```go
package main

import "github.com/veandco/go-sdl2/sdl"

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 200, 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}
```

For more complete examples, see https://github.com/veandco/go-sdl2-examples. You can run any of the `.go` files with `go run`.


# Requirements
* [SDL2](http://libsdl.org/download-2.0.php)
* [SDL2_image (optional)](http://www.libsdl.org/projects/SDL_image/)
* [SDL2_mixer (optional)](http://www.libsdl.org/projects/SDL_mixer/)
* [SDL2_ttf (optional)](http://www.libsdl.org/projects/SDL_ttf/)
* [SDL2_gfx (optional)](http://www.ferzkopp.net/wordpress/2016/01/02/sdl_gfx-sdl2_gfx/)

Below is some commands that can be used to install the required packages in
some Linux distributions. Some older versions of the distributions such as
Ubuntu 13.10 may also be used but it may miss an optional package such as
_libsdl2-ttf-dev_ on Ubuntu 13.10's case which is available in Ubuntu 14.04.

On __Ubuntu 14.04 and above__, type:\
`apt install libsdl2{,-image,-mixer,-ttf,-gfx}-dev`

On __Fedora 25 and above__, type:\
`yum install SDL2{,_image,_mixer,_ttf,_gfx}-devel`

On __Arch Linux__, type:\
`pacman -S sdl2{,_image,_mixer,_ttf,_gfx}`

On __Gentoo__, type:\
`emerge -av libsdl2 sdl2-{image,mixer,ttf,gfx}`

On __macOS__, install SDL2 via [Homebrew](http://brew.sh) like so:\
`brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config`

On __Windows__,
1. Install mingw-w64 from [Mingw-builds](http://mingw-w64.org/doku.php/download/mingw-builds)
    * Version: latest (at time of writing 6.3.0)
    * Architecture: x86_64
    * Threads: win32
    * Exception: seh
    * Build revision: 1
    * Destination Folder: Select a folder that your Windows user owns
2. Install SDL2 http://libsdl.org/download-2.0.php
    * Extract the SDL2 folder from the archive using a tool like [7zip](http://7-zip.org)
    * Inside the folder, copy the `i686-w64-mingw32` and/or `x86_64-w64-mingw32` depending on the architecture you chose into your mingw-w64 folder e.g. `C:\Program Files\mingw-w64\x86_64-6.3.0-win32-seh-rt_v5-rev1\mingw64`
3. Setup Path environment variable
    * Put your mingw-w64 binaries location into your system Path environment variable. e.g. `C:\Program Files\mingw-w64\x86_64-6.3.0-win32-seh-rt_v5-rev1\mingw64\bin` and `C:\Program Files\mingw-w64\x86_64-6.3.0-win32-seh-rt_v5-rev1\mingw64\x86_64-w64-mingw32\bin`
4. Open up a terminal such as `Git Bash` and run `go get -v github.com/veandco/go-sdl2/sdl`.
5. (Optional) You can repeat __Step 2__ for [SDL_image](https://www.libsdl.org/projects/SDL_image), [SDL_mixer](https://www.libsdl.org/projects/SDL_mixer), [SDL_ttf](https://www.libsdl.org/projects/SDL_ttf)
    * NOTE: pre-build the libraries for faster compilation by running `go install github.com/veandco/go-sdl2/{sdl,img,mix,ttf}`

* Or you can install SDL2 via [Msys2](https://msys2.github.io) like so:
`pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-SDL2{,_image,_mixer,_ttf,_gfx}`


# Installation
To get the bindings, type:\
`go get -v github.com/veandco/go-sdl2/sdl`\
`go get -v github.com/veandco/go-sdl2/img`\
`go get -v github.com/veandco/go-sdl2/mix`\
`go get -v github.com/veandco/go-sdl2/ttf`\
`go get -v github.com/veandco/go-sdl2/gfx`

or type this if you use Bash terminal:\
`go get -v github.com/veandco/go-sdl2/{sdl,img,mix,ttf}`

Due to `go-sdl2` being under active development, a lot of breaking changes are going to happen during v0.x. With [versioning system](https://github.com/golang/proposal/blob/master/design/24301-versioned-go.md) coming to Go soon, we'll make use of semantic versioning to ensure stability in the future.


# Cross-compiling
### Linux to Windows
1. Install MinGW toolchain.
   * On **Arch Linux**, it's simply `pacman -S mingw-w64`.
2. Download the SDL2 development package for MinGW [here](http://libsdl.org/download-2.0.php) (and the others like *SDL_image*, *SDL_mixer*, etc.. [here](https://www.libsdl.org/projects/) if you use them).
3. Extract the SDL2 development package and copy the `x86_64-w64-mingw32` folder inside recursively to the system's MinGW `x86_64-w64-mingw32` folder. You may also do the same for the `i686-w64-mingw32` folder.
   * On **Arch Linux**, it's `cp -r x86_64-w64-mingw32 /usr`.
4. Now you can start cross-compiling your Go program by running `env CGO_ENABLED="1" CC="/usr/bin/x86_64-w64-mingw32-gcc" GOOS="windows" CGO_LDFLAGS="-lmingw32 -lSDL2" CGO_CFLAGS="-D_REENTRANT" go build -x main.go`. You can change some of the parameters if you'd like to. In this example, it should produce a `main.exe` executable file.
5. Before running the program, you need to put `SDL2.dll` from the [SDL2 runtime package](http://libsdl.org/download-2.0.php) (For others like *SDL_image*, *SDL_mixer*, etc.., look for them [here](https://www.libsdl.org/projects/)) for Windows in the same folder as your executable.
6. Now you should be able to run the program using Wine or Windows!

### macOS to Windows
1. Install [Homebrew](https://brew.sh)
2. Install MinGW through Homebrew via `brew install mingw-w64`
3. Download the SDL2 development package for MinGW [here](http://libsdl.org/download-2.0.php) (and the others like *SDL_image*, *SDL_mixer*, etc.. [here](https://www.libsdl.org/projects/) if you use them).
4. Extract the SDL2 development package and copy the `x86_64-w64-mingw` folder inside recursively to the system's MinGW `x86_64-w64-mingw32 folder`. You may also do the same for the `i686-w64-mingw32` folder. The path to MinGW may be slightly different but the command should look something like `cp -r x86_64-w64-mingw32 /usr/local/Cellar/mingw-w64/5.0.3/toolchain-x86_64`.
5. Now you can start cross-compiling your Go program by running `env CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows CGO_LDFLAGS="-L/usr/local/Cellar/mingw-w64/5.0.3/toolchain-x86_64/x86_64-w64-mingw32/lib -lSDL2" CGO_CFLAGS="-I/usr/local/Cellar/mingw-w64/5.0.3/toolchain-x86_64/x86_64-w64-mingw32/include -D_REENTRANT" go build -x main.go`. You can change some of the parameters if you'd like to. In this example, it should produce a `main.exe` executable file.
6. Before running the program, you need to put `SDL2.dll` from the [SDL2 runtime package](http://libsdl.org/download-2.0.php) (For others like *SDL_image*, *SDL_mixer*, etc.., look for them [here](https://www.libsdl.org/projects/)) for Windows in the same folder as your executable.
7. Now you should be able to run the program using Wine or Windows!


# FAQ
__Why does the program not run on Windows?__
Try putting the [runtime libraries](http://libsdl.org/download-2.0.php) (e.g. `SDL2.dll` and friends) in the same folder as your program.

__Why does my program crash randomly or hang?__
Putting `runtime.LockOSThread()` at the start of your main() usually solves the problem (see [SDL2 FAQ](https://wiki.libsdl.org/FAQDevelopment) about multi-threading).

UPDATE: Recent update added a call queue system where you can put thread-sensitive code and have it called synchronously on the same OS thread. See the `render_queue` or `render_goroutines` examples from https://github.com/veandco/go-sdl2-examples to see how it works.

__Why can't SDL_mixer seem to play MP3 audio file?__
Your installed SDL_mixer probably doesn't support MP3 file.

On __macOS__, this is easy to correct. First remove the faulty mixer: `brew remove sdl2_mixer`, then reinstall it with the MP3 option: `brew install sdl2_mixer --with-flac --with-fluid-synth --with-libmikmod --with-libmodplug --with-smpeg2`. If necessary, check which options you can enable with `brew info sdl2_mixer`. You could also try installing sdl2\_mixer with mpg123 by running `brew install sdl2_mixer --with-mpg123`.

On __Other Operating Systems__, you will need to compile smpeg and SDL_mixer from source with the MP3 option enabled. You can find smpeg in the `external` directory of SDL_mixer. Refer to issue [#148](https://github.com/veandco/go-sdl2/issues/148) for instructions.

_Note that there seems to be a problem with SDL_mixer 2.0.2 so you can also try to revert back to 2.0.1 and see if it solves your problem_

__Does go-sdl2 support compiling on mobile platforms like Android and iOS?__
For Android, see https://github.com/gen2brain/go-sdl2-android-example.

There is currently no support for iOS yet.

__Why does my window not immediately render after creation?__
It appears the rendering subsystem needs some time to be able to present the drawn pixels. This can be workaround by adding delay using `sdl.Delay()` or put the rendering code inside a draw loop.

# License
Go-SDL2 is BSD 3-clause licensed.

/*
Copyright (C) 2017-2018 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

//go:generate go run ../common/data/generate.go

package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/andreas-jonsson/fantasim-pub/frontends/common/game"
	sys "github.com/andreas-jonsson/fantasim-pub/frontends/pocket/platform"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"golang.org/x/net/websocket"
)

const lobbyURL = "http://lobby.fantasim.net"

type Config struct {
	Host, Key string
}

func main() {
	app.Main(func(a app.App) {
		var (
			glctx  gl.Context
			sz     size.Event
			images *glutil.Images
			glimg  *glutil.Image
			cfg    = Config{"localhost", "00000000"}
		)

		const configPath = "/storage/emulated/0/Download/fantasim.json"
		newConfig := func() {
			data, _ := json.MarshalIndent(cfg, "", "\t")
			if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
				log.Println(err)
			}
		}

		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			log.Println(err)
			newConfig()
			return
		}

		if err = json.Unmarshal(data, &cfg); err != nil {
			log.Println(err)
			newConfig()
			return
		}

		//log.Println(jni.Test())

		paintDoneChan := make(chan struct{})
		exitChan := make(chan struct{})

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx = e.DrawContext.(gl.Context)
					images = glutil.NewImages(glctx)

					//serverAddress := "79.102.38.201:8080"
					//playerKey := "787C29A4B5296915"

					serverAddress := cfg.Host
					playerKey := cfg.Key

					log.Println("start cmd")
					_, err := os.StartProcess("am", []string{"start", "-a", "android.intent.action.VIEW", "-d", lobbyURL}, &os.ProcAttr{Dir: "/system/bin"})
					if err != nil {
						log.Println(err)
					}
					/*
						cmd := exec.Command("am", "start", "-a", "android.intent.action.VIEW", "-d", lobbyURL)
						if err := cmd.Start(); err != nil {
							log.Println(err)
						}
						log.Println("done!")
					*/

					go func() {
						defer func() {
							exitChan <- struct{}{}
						}()

						log.Println("Waiting for resize event...")

						//for sz := range sys.InputEventChan {
						//if ev, ok := sz.(size.Event); ok {
						//glimg = images.NewImage(sz.X, sz.Y)
						glimg = images.NewImage(640, 400)
						sys.ExternalBackBuffer = glimg.RGBA
						//break
						//}
						//}

						fmt.Println(sys.ExternalBackBuffer.Bounds().Max)

						origin := fmt.Sprintf("http://%s/", serverAddress)
						apiWs, err := websocket.Dial(fmt.Sprintf("ws://%s/api", serverAddress), "", origin)
						if err != nil {
							log.Fatalln(err)
						}
						defer apiWs.Close()

						infoWs, err := websocket.Dial(fmt.Sprintf("ws://%s/info", serverAddress), "", origin)
						if err != nil {
							log.Fatalln(err)
						}
						defer infoWs.Close()

						if err := json.NewEncoder(io.MultiWriter(apiWs, infoWs)).Encode(playerKey); err != nil {
							log.Fatalln(err)
						}

						if err := json.NewEncoder(apiWs).Encode("gob"); err != nil {
							log.Fatalln(err)
						}

						enc := gob.NewEncoder(apiWs)
						dec := gob.NewDecoder(apiWs)
						decInfo := gob.NewDecoder(infoWs)

						s, err := sys.InitPocket()
						if err != nil {
							log.Fatalln(err)
						}
						defer s.Quit()

						if err := game.Initialize(s, s); err != nil {
							log.Fatalln(err)
						}

						if err := game.Start(enc, dec, decInfo); err != nil {
							log.Fatalln(err)
						}
					}()

					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					sys.InputEventChan <- nil

				loop:
					for {
						select {
						case sys.PaintEventChan <- paintDoneChan:
							<-paintDoneChan
						case <-exitChan:
							break loop
						}
					}

					glimg.Release()
					images.Release()
					glctx = nil
				}
			case paint.Event:
				if glctx == nil || e.External {
					continue
				}

				sys.PaintEventChan <- paintDoneChan
				glimg.Upload()
				<-paintDoneChan

				glctx.ClearColor(0, 0, 0, 1)
				glctx.Clear(gl.COLOR_BUFFER_BIT)

				glimg.Draw(sz, geom.Point{0, 0}, geom.Point{geom.Pt(sz.WidthPx) / geom.Pt(sz.PixelsPerPt), 0}, geom.Point{0, geom.Pt(sz.HeightPx) / geom.Pt(sz.PixelsPerPt)}, glimg.RGBA.Bounds())

				a.Publish()
				a.Send(paint.Event{})
			case size.Event:
				sz = e
				select {
				case sys.InputEventChan <- e:
				default:
				}
			case touch.Event, key.Event:
				select {
				case sys.InputEventChan <- e:
				default:
				}
			}
		}
	})
}

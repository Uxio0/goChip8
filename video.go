package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

/*func main() {
	s := SDLWindow{}
	s.Init()
}*/

type SDLWindow struct {
	window  *sdl.Window
	surface *sdl.Surface
}

func (this *SDLWindow) Init() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	window.UpdateSurface()

	/*
		rect := sdl.Rect{0, 0, 200, 200}
		surface.FillRect(&rect, 0xffff0000)
		window.UpdateSurface()
		rect = sdl.Rect{0, 0, 100, 200}
		surface.FillRect(&rect, 0xff00ff00)
		window.UpdateSurface()
	*/
	this.window = window
	this.surface = surface

	sdl.Delay(1000)
}

func (this *SDLWindow) Close() {
	sdl.Quit()
	this.window.Destroy()
}

func (this *SDLWindow) Draw(sprites [64][32]bool) {
	for x, column := range sprites {
		for y, element := range column {
			//rect := sdl.Rect{int32(x), int32(y), int32(x), int32(y)}
			//this.surface.FillRect(&rect, 0xffffffff)
			var color byte
			if element {
				color = 0xff
			} else {
				color = 0
			}
			for i := x * 4; i < x*4+4; i++ {
				for j := y * 4; j < y*4+4; j++ {
					this.drawPoint(int32(i), int32(j), color)
				}
			}
		}
	}
	this.window.UpdateSurface()
}

func (this *SDLWindow) drawPoint(x int32, y int32, value byte) {
	pixels := this.surface.Pixels()
	address := ((y * this.surface.W) + x) * int32(this.surface.Format.BytesPerPixel)
	pixels[address] = value
	pixels[address+1] = value
	pixels[address+2] = value
	pixels[address+3] = value
}

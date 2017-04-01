package main

import "github.com/veandco/go-sdl2/sdl"

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

func Draw(sprites [][]byte) {

}

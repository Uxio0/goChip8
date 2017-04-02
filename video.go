package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

/*func main() {
	s := SDLWindow{}
	s.Init()
}*/

const (
	Scale  = 8
	Width  = 64
	Height = 32
)

type SDLWindow struct {
	window  *sdl.Window
	surface *sdl.Surface
}

func (this *SDLWindow) ManageEvents() {
	for event := sdl.PollEvent(); ; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.KeyUpEvent:
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		case *sdl.KeyDownEvent:
			fmt.Printf("[%d ms] Keyboard\ttype:%d\tsym:%c\tmodifiers:%d\tstate:%d\trepeat:%d\n",
				t.Timestamp, t.Type, t.Keysym.Sym, t.Keysym.Mod, t.State, t.Repeat)
		}
	}
}

func (this *SDLWindow) CheckKeyPress(key byte) bool {
	for i := 0; i < 15; i++ {
		event := sdl.PollEvent()
		switch t := event.(type) {
		case *sdl.KeyDownEvent:
			if string(t.Keysym.Sym) == fmt.Sprintf("%x", key) {
				return true
			}
		}
	}
	return false
}

func (this *SDLWindow) WaitUntilKeyPress() byte {
	for event := sdl.PollEvent(); ; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.KeyDownEvent:
			pressed, err := strconv.ParseInt(string(t.Keysym.Sym), 16, 8)
			if err == nil && pressed >= 0 && pressed < 16 {
				return byte(pressed)
			}
		}
	}
}

func (this *SDLWindow) Clear() {
	rect := sdl.Rect{0, 0, Scale * Width, Scale * Height}
	this.surface.FillRect(&rect, 0)
	this.window.UpdateSurface()
}

func (this *SDLWindow) Init() {
	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		Scale*Width, Scale*Height, sdl.WINDOW_SHOWN)
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

func (this *SDLWindow) Draw(sprites [Width][Height]bool) {
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
			for i := x * Scale; i < x*Scale+Scale; i++ {
				for j := y * Scale; j < y*Scale+Scale; j++ {
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
	pixels[address] = value   //Blue
	pixels[address+1] = value //Gren
	pixels[address+2] = value //Red
	//pixels[address+3] = value //Alpha?
}

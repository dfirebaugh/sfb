package sfb

/*
#cgo LDFLAGS: -L. -lsfb -lSDL2
#include "sfb.h"
#include <stdlib.h>

void set_pixel(SFB_Context *context, int x, int y, Pixel color) {
    sfb_set_pixel(context, x, y, color);
}

void set_window_title(SFB_Context *context, const char *title) {
    sfb_set_window_title(context, title);
}

*/
import "C"

import (
	"fmt"
	"image/color"
	"unsafe"
)

type displayer interface {
	SetPixel(x, y int16, c color.RGBA)
	Display() error
	Size() (int16, int16)
}

type Pixel struct {
	R, G, B, A uint8
}

const (
	FB_QUIT       = C.SDL_QUIT
	FB_KEYDOWN    = C.SDL_KEYDOWN
	FB_KEY_ESCAPE = C.SDLK_ESCAPE
)

type Screen interface {
	displayer
	SetWidth(width int16)
	SetHeight(height int16)
	SetTitle(title string)
	EnableFPS()
}

type display struct {
	context          *C.SFB_Context
	width            int16
	height           int16
	title            string
	framebuffer      []uint32
	scaleFactor      int
	fpsEnabled       bool
	transparentColor Pixel
}

func (d *display) SetWidth(width int16) {
	if width <= 0 {
		panic(fmt.Sprintf("invalid width: %d", width))
	}
	d.width = width
	d.updateFramebuffer()
}

func (d *display) SetHeight(height int16) {
	if height <= 0 {
		panic(fmt.Sprintf("invalid height: %d", height))
	}
	d.height = height
	d.updateFramebuffer()
}

func (d *display) SetTitle(title string) {
	d.title = title
	titleC := C.CString(title)
	defer C.free(unsafe.Pointer(titleC))
	C.set_window_title(d.context, titleC)
}

func (d *display) EnableFPS() {
	d.fpsEnabled = true
}

func (d *display) updateFramebuffer() {
	if d.width > 0 && d.height > 0 {
		d.framebuffer = make([]uint32, d.width*d.height)
		titleC := C.CString(d.title)
		defer C.free(unsafe.Pointer(titleC))
		if d.context == nil {
			d.context = (*C.SFB_Context)(C.malloc(C.sizeof_SFB_Context))
		}
		if C.sfb_init(d.context, titleC, 0, 0, C.int(d.width), C.int(d.height), 0) != 0 {
			C.free(unsafe.Pointer(d.context))
			d.context = nil
			panic("failed to initialize display")
		}
	}
}

func (d *display) Clear() {
	clearColor := Pixel{R: 0, G: 0, B: 0, A: 255}
	for i := range d.framebuffer {
		d.framebuffer[i] = uint32(clearColor.R)<<24 | uint32(clearColor.G)<<16 | uint32(clearColor.B)<<8 | uint32(clearColor.A)
	}
}

func (d *display) Destroy() {
	if d.context != nil {
		C.sfb_destroy(d.context)
		C.free(unsafe.Pointer(d.context))
		d.context = nil
	}
}

func (d *display) Display() error {
	if len(d.framebuffer) == 0 {
		d.SetWidth(240)
		d.SetHeight(120)
	}
	C.sfb_render(d.context)
	if d.fpsEnabled {
		C.sfb_enable_fps(d.context)
	}
	return nil
}

func (d *display) Height() int16 {
	return d.height
}

func (d *display) Resize(width, height int16) {
	if width <= 0 || height <= 0 {
		panic(fmt.Sprintf("invalid dimensions: width=%d, height=%d", width, height))
	}
	d.width = width
	d.height = height
	d.framebuffer = make([]uint32, width*height)
	C.sfb_resize(d.context, C.int(width), C.int(height))
}

func (d *display) SetPixel(x, y int16, c color.RGBA) {
	if x < 0 || x >= d.width || y < 0 || y >= d.height {
		return
	}

	color := Pixel{R: c.R, G: c.G, B: c.B, A: c.A}
	C.set_pixel(d.context, C.int(x), C.int(y), *(*C.Pixel)(unsafe.Pointer(&color)))
}

func (d *display) Size() (int16, int16) {
	return d.width, d.height
}

func (d *display) Width() int16 {
	return d.width
}

func Run(update func(s Screen)) {
	display := &display{
		scaleFactor:      1,
		fpsEnabled:       false,
		transparentColor: Pixel{R: 0, G: 0, B: 0, A: 0},
	}

	update(display)
	display.Destroy()
}

func PollEvent() (eventType, keySym uint32, eventAvailable bool) {
	var cEventType C.uint32_t
	var cKeySym C.uint32_t
	if C.sfb_poll_event(&cEventType, &cKeySym) != 0 {
		return uint32(cEventType), uint32(cKeySym), true
	}
	return 0, 0, false
}

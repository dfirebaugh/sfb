package main

import (
	"runtime"

	"github.com/dfirebaugh/sfb"

	"github.com/dfirebaugh/hlg/pkg/draw"
	"github.com/dfirebaugh/hlg/pkg/math/geom"
	"golang.org/x/image/colornames"
)

func handleEvent(eventType uint32, keySym uint32) bool {
	switch eventType {
	case sfb.FB_QUIT:
		return false
	case sfb.FB_KEYDOWN:
		if keySym == sfb.FB_KEY_ESCAPE {
			return false
		}
	}
	return true
}

func main() {
	runtime.LockOSThread()
	width, height := 240, 120
	title := "Simple Example"

	tri := geom.MakeTriangle([3]geom.Vector{
		geom.MakeVector(120, 30),
		geom.MakeVector(90, 90),
		geom.MakeVector(150, 90),
	})

	sfb.Run(func(s sfb.Screen) {
		s.SetWidth(int16(width))
		s.SetHeight(int16(height))
		s.SetTitle(title)
		running := true

		for running {
			for {
				eventType, keySym, eventAvailable := sfb.PollEvent()
				if !eventAvailable {
					break
				}
				if !handleEvent(eventType, keySym) {
					running = false
					break
				}
			}

			draw.Triangle(tri).Fill(s, colornames.Purple)
			s.Display()
		}
	})
}

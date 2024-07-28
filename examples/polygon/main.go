package main

import (
	"github.com/dfirebaugh/hlg/pkg/draw"
	"github.com/dfirebaugh/hlg/pkg/math/geom"
	"github.com/dfirebaugh/sfb"
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
	width, height := 240, 120
	title := "Polygon Example"

	points := []geom.Point{
		geom.MakePoint(60, 40),
		geom.MakePoint(120, 40),
		geom.MakePoint(150, 60),
		geom.MakePoint(120, 80),
		geom.MakePoint(60, 80),
		geom.MakePoint(30, 60),
	}

	polygon := geom.MakePolygon(points...)
	background := geom.MakeRect(0, 0, float32(width), float32(height))

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

			draw.Rect(background).Fill(s, colornames.Black)
			draw.Polygon(polygon).Fill(s, colornames.Green)
			s.Display()
		}
	})
}

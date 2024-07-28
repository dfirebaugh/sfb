package draw

import (
	"image/color"

	"github.com/dfirebaugh/sfb/pkg/geom"
	"golang.org/x/image/colornames"
)

type Line geom.Line

func (l Line) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	l.line(d, int16(l.Segment.V0[0]), int16(l.Segment.V0[1]), int16(l.Segment.V1[0]), int16(l.Segment.V1[1]), color)
}

func (l Line) line(display displayer, x0 int16, y0 int16, x1 int16, y1 int16, color color.RGBA) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := int16(1)
	sy := int16(1)
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	for {
		display.SetPixel(x0, y0, color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

package draw

import (
	"image/color"
	"sync"

	"github.com/dfirebaugh/sfb/pkg/geom"
	"golang.org/x/image/colornames"
)

type Triangle geom.Triangle

func (t Triangle) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	t.draw(d, int16(t[0][0]), int16(t[0][1]), int16(t[1][0]), int16(t[1][1]), int16(t[2][0]), int16(t[2][1]), color)
}

func (t Triangle) draw(display displayer, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA) {
	Line{}.line(display, x0, y0, x1, y1, color)
	Line{}.line(display, x0, y0, x2, y2, color)
	Line{}.line(display, x1, y1, x2, y2, color)
}

func (t Triangle) Fill(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	t.fill(d, int16(t[0][0]), int16(t[0][1]), int16(t[1][0]), int16(t[1][1]), int16(t[2][0]), int16(t[2][1]), color)
}

func (t Triangle) fill(display displayer, x0 int16, y0 int16, x1 int16, y1 int16, x2 int16, y2 int16, color color.RGBA) {
	var wg sync.WaitGroup
	numWorkers := 4
	tasks := make(chan struct{ x1, y1, x2, y2 int16 }, numWorkers*2)

	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}
	if y1 > y2 {
		x1, y1, x2, y2 = x2, y2, x1, y1
	}
	if y0 > y1 {
		x0, y0, x1, y1 = x1, y1, x0, y0
	}

	if y0 == y2 {
		a := x0
		b := x0
		if x1 < a {
			a = x1
		} else if x1 > b {
			b = x1
		}
		if x2 < a {
			a = x2
		} else if x2 > b {
			b = x2
		}
		Line{}.line(display, a, y0, b, y0, color)
		return
	}

	dx01 := x1 - x0
	dy01 := y1 - y0
	dx02 := x2 - x0
	dy02 := y2 - y0
	dx12 := x2 - x1
	dy12 := y2 - y1

	sa := int16(0)
	sb := int16(0)
	a := int16(0)
	b := int16(0)

	last := y1 - 1
	if y1 == y2 {
		last = y1
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				Line{}.line(display, task.x1, task.y1, task.x2, task.y2, color)
			}
		}()
	}

	// Upper
	y := y0
	for ; y <= last; y++ {
		a = x0 + sa/dy01
		b = x0 + sb/dy02
		sa += dx01
		sb += dx02
		tasks <- struct{ x1, y1, x2, y2 int16 }{a, y, b, y}
	}

	// Lower
	sa = dx12 * (y - y1)
	sb = dx02 * (y - y0)
	for ; y <= y2; y++ {
		a = x1 + sa/dy12
		b = x0 + sb/dy02
		sa += dx12
		sb += dx02
		tasks <- struct{ x1, y1, x2, y2 int16 }{a, y, b, y}
	}

	close(tasks)
	wg.Wait()
}

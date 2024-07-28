package draw

import (
	"errors"
	"image/color"
	"sync"

	"github.com/dfirebaugh/sfb/pkg/geom"
	"golang.org/x/image/colornames"
)

var frameCounter uint64

type Rect geom.Rect

func (r Rect) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	r.draw(d, int16(r[0]), int16(r[1]), int16(r[2]), int16(r[3]), color)
}

func (r Rect) draw(display displayer, x int16, y int16, w int16, h int16, color color.RGBA) error {
	if w <= 0 || h <= 0 {
		return errors.New("empty rectangle")
	}
	var wg sync.WaitGroup

	numWorkers := 12
	tasks := make(chan struct{ x1, y1, x2, y2 int16 }, numWorkers*2)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				Line{}.line(display, task.x1, task.y1, task.x2, task.y2, color)
			}
		}()
	}

	tasksToDraw := []struct {
		x1, y1, x2, y2 int16
	}{
		{x, y, x + w - 1, y},
		{x, y, x, y + h - 1},
		{x + w - 1, y, x + w - 1, y + h - 1},
		{x, y + h - 1, x + w - 1, y + h - 1},
	}

	for _, task := range tasksToDraw {
		tasks <- task
	}

	close(tasks)
	wg.Wait()
	return nil
}

func (r Rect) Fill(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	r.fill(d, int16(r[0]), int16(r[1]), int16(r[2]), int16(r[3]), color)
}

func (r Rect) fill(display displayer, x int16, y int16, w int16, h int16, color color.RGBA) error {
	if w <= 0 || h <= 0 {
		return errors.New("empty rectangle")
	}
	var wg sync.WaitGroup

	numWorkers := 4
	tasks := make(chan struct{ startX, endX, startY, endY int16 }, numWorkers*2)

	// drawEven := atomic.LoadUint64(&frameCounter)%2 == 0

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range tasks {
				for x := task.startX; x <= task.endX; x++ {
					// if (x%2 == 0) != drawEven {
					// 	continue
					// }
					Line{}.line(display, x, task.startY, x, task.endY, color)
				}
			}
		}()
	}

	halfW, halfH := w/2, h/2
	quadrants := []struct{ startX, endX, startY, endY int16 }{
		{x, x + halfW - 1, y, y + halfH - 1},
		{x + halfW, x + w - 1, y, y + halfH - 1},
		{x, x + halfW - 1, y + halfH, y + h - 1},
		{x + halfW, x + w - 1, y + halfH, y + h - 1},
	}

	for _, quad := range quadrants {
		tasks <- quad
	}

	close(tasks)
	wg.Wait()
	return nil
}

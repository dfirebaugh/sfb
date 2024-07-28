package draw

import (
	"image/color"
	"sync"

	"github.com/dfirebaugh/sfb/pkg/geom"
	"golang.org/x/image/colornames"
)

type Circle geom.Circle

func (c Circle) Draw(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	c.draw(d, int16(c.X), int16(c.Y), int16(c.R), color)
}

func (c Circle) draw(display displayer, x0 int16, y0 int16, r int16, color color.RGBA) {
	var wg sync.WaitGroup

	setCirclePixels := func(x, y int16) {
		display.SetPixel(x0+x, y0+y, color)
		display.SetPixel(x0-x, y0+y, color)
		display.SetPixel(x0+x, y0-y, color)
		display.SetPixel(x0-x, y0-y, color)
		display.SetPixel(x0+y, y0+x, color)
		display.SetPixel(x0-y, y0+x, color)
		display.SetPixel(x0+y, y0-x, color)
		display.SetPixel(x0-y, y0-x, color)
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		display.SetPixel(x0, y0+r, color)
		display.SetPixel(x0, y0-r, color)
		display.SetPixel(x0+r, y0, color)
		display.SetPixel(x0-r, y0, color)

		f := 1 - r
		ddfx := int16(1)
		ddfy := -2 * r
		x := int16(0)
		y := r

		numWorkers := 4
		tasks := make(chan struct{ x, y int16 }, numWorkers)
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range tasks {
					if (task.x+task.y)%2 == 0 {
						setCirclePixels(task.x, task.y)
					}
				}
			}()
		}

		for x < y {
			if f >= 0 {
				y--
				ddfy += 2
				f += ddfy
			}
			x++
			ddfx += 2
			f += ddfx

			tasks <- struct{ x, y int16 }{x, y}
		}

		close(tasks)
	}()

	go func() {
		defer wg.Done()
		display.SetPixel(x0, y0+r, color)
		display.SetPixel(x0, y0-r, color)
		display.SetPixel(x0+r, y0, color)
		display.SetPixel(x0-r, y0, color)

		f := 1 - r
		ddfx := int16(1)
		ddfy := -2 * r
		x := int16(0)
		y := r

		numWorkers := 4
		tasks := make(chan struct{ x, y int16 }, numWorkers)
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range tasks {
					if (task.x+task.y)%2 != 0 {
						setCirclePixels(task.x, task.y)
					}
				}
			}()
		}

		for x < y {
			if f >= 0 {
				y--
				ddfy += 2
				f += ddfy
			}
			x++
			ddfx += 2
			f += ddfx

			tasks <- struct{ x, y int16 }{x, y}
		}

		close(tasks)
	}()

	wg.Wait()
}

func (c Circle) Fill(d displayer, clr color.Color) {
	color, ok := clr.(color.RGBA)
	if !ok {
		color = colornames.Black
	}
	c.fill(d, int16(c.X), int16(c.Y), int16(c.R), color)
}

func (c Circle) fill(display displayer, x0 int16, y0 int16, r int16, color color.RGBA) {
	var wg sync.WaitGroup

	f := 1 - r
	ddfx := int16(1)
	ddfy := -2 * r
	x := int16(0)
	y := r

	numWorkers := 4
	oddTasks := make(chan struct{ x, y int16 }, numWorkers)
	evenTasks := make(chan struct{ x, y int16 }, numWorkers)

	// Worker pool for odd pixels
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range oddTasks {
				Line{}.line(display, x0+task.x, y0-task.y, x0+task.x, y0+task.y, color)
				Line{}.line(display, x0+task.y, y0-task.x, x0+task.y, y0+task.x, color)
				Line{}.line(display, x0-task.x, y0-task.y, x0-task.x, y0+task.y, color)
				Line{}.line(display, x0-task.y, y0-task.x, x0-task.y, y0+task.x, color)
			}
		}()
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range evenTasks {
				Line{}.line(display, x0+task.x, y0-task.y, x0+task.x, y0+task.y, color)
				Line{}.line(display, x0+task.y, y0-task.x, x0+task.y, y0+task.x, color)
				Line{}.line(display, x0-task.x, y0-task.y, x0-task.x, y0+task.y, color)
				Line{}.line(display, x0-task.y, y0-task.x, x0-task.y, y0+task.x, color)
			}
		}()
	}

	Line{}.line(display, x0, y0-r, x0, y0+r, color)

	for x < y {
		if f >= 0 {
			y--
			ddfy += 2
			f += ddfy
		}
		x++
		ddfx += 2
		f += ddfx

		task := struct{ x, y int16 }{x, y}
		if (x+y)%2 == 0 {
			evenTasks <- task
		} else {
			oddTasks <- task
		}
	}

	close(oddTasks)
	close(evenTasks)
	wg.Wait()
}

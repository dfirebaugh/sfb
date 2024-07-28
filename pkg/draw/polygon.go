package draw

import (
	"image/color"
	"sort"
	"sync"

	"github.com/dfirebaugh/sfb/pkg/geom"
)

type Polygon geom.Polygon

func (p Polygon) Draw(d displayer, clr color.Color) {
	numPoints := len(p)
	for i := 0; i < numPoints; i++ {
		start := p[i]
		end := p[(i+1)%numPoints]

		line := Line(geom.MakeLine(start, end))
		line.Draw(d, clr)
	}
}

func (p Polygon) Fill(d displayer, clr color.Color) {
	minY, maxY := p[0].Y, p[0].Y
	for _, pt := range p {
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}

	var wg sync.WaitGroup
	numWorkers := 4
	tasks := make(chan int, numWorkers*2)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for y := range tasks {
				var intersections []float64
				for i := 0; i < len(p); i++ {
					next := (i + 1) % len(p)
					if intersects(p[i], p[next], float32(y)) {
						x := float64(intersectionX(p[i], p[next], float32(y)))
						intersections = append(intersections, x)
					}
				}

				sort.Float64s(intersections)

				for i := 0; i < len(intersections)-1; i += 2 {
					start := geom.MakePoint(float32(intersections[i]), float32(y))
					end := geom.MakePoint(float32(intersections[i+1]), float32(y))
					line := Line(geom.MakeLine(start, end))
					line.Draw(d, clr)
				}
			}
		}()
	}

	for y := int(minY); y <= int(maxY); y++ {
		tasks <- y
	}

	close(tasks)
	wg.Wait()
}

func intersects(p1, p2 geom.Point, y float32) bool {
	return (p1.Y <= y && p2.Y > y) || (p1.Y > y && p2.Y <= y)
}

func intersectionX(p1, p2 geom.Point, y float32) float32 {
	if p1.Y == p2.Y {
		return p1.X
	}
	return p1.X + (y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
}

package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"sync"

	"github.com/dfirebaugh/sfb"
	"github.com/dfirebaugh/sfb/pkg/draw"
	"github.com/dfirebaugh/sfb/pkg/geom"
	"golang.org/x/image/colornames"
)

const (
	numCircles = 45
)

type MovingCircle struct {
	geom.Circle
	Velocity geom.Point
	Color    color.RGBA
}

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

func (c *MovingCircle) Update(screenWidth, screenHeight int, circles []MovingCircle, index int) {
	c.X += c.Velocity.X
	c.Y += c.Velocity.Y

	if c.X-c.R <= 0 || c.X+c.R >= float32(screenWidth) {
		c.Velocity.X = -c.Velocity.X
	}
	if c.Y-c.R <= 0 || c.Y+c.R >= float32(screenHeight) {
		c.Velocity.Y = -c.Velocity.Y
	}

	for i := range circles {
		if i != index {
			dx := circles[i].X - c.X
			dy := circles[i].Y - c.Y
			distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			if distance < c.R+circles[i].R {
				overlap := 0.5 * (c.R + circles[i].R - distance)

				c.X -= overlap * (dx / distance)
				c.Y -= overlap * (dy / distance)
				circles[i].X += overlap * (dx / distance)
				circles[i].Y += overlap * (dy / distance)

				normalX, normalY := dx/distance, dy/distance
				tangentX, tangentY := -normalY, normalX

				dotProductTangent1 := c.Velocity.X*tangentX + c.Velocity.Y*tangentY
				dotProductTangent2 := circles[i].Velocity.X*tangentX + circles[i].Velocity.Y*tangentY

				dotProductNormal1 := c.Velocity.X*normalX + c.Velocity.Y*normalY
				dotProductNormal2 := circles[i].Velocity.X*normalX + circles[i].Velocity.Y*normalY

				momentum1 := (dotProductNormal1*(c.R-circles[i].R) + 2*circles[i].R*dotProductNormal2) / (c.R + circles[i].R)
				momentum2 := (dotProductNormal2*(circles[i].R-c.R) + 2*c.R*dotProductNormal1) / (c.R + circles[i].R)

				c.Velocity.X = tangentX*dotProductTangent1 + normalX*momentum1
				c.Velocity.Y = tangentY*dotProductTangent1 + normalY*momentum1
				circles[i].Velocity.X = tangentX*dotProductTangent2 + normalX*momentum2
				circles[i].Velocity.Y = tangentY*dotProductTangent2 + normalY*momentum2
			}
		}
	}
}

func NewMovingCircle(screenWidth, screenHeight int) MovingCircle {
	radius := 20 + rand.Float32()*30
	x := radius + rand.Float32()*(float32(screenWidth)-2*radius)
	y := radius + rand.Float32()*(float32(screenHeight)-2*radius)

	return MovingCircle{
		Circle: geom.Circle{X: x, Y: y, R: radius},
		Velocity: geom.Point{
			X: rand.Float32()*4 - 2,
			Y: rand.Float32()*4 - 2,
		},
		Color: color.RGBA{
			R: uint8(rand.Intn(256)),
			G: uint8(rand.Intn(256)),
			B: uint8(rand.Intn(256)),
			A: 255,
		},
	}
}

func main() {
	runtime.LockOSThread()

	screenWidth, screenHeight := 800, 600
	circles := make([]MovingCircle, numCircles)
	for i := range circles {
		circles[i] = NewMovingCircle(screenWidth, screenHeight)
	}

	background := geom.MakeRect(0, 0, float32(screenWidth), float32(screenHeight))
	sfb.Run(func(s sfb.Screen) {
		s.SetWidth(int16(screenWidth))
		s.SetHeight(int16(screenHeight))
		s.SetTitle("Moving Circles")
		s.EnableFPS()

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
			var wg sync.WaitGroup
			for i := range circles {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					draw.Circle(circles[i].Circle).Fill(s, circles[i].Color)
					circles[i].Update(screenWidth, screenHeight, circles, i)
				}(i)
			}

			wg.Wait()

			if err := s.Display(); err != nil {
				fmt.Printf("Error displaying framebuffer: %v\n", err)
				running = false
			}
		}
	})
}

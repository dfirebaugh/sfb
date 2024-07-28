package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/dfirebaugh/sfb"
	"github.com/dfirebaugh/sfb/pkg/draw"
	"github.com/dfirebaugh/sfb/pkg/geom"
	"golang.org/x/image/colornames"
)

const (
	numCircles = 50
	rectWidth  = 50
	rectHeight = 50
)

type MovingCircle struct {
	circle geom.Circle
	dx     float32
	dy     float32
}

type MovingRect struct {
	rect geom.Rect
	dx   float32
	dy   float32
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

func main() {
	runtime.LockOSThread()

	screenWidth, screenHeight := 800, 600

	circles := make([]MovingCircle, numCircles)
	for i := range circles {
		circles[i] = MovingCircle{
			circle: geom.MakeCircle(
				rand.Float32()*float32(screenWidth),
				rand.Float32()*float32(screenHeight),
				5+rand.Float32()*10),
			dx: rand.Float32()*4 - 2,
			dy: rand.Float32()*4 - 2,
		}
	}

	rect := MovingRect{
		rect: geom.MakeRect(
			rand.Float32()*float32(screenWidth-rectWidth),
			rand.Float32()*float32(screenHeight-rectHeight),
			rectWidth,
			rectHeight,
		),
		dx: rand.Float32()*10 - 2,
		dy: rand.Float32()*10 - 2,
	}

	sfb.Run(func(s sfb.Screen) {
		s.SetWidth(int16(screenWidth))
		s.SetHeight(int16(screenHeight))
		s.SetTitle("Moving and Rect")
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

			draw.Rect(rect.rect).Fill(s, colornames.Black)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				draw.Rect(rect.rect).Fill(s, colornames.Red)
			}()

			rect.rect[0] += rect.dx
			rect.rect[1] += rect.dy

			if rect.rect[0] <= 0 || rect.rect[0]+rect.rect[2] >= float32(screenWidth) {
				rect.dx = -rect.dx
			}
			if rect.rect[1] <= 0 || rect.rect[1]+rect.rect[3] >= float32(screenHeight) {
				rect.dy = -rect.dy
			}

			wg.Wait()

			if err := s.Display(); err != nil {
				fmt.Printf("Error displaying framebuffer: %v\n", err)
				running = false
			}
			time.Sleep(3 * time.Millisecond)
		}
	})
}

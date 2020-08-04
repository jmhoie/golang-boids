package main

import (
	"image/color"
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth, screenHeight = 640, 360
	boidCount                 = 500   // amount of boids we're creating
	viewRadius                = 13    // amount of pixels each boid can see
	adjRate                   = 0.015 // adjustment rate for the acceleration, this smoothens out the boid movement
)

var (
	green   = color.RGBA{10, 255, 50, 255} // green color used for rendering the boids
	boids   [boidCount]*Boid
	boidMap [screenWidth + 1][screenHeight + 1]int // 2D array representing the position of each boid - this is memory shared between all boids. This will get upadated with the boids id, to represent that the spot is taken up.
	rWlock  = sync.RWMutex{}                       // Readers-Writers mutex used for thread synchronization
)

func update(screen *ebiten.Image) error {
	if !ebiten.IsDrawingSkipped() {
		for _, boid := range boids { // renders 4 pixels for each boid in a diamond shape
			screen.Set(int(boid.position.x+1), int(boid.position.y), green)
			screen.Set(int(boid.position.x-1), int(boid.position.y), green)
			screen.Set(int(boid.position.x), int(boid.position.y-1), green)
			screen.Set(int(boid.position.x), int(boid.position.y+1), green)
		}
	}
	return nil
}

func main() {
	for i, row := range boidMap { // initialize the boidMap with -1 (to represent an empty spot)
		for j := range row {
			boidMap[i][j] = -1
		}
	}

	for i := 0; i < boidCount; i++ {
		createBoid(i)
	}
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Boids in a box"); err != nil {
		log.Fatal(err)
	}
}

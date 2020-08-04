package main

import (
	"math"
	"math/rand"
	"time"
)

// Boid is the structure of each boid (bird-android). It has a position vector, a velocity vector and an id.
type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

func (b *Boid) calcAcceleration() Vector2D { // find the enclosing window (the view space) for the boid, i.e. all the positions that a given boid can see
	upper, lower := b.position.AddV(viewRadius), b.position.AddV(-viewRadius) // upper right and lower left of the view box
	avgPosition, avgVelocity, separation := Vector2D{0, 0}, Vector2D{0, 0}, Vector2D{0, 0}
	count := 0.0 // count of boids closer than the viewradius, i.e. all the other boids inside the view box

	rWlock.RLock() // aquire the readers lock, which does NOT restrict reading access to other boids.
	for i := math.Max(lower.x, 0); i <= math.Min(upper.x, screenWidth); i++ {
		for j := math.Max(lower.y, 0); j <= math.Min(upper.y, screenHeight); j++ {
			if otherBoidID := boidMap[int(i)][int(j)]; otherBoidID != -1 && otherBoidID != b.id { // finds all the other boids in the view box, also avoids checking itself.
				if dist := boids[otherBoidID].position.Distance(b.position); dist < viewRadius {
					count++                                                                                       // update the number of boids encountered
					avgVelocity = avgVelocity.Add(boids[otherBoidID].velocity)                                    // update the average velocity of all boids in the view space
					avgPosition = avgPosition.Add(boids[otherBoidID].position)                                    // update the average position of all boids in the view space
					separation = separation.Add(b.position.Subtract(boids[otherBoidID].position).DivisionV(dist)) // gives an average acceleration to move away from other boids - i.e. the separation acceleration
				}
			}
		}
	}
	rWlock.RUnlock() // unlock the readers lock

	accel := Vector2D{ // calculates how far we should bounce, depending on how far we are from the edge/wall
		b.borderBounce(b.position.x, screenWidth),
		b.borderBounce(b.position.y, screenHeight),
	}
	if count > 0 {
		avgPosition, avgVelocity = avgPosition.DivisionV(count), avgVelocity.DivisionV(count)
		accelAlignment := avgVelocity.Subtract(b.velocity).MultiplyV(adjRate)
		accelCohesion := avgPosition.Subtract(b.position).MultiplyV(adjRate)
		accelSeparation := separation.MultiplyV(adjRate)
		accel = accel.Add(accelAlignment).Add(accelCohesion).Add(accelSeparation) // total acceleration
	}

	return accel
}

// borderBounce calculates the acceleration which a boid will need to turn away when approaching a wall
func (b *Boid) borderBounce(pos, maxBorderPos float64) float64 {
	if pos < viewRadius {
		return 1 / pos
	} else if pos > maxBorderPos-viewRadius {
		return 1 / (pos - maxBorderPos)
	}
	return 0
}

// moveOne is the method to move the boid once. It adds the velocity vector to the position vector.
func (b *Boid) moveOne() {
	acceleration := b.calcAcceleration()                   // since calcAcceleration also locks the mutex, this must be done before locking the mutex again below. If this isn't done, then the boid won't be able to extract the value, sinc the mutex is already locked.
	rWlock.Lock()                                          // aquire and lock the mutex - this is the writers lock, which locks the whole mutex, so neither reading nor writing can be accessed by other boids
	b.velocity = b.velocity.Add(acceleration).limit(-1, 1) // updates the velocity based on the acceleration, also limits it so it doesn't jump more than 1 pixel at a time (helps gettin smooth movement).
	boidMap[int(b.position.x)][int(b.position.y)] = -1     // when a boid is going to move, the old position in the boidMap is updated with -1
	b.position = b.position.Add(b.velocity)
	boidMap[int(b.position.x)][int(b.position.y)] = b.id // when a boid has moved, the new position in the boidMap is updated with the boid id

	rWlock.Unlock() // unlock the mutex, releasing it for other boids to use - both reading and writing.
}

// start is the method to start the boids movement. It's an infinite loop which moves the boid one space then waits for 5ms, and then repeats.
func (b *Boid) start() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

// createBoid is the constructor function for new a new boid.
func createBoid(bid int) {
	b := Boid{
		position: Vector2D{rand.Float64() * screenWidth, rand.Float64() * screenHeight}, // rand.Float64() gives a random number between 0 and 1, multiply by height and width to get a random position on de screen
		velocity: Vector2D{(rand.Float64() * 2) - 1.0, (rand.Float64() * 2) - 1.0},      // want less than 1 pixel, to make the motion fluid (boids jumping more than one pixel would look weird). We get a random number on the interval [-1, 1]
		id:       bid,
	}
	boids[bid] = &b
	boidMap[int(b.position.x)][int(b.position.y)] = b.id // updates the boidMap to with the boid id - that spot is now occupied by a boid
	go b.start()                                         // "go" is the keyword to start a new thread. So this function creates a boid and starts it's movement in it's own thread.
}

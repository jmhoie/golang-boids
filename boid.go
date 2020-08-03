package main

// Boid is the structure of each boid. It has a position vector, a velocity vector and an id.
type Boid struct {
	position Vector2D
	velocity Vector2D
	id       int
}

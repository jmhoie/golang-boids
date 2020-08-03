package main

import "math"

// Vector2D represents a vector with an x- and y-coordinate. This is used for both position and velocity.
type Vector2D struct {
	x float64
	y float64
}

// Add adds two 2D vectors.
func (v1 Vector2D) Add(v2 Vector2D) Vector2D {
	return Vector2D{v1.x + v2.x, v1.y + v2.y}
}

// Subtract subtracts two 2D vectors.
func (v1 Vector2D) Subtract(v2 Vector2D) Vector2D {
	return Vector2D{v1.x - v2.x, v1.y - v2.y}
}

// Multiply multiplies two 2D vectors.
func (v1 Vector2D) Multiply(v2 Vector2D) Vector2D {
	return Vector2D{v1.x * v2.x, v1.y * v2.y}
}

// AddV adds a value to a 2D vector.
func (v1 Vector2D) AddV(d float64) Vector2D {
	return Vector2D{v1.x + d, v1.y + d}
}

// MultiplyV multiplies a 2D vector with a value.
func (v1 Vector2D) MultiplyV(d float64) Vector2D {
	return Vector2D{v1.x * d, v1.y * d}
}

// DivisionV divides a value from a 2D vector.
func (v1 Vector2D) DivisionV(d float64) Vector2D {
	return Vector2D{v1.x / d, v1.y / d}
}

// limit restricts x- and y-position to a lower and upper bound.
func (v1 Vector2D) limit(lower, upper float64) Vector2D {
	return Vector2D{math.Min(math.Max(v1.x, lower), upper), math.Min(math.Max(v1.y, lower), upper)}
}

// Distance calculates the distance between two vectors using Pythagoras' Theorem
func (v1 Vector2D) Distance(v2 Vector2D) float64 {
	return math.Sqrt(math.Pow(v1.x-v2.x, 2) + math.Pow(v1.y-v2.y, 2))
}

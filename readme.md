# Boids Simulation  
***This program was made by following the Udemy course [*Mastering Multithreaded Programming with Go (Golang)*](https://www.udemy.com/course/multithreading-in-go-lang/) by James Cutajar.***  

In this project we will be using GoRoutines to model Boids. Using concurrent programming to model a problem can make things significantly simpler.

Instead of having a complicated simulation with alot of parameters, we will be using threads where each thread can interract with other threads based on defined properties and ruleset.

***Goal***: Simulate the flocking behaviour of birds in large groups.

## Implementation
We'll be plotting dots in a 2D coordinate system, where each dot will represent a boid.  
For graphics, we're using the 2D game engine [*Ebiten.*](https://github.com/hajimehoshi/ebiten)

A boid has 3 properties; position vector, velocity vector and an id.  
To simulate the movement and behaviour of a flock of birds, we'll be implementing 3 simple rulesets; *Alignment*, *Cohesion* and *Separation*.

In order to implement these rulesets, the boids need to be able to get the position and velocity of other nearby boids. We will accomplish this with memory sharing by mapping all the possible positions (the whole graphics window) in a 2D array. Then we define a *view-box* as a smaller part of the position map which surrounds a given boid. I.e. surrounding each boid, there is a square which represents how far that boid can see. An individual boid will be able to read other boids properties as long as they're within the given boids view. This is much more efficient than checking every boid on the screen and calculating the distance (the latter being an O(n^2) formula).

### Alignment
The *alignment* will make the boids form groups by averaging out the current velocities of all the boids in a given view-space. I.e. boids close to eachother will move in the same direction.  
This is accomplished by taking the target vector (where a boid is trying to move) and subtracting it by the current velocity vector. This results in an acceleration vector, which we then add to the original velocity vector. To smoothen out this movement, we also multiply that acceleration vector by an adjustment rate (a float between 0 and 1). This will slow down the movement which makes it less jittery.  
$target(0,4) - current(3,3) = acceleration(-3,1)$


This is accomplished by taking the target vector (where it a boid is trying to move) and subtracting it by the current velocity vector. This results in an acceleration vector, which we then add to the original velocity vector. To smoothen out the movement, we also multiply the acceleration vector by an ajudstment rate (number between 0 and 1). This slows down the movement which makes it less "jittery".  
Example: target(0,4) - current(3,3) = acceleration(-3,1) -> acceleration(-3,1)*adjRate(0.5) + current(3,3) = new(1.5,3.5)

3 rules: Alignment, Cohesion, Separation

Cohesion: Moving a boid to a closer position to the other boids (grouping the together). Do this by finding the average position of nearby boids, and calculate the vector which will bring it there. The new positiong is calculated by: target-positon - current-position = acceleration (which we multiply by the adjustment rate). Then we add the acceleration this acceleration to the alignment acceleration and the current velocity.
I.e: cohesion-acceleration + alignment-acceleration + current-velocity = new-velocity.
    
Race condition -> uncysnced threads results in inaccuracies because calculations aren't restricted to a set of boids at a time. In this program, when scanning the view-box for nearby boids, we do it one position at a time (in a 2D array). However it is possible for a boid to move to a new position before we have completed the scanning, which can result in a specific boid getting counted multiple times or alternatively not counted at all. This gives very inaccurate calculations for the average flock velocity. Since boids change their velocities, it is also possible for the scan to get the velocity of a boid right before it changes, which also results in inaccurate calculations.

Thread synchronization using mutexes. Locks a particular piece of execution. Gives and restrict access. Guarrantees that only one thread can update a value at a time. A mutex can only be held by a single thread at a time. In this program we will use a mutex when: reading or updating the map, updating the velocity, or updating the position.

Readers-Writers lock: Special mutext that allows multiple readers but ONLY one writer. If the writer-lock is NOT in use, then multiple threads can read from the reader lock. When the writer lock IS in use, this also blocks the reader-lock so that it cannot be accessed while values are being updated. This increases efficiency, especially in this program, since the write operations are much faster than the read operations (read -> loops through the whole 2D position array. write -> updates position, velocity and map).

Separation: Move away from boids that are too close. Take the current position vector and for each nearby boid subtract it's position. This results in an acceleration which will move the current boid away from the group. We divide this by the distance to the boid, so that if we're closer the acceleration will be higher than if we're further apart. Do this for each boid in the view-space and add the resulting accelerations together to get the new vector. This new vector is the separation acceleration, which we then multiply by the adjustment rate. Then we add all the accelerations together (alignment, cohesion and separation) to get the total acceleration. By adding this to the current velocity we get the new velocity for the boid.

Updating wall-bounce: Steer away from wall (horizontal or vertical) instead of just flipping around when it hits the wall. This is much more realistic. We calculate the acceleration to move away from the wall by getting the distance between the boid and the wall (both x- and y-direction), then calculate the bounce acceleration. This acceleration is equal to the reciprical of the x- and y-distance (1/distance_x, 1/distance_y). This bounce acceleration will be larger, the closer we are to the wall. We add this to the other accelerations and the current velocity. This is the final new velocity.

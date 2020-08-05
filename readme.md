# Boids Simulation <!-- omit in toc -->
***This simulation was made by following the Udemy course [*Mastering Multithreaded Programming with Go (Golang)*](https://www.udemy.com/course/multithreading-in-go-lang/) by James Cutajar.***  

For instructions on how to run this program, see [*Running The Simulation*](#2-running-the-simulation)

In this project we will be using GoRoutines to model Boids. Using concurrent programming to model a problem can make things significantly simpler.

Instead of having a complicated simulation with alot of parameters, we will be using threads where each thread can interract with other threads based on defined properties and ruleset.

***Goal***: Simulate the flocking behaviour of birds in large groups.

- [1. Implementation](#1-implementation)
  - [1.1. Thread Synchronization](#11-thread-synchronization)
  - [1.2. Alignment](#12-alignment)
  - [1.3. Cohesion](#13-cohesion)
  - [1.4. Separation](#14-separation)
  - [1.5. Wall-Bounce](#15-wall-bounce)
- [2. Running The Simulation](#2-running-the-simulation)

## 1. Implementation
We'll be plotting dots in a 2D coordinate system, where each dot will represent a boid.  
For graphics, we're using the 2D game engine [*Ebiten.*](https://github.com/hajimehoshi/ebiten)

A boid has 3 properties; position vector, velocity vector and an id.  
To simulate the movement and behaviour of a flock of birds, we'll be implementing 3 simple rulesets; *Alignment*, *Cohesion* and *Separation*.

In order to implement these rulesets, the boids need to be able to get the position and velocity of other nearby boids. We will accomplish this with memory sharing by mapping all the possible positions (the whole graphics window) in a 2D array. Then we define a *view-box* as a smaller part of the position map which surrounds a given boid. I.e. surrounding each boid, there is a square which represents how far that boid can see. An individual boid will be able to read other boids properties as long as they're within the given boids view. This is much more efficient than checking every boid on the screen and calculating the distance (the latter being an O(n^2) formula).

### 1.1. Thread Synchronization
In this program, there is possible to encounter a race condition; Unsynced threads can result in inaccuracies because calculations are not restricted to one boid at a time. Here we're scanning the view-box for nearby boids one position at atime (in the 2D position array). However boids will keep moving during this time, which may result in a single boid being counted multiple, or alternatively, not counted at all. This will result in inaccurate calculations. It is also possible that a boid changes it velocity right after it is scanned, which also result in inaccuracies.

To avoid this, we're employing thread synchronization using mutexes. A mutex will manage the access to a particular part of execution. It will lock the piece of execution when it is being used, restricting read and write access to only one boid at a time. This is used when reading or updating the 2D position array and when updating the velocity- and/or position-vectors.

This will work, especially since there are a relatively small number of boids, however, it is not very efficient. Considering that the parts of the program where a boid needs writing access executes much faster than the part that needs read access (**read**: loops through the whole 2D position array, **write**: update vales for position, velocity and map), this type of mutex will slow down the program.  
We will therefore implement a special type of mutex; the *Readers-Writers lock*.

This lock allows multiple readers but only **one** writer. Additionally, when the writer lock is in use, it also blocks all reading access. In other words, as long as the writer lock is **not** in use, all the boids can read from the readers lock. When a boid uses the writers lock, this blocks **both** the readers- and writers-lock, such that no boid can read new data whilst it is gettin written.

### 1.2. Alignment
The *alignment* rule will make close-by boids to move in the same direction by averaging out the current velocities of all the boids in a given view-space. 
This is accomplished by taking the target vector (where a boid is trying to move) and subtracting it by the current velocity vector. This results in an acceleration vector, which we then add to the original velocity vector which gives us the new velocity vector. To smoothen out this movement, we also multiply that acceleration vector by an adjustment rate (a float between 0 and 1). This will slow down the movement and make it less jittery.


![$\large targetVelocity - currentVelocity = acceleration$](https://latex.codecogs.com/png.latex?%5Cinline%20%5Cdpi%7B150%7D%20%5Cbg_white%20%5Clarge%20targetVelocity%20-%20currentVelocity%20%3D%20acceleration)

![$\large acceleration*adjustmentRate + currentVelocity = newVelocity$ ](https://latex.codecogs.com/png.latex?%5Cinline%20%5Cdpi%7B150%7D%20%5Cbg_white%20%5Clarge%20acceleration*adjustmentRate%20&plus;%20currentVelocity%20%3D%20newVelocity)

### 1.3. Cohesion
The *cohesion* rule will make boids move closer to nearby boids, i.e. the boids will form groups. We will acheive this by finding the average position of nearby boids, and calculate the vector which will bring the current boid to that position. This is done by finding the acceleration vector which will bring the boids together, and then adding it to the already mentioned alignment acceleration and the current boid velocity.

![$\large targetPosition - currentPosition = acceleration$ ](https://latex.codecogs.com/png.latex?%5Cinline%20%5Cdpi%7B150%7D%20%5Cbg_white%20%5Clarge%20targetPosition%20-%20currentPosition%20%3D%20acceleration)

We multiply the acceleration by the adjustment rate and add the result to the previously mentioned *alignment acceleration*.

![$acceleration*adjustmentRate + alignmentAcceleration + currentVelocity = newVelocity$ ](https://latex.codecogs.com/png.latex?%5Cinline%20%5Cdpi%7B150%7D%20%5Cbg_white%20acceleration*adjustmentRate%20&plus;%20alignmentAcceleration%20&plus;%20currentVelocity%20%3D%20newVelocity)


### 1.4. Separation
The *separation* rule will make a boid move away from other boids that are too close. We find the separation acceleration by going through each nearby boid and subtracting the it's position vector from the current position vectors. This will result in multiple vectors (equal to the amount of nearby boids) which we will add together. We divide this result by the distance to the boid, so that the closer it is, the faster it moves away. Multiply the acceleration with the adjustment rate, then add all the accelerations together (alignment, cohesion and separation) and finally add this total acceleration to the current velocity vector to get the new velocity.

### 1.5. Wall-Bounce
The fourth and final element of this program is the wall-bounce. When a boid approaches the end of the screen (the edge/wall), it needs to steer away. We therefore need to calculate the acceleration to move away from the edge. This is done by getting the reciprical distance from the boid to the wall.

![$\large (1/xDistance, 1/yDistance)$ ](https://latex.codecogs.com/png.latex?%5Cinline%20%5Cdpi%7B150%7D%20%5Cbg_white%20%5Clarge%20%281/xDistance%2C%201/yDistance%29)

This bounce acceleration will be larger the closer we are to the wall. We add this to all the other accelerations (alignment, cohesion, separation) and the current velocity. This is the final new velocity.

## 2. Running The Simulation
You need a Go compiler (gc or gccgo).  
You also need the [Erbiten](https://github.com/hajimehoshi/ebiten) library to run this program. 

Compile and create an executable file:  
`go build`  

Then run that executable.  


Example:  
`.\golang-boids`
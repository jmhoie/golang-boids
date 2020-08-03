This program is from the Udemy course: "Mastering Multithreaded Programming with Go (Golang)" by James Cutajar
Link: https://www.udemy.com/course/multithreading-in-go-lang/

In this project we will be using GoRoutines to model Boids. Using concurrent programming to model a problem can make things significantly simpler.

We will be using threads to model birds, each thread represents a distinct bird simulation.
Instead of having a complicated simulation with alot of parameters, we will be using threads where each thread can interract with other threads based on defined properties and ruleset.

Goal: Simulate the flocking behaviour of birds in large groups.

Boid = bird + android

Implementation:
    Plotting dots in a 2D coordinate stystem, where each dot represents a bird.
    Graphics -> Ebiten (2D game engine)
    Boid -> position, velocity, id
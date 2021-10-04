# Express and Typescript

This is a small set of code that runs a stateless JavaScript application, with
no remote storage. For now, all the values will be hardcoded.

## Overview

This purpose of this code is to reproduce some unexpected behavior when using
Linkerd with a node application. Specifically, under load requests to the node
application, return 502 errors and connections appear to be being closed from
the _client_ side. So, in order to properly test this, we'll need a node
application that calls another node application.

## Implementation

Every Friday, I host a [7 minute workout](https://well.blogs.nytimes.com/2013/05/09/the-scientific-7-minute-workout/)
over video call with my coworkers, which currently involves using one device
for the video call, one device as a timer, and a third device to keep track of
the exercises.

I'd like to be able to host the workout in a simpler way, so this application
displays a single page with a button to start the workout. When the "Go!" button
is pressed, there will be a short delay to get set, then each of the exercises
will run sequentially and the JavaScript code will keep track of the active and
rest times.

With this application in place, I can share my screen on the video call and the
team can follow along without me having to manually manage the active and rest
timers.

This application consists of two services running express servers, and one
service will act as a gateway by handling external requests and calling another
service that provides the workout functionality.

The JavaScript code will use TypeScript, because I'd like to get some hands-on
experience with it. :)

_It is entirely possible to do this with one service. The goal is to reproduce
the 502 errors by having a separate application that routes requests to the
underlying "api" service_

### Requirements
* node
* express
* ...

### Deployment
* docker
* kubernetes
* ...

## Future Work

- Maybe hook the service up to persistent storage

package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// define our game states
type State string

const IS_RUNNING State = "running"
const IS_STOPPED State = "stopped"

// define our stop reasons
type StopReason string

const FINISHED StopReason = "user finished the game"
const TIMEOUT StopReason = "game timed out"
const STOPPED StopReason = "user stoped the game"

// current state of the game
var gameState State

// define the time the user has to complete the game
var timeout time.Duration = 5 * time.Second

// define some communication channels
var contactChannel chan bool      // a channel to register all contacts with the wire
var finishChannel chan StopReason // a channel to register finish events (from user or timeout)

// lock concurrent access to the shared variables
var Mutex *sync.Mutex

func init() {
	// create a new lock for concurrent access
	Mutex = &sync.Mutex{}

	// initialize the game state
	setState(IS_STOPPED)

	// initialize our communication channels
	contactChannel = make(chan bool)
	finishChannel = make(chan StopReason)
}

// --- EVENT HANDLERS FOR USER INTERACTION ---

// handleButtonPress will be called on every button press
// (i.e. start, stop, restart)
func handleButtonPress(s interface{}) {
	log.Println("button pressed")
	currentState := getState()

	// start or stop the game
	if currentState == IS_STOPPED {
		go startGame()
		return
	}

	// send a message into our finish channel
	// if the game is running
	select {
	case finishChannel <- STOPPED:
	default:
	}

}

// handleWireContact will be called whenever the user
// touches the wire
func handleWireContact(s interface{}) {
	log.Println("wire touched")
	currentState := getState()

	if currentState != IS_RUNNING {
		fmt.Println("the timer is currently not running")
		return
	}

	// send a contact event to our communication channel
	select {
	case contactChannel <- true:
	default:
	}
}

// handleGameFinished will be called whenever the user
// is touching the finish platform
func handleGameFinished(s interface{}) {
	log.Println("game finished")
	currentState := getState()

	if currentState != IS_RUNNING {
		fmt.Println("the timer is currently not running")
		return
	}

	// send a message into our finish channel
	select {
	case finishChannel <- FINISHED:
	default:
	}

}

// --- ACTUAL GAME ACTIONS ---

// startGame will start a new round
func startGame() {

	// initialize the start time and touch counter
	startTime := time.Now()
	touchCounter := 0

	// set the game state
	setState(IS_RUNNING)

	fmt.Print("\033[H\033[2J")

	done := make(chan struct{})

	// create a separate go-routine for our ticker
	go func(startTime time.Time, done <-chan struct{}) {
		fmt.Println("game started")
		for {
			select {
			// create a ticker to check the time in regular timespans
			case <-time.After(200 * time.Millisecond):
				timeElapsed := time.Now().Sub(startTime)
				fmt.Print("\rTICK: ", timeElapsed.Seconds())

				if timeElapsed >= timeout {
					finishChannel <- TIMEOUT
				}

			case <-done:
				return
			}
		}

	}(startTime, done)

	// increase the touch counter on every touch
	go func(startTime time.Time, counter int, done chan struct{}) {
		for {
			select {
			case <-contactChannel:
				counter += 1

			case reason := <-finishChannel:
				timeElapsed := time.Now().Sub(startTime)

				// set the state of the game to stopped
				setState(IS_STOPPED)

				// print the results
				fmt.Printf(resultLog, reason, timeElapsed.Seconds(), counter)

				close(done)
				return

			case <-done:
				return
			}
		}

	}(startTime, touchCounter, done)

}

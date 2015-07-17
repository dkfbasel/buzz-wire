package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// define our game states
type State int

const (
	RUNNING State = 1 + iota
	STOPPED
)

var gameStart time.Time // start of the game
var gameState State     // current state of the game
var touchCounter int64  // number of touches

var timeout time.Duration = 5 * time.Second // define the time the user has to complete the game
var cancelTimeout chan bool                 // cancel the timeout

var Mutex *sync.Mutex // lock concurrent access to the shared variables

func init() {
	// create a new lock for concurrent access
	Mutex = &sync.Mutex{}

	// initialize the game state
	gameState = STOPPED

	// create a channel to cancel a timeout
	cancelTimeout = make(chan bool)
}

// --- EVENT HANDLERS FOR USER INTERACTION ---

// handleButtonPress will be called on every button press
// (i.e. start, stop, restart)
func handleButtonPress(s interface{}) {
	log.Println("button pressed")
	state := readOutState()

	// start or stop the game
	if state == STOPPED {
		go startGame()
	} else {
		go stopGame()
	}
}

// handleWireContact will be called whenever the user
// touches the wire
func handleWireContact(s interface{}) {
	log.Println("wire touched")
	state := readOutState()

	// only do something if the time is running
	if state == RUNNING {
		go increaseCounter()
	} else {
		fmt.Println("the timer is currently not running")
	}
}

// handleGameFinished will be called whenever the user
// is touching the finish platform
func handleGameFinished(s interface{}) {
	log.Println("game finished")
	state := readOutState()

	if state == RUNNING {
		// game was finished before the timeout
		go finishGame(true)
	} else {
		fmt.Println("the timer is currently not running")
	}
}

func readOutState() State {
	Mutex.Lock()
	state := gameState
	Mutex.Unlock()
	return state
}

// --- ACTUAL GAME ACTIONS ---

// startGame will start a new round
func startGame() {
	Mutex.Lock()

	// get the current time
	currentTime := time.Now()

	// set our game parameters
	gameState = RUNNING
	gameStart = currentTime
	touchCounter = 0

	// simulatedEvents <- "enableLed"

	Mutex.Unlock()

	// create a ticker-channel that can be
	// closed once the game is finished or stopped
	tickerChannel := make(chan bool)

	// create a separate go-routine for our ticker
	go func() {
		for {
			select {
			// TODO: set a millisecond interval that results in a smooth ticker
			case <-time.After(200 * time.Millisecond):
				fmt.Println("TICK: ", time.Now().Sub(currentTime))
			case <-tickerChannel:
				return
			}
		}
	}()

	// wait for timeout or cancellation of the timeout
	select {
	case <-time.After(timeout):
		fmt.Println("The game timed out")
		// game was not finished before the timeout
		go finishGame(false)

		// stop the clock ticker
		close(tickerChannel)

	case <-cancelTimeout:
		// someone cancelled the timeout
		// either by stopping the game or finishing
		fmt.Println("Timeout was canceled")

		// stop the clock ticker
		close(tickerChannel)
	}

}

// stopGame will stop/abort the current round
func stopGame() {
	Mutex.Lock()
	defer Mutex.Unlock()

	// check if the game was already stopped before
	if gameState == STOPPED {
		fmt.Println("Game was already stopped")
		return
	}

	// simulatedEvents <- "disableLed"

	// cancel the timeout from already running game
	// note: this will not block if channel cannot receive
	select {
	case cancelTimeout <- true:
	default:
	}

	// set the game state to stopped
	gameState = STOPPED
}

func increaseCounter() {
	Mutex.Lock()

	// increase the count of contacts
	touchCounter += 1

	Mutex.Unlock()
}

func finishGame(beforeTimeout bool) {

	// get the current time
	finishTime := time.Now()

	// cancel the timeout
	// note: this will not block if channel cannot receive
	select {
	case cancelTimeout <- true:
	default:
	}

	// lock concurrent access
	Mutex.Lock()

	// check if the game was already stopped before
	if gameState == STOPPED {
		fmt.Println("Game was already stopped")
		return
	}

	// simulatedEvents <- "disableLed"

	// set the game state to stopped
	gameState = STOPPED

	// calculate the time that elapsed from start to finish
	timeElapsed := finishTime.Sub(gameStart)
	contacts := touchCounter

	// unlock concurrent access
	Mutex.Unlock()

	// print out the results
	fmt.Printf(`
RESULTS:
---------
Duration:  %v
Contacts:  %v
Timeout:   %v
`, timeElapsed.Seconds(), contacts, !beforeTimeout)

}

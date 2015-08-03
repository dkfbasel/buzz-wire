package main

import (
	"fmt"
	"log"
	"strconv"
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

// define the genders available in our application
type Gender string

const MALE Gender = "male"
const FEMALE Gender = "female"

// current state of the game
var GameState State

// define the time the user has to complete the game
var Timeout time.Duration

// define some communication channels
var contactChannel chan bool          // a channel to register all contacts with the wire
var contactChannelDebounced chan bool // a debouncer for our contact channel
var finishChannel chan StopReason     // a channel to register finish events (from user or timeout)

var DebounceContact time.Duration // how long should the contact channel be debounced

// lock concurrent access to the shared variables
var Mutex *sync.Mutex

// --- EVENT HANDLERS FOR USER INTERACTION ---

// handleButtonPress will be called on every button press
// - i.e. start, stop, restart - and should be passed the gender
// that the button is allocated to (male, female)
func handleButtonPress(gender Gender) {
	log.Println("button pressed:", gender)
	currentState := getState()

	// start or stop the game
	if currentState == IS_STOPPED {
		go startGame(gender)
		return
	}

	// send a message into our finish channel
	// if the game is running
	select {
	case finishChannel <- STOPPED:
	default:
	}

}

// handleStartContact will make sure, that the users is starting from the
// beginning of the wire
func handleStartContact(s interface{}) {
	// TODO: make sure the user starts from the start area
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
	case contactChannelDebounced <- true:
	default:
	}
}

// handleFinishContact will be called whenever the user
// is touching the finish platform
func handleFinishContact(s interface{}) {
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
func startGame(gender Gender) {

	// signal the webserver that we are about to start the game
	signalChannel <- "game::start::" + string(gender)

	// initialize the start time and touch counter
	startTime := time.Now()
	touchCounter := 0

	// set the game state
	setState(IS_RUNNING)

	fmt.Print("\033[H\033[2J")

	done := make(chan struct{})

	// define the led event to use (male or female)
	var ledEvent string
	if gender == FEMALE {
		ledEvent = "enableLedWoman"
	} else {
		ledEvent = "enalbeLedMan"
	}

	// enable the led light
	select {
	case GameEvents <- ledEvent:
	default:
	}

	// create a separate go-routine for our ticker
	go func(startTime time.Time, done <-chan struct{}) {
		fmt.Println("game started")
		for {
			select {
			// create a ticker to check the time in regular timespans
			case <-time.After(200 * time.Millisecond):
				timeElapsed := time.Now().Sub(startTime)
				// fmt.Print("\rTICK: ", timeElapsed.Seconds())

				if timeElapsed >= Timeout {
					finishChannel <- TIMEOUT
				}

			case <-done:
				return
			}
		}

	}(startTime, done)

	// increase the touch counter on every touch
	go func(startTime time.Time, counter int, gender Gender, done chan struct{}) {
		for {
			select {
			case <-contactChannel:
				fmt.Println("register contact")
				counter += 1

				// signal the webserver that the wire was touched
				signalChannel <- "game::contact::" + strconv.Itoa(counter)

				// sound the buzzer
				select {
				case GameEvents <- "soundBuzzer":
				default:
				}

			case reason := <-finishChannel:

				timeElapsed := time.Now().Sub(startTime)

				// signal the webserver that the game was finished
				signalChannel <- "game::finished::" + strconv.FormatFloat(timeElapsed.Seconds(), 'f', 3, 64)

				// set the state of the game to stopped
				setState(IS_STOPPED)

				// print the results
				fmt.Printf(resultLog, gender, reason, timeElapsed.Seconds(), counter)

				// close our waiting channel
				close(done)

			case <-done:
				fmt.Println("game was finished..")
				select {
				case GameEvents <- "ledOff":
				default:
				}
				return
			}
		}

	}(startTime, touchCounter, gender, done)

}

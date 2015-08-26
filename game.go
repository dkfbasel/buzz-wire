package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

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

// handleStartContact will make sure, that the user is starting from the
// beginning of the wire
func handleStartContact(s interface{}) {
	// send a signal down the start channel
	select {
	case startChannel <- true:
	default:
	}
}

// handleWireContact will be called whenever the user
// touches the wire
func handleWireContact(s interface{}) {
	debug("wire touched")
	currentState := getState()

	if currentState != IS_RUNNING {
		debug("the timer is currently not running")
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
	debug("game finished")
	currentState := getState()

	if currentState != IS_RUNNING {
		debug("the timer is currently not running")
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

	// clear the console output
	fmt.Print("\033[H\033[2J")

	// generate a new study id
	studyNumber := strconv.Itoa(random(1000, 9999))
	studyID := StartNumberOfID + studyNumber[0:2] + "-" + studyNumber[2:4]

	// signal the webserver that we are about to start the game
	signal("game::countdown::" + string(gender) + "::" + studyID)
	debug("game countdown started")

	// initialize the start time and touch counter
	touchCounter := 0

	// set the game state
	setState(IS_RUNNING)

	// initialize a channel to end the game-round
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

	// wait three seconds for the counter to finish
	startTime := time.Now().Add(3 * time.Second)

	// create a separate go-routine for our ticker
	go func(startTime time.Time, studyID string, done <-chan struct{}) {
		debug("game started")
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

	}(startTime, studyID, done)

	// increase the touch counter on every touch and react to finish events
	go func(startTime time.Time, counter int, gender Gender, done chan struct{}) {

		// check if the user touched the start of the wire
		var startTouched bool

		for {
			select {
			case <-startChannel:

				// ignore touching the start channel if the game is already started
				if startTouched == true {
					break
				}

				if time.Now().Before(startTime) {
					debug("wait for timer to finish")
					break
				}

				// start region was touched after the timer finished
				debug("start region touched")
				startTouched = true

				// signal the server, that the game should start
				signalChannel <- "game::start::" + string(gender)

			case <-contactChannel:
				debug("register contact")

				// check if the start region has been touched before
				if startTouched == false {
					break
				}

				// increase the touch counter
				counter++

				// signal the webserver that the wire was touched
				signalChannel <- "game::contact::" + strconv.Itoa(counter)

				// sound the buzzer
				select {
				case GameEvents <- "soundBuzzer":
				default:
				}

			case reason := <-finishChannel:

				// only react to events if the run was correctly started
				if startTouched == false && reason == FINISHED {
					break
				}

				timeElapsed := time.Now().Sub(startTime)

				// signal the webserver that the game was finished
				signalChannel <- "game::finished::" + string(reason) + "::" + strconv.FormatFloat(timeElapsed.Seconds(), 'f', 3, 64)

				// set the state of the game to stopped
				setState(IS_STOPPED)

				// print the results
				fmt.Printf(resultLog, gender, reason, timeElapsed.Seconds(), counter)

				// close our waiting channel
				close(done)

			case <-done:
				debug("game was finished..")
				select {
				case GameEvents <- "ledOff":
				default:
				}
				return
			}
		}

	}(startTime, touchCounter, gender, done)

}

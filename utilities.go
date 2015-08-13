package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// getState will get the game state safe for concurrency
func getState() State {
	Mutex.RLock()
	state := GameState
	Mutex.RUnlock()
	return state
}

// setState will set the game state safe for concurrency
func setState(state State) {
	Mutex.Lock()
	GameState = state
	Mutex.Unlock()
}

// signal will send an asynchronous signal on the signal channel
func signal(message string) {
	go func() {
		select {
		case signalChannel <- message:
		default:
		}
	}()
}

// debounceContactChannel will put of responding to a given channel for the specified interval
// i.e. to debounce user contact with the wire
func debounceContactChannel(interval time.Duration, output chan bool) chan bool {

	// initialize the channel to receive input on
	input := make(chan bool)

	go func() {
		var buffer bool
		var ok bool

		// we wait until our input gets called at least once or the input channel is closed
		buffer, ok = <-input
		if !ok {
			return
		}

		// we pass the value from the initial call to our output channel
		output <- buffer

		// we initialize a wait function
		for {
			select {
			// we wait for a signal or closing of the input channel
			case buffer, ok = <-input:
				fmt.Println("debounced ..")

				// exit if the input channel was closed
				if !ok {
					return
				}

			// wait for the given time interval
			case <-time.After(interval):
				// send the data once to the output channel and start waiting again
				output <- <-input
			}
		}
	}()

	return input
}

// random will generate a random number in the given range
func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// --- TESTING UTILITIES ---

// define the string we use to log our results
var resultLog = `
RESULTS:
---------
Gender:      %v
Stop Reason: %v
Duration:    %v
Contacts:    %v
`

var commandsInfo = `
COMMANDS:
---------
button:  User pressed button
contact: User made contact to wire
finish:  User finished the game
`

// parseMessages will parse the command line and send
// the input to a channel
func parseCommandLine(messages chan<- string, done chan<- bool) {

	// use the bufio-scanner to read command line with spaces
	scanner := bufio.NewScanner(os.Stdin)
	var message string

	// keep on scanning
	for scanner.Scan() {
		message = scanner.Text()

		if message == "quit" || message == "exit" {
			// explore option to close channel with close(done)
			done <- true
			continue
		}

		if message == "clear" || message == "cls" {
			// clear the console
			fmt.Print("\033[H\033[2J")
			continue
		}

		if message == "c" {
			message = "contact"
		}

		messages <- message
	}
}

// drawInterface will draw a nice textural gui for our game
func drawInterface(buttonState bool, timeElapsed string) {

	buttonText := " "
	if buttonState == true {
		buttonText = "x"
	}

	timeText := padLeft(timeElapsed, "0", 4)

	// clear the screen
	fmt.Print("\033[H\033[2J")

	// print the current status
	fmt.Printf(`
  ┌────────┐┌───┐     ┌──────────────┐
  │ BUTTON ││ %s │     │ TIMER: %s  │
  └────────┘└───┘     └──────────────┘

  `, buttonText, timeText)

}

// padLeft will pad the given string with the given character until
// the overal given length is reached
func padLeft(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

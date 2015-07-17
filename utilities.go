package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// define the string we use to log our results
var resultLog string = `
RESULTS:
---------
Stop Reason: %v
Duration:    %v
Contacts:    %v
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

// getState will get the game state safe for concurrency
func getState() State {
	Mutex.Lock()
	state := gameState
	Mutex.Unlock()
	return state
}

// setState will set the game state safe for concurrency
func setState(state State) {
	Mutex.Lock()
	gameState = state
	Mutex.Unlock()
}

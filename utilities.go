package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func padLeft(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

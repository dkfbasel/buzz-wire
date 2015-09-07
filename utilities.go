package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
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
				debug("debounced contact")

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

// debug will print a message to the console if in debug mode
func debug(text string) {
	if Mode == "debug" {
		fmt.Println(text)
	}
}

// saveResultToDisk will save the given result to the file system
func saveResultToDisk(studyID string, gender string, duration string, hits string, stopReason string) {

	var file *os.File

	// try to get file statistics (to check whether file exists)
	_, err := os.Stat(outputFile)

	// create a new file or append the exisiting
	if os.IsNotExist(err) {
		if file, err = os.Create(outputFile); err != nil {
			debug("could not create the output file: " + err.Error())
			return
		}
		defer file.Close()
		// add the column description as first entry
		file.WriteString("studyid;gender;duration;hits;stopreason\n")

	} else {
		if file, err = os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0600); err != nil {
			debug("could not open the output file: " + err.Error())
			return
		}
		defer file.Close()
	}

	if _, err = file.WriteString(studyID + ";" + gender + ";" + duration + ";" + hits + ";" + stopReason + "\n"); err != nil {
		debug("could not write to the output file: " + err.Error())
	}

}

// parseExistingStudyResults will try to read all existing study results from
// the result file and create a map of study-ids already in use
func parseExistingStudyResults() (map[string]bool, error) {

	idMap := make(map[string]bool)

	content, err := ioutil.ReadFile(outputFile)
	if err != nil {
		return idMap, err
	}

	// split the file content into separate lines
	lines := strings.Split(string(content), "\n")

	// extract the study id from each line
	for index, line := range lines {

		// skip the first line
		if index == 0 {
			continue
		}

		// split the line into separate rows
		rows := strings.Split(line, ";")

		// select the study id (first row)
		idMap[rows[0]] = true
	}

	return idMap, nil
}

// generateNewStudyId will generate a new study id that is not already in use
func generateNewStudyID() string {
	studyNumber := strconv.Itoa(random(1000, 9999))
	studyID := StartNumberOfID + studyNumber[0:2] + "-" + studyNumber[2:4]

	// make sure the id is not already in use
	if _, ok := idsInUse[studyID]; ok == true {
		studyID = generateNewStudyID()
	}

	// add the study id to the map of already used ids
	idsInUse[studyID] = true

	return studyID
}

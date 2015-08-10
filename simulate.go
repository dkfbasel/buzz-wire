package main

import (
	"fmt"
	"os"
)

// simulate can be used to simulate raspberry pi events
// through console input
func simulate(eventChannel chan<- string) {

	fmt.Println(commandsInfo)

	messageChan := make(chan string)
	done := make(chan bool)

	// start parsing the commandline
	go parseCommandLine(messageChan, done)

	for {
		select {
		case message := <-messageChan:
			if message == "button" || message == "contact" || message == "finish" || message == "off" || message == "start" {
				eventChannel <- message
			}

		case <-done:
			fmt.Println("shutting down")
			os.Exit(0)

		}
	}

}

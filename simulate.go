package main

import (
	"fmt"
	"os"
	"time"
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
			if message == "button" || message == "contact" || message == "finish" || message == "off" {
				eventChannel <- message
			}

		case <-done:
			fmt.Println("shutting down")
			os.Exit(0)

		case <-time.After(30 * time.Second):
			fmt.Println("simulation time has expired, we are shutting down ..")
			os.Exit(0)
		}
	}

}

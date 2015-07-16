package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hybridgroup/gobot"
)

func main() {

	// initialize a new gobot
	gbot := gobot.NewGobot()

	// // initialize a raspberry pi adaptor
	// raspberry := raspi.NewRaspiAdaptor("raspi")

	// // initialize our connection with the corresponding pins
	// rLed := gpio.NewLedDriver(r, "led", "7")
	// rButton := gpio.NewButtonDriver(raspberry, "button", "8")

	// define the events that we will listen to
	buttonPressEvent := gobot.NewEvent()
	contactEvent := gobot.NewEvent()
	finishEvent := gobot.NewEvent()

	// register the events with gobot
	gobot.On(buttonPressEvent, handleButtonPress)
	gobot.On(contactEvent, handleWireContact)
	gobot.On(finishEvent, handleGameFinished)

	// create a channel for simulated events
	simulatedEvents := make(chan string)

	// simulate events from the pi
	go simulate(simulatedEvents)

	// define the work to be done by the robot
	work := func() {
		// don't do anything yet
		// TODO: add handlers for physical events from the raspberry pi
		for event := range simulatedEvents {

			if event == "button" {
				gobot.Publish(buttonPressEvent, nil)
				continue
			}

			if event == "contact" {
				gobot.Publish(contactEvent, nil)
				continue
			}

			if event == "finish" {
				gobot.Publish(finishEvent, nil)
			}

		}
	}

	// define a base robot
	robot := gobot.NewRobot("test", []gobot.Connection{raspberry}, []gobot.Device{}, work)

	// add the robot to the fleet
	gbot.AddRobot(robot)

	// start the robot
	gbot.Start()
}

// simulate can be used to simulate raspberry pi events
// through console input
func simulate(eventChannel chan<- string) {

	fmt.Println(`
COMMANDS:
---------
button:  User pressed button
contact: User made contact to wire
finish:  User finished the game
`)

	messageChan := make(chan string)
	done := make(chan bool)

	// start parsing the commandline
	go parseCommandLine(messageChan, done)

	for {
		select {
		case message := <-messageChan:
			if message == "button" || message == "contact" || message == "finish" {
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

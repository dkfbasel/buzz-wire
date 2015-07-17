package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
)

var simulatedEvents chan string

func main() {

	// initialize a new gobot
	gbot := gobot.NewGobot()

	// initialize a raspberry pi adaptor
	raspberry := raspi.NewRaspiAdaptor("raspi")

	// initialize our connection with the corresponding pins
	led := gpio.NewLedDriver(raspberry, "led", "7")
	button := gpio.NewButtonDriver(raspberry, "button", "11")
	contact := gpio.NewButtonDriver(raspberry, "contact", "15")

	// create a channel for simulated events
	simulatedEvents = make(chan string)

	// simulate events from the pi
	go simulate(simulatedEvents)

	// define the work to be done by the robot
	work := func() {

		// user pushed the start/stop button
		gobot.On(button.Event("push"), handleButtonPress)

		// user made contact with wire
		gobot.On(contact.Event("push"), handleWireContact)

		gobot.Every(1*time.Second, func() {
			fmt.Println("BUTTON IS ACTIVE:", button.Active)
			led.Toggle()
		})

		go func() {

			for event := range simulatedEvents {

				if event == "button" {
					led.Toggle()
					handleButtonPress(nil)
				}

				if event == "contact" {
					handleWireContact(nil)
				}

				if event == "finish" {
					handleGameFinished(nil)
				}

				if event == "enableLed" {
					led.On()
				}

				if event == "disableLed" {
					led.Off()
				}
			}
		}()

	}

	// robot := gobot.NewRobot("buzzwire",
	// 	[]gobot.Connection{},
	// 	[]gobot.Device{},
	// 	work)

	// define a base robot
	robot := gobot.NewRobot("buzzwire",
		[]gobot.Connection{raspberry},
		[]gobot.Device{button, contact, led},
		work)

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

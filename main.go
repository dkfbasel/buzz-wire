package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/raspi"
	"github.com/spf13/viper"
)

const (
	MODE_DEBUG string = "debug"
	MODE_PI    string = "pi"
)

var Mode string = MODE_PI
var GameEvents chan string

func main() {

	initConfiguration()

	// initialize a new gobot
	gbot := gobot.NewGobot()

	// initialize a raspberry pi adaptor
	raspberry := raspi.NewRaspiAdaptor("raspi")

	// start/stop button for a woman
	buttonWoman := gpio.NewButtonDriver(raspberry, "buttonWoman", "11")
	ledWoman := gpio.NewLedDriver(raspberry, "ledWoman", "25") // or 12

	// start/stop buttom for a man
	// buttonMan := gpio.NewButtonDriver(raspberry, "buttonMan", "32")
	// ledMan := gpio.NewLedDriver(raspberry, "ledMan", "18")

	// contact with the wire (start- and finish-area)
	// contactStart := gpio.NewButtonDriver(raspberry, "contactStart", "33")
	contactFinish := gpio.NewButtonDriver(raspberry, "contactFinish", "35")

	// user made contact with the wire (use buzzer to indicate audible)
	contactWire := gpio.NewButtonDriver(raspberry, "contactWire", "15")
	buzzer := gpio.NewLedDriver(raspberry, "buzzer", "16")

	// create a channel for game events
	GameEvents = make(chan string)

	// simulate events with keyboard interaction
	go simulate(GameEvents)

	// define the work to be done by the robot
	work := func() {

		// user pushed the start/stop button
		gobot.On(buttonWoman.Event("push"), func(data interface{}) {
			handleButtonPress(FEMALE)
		})

		// gobot.On(buttonMan.Event("push"), func(data interface{}) {
		// 	handleButtonPress(MALE)
		// })

		// user made contact with wire
		gobot.On(contactWire.Event("push"), handleWireContact)

		// user is starting the game (must touch the starting area)
		// TODO: add handler for starting event
		// gobot.On(contactStart.Event("push"), handleStartContact)

		// user finished the game (touched finish area)
		gobot.On(contactFinish.Event("push"), handleFinishContact)

		go func() {

			for event := range GameEvents {

				switch event {
				// sound the buzzer
				case "soundBuzzer":
					go func() {
						buzzer.On()
						<-time.After(300 * time.Millisecond)
						buzzer.Off()
					}()

				// enable/disable the led for the woman button
				case "enableLedWoman":
					ledWoman.On()
				case "disableLedWoman":
					ledWoman.Off()

				// enable/disable the lef for the man button
				// case "enableLedMan":
				// 	ledMan.On()
				// case "disableLedMan":
				// 	ledMan.Off()

				// disable all leds
				case "ledOff":
					ledWoman.Off()
					// ledMan.Off()

				// simulated events
				case "button":
					handleButtonPress(FEMALE)
				case "contact":
					handleWireContact(nil)
				case "start":
					handleStartContact(nil)
				case "finish":
					handleFinishContact(nil)
				}
			}
		}()

	}

	// we need to define a robot to be used with gobot
	var robot *gobot.Robot

	// switch cases depending on the mode
	if Mode == MODE_DEBUG {
		// debug mode is run on the mac without physical connections
		robot = gobot.NewRobot("buzzwire",
			[]gobot.Connection{},
			[]gobot.Device{},
			work)
	} else {
		// all other modes are run on the pi with physical connections
		robot = gobot.NewRobot("buzzwire",
			[]gobot.Connection{raspberry},
			[]gobot.Device{buttonWoman, ledWoman, contactWire, buzzer, contactFinish},
			work)
	}

	// add the robot to the fleet
	gbot.AddRobot(robot)

	// start the webserver in a separate go routine
	go startServer("localhost:8484")

	// start the robot (blocking)
	gbot.Start()
}

// initialize the configuration
func initConfiguration() {

	// try to load some configuration information from the config file
	viper.SetConfigName("config")
	viper.SetDefault("mode", MODE_PI)
	viper.SetDefault("timeout", 5*time.Second)
	viper.SetDefault("debounceContact", 1000*time.Millisecond)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// read configuration from file or use default
	Timeout = viper.GetDuration("timeout")
	DebounceContact = viper.GetDuration("debounceContact")
	Mode = viper.GetString("mode")

	// create a new lock for concurrent access
	Mutex = &sync.Mutex{}

	// initialize the game state
	setState(IS_STOPPED)

	// initialize our communication channels
	// note: the contactChannel is debounced to handle continuous triggering of the contact
	contactChannel = make(chan bool)
	contactChannelDebounced = debounceContactChannel(DebounceContact, contactChannel)

	// we also need a signal channel to communicate events to the webserver
	signalChannel = make(chan string)

	// handle finish events
	finishChannel = make(chan StopReason)

}

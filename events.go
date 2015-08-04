package main

import (
	"time"

	"github.com/hybridgroup/gobot"
)

func eventHandling() {

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

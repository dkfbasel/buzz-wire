package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

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
	Mutex = &sync.RWMutex{}

	// initialize the game state
	setState(IS_STOPPED)

	// initialize our communication channels
	// note: the contactChannel is debounced to handle continuous triggering of the contact
	contactChannel = make(chan bool)
	contactChannelDebounced = debounceContactChannel(DebounceContact, contactChannel)

	// we need a signal channel to communicate events to the webserver
	signalChannel = make(chan string, 100)

	// handle start events
	startChannel = make(chan bool)

	// handle finish events
	finishChannel = make(chan StopReason)

}

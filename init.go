package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const outputFile = "./results.csv"

var idsInUse map[string]bool

// define our game states
type State string

const IS_RUNNING State = "running"
const IS_STOPPED State = "stopped"

// define our stop reasons
type StopReason string

const FINISHED StopReason = "finished"
const TIMEOUT StopReason = "timeout"
const STOPPED StopReason = "stopped"

// define the genders available in our application
type Gender string

const MALE Gender = "male"
const FEMALE Gender = "female"

// StartNumberOfID depicts the first number to be used for the patient id
var StartNumberOfID string

// GameState reflects current state of the game
var GameState State

// Timeout is used to define the time the user has to complete the game
var Timeout time.Duration

// define some communication channels
var contactChannel chan bool          // a channel to register all contacts with the wire
var contactChannelDebounced chan bool // a debouncer for our contact channel
var startChannel chan bool            // a channel to register touching the start area
var finishChannel chan StopReason     // a channel to register finish events (from user or timeout)

// DebounceContact defines how long the contact channel should not react to new events
var DebounceContact time.Duration

// Mutex will lock concurrent access to the shared variables
var Mutex *sync.RWMutex

// initialize the configuration
func initConfiguration() {

	// try to load some configuration information from the config file
	viper.SetConfigName("config")
	viper.SetDefault("mode", MODE_PI)
	viper.SetDefault("timeout", 5*time.Second)
	viper.SetDefault("debounceContact", 1000*time.Millisecond)
	viper.SetDefault("startNumber", "1")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// read configuration from file or use default
	Timeout = viper.GetDuration("timeout")
	DebounceContact = viper.GetDuration("debounceContact")
	Mode = viper.GetString("mode")
	StartNumberOfID = viper.GetString("startNumber")

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

	// parse all existing results
	idsInUse, _ = parseExistingStudyResults()

}

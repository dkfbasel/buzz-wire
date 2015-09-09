package main

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pborman/uuid"
	"golang.org/x/net/websocket"
)

var signalChannel chan string
var socketPool map[string]chan string
var socketMutex *sync.Mutex

// startServer will start a webserver on the given address and open the default
// browser on the given address
func startServer(address string) {

	server := echo.New()
	// server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	// server static files
	server.Index("website/index.html")

	server.ServeFile("/velocity.min.js", "website/velocity.min.js")
	server.ServeFile("/lodash.min.js", "website/lodash.min.js")
	server.ServeFile("/animations.js", "website/animations.js")
	server.ServeFile("/application.js", "website/application.js")

	server.ServeFile("/style.css", "website/style.css")
	server.ServeFile("/style.css.map", "website/style.css.map")
	server.ServeFile("/style.scss", "website/style.scss")

	// handle websocket connections
	server.WebSocket("/ws", establishSocketConnection)

	// open the page in the default browser
	// go openInDefaultBrowser(address)

	// define a mutex to be used to register sockets
	socketMutex = &sync.Mutex{}

	// initialize our map of socket connections
	socketPool = make(map[string]chan string)

	// run a pool to address all open socket connections
	go handleSocketPool()

	// run the server (blocking)
	server.Run(address)
}

// handleSocketPool will send all messages received on the signalChannel to
// to all registered channels
func handleSocketPool() {

	// wait for signals to be sent to the signalChannel
	for signal := range signalChannel {

		// send the signal to all registered socket channels
		for _, socketChannel := range socketPool {
			select {
			case socketChannel <- signal:
			default:
			}
		}
	}

}

// registerInSocketPool will register the given channel in the map of socket
// connections using the uniqueSocketID as identifier
func registerInSocketPool(uniqueSocketID string, socketChannel chan string) {
	debug("SOCKET ADDED: " + uniqueSocketID)
	socketMutex.Lock()
	socketPool[uniqueSocketID] = socketChannel
	socketMutex.Unlock()
}

// removeFromSocketPool will delete the connection with the given id from the
// pool of socket connections
func removeFromSocketPool(uniqueSocketID string) {
	debug("SOCKET REMOVED: " + uniqueSocketID)
	socketMutex.Lock()
	delete(socketPool, uniqueSocketID)
	socketMutex.Unlock()
}

// establishSocketConnection will handle the live connection between the webpage
// and our game
func establishSocketConnection(c *echo.Context) error {

	// upgrade the connection to a socket
	ws := c.Socket()

	// create a new channel to receive messages
	socketChannel := make(chan string)

	// create a unique id for the socket Connection
	uniqueSocketID := uuid.New()

	// register the socket in the pool
	registerInSocketPool(uniqueSocketID, socketChannel)

	// send signals to client if something is put on the socketChanel
	// note: socketChannel will block until the next receive
	for signal := range socketChannel {

		if ws != (*websocket.Conn)(nil) {
			debug("SIGNAL: " + signal)

			// try to send the signal on the websocket
			err := websocket.Message.Send(ws, signal)

			if err != nil {
				fmt.Println("Socket connection lost")
				// remove the socket from the pool
				removeFromSocketPool(uniqueSocketID)
				// close the receiverChannel (thus also exiting the loop)
				close(socketChannel)
			}

		} else {
			debug("NO CLIENT CONNECTED: " + signal)

			// remove the socket from the pool
			removeFromSocketPool(uniqueSocketID)

			// close the receiverChannel (thus also exiting the loop)
			close(socketChannel)
		}

	}

	return nil
}

// openInDefaultBrowser will open the given address in the users default browser
// (after a short timeout)
func openInDefaultBrowser(address string) {
	<-time.After(1 * time.Second)
	cmd := exec.Command("iceweasel", "--display=:0", "http://localhost:8484", "--fullscreen")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Bitte öffnen Sie Ihren Browser auf der Adresse", address)
	}
}

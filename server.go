package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/net/websocket"
)

var signalChannel chan string

// startServer will start a webserver on the given address and open the default
// browser on the given address
func startServer(address string) {

	server := echo.New()
	// server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	// server static files
	server.Index("website/index.html")
	server.ServeFile("/websocket.min.js", "website/websocket.min.js")
	server.ServeFile("/lodash.min.js", "website/lodash.min.js")
	server.ServeFile("/style.css", "website/style.css")
	server.ServeFile("/style.css.map", "website/style.css.map")
	server.ServeFile("/style.scss", "website/style.scss")

	// handle websocket connections
	server.WebSocket("/ws", handleSocketConnection)

	// run the server (blocking)
	server.Run(address)
}

// showIndexPage will display the main page showing an interactive page
func showIndexPage(c *echo.Context) error {
	return c.String(http.StatusOK, "This is the start-page")
}

// handleSocketConnection will handle the live connection between the webpage
// and our game
func handleSocketConnection(c *echo.Context) error {

	// upgrade the connection to a socket
	ws := c.Socket()

	for {
		// send signals if something is put on the signalChannel
		signal := <-signalChannel
		websocket.Message.Send(ws, signal)

	}

}

/* globals hud, ambulance */

// define some variables to have a ticking timer
var startTime;
var timerReference;


function remainingTime(limit, durationInMillisecs) {

	var timeAsNumber = (limit * 1000) - durationInMillisecs;

	var timeString;
	var timeAsString = timeAsNumber.toFixed(0).toString();

	if (timeAsNumber >= 10000) {
		timeString = timeAsString.substr(0, 2) + ':' + timeAsString.substr(2, 2);

	} else if (timeAsNumber < 10000 && timeAsNumber >= 1000) {
		timeString = '0' + timeAsString.substr(0, 1) + ':' + timeAsString.substr(1, 2);

	} else if (timeAsNumber < 1000 && timeAsNumber >= 100) {
		timeString = '00:' + timeAsString.substr(0,2);

	} else if (timeAsNumber < 100 && timeAsNumber > 0) {
		timeString = '00:0' + timeAsString.substr(0,1);

	} else if (timeAsNumber <= 0) {
		timeString = '00:00';

	}

	return timeString;
}

function showTimer() {
	var currentTime = new Date().getTime();
	var durationInMillisecs = currentTime - startTime;

	hud.elements.time.textContent = remainingTime(20, durationInMillisecs);
}


function handleMessageReceived(event) {

	// get the message of the event
	var message = event.data;

	// split the message into separate tokens
	var tokens = message.split("::");

	switch (tokens[1]) {
		case "start":

			// TODO: set the study id

			// set the number of hits
			hud.elements.trauma.textContent = "0";

			// start a ticking clock
			startTime = new Date().getTime();
			window.clearInterval(timerReference);
			timerReference = window.setInterval(showTimer, 10);

			// start the ambulance
			ambulance.start();
			break;

		case "contact":
			// show the number of contacts with the wire
			hud.elements.trauma.textContent = tokens[2];

			// bump the ambulance
			ambulance.bump();
			break;

		case "finished":

			// clear the ticking clock
			window.clearInterval(timerReference);

			switch (tokens[2]) {
				case "stopped":
					ambulance.stopEarly();
					break;

				case "finished":
					ambulance.stopFinish();

					var elapsedTime = parseFloat(tokens[3]) * 1000;
					hud.elements.time.textContent = remainingTime(20, elapsedTime);

					break;

				case "timeout":
					ambulance.stopTimeout();
					break;
			}
	}
}


// --- WEBSOCKET - HANDLING ---

// create a reconnecting websocket
var loc = window.location;
var uri = 'ws:';

uri += '//' + loc.host;
uri += loc.pathname + 'ws';

function connectSocket() {

	// create a new websocket
	var ws = new WebSocket(uri);

	// attach a function on opening the socket
	ws.onopen = function() {
		console.log('Connected to the game engine')
	}

	// attach a function on closing the socket
	ws.onclose = function() {
		console.log('Closed connection to the game engine');
	}

	ws.onmessage = handleMessageReceived;

}

connectSocket();

// --- TEST FUNCTIONS ---

function start() {
	// start a ticking clock
	startTime = new Date().getTime();
	window.clearInterval(timerReference);
	timerReference = window.setInterval(showTimer, 10);

	// start the ambulance
	ambulance.start();
};

function stop() {
	// clear the ticking clock
	window.clearInterval(timerReference);
};

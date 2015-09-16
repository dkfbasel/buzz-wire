/* globals hud, ambulance, counter, main, _ , hideResult */

// define some variables to have a ticking timer
var startTime;
var timerReference;

var info ={
	id: null,
	remaining: null,
	contacts: null,
	alive: null
};

// the amount of time available to complete the game
var timeLimit = 25;

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

	hud.elements.time.textContent = remainingTime(timeLimit, durationInMillisecs);
}

var debouncedBumping = _.debounce(ambulance.bump, 200, {
	'leading': true,
	'trailing': false
});

function handleMessageReceived(event) {

	// get the message of the event
	var message = event.data;

	// log message to the console
	// console.info("Socket:", message);

	// split the message into separate tokens
	var tokens = message.split("::");

	switch (tokens[1]) {
		case "countdown":

			hideResult();

			info.id = tokens[3];
			info.alive = 1;
			info.contacts = 0;

			// set the study id
			hud.elements.id.textContent = tokens[3];
			hud.elements.resultId.textContent = tokens[3];
			hud.elements.resultPoints.textContent = "25";

			// set the number of hits
			hud.elements.trauma.textContent = "0";

			counter.start(function() {
				// start a ticking clock
				startTime = new Date().getTime();
				window.clearInterval(timerReference);
				timerReference = window.setInterval(showTimer, 10);
			});

			break;

		case "start":
			// start the ambulance, once the user touched the start region
			ambulance.start(tokens[2]);
			break;

		case "contact":
			// show the number of contacts with the wire
			hud.elements.trauma.textContent = tokens[2];
			info.contacts = tokens[2];

			// bump the ambulance
			// ambulance.bump();
			debouncedBumping();
			break;

		case "finished":

			// clear the ticking clock
			window.clearInterval(timerReference);

			switch (tokens[2]) {
				case "stopped":
					ambulance.stopEarly();
					break;

				case "finished":
					info.alive = 1;
					var elapsedTime = parseFloat(tokens[3]) * 1000;
					info.remaining = remainingTime(timeLimit, elapsedTime);
					hud.elements.time.textContent = info.remaining;

					ambulance.stopFinish();
					break;

				case "timeout":
					info.alive = 0;
					info.remaining = remainingTime(timeLimit, timeLimit * 1000);
					hud.elements.time.textContent = info.remaining;
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
		main.showCounter();
		console.log('Connected to the game engine');
	};

	// attach a function on closing the socket
	ws.onclose = function() {
		main.showDisconnect();
		console.log('Closed connection to the game engine');

		// try to reconnect to the socket
		connectSocket();
	};

	ws.onmessage = handleMessageReceived;
}

// automatically connect when opening the page
connectSocket();

// --- TEST FUNCTIONS ---

function start() {
	// start a ticking clock
	startTime = new Date().getTime();
	window.clearInterval(timerReference);
	timerReference = window.setInterval(showTimer, 10);

	// start the ambulance
	ambulance.start();
}

function stop() {
	// clear the ticking clock
	window.clearInterval(timerReference);
}

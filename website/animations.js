/* globals Velocity, console, info */

// --- INITIALIZE ELEMENTS ---

var main = {
	elements: {
		canvas: document.getElementById('canvas'),
		counter: document.getElementById('counter'),
		disconnect: document.getElementById('disconnect')
	},
	showCanvas: null,
	showCounter: null,
	showDisconnect: null
};

var hud = {
	elements: {
		hud: document.getElementById('hud'),
		hudBackground: document.getElementById('hud-background'),

		time: document.getElementById('time-content'),
		trauma: document.getElementById('trauma-content'),

		results: document.getElementById('results'),
		resultId: document.getElementById('results-study-id-content'),
		resultPoints: document.getElementById('results-points-content'),

		id: document.getElementById('id-content'),
		ecg: document.getElementById('ecg-curve-line')
	},
	show: null,
	hide: null
};

var counter = {
	elements: {
		number1: document.getElementById('counter-number-1'),
		number2: document.getElementById('counter-number-2'),
		number3: document.getElementById('counter-number-3')
	},
	start: null,
	hide: null
};

// --- MAIN ELEMENTS ---

// show disconnected message
main.showDisconnect = function() {
	main.elements.disconnect.style.opacity = '1';
	main.elements.counter.style.opacity = '0';
	main.elements.canvas.style.opacity = '0';
};

// show disconnected message
main.showCounter = function() {
	main.elements.disconnect.style.opacity = '0';
	main.elements.counter.style.opacity = '1';
	main.elements.canvas.style.opacity = '0';

	// set the opacity on our numbers
	counter.elements.number3.style.opacity = '1';
	counter.elements.number2.style.opacity = '0';
	counter.elements.number1.style.opacity = '0';
};

main.showCanvas = function() {
	main.elements.disconnect.style.opacity = '0';
	main.elements.counter.style.opacity = '0';
	main.elements.canvas.style.opacity = '1';
};

// @codekit-append "animations/_counter.js";
// @codekit-append "animations/_results.js";
// @codekit-append "animations/_ambulance.js";

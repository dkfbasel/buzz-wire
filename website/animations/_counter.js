var counter = {
	elements: {
		number1: document.getElementById('counter-number-1'),
		number2: document.getElementById('counter-number-2'),
		number3: document.getElementById('counter-number-3')
	},
	start: null,
	hide: null
};

// --- COUNTER ---

// start the counter
counter.start = function(callback) {

	// hide the main canvas
	main.elements.disconnect.style.opacity = '0';
	main.elements.counter.style.opacity = '1';
	main.elements.canvas.style.opacity = '0';

	// set the opacity on our numbers
	counter.elements.number3.style.opacity = '1';
	counter.elements.number2.style.opacity = '0';
	counter.elements.number1.style.opacity = '0';

	window.setTimeout(function() {
		counter.elements.number3.style.opacity = '0';
		counter.elements.number2.style.opacity = '1';
		counter.elements.number1.style.opacity = '0';

		window.setTimeout(function() {
			counter.elements.number3.style.opacity = '0';
			counter.elements.number2.style.opacity = '0';
			counter.elements.number1.style.opacity = '1';

			window.setTimeout(function() {

				main.elements.counter.style.opacity = '0';
				main.elements.canvas.style.opacity = '1';

				if (callback) {
					callback();
				}
			}, 1000);

		}, 1000);

	}, 1000);

};

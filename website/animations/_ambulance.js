/* globals Velocity, main, showResult, hud */

var ambulance = {
	elements: {
		group: document.getElementById('ambulance'),
		protest: document.getElementById('protest'),
		car: document.getElementById('car-animation'),
		driving: document.getElementById('car-driving-animation'),
		bump: document.getElementById('car-bump-animation'),
		siren: document.getElementById('siren'),

		carBackground: document.getElementById('car-background'),
		carBandHorizontal: document.getElementById('band-horizontal'),
		carBandVertical: document.getElementById('band-vertical'),
		carBandVerticalPart: document.getElementById('band-vertical-part'),
		carWindow: document.getElementById('window'),

		carSignVertical: document.getElementById('sign-vertical'),
		carSignRotateLeft: document.getElementById('sign-rotate-left'),
		carSignRotateRight: document.getElementById('sign-rotate-right'),

		wheelFront: document.getElementById('wheel-front'),
		wheelBack: document.getElementById('wheel-back'),

		wheelFrontRim: document.getElementById('wheel-front-rim'),
		wheelFrontRimHighlight: document.getElementById('wheel-front-rim-highlight'),
		wheelBackRim: document.getElementById('wheel-back-rim'),
		wheelBackRimHighlight: document.getElementById('wheel-back-rim-highlight'),

		driver: document.getElementById('driver'),
		driverMan: document.getElementById('man'),
		driverWoman: document.getElementById('woman')

	},
	start: null,
	bump: null,
	stopEarly: null,
	stopTimeout: null,
	stopFinish: null,
	isStarted: null
};

var hospital = {
	elements: {
		group: document.getElementById('hospital'),
		door: document.getElementById('door')
	}
};

// --- INITIALIZE THE START POSITION OF OUR CANVAS

// set transform origin on wheels to allow for rotation
Velocity.hook(
	[ambulance.elements.wheelBack, ambulance.elements.wheelFront],
	'transformOrigin', '36px 35px');

// transform origin of rotated sign elements
Velocity.hook(
	[ambulance.elements.carSignRotateLeft, ambulance.elements.carSignRotateRight],
	'transformOrigin', '37px 37px'
);

Velocity.hook(ambulance.elements.carSignRotateLeft, 'rotateZ', '-60deg');
Velocity.hook(ambulance.elements.carSignRotateRight, 'rotateZ', '60deg');


// transform car origin for rotation
Velocity.hook(ambulance.elements.car, 'transformOrigin', '195px 260px');

// --- AMBULANCE ---

// make the wheels rotating
var rotateWheel = function(element, startValue) {
	Velocity(element, {
		rotateZ: [720 + startValue, startValue]
	}, {
		duration: 1200,
		easing: "linear",
		queue: false,
		complete: function() {
			rotateWheel(element, startValue + 720);
		}
	});
};


// start the ambulance
ambulance.start = function(gender) {

	ambulance.isStarted = true;

	main.elements.disconnect.style.opacity = '0';
	main.elements.counter.style.opacity = '0';
	main.elements.canvas.style.opacity = '1';

	// reset the ecg
	Velocity.hook(hud.elements.ecg, 'translateX', '0');

	// set woman and man driver to hidden
	ambulance.elements.driverWoman.style.visibility = 'hidden';
	ambulance.elements.driverMan.style.visibility = 'hidden';

	switch (gender) {
		case "female":
			ambulance.elements.driverWoman.style.visibility = 'visible';
			break;

		case "male":
			ambulance.elements.driverMan.style.visibility = 'visible';
			break;
	}

	// reset the car if it was stopped through timeout
	Velocity.hook(ambulance.elements.carSignVertical, 'opacity', '1');

	ambulance.elements.carSignRotateLeft.style.fill = '';
	ambulance.elements.carSignRotateRight.style.fill = '';

	Velocity.hook(ambulance.elements.carSignRotateLeft, 'rotateZ', '-60deg');
	// Velocity.hook(ambulance.elements.carSignRotateLeft, 'fill', '#58AACF');

	Velocity.hook(ambulance.elements.carSignRotateRight, 'rotateZ', '60deg');
	// Velocity.hook(ambulance.elements.carSignRotateRight, 'fill', '#58AACF');
	Velocity.hook(ambulance.elements.carSignRotateRight, 'translateX', '0');
	Velocity.hook(ambulance.elements.carSignRotateRight, 'scaleY', '1');


	ambulance.elements.carBackground.style.fill = '';
	// Velocity.hook(ambulance.elements.carBackground, 'fill', '#FEF5E4');

	Velocity.hook(ambulance.elements.siren, 'translateY', '0px');

	ambulance.elements.carBandHorizontal.style.fill = '';
	ambulance.elements.carBandVertical.style.fill = '';
	ambulance.elements.carBandVerticalPart.style.fill = '';

	// Velocity.hook(ambulance.elements.carBandHorizontal, 'fill', '#EB5151');
	// Velocity.hook(ambulance.elements.carBandVertical, 'fill', '#EFE4CE');
	// Velocity.hook(ambulance.elements.carBandVerticalPart, 'fill', '#C84648');

	ambulance.elements.wheelFrontRim.style.fill = '';
	ambulance.elements.wheelBackRim.style.fill = '';
	ambulance.elements.wheelFrontRimHighlight.style.fill = '';
	ambulance.elements.wheelBackRimHighlight.style.fill = '';
	// Velocity.hook([ambulance.elements.wheelBackRim, ambulance.elements.wheelFrontRim], 'fill', '#453D30');
	// Velocity.hook([ambulance.elements.wheelBackRimHighlight, ambulance.elements.wheelFrontRimHighlight], 'fill', '#5E564B');

	ambulance.elements.carWindow.style.fill = '';


	// initialize start position of the car
	Velocity.hook(ambulance.elements.group, 'translateX', '0px');
	Velocity.hook(ambulance.elements.car, 'rotateY', '0deg');
	Velocity.hook(ambulance.elements.car, 'scale', '1');

	// add horizontal motion to the car
	Velocity(ambulance.elements.group, {
		translateX: [760, 0],
		translateY: [0, 0]
	}, {
		duration: 1000,
		easing: "ease-in"
	});

	// make car vibrate
	Velocity(ambulance.elements.driving, {
		translateY: [1.5, 0]
	}, {
		duration: 80,
		loop: true,
		easing: "swing"
	});

	// make the wheel rotating
	rotateWheel(ambulance.elements.wheelFront, 0);
	rotateWheel(ambulance.elements.wheelBack, 0);

};

// bump the ambulance on wire contact
ambulance.bump = function() {

	if (!ambulance.isStarted) {
		return;
	}

	Velocity([ambulance.elements.protest, ambulance.elements.bump], 'stop');

	// NOTE: somehow there is a problem with the car rotation when
	// trying to animate a car bump. In addition, the RaspberryPi
	// seems to have performance problems with it.

	// Velocity.hook(ambulance.elements.car, 'rotateY', '0deg');
	// Velocity.hook(ambulance.elements.car, 'scale', '1');

	Velocity(ambulance.elements.bump, {
		translateY: [-8, 0]
	}, {
		duration: 80,
		easing: "ease-out-bounce",
		loop: 2
	});

	Velocity(ambulance.elements.protest, {
		opacity: [1, 0]
	}, {
		duration: 200,
		delay: 50,
		queue: false,
		easing: "swing",
		complete: function() {
			Velocity(ambulance.elements.protest, {
				opacity: 0
			}, {
				duration: 200,
				delay: 300,
				queue: false
			});
		}
	});
};

// user stopped game early (turn ambulance around)
ambulance.stopEarly = function(rotationDuration, returnDuration, callback) {

	// prevent animation if ambulance is not started
	if (!ambulance.isStarted) {
		return;
	}
	ambulance.isStarted = false;

	if (!returnDuration) {
		returnDuration = 1000;
	}

	if (!rotationDuration) {
		rotationDuration = 120;
	}

	Velocity([ambulance.elements.car, ambulance.elements.protest, ambulance.elements.bump], 'finish');

	Velocity(ambulance.elements.car, {
		rotateY: [180, 0],
	}, {
		duration: rotationDuration,
		easing: "ease-in"
	});

	Velocity(ambulance.elements.group, {
		translateX: 0,
	}, {
		duration: returnDuration,
		easing: "ease-in-out",
		delay: 80,
		complete: function() {
			Velocity(ambulance.elements.driving, 'stop');
			if (callback) {
				callback();
			}
		}
	});

	// hide the protest sign
	Velocity(ambulance.elements.protest, {
		opacity: 0
	}, {
		duration: 500,
		complete: function() {
			Velocity(ambulance.elements.protest, 'stop');
		}
	});
};


// games is stopped due to timeout (convert ambulance into hearse-car)
ambulance.stopTimeout = function() {

	// prevent animation if ambulance is not started
	if (!ambulance.isStarted) {
		return;
	}

	Velocity([ambulance.elements.car, ambulance.elements.protest, ambulance.elements.bump], 'finish');

	Velocity.hook(ambulance.elements.carSignVertical, 'opacity', '0');

	Velocity(hud.elements.ecg, {
		translateX: [-274, 0]
	}, {
		duration: 900,
		easing: "linear"
	});

	Velocity(ambulance.elements.carSignRotateLeft, {
		rotateZ: [720, -60],
		fill: '#ffffff'
	}, {
		duration: 400,
	});

	Velocity(ambulance.elements.carSignRotateRight, {
		rotateZ: [810, 60],
		fill: '#ffffff'
	}, {
		duration: 400,
		queue: false
	});


	Velocity(ambulance.elements.carSignRotateRight, {
		translateX: -10,
		scaleY: 0.7,
	}, {
		duration: 100,
		delay: 200,
		queue: false
	});

	Velocity(ambulance.elements.carBackground, {
		fill: '#565656'
	}, {
		duration: 800
	});

	Velocity(ambulance.elements.siren, {
		translateY: 40
	}, {
		duration: 800,
		easing: "ease-out"
	});

	Velocity(ambulance.elements.carBandHorizontal, {
		fill: '#413050'
	}, {
		duration: 800
	});

	Velocity(ambulance.elements.carBandVertical, {
		fill: '#6C6C6C'
	}, {
		duration: 800
	});

	Velocity(ambulance.elements.carBandVerticalPart, {
		fill: '#512B72'
	}, {
		duration: 800
	});

	Velocity([ambulance.elements.wheelBackRim, ambulance.elements.wheelFrontRim], {
		fill: '#9B9B9B'
	}, {
		duration: 800
	});

	Velocity([ambulance.elements.wheelBackRimHighlight, ambulance.elements.wheelFrontRimHighlight], {
		fill: '#BBBBBB'
	}, {
		duration: 800
	});

	Velocity(ambulance.elements.carWindow, {
		fill: '#667586'
	}, {
		duration: 800
	});

	Velocity(ambulance.elements.group, {
		translateX: 840,
	}, {
		duration: 600,
		easing: 'linear',
		delay: 100,
		complete: function() {
			ambulance.stopEarly(300, 2500, showResult);
		}
	});

	// hide the protest sign
	Velocity(ambulance.elements.protest, {
		opacity: 0
	}, {
		duration: 500,
		complete: function() {
			Velocity(ambulance.elements.protest, 'stop');
		}
	});

};

ambulance.stopFinish = function() {

	// prevent animation if ambulance is not started
	if (!ambulance.isStarted) {
		return;
	}
	ambulance.isStarted = false;

	Velocity([ambulance.elements.car, ambulance.elements.protest, ambulance.elements.bump], 'finish');

	// open the hospital door
	Velocity(hospital.elements.door, {
		translateY: [-100, 0],
	}, {
		duration: 300,
		easing: 'ease-in',
		queue: false
	});

	// add horizontal motion to the car
	Velocity(ambulance.elements.group, {
		translateX: 1760,
		translateY: -110
	}, {
		duration: 800,
		easing: [0.455, 0.03, 0.515, 0.955],
		queue: false,
		delay: 300,
		complete: function() {

			// close the hospital door
			Velocity(hospital.elements.door, {
				translateY: 0,
			}, {
				duration: 200,
				easing: 'ease-out-quad',
			});

			showResult();
		}
	});

	Velocity(ambulance.elements.car, {
		scale: 0.3,
	}, {
		duration: 800,
		queue: false,
		easing: [0.455, 0.03, 0.515, 0.955],
		delay: 280,
		complete: function() {
			Velocity(ambulance.elements.driving, 'stop');
			Velocity(ambulance.elements.wheelFront, 'stop');
			Velocity(ambulance.elements.wheelBack, 'stop');
		}
	});

	// hide the protest sign
	Velocity(ambulance.elements.protest, {
		opacity: 0
	}, {
		duration: 500,
		queue: false,
		complete: function() {
			Velocity(ambulance.elements.protest, 'stop');
		}
	});

};

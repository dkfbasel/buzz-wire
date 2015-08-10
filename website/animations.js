/* globals Velocity, console */

var ambulance = {
	elements: {
		group: document.getElementById('ambulance'),
		protest: document.getElementById('protest'),
		car: document.getElementById('car-animation'),
		driving: document.getElementById('car-driving-animation'),
		siren: document.getElementById('siren'),

		carBackground: document.getElementById('background'),
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

	},
	start: null,
	bump: null,
	stopEarly: null,
	stopTimeout: null,
	stopFinish: null
};

var hospital = {
	elements: {
		group: document.getElementById('hospital')
	}
};

var hud = {
	elements: {
		hud: document.getElementById('hud'),
		time: document.getElementById('time-content'),
		trauma: document.getElementById('trauma-content'),
		id: document.getElementById('id-content'),
		ecg: document.getElementById('ecg-content')
	},
	show: null,
	hide: null
};

// set transform origin on wheels to allow for rotation
ambulance.elements.wheelFront.style.transformOrigin = '36px 35px';
ambulance.elements.wheelBack.style.transformOrigin = '36px 35px';

// transform origin of rotated sign elements
ambulance.elements.carSignRotateLeft.style.transformOrigin = "50% 50%";
ambulance.elements.carSignRotateRight.style.transformOrigin = "50% 50%";
Velocity.hook(ambulance.elements.carSignRotateLeft, 'rotateZ', '-60deg');
Velocity.hook(ambulance.elements.carSignRotateRight, 'rotateZ', '60deg');


// transform car origin for rotation
ambulance.elements.car.style.transformOrigin = '50% 100%';

// transform origin of hospital for scaling
hospital.elements.group.style.transformOrigin = '100% 100%';
hospital.elements.group.style.transform = 'scale(0.5)';
hospital.elements.group.style.visibility = 'visible';

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
ambulance.start = function() {

	// initialize start position of the car
	Velocity.hook(ambulance.elements.group, 'translateX', '0px');
	Velocity.hook(ambulance.elements.car, 'rotateY', '0deg');
	Velocity.hook(ambulance.elements.car, 'scale', '1');

	// add horizontal motion to the car
	Velocity(ambulance.elements.group, {
		translateX: [760, 0]
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

	// increase the size of the hospital
	Velocity(hospital.elements.group, {
		scale: [1, 0.5]
	}, {
		duration: 9000,
		delay: 500,
		easing: 'linear'
	});

};

// bump the ambulance on wire contact
ambulance.bump = function() {

	Velocity.hook(ambulance.elements.car, 'rotateY', '0deg');
	Velocity.hook(ambulance.elements.car, 'scale', '1');

	Velocity(ambulance.elements.car, {
		translateY: [-15, 0],
	}, {
		duration: 70,
		easing: "ease-out-bounce",
		loop: 2
	});

	Velocity(ambulance.elements.protest, {
		opacity: [1, 0]
	}, {
		duration: 200,
		delay: 50,
		easing: "swing",
		complete: function() {
			Velocity(ambulance.elements.protest, 'reverse', {
				duration: 200,
				delay: 400
			});
		}
	});
};


// user stopped game early (turn ambulance around)
ambulance.stopEarly = function(rotationDuration, returnDuration) {

	if (!returnDuration) {
		returnDuration = 1000;
	}

	if (!rotationDuration) {
		rotationDuration = 120;
	}

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
		delay: 80
	});

	Velocity(hospital.elements.group, 'stop');

	// increase the size of the hospital
	Velocity(hospital.elements.group, {
		scale: 0.5
	}, {
		duration: returnDuration - 100,
		easing: 'linear'
	});
};


// games is stopped due to timeout (convert ambulance into hearse-car)
ambulance.stopTimeout = function() {

	Velocity.hook(ambulance.elements.carSignVertical, 'opacity', '0');

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
			ambulance.stopEarly(500, 3000);

			Velocity(hospital.elements.group, 'stop');
			// increase the size of the hospital
			Velocity(hospital.elements.group, {
				scale: 0.5
			}, {
				duration: 2000,
				easing: 'linear'
			});
		}
	});

};

ambulance.stopFinish = function() {

	// add horizontal motion to the car
	Velocity(ambulance.elements.group, {
		translateX: 1860,
		translateY: -50,
	}, {
		duration: 800,
		easing: [0.455, 0.03, 0.515, 0.955]
	});

	// stop all previous hopsital group animations
	Velocity(hospital.elements.group, 'stop');

	// scale the hospital to full size
	Velocity(hospital.elements.group, {
		scale: 1
	}, {
		duration: 800,
		easing: [0.455, 0.03, 0.515, 0.955]
	});

	Velocity(ambulance.elements.car, {
		scale: 0.5,
	}, {
		duration: 800,
		easing: [0.455, 0.03, 0.515, 0.955],
		complete: function() {
			Velocity(ambulance.elements.driving, 'stop');
			Velocity(ambulance.elements.wheelFront, 'stop');
			Velocity(ambulance.elements.wheelBack, 'stop');
		}
	});

};

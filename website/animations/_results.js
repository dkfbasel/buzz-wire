/* globals Velocity, hud, info */

// --- RESULTS ---

function animateScore(start, stop, callbackFn) {
	var difference = Math.abs(start - stop);

	Velocity(hud.elements.resultPoints, {
		tween: [stop, start]
	}, {
		duration: (difference * 80),
		easing: "ease-in",
		delay: 250,
		progress: function(elements, complete, remaining, startTime, tweenValue) {
			var value = parseInt(tweenValue, 10);
			elements[0].textContent = value;
		},
		complete: callbackFn
	});
}


function hideResult() {
	Velocity.hook(hud.elements.results, 'opacity', '0');
}

function showResult() {

	// show the results view
	Velocity(hud.elements.results, {
		opacity: [1, 1],
		translateX: [50, -1200],
		translateY: [195, 195]
	}, {
		duration: 400,
		easing: [80, 11],
		complete: function() {

			var remainingSeconds;

			if (!info.remaining) {
				remainingSeconds = 0;
			} else {
				remainingSeconds = parseInt(info.remaining.split(":")[0], 10);
			}
			// add seconds to score
			var score1 = 25 + remainingSeconds;

			animateScore(25, score1, function() {

				// remove points for contacts
				var score2 = score1 - (2 * info.contacts);

				if (score2 < 0) {
					score2 = 0;
				}

				animateScore(score1, score2, function() {

					var score3 = score2;

					// remove points for dead patient
					if (info.alive === 0) {
						score3 = score2 - 20;
					}

					// score cannot be below zero
					if (score3 < 0) {
						score3 = 0;
					}

					if (score3 !== score2) {
						animateScore(score2, score3);
					}

				});

			});

		}
	});

}

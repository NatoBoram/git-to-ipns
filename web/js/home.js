if (!!!gi) var gi = {};
gi.home = {
	submit: function () {
		let link = document.querySelector("#gitURLInput").value;

		document.querySelector("#card_result").innerHTML = templates.spinner_grow.render();

		// POST
		fetch("/api/add/", {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ URL: link }),
		})
			.then(response => response.json())
			.then(response => {

				// Apply result
				document.querySelector("#card_result").innerHTML = templates.result_card.render(response);

			})
			.catch(error => console.error(error));
	}
};
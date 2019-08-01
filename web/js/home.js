if (!!!gipns) var gipns = {};
gipns.home = {

	submit: function () {
		const link = document.querySelector("#gitURLInput").value;
		document.querySelector("#card_result").innerHTML = templates.spinner_grow.render();

		// POST
		fetch("/api/repos/", {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ URL: link }),
		}).then(response => {

			// Status
			if (response.status !== 200) {
				return gipns.home.error(response, "#card_result");
			}

			// Apply
			response.json().then(response => {
				document.querySelector("#card_result").innerHTML = templates.result_card.render(response);
				gipns.home.list();
			});
		}).catch(error => console.error(error));
	},

	list: function () {

		// GET
		fetch("/api/repos/", {
			method: "GET",
		}).then(response => {

			// Status
			if (response.status !== 200) {
				return gipns.home.error(response, "#list_alert");
			}

			// Apply
			response.json().then(response => {
				document.querySelector("#repos_table").innerHTML = templates.repos_table.render({ repos: response });
			});
		}).catch(error => console.error(error));
	},

	delete: function (url) {
		document.querySelector("#list_alert").innerHTML = templates.spinner_grow.render();

		// DELETE
		fetch("/api/repos/" + url, {
			method: "DELETE",
		}).then(response => {

			// Status
			if (response.status !== 200) {
				return gipns.home.error(response, "#list_alert");
			}

			// Apply
			document.querySelector("#list_alert").innerHTML = templates.alert_success.render({ message: "Successfully deleted the repo." });
			gipns.home.list();
		}).catch(error => console.error(error));
	},

	error: function (response, html_id) {
		console.error(response);
		const alert_content = {
			heading: response.status + " : " + response.statusText,
			message: ""
		};

		response.json().then(response => {
			alert_content.message = response.message;
			document.querySelector(html_id).innerHTML = templates.alert_danger.render(alert_content);
		});
	}
};

// Startup
gipns.home.list();

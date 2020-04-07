const host = window.location;
const status_list_uri = 'status/list';
const url = host + status_list_uri;


function askAPI(url) {

	fetch(url)
		.then(function(response) {
		return response.json();
		})
		.then(function(data) {
			prepareData(data);
	});

}

askAPI(url);

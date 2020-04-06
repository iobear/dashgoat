const protocol = window.location.protocol + '//';
const host = window.location.host;
const status_list_uri = '/dg/status/list';
let url = protocol + host + status_list_uri;

function askAPI(url) {
	//console.log(new Date());

	fetch(url)
		.then(function(response) {
		return response.json();
		})
		.then(function(data) {
			prepareData(data);
	});

}

askAPI(url);

//setInterval(askAPI(url), 2000);


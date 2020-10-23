const host = window.location;
const status_list_uri = 'status/list';
const health_uri = 'health';


function askAPI()
{
	const url = host + status_list_uri;

	fetch(url)
		.then(function(response)
		{
		return response.json();
		})
		.then(function(data)
		{
				prepareData(data);
		});

}


function askHealth()
{
	const url = host + health_uri;

	fetch(url)
		.then(function(response)
		{
		return response.json();
		})
		.then(function(data)
		{
				updateVersion(data);
		});

}


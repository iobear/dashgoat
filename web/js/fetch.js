const host = window.location;
const status_list_uri = 'status/list';
const url = host + status_list_uri;


function askAPI()
{
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

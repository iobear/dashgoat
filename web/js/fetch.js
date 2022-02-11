const host = window.location;
const status_list_uri = 'status/list';
const health_uri = 'health';
const sleep = (sec) => {
	return new Promise(resolve => setTimeout(resolve, sec * 1000))
}

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
		}).catch(()=>{
			waitForBackend();
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


function waitForBackend()
{
	console.log('Waiting for backend to come alive');

	sleep(60).then(() => {
		askAPI();
	})

}

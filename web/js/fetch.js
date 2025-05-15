/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

const host = window.location.origin + window.location.pathname;
let status_list_uri = 'status/list';
const health_uri = 'health';
const sleep = (sec) => {
	return new Promise(resolve => setTimeout(resolve, sec * 1000))
}

function askAPI()
{
	search = getSearchQuery();
	var url = host + status_list_uri;

	if (search != '') {
		url = 'status/listsearch/' + getSearchQuery();
	}

	fetch(url)
		.then(function(response)
		{
			if (response.status == 204){
				if (search === '') {
					return tellDashboard("Waiting for first update", "info");
				} else {
					return tellDashboard("Empty search result", "info");
				}
			}
			return response.json();
		})
		.then(function(data)
		{
			prepareData(data);
		}).catch(err =>{
			console.log(err);
			waitForBackend();
		});

}

function tellDashboard(message, status) {
	let result = {}
	let msg = {}

	msg["message"] = message;
	msg["status"] = status;
	msg["service"] = "JS";
	msg["host"] = "localhost";
	msg["severity"] = "info";
	msg["change"] = 0;
	msg["probe"] = 0;

	result["localhostJS"] = msg;

	return result;
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
			insertBackendAppData(data);
		});

}

function waitForBackend()
{
	msg = 'Waiting for backend to come alive'

	prepareData(tellDashboard(msg, 'warning'), false);
	console.log('Waiting for backend to come alive');

	sleep(4).then(() => {
		askAPI();
	})

}

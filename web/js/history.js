/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

const metrics_uri = 'metricshistory/';


function drawTimelineTd(data, rowid) {

	let rowToUpdate = document.getElementById(rowid);

	// Create a wrapper div
	let wrapperDiv = document.createElement('div');
	wrapperDiv.classList.add('timeline-wrapper');
	wrapperDiv.id = `${rowid}c`;

	// Create 288 fields (288 * 5 min is 24 hours) inside the wrapper div
	for (let i = 0; i < 288; i++)
	{
		const field = document.createElement('div');
		field.classList.add('timeline-field');
		field.title = 'No data';
		wrapperDiv.appendChild(field); // Append each field to the wrapper div
	}
	rowToUpdate.appendChild(wrapperDiv);

	const startOfDay = new Date();
	startOfDay.setHours(0, 0, 0, 0); // Set to 00:00:00 of the current day

	// Update fields based on data
	data.forEach(item =>
	{
		const timeDiff = item.Timestamp * 1000 - startOfDay.getTime(); // Assuming Timestamp is in milliseconds
		const index = Math.floor(timeDiff / (5 * 60 * 1000)); // Convert to 5-minute intervals
	
		const field = wrapperDiv.children[index];
		if (field)
		{
			field.style.backgroundColor = getStatusColor(item.Status);
			// field.title = `${item.Status} at ${new Date(item.Timestamp * 1000).toLocaleString()}`;
			field.title = `${item.Status} at ${unix2timeDay(item.Timestamp)}`;
		}
	});
}

function getStatusColor(status)
{
	switch (status) {
		case 'critical': return '#ef1b11'; // Red
		case 'error': return '#db7600'; // Orange
		case 'warning': return '#e2dc0c'; // Yellow
		case 'info': return '#14a7c9'; // Blue
		case 'ok': return '#00c080'; // Green
		default: return '#808080'; // Gray
	}
}

function askMetrics(hostservice, rowid)
{
	let url = host + metrics_uri + hostservice + '/24'; //24 Hour hard coded

	fetch(url)

		.then(response => response.json())

		.then(data =>
		{
			drawTimelineTd(data, rowid + 'h');
		});
}

function displayHistory(service, rowid)
{
	if (backend_app_data['metrics_history'] == false)
	{
		alert('Metrics history not enabled');
		return;
	}

	if (service == "")
	{
		console.log('Missing service data');
	}

	if (rowid == "")
	{
		console.log('Missing rowid');
	}

	//Toggle history display
	let rowToUpdate = document.getElementById(rowid+'h');

	if (document.getElementById(rowid+'hc'))
	{
		//Already shown, removing..
		rowToUpdate.innerHTML = "";
		rowToUpdate.setAttribute('style', 'display: none;');
		return;
	}
	else
	{
		//None found, building..
		rowToUpdate.removeAttribute('style', 'display: none;');
	}

	hostservice =  service['host']+service['service'];
	askMetrics(hostservice, rowid);
}

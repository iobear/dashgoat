/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

function timeDiff(unixtime)
{
	const endDate = new Date();
	const startDate = new Date(unixtime * 1000);

	var diff = endDate.getTime() - startDate.getTime();
	var days = Math.floor(diff / (1000 * 60 * 60 * 24));
	var hours = Math.floor(diff / (1000 * 60 * 60));
	var minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

	if (days >= 1)
	{
		hours = Math.ceil((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
		return (days <= 9 ? "0" : "") + days + "D" + (hours <= 9 ? "0" : "") + hours + "H";
	}
	else
	{
		return (hours <= 9 ? "0" : "") + hours + "H" + (minutes <= 9 ? "0" : "") + minutes;
	}
}

function unix2timeDay(unixTimestamp) {
	const daysOfWeek = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
	const date = new Date(unixTimestamp * 1000);
	const hours = date.getHours().toString().padStart(2, '0');
	const minutes = date.getMinutes().toString().padStart(2, '0');
	const dayName = daysOfWeek[date.getDay()];

	const formattedDateTime = `${hours}:${minutes} ${dayName}`;

	return formattedDateTime;
}


function updateFaviconColor(color) {
	var favicon = document.getElementById('dynamic-favicon');
	var svgMarkup = `
		<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100">
			<circle cx="50" cy="50" r="40" stroke="black" stroke-width="3" fill="${color}" />
		</svg>
	`;
	favicon.href = "data:image/svg+xml," + encodeURIComponent(svgMarkup);
}

function getStatusColor(statusStr) {
	switch (statusStr) {
		case 'critical': return '#ef1b11'; // Red
		case 'error': return '#ef8511'; // Orange
		case 'warning': return '#e2dc0c'; // Yellow
		case 'info': return '#14a7c9'; // Blue
		case 'ok': return '#33ff00'; // Green
		default: return '#808080'; // Gray
	}
}

function onDashboardStateChange(newState)
{
	var color;
	color = getStatusColor(newState);
	updateFaviconColor(color);
}

function isInt(value)
{
	if (isNaN(value))
	{
		return false;
	}
	var x = parseFloat(value);

	return (x | 0) === x;
}

function lowerCase(item)
{
	if (item == "")
	{
		return item;
	}

	if (item == undefined)
	{
		return item;
	}

	if (isInt(item))
	{
		return item;
	}

	return item.toLowerCase();
}

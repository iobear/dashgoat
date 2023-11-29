function unix2date(unixtime)
{
	const datetime = new Date(unixtime*1000);

	const year = datetime.getFullYear();
	const month = "0" + (datetime.getMonth() + 1);
	const day = "0" + datetime.getDate();
	const hour = datetime.getHours();
	const min = "0" + datetime.getMinutes();
	const sec = "0" + datetime.getSeconds();
	const result = `${year}-${month.substr(-2)}-${day.substr(-2)} ${hour}:${min.substr(-2)}:${sec.substr(-2)}`;

	return result;
}

function timeDiff(unixtime) {
	const endDate = new Date();
	const startDate = new Date(unixtime * 1000);

	var diff = endDate.getTime() - startDate.getTime();
	var days = Math.floor(diff / (1000 * 60 * 60 * 24));
	var hours = Math.floor(diff / (1000 * 60 * 60));
	var minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

	if (days >= 1) {
		hours = Math.ceil((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
		return (days <= 9 ? "0" : "") + days + "D" + (hours <= 9 ? "0" : "") + hours + "H";
	} else {
		return (hours <= 9 ? "0" : "") + hours + "H" + (minutes <= 9 ? "0" : "") + minutes;
	}
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

function onDashboardStateChange(newState) {
	var color;
	switch(newState) {
		case 'critical':
			color = '#ef1b11'; // Red
			break;
		case 'error':
			color = '#ef8511'; // Orange
			break;
		case 'warning':
			color = '#e2dc0c'; // Yellow
			break;
		case 'info':
			color = '#14a7c9'; // Blue
			break;
		case 'ok':
			color = '#33ff00'; // Green
			break;
		default:
			color = '#808080'; // Gray
	}
	updateFaviconColor(color);
}

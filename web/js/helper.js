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

/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

const default_items = ['host', 'service', 'status', 'message', 'change', 'probe'];
const footer_items = ['dashname','','','','','dashversion'];
let dash_items = [];
let oldrows = 0;
let printList = [];


function createHeader()
{
	let cell = document.createElement("thead");
	dashtable.appendChild(cell).id = "dashheader";

	for (let item of dash_items)
	{
		let cell = document.createElement("th");
		cell.setAttribute("onclick", "setSortBy(this.innerText.toLowerCase())");

		dashheader.appendChild(cell).classList.add(item)
		dashheader.appendChild(cell).innerText = item;
	}

}


function createTableBody()
{
	let cell = document.createElement("tbody");
	dashtable.appendChild(cell).id = "dashbody";

}


function createFooter()
{
	let cell = document.createElement("tfoot");
	dashtable.appendChild(cell).id = "dashfooter";

	for (let item of footer_items)
	{
		let cell2 = document.createElement("th");

		if (item)
		{
			dashfooter.appendChild(cell2).id = item;
		}
		dashtable.appendChild(cell2).classList = "dashfooter";
		dashfooter.appendChild(cell2).innerText = item;
	}
}


function removeList()
{
	let prelist = document.getElementById("dashtable");

	if (prelist.childElementCount > 0)
	{
		document.getElementById("dashheader").remove();
		document.getElementById("dashbody").remove();
	}

}


function createRows(startrow)
{
	for (let c = startrow; c < rows; c++)
	{
		rowid = `row${c}`;
		let cell = document.createElement("tr");
		dashbody.appendChild(cell).id = rowid;

		createColumns(rowid);
		createHistoryRows(rowid);
	}
}

function createHistoryRows(rowid)
{
	history_rowid = rowid + 'h';
	cell = document.createElement("tr");
	cell.setAttribute('style', 'display: none;');
	dashbody.appendChild(cell).id = history_rowid;
}


function removeRows(startrow)
{
	for (let c = startrow; c < oldrows; c++)
	{
		rowid = `row${c}`;
		document.getElementById(rowid).remove();
	}
}


function createColumns(rowid)
{
	let rowToUpdate = document.getElementById(rowid);

	for (let item of dash_items)
	{
		let cell = document.createElement("td");

		rowToUpdate.appendChild(cell).classList.add(item)
		rowToUpdate.appendChild(cell).id = `${rowid}${item}`;
	}
}


function updateRows(refresh = true)
{
	const refreshSec = 4;
	let count = 0;

	for (let service of print_list)
	{
		for (let item of dash_items)
		{
			let row = "row" + count;
			let toUpdate = row + item;

			let container = document.getElementById(toUpdate);

			if (lowerCase(container.innerText) != lowerCase(service[item]))
			{
				container.innerText = service[item];
			}

			if (item == 'change' || item == 'probe')
			{
				if (service[item])
				{
					container.innerText = timeDiff(service[item]);
					container.onclick = function() {
						displayHistory(service, row);
						console.log("Clicked on:", row);
					};
				}
			}

			if (item == 'status')
			{
				if (count == 0)
				{
					onDashboardStateChange(service[item]);
				}
				if (service[item] != 'ok')
				{
					changeRowColor("row" + count, service[item]);
				}
				else
				{
					changeRowColor("row" + count, "");
				}
			}
		}
		count++;
	}

	if (refresh) {
		setTimeout(() => {  askAPI(); }, refreshSec * 1000);
	}
}


function changeRowColor(rowid, status)
{
	document.getElementById(rowid).className = status;
}


function changeRowAmount(rows)
{
	const rowdiff = rows - oldrows;

	if (rowdiff > 0)
	{ //add rows
		createRows(oldrows);
		oldrows = rows;
		return true
	}

	if (rowdiff < 0)
	{ //remove rows
		removeRows(rows);
		oldrows = rows;
		return true
	}
}


function selectHeader()
{
	dash_items = default_items;
	createHeader();
}


function updateVersion(data)
{
	document.title = data['DashName'];
	document.getElementById("dashname").textContent = data['DashName'];
	document.getElementById("dashversion").textContent = data['DashGoatVersion'];
}


function showDash()
{
	oldrows = 0;
	removeList();
	createTableBody();
	selectHeader();
	askAPI();
	createFooter();
	askHealth();
}

showDash();

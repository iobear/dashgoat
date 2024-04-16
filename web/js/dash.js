/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

const default_header = ['host', 'service', 'status', 'message', 'change', 'probe'];
const table_footer = ['dashname','','','','','dashversion'];
let table_header = [];
let oldrows = 0;
let printList = [];
let backend_app_data = {};
backend_app_data['metrics_history'] = true;

function createHeader()
{
	let cell = document.createElement("div");
	cell.classList.add('divTableHeader');
	cell.classList.add('divTableRow');
	dashtable.appendChild(cell).id = "dashheader";

	for (let item of table_header)
	{
		let cell = document.createElement("div");

		cell.setAttribute("onclick", "setSortBy(this.innerText)");
		dashheader.appendChild(cell).classList.add('cell'+item)
		dashheader.appendChild(cell).innerText = item;
	}

}


function createTableBody()
{
	let cell = document.createElement("div");

	cell.classList.add('divTableBody');
	dashtable.appendChild(cell).id = "dashbody";
}


function createFooter()
{
	let cell = document.createElement("div");
	cell.classList.add('divTableFoot');
	cell.classList.add('divTableRow');
	dashtable.appendChild(cell).id = "dashfooter";

	let count = 0;
	for (let item of table_footer)
	{
		let cell2 = document.createElement("div");

		if (item)
		{
			dashfooter.appendChild(cell2).id = item;
		}
		dashtable.appendChild(cell2).classList = "dashfooter";
		cell2.classList.add('cell'+table_header[count]);
		dashfooter.appendChild(cell2).innerText = item;
		count++;
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
		let cell = document.createElement("div");
		cell.classList.add('divTableRow');
		dashbody.appendChild(cell).id = rowid;

		createColumns(rowid);
		if (backend_app_data['metrics_history'] == true)
		{
			createHistoryRows(rowid);
		}
		else
		{
			console.log('no history');
		}
	}
}

function createHistoryRows(rowid)
{
	history_rowid = rowid + 'h';
	cell = document.createElement("div");

	cell.setAttribute('style', 'display: none;');
	cell.classList.add('timeline');
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


function createColumns(rowid) {
    let rowToUpdate = document.getElementById(rowid);

	for (let item of table_header)
	{
        let tdDiv = document.createElement('div');

		tdDiv.id = `${rowid}${item}`;
		tdDiv.classList.add('cell'+item);

        rowToUpdate.appendChild(tdDiv);
    }
}


function updateRows(refresh = true)
{
	const refreshSec = 4;
	let count = 0;

	for (let service of print_list)
	{
		for (let item of table_header)
		{
			let row = "row" + count;
			let toUpdate = row;
			let container = document.getElementById(toUpdate);

			// Target the div inside the container (td)
			let divInsideContainer = container.querySelector('#'+row+item);

			if (!divInsideContainer)
			{
				console.log('Empty container');
				continue;
			}

			if (lowerCase(divInsideContainer.innerText) != lowerCase(service[item]))
			{
				divInsideContainer.innerText = service[item];
			}

			if (item == 'change' || item == 'probe')
			{
				if (service[item])
				{
					divInsideContainer.innerText = timeDiff(service[item]);
					divInsideContainer.onclick = function()
					{
						displayHistory(service, row);
					};
				}
			}

			if (item == 'status')
			{
				if (count == 0)
				{
					onDashboardStateChange(service[item]);
				}
				changeRowColor("row" + count, service[item]);
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
	document.getElementById(rowid).className = ('divTableRow ' + status);

	if (backend_app_data['metrics_history'] == true) {
		document.getElementById(rowid+'h').className = ('timeline ' + status);
	}
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
	table_header = default_header;
	createHeader();
}


function insertBackendAppData(data)
{
	backend_app_data['metrics_history'] = data['Prometheus'];
	document.title = data['DashName'];
	document.getElementById("dashname").textContent = data['DashName'];
	document.getElementById("dashversion").textContent = data['DashGoatVersion'];
}


function showDash()
{
	oldrows = 0;
	removeList();
	selectHeader();
	createTableBody();
	askAPI();
	createFooter();
	askHealth();
}

showDash();


const container = document.getElementById("container");
const default_items = ['host', 'service', 'status', 'message', 'change', 'probe'];
let dash_items = [];
let oldrows = 0;
let printList = [];
let sort_by = '';
let current_job = '';


function createHeader()
{
	let cell = document.createElement("thead");
	dashtable.appendChild(cell).id = "dashheader";

	for (let item of dash_items)
	{
		let cell = document.createElement("th");
		cell.setAttribute("onclick", "updateRows(this.innerText.toLowerCase())");

		dashheader.appendChild(cell).classList.add(item)
		dashheader.appendChild(cell).innerText = item;
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
	let cell = document.createElement("tbody");
	dashtable.appendChild(cell).id = "dashbody";

	for (let c = startrow; c < rows; c++)
	{
		rowid = `row${c}`;
		let cell = document.createElement("tr");
		dashbody.appendChild(cell).id = rowid;

		createColumns(rowid);
	}

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


function updateRows(sort = '')
{
	if (sort)
	{
		sort_by = sort;
	}

	let rowcount = 0;
	count = 0;

	sortPrintList();

	for (let service of print_list)
	{
		for (let item of dash_items)
		{
			let toUpdate = "row" + count + item;

			let container = document.getElementById(toUpdate);
			container.innerText = "";
			container.innerText = service[item];

			if (item == 'change' || item == 'probe') {

				if (service[item]) {
					container.innerText = timeDiff(service[item]);
				}

			} else {
				container.innerText = service[item];
			}

			if (item == 'status')
			{
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
	setTimeout(() => {  askAPI(); }, 4000);
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


function selectHeader(job)
{
	if (job == 'default')
	{
		dash_items = default_items;
	}

	createHeader();
}


function prepareData(data)
{
	print_list = [];
	keys = Object.keys( data );
	rows = keys.length;

	//If empty API
	if ((rows == 2) && keys.includes('status'))
	{
		alert(data['message']);
		return true;
	}

	for (var value of keys)
	{
		if (value)
		{
			print_list.push(data[value]);
		}
		else
		{
			rows = rows - 1;
		}
	}

	changeRowAmount(rows);
	updateRows();
}


function sortPrintList()
{
	if (sort_by == '')
	{
		return '';
	}

	var item_list = print_list.map(function (item) {
		return item[sort_by];
	});

	org_list = item_list.slice(); //copy list

	if (Number.isInteger(org_list[0]))
	{
		item_list.sort(function(a, b){return a-b});
	}
	else
	{
		item_list.sort();
	}

	let new_print_list = [];
	for (let sort_key in item_list)
	{
		let index_to_find = org_list.indexOf(item_list[sort_key]);
		org_list[index_to_find] = '';
		new_print_list.push(print_list[index_to_find]);
	}

	print_list = new_print_list;
	new_print_list = []; //clean list
	org_list = []; //clean list
}


function show(job, sort='')
{
	current_job = job;
	if (sort != '')
	{
		sort_by = sort;
	}

	oldrows = 0;
	removeList();
	selectHeader(job);
	askAPI();
}

show('default','status')
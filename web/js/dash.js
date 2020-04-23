
const container = document.getElementById("container");
const dashItems = ['host', 'service', 'status', 'message', 'change', 'seen'];
let oldrows = 0;
let printList = [];


function createHeader() {
	let cell = document.createElement("thead");
	dashtable.appendChild(cell).id = "dashheader";


	for (let item of dashItems) {

		let cell = document.createElement("th");

		dashheader.appendChild(cell).classList.add(item)
		dashheader.appendChild(cell).innerText = item;
	}

}


function createRows(startrow) {
	let cell = document.createElement("tbody");
	dashtable.appendChild(cell).id = "dashbody";

	for (let c = startrow; c < rows; c++) {

		rowid = `row${c}`;
		let cell = document.createElement("tr");

		dashbody.appendChild(cell).id = rowid;
		createColumns(rowid);
	}

}


function removeRows(startrow) {

	for (let c = startrow; c < oldrows; c++) {

		rowid = `row${c}`;
		document.getElementById(rowid).remove();
	}

}


function createColumns(rowid) {

	let rowToUpdate = document.getElementById(rowid);

	for (let item of dashItems) {

		let cell = document.createElement("td");

		rowToUpdate.appendChild(cell).classList.add(item)
		rowToUpdate.appendChild(cell).id = `${rowid}${item}`;
	}
}


function updateRows() {

	let rowcount = 0;
	count = 0;
	
	for (let service of printList) {

		for (let item of dashItems) {
			let toUpdate = "row" + count + item;

			let container = document.getElementById(toUpdate);
			container.innerText = "";

			if (item == 'change' || item == 'seen') {

				if (service[item]) {
					container.innerText = timeDiff(service[item]);
				}

			} else {
				container.innerText = service[item];

			}

			if (item == 'status') {

				if (service[item] != 'ok') {
					changeRowColor("row" + count, service[item]);
				} else {
					changeRowColor("row" + count, "");
				}

			}
		}

		count++;
	}
	setTimeout(() => {  askAPI(url); }, 4000);
}


function changeRowColor(rowid, status) {

	document.getElementById(rowid).className = status;
	
}

function changeRowAmount(rows) {
	const rowdiff = rows - oldrows;

	if (rowdiff > 0) { //add rows
		createRows(oldrows);
		oldrows = rows;
		return true
	}

	if (rowdiff < 0) { //remove rows
		removeRows(rows);
		oldrows = rows;
		return true
	}
}


function prepareData(data) {
	printList = [];
	var printData = {};
	keys = Object.keys( data );
	rows = keys.length;

	//create status objects
	for (var value of keys) {
		if (value == "") {
			value = "empty-key";
		}

		var status = data[value].status;
		printData[status] = {};

		//sort arr
		if (status == 'ok') {
			printList.push(data[value]);

		} else { //non status ok on top..
			printList.unshift(data[value]);

		}
	}

	//check for error status
	for (var value of keys) {

		var status = data[value].status;
		printData[status][value] = data[value];

	}

	changeRowAmount(rows);
	updateRows();
}

createHeader();


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


function createRows(rows, classname) {
	let cell = document.createElement("tbody");
	dashtable.appendChild(cell).id = "dashbody";

	for (let c = 0; c < rows; c++) {

		rowid = `row${c}`;
		let cell = document.createElement("tr");

		dashbody.appendChild(cell).id = rowid;
		createColumns(rowid);
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
	//console.log(printList);

	var rowcount = 0;
	count = 0;
	
	for (var service of printList) {

		for (let item of dashItems) {
			let toUpdate = "row" + count + item;

			let container = document.getElementById(toUpdate);

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


function prepareData(data) {
	printList = [];
	var printData = {};
	keys = Object.keys( data );
	rows = keys.length;

	//create status objects
	for (var value of keys) {

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
	
	if (oldrows == 0){
		createRows(rows);
		oldrows = rows;
	}	
	updateRows();
}

createHeader();

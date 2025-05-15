/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

let sort_by = 'change';
let sort_reverse = true;
let status_maps = {};
let status_arr = ['critical','error','warning','info','ok'];
let print_list = [];


function initStatusItem(value)
{
	status_maps[value] = [];
}


function setSortBy(sort)
{
	if (sort == sort_by)
	{
		sort_reverse = !sort_reverse;
	} else {
		sort_by = sort;
	}

	updateRows();
}


function isStatus(value)
{
	if (!status_arr.includes(map_key))
	{
		status_arr.unshift(map_key);
		status_maps[map_key] = [];
	}
}


function sortByStatus()
{

	status_count = status_arr.length;
	for (let si = 0; si < status_count; si++)
	{ //iter status'

		list_to_add = status_maps[status_arr[si]];
		event_count = list_to_add.length;

		if (event_count > 1)
		{
			list_to_add = sortListByValue(list_to_add, event_count);
		}

		for (let i = 0; i < event_count; i++)
		{
			print_list.push(list_to_add[i]);
		}

	}
}


function prepareData(data, refresh = true)
{
	status_arr.forEach(initStatusItem); //init empty status_map
	print_list = [];
	keys = Object.keys( data );
	rows = keys.length;

	if ((rows == 2) && keys.includes('status')) //if empty API
	{
		alert(data['message']);
		return true;
	}

	for (var value of keys)
	{
		if (value)
		{
			map_key = data[value].status
			isStatus(map_key);
			status_maps[map_key].push(data[value]);
		}
		else
		{
			rows = rows - 1;
		}
	}

	sortByStatus();
	changeRowAmount(rows);
	updateRows(refresh);
}


function sortListByValue(list_to_sort, list_count)
{
	var item_list = list_to_sort.map(function (item)
	{
		return item[sort_by];
	});

	org_list = item_list.slice(); //copy list

	if (Number.isInteger(org_list[0]))
	{
		tmp_zero_list = [];
		tmp_sort_list = [];

		for (let i = 0; i < list_count; i++)
		{
			if (item_list[i] == 0)
			{
				tmp_zero_list.push(item_list[i]);
			}
			else
			{
				tmp_sort_list.push(item_list[i]);
			}
		}

		if (tmp_sort_list.length > 1)
		{
			if (sort_reverse) {
				tmp_sort_list = tmp_sort_list.sort((a, b) => b - a);
			} else {
				tmp_sort_list = tmp_sort_list.sort((a, b) => a - b);
			}
		}

		var item_list = tmp_sort_list.concat(tmp_zero_list);

	}
	else
	{
		if (sort_reverse) {
			item_list.reverse();
		} else {
			item_list.sort();
		}

	}

	let new_list = [];
	for (let sort_key in item_list)
	{
		let index_to_find = org_list.indexOf(item_list[sort_key]);
		org_list[index_to_find] = '';
		new_list.push(list_to_sort[index_to_find]);
	}

	return new_list;
}

function getSearchQuery() {

	const urlParams = new URLSearchParams(window.location.search);
	let searchParam = urlParams.get('search');
	if (searchParam === null) {
		console.log("No search parameter found");
		return '';
	} else {
		searchParam = decodeURIComponent(searchParam); // Decode the URI component to handle special characters
	}

	searchParam = searchParam.replace(/ /g, '+'); // Replace spaces with plus signs for URL encoding

	// allow only  a-zA-Z0-9 and +,-_, delete non ascii characters

	let new_search_param = '';
	for (let i = 0; i < searchParam.length; i++) {
		if ((searchParam[i] >= 'A' && searchParam[i] <= 'Z') || (searchParam[i] >= 'a' && searchParam[i] <= 'z') || (searchParam[i] >= '0' && searchParam[i] <= '9') || searchParam[i] === '+' || searchParam[i] === '-' || searchParam[i] === '_' || searchParam[i] === ',') {
			new_search_param += searchParam[i];
		}
	}

	return new_search_param;

}
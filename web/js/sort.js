let sort_by = 'change';
let sort_reverse = true;
let status_maps = {};
let status_arr = ['critical','error','warning','info','ok'];


function iterStatus(value)
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


function prepareData(data)
{
	status_arr.forEach(iterStatus); //init empty status_map
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
	updateRows();
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

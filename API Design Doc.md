API Design

| REST verb | URI                          | Note                     |
| --------- | ---------------------------- | ------------------------ |
| GET       | /login                       | 进入登陆页面               |
| POST      | /session                     | 登陆验证                  |
| DELETE    | /session                     | 注销登陆                  |
| GET       | /report/catalog              | 获取所有巡检结果的文件名     |
| GET       | /report/id/{id}              | 获取指定的某次结果          |
| GET       | /report/latest               | 获取最后一次巡检结果        |
| GET       | /report/meta                 | 获取巡检结果元信息          |
| POST      | /report                      | 执行一次巡检              |
| GET       | /report/download/all         | 下载所有巡检报告           |
| GET       | /report/download/{id}        | 下载指定的巡检报告          |
| GET       | /script/local                | 查看所有远程仓库脚本        |
| GET       | /script/remote               | 查看所有远程仓库脚本        |
| GET       | /script/remote/readme/:name  | 查看指定脚本名的介绍        |
| POST      | /script/remote/download/:name| 下载远程仓库的脚本到本地     |

### 登陆验证

`POST /session`

- Request Form Data:

```json
{
  "username": "tidb",
  "password": "password"
}
```

- Response Body - 200

```json
{
  "token": "xxx"
}
```

- Response Body - 401

```json
{
  "error": "login failed"
}
```

### 获取巡检历史记录

`GET /report/catalog`

- Request Body:

```json
{
  "start": 1,
  "length": 10
}
```

Response Body:

```json
{
	"draw": 1, //当前页码
    "recordsFiltered": 2, //被filter后的行数，可以忽略
	"recordsTotal": 3, //总行数
	"data": [
		{
			"id": "20211228092830",
			"starttime":"2021-12-28 09:28:03",
			"duration": 30,
			"success": 29,
			"warning": 1,
			"total": 30
		},
		{
			"id": "20211229092830",
			"starttime":"2021-12-28 09:28:03",
			"duration": 25,
			"success": 28,
			"warning": 2,
			"total": 30
		},
		...
	]
}
```

### 获取指定的某次结果/ 获取最后一次巡检结果

`GET /report/{id}` `GET /report/last`

- Response Body - 200:

```json
{
  "total": 2,
  "id": "20211229092830",
  "data": [
    {
      "check_class": "集群",
      "check_name": "存活的TiDB数量",
      "check_item": "TiDB节点数量",
      "check_status": "正常",
      "check_value": "3",
      "threshold": "等于3"
    },
    {
      "check_class": "集群",
      "check_name": "存活的TiKV数量",
      "check_item": "TiKV节点数量",
      "check_status": "异常_新增",
      "check_value": "4",
      "threshold": "等于5"
    }
  ]
}
```

- Response Body - 400:

```json
{
  "id": "18881223083050",
  "error": "file not found"
}
```

### 获取巡检结果元信息

`GET /report/meta`

- Response Body:

```json
{
	"cluster_status": {
			"total_execution_count": 500,
			"total_checked_item": 1500,
			"last_execution": "20211228092830",
			"cluster_health": 88
		},
	"recent_warnings_total_check": 10,
	"recent_warnings": [
		{
			"check_time": "20211228092830",
			"warning_count": 4
		},
		{
			"check_name": "20211228092829",
			"warning_count": 2
		},
        ...
	],
	"normal_week_total": 30,
	"normal_week": [
		{
			"check_name": "存活的TiDB数量",
			"check_item": "TiDB节点数量",
			"count": 30
		},
		{
			"check_name": "存活的TiKV数量",
			"check_item": "TiKV节点数量",
			"count": 30
		},
		...
	],
	"normal_month_total": 30,
	"normal_month": [
		{
			"check_name": "存活的TiDB数量",
			"check_item": "TiDB节点数量",
			"count": 120
		},
		{
			"check_name": "存活的TiKV数量",
			"check_item": "TiKV节点数量",
			"count": 120
		},
		...
	],
	"normal_year_total": 30,
	"normal_year": [
		{
			"check_name": "存活的TiDB数量",
			"check_item": "TiDB节点数量",
			"count": 1440
		},
		{
			"check_name": "存活的TiKV数量",
			"check_item": "TiKV节点数量",
			"count": 1440
		},
		...
	]
}
```

### 执行一次巡检

`POST /report`

- Response Body (prior to last):

```json
{
  "finished": false,
  "total": 0,
  "check_class": "集群",
  "check_name": "存活的TiDB数量",
  "check_item": "TiDB节点数",
  "check_result": "正常",
  "check_value": 5,
  "check_threshold": "等于5",
  "check_time": 20211221063030
}
```

- Response Body (last response):

```json
{
  "finished": true,
  "check_class": "集群",
  "check_name": "存活的TiDB数量",
  "check_item": "TiDB节点数",
  "check_result": "正常",
  "check_value": 5,
  "check_threshold": "等于5",
  "check_time": 20211221063030
}
```

### 查看本地脚本列表

`GET /script/local`

- Request Query Param
```
  // The start default value is 0
  // The length default value is 10
  ?start=0&length=2
```

- Response Body

```json
{
  "total": 2, 
  "script_list": [
    {
      "name": "alive_pd_number"
    },{
      "name": "alive_tidb_number"
    }
  ]
}
```

### 查看所有远程仓库脚本

`GET /script/remote`

- Request Query Param
```
  // The start default value is 0
  // The length default value is 10
  ?start=0&length=2
```
- 
- Response Body :

```json
{
  "total": 2,
  "script_list": [
    {
      "name": "alive_pd_number",
      "download": true
    },{
      "name": "alive_tidb_number",
      "download": false
    }
  ]
}
```

### 查看指定脚本名的介绍

`GET /script/remote/readme/:name`

- Response Body :

```json
{
    "readme": ""
}
```

### 下载远程仓库的脚本到本地

`POST /script/remote/download/:name`

- Response Body :

```json
{
  "status": "ok"
}
```
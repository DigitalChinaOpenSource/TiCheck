API Design

| REST verb | URI                   | Note                     |
| --------- | --------------------- | ------------------------ |
| GET       | /login                | 进入登陆页面             |
| POST      | /session              | 登陆验证                 |
| DELETE    | /session              | 注销登陆                 |
| GET       | /report/catalog       | 获取所有巡检结果的文件名 |
| GET       | /report/id/{id}       | 获取指定的某次结果       |
| GET       | /report/latest        | 获取最后一次巡检结果     |
| GET       | /report/meta          | 获取巡检结果元信息       |
| POST      | /report               | 执行一次巡检             |
| GET       | /report/download/all  | 下载所有巡检报告         |
| GET       | /report/download/{id} | 下载指定的巡检报告       |

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

### 获取最近巡检

`GET /report/catalog`

- Request Body:

```json
{
  "page": 1,
  "per_page": 3
}
```

Response Body:

```json
{
    "total_pages": 2,
	"data_total": 3,
	"data": [
		{
			"id": "20211228092830",
			"duration": 30,
			"success": 29,
			"warning": 1,
			"total": 30
		},
		{
			"id": "20211229092830",
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
	"recent_warnings_total_check": 4,
	"recent_warnings": [
		{
			"check_name": "存活的TiDB数量",
			"warning_count": 4
		},
		{
			"check_name": "存活的TiKV数量",
			"warning_count": 2
		},
        ...
	],
	"normal_week_total": 30,
	"normal_week": [
		{
			"check_name": "存活的TiDB数量",
			"count": 30
		},
		{
			"check_name": "存活的TiKV数量",
			"count": 30
		},
		...
	],
	"normal_month_total": 30,
	"normal_month": [
		{
			"check_name": "存活的TiDB数量",
			"count": 120
		},
		{
			"check_name": "存活的TiKV数量",
			"count": 120
		},
		...
	],
	"normal_year_total": 30,
	"normal_year": [
		{
			"check_name": "存活的TiDB数量",
			"count": 1440
		},
		{
			"check_name": "存活的TiKV数量",
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

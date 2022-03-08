API Design
| REST verb | URI                          | Note                     |
| --------- | ---------------------------- | ------------------------ |
| GET       | /login                       | 进入登陆页面               |
| POST      | /session                     | 登陆验证                  |
| POST      | /session/{user}              | 注销登陆                  |
| GET       | /session/{user}              | 获得用户的信息(集群列表)     |
| GET       | /cluster/{clusterID}         | 获得某个集群的详细信息       |
| POST      | /cluster                     | 添加一个集群               |
| POST      | /cluster/{clusterID}         | 修改单个集群的配置          |
| DELETE    | /cluster/{clusterID}         | 删除一个集群               |
| GET       | /cluster/user/{user}         | 获得某个用户管理的所有集群详细信息 |

| GET       | /cluster/report/{reportID}           | 获取指定的某次巡检结果       |
| POST      | /cluster/report/{clusterID}          | 执行一次巡检               |
| GET       | /cluster/report/all/{clusterID}      | 获取某个集群历史巡检记录     |
| GET       | /cluster/report/latest/{clusterID}   | 获取某个集群最后一次巡检结果  |
| GET       | /cluster/report/meta/{clusterID}     | 获取某个集群巡检结果元信息    |
| GET       | /cluster/report/download/{reportID}  | 下载指定的巡检报告          |

| GET       | /cluster/task                        | 查看任务列表               |
| GET       | /cluster/task/{taskID}               | 通过ID查看指定任务            |
| POST      | /cluster/task                        | 创建一个定时任务            |
| Delete    | /cluster/task/{taskID}               | 删除一个定时任务            |

| GET       | /script/cluster                      | 查看某个集群的脚本指标配置   |
| POST      | /script/status                       | 启动或者禁用指定脚本        | 
| POST      | /script/install                      | 将一个脚本加到当前集群      |
| GET       | /script/local                        | 查看所有本地仓库脚本        |
| POST      | /script/local                        | 上传一个脚本到本地仓库      |
| GET       | /script/remote                       | 查看所有远程仓库脚本        |
| POST      | /script/remote/download              | 下载远程仓库的脚本到本地     |

### 访问失败，统一返回

- Response Body - 401

```json
{
  "err": "err message"
}
```

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

### 注销登录
`POST /session/{userID}`

- Request Body
null

- Response Body - 200

```json
{
  "status": "OK"
}
```

### 获得用户的信息(集群列表)
`GET /session/{user}`

- Request Query Param
```
  // The start default value is 0
  // The length default value is 10
  ?start=0&length=2
```

— Response Body:

```json
{
  "user": "root",
  "start": 1,
  "length": 10,
  "total": 20,
  "cluster": [
    {
      "cluster_id": 1,
      "cluster_name": "cluster",
      "prometheus_url": "url",
      "grafana_url": "url",
      "description": "description",
      "create_time": "create_time",
      "last_check_time": "last_check_time",
      "tidb_number": 10,
      "tidb_down_number": 0,
      "tikv_number": 10,
      "tikv_down_number": 0,
      "pd_number": 10,
      "pd_down_number": 0
    },{ 
    
    }
  ]
}
```

### 获得某个集群的详细信息
`GET /cluster/{clusterID}`

- Response Body -200
```json
{
  "cluster_id": 1,
  "cluster_name": "cluster",
  "cluster_owner": "cluster_owner",
  "tidb_version": "tidb_version",
  "description": "description",
  "create_time": 1,
  "last_check_time": 1,
  "total_execution_count": 666,
  "total_checked_item": 1500,
  "cluster_health": 88
}
```

### 添加一个集群
`POST /cluster`

- Request Body
```json
{
  "cluster_name": "cluster_name",
  "prometheus_url": "prometheus_url",
  "tidb_user": "root",
  "tidb_pw": "root",
  "description": "description"
}
```

- Response Body -200
```json
{
  "cluster_id": 2
}
```

### 修改单个集群的配置
`POST /cluster/{clusterID}`

- Request Body
```json
{
  "cluster_name": "cluster_name",
  "prometheus_url": "prometheus_url",
  "tidb_user": "root",
  "tidb_pw": "root",
  "description": "description"
}
```

- Response Body -200
```json
{
  "status": "OK"
}
```

### 删除一个集群
`DELETE /cluster/{clusterID}`

- Response Body -200
```json
{
  "status": "OK"
}
```

### 获得某个用户管理的所有集群详细信息
目前与获取用户的信息（集群列表接口返回相同）

### 获取指定的某次巡检结果
`GET /cluster/report/{reportID}`

- Response Body -200

```json
{
  "report_id": "report_id",
  "cluster_name": "cluster_name",
  "check_time": "check_time",
  "normal_items": "normal_items",
  "warning_items": "warning_items",
  "total_items": "total_items",
  "duration": "duration",
  "check_items": [
    {
      "check_tag": "check_tag",
      "check_name": "check_name",
      "check_item": "check_item",
      "operator": "operator",
      "threshold": "threshold",
      "duration": "duration",
      "check_value": "check_value",
      "check_status": "check_status"
    },{
      
    }
  ]
}
```

### 执行一次巡检
`POST /cluster/report/{clusterID}`

- Response Body (prior to last):

```json
{
  "finished": false,
  "check_tag": "check_tag",
  "check_name": "check_name",
  "check_item": "check_item",
  "operator": "operator",
  "threshold": "threshold",
  "duration": "duration",
  "check_value": "check_value",
  "check_status": "check_status"
}
```

- Response Body (last response):

```json
{
  "finished": true,
  "check_tag": "check_tag",
  "check_name": "check_name",
  "check_item": "check_item",
  "operator": "operator",
  "threshold": "threshold",
  "duration": "duration",
  "check_value": "check_value",
  "check_status": "check_status"
}
```

### 获取某个集群历史巡检记录
`GET /cluster/report/all/{clusterID}`

- Request Query Param
```
  // The page_num default value is 1
  // The page_size default value is 10
  ?page_size=10&page_num=1
```

- Response Body -200

```json
{
  "page_num": 1,
  "page_size": 10,
  "total": 10,
  "data": [{
        "id": "20211228092830",
        "start_time":"2021-12-28 09:28:03",
        "duration": 30,
        "normal": 29,
        "warning": 1,
        "total": 30
      },{
      }
  ]
}
```

### 获取某个集群最后一次巡检结果
`GET /cluster/report/latest/{clusterID}`

- Response Body -200
```json
{
  "report_id": "report_id",
  "cluster_name": "cluster_name",
  "check_time": "check_time",
  "normal_items": "normal_items",
  "warning_items": "warning_items",
  "total_items": "total_items",
  "duration": "duration",
  "check_items": [
    {
      "check_tag": "check_tag",
      "check_name": "check_name",
      "check_item": "check_item",
      "operator": "operator",
      "threshold": "threshold",
      "duration": "duration",
      "check_value": "check_value",
      "check_status": "check_status"
    },{

    }
  ]
}
```

### 获取某个集群巡检结果元信息
`/cluster/report/meta/{clusterID}`

- Response Body -200
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
		}
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
		}
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
		}
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
		}
	]
}
```

### 下载指定的巡检报告
`GET /cluster/report/download/{reportID}`

### 查看任务列表
`GET /cluster/task/`

- Request Query Param
```
  // The cluster_id default value is nil
  // The start default value is 0
  // The length default value is 10
  ?cluster_id=1&start=0&length=2
```

- Response Body -200 
```json
{
  "start": 1,
  "length": 10,
  "total": 20,
  "task_list":[{
    "scheduler_id": 1,
    "cluster_id": 1,
    "title": "title",
    "cron_expression": "cron_expression",
    "is_active": 1,
    "creator": "creator",
    "create_time": 1,
    "run_count": 1
  },{
    
  }]
}
```

### 通过ID查看指定任务
`GET /cluster/task/{taskID}`

- Response Body -200
```json
{
  "scheduler_id": 1,
  "cluster_id": 1,
  "title": "title",
  "cron_expression": "cron_expression",
  "is_active": 1,
  "creator": "creator",
  "create_time": 1,
  "run_count": 1
}
```

### 创建一个定时任务
`POST /cluster/task`

- Request Body
```json
{
  "cluster_id": 1,
  "title": "title",
  "cron_expression": "cron_expression",
  "is_active": 1,
  "creator": "creator"
}
```

- Response Body -200
```json
{
  "status": "OK"
}
```

### 删除一个定时任务
`DELETE /cluster/task/{taskID}`

- Response Body -200
```json
{
  "status": "OK"
}
```

### 查看某个集群的脚本配置
`GET /script/cluster`

- Request Query Param
```
  // The cluster_id cannot be empty
  // The start default value is 0
  // The length default value is 10
  ?cluster_id=1&start=0&length=2
```

- Response 
```json
{
  "start": 1,
  "total": 10,
  "length": 10,
  "checklist_id": 1,
  "scripts": [{
    "script_id": 1,
    "script_name": "script_name",
    "script_tag": "script_tag",
    "is_enabled": "is_enabled",
    "description": "description",
    "threshold_operator": 1,
    "threshold_value": "threshold_value",
    "threshold_args": "threshold_args"
  },{
    
  }]
}
```

### 启动或者禁用指定脚本
`POST /script/status` 

```json
{
  "cluster_id": 1,
  "script_id": 1,
  "enable": true
}
```

- Response Body -200
```json
{
  "status": "OK"
}
```

### 将一个脚本加到当前集群
`POST /script/install`

```json
{
  "cluster_id": 1,
  "script_id": 1
}
```

- Response Body -200
```json
{
  "status": "OK"
}
```

### 查看所有本地仓库脚本
`GET /script/local`

- Request Query Param
```
  // The start default value is 0
  // The length default value is 10
  ?&start=0&length=2
```

- Response Body -200

```json
{
  "start": 1,
  "length": 10,
  "total": 20,
  "script_list": [{
    "script_id": 1,
    "script_name": "script_name",
    "script_tag": "script_tag",
    "description": "description",
    "threshold_operator": 1,
    "threshold_value": "threshold_value",
    "threshold_args": "threshold_args",
    "is_system": 1,
    "script_creator": "script_creator",
    "script_create_time": 1,
    "script_update_time": 1
  },{
    
  }]
}
```

### 上传一个脚本到本地仓库
`POST /script/local`

- Request Body

```json
{
  "script_name": "script_name",
  "script_tag": "script_tag",
  "description": "description",
  "script_creator": "script_creator"
}
```

- Response Body -200

```json
{
  "status": "ok"
}
```

### 查看所有远程仓库脚本
`GET /script/remote`

- Request Query Param

```
  // The start default value is 0
  // The length default value is 10
  ?&start=0&length=2
```

- Response Body -200

```json
{
  "start": 1,
  "length": 10,
  "total": 20,
  "script_list": [{
    "script_name": "script_name",
    "script_tag": "script_tag",
    "description": "description"
  },{

  }]
}
```

### 下载远程仓库的脚本到本地
`POST /cluster/script/remote/download`

- Response Body -200
```json
{
  "status": "ok"
}
```
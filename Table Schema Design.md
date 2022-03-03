# Table Schema Design

> ***注意事项：***
>
> -  所有时间类型字段默认使用`INTEGER`，也就是秒级时间戳，除非有特殊说明要用`TEXT`类型。
> -  xxx

### 系统用户信息

- Table Name: `user`
- Table Schema:

| column name   | data type | explanation             | nullable | other notes |
| ------------- | --------- | ----------------------- | -------- | ----------- |
| id    | INTEGER   | 用户id | 不可为空 | 主键    |
| user_name    | TEXT   | 用户名，登录名 | 不可为空 | 唯一值    |
| user_password    | TEXT   | 用户登录密码 | 不可为空 |     |
| full_name    | TEXT   | 用户真实姓名 | 不可为空 |     |
| email  | TEXT   |   | 不可为空 |  用户邮箱地址 |
| is_enabled | INTEGER   | 启用状态     | 不可为空 |  0-禁用，1-启用     |
| creator      | TEXT   | 创建人        | 不可为空 |   |
| create_time      | TEXT   | 创建时间     | 不可为空 |   |


### 脚本（组件/插件/巡检项）信息

- Table Name: `script`
- Table Schema:

| column name   | data type | explanation             | nullable | other notes |
| ------------- | --------- | ----------------------- | -------- | ----------- |
| id    | INTEGER   | 脚本id | 不可为空 | 主键    |
| script_name    | TEXT   | 脚本名称， 也就是这个巡检项的名称 | 不可为空 |     |
| file_name  | TEXT   | 关联的脚本文件  | 不可为空 | 文件名需要唯一，更新的时候根据文件名去查找            |
| tag | TEXT   | 脚本所属分类，即以前的检测类别                  | 不可为空 |             |
| description | TEXT  | 脚本简介     | 可为空 |   |
| Operator   | INTEGER   | 阈值检查方式，默认值 | 可为空 |  0-无，1-等于，2-大于，3-大于等于，4-小于，5-小于等于|
| Threshold      | TEXT   | 阈值检测值，默认值 | 可为空 |   |
| args    | TEXT   | 阈值检测参数，默认值 | 可为空 |     |
| is_system      | INTEGER   | 是否系统自带脚本        | 不可为空 |   |
| creator      | TEXT   | 脚本创建人        | 不可为空 | 系统脚本固定为system，自定义脚本为上传用户  |
| create_time      | TEXT   | 脚本创建时间     | 不可为空 |   |
| update_time      | TEXT   | 脚本更新时间     | 可为空 |   |


### 集群信息

- Table Name: `cluster`
- Table Schema:

| column name   | data type | explanation             | nullable | other notes |
| ------------- | --------- | ----------------------- | -------- | ----------- |
| id    | INTEGER   | 集群id | 不可为空 | 主键    |
| name    | TEXT   | 集群名称 | 不可为空 |     |
| prometheus_url  | TEXT   | TIDB集群的Prometheus地址 | 不可为空 |      |
| tidb_username | TEXT   | tidb登录用户  | 不可为空 |             |
| tidb_password   | TEXT   | tidb登录用户密码  | 不可为空 |             |
| description      | TEXT   | 集群描述        | 可为空 |   |
| create_time    | TEXT   | 创建时间 | 不可为空 | 精确到秒    |
| owner   | TEXT   | 集群owner，默认为创建人    | 不可为空 |   |
| tidb_version | TEXT   | tidb版本      | 不可为空 |   |
| dashboard_url | TEXT   | dashboard地址       | 可为空 |   |
| grafana_url | TEXT   | grafana地址         | 可为空 |   |
| last_check_time | TEXT   | 最后巡检时间      | 可为空 |  |
| cluster_health | INTEGER   | 集群健康度      | 可为空 |  |

### 集群巡检项

- Table Name: `cluster_checklist`
- Table Schema:

| column name   | data type | explanation             | nullable | other notes |
| ------------- | --------- | ----------------------- | -------- | ----------- |
| checklist_id    | INTEGER   | id | 不可为空 | 主键    |
| cluster_id    | INTEGER   | 集群id | 不可为空 |     |
| script_id    | INTEGER   | 脚本id | 不可为空 |     |
| is_enabled    | INTEGER   | 启用状态 | 不可为空 |  0-未启用，1-已启用   |
| threshold_operator   | INTEGER   | 阈值检查方式  | 可为空 |  0-无，1-等于，2-大于，3-大于等于，4-小于，5-小于等于|
| threshold_value      | TEXT   | 阈值检测值    | 可为空 |   |
| threshold_args    | TEXT   | 阈值检测参数 | 可为空 |     |

### 集群巡检定时器

- Table Name: `scheduler`
- Table Schema:

| column name   | data type | explanation             | nullable | other notes |
| ------------- | --------- | ----------------------- | -------- | ----------- |
| id    | INTEGER   | id | 不可为空 | 主键    |
| cluster_id    | INTEGER   | 集群id | 不可为空 |     |
| name    | TEXT   | 定时器名称 | 不可为空 |     |
| cron_expression    | TEXT   | cron字符串 | 不可为空 |     |
| is_enabled   | INTEGER   | 激活状态  | 不可为空 |  0-未激活，1-已激活 |
| creator    | TEXT   | 定时器创建人 | 不可为空 |     |
| create_time    | TEXT   | 创建时间 | 不可为空 | 精确到秒    |
| run_count    | INTEGER   | 巡检次数 | 不可为空 |     |

### 集群巡检记录

- Table Name: `check_history`
- Table Schema:

| column name   | data type | explanation             | nullable | other notes |
| ------------- | --------- | ----------------------- | -------- | ----------- |
| id    | INTEGER   | 巡检开始时间戳，作为 id | 不可为空 | 精确到秒 |
| check_time    | TEXT   | 巡检开始时间戳，作为 id | 不可为空 | 精确到秒 |
| cluster_id    | INTEGER   | 集群id | 不可为空 |     |
| scheduler_id  | INTEGER   | 定时器id | 可为空 | 如果是定时器触发就关联定时器id，手动运行为空  |
| normal_items  | INTEGER   | 正常项                  | 不可为空 |             |
| warning_items | INTEGER   | 告警项                  | 不可为空 |             |
| total_items   | INTEGER   | 总检查项                | 不可为空 |             |
| duration      | TEXT   | 本次巡检累计耗时   | 不可为空 | 单位为毫秒  |

### 集群巡检结果

- Table Name: `check_data`
- Table Schema:

| column name  | data type | explanation          | nullable                     | other notes                           |
| ------------ | --------- | -------------------- | ---------------------------- | ------------------------------------- |
| id           | INTEGER   | 自增主键             | 不可为空                     |                                       |
| history_id   | INTEGER   | 巡检记录id         | 不可为空                     | 关联到check_history的id                   |
| cluster_id    | INTEGER   | 集群id | 不可为空 |     |
| check_tag    | TEXT      | 检查类别             | 不可为空                     |                                       |
| check_name   | TEXT      | 检查项目             | 不可为空                     |                                       |
| operator     | TEXT      | 比较方式             | 不可为空                     | 参考tck_addons.threshold_operator |
| threshold    | REAL      | 阈值                 |  |        |
| duration     | TEXT   | 该项巡检用了多少时间   | 不可为空                     | 单位为毫秒   |
| check_item   | TEXT      | 检查指标             | 不可为空                     |                                       |
| check_value  | REAL      | 指标数值             | 可为空，代表脚本无输出       |                                       |
| check_status | TEXT      | 检查结果             | 不可为空                     | `正常`/`异常_已有`/`异常_新增`        |

逻辑包含关系为:
check_time >= check_class >= check_name = operator = threshold = duration >= check_item = check_status = check_value
即为 `次`级 >= `类`级 >= `项`级 >= `指标`级
需要注意的是单个检测脚本为`项`级

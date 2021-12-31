# Table Schema Design

巡检数据储存

- Table Name: `check_history`
- Table Schema:

| column name   | data type | explanation             | nullable | other notes |
| ------------- | --------- | ----------------------- | -------- | ----------- |
| check_time    | INTEGER   | 巡检开始时间戳，作为 id | 不可为空 | 精确到秒    |
| normal_items  | INTEGER   | 正常项                  | 不可为空 |             |
| warning_items | INTEGER   | 告警项                  | 不可为空 |             |
| total_items   | INTEGER   | 总检查项                | 不可为空 |             |
| duration      | INTEGER   | 本次巡检累计耗时        | 不可为空 | 单位为毫秒  |

- Table Name: `check_data`
- Table Schema:

| column name  | data type | explanation          | nullable                     | other notes                           |
| ------------ | --------- | -------------------- | ---------------------------- | ------------------------------------- |
| id           | INTEGER   | 自增主键             | 不可为空                     |                                       |
| check_time   | INTEGER   | 巡检开始时间         | 不可为空                     | 精确到秒                              |
| check_class  | TEXT      | 检查类别             | 不可为空                     |                                       |
| check_name   | TEXT      | 检查项目             | 不可为空                     |                                       |
| operator     | TEXT      | 比较方式             | 不可为空                     | `等于`/`大于等于`/`小于等于`/`无数据` |
| threshold    | REAL      | 阈值                 | 在比较方式为`无数据`时可为空 |                                       |
| duration     | INTEGER   | 该项巡检用了多少时间 | 不可为空                     | 单位为毫秒                            |
| check_item   | TEXT      | 检查指标             | 不可为空                     |                                       |
| check_value  | REAL      | 检查结果             | 可为空，代表脚本无输出       |                                       |
| check_status | TEXT      | 检查结果             | 不可为空                     | `正常`/`异常_已有`/`异常_新增`        |

逻辑包含关系为:
check_time >= check_class >= check_name = operator = threshold = duration >= check_item = check_status = check_value
即为 `次`级 >= `类`级 >= `项`级 >= `指标`级
需要注意的是单个检测脚本为`项`级

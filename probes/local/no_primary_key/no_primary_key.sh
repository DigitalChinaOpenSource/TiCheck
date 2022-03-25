#!/bin/bash

BASE_PATH=$1
MYSQL_LOGIN_PATH=$2
PROMETHEUS_ADDRESSES=$3

# 检查表是否没有主键
# 打印出所有没有主键副本的库名及表名

sql_command=\
"select table_schema,table_name from information_schema.tables
where (table_schema,table_name) not in(
    select distinct table_schema,table_name from information_schema.columns where COLUMN_KEY='PRI'
)
and table_schema not in (
    'METRICS_SCHEMA','mysql','INFORMATION_SCHEMA','PERFORMANCE_SCHEMA','test');"

result=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$sql_command" -ss)

IFS=$'\n'
for i in $result ; do
    echo "$i" | awk '{print $tck_result: $1"."$2=无主键}'
done

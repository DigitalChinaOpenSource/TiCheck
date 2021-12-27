#!/bin/bash

BASE_PATH=$1
MYSQL_LOGIN_PATH=$2
PROMETHEUS_ADDRESSES=$3

# 检查表是否没有TiFlash Replica
# 打印出所有没有TiFlash副本的库名及表名

sql_command=\
"select -- tr.TABLE_SCHEMA,tr.TABLE_NAME,
t.TABLE_SCHEMA,t.TABLE_NAME
from information_schema.\`TIFLASH_REPLICA\` as tr
right join information_schema.\`TABLES\` as t
on tr.TABLE_SCHEMA = t.TABLE_SCHEMA
and tr.TABLE_NAME = t.TABLE_NAME
where t.TABLE_SCHEMA not in ('INFORMATION_SCHEMA','METRICS_SCHEMA','PERFORMANCE_SCHEMA','mysql','test')
and tr.TABLE_NAME is null;"

result=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$sql_command" -ss)

IFS=$'\n'
for i in $result ; do
    echo "$i" | awk '{print $1"."$2}'
    echo "无TiFlash副本"
done

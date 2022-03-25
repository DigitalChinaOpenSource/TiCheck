#!/bin/bash

BASE_PATH=$1
MYSQL_LOGIN_PATH=$2
PROMETHEUS_ADDRESSES=$3

# 今天以内执行花费5s以上并且走的是tikv的SQL
# 会打印出JOB ID以及所用时间

print_sql=\
"select
JOB_ID, TIMESTAMPDIFF(SECOND,START_TIME,END_TIME) as cost_time
from information_schema.DDL_JOBS
where START_TIME between (now() - INTERVAL 1 DAY) and now()
having cost_time > 3600
order by cost_time desc
limit 100;"

to_print=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$print_sql" -ss)

IFS=$'\n'
for i in $to_print ; do
    echo "$i" | awk '{print "$tck_result: DDL JOB_ID:"$1"=花费时间"$2"秒"}'
done

# save_sql=\
# "select
# JOB_ID,DB_NAME,TABLE_NAME,JOB_TYPE,STATE,QUERY,START_TIME,END_TIME,TIMESTAMPDIFF(SECOND,START_TIME,END_TIME) as cost_time
# from information_schema.DDL_JOBS
# where START_TIME between (now() - INTERVAL 1 DAY) and now()
# having cost_time > 3600
# order by cost_time desc
# limit 100 \G"

# to_save=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$save_sql")
# if [ -n "$to_save" ]; then
#   echo "$to_save" > "$BASE_PATH"/err_report/long_ddl_job
# fi

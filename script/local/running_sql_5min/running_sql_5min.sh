#!/bin/bash

BASE_PATH=$1
MYSQL_LOGIN_PATH=$2
PROMETHEUS_ADDRESSES=$3

# 正在执行时间超过5min的SQL
# 打印出DIGEST与花费时间

print_sql=\
"select
\`DIGEST\`,\`TIME\`
from information_schema.cluster_processlist
where \`COMMAND\` <> 'Sleep'
and \`TIME\`> 300;"

to_print=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$print_sql" -ss)

IFS=$'\n'
for i in $to_print ; do
    echo "$i" | awk '{print "SQL ID:"$1"\n已经运行"$2"秒"}'
done

save_sql=\
"select *
from information_schema.cluster_processlist
where \`COMMAND\` <> 'Sleep'
and \`TIME\`> 300 \G"

to_save=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$save_sql")
if [ -n "$to_save" ]; then
  echo "$to_save" > "$BASE_PATH"/err_report/running_sql_5min
fi

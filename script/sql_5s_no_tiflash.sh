#!/bin/bash

BASE_PATH=$1
MYSQL_LOGIN_PATH=$2
PROMETHEUS_ADDRESSES=$3

# 今天以内执行话费5s以上并且走的是tikv的SQL
# 会打印出DIGEST以及花费时间

print_sql=\
"set @@sql_mode='';

select Digest,Query_time
from information_schema.CLUSTER_SLOW_QUERY
where \`User\` not in ('','tidb') and \`Query\` <> ''
and \`User\`<> 'root'
and Query_time > 5
and Time between (now() - INTERVAL 1 DAY) and now()
and Plan like '%[tikv]%'
group by Digest
order by Query_time desc
limit  100;"
to_print=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$print_sql" -ss)

IFS=$'\n'
for i in $to_print ; do
    echo "$i" | awk '{print "SQL ID:"$1"\n花费时间"$2"秒"}'
done

save_sql=\
"set @@sql_mode='';

select INSTANCE,Time,Txn_start_ts,Query_time,Digest,Plan,max(Query)
from information_schema.CLUSTER_SLOW_QUERY
where \`User\`not in ('','tidb') and \`Query\` <> ''
and \`User\`<> 'root'
and Query_time > 5
and Time between (now() - INTERVAL 1 DAY) and now()
and Plan like '%[tikv]%'
group by Digest
order by Query_time desc
limit  100 \G"

to_save=$(mysql --login-path="${MYSQL_LOGIN_PATH}" -e "$save_sql")
if [ -n "$to_save" ]; then
  echo "$to_save" > "$BASE_PATH"/err_report/sql_5s_no_tiflash
fi

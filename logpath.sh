#!/bin/bash

name=$1
user=$2
host=$3
port=$4
passwd=$5

expect -c "
spawn mysql_config_editor set --login-path=${name} --user=${user}  --host=${host} --port=${port} --password
expect -nocase \"Enter Password:\"
send \"${passwd}\r\"
interact
"

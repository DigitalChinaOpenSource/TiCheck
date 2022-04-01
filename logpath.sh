#!/bin/bash

##==============================================================================================
## load login-path variables
##==============================================================================================
LOGIN_PATH_NAME=$1
LOGIN_USER=$2
HOST=$3
PORT=$4
PASSWD=$5

##==============================================================================================
## remove MySQL login-path first
##==============================================================================================
mysql_config_editor remove --login-path="${LOGIN_PATH_NAME}"

##==============================================================================================
## generate MySQL login-path
##==============================================================================================
expect -c "
spawn mysql_config_editor set --login-path=${LOGIN_PATH_NAME} --user=${LOGIN_USER} --host=${HOST} --port=${PORT} --password
expect -nocase \"Enter Password:\"
send \"${PASSWD}\r\"
interact
"

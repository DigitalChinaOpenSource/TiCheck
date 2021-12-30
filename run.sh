#!/bin/bash

##==============================================================================================
## load run-time variables
##==============================================================================================

# load current path as base path
BASE_PATH=$(cd "$(dirname "$0")" || exit 1; pwd)

# check if run-time variable exist
if [ ! -d "${BASE_PATH}"/config ] || [ ! -f "${BASE_PATH}"/config/load_vars.sh ]; then
    echo "[run.sh] error! fail to load runtime variables"
    echo "[run.sh] failed to run!"
    exit 1
fi

source "${BASE_PATH}"/config/load_variables.sh

##==============================================================================================
## check script runner to be the same as specified in the config file
##==============================================================================================

if [[ $(whoami) != "${OS_USER}" ]]; then
  echo "[run.sh] error! current user is not ${OS_USER}!"
  echo "[run.sh] failed to run!"
  exit 1;
fi

##==============================================================================================
## executing all check scripts formats print
## this will also generate report in /report and error report in err_report/
##==============================================================================================
if python "${BASE_PATH}"/generate_report.py "$BASE_PATH" "$MYSQL_LOGIN_PATH" "$PROMETHEUS_ADDRESSES"
then
  echo "[run.sh] success!"
else
  echo "[run.sh] error! failed to run generate_report.py"
  exit 1
fi

#!/bin/bash

##==============================================================================================
## load run-time variables
##==============================================================================================

# load current path as base path
BASE_PATH=$(cd "$(dirname "$0")" || exit 1; pwd)

# check if run-time variable exist
if [ ! -d "${BASE_PATH}"/config ] || [ ! -f "${BASE_PATH}"/config/load_vars.sh ]; then
    echo -e "error!\t fail to find runtime variables, please check if the file 'config/load_variables.sh' exist in the current path."
    echo "failed to run!"
    exit 1
fi

source "${BASE_PATH}"/config/load_variables.sh

##==============================================================================================
## check script runner to be the same as specified in the config file
##==============================================================================================

if [[ $(whoami) != "${OS_USER}" ]]; then
  echo -e "error!\tcurrent user is not ${OS_USER}!"
  echo -e "\tplease use ${OS_USER} to run this script."
  echo "failed to run!"
  exit 1;
fi

##==============================================================================================
## set up environment
##==============================================================================================

# add /usr/bin to PATH in case we can't find mysql binary when executing from crontab
export PATH=/usr/bin:$PATH

##==============================================================================================
## cleanup old files
##==============================================================================================
sh "${BASE_PATH}"/cleanup.sh "${BASE_PATH}" || exit 1

##==============================================================================================
## executing all check scripts formats print
## this will also generate report in /report and error report in err_report/
##==============================================================================================
if python "${BASE_PATH}"/geenrate_report.py "$BASE_PATH" "$MYSQL_LOGIN_PATH" "$PROMETHEUS_ADDRESSES"
then
  echo "success!"
else
  echo -e "error!\t fail to execute check scripts."
  exit 1
fi

##==============================================================================================
## update baselines with the latest report
##==============================================================================================
sh "${BASE_PATH}"/update_baselines.sh "${BASE_PATH}" || exit 1

##==============================================================================================
## upload to ftp server
##==============================================================================================
sh "${BASE_PATH}"/upload_to_ftp.sh "${BASE_PATH}" "${FTP_SERVER}"|| exit 1

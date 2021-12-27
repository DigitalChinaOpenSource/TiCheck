#!/bin/bash

BASE_PATH=$1
REPORT_CLEANING_INTERVAL=$2

##==============================================================================================
## cleanup err_report/
##==============================================================================================
if [ -d "${BASE_PATH}"/err_report ]; then
  rm -f "${BASE_PATH}"/err_report/*
else
    echo -e "error!\t can't find err_report directory"
fi

##==============================================================================================
## cleanup report/ according to cleaning interval
##==============================================================================================
if [ -d "${BASE_PATH}"/report ];then
    find "${BASE_PATH}"/report -name "report_*.csv" -mtime +"$REPORT_CLEANING_INTERVAL" -exec rm {} \;
else
    echo -e "error!\t can't find report directory"
fi

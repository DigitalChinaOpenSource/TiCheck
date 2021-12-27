#!/bin/bash

BASE_PATH=$1

LATEST_REPORT=$(ls -lrt "$BASE_PATH"| tail -1 | awk '{ print $9}')

cp "$BASE_PATH"/"$LATEST_REPORT" "$BASE_PATH"/config/baseline.csv

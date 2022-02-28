#!/bin/bash

# load current path as base path
BASE_PATH=$(cd "$(dirname "$0")" || exit 1; pwd)

code "$BASE_PATH"/../config/load_variables.sh
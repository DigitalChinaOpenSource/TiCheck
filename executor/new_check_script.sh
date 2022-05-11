#!/bin/bash

# load current path as base path
BASE_PATH=$(cd "$(dirname "$0")" || exit 1; pwd)

cat << EOF > ../script/new_script
# parameters to this script:
#   [path to TiCheck] [MySQL login-path] [Prometheus address] [check arguments from config]
#
# DONT FORGET TO CHANGE FILENAME AND SAVE IT
# DONT FORGET TO ENABLE IT IN EXECUTION CONFIG
EOF

code "$BASE_PATH"/../script/new_script
#!/bin/bash
RED='\033[0;31m'
NC='\033[0m'
GREEN='\033[32m'
YELLOW='\033[33m'

check_dependency(){
    echo "+--------------------------------+-------------+"
    if ! type $1 >/dev/null 2>&1; then
        printf "+ %-30s |${RED} %-10s ${NC} +\n" $1  "fail";
    else
        printf "+ %-30s |${GREEN} %-10s ${NC} +\n" $1  "pass";
    fi
}


main(){
    # dependency install check
    echo "+--------------------------------+-------------+"
    printf "+ %-30s | %-10s  +\n" "Dependency detection"  "State";
    check_dependency npm
    check_dependency python3
    check_dependency pip
    check_dependency mysql
    check_dependency mysql_config_editor
    check_dependency expect
    echo "+--------------------------------+-------------+"
    echo -e "${YELLOW}Warning: ${NC}If the above items fail to be detected, it may affect the normal operation of the inspection function. \n"

    # build project
    make build

    # run server
    result=$(ps -aux | grep ticheck-server | grep -v "grep" | wc -l)
    if [[ $result -gt 0 ]];then
        kill -9 $(ps aux | grep ticheck-server | grep -v "grep" | awk '{print $2}')
    fi
    export GIN_MODE=release
    nohup ./bin/ticheck-server --work-dir=$(pwd)'/' --port=8066 &
    echo -e "${GREEN}TiCheck Successfuly installed. Now visit http://localhost:8066 to enjoy it ${NC}"
}

main;

exit 0
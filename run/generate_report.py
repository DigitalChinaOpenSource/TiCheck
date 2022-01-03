# coding=utf-8

# usage: generate_report.py base_path mysql_login_path prometheus_address

# config format:
# enable:check_class:check_name:check_script:operator:threshold:check_args

# writes result into two tables in sqlite3 database at report/report.db
# check_history and check_data
# for more details, read doc/table_schema_design.md

import os
import sys
import sqlite3
import timeit

base_path = sys.argv[1]
mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]
check_time = sys.argv[4]

config_file = base_path + "/config/execution_config.csv"
script_dir = base_path + "/script"
target_db = base_path + "report/report.db"
target_file = base_path + "/report/" + check_time + ".csv"


# write one line to log file
def write_check_data_csv(check_data_list):
    with open(target_file, 'a') as log:
        text = "{0[0]}||{0[1]}||{0[2]}||{0[3]}||{0[4]}||{0[5]}||{0[6]}||{0[7]}||{0[8]}".format(check_data_list)
        log.write(text)


def write_check_data_db(cursor, check_data_list):
    sql = "INSERT INTO check_data VALUES (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
    cursor.execute(sql, check_data_list)


def write_check_history_db(cursor, check_history_list):
    sql = "INSERT INTO check_history VALUES (?, ?, ?, ?, ?)"
    cursor.execute(sql, check_history_list)


def write_check_data(conn, data):
    write_check_data_csv(data)

    cursor = conn.cursor()
    write_check_data_db(cursor, data)
    cursor.close()
    conn.commit()


def write_check_history(conn, data):
    cursor = conn.cursor()
    write_check_history_db(cursor, data)
    cursor.close()
    conn.commit()


# initialize both check_data and check_history
# it should be idempotent
def initialize_db(conn):
    cursor = conn.cursor()
    create_check_data = "CREATE TABLE IF NOT EXISTS check_data (" \
                        "id INTEGER PRIMARY KEY AUTOINCREMENT," \
                        "check_time INTEGER NOT NULL," \
                        "check_class TEXT NOT NULL," \
                        "check_name TEXT NOT NULL," \
                        "operator TEXT NOT NULL," \
                        "threshold REAL NOT NULL," \
                        "duration INTEGER NOT NULL," \
                        "check_item TEXT NOT NULL," \
                        "check_value REAL," \
                        "check_status TEXT NOT NULL) "
    create_check_history = "CREATE TABLE IF NOT EXISTS check_history (" \
                           "check_time INTEGER NOT NULL," \
                           "normal_items INTEGER NOT NULL," \
                           "warning_items INTEGER NOT NULL," \
                           "total_items INTEGER NOT NULL," \
                           "duration INTEGER NOT NULL) "

    cursor.execute(create_check_data)
    cursor.execute(create_check_history)
    cursor.close()
    conn.commit()


# read config file into a 2d array
def read_config():
    return_list = []
    with open(config_file, "r") as fd:
        for line in fd.readlines():
            # skip empty line
            if not line.strip('\n'):
                continue
            (enable, check_class, check_name, script_name, operator, threshold, check_args) = \
                line.strip('\n').split(':')
            # skip disabled options
            if enable != "ENABLE":
                continue
            return_list.append([check_class, check_name, script_name, operator, threshold, check_args])
    return return_list


# read baseline file, if baseline file contains abnormal result, add last appearance time to dictionary
# return dictionary
def read_baseline(cursor):
    return_dict = {}
    result = cursor.execute("SELECT check_time FROM check_history DESC LIMIT 1")
    if len(result) == 0:
        return return_dict
    last_check_time = result[0][0]
    result = cursor.execute("SELECT check_name, check_item "
                            "FROM check_data " +
                            "WHERE check_time =" + last_check_time +
                            "AND check_status != '正常'")
    for row in result:
        return_dict.update({row[0] + row[1]: None})
    return return_dict


def run_script(script_name, check_args):
    script_result = None
    script_file = script_dir + '/' + script_name

    if not os.path.exists(script_file):
        print("error!\tscript: " + script_file + " not exits")
        return 0, None

    start_time = timeit.default_timer()
    if script_name.endswith(".sh"):
        script_result = os.popen(
            "sh %s %s %s %s %s" % (script_file, base_path, mysql_login_path, prometheus_address, check_args)
        ).read().splitlines()
    elif script_name.endswith(".py"):
        script_result = os.popen(
            "python %s %s %s %s %s" % (script_file, base_path, mysql_login_path, prometheus_address, check_args)
        ).read().splitlines()
    else:
        print("error!\tscript: " + script_file + "not supported")
    elapsed_time = timeit.default_timer() - start_time
    elapsed_int = int(elapsed_time * 1000)
    return elapsed_int, script_result


# check if script result within threshold
# only support datatype that can be turned into float
def compare_threshold(operator, threshold, script_result):
    if threshold.endswith('%'):
        threshold = float(threshold.strip('%')) / 100
    else:
        threshold = float(threshold)

    if script_result.endswith('%'):
        script_result = float(script_result.strip('%')) / 100
    else:
        script_result = float(script_result)

    if operator == "小于等于":
        return script_result <= threshold
    elif operator == "大于等于":
        return script_result >= threshold
    elif operator == "等于":
        return script_result == threshold
    else:
        print("error!\tcomparator: " + operator + " not supported")
        sys.exit(1)


def run_all():
    # if no config file, exit
    if not os.path.exists(config_file):
        print("error\tconfig file: " + config_file + " not exists")
        sys.exit(1)

    normal_items = 0
    warning_items = 0
    total_items = 0
    total_duration = 0

    conn = sqlite3.connect(target_db)
    past_abnormality = read_baseline(conn.cursor())
    # read config file
    check_list = read_config()

    # csv header
    write_check_data_csv(["开始时间", "检查类别", "检查项", "比较方式", "阈值", "巡检用时", "检查指标", "指标数值", "检查结果"])

    for check_class, check_name, script_name, operator, threshold, check_args in check_list:
        # execute script and get result
        # this result should be a list with even index as check_item and odd index as check_value
        duration, script_result = run_script(script_name, check_args)

        # skip to next item if some error occurs
        if script_result is None:
            continue

        result_len = len(script_result)

        # check if error is new by checking if error exists in past_abnormality
        if result_len == 0:
            check_item = script_name  # when no return, use script_name as check_item
            check_value = "无数据"  # when no return, use "无数据" as check_value
            abnormality_key = check_name + check_item  # remember key in past_abnormality is check_name + check_item
            if operator == "无数据":
                check_status = "正常"
                normal_items += 1
            # if key exist in past_abnormality, then its existing error, else new error
            elif abnormality_key in past_abnormality:
                check_status = "异常_已有"
                warning_items += 1
            else:
                check_status = "异常_新增"
                warning_items += 1

            total_items += 1
            total_duration += duration

            data_to_write = [check_time,
                             check_class,
                             check_name,
                             operator,
                             threshold,
                             duration,
                             check_item,
                             check_value,
                             check_status]

            write_check_data(conn, data_to_write)
            continue

        # when there are result, check result every two lines
        # first line is check_item, second line is check_value
        # when operator is "无数据", then check_status is "异常"
        index = 0
        while index < len(script_result):
            check_item = script_result[index]
            check_value = script_result[index + 1]
            abnormality_key = check_name + check_item
            if operator == "无数据":
                check_status = "异常_已有" if abnormality_key in past_abnormality else "异常_新增"
                warning_items += 1
            elif operator == "大于等于" or operator == "小于等于" or operator == "等于":
                if compare_threshold(operator, threshold, check_value):
                    check_status = "正常"
                    normal_items += 1
                else:
                    check_status = "异常_已有" if abnormality_key in past_abnormality else "异常_新增"
                    warning_items += 1
            else:
                print("error!\toperator: " + operator + " not supported")
                check_status = "operator error"

            total_items += 1
            total_duration += duration
            data_to_write = [check_time,
                             check_class,
                             check_name,
                             operator,
                             threshold,
                             duration,
                             check_item,
                             check_value,
                             check_status]
            write_check_data(conn, data_to_write)
            index += 2

    write_check_history_db(conn, [check_time,
                                  normal_items,
                                  warning_items,
                                  total_items,
                                  total_duration])
    conn.close()


if __name__ == '__main__':
    run_all()

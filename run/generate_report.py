# coding=utf-8

# Usage: generate_report.py base_path mysql_login_path prometheus_addresses
# will read baselines.csv under base_path/config/
# will produce report_[datetime].csv under base_path/report
# will read execution_config.csv under base_path/config

# config format:
#   ENABLE:check_class:check_name:script_name:operator:threshold:check_args
# printout format:
#   check_class check_name check_item colorized_check_status check_value operator+threshold
# log/baseline:
#   检查类别||检查项||检查对象||检查结果||检查值||检查阈值||出现时间
#   check_class||check_name||check_item||check_status||check_value||operator+threshold||last_appearance

import os
import sys
import requests

base_path = sys.argv[1]
mysql_login_path = sys.argv[2]
prometheus_addresses = sys.argv[3]
current_time = os.popen('date "+%Y%m%d%H%M%S"').read().strip('\n')
config_file = base_path + "/config/execution_config.csv"
script_dir = base_path + "/script"
baseline_file = base_path + "/config/baseline.csv"
# output target
output_log = base_path + "/report/report_" + current_time + ".csv"

# printout column width:
a0 = 8
a1 = 25
a2 = 30
a3 = 10
a4 = 10
a5 = 13


# check if a unicode character is chinese
def is_chinese(uchar):
    if u'\u4e00' <= uchar <= u'\u9fa5':
        return True
    else:
        return False


# fill text to width, will not cut
# left-align
def format_text(text, width, colored=False):
    str_text = str(text)
    unicode_text = str_text.decode("utf-8")
    cn_count = 0
    color_char = 0
    if colored:
        color_char = 9
    for u in unicode_text:
        if is_chinese(u):
            cn_count += 1

    return str_text + (width + color_char - cn_count - len(unicode_text)) * " "


# format and print one line
def format_line(str0, str1, str2, str3, str4, str5, colored_result=False):
    fmt = "| {0} | {1} | {2} | {3} | {4} | {5} |"
    return fmt.format(format_text(str0, a0),
                      format_text(str1, a1),
                      format_text(str2, a2),
                      format_text(str3, a3, colored_result),
                      format_text(str4, a4),
                      format_text(str5, a5)
                      )


# print line break
def print_linebreak():
    print(("+ {0:{0}<%d} + {0:{0}<%d} + {0:{0}<%d} + {0:{0}<%d} + {0:{0}<%d} + {0:{0}<%d} + " % (
        a0, a1, a2, a3, a4, a5)).format('-'))


# print titile
def print_title():
    print_linebreak()
    print(format_line("检查类别", "检查项", "检查对象", "检查结果", "检查值", "检查阈值"))
    print_linebreak()


# colorize result and printout onto screen
def print_text(check_class, check_name, check_item, check_status, check_value, operator, threshold):
    # define check_result color
    if check_status == "正常":
        check_status = "\033[32m%s\033[0m" % check_status
    else:
        check_status = "\033[31m%s\033[0m" % check_status
    print(format_line(check_class,
                      check_name,
                      check_item,
                      check_status,
                      check_value,
                      operator + threshold,
                      colored_result=True))
    print_linebreak()


# write one line to log file
def write_log(check_class, check_name, check_item, check_status, check_value, operator, threshold, last_appearance):
    with open(output_log, 'a') as log:
        text = "%s||%s||%s||%s||%s||%s||%s\n" % (check_class,
                                                 check_name,
                                                 check_item,
                                                 check_status,
                                                 check_value,
                                                 operator + threshold,
                                                 last_appearance)
        log.write(text)


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
def read_baseline():
    return_dict = {}
    with open(baseline_file, "r") as fd:
        for line in fd.readlines():
            _, check_name, check_item, status, _, _, last_appearance = line.strip('\n').split('||')
            if status != "正常":
                # use check_name + check_item as key to avoid collision of check_item
                return_dict.update({check_name + check_item: last_appearance})
    return return_dict


def run_script(script_name, prometheus_address, check_args):
    script_result = None
    script_file = script_dir + '/' + script_name

    if not os.path.exists(script_file):
        print("error!\tscript: " + script_file + " not exits")
        return None

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

    return script_result


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


# check if specific prometheus is alive
def prom_alive(address):
    try:
        response = requests.get('http://%s/api/v1/query' % address, params={'query': 'probe_success'})
        if response.json()["data"]['result']:
            return True
        else:
            return False
    except:
        return False


# try to find a prometheus that's alive, error out if no one is alive
def find_alive_prometheus(prom_addresses):
    for address in prom_addresses.split(','):
        if prom_alive(address):
            return address
    print("error!\tfailed to find available prometheus")
    sys.exit(1)


def run_all():
    # if no config file, exit
    if not os.path.exists(config_file):
        print("error\tconfig file: " + config_file + " not exists")
        sys.exit(1)

    baseline_exist = os.path.exists(baseline_file)
    if baseline_exist:
        past_abnormality = read_baseline()
    else:
        print("\tcan't find baseline file, please make sure this is a fresh run of the script")
        past_abnormality = {}

    alive_prometheus = find_alive_prometheus(prometheus_addresses)
    # read config file
    check_list = read_config()

    # print title and write title to log file
    print_title()
    write_log("检查类别", "检查项", "检查对象", "检查结果", "检查值", "检查阈值", "", "出现时间")

    for check_class, check_name, script_name, operator, threshold, check_args in check_list:
        # execute script and get result
        # this result should be a list with even index as check_item and odd index as check_value
        script_result = run_script(script_name, alive_prometheus, check_args)

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
            # if key exist in past_abnormality, then its existing error, else new error
            elif abnormality_key in past_abnormality:
                check_status = "异常_已有"
            else:
                check_status = "异常_新增"
            write_log(check_class, check_name, check_item, check_status, check_value, operator, threshold,
                      current_time)
            print_text(check_class, check_name, check_item, check_status, check_value, operator, threshold)
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
            elif operator == "大于等于" or operator == "小于等于" or operator == "等于":
                if compare_threshold(operator, threshold, check_value):
                    check_status = "正常"
                else:
                    check_status = "异常_已有" if abnormality_key in past_abnormality else "异常_新增"
            else:
                print("error!\toperator: " + operator + " not supported")
                check_status = "operator error"
            write_log(check_class, check_name, check_item, check_status, check_value, operator, threshold,
                      current_time)
            print_text(check_class, check_name, check_item, check_status, check_value, operator, threshold)
            index += 2


if __name__ == '__main__':
    run_all()

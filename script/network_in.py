# coding=utf-8

import sys
import requests

base_path = sys.argv[1]
mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]

check_args = sys.argv[3:] if len(sys.argv) > 3 else ""

pql = 'max_over_time( ( irate( node_network_receive_bytes_total{device!="lo"}[5m] ) )[1d:5m] ) / 1024 / 1024'

if check_args:
    device_filter = ""
    for device in check_args:
        device_filter = device_filter + device + "|"
    device_filter = device_filter.strip("|")
    pql = 'max_over_time( ( irate( node_network_receive_bytes_total{device=~"%s"}[5m] ) )[1d:5m] ) / 1024 / 1024' \
          % device_filter

try:
    response = requests.get('http://%s/api/v1/query' % prometheus_address, params={'query': pql})

    for result in response.json()['data']['result']:
        result_instance = result['metric']['instance']
        result_device = result['metric']['device']
        result_value = float(result['value'][1])
        print(result_instance + '/' + result_device)
        print("%.3f" % result_value)
except:
    sys.exit(1)

# coding=utf-8

import sys
import requests

base_path = sys.argv[1]
mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]

pql = 'sum(increase(tidb_server_execute_error_total[1d])) by (type, instance) > 0'

try:
    response = requests.get('http://%s/api/v1/query' % prometheus_address, params={'query': pql})

    for result in response.json()['data']['result']:
        result_instance = result['metric']['instance']
        result_type = result['metric']['type']
        result_value = result['value'][1].split('.')[0]
        print(result_instance + result_type)
        print(result_value)
except:
    sys.exit(1)

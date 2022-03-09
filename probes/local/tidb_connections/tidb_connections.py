# coding=utf-8

import sys
import requests

base_path = sys.argv[1]
mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]

pql = 'tidb_server_connections'

try:
    response = requests.get('http://%s/api/v1/query' % prometheus_address, params={'query': pql})

    for result in response.json()['data']['result']:
        result_instance = result['metric']['instance']
        result_value = result['value'][1]
        print(result_instance)
        print(result_value)
except:
    sys.exit(1)

# coding=utf-8

import sys
import requests

base_path = sys.argv[1]
mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]

pql = 'sum(increase(tidb_server_query_total{result="OK"}[1d])) by (result)'

try:
    response = requests.get('http://%s/api/v1/query' % prometheus_address, params={'query': pql})

    for result in response.json()['data']['result']:
        result_value = result['value'][1].split('.')[0]
        print('今天成功的Query数')
        print(result_value)
except:
    sys.exit(1)

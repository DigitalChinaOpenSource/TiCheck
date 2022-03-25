# coding=utf-8

import sys
import requests

# base_path = sys.argv[1]
# mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]
pql = 'sum(probe_success{group="pd"})'

try:
    response = requests.get('%s/api/v1/query' % prometheus_address, params={'query': pql})

    for result in response.json()['data']['result']:
        result_value = result['value'][1]
        print ("[tck_result:] PD节点数量="+result_value)
        
except :
    sys.exit(1)

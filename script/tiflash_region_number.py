# coding=utf-8

import sys
import requests

base_path = sys.argv[1]
mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]

pql = 'sum(tiflash_proxy_tikv_raftstore_region_count{type="region"}) by (instance)'

try:
    response = requests.get('http://%s/api/v1/query' % prometheus_address, params={'query': pql})

    for result in response.json()['data']['result']:
        result_address = result['metric']['instance']
        result_value = result['value'][1]
        print(result_address)
        print(result_value)
except:
    sys.exit(1)

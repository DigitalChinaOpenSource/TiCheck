# coding=utf-8

import sys
import requests

# base_path = sys.argv[1]
# mysql_login_path = sys.argv[2]
prometheus_address = sys.argv[3]

pql = 'node_memory_MemAvailable_bytes / 1024 / 1024'

try:
    response = requests.get('%s/api/v1/query' % prometheus_address, params={'query': pql})

    for result in response.json()['data']['result']:
        result_instance = result['metric']['instance']
        result_value = float(result['value'][1])
        print ("[tck_result:] "+result_instance+"="+"{:.3f}".format(result_value))
except:
    sys.exit(1)

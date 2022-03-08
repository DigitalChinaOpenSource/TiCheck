import request from '@/utils/request'

const clusterApi = {
  getClusterList: '/cluster/list',
  getClusterInfo: '/cluster/info/'
}

export function getClusterList () {
  return request({
    url: clusterApi.getClusterList,
    method: 'get'
  })
}

export function getClusterInfo (param) {
  return request({
    url: clusterApi.getClusterInfo + param,
    method: 'get'
  })
}

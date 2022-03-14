import request from '@/utils/request'

const clusterApi = {
  getClusterList: '/cluster/list',
  getClusterInfo: '/cluster/info/',
  addCluster: '/cluster/add',
  updateCluster: '/cluster/update/',
  getSchedulerList: '/cluster/scheduler/'
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

export function addCluster (parameter) {
  return request({
    url: clusterApi.addCluster,
    method: 'post',
    data: parameter
  })
}

export function updateCluster (id, parameter) {
  return request({
    method: 'post',
    url: clusterApi.updateCluster + id,
    data: parameter
  })
}

export function getSchedulerList (id) {
  return request({
    method: 'get',
    url: clusterApi.getSchedulerList + id
  })
}

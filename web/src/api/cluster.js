import request from '@/utils/request'

const clusterApi = {
  getClusterList: '/cluster/list',
  getClusterInfo: '/cluster/info/',
  getInitialCluster: '/cluster/initial/',
  addCluster: '/cluster/add',
  updateCluster: '/cluster/update/',
  getSchedulerList: '/cluster/scheduler/',
  addScheduler: '/cluster/scheduler/add',
  updateScheduler: '/cluster/scheduler/update',
  deleteScheduler: '/cluster/scheduler/delete/'
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

export function getInitialCluster (param) {
  return request({
    url: clusterApi.getInitialCluster + param,
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
    method: 'put',
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

export function addScheduler (parameter) {
  return request({
    url: clusterApi.addScheduler,
    method: 'post',
    data: parameter
  })
}

export function updateScheduler (parameter) {
  return request({
    method: 'put',
    url: clusterApi.updateScheduler,
    data: parameter
  })
}

export function deleteScheduler (id) {
  return request({
    method: 'delete',
    url: clusterApi.deleteScheduler + id
  })
}

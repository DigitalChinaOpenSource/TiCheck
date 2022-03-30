import request from '@/utils/request'
import i18n from '@/locales/index'

const checkApi = {
    checkHistory: '/cluster/report/all/',
    downloadReport: '/cluster/report/download/',
    getReportDetail: '/cluster/report/id/',
    getProbeList: '/cluster/probe/',
    getAddProbeList: '/cluster/probe/add/',
    addProbe: '/cluster/probe',
    deleteProbe: '/cluster/probe/',
    changeProbeStatus: '/cluster/probe/status',
    updateProbeConfig: '/cluster/probe/config',
    getExecuteInfo: '/cluster/check/info/',
    runExecute: '/cluster/check/run/',
}

export function getCheckHistoryByClusterID (clusterID, page, pageSize, startTime, endTime) {
    return request({
        url: checkApi.checkHistory + clusterID + '?page_num=' + page + '&page_size=' + pageSize + '&start_time=' + startTime + '&end_time=' + endTime,
        method: 'get'
    })
}

export function downloadReport (reportID) {
    return request({
        url: checkApi.downloadReport + reportID,
        method: 'get'
    })
}


export function getReportDetail(reportID) {
    return request({
        url: checkApi.getReportDetail + reportID,
        method: 'get'
    })
}

export function getProbeList(clusterID) {
    return request({
        url: checkApi.getProbeList + clusterID,
        method: 'get'
    })
}

export function getAddProbeList(clusterID) {
    return request({
        url: checkApi.getAddProbeList + clusterID,
        method: 'get'
    })
}

export function addProbe(params) {
    return request({
        url: checkApi.addProbe,
        method: 'post',
        data: params
    })
}

export function changeProbeStatus(params) {
    return request({
        url: checkApi.changeProbeStatus,
        method: 'put',
        data: params
    })
}

export function getExecuteInfo(clusterID) {
    return request({
        url: checkApi.getExecuteInfo + clusterID,
        method: 'get'
    })
}

export function runExecute(clusterID) {
    return new WebSocket("ws://localhost:8081" + checkApi.runExecute + clusterID)
}

export function updateProbeConfig(params) {
    return request({
        url: checkApi.updateProbeConfig,
        method: 'put',
        data: params
    })
}

export function deleteProbe(id) {
    return request({
        url: checkApi.deleteProbe + id,
        method: 'delete'
    })
}

export function mapOperatorValue (operator) {
    switch (operator) {
    case 0:
      return 'NA'
    case 1:
      return '='
    case 2:
      return '>'
    case 3:
      return '>='
    case 4:
      return '<'
    case 5:
      return '<='
    default:
      return 'NA'
  }
}

export function mapEnableValue (enable) {
    switch (enable) {
    case 0:
      return false
    case 1:
      return true
    default:
      return false
  }
}

export function mapTagText (tag) {
    switch (tag.toLowerCase()) {
    case 'cluster':
        return i18n.t('check.probe.tag.cluster')
    case 'network':
        return i18n.t('check.probe.tag.network')
    case 'running_state':
        return i18n.t('check.probe.tag.running_state')
    case 'others':
        return i18n.t('check.probe.tag.others')
    }
}
import request from '@/utils/request'

const checkApi = {
    checkHistory: '/cluster/report/all/',
    downloadReport: '/cluster/report/download/'
}

export function getCheckHistoryByClusterID (parameter) {
    return request({
        url: checkApi.checkHistory + parameter,
        method: 'get'
    })
}

export function downloadReport (parameter) {
    return request({
        url: checkApi.downloadReport + parameter,
        method: 'get'
    })
}

import request from '@/utils/request'

const checkApi = {
    checkHistory: '/cluster/report/all/',
    downloadReport: '/cluster/report/download/',
    getReportDetail: '/cluster/report/id/'
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
        url: getReportDetail + reportID,
        method: 'get'
    })
}
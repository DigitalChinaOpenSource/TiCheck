<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240)"
      title= "cluster.detail.check.history"
    />
    <a-table :columns="columns" :data-source="data" rowKey="id">
      <a slot="ID" slot-scope="text">{{ text }}</a>
      <span slot="customTitle"><a-icon type="smile-o" /> ID </span>
      <span slot="normal_items" slot-scope="normal_items">
        <a-tag
          :key="normal_items"
          :color="'green'"
        >
          {{ normal_items }}
        </a-tag>
      </span>
      <span slot="warning_items" slot-scope="warning_items">
        <a-tag
          :key="warning_items"
          :color="'volcano'"
        >
          {{ warning_items }}
        </a-tag>
      </span>
      <span slot="action" slot-scope="record">
        <a @click="downloadReport(record.id)">Download</a>
      </span>
    </a-table>
  </div>
</template>
<script>
import { getCheckHistoryByClusterID, downloadReport } from '@/api/check'

const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    key: 'id',
    slots: { title: 'customTitle' },
    scopedSlots: { customRender: 'ID' }
  },
  {
    title: 'Time',
    dataIndex: 'check_time',
    key: 'check_time'
  },
  {
    title: 'Duration',
    dataIndex: 'duration',
    key: 'duration'
  },
  {
    title: 'Normal',
    key: 'normal_items',
    dataIndex: 'normal_items',
    scopedSlots: { customRender: 'normal_items' }
  },
  {
    title: 'Warning',
    key: 'warning_items',
    dataIndex: 'warning_items',
    scopedSlots: { customRender: 'warning_items' }
  },
  {
    title: 'Total',
    key: 'total_items',
    dataIndex: 'total_items'
  },
  {
    title: 'Action',
    key: 'action',
    scopedSlots: { customRender: 'action' }
  }
]

const data = []

export default {
  data () {
    return {
      data,
      columns
    }
  },
  mounted () {
    this.getHistoryList()
  },
  methods: {
    getHistoryList () {
      getCheckHistoryByClusterID(1).then(res => {
        this.data = res
      })
    },
    downloadReport (params) {
      console.log(params)
      downloadReport(params)
    }
  }
}
</script>

<template>
  <a-table :columns="columns" :data-source="data">
    <a slot="ID" slot-scope="text">{{ text }}</a>
    <span slot="customTitle"><a-icon type="smile-o" /> ID </span>
    <span slot="normal" slot-scope="normal">
      <a-tag
        :key="normal"
        :color="'green'"
      >
        {{ normal }}
      </a-tag>
    </span>
    <span slot="warning" slot-scope="warning">
      <a-tag
        :key="warning"
        :color="'volcano'"
      >
        {{ warning }}
      </a-tag>
    </span>
    <span slot="action" slot-scope="record">
      <a>Download - {{record.id}}</a>
    </span>
  </a-table>
</template>
<script>
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
    dataIndex: 'start_time',
    key: 'start_time'
  },
  {
    title: 'Duration',
    dataIndex: 'duration',
    key: 'duration'
  },
  {
    title: 'Normal',
    key: 'normal',
    dataIndex: 'normal',
    scopedSlots: { customRender: 'normal' }
  },
  {
    title: 'Warning',
    key: 'warning',
    dataIndex: 'warning',
    scopedSlots: { customRender: 'warning' }
  },
  {
    title: 'Total',
    key: 'total',
    dataIndex: 'total',
    scopedSlots: { customRender: 'total' }
  },
  {
    title: 'Action',
    key: 'action',
    scopedSlots: { customRender: 'action' }
  }
]

const data = [
  {
    id: 1,
    start_time: '2022-10-10 23:00:00',
    duration: 32,
    normal: 12,
    warning: 21,
    total: 33
  }
]

export default {
  data () {
    return {
      data,
      columns
    }
  }
}

</script>

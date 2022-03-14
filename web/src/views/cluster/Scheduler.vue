<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px;"
      :title="$t('cluster.scheduler.title')"
    />
    <a-button type="primary" @click="addScheduler" style="float: right">
      {{ $t('cluster.scheduler.btn.add') }}
    </a-button>

    <a-table
      :columns="columns"
      :rowKey="(scheduler) => scheduler.id"
      :data-source="schedulerList"
      :pagination="paginationOpt"
      style="padding-top: 60px"
    >
      <span slot="status" slot-scope="status">
        {{ mapStatusValue(status) }}
      </span>
      <span slot="action" slot-scope="scheduler">
        <a @click="jump2History(scheduler.id)" target="_blank" style="color: #40a9ff">{{ $t('cluster.scheduler.table.action.history') }}</a>
      </span>
    </a-table>
  </div>
</template>

<script>
import { getSchedulerList } from '@/api/cluster'

const schedulerList = []

const paginationOpt = {
  showTotal: total => `Total ${total} items`,
  showSizeChanger: true,
  pageSizeOptions: ['10', '30', '50', '100'],
  defaultPageSize: 100
}

export default {
  name: 'ClusterScheduler',
  data () {
    return {
      schedulerList,
      paginationOpt,
      columns: [
        {
          title: this.$t('cluster.scheduler.table.id'),
          dataIndex: 'id',
          key: 'id',
          hide: true
        },
        {
          title: this.$t('cluster.scheduler.table.name'),
          dataIndex: 'name',
          key: 'name'
        },
        {
          title: this.$t('cluster.scheduler.table.create'),
          dataIndex: 'create_time',
          key: 'create_time'
        },
        {
          title: this.$t('cluster.scheduler.table.cron'),
          dataIndex: 'cron',
          key: 'cron'
        },
        {
          title: this.$t('cluster.scheduler.table.status'),
          key: 'status',
          dataIndex: 'status',
          scopedSlots: { customRender: 'status' }
        },
        {
          title: this.$t('cluster.scheduler.table.count'),
          key: 'count',
          dataIndex: 'count'
        },
        {
          title: this.$t('cluster.scheduler.table.action'),
          key: 'action',
          scopedSlots: { customRender: 'action' }
        }
      ]
    }
  },
  methods: {
    getSchedulers () {
      const id = this.$route.params.id?.toString()
      getSchedulerList(id).then(res => {
        this.schedulerList = res.data
        console.log('list=>', this.schedulerList)
      })
    },
    mapStatusValue (status) {
      switch (status) {
        case 0:
          return this.$t('cluster.scheduler.table.status.no')
        case 1:
          return this.$t('cluster.scheduler.table.status.yes')
        default:
          return this.$t('cluster.scheduler.table.status.null')
      }
    },
    jump2History (schedulerID) {
      console.log('jump to history detail view', schedulerID)
    },
    addScheduler () {
      console.log('add scheduler')
    }
  },
  mounted () {
    this.getSchedulers()
  }
}

</script>

<template>
  <div>
    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Cluster Info"
    />
    <a-card :bordered="false" :style="{ marginTop: '5px' }" :data-sourece="clusterInfo">
      <a-descriptions title="" size="middle">
        <a-descriptions-item label="Cluster Name">{{ clusterInfo[0].cluster_name }}</a-descriptions-item>
        <a-descriptions-item label="Cluster Version">{{ clusterInfo[0].cluster_version }}</a-descriptions-item>
        <a-descriptions-item label="Create Time">{{ clusterInfo[0].create_time }}</a-descriptions-item>
        <a-descriptions-item label="Cluster User" span="3">{{ clusterInfo[0].cluster_owner }}</a-descriptions-item>
        <a-descriptions-item label="Cluster Description">{{ clusterInfo[0].cluster_description }}</a-descriptions-item>
      </a-descriptions>
    </a-card>

    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Cluster Status"
    />
    <a-row :gutter="24" :style="{marginTop:'5px'}">
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <chart-card :loading="loading" title="累积巡检次数" total="123">
        </chart-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <chart-card :loading="loading" title="累积检查项目" :total="8846 | NumberFormat">
        </chart-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <chart-card :loading="loading" title="最后巡检时间" :total="6560 | NumberFormat">
        </chart-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <chart-card :loading="loading" title="集群健康度" total="78%">
        </chart-card>
      </a-col>
    </a-row>

    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Recent Alerts"
    />
    <a-card :loading="loading" :bordered="false" :body-style="{padding: '0'}">
      <div class="salesCard">
        <a-row>
          <a-col :xl="16" :lg="12" :md="12" :sm="24" :xs="24">
            <bar :data="barData" :title="$t('dashboard.analysis.sales-trend')" :style="{ marginTop: '24px'}"/>
          </a-col>
          <a-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24">
            <rank-list :title="$t('dashboard.analysis.sales-ranking')" :list="rankList"/>
          </a-col>
        </a-row>
      </div>
    </a-card>
  </div>
</template>

<script>
import { getClusterInfo } from '@/api/cluster'
import {
  ChartCard,
  RankList,
  Bar
} from '@/components'
// import { getClusterInfo } from '@/api/manage'
const clusterInfo = []
export default {
  name: 'ClusterInfo',
  clusterID: '',
  components: {
    ChartCard,
    RankList,
    Bar
  },
  data () {
    return {
      clusterInfo
    }
  },
  created () {
    this.clusterID = this.$route.params.id?.toString()
    this.localClusterInfo()
    setTimeout(() => {
      this.loading = !this.loading
    }, 1000)
  },
  methods: {
    localClusterInfo () {
      getClusterInfo(this.clusterID).then(res => {
        this.clusterInfo[0] = res.data
      })
    }
  }
}
</script>

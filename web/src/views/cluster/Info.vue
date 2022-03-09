<template>
  <div>
    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Cluster Info"
    />
    <a-card :bordered="true" :style="{ marginTop: '5px' }">
      <a-descriptions title="" size="middle">
        <a-descriptions-item label="Cluster Name">{{ clusterInfo.name }}</a-descriptions-item>
        <a-descriptions-item label="Cluster Version">{{ clusterInfo.version }}</a-descriptions-item>
        <a-descriptions-item label="Create Time">{{ clusterInfo.create_time }}</a-descriptions-item>
        <a-descriptions-item label="Cluster User" span="3">{{ clusterInfo.owner }}</a-descriptions-item>
        <a-descriptions-item label="Cluster Description">{{ clusterInfo.description }}</a-descriptions-item>
      </a-descriptions>
    </a-card>

    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Cluster Status"
    />
    <a-row :gutter="24" :style="{marginTop:'5px'}">
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <a-card title="累积巡检次数" style="text-align: center">
          <span>
            {{ clusterInfo.check_count }}
          </span>
        </a-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <a-card title="累积检查项目" style="text-align: center">
          <span>
            {{ clusterInfo.check_total }}
          </span>
        </a-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <a-card title="最后巡检时间" style="text-align: center">
          <span>
            {{ clusterInfo.last_check_time }}
          </span>
        </a-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <a-card title="集群健康度" style="text-align: center">
          <span>
            {{ clusterInfo.cluster_health }}
          </span>
        </a-card>
      </a-col>
    </a-row>
    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Recent Alerts"
    />
    <a-card :bordered="false" :body-style="{padding: '0'}">
      <div class="salesCard">
        <a-row>
          <a-col :xl="16" :lg="12" :md="12" :sm="24" :xs="24">
            <bar :title="$t('dashboard.analysis.sales-trend')" :style="{ marginTop: '24px'}" :data="clusterInfo.recent_warning_items" />
          </a-col>
          <a-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24">
            <a-tabs default-active-key="1" size="large" :tab-bar-style="{marginBottom: '24px', paddingLeft: '16px'}" style="color: #40a9ff">
              <a-tab-pane loading="true" tab="weekly" key="1">
                <rank-list :list="clusterInfo.weekly_history_warnings" :style="{ marginTop: '24px'}"/>
              </a-tab-pane>
              <a-tab-pane tab="monthly" key="2">
                <rank-list :list="clusterInfo.monthly_history_warnings" :style="{ marginTop: '24px'}"/>
              </a-tab-pane>
              <a-tab-pane tab="yearly" key="3">
                <rank-list :list="clusterInfo.yearly_history_warnings" :style="{ marginTop: '24px'}"/>
              </a-tab-pane>
            </a-tabs>
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
const clusterInfo = {}
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
        this.clusterInfo = res.data
      })
    }
  }
}
</script>

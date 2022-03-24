<template>
  <div>
    <a-page-header
      :ghost="false"
      :style="{ marginTop: '24px'}"
      :title="$t('cluster.info.cluster-info')"
    >
      <template slot="extra">
        <a-button type="primary">{{ $t('cluster.info.check') }}</a-button>
      </template>
      <a-card :bordered="true" :style="{ marginTop: '5px' }">
        <a-descriptions title="" size="middle">
          <a-descriptions-item :label="$t('cluster.info.name')">{{ clusterInfo.name }}</a-descriptions-item>
          <a-descriptions-item :label="$t('cluster.info.version')">{{ clusterInfo.version }}</a-descriptions-item>
          <a-descriptions-item :label="$t('cluster.info.create-time')">{{ clusterInfo.create_time | moment }}</a-descriptions-item>
          <a-descriptions-item :label="$t('cluster.info.owner')" span="3">{{ clusterInfo.owner }}</a-descriptions-item>
          <a-descriptions-item :label="$t('cluster.info.description')">{{ clusterInfo.description }}</a-descriptions-item>
        </a-descriptions>
      </a-card>
    </a-page-header>
    <a-page-header
      :style="{ marginTop: '24px'}"
      :title="$t('cluster.info.status')"
    >
      <a-row :gutter="24" :style="{marginTop:'5px'}">
        <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
          <a-card :title="$t('cluster.info.status.count')" style="text-align: center">
            <span>
              {{ clusterInfo.check_count }}
            </span>
          </a-card>
        </a-col>
        <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
          <a-card :title="$t('cluster.info.status.total')" style="text-align: center">
            <span>
              {{ clusterInfo.check_total }}
            </span>
          </a-card>
        </a-col>
        <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
          <a-card :title="$t('cluster.info.status.last')" style="text-align: center">
            <span>
              {{ clusterInfo.last_check_time | moment }}
            </span>
          </a-card>
        </a-col>
        <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
          <a-card :title="$t('cluster.info.status.healthy')" style="text-align: center">
            <span>
              {{ clusterInfo.cluster_health }}
            </span>
          </a-card>
        </a-col>
      </a-row>
    </a-page-header>
    <a-page-header
      :style="{ marginTop: '24px'}"
      :title="$t('cluster.info.recent')"
    >
      <a-card :bordered="false" :body-style="{padding: '0'}">
        <div class="salesCard">
          <a-row>
            <a-col :xl="16" :lg="12" :md="12" :sm="24" :xs="24">
              <bar :style="{ marginTop: '24px'}" :data="clusterInfo.recent_warning_items" />
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
    </a-page-header>
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
      getClusterInfo(this.clusterID)
        .then(res => { this.clusterInfo = res.data })
        .catch(() => {
          this.$router.push({ path: '/' })
        })
    }
  }
}
</script>

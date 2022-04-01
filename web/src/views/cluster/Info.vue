<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      :title="$t('cluster.info.cluster-info')"
    >
      <template slot="extra">
        <a-button type="primary" style="float: right" @click="doCheck">{{
          $t("cluster.info.check")
        }}</a-button>
      </template>
    </a-page-header>
    <a-row :gutter="24" :style="{ margin: '5px 0px 0px 0px',background:'#ECECEC' }">
      <a-col :sm="24" :md="24" :xl="24" :style="{ marginBottom: '4px' }">
        <div>
          <a-card :style="{ marginTop: '12px' }">
            <a-descriptions title="" size="default">
              <a-descriptions-item :label="$t('cluster.info.name')">{{
                clusterInfo.name
              }}</a-descriptions-item>
              <a-descriptions-item :label="$t('cluster.info.version')">{{
                clusterInfo.version
              }}</a-descriptions-item>
              <a-descriptions-item :label="$t('cluster.info.create-time')">{{
                clusterInfo.create_time | moment
              }}</a-descriptions-item>
              <a-descriptions-item :label="$t('cluster.info.owner')" span="3">{{
                clusterInfo.owner
              }}</a-descriptions-item>
              <a-descriptions-item :label="$t('cluster.info.description')" span="3">{{
                clusterInfo.description
              }}</a-descriptions-item>
            </a-descriptions>
          </a-card>
        </div>
      </a-col>
    </a-row>
    <!-- <a-page-header
      style="
        margin-top: 24px;
        border: 1px solid rgb(235, 237, 240);
        margin-bottom: 20px;
      "
      :title="$t('cluster.info.status')"
    /> -->
    <a-row :gutter="24" :style="{ margin: '0px',background:'#ECECEC', padding:'15px 0px' }">
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <!-- <a-card
          :title="$t('cluster.info.status.count')"
          style="text-align: center"
        >
          <span>
            {{ clusterInfo.check_count }}
          </span>
        </a-card> -->
        <chart-card :loading="false" :title="$t('cluster.info.status.count')" :total="clusterInfo.check_count | NumberFormat">
          <a-tooltip :title="$t('dashboard.analysis.introduce')" slot="action">
            <a-icon type="info-circle-o" />
          </a-tooltip>
          <div>
            <mini-area />
          </div>
          <template slot="footer">{{ $t('cluster.info.status.count-today') }}：<span> {{ clusterInfo.today_check_count | NumberFormat }}</span></template>
        </chart-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <!-- <a-card
          :title="$t('cluster.info.status.total')"
          style="text-align: center"
        >
          <span>
            {{ clusterInfo.check_total }}
          </span>
        </a-card> -->
        <chart-card :loading="false" :title="$t('cluster.info.status.total')" :total="clusterInfo.check_total | NumberFormat">
          <a-tooltip :title="$t('dashboard.analysis.introduce')" slot="action">
            <a-icon type="info-circle-o" />
          </a-tooltip>
          <div>
            <mini-area />
          </div>
          <template slot="footer">{{ $t('cluster.info.status.total-today') }}：<span> {{ clusterInfo.today_check_total | NumberFormat }}</span></template>
        </chart-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <!-- <a-card
          :title="$t('cluster.info.status.last')"
          style="text-align: center"
        >
          <span>{{ clusterInfo.last_check_time | moment }}</span>
        </a-card> -->
        <chart-card :loading="false" :title="$t('cluster.info.status.last')" :total="clusterInfo.last_check_time | moment">
          <a-tooltip :title="$t('dashboard.analysis.introduce')" slot="action">
            <a-icon type="info-circle-o" />
          </a-tooltip>
          <template slot="footer">
            <span flag="down" style="margin-right: 16px;">
              <span slot="term">{{ $t('cluster.info.status.last-normal') }}：</span>
              {{ clusterInfo.last_check_normal }}
            </span>
            <span flag="up">
              <span slot="term">{{ $t('cluster.info.status.last-warning') }}：</span>
              {{ clusterInfo.last_check_warning }}
            </span>
          </template>
        </chart-card>
      </a-col>
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <!-- <a-card
          :title="$t('cluster.info.status.healthy')"
          style="text-align: center"
        >
          <span>
            {{ clusterInfo.cluster_health }}
          </span>
        </a-card> -->
        <chart-card :loading="false" :title="$t('cluster.info.status.healthy')" :total="clusterInfo.cluster_health">
          <a-tooltip :title="$t('dashboard.analysis.introduce')" slot="action">
            <a-icon type="info-circle-o" />
          </a-tooltip>
          <div>
            <mini-progress color="rgb(19, 194, 194)" :target="80" :percentage="clusterInfo.cluster_health" height="8px" />
          </div>
          <template slot="footer">
            <span>
              <span slot="term">{{ $t('cluster.info.status.healthy-update') }}：</span>
              {{ clusterInfo.health_update_time | moment }}
            </span>
          </template>
        </chart-card>
      </a-col>
    </a-row>
    <!-- <a-page-header
      style="
        margin-top: 24px;
        border: 1px solid rgb(235, 237, 240);
        margin-bottom: 20px;
      "
      :title="$t('cluster.info.recent')"
    /> -->
    <a-card :bordered="false" style="marginTop:20px" :body-style="{ padding: '0' }" :title="$t('cluster.info.recent')">
      <div class="salesCard">
        <a-row>
          <a-col :xl="16" :lg="12" :md="12" :sm="24" :xs="24">
            <bar
              :style="{ padding: '0' }"
              :data="clusterInfo.recent_warning_items"
            />
          </a-col>
          <a-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24">
            <a-tabs
              default-active-key="1"
              size="large"
              :tab-bar-style="{ margin: '0 0 0 20px' }"
              style="color: #40a9ff"
            >
              <a-tab-pane loading="true" tab="weekly" key="1">
                <rank-list
                  :list="clusterInfo.weekly_history_warnings"
                  :style="{ padding: '0 22px' }"
                />
              </a-tab-pane>
              <a-tab-pane tab="monthly" key="2">
                <rank-list
                  :list="clusterInfo.monthly_history_warnings"
                  :style="{ padding: '0 22px' }"
                />
              </a-tab-pane>
              <a-tab-pane tab="yearly" key="3">
                <rank-list
                  :list="clusterInfo.yearly_history_warnings"
                  :style="{ padding: '0 22px' }"
                />
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
import { ChartCard, RankList, Bar, MiniProgress, MiniArea } from '@/components'
import moment from 'moment'
const clusterInfo = {}
export default {
  name: 'ClusterInfo',
  clusterID: '',
  components: {
    ChartCard,
    RankList,
    Bar,
    MiniProgress,
    MiniArea
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
        .then((res) => {
          this.clusterInfo = res.data
          for (
            let i = 0;
            i < this.clusterInfo.recent_warning_items.length;
            i++
          ) {
            this.clusterInfo.recent_warning_items[i].time = moment(
              this.clusterInfo.recent_warning_items[i].time
            ).format('MM-DD HH:mm:ss')
          }
        })
        .catch(() => {
          this.$router.push({ path: '/' })
        })
    },
    doCheck () {
      this.$router.push({
        name: 'ExecuteCheck',
        params: { id: this.clusterID },
        query: { id: this.clusterID }
      })
    }
  },
  mounted () {
    this.localClusterInfo()
    setTimeout(() => {
      this.loading = !this.loading
    }, 1000)
  }
}
</script>

<template>
  <div>
    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Cluster Info"
    />
    <a-card :bordered="false" :style="{ marginTop: '5px' }">
      <a-descriptions title="" size="middle">
        <a-descriptions-item label="Cluster Name">cluster name</a-descriptions-item>
        <a-descriptions-item label="Cluster Version">v5.4.4</a-descriptions-item>
        <a-descriptions-item label="Create Time">2021-03-04 14:00:00</a-descriptions-item>
        <a-descriptions-item label="Cluster User" span="3">qiruia@qiruia</a-descriptions-item>
        <a-descriptions-item label="Cluster Description">this is a simple cluster!</a-descriptions-item>
      </a-descriptions>
    </a-card>

    <a-page-header
      :style="{ marginTop: '24px'}"
      title="Cluster Status"
    />
    <a-row :gutter="24" :style="{marginTop:'5px'}">
      <a-col :sm="24" :md="12" :xl="6" :style="{ marginBottom: '4px' }">
        <chart-card :loading="loading" title="累积巡检次数" total="78">
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
import {
  ChartCard,
  RankList,
  Bar
} from '@/compoments'
// import { getClusterInfo } from '@/api/manage'
export default {
  name: 'ClusterInfo',
  components: {
    ChartCard,
    RankList,
    Bar
  },
  created () {
    setTimeout(() => {
      this.loading = !this.loading
    }, 1000)
  }
}
</script>

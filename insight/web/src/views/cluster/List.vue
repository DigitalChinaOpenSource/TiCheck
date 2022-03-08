<template>
  <page-header-wrapper
    title="CLuster List"
  >
    <!-- actions -->
    <template v-slot:extra>
      <a-button type="primary" >Add Cluster</a-button>
    </template>
    <a-list
      rowKey="id"
      :grid="{ gutter: 16, xs: 1, sm: 1, md: 1, lg: 2, xl: 2, xxl: 2 }"
      :dataSource="dataSource"
    >
      <a-list-item
        slot="renderItem"
        slot-scope="item"
      >
        <template>
          <a-card
            :hoverable="true"
            :title="false"
          >
            <div>
              <a-row type="flex" justify="space-between">
                <a-col>
                  CLuster Name
                </a-col>
                <a-col>
                  Add Time :  {{ item.create_time }}
                </a-col>
              </a-row>
            </div>
            <div>
              <a-row type="flex">
                <div>
                  <a
                    class="text-black-bold"
                    style="font-size:30px"
                    @click="jump2Info(item)"
                  >
                    {{item.cluster_name}}
                  </a>
                </div>
              </a-row>
              <div style="margin-top: 15px;margin-bottom: 5px">
                Cluster Description
              </div>
              <div>
                <div v-if="item.description">
                  <div>
                    {{ item.description }}
                  </div>
                </div>
                <div v-else>
                  暂无项目简介
                </div>
              </div>
              <div style="margin-top: 15px">
                <a-row type="flex" justify="space-between">
                  <a-col>
                    <span>Node Info</span>
                  </a-col>
                  <a-col>
                    <span>
                      <a style="margin-right: 15px">Grafana</a>
                      <a>Dashboard</a>
                    </span>
                  </a-col>
                </a-row>
              </div>
            </div>
            <div>
              <a-list :grid="{ gutter: 16, column: 4 }" :data-source="nodeNum" style="margin-top: 25px">
                <a-list-item slot="renderItem" slot-scope="items">
                  <a-card :title="items.title" style="text-align: center">
                    {{ items.num }}
                  </a-card>
                </a-list-item>
              </a-list>
            </div>
            <div>
              <a-row type="flex" justify="end">
                <a-col>
                  Last Check Time: {{ item.last_check_time }}
                </a-col>
              </a-row>
            </div>
          </a-card>
        </template>
      </a-list-item>
    </a-list>
  </page-header-wrapper>
</template>

<script>
import { getClusterList } from '@/api/cluster'
import {
  ChartCard,
  RankList,
  Bar
} from '@/components'
// import { getClusterInfo } from '@/api/manage'

const dataSource = []
const nodeNum = [
  {
    id: 1,
    title: 'PD',
    num: 3
  },
  {
    id: 2,
    title: 'TiDB',
    num: 3
  },
  {
    id: 3,
    title: 'TiKV',
    num: 3
  },
  {
    id: 4,
    title: 'TiFlash',
    num: 0
  }
]

export default {
  name: 'ClusterList',
  components: {
    ChartCard,
    RankList,
    Bar
  },
  data () {
    return {
      extraImage: 'https://gw.alipayobjects.com/zos/rmsportal/RzwpdLnhmvDJToTdfDPe.png',
      dataSource,
      nodeNum
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    getList () {
      getClusterList().then(res => {
        this.dataSource = res.data
      })
    },
    jump2Info (item) {
      this.$router.push({ name: 'ClusterInfo', params: { id: item.id } })
    }
  },
  created () {
    setTimeout(() => {
      this.loading = !this.loading
    }, 1000)
  }
}
</script>

<style lang="less" scoped>

  .text-black-bold {
    color: black;
  }

  a.text-black-bold:hover {
    color: #40a9ff;
  }

</style>

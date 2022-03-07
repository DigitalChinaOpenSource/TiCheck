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
              <a-row>
                <a-col :span="8">
                  CLuster Name
                </a-col>
                <a-col :span="12" :offset="4">
                  Add Time :  {{ item.createtTime }}
                </a-col>
              </a-row>
            </div>
            <div>
              <a-row type="flex">
                <div>
                  <a
                    class="text-black-bold"
                    style="font-size:30px"
                  >
                    {{item.title}}
                  </a>
                </div>
              </a-row>
              <div style="margin-top: 15px;margin-bottom: 5px">
                Cluster Description
              </div>
              <div>
                <div v-if="item.content">
                  <div>
                    {{ item.content }}
                  </div>
                </div>
                <div v-else>
                  暂无项目简介
                </div>
              </div>
              <div style="margin-top: 15px">
                <a-row>
                  <a-col :span="8">
                    <span>Node Info</span>
                  </a-col>
                  <a-col :span="8" :offset="8">
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
              <a-row>
                <a-col :span="14" :offset="10">
                  Last Check Time: {{ item.createtTime }}
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
import {
  ChartCard,
  RankList,
  Bar
} from '@/components'
// import { getClusterInfo } from '@/api/manage'

const dataSource = [
  {
    id: 1,
    title: 'Alipay',
    content: 'this is a single tidb for test!',
    createtTime: '2022-01-01 22:22:22'
  },
  {
    id: 2,
    title: 'Alipay',
    content: 'this is a single tidb for test!',
    createtTime: '2022-01-01 22:22:22'
  }
]

const nodeNum = [
  {
    id: 1,
    title: 'PD Num',
    num: 3
  },
  {
    id: 2,
    title: 'TiDB Num',
    num: 3
  },
  {
    id: 3,
    title: 'TiKV Num',
    num: 3
  },
  {
    id: 4,
    title: 'TiFlash Num',
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
  methods: {
    testFun () {
      this.$message.info('快速开始被点击！')
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

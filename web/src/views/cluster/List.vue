<template>
  <page-header-wrapper
    :title="$t('cluster.list.cluster-list')"
  >
    <!-- actions -->
    <template v-slot:extra>
      <div>
        <a-button type="primary" @click="showModal" >{{ $t('cluster.list.add-cluster') }}</a-button>
        <a-modal v-model="modalVisible" :title="$t('cluster.list.add-cluster')" @ok="handleOk" @cancel="modalCancel" width="70%">
          <a-form :form="clusterForm">
            <a-form-item
              :label="$t('cluster.list.name')"
              :labelCol="{lg: {span: 7}, sm: {span: 7}}"
              :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
              <a-input
                v-decorator="['name',{rules: [{ required: true }]}]"
                :placeholder="$t('cluster.list.input.name')"
                name="name" />
            </a-form-item>
            <a-form-item
              :label="$t('cluster.list.prometheus')"
              :labelCol="{lg: {span: 7}, sm: {span: 7}}"
              :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
              <a-input
                name="url"
                addon-before="http://"
                :placeholder="$t('cluster.list.input.prometheus')"
                v-decorator="[
                  'url',
                  {
                    rules: [
                      {
                        required: true,
                        message: ''
                      }]
                  }]" />
            </a-form-item>
            <a-form-item
              :label="$t('cluster.list.user')"
              :labelCol="{lg: {span: 7}, sm: {span: 7}}"
              :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
              <a-input
                name="user"
                :placeholder="$t('cluster.list.input.user')"
                v-decorator="['user',{rules: [{ required: true }]}]" />
            </a-form-item>
            <a-form-item
              :label="$t('cluster.list.passwd')"
              :labelCol="{lg: {span: 7}, sm: {span: 7}}"
              :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
              <a-input-password
                name="passwd"
                :placeholder="$t('cluster.list.input.passwd')"
                v-decorator="['passwd',{rules: [{ required: true }]}]" />
            </a-form-item>
            <a-form-item
              :label="$t('cluster.list.description')"
              :labelCol="{lg: {span: 7}, sm: {span: 7}}"
              :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
              <a-textarea
                :auto-size="{ minRows: 4, maxRows: 6 }"
                name="description"
                :placeholder="$t('cluster.list.input.description')"
                v-decorator="['description']" />
            </a-form-item>
            <a-form-item
              :label="$t('cluster.list.checkItem')"
              :labelCol="{lg: {span: 7}, sm: {span: 7}}"
              :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
              <a-select
                mode="multiple"
                placeholder="please select default check item"
                v-decorator="['check_items',{rules: [{ required: true }]}]"
              >
                <a-select-option v-for="i in checkItems" :key="i" :value="i">
                  {{ i }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-form>
        </a-modal>
      </div>
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
            :title="false"
          >
            <div>
              <a-row type="flex" justify="space-between">
                <a-col>
                  {{ $t('cluster.list.name') }}
                </a-col>
                <a-col>
                  {{ $t('cluster.list.add-time') }} :  {{ item.create_time | moment }}
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
                    {{ item.cluster_name }}
                  </a>
                </div>
              </a-row>
              <div style="margin-top: 15px;margin-bottom: 5px">
                {{ $t('cluster.list.description') }}
              </div>
              <div>
                <div v-if="item.description">
                  <div>
                    {{ item.description }}
                  </div>
                </div>
                <div v-else>
                  {{ $t('cluster.list.description.else') }}
                </div>
              </div>
              <div style="margin-top: 15px">
                <a-row type="flex" justify="space-between">
                  <a-col>
                    <span>{{ $t('cluster.list.node-info') }}</span>
                  </a-col>
                  <a-col>
                    <span>
                      <a style="margin-right: 15px;color: #40a9ff" :href="item.grafana_url" target="_blank">Grafana</a>
                      <a style="color: #40a9ff" :href="item.dashboard_url" target="_blank">Dashboard</a>
                    </span>
                  </a-col>
                </a-row>
              </div>
            </div>
            <div>
              <a-list :grid="{ gutter: 16, column: 4 }" :data-source="item.nodes" style="margin-top: 25px" v-if="item.normal">
                <a-list-item slot="renderItem" slot-scope="node">
                  <a-card :title="node.type" style="text-align: center">
                    {{ node.normal }} / {{ node.count }}
                  </a-card>
                </a-list-item>
              </a-list>
              <div style="margin-top: 25px;" v-else>
                <a-card title="Warning" style="text-align: center;margin-bottom: 16px">
                  {{ $t('cluster.list.prometheus.warning') }}
                </a-card>
              </div>
            </div>
            <div>
              <a-row type="flex" justify="end">
                <a-col>
                  {{ $t('cluster.list.last-check-time') }}: {{ item.last_check_time | moment }}
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
import { getClusterList, addCluster } from '@/api/cluster'
import {
  ChartCard,
  RankList,
  Bar
} from '@/components'
// import { getClusterInfo } from '@/api/manage'

const dataSource = []
const checkItems = [
  'alive_pd_number',
  'alive_tidb_number',
  'alive_tikv_number',
  'available_memory',
  'failed_query_type',
  'long_ddl_job',
  'no_primary_key',
  'tidb_connections',
  'running_sql_5min',
  'tikv_region_number'
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
      checkItems,
      dataSource,
      modalVisible: false,
      clusterForm: this.$form.createForm(this)
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    getList () {
      getClusterList()
        .then(res => {
          this.dataSource = res.data
      }).catch(res => {
        this.failed(res)
      })
    },
    jump2Info (item) {
      this.$router.push({ name: 'ClusterInfo', params: { id: item.id } })
    },
    showModal () {
      this.modalVisible = true
    },
    modalCancel () {
      this.modalVisible = false
      this.clusterForm.resetFields()
    },
    handleOk () {
      this.clusterForm.validateFields((err, values) => {
        if (err) {
          this.failed(err)
        }
        console.log('values =>', values)
        addCluster(values)
        .then(res => {
          this.ifSuccess()
          this.modalVisible = false
          this.clusterForm = this.$form.createForm(this)
          }
        )
        .catch(res => this.failed(res))
        .finally(() => {
        })
      })
    },
    ifSuccess () {
      this.$notification.success({
        message: 'success',
        description: `success`
      })
      this.getList()
    },
    failed (res) {
      this.$notification['error']({
        message: 'error',
        description: res.error.msg,
        duration: 4
      })
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

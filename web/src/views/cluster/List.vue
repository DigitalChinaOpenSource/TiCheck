<template>
  <page-header-wrapper
    :title="$t('cluster.list.cluster-list')"
  >
    <!-- actions -->
    <template v-slot:extra>
      <div>
        <a-button type="primary" @click="showModal" >{{ $t('cluster.list.add-cluster') }}</a-button>
        <a-modal v-model="modalVisible" :title="$t('cluster.list.add-cluster')" @ok="handleOk" width="70%">
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
                  {{ $t('cluster.list.add-time') }} :  {{ item.create_time }}
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
              <a-list :grid="{ gutter: 16, column: 4 }" :data-source="item.nodes" style="margin-top: 25px">
                <a-list-item slot="renderItem" slot-scope="node">
                  <a-card :title="node.type" style="text-align: center">
                    {{ node.count }}
                  </a-card>
                </a-list-item>
              </a-list>
            </div>
            <div>
              <a-row type="flex" justify="end">
                <a-col>
                  {{ $t('cluster.list.last-check-time') }}: {{ item.last_check_time }}
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

export default {
  name: 'ClusterList',
  components: {
    ChartCard,
    RankList,
    Bar
  },
  data () {
    return {
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
      getClusterList().then(res => {
        this.dataSource = res.data
      })
    },
    jump2Info (item) {
      this.$router.push({ name: 'ClusterInfo', params: { id: item.id } })
    },
    showModal () {
      this.modalVisible = true
    },
    handleOk () {
      this.clusterForm.validateFields((err, values) => {
        if (err) {
          this.addFailed()
        }
        values.owner = 'test'
        addCluster(values)
        .then(res => this.addSuccess())
        .catch(res => this.addFailed())
        .finally(() => {
          this.modalVisible = false
          this.clusterForm = this.$form.createForm(this)
        })
      })
    },
    addSuccess () {
      this.$notification.success({
        message: 'add cluster success',
        description: `success`
      })
      this.getList()
    },
    addFailed () {
      this.$notification['error']({
        message: 'error',
        description: `error`,
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

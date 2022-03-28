<template>
  <div>
    <a-page-header
      :ghost="false"
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px;"
      :title="$t('cluster.list.setting.title')"
    />
    <a-form :form="updateForm">
      <a-form-item
        :label="$t('cluster.list.name')"
        :labelCol="{lg: {span: 7}, sm: {span: 7}}"
        :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
        <a-input
          v-decorator="['name',{rules: [{ required: true }], initialValue: this.initialCluster.name }]"
          :placeholder="$t('cluster.list.input.name')"
          name="name"/>
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
                }],
              initialValue: this.initialCluster.url
            }]"/>
      </a-form-item>
      <a-form-item
        :label="$t('cluster.list.user')"
        :labelCol="{lg: {span: 7}, sm: {span: 7}}"
        :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
        <a-input
          name="user"
          :placeholder="$t('cluster.list.input.user')"
          v-decorator="['user',{rules: [{ required: true }], initialValue: this.initialCluster.user}]"/>
      </a-form-item>
      <a-form-item
        :label="$t('cluster.list.passwd')"
        :labelCol="{lg: {span: 7}, sm: {span: 7}}"
        :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
        <a-input-password
          name="passwd"
          :placeholder="$t('cluster.list.input.passwd')"
          v-decorator="['passwd',{rules: [{ required: true }]}]"/>
      </a-form-item>
      <a-form-item
        :label="$t('cluster.list.description')"
        :labelCol="{lg: {span: 7}, sm: {span: 7}}"
        :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
        <a-textarea
          :auto-size="{ minRows: 4, maxRows: 6 }"
          name="description"
          :placeholder="$t('cluster.list.input.description')"
          v-decorator="['description', {initialValue: this.initialCluster.description}]"/>
      </a-form-item>
      <a-form-item
        label=" "
        :colon="false"
        :labelCol="{lg: {span: 7}, sm: {span: 7}}">
        <a-button type="primary" @click="showConfirm">{{ $t('cluster.list.setting.btn') }}</a-button>
      </a-form-item>
    </a-form>
  </div>
</template>

<script>
import { updateCluster, getInitialCluster } from '@/api/cluster'

const initialCluster = {}
export default {
  name: 'ClusterSetting',
  clusterID: '',
  data () {
    return {
      initialCluster,
      updateForm: this.$form.createForm(this)
    }
  },
  computed: {
    userInfo () {
      return this.$store.getters.userInfo
    }
  },
  methods: {
    showConfirm () {
      this.$confirm({
        title: this.$t('cluster.list.setting.confirm.title'),
        content: this.$t('cluster.list.setting.confirm.content'),
        onOk: () => {
          this.handleUpdate()
        },
        onCancel () {
        }
      })
    },
    handleUpdate () {
      this.updateForm.validateFields((err, values) => {
        if (err) {
          this.updateFailed()
        }
        updateCluster(this.clusterID, values)
          .then(res => this.updateSuccess())
          .catch(res => this.updateFailed())
          .finally(() => {
            this.InitialClusterInfo()
            this.updateForm = this.$form.createForm(this)
          })
      })
    },
    updateSuccess () {
      this.$notification.success({
        message: 'success',
        description: 'update'
      })
    },
    updateFailed () {
      this.$notification['error']({
        message: 'failed'
      })
    },
    InitialClusterInfo () {
      console.log('test')
      getInitialCluster(this.clusterID)
        .then(res => {
          this.initialCluster = res.data
        })
        .catch(() => {
          this.$router.push({ path: '/' })
        })
    }
  },
  created () {
    this.clusterID = this.$route.params.id
    this.InitialClusterInfo()
    console.log(this.initialCluster)
  }
}

</script>

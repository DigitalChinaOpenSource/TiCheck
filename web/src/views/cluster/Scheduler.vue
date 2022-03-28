<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px;"
      :title="$t('cluster.scheduler.title')"
    />
    <a-button type="primary" @click="showAddModal" style="float: right">
      {{ $t('cluster.scheduler.btn.add') }}
    </a-button>
    <a-modal v-model="scheModalVisible" :title="$t('cluster.scheduler.modal.title')" @ok="handleAddScheduler" @cancel="addCancel" width="70%">
      <a-form :form="schedulerForm">
        <a-form-item
          :label="$t('cluster.scheduler.modal.name')"
          :labelCol="{lg: {span: 5}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 15}, sm: {span: 17} }">
          <a-input
            v-decorator="['name',{rules: [{ required: true }]}]"
            :placeholder="$t('cluster.scheduler.modal.place.name')"
            name="name" />
        </a-form-item>
        <a-form-item
          :label="$t('cluster.scheduler.modal.cron')"
          :labelCol="{lg: {span: 5}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 15}, sm: {span: 17} }">
          <a-radio-group v-model="radioValue" @change="radioChange">
            <a-radio :value="1">
              {{ $t('cluster.scheduler.modal.radio.1') }}
            </a-radio>
            <a-radio :value="2">
              {{ $t('cluster.scheduler.modal.radio.2') }}
            </a-radio>
            <a-radio :value="3">
              {{ $t('cluster.scheduler.modal.radio.3') }}
            </a-radio>
          </a-radio-group>
          <a-input
            v-decorator="['cron',{rules: [{ required: true }],initialValue: '* * * * * ? *'}]"
            @change="radioInputChange"
            name="cron"/>
        </a-form-item>
        <a-form-item
          :label="$t('cluster.scheduler.modal.active')"
          :labelCol="{lg: {span: 5}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 15}, sm: {span: 17} }">
          <a-switch
            v-decorator="['status', {initialValue: true}]"
            :checked="isChecked"
            @change="switchOnChange"
            :checked-children="$t('cluster.scheduler.switch.child.yes')"
            :un-checked-children="$t('cluster.scheduler.switch.child.no')" />
        </a-form-item>
      </a-form>
    </a-modal>
    <a-table
      :columns="columns"
      :rowKey="(scheduler) => scheduler.id"
      :data-source="schedulerList"
      :pagination="paginationOpt"
      style="padding-top: 60px"
    >
      <span slot="status" slot-scope="status">
        {{ mapStatusValue(status) }}
      </span>
      <span slot="action" slot-scope="scheduler">
        <a @click="jump2History(scheduler.id)" target="_blank" style="color: #40a9ff">{{ $t('cluster.scheduler.table.action.history') }}</a>
      </span>
      <span slot="operation" slot-scope="text, record">
        <a @click="showEditModal(record)">{{ $t('cluster.scheduler.btn.edit') }}</a>
        <a-divider type="vertical" />
        <a-popconfirm
          :title="$t('cluster.scheduler.delete.msg')"
          :ok-text="$t('cluster.scheduler.delete.ok')"
          :cancel-text="$t('cluster.scheduler.delete.cancel')"
          @confirm="handleDelete(record)"
        >
          <a href="#">{{ $t('cluster.scheduler.btn.delete') }}</a>
        </a-popconfirm>
      </span>
    </a-table>
    <a-modal v-model="editModalVisible" :title="$t('cluster.scheduler.modal.edit.title')" @ok="handleEditScheduler" width="70%">
      <a-form :form="editSchedulerForm">
        <a-form-item>
          <span v-decorator="['id']"></span>
        </a-form-item>
        <a-form-item
          :label="$t('cluster.scheduler.modal.name')"
          :labelCol="{lg: {span: 5}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 15}, sm: {span: 17} }">
          <a-input
            v-decorator="['name',{rules: [{ required: true }]}]"
            :placeholder="$t('cluster.scheduler.modal.place.name')"
            name="name" />
        </a-form-item>
        <a-form-item
          :label="$t('cluster.scheduler.modal.cron')"
          :labelCol="{lg: {span: 5}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 15}, sm: {span: 17} }">
          <a-radio-group v-model="radioValue" @change="editRadioChange">
            <a-radio :value="1">
              {{ $t('cluster.scheduler.modal.radio.1') }}
            </a-radio>
            <a-radio :value="2">
              {{ $t('cluster.scheduler.modal.radio.2') }}
            </a-radio>
            <a-radio :value="3">
              {{ $t('cluster.scheduler.modal.radio.3') }}
            </a-radio>
          </a-radio-group>
          <a-input
            v-decorator="['cron',{rules: [{ required: true }],initialValue: '* * * * * ? *'}]"
            @change="editRadioInputChange"
            name="cron"/>
        </a-form-item>
        <a-form-item
          :label="$t('cluster.scheduler.modal.active')"
          :labelCol="{lg: {span: 5}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 15}, sm: {span: 17} }">
          <a-switch
            v-decorator="['status', {initialValue: true}]"
            :checked="isChecked"
            @change="editSwitchOnChange"
            :checked-children="$t('cluster.scheduler.switch.child.yes')"
            :un-checked-children="$t('cluster.scheduler.switch.child.no')" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script>
import { getSchedulerList, addScheduler, updateScheduler, deleteScheduler } from '@/api/cluster'
import moment from "moment";

const schedulerList = []

const paginationOpt = {
  showTotal: total => `Total ${total} items`,
  showSizeChanger: true,
  pageSizeOptions: ['10', '30', '50', '100'],
  defaultPageSize: 100
}

export default {
  name: 'ClusterScheduler',
  data () {
    return {
      schedulerList,
      paginationOpt,
      columns: [],
      clusterID: this.$route.params.id?.toString(),
      editModalVisible: false,
      scheModalVisible: false,
      isChecked: true,
      radioValue: 1,
      schedulerForm: this.$form.createForm(this, { name: 'addForm' }),
      editSchedulerForm: this.$form.createForm(this, { name: 'editForm' })
    }
  },
  methods: {
    getSchedulers () {
      getSchedulerList(this.clusterID).then(res => {
        this.schedulerList = res.data
        for (let i = 0; i < this.schedulerList.length; i++) {
          this.schedulerList[i].create_time = moment(this.schedulerList[i].create_time).format('YYYY-MM-DD HH:mm:ss')
        }
      })
    },
    mapStatusValue (status) {
      switch (status) {
        case 0:
          return this.$t('cluster.scheduler.table.status.no')
        case 1:
          return this.$t('cluster.scheduler.table.status.yes')
        default:
          return this.$t('cluster.scheduler.table.status.null')
      }
    },
    jump2History (schedulerID) {
      console.log('jump to history detail view', schedulerID)
    },
    showAddModal () {
      this.scheModalVisible = true
    },
    switchOnChange (checked) {
      this.isChecked = checked
      this.schedulerForm.setFieldsValue({
        status: checked
      })
    },
    radioChange () {
      if (this.radioValue === 1) {
        this.schedulerForm.setFieldsValue({
          cron: '* * * * * ? *'
        })
      }
      if (this.radioValue === 2) {
        this.schedulerForm.setFieldsValue({
          cron: '0 0 20 * * ?'
        })
      }
      if (this.radioValue === 3) {
        this.schedulerForm.setFieldsValue({
          cron: '0 0 0 ? * SUN'
        })
      }
    },
    radioInputChange () {
      this.radioValue = 1
    },
    editSwitchOnChange (checked) {
      this.isChecked = checked
      this.editSchedulerForm.setFieldsValue({
        status: checked
      })
    },
    editRadioChange () {
      if (this.radioValue === 1) {
        this.editSchedulerForm.setFieldsValue({
          cron: '* * * * * ? *'
        })
      }
      if (this.radioValue === 2) {
        this.editSchedulerForm.setFieldsValue({
          cron: '0 0 20 * * ?'
        })
      }
      if (this.radioValue === 3) {
        this.editSchedulerForm.setFieldsValue({
          cron: '0 0 0 ? * SUN'
        })
      }
    },
    editRadioInputChange () {
      this.radioValue = 1
    },
    handleAddScheduler () {
      this.schedulerForm.validateFields((err, values) => {
        if (err) {
          this.addFailed()
        }
        values.cluster_id = this.clusterID
        addScheduler(values)
        .then(res => this.addSuccess())
        .catch(res => this.addFailed())
        .finally(() => {
          this.scheModalVisible = false
          this.radioValue = 1
          this.isChecked = true
          this.schedulerForm.resetFields()
        })
      })
    },
    addCancel () {
      this.scheModalVisible = false
      this.schedulerForm.resetFields()
    },
    addFailed () {
      this.$notification['error']({
        message: 'error',
        description: `error`,
        duration: 4
      })
    },
    addSuccess () {
      this.$notification.success({
        message: 'add scheduler success',
        description: `success`
      })
      this.getSchedulers()
    },
    showEditModal (record) {
      setTimeout(() => {
        this.editSchedulerForm.setFieldsValue({
          'id': record.id,
          'name': record.name,
          'cron': record.cron
        })
      }, 100)
      this.editModalVisible = true
    },
    handleEditScheduler () {
      this.editSchedulerForm.validateFields((err, values) => {
        if (err) {
          this.addFailed()
        }
        values.cluster_id = this.clusterID
        updateScheduler(values)
          .then(res => this.addSuccess())
          .catch(res => this.addFailed())
          .finally(() => {
            this.editModalVisible = false
            this.isChecked = true
            this.radioValue = 1
            this.editSchedulerForm.resetFields()
          })
      })
    },
    handleDelete (record) {
      deleteScheduler(record.id)
        .then(() => {
        this.getSchedulers()
      })
        .catch(() => {
          this.$notification['error']({
            message: 'err',
            description: 'failed to delete',
            duration: 4
          })
        })
    }
  },
  mounted () {
    this.getSchedulers()
  },
  beforeUpdate () {
    this.columns = [
      {
        title: this.$t('cluster.scheduler.table.id'),
        dataIndex: 'index',
        key: 'index'
      },
      {
        title: this.$t('cluster.scheduler.table.name'),
        dataIndex: 'name',
        key: 'name'
      },
      {
        title: this.$t('cluster.scheduler.table.create'),
        dataIndex: 'create_time',
        key: 'create_time'
      },
      {
        title: this.$t('cluster.scheduler.table.cron'),
        dataIndex: 'cron',
        key: 'cron'
      },
      {
        title: this.$t('cluster.scheduler.table.status'),
        key: 'status',
        dataIndex: 'status',
        scopedSlots: { customRender: 'status' }
      },
      {
        title: this.$t('cluster.scheduler.table.count'),
        key: 'count',
        dataIndex: 'count'
      },
      {
        title: this.$t('cluster.scheduler.table.action'),
        key: 'action',
        scopedSlots: { customRender: 'action' }
      },
      {
        title: this.$t('cluster.scheduler.table.operation'),
        key: 'operation',
        scopedSlots: { customRender: 'operation' }
      }
    ]
  }
}

</script>

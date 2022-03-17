<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      title="$t('cluster.detail.check.probe.title')"
    />
    <a-button type="primary" @click="addProbe" style="float: right">
      {{ $t("cluster.detail.check.probe.add") }}
    </a-button>

    <a-table
      :columns="columns"
      :data-source="data"
      :rowKey="(record) => record.id"
      :pagination="paginationOpt"
      style="padding-top: 60px"
    >
      <div
        slot="filterDropdown"
        slot-scope="{
          setSelectedKeys,
          selectedKeys,
          confirm,
          clearFilters,
          column,
        }"
        style="padding: 8px"
      >
        <a-input
          v-ant-ref="(c) => (searchInput = c)"
          :placeholder="`Search ${column.dataIndex}`"
          :value="selectedKeys[0]"
          style="width: 188px; margin-bottom: 8px; display: block"
          @change="
            (e) => setSelectedKeys(e.target.value ? [e.target.value] : [])
          "
          @pressEnter="
            () => handleSearch(selectedKeys, confirm, column.dataIndex)
          "
        />
        <a-button
          type="primary"
          icon="search"
          size="small"
          style="width: 90px; margin-right: 8px"
          @click="() => handleSearch(selectedKeys, confirm, column.dataIndex)"
        >
          Search
        </a-button>
        <a-button
          size="small"
          style="width: 90px"
          @click="() => handleReset(clearFilters)"
        >
          Reset
        </a-button>
      </div>
      <span slot="operator" slot-scope="operator">
        {{ mapOperatorValue(operator) }}
      </span>
      <a-switch
        slot="is_enabled"
        slot-scope="value, record, index"
        :loading="loadingIds.includes(record.id)"
        :default-checked="mapEnableValue(value)"
        @change="() => onEnableChange(record, index)"
      />
      <span slot="action" slot-scope="text, record">
        <a @click="handleEdit(record)">Edit</a>
        <a-divider type="vertical" />
        <a @click="handleDelete(record)">Delete</a>
      </span>
    </a-table>

    <div>
      <a-modal
        title="Add Check Item"
        :visible="modalVisible"
        :confirm-loading="confirmLoading"
        :footer="null"
        @cancel="handleCancel"
      >
        <a-form
          :form="form"
          :label-col="{ span: 8 }"
          :wrapper-col="{ span: 12 }"
          @submit="(e) => handleUpdateProbe(e)"
        >
          <a-form-item label="Probe ID">
            <span class="ant-form-text" v-decorator="['probe_id', {}]">
              {{ form.getFieldValue("probe_id") }}
            </span>
          </a-form-item>
          <a-form-item label="Script Name">
            <span class="ant-form-text" v-decorator="['script_name', {}]">
              {{ form.getFieldValue("script_name") }}
            </span>
          </a-form-item>
          <a-form-item label="Operator">
            <a-select
              v-decorator="[
                'operator',
                {
                  rules: [
                    {
                      required: true,
                      message: 'Please select probe\'s operator',
                    },
                  ],
                },
              ]"
              placeholder="Select a operator"
            >
              <a-select-option value="0">
                {{ mapOperatorValue(0) }}
              </a-select-option>
              <a-select-option value="1">
                {{ mapOperatorValue(1) }}
              </a-select-option>
              <a-select-option value="2">
                {{ mapOperatorValue(2) }}
              </a-select-option>
              <a-select-option value="3">
                {{ mapOperatorValue(3) }}
              </a-select-option>
              <a-select-option value="4">
                {{ mapOperatorValue(4) }}
              </a-select-option>
              <a-select-option value="5">
                {{ mapOperatorValue(5) }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="Threshold">
            <a-input
              v-decorator="[
                'threshold',
                {
                  rules: [
                    { required: true, message: 'Please input probe threshold' },
                  ],
                },
              ]"
            />
          </a-form-item>
          <a-form-item :wrapper-col="{ span: 16, offset: 7 }">
            <a-button :loading="confirmLoading" html-type="submit">
              Submit
            </a-button>
            <a-button
              key="back"
              @click="handleCancel"
              style="margin-left: 50px"
            >
              Cannel
            </a-button>
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
  </div>
</template>
<script>
import {
  getProbeList,
  mapOperatorValue,
  mapEnableValue,
  changeProbeStatus,
  updateProbeConfig,
  deleteProbe,
} from "@/api/check";

const columns = [
  {
    title: "ID",
    dataIndex: "id",
    key: "id",
    hide: true,
  },
  {
    title: "probe_id",
    dataIndex: "probe_id",
    key: "probe_id",
    slots: { title: "customTitle" },
    scopedSlots: { customRender: "probe_id" },
  },
  {
    title: "script_name",
    dataIndex: "script_name",
    key: "script_name",
    scopedSlots: {
      filterDropdown: "filterDropdown",
      filterIcon: "filterIcon",
      customRender: "customRender",
    },
    onFilter: (value, record) =>
      record.script_name.toString().toLowerCase().includes(value.toLowerCase()),
    onFilterDropdownVisibleChange: (visible) => {
      if (visible) {
        setTimeout(() => {
          this.searchInput.focus();
        }, 0);
      }
    },
  },
  {
    title: "file_name",
    dataIndex: "file_name",
    key: "file_name",
  },
  {
    title: "tag",
    key: "tag",
    dataIndex: "tag",
    scopedSlots: { customRender: "tag" },
    filters: [
      { text: "集群", value: "集群" },
      { text: "网络", value: "网络" },
      { text: "运行状态", value: "运行状态" },
    ],
    filterMultiple: false,
    onFilter: (value, record) => record.tag.indexOf(value) === 0,
  },
  {
    title: "description",
    key: "description",
    dataIndex: "description",
  },
  {
    title: "operator",
    key: "operator",
    dataIndex: "operator",
    scopedSlots: { customRender: "operator" },
  },
  {
    title: "threshold",
    key: "threshold",
    dataIndex: "threshold",
  },
  {
    title: "is_enabled",
    key: "is_enabled",
    dataIndex: "is_enabled",
    scopedSlots: { customRender: "is_enabled" },
    filters: [
      { text: "disable", value: 0 },
      { text: "enable", value: 1 },
    ],
    filterMultiple: false,
    onFilter: (value, record) => record.is_enabled === value,
  },
  {
    title: "Action",
    key: "action",
    scopedSlots: { customRender: "action" },
  },
];

const data = [];

const paginationOpt = {
  showTotal: (total) => `Total ${total} items`,
  showSizeChanger: true,
  pageSizeOptions: ["10", "30", "50", "100"],
  defaultPageSize: 100,
};

export default {
  data() {
    return {
      data,
      columns,
      paginationOpt,
      start_time: "",
      end_time: "",
      clusterID: this.$route.params.id,
      loadingIds: [],

      // search
      searchText: "",
      searchInput: null,
      searchedColumn: "",

      // modal
      form: this.$form.createForm(this, { name: "probeForm" }),
      modalVisible: false,
      confirmLoading: false,
      editProbeRecord: {},
    };
  },
  activated() {
    this.clusterID = this.$route.params.id;
    this.getProbeListByClusterID();
  },
  methods: {
    getProbeListByClusterID() {
      getProbeList(this.clusterID).then((res) => {
        this.data = res.data;
      }).catch(err => {
        this.$router.push({ name: "cluster" });
      });
    },
    addProbe() {
      this.$router.push({ name: "ProbeAdd", params: { id: this.clusterID } });
    },
    onEnableChange(record, index) {
      console.log(index + "onEnableChange" + record.is_enabled);
      const status = record.is_enabled;
      this.loadingIds = [...this.loadingIds, record.id];
      changeProbeStatus({
        id: record.id,
        is_enabled: status === 1 ? 0 : 1,
      }).finally(() => {
        this.loadingIds.splice(this.loadingIds.indexOf(record.id), 1);
        this.getProbeListByClusterID();
      });
    },
    handleEdit(record) {
      this.modalVisible = true;
      this.editProbeRecord = record;
      this.$nextTick(() => {
        this.form.setFieldsValue({
          probe_id: record.probe_id,
          script_name: record.script_name,
          operator: mapOperatorValue(record.operator),
          threshold: record.threshold,
        });
      });
    },
    handleDelete(record) {
      deleteProbe(record.id).then(() => {
        this.getProbeListByClusterID();
      });
    },
    handleUpdateProbe(e) {
      e.preventDefault();
      this.confirmLoading = true;
      var a = {
        id: this.editProbeRecord.id,
        probe_id: this.form.getFieldValue("probe_id"),
        operator: parseInt(this.form.getFieldValue("operator")),
        threshold: this.form.getFieldValue("threshold"),
      };
      updateProbeConfig(a)
        .then((res) => {
          this.$message.success("Success");
          this.getProbeListByClusterID();

          setTimeout(() => {
            this.modalVisible = false;
            this.confirmLoading = false;
          }, 1000);
        })
        .catch((err) => {
          this.confirmLoading = false;
        });
    },
    handleCancel() {
      this.modalVisible = false;
    },
    mapOperatorValue(operator) {
      return mapOperatorValue(operator);
    },
    mapEnableValue(is_enabled) {
      return mapEnableValue(is_enabled);
    },

    // search
    handleSearch(selectedKeys, confirm, dataIndex) {
      confirm();
      this.searchText = selectedKeys[0];
      this.searchedColumn = dataIndex;
    },

    handleReset(clearFilters) {
      clearFilters();
      this.searchText = "";
    },
  },
};
</script>

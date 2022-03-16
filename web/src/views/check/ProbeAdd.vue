<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      title="$t('cluster.detail.check.probe.add.title')"
      @back="back"
    />

    <a-table
      :columns="columns"
      :data-source="data"
      :rowKey="(record) => record.id"
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
      <a-button
        type="primary"
        slot="action"
        slot-scope="record"
        @click="addProbe(record)"
      >
        {{ $t("add") }}
      </a-button>
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
          @submit="(e) => handleAddProbe(e)"
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
import { getAddProbeList, mapOperatorValue, addProbe } from "@/api/check";
const columns = [
  {
    title: "ID",
    dataIndex: "id",
    key: "id",
    slots: { title: "customTitle" },
    scopedSlots: { customRender: "id" },
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
    title: "operator",
    key: "operator",
    dataIndex: "operator",
    scopedSlots: { customRender: "operator" },
  },
  {
    title: "threshold",
    key: "threshold",
    dataIndex: "threshold",
    scopedSlots: { customRender: "threshold" },
  },
  {
    title: "creator",
    key: "creator",
    dataIndex: "creator",
  },
  {
    title: "update_time",
    key: "update_time",
    dataIndex: "update_time",
  },
  {
    title: "action",
    key: "action",
    scopedSlots: { customRender: "action" },
  },
];

const data = [];

export default {
  data() {
    return {
      data,
      columns,
      clusterID: this.$route.params.id,

      // search
      searchText: "",
      searchInput: null,
      searchedColumn: "",

      // modal
      form: this.$form.createForm(this, { name: "probeForm" }),
      modalVisible: false,
      confirmLoading: false,
    };
  },
  activated() {
    this.getAddProbeList();
  },
  methods: {
    getAddProbeList() {
      getAddProbeList(this.clusterID).then((res) => {
        this.data = res.data;
      });
    },
    back() {
      console.log("back");
      this.$router.go(-1);
    },
    mapOperatorValue(operator) {
      return mapOperatorValue(operator);
    },
    addProbe(record) {
      this.modalVisible = true;
      this.$nextTick(() => {
        this.form.setFieldsValue({
          probe_id: record.id,
          script_name: record.script_name,
          operator: mapOperatorValue(record.operator),
          threshold: record.threshold,
        });
      });
    },
    handleAddProbe(e) {
      e.preventDefault();
      this.confirmLoading = true;
      addProbe({
        cluster_id: parseInt(this.clusterID),
        probe_id: this.form.getFieldValue("probe_id"),
        operator: parseInt(this.form.getFieldValue("operator")),
        threshold: this.form.getFieldValue("threshold"),
      })
        .then((res) => {
          this.$message.success("Success");
          this.getAddProbeList();

          setTimeout(() => {
            this.confirmLoading = false;
            this.modalVisible = false;
          }, 1000);
        })
        .catch((err) => {
          this.confirmLoading = false;
        });
    },
    handleCancel() {
      this.modalVisible = false;
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

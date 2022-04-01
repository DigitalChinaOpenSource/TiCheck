<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      :title="$t('check.probe.add.title')"
      @back="back"
    />
    <!-- <div style="float: right">
      <a-radio-group default-value="local" button-style="solid" size="large">
        <a-radio-button value="local" @click="getLocalData">
          {{ $t("check.probe.add.local") }}
        </a-radio-button>
        <a-radio-button value="remote" @click="getRemoteData">
          {{ $t("check.probe.add.remote") }}
        </a-radio-button>
        <a-radio-button value="custom" @click="getCustomData">
          {{ $t("check.probe.add.custom") }}
        </a-radio-button>
      </a-radio-group>
    </div> -->
    <a-table
      :columns="columns"
      :data-source="data"
      :rowKey="(record) => record.id"
      :pagination="pagination"
      style="padding-top: 50px"
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
          {{ $t("check.probe.add.search") }}
        </a-button>
        <a-button
          size="small"
          style="width: 90px"
          @click="() => handleReset(clearFilters)"
        >
          {{ $t("check.probe.add.reset") }}
        </a-button>
      </div>
      <span slot="tag" slot-scope="tag">
        {{ mapTagText(tag) }}
      </span>
      <span slot="operator" slot-scope="operator">
        {{ mapOperatorValue(operator) }}
      </span>
      <span slot="update_time" slot-scope="update_time">
        {{ update_time | moment }}
      </span>
      <a-button
        type="primary"
        slot="action"
        slot-scope="record"
        @click="addProbe(record)"
      >
        {{ $t("check.probe.add.table.add") }}
      </a-button>
    </a-table>

    <div>
      <a-modal
        :title="$t('check.probe.module.title')"
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
          <a-form-item :label="$t('check.probe.module.id')">
            <span class="ant-form-text" v-decorator="['probe_id', {}]">
              {{ form.getFieldValue("probe_id") }}
            </span>
          </a-form-item>
          <a-form-item :label="$t('check.probe.module.name')">
            <span class="ant-form-text" v-decorator="['script_name', {}]">
              {{ form.getFieldValue("script_name") }}
            </span>
          </a-form-item>
          <a-form-item :label="$t('check.probe.module.operator')">
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
              :labelInValue="true"
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
          <a-form-item :label="$t('check.probe.module.threshold')">
            <a-input
              v-decorator="[
                'threshold',
                {
                  rules: [
                    { required: true, message: 'Please input a number' },
                  ],
                },
              ]"
            />
          </a-form-item>
          <a-form-item :wrapper-col="{ span: 16, offset: 7 }">
            <a-button :loading="confirmLoading" html-type="submit">
              {{ $t("check.probe.module.submit") }}
            </a-button>
            <a-button
              key="back"
              @click="handleCancel"
              style="margin-left: 50px"
            >
              {{ $t("check.probe.module.cancel") }}
            </a-button>
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
  </div>
</template>
<script>
import { getAddProbeList, mapOperatorValue, addProbe, mapTagText } from "@/api/check";

export default {
  data() {
    return {
      data: [],
      columns: [],
      pagination: {
        showTotal: (total) => `Total ${total} items`,
        showSizeChanger: true,
        pageSizeOptions: ["10", "20", "30", "40"],
      },
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
      getAddProbeList(this.clusterID)
        .then((res) => {
          this.data = res.data;
        })
        .catch((err) => {
          this.$router.push({ name: "cluster" });
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
          operator: {
            key: record.operator,
            label: mapOperatorValue(record.operator),
          },
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
        operator: parseInt(this.form.getFieldValue("operator").key),
        threshold: this.form.getFieldValue("threshold"),
      })
        .then((res) => {
          this.$message.success("Success");
          this.getAddProbeList();

          setTimeout(() => {
            this.confirmLoading = false;
            this.modalVisible = false;
          }, 500);
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

    mapTagText(tag) {
      return mapTagText(tag);
    },
  },
  beforeUpdate() {
    this.columns = [
      {
        title: this.$t("check.probe.add.table.probeID"),
        dataIndex: "id",
        key: "id",
        slots: { title: "customTitle" },
        scopedSlots: { customRender: "id" },
      },
      {
        title: this.$t("check.probe.add.table.name"),
        dataIndex: "script_name",
        key: "script_name",
        scopedSlots: {
          filterDropdown: "filterDropdown",
          filterIcon: "filterIcon",
          customRender: "customRender",
        },
        onFilter: (value, record) =>
          record.script_name
            .toString()
            .toLowerCase()
            .includes(value.toLowerCase()),
        onFilterDropdownVisibleChange: (visible) => {
          if (visible) {
            setTimeout(() => {
              this.searchInput.focus();
            }, 0);
          }
        },
      },
      {
        title: this.$t("check.probe.table.tag"),
        key: "tag",
        dataIndex: "tag",
        scopedSlots: { customRender: "tag" },
        filters: [
          { text: this.$t("check.probe.tag.cluster"), value: "cluster" },
          { text: this.$t("check.probe.tag.network"), value: "network" },
          { text: this.$t("check.probe.tag.running_state"), value: "running_state" },
          { text: this.$t("check.probe.tag.others"), value: "others" },
        ],
        filterMultiple: false,
        onFilter: (value, record) => record.tag.indexOf(value) === 0,
      },
      {
        title: this.$t("check.probe.add.table.operator"),
        key: "operator",
        dataIndex: "operator",
        scopedSlots: { customRender: "operator" },
      },
      {
        title: this.$t("check.probe.add.table.threshold"),
        key: "threshold",
        dataIndex: "threshold",
        scopedSlots: { customRender: "threshold" },
      },
      {
        title: this.$t("check.probe.add.table.updateTime"),
        key: "update_time",
        dataIndex: "update_time",
        scopedSlots: { customRender: "update_time" },
      },
      {
        title: this.$t("check.probe.add.table.action"),
        key: "action",
        scopedSlots: { customRender: "action" },
      },
    ];
  },
};
</script>

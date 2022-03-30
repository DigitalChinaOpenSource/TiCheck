<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      :title="$t('check.history.detail.title')"
      @back="back"
    />
    <a-button
      type="primary"
      :loading="loading"
      @click="download"
      style="float: right"
    >
      {{ $t("check.history.detail.download") }}
    </a-button>
    <a-table
      :columns="columns"
      :data-source="data"
      :rowKey="(record) => record.id"
      :pagination="paginationOpt"
      style="margin-top: 60px"
    >
      <span slot="check_tag" slot-scope="check_tag">
        {{ mapTagText(check_tag) }}
      </span>
      <span slot="check_status" slot-scope="check_status">
        <a-tag
          :key="check_status"
          :color="
            check_status === 2
              ? 'red'
              : check_status === 0
              ? 'green'
              : check_status === -1
              ? '#FF00FF' : 'volcano'
          "
        >
          {{ mapStatusText(check_status) }}
        </a-tag>
      </span>
      <span slot="operator" slot-scope="operator">
        {{ mapOperatorValue(operator) }}
      </span>
    </a-table>
  </div>
</template>

<script>
import { getReportDetail, downloadReport, mapTagText, mapStatusText } from "@/api/check";

const paginationOpt = {
  showTotal: (total) => `Total ${total} items`,
  showSizeChanger: true,
  pageSizeOptions: ["10", "30", "50", "100"],
  defaultPageSize: 100,
};

export default {
  data() {
    return {
      loading: false,
      reportID: this.$route.params.id,
      report: {},
      columns: [],
      data: [],
      paginationOpt,
    };
  },
  activated() {
    this.reportID = this.$route.params.id;
    this.getReport();
  },
  methods: {
    getReport() {
      getReportDetail(this.reportID).then((response) => {
        const result = response.data;
        this.data = result;
      });
    },
    back() {
      console.log("back");
      this.$router.go(-1);
    },
    download() {
      this.loading = true;
      console.log("downloading" + this.reportID);
      downloadReport(this.reportID).then((response) => {
        this.loading = false;
      });
    },
    mapStatusValue(status) {
      switch (status) {
        case -1:
          return "script error"
        case 0:
          return "normal";
        case 1:
          return "existing abnormal";
        case 2:
          return "new abnormal";
        default:
          return "unknown";
      }
    },
    mapOperatorValue(operator) {
      switch (operator) {
        case 0:
          return "NA";
        case 1:
          return "equal to";
        case 2:
          return "greater than";
        case 3:
          return "greater than or equal to";
        case 4:
          return "less than";
        case 5:
          return "less than or equal to";
        default:
          return "NA";
      }
    },

    mapTagText(tag) {
      return mapTagText(tag);
    },

    mapStatusText (status) {
      return mapStatusText(status)
    }
  },
  beforeUpdate() {
    this.columns = [
      {
        title: this.$t("check.history.detail.table.name"),
        dataIndex: "check_name",
        key: "check_name",
      },
      {
        title: this.$t("check.history.detail.table.tag"),
        dataIndex: "check_tag",
        key: "check_tag",
        scopedSlots: { customRender: "check_tag" },
        filters: [
          { text: this.$t("check.probe.tag.cluster"), value: "cluster" },
          { text: this.$t("check.probe.tag.network"), value: "network" },
          { text: this.$t("check.probe.tag.running_state"), value: "running_state" },
          { text: this.$t("check.probe.tag.others"), value: "others" },
        ],
        filterMultiple: false,
        onFilter: (value, record) => record.check_tag.indexOf(value) === 0,
      },
      {
        title: this.$t("check.history.detail.table.item"),
        dataIndex: "check_item",
        key: "check_item",
      },
      {
        title: this.$t("check.history.detail.table.status"),
        dataIndex: "check_status",
        key: "check_status",
        scopedSlots: { customRender: "check_status" },
        filters: [
          { text: this.$t("check.history.detail.table.status.normal"), value: 0 },
          { text: this.$t("check.history.detail.table.status.new_abnormal"), value: 2 },
          { text: this.$t("check.history.detail.table.status.existing_abnormal"), value: 1 },
          { text: this.$t("check.history.detail.table.status.script_error"), value: -1 }
        ],
        filterMultiple: false,
        onFilter: (value, record) =>
          record.check_status.toString().indexOf(value) === 0,
      },
      {
        title: this.$t("check.history.detail.table.operator"),
        dataIndex: "operator",
        key: "operator",
        scopedSlots: { customRender: "operator" },
      },
      {
        title: this.$t("check.history.detail.table.threshold"),
        dataIndex: "threshold",
        key: "threshold",
      },
      {
        title: this.$t("check.history.detail.table.actual"),
        dataIndex: "check_value",
        key: "check_value",
      },
      {
        title: this.$t("check.history.detail.table.duration")+"(ms)",
        dataIndex: "duration",
        key: "duration",
      },
    ]
  },
};
</script>

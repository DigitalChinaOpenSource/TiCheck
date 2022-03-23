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
      <span slot="check_status" slot-scope="check_status">
        <a-tag
          :key="check_status"
          :color="
            check_status === 2
              ? 'red'
              : check_status === 0
              ? 'green'
              : 'volcano'
          "
        >
          {{ mapStatusValue(check_status) }}
        </a-tag>
      </span>
      <span slot="operator" slot-scope="operator">
        {{ mapOperatorValue(operator) }}
      </span>
    </a-table>
  </div>
</template>

<script>
import { getReportDetail, downloadReport } from "@/api/check";

const columns = [];
const data = [];
const report = {};

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
      report,
      columns,
      data,
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
        case 0:
          return "normal";
        case 1:
          return "existing error";
        case 2:
          return "new error";
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
        filters: [
          { text: "集群", value: "集群" },
          { text: "网络", value: "网络" },
          { text: "运行状态", value: "运行状态" },
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
          { text: "normal", value: 0 },
          { text: "new error", value: 2 },
          { text: "existing error", value: 1 },
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
        title: this.$t("check.history.detail.table.duration"),
        dataIndex: "duration",
        key: "duration",
      },
    ]
  },
};
</script>

<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      :title="$t('check.history.title')"
    />
    <div style="float: right">
      <span>{{ $t('check.history.timeSelect') }}</span>
      <a-range-picker @change="onTimeChange" style="margin-left: 30px" />
    </div>

    <a-table
      :columns="columns"
      :data-source="data"
      :rowKey="(record) => record.id"
      :pagination="pagination"
      @change="handleChange"
      style="padding-top: 50px"
    >
      <span slot="check_time" slot-scope="check_time">
        {{ check_time | moment }}
      </span>
      <span slot="normal_items" slot-scope="normal_items">
        <a-tag :key="normal_items" :color="'green'">
          {{ normal_items }}
        </a-tag>
      </span>
      <span slot="warning_items" slot-scope="warning_items">
        <a-tag :key="warning_items" :color="'volcano'">
          {{ warning_items }}
        </a-tag>
      </span>
      <span slot="action" slot-scope="record">
        <a @click="getReportDetail(record.id)">
          {{ $t('check.history.table.detail') }}
        </a>
        <a-divider type="vertical" />
        <a @click="downloadReport(record)">
          {{ $t('check.history.table.download') }}
        </a>
      </span>
    </a-table>
  </div>
</template>
<script>
import { getCheckHistoryByClusterID, downloadReport } from "@/api/check";

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
      start_time: "",
      end_time: "",
      clusterID: this.$route.params.id,
    };
  },
  activated() {
    this.clusterID = this.$route.params.id;
    this.getHistoryList();
  },
  methods: {
    handleChange(pagination) {
      this.getHistoryList(pagination);
    },
    getHistoryList(
      pagination = {
        current: 1,
        pageSize: 10,
      }
    ) {
      this.pagination = pagination;
      getCheckHistoryByClusterID(
        this.clusterID,
        pagination.current,
        pagination.pageSize,
        this.start_time,
        this.end_time
      )
        .then((res) => {
          const result = res.data;
          this.data = result.data;
          this.pagination = {
            ...this.pagination,
            total: result.total,
          };
        })
        .catch((err) => {
          this.$router.push({ name: "cluster" });
        });
    },
    downloadReport(reportID) {
      downloadReport(reportID);
    },
    onTimeChange(date, dateString) {
      this.start_time = dateString[0];
      this.end_time = dateString[1];
      this.getHistoryList();
    },
    getReportDetail(reportID) {
      this.$router.push({ name: "ReportDetail", params: { id: reportID } });
    },
    // onChange: (current, pageSize) => {
    //   debugger

    //   this.getHistoryList(current, pageSize)
    // },
    // onShowSizeChange: (current, pageSize) => {
    //   this.pageSize = pageSize
    // }
  },
  beforeUpdate() {
    this.columns = [
      {
        title: this.$t("check.history.table.checkTime"),
        dataIndex: "check_time",
        key: "check_time",
        scopedSlots: { customRender: "check_time" },
      },
      {
        title: this.$t("check.history.table.duration"),
        dataIndex: "duration",
        key: "duration",
      },
      {
        title: this.$t("check.history.table.normal"),
        key: "normal_items",
        dataIndex: "normal_items",
        scopedSlots: { customRender: "normal_items" },
      },
      {
        title: this.$t("check.history.table.warning"),
        key: "warning_items",
        dataIndex: "warning_items",
        scopedSlots: { customRender: "warning_items" },
      },
      {
        title: this.$t("check.history.table.total"),
        key: "total_items",
        dataIndex: "total_items",
      },
      {
        title: this.$t("check.history.table.action"),
        key: "action",
        scopedSlots: { customRender: "action" },
      },
    ];
  },
};
</script>

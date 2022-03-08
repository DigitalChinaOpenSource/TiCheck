<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px;"
      title="cluster.detail.check.history"
    />
    <div style="float: right">
        <span>{{ $t('cluster.detail.check.history.timeSelect') }}</span>
        <a-range-picker @change="onTimeChange" style="margin-left: 20px" />
    </div>

    <a-table
      :columns="columns"
      :data-source="data"
      :rowKey="(record) => record.id"
      :pagination="pagination"
      @change="handleChange" 
      style="padding-top: 20px"
    >

      <a @click="getReportDetail(text)" slot="id" slot-scope="text">{{ text }}</a>
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
        <a @click="downloadReport(record.id)">cluster.detail.check.history.download</a>
      </span>
    </a-table>
  </div>
</template>
<script>
import { getCheckHistoryByClusterID, downloadReport } from '@/api/check';

const columns = [
  {
    title: "ID",
    dataIndex: "id",
    key: "id",
    slots: { title: "customTitle" },
    scopedSlots: { customRender: "id" },
  },
  {
    title: "Time",
    dataIndex: "check_time",
    key: "check_time",
  },
  {
    title: "Duration",
    dataIndex: "duration",
    key: "duration",
  },
  {
    title: "Normal",
    key: "normal_items",
    dataIndex: "normal_items",
    scopedSlots: { customRender: "normal_items" },
  },
  {
    title: "Warning",
    key: "warning_items",
    dataIndex: "warning_items",
    scopedSlots: { customRender: "warning_items" },
  },
  {
    title: "Total",
    key: "total_items",
    dataIndex: "total_items",
  },
  {
    title: "Action",
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
      pagination: {},
      start_time: "",
      end_time: "",
    };
  },

  watch: {
    current(val) {
      debugger;
    },
  },
  mounted() {
    this.getHistoryList();
  },
  methods: {
    handleChange(pagination) {
      this.pagination = pagination;
      this.getHistoryList(pagination)
    },
    getHistoryList(
      pagination = {
        current: 1,
        pageSize: 10,
      }
    ) {
      getCheckHistoryByClusterID(1, pagination.current, pagination.pageSize, this.start_time, this.end_time).then(
        (res) => {
          this.data = res.data;
          this.pagination = {
            ...this.pagination,
            total: res.total
          }
        }
      );
    },
    downloadReport(params) {
      console.log(params);
      downloadReport(params);
    },
    onTimeChange(date, dateString) {
      this.start_time = dateString[0];
      this.end_time = dateString[1];
      this.getHistoryList();      
    },
    getReportDetail(reportID) {
      console.log(reportID);
      this.$router.push({name: 'ReportDetail', params: { id: reportID }});
    }
    // onChange: (current, pageSize) => {
    //   debugger

    //   this.getHistoryList(current, pageSize)
    // },
    // onShowSizeChange: (current, pageSize) => {
    //   this.pageSize = pageSize
    // }
  },
};
</script>

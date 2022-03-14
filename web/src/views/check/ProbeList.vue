<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px;"
      title="$t('cluster.detail.check.probe.title')"
    />
    <a-button type="primary" :loading="loading" @click="addProbe" style="float: right">
      {{$t('cluster.detail.check.probe.add')}}
    </a-button>

    <a-table
      :columns="columns"
      :data-source="data"
      :rowKey="(record) => record.id"
      :pagination="paginationOpt"
      style="padding-top: 60px"
    >
      <span slot="operator" slot-scope="operator">
          {{ mapOperatorValue(operator) }}
      </span>
    </a-table>
  </div>
</template>
<script>
import { getProbeList } from '@/api/check';

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
    scopedSlots: { customRender: "tag" }
  },
  {
    title: "description",
    key: "description",
    dataIndex: "description"
  },
  {
    title: "operator",
    key: "operator",
    dataIndex: "operator",
    scopedSlots: { customRender: 'operator' }
  },
  {
    title: "threshold",
    key: "threshold",
    dataIndex: "threshold"
  },
  {
    title:"is_enabled",
    key:"is_enabled",
    dataIndex:"is_enabled",
    scopedSlots: { customRender: 'is_enabled' }
  }
];

const data = [];

const paginationOpt = {
  showTotal: total => `Total ${total} items`,
  showSizeChanger: true,
  pageSizeOptions: ['10', '30', '50', '100'],
  defaultPageSize: 100
}

export default {
  data() {
    return {
      data,
      columns,
      paginationOpt,
      start_time: "",
      end_time: "",
      clusterID: this.$route.params.id
    };
  },
  activated() {
    this.clusterID = this.$route.params.id;
    this.getProbeListByClusterID();
  },
  methods: {
    getProbeListByClusterID(){
      getProbeList(this.clusterID).then(
        (res) => {
          this.data = res.data;
        }
      );
    },
    addProbe() {
      this.$router.push({name: 'ProbeAdd', params: { id: this.clusterID }});
    },
    mapOperatorValue (operator) {
      switch (operator) {
        case 0:
          return 'NA'
        case 1:
          return 'equal to'
        case 2:
          return 'greater than'
        case 3:
          return 'greater than or equal to'
        case 4:
          return 'less than'
        case 5:
          return 'less than or equal to'
        default:
          return 'NA'
      }
    }
  }
}
</script>

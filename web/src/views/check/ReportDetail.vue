<template>
  <div>
    <a-page-header
        style="border: 1px solid rgb(235, 237, 240)"
        title="Title"
        sub-title="This is a subtitle"
        @back='back'
    />
  </div>
</template>

<script>
import { getReportDetail } from '@/api/check'

const report = {};

export default {
  data () {
    reportID = this.$route.params.id,
    report
  },
  mounted() {
    this.getReport();
  },
  methods: {
    getReport () {
      getReportDetail(this.reportID).then(response => {
        console.log(this.reportID);
        this.report = response
      })
    },
    back () {
      console.log('back')
      this.$router.go(-1)
    }
  }
}
</script>

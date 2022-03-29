<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      :title="$t('check.execute.title')"
    />
    <div style="float: left">
      <span> {{ $t("check.execute.times") }}: {{ checkTimes }} </span>
      <a-divider type="vertical" />
      <span>
        {{ $t("check.execute.lastTime") }}: {{ lastCheckTime | moment }}</span
      >
    </div>
    <div style="float: right">
      <a-button
        class="button"
        type="primary"
        :loading="isRunning"
        @click="confirmRunCheck"
      >
        {{ $t("check.execute.runCheck") }}
      </a-button>
      <a-button class="button" type="primary" @click="editProbe">
        {{ $t("check.execute.editProbe") }}
      </a-button>
      <a-button class="button" type="primary" :disabled="!isCompleted" @click="getResultReport">
        {{ $t("check.execute.checkResult") }}
      </a-button>
    </div>
    <div>
      <a-progress
        style="margin-top: 20px"
        :percent="getTotalPercent()"
        :status="progressStatus"
        :strokeWidth="12"
        :stroke-color="{
          from: '#108ee9',
          to: '#87d068',
        }"
      />
    </div>
    <!-- just display when there is no task -->
    <div v-if="!isRunning && !isCompleted">
      <div class="progress">
        <a-progress
          type="dashboard"
          :percent="getClusterPercent()"
          :width="200"
        >
          <template #format="">
            <div class="large">
              <span
                >{{ checkTags.cluster.current }}/{{ checkTags.cluster.total }}
              </span>
              <br /><br />
              <span> {{ $t("check.probe.tag.cluster") }} </span>
            </div>
          </template>
        </a-progress>
      </div>
      <div class="progress">
        <a-progress
          type="dashboard"
          :percent="getNetworkPercent()"
          :width="200"
        >
          <template #format="">
            <div class="large">
              <span
                >{{ checkTags.network.current }}/{{ checkTags.network.total }}
              </span>
              <br /><br />
              <span> {{ $t("check.probe.tag.network") }} </span>
            </div>
          </template>
        </a-progress>
      </div>
      <div class="progress">
        <a-progress type="dashboard" :percent="getStatePercent()" :width="200">
          <template #format="">
            <div class="large">
              <span
                >{{ checkTags.running_state.current }}/{{
                  checkTags.running_state.total
                }}
              </span>
              <br /><br />
              <span> {{ $t("check.probe.tag.running_state") }} </span>
            </div>
          </template>
        </a-progress>
      </div>
      <div class="progress">
        <a-progress type="dashboard" :percent="getOthersPercent()" :width="200">
          <template #format="">
            <div class="large">
              <span
                >{{ checkTags.others.current }}/{{ checkTags.others.total }}
              </span>
              <br /><br />
              <span> {{ $t("check.probe.tag.others") }} </span>
            </div>
          </template>
        </a-progress>
      </div>
    </div>
    <!-- just display when task is running-->
    <div style="float: left; margin-top: 20px" v-if="isRunning || isCompleted">
      <a-progress type="circle" :percent="getClusterPercent()" :width="100">
        <template #format="">
          <div class="small">
            <span>
              {{ checkTags.cluster.current }}/{{ checkTags.cluster.total }}
            </span>
            <br /><br />
            <span> {{ $t("check.probe.tag.cluster") }} </span>
          </div>
        </template>
      </a-progress>
      <br />
      <br />
      <a-progress type="circle" :percent="getNetworkPercent()" :width="100">
        <template #format="">
          <div class="small">
            <span>
              {{ checkTags.network.current }}/{{ checkTags.network.total }}
            </span>
            <br /><br />
            <span> {{ $t("check.probe.tag.network") }} </span>
          </div>
        </template>
      </a-progress>
      <br />
      <br />
      <a-progress type="circle" :percent="getStatePercent()" :width="100">
        <template #format="">
          <div class="small">
            <span
              >{{ checkTags.running_state.current }}/{{
                checkTags.running_state.total
              }}
            </span>
            <br /><br />
            <span> {{ $t("check.probe.tag.running_state") }} </span>
          </div>
        </template>
      </a-progress>
      <br />
      <br />
      <a-progress type="circle" :percent="getOthersPercent()" :width="100">
        <template #format="">
          <div class="small">
            <span
              >{{ checkTags.others.current }}/{{ checkTags.others.total }}
            </span>
            <br /><br />
            <span> {{ $t("check.probe.tag.others") }} </span>
          </div>
        </template>
      </a-progress>
    </div>

    <!-- data list -->
    <div
      style="padding-left: 140px; margin-top: 20px"
      v-if="isRunning || isCompleted"
    >
      <a-list :data-source="data">
        <a-list-item slot="renderItem" slot-scope="item">
          <a-list-item-meta :description="item.check_item">
            <span slot="title">{{ item.check_name }}</span>
          </a-list-item-meta>

          <a-list-item-meta>
            <span slot="title">{{ mapTagText(item.check_tag) }}</span>
          </a-list-item-meta>

          <a-list-item-meta>
            <span slot="title">{{
              $t("check.execute.list.check_threshold") +
              ": " +
              mapOperatorValue(item.operator) +
              " " +
              item.threshold
            }}</span>
          </a-list-item-meta>

          <a-list-item-meta>
            <span slot="title">{{
              $t("check.execute.list.result") + ": " + item.check_value
            }}</span>
          </a-list-item-meta>

          <a-icon :type="mapStatusIconType(item.check_status)" :style="mapStatusIconStyle(item.check_status)" />
        </a-list-item>
      </a-list>
    </div>

    <!-- modal -->
    <div>
      <a-modal
        v-model="visible"
        :title="$t('check.execute.modal.title')"
        @ok="runCheck"
      >
        <p>{{ $t("check.execute.modal.contents") }}</p>
      </a-modal>
    </div>
  </div>
</template>

<style>
.button {
  margin-left: 10px;
}
.progress {
  float: left;
  margin-top: 50px;
  margin-left: 10%;
}
.large {
  margin-top: 20px;
  font-size: 20px;
}
.small {
  margin-top: 10px;
  font-size: 12px;
}
</style>

<script>
import {
  getExecuteInfo,
  runExecute,
  mapOperatorValue,
  mapTagText,
} from "@/api/check";

const checkTags = {
  cluster: {
    total: 0,
    current: 0,
  },
  network: {
    total: 0,
    current: 0,
  },
  running_state: {
    total: 0,
    current: 0,
  },
  others: {
    total: 0,
    current: 0,
  },
};

export default {
  data() {
    return {
      current: 0,
      total: 0,
      checkTimes: 0,
      lastCheckTime: "xxxx-xx-xx",
      progressStatus: "active",
      data: [],
      visible: false,
      isRunning: false,
      isCompleted: false,
      clusterID: this.$route.params.id,
      reportID: 0,
      checkTags: checkTags,

      // use to record haved checked probe.
      probeID: [],
    };
  },
  activated() {
    this.clusterID = this.$route.params.id;
    this.getExecuteInfo();
  },
  methods: {
    getExecuteInfo() {
      console.log(this.checkTags);
      getExecuteInfo(this.clusterID)
        .then((res) => {
          const data = res.data;
          const checkItems = data.check_items;

          this.checkTimes = data.check_times;
          this.lastCheckTime = data.last_check_time;
          this.total = checkItems.total;
          if (this.total == 0) {
            this.$notification["info"]({
              title: this.$t("check.execute.notification.no_probe.title"),
              message: this.$t("check.execute.notification.no_probe"),
            });
          }
          this.checkTags.cluster.total = checkItems.tags.cluster
            ? checkItems.tags.cluster
            : 0;
          this.checkTags.network.total = checkItems.tags.network
            ? checkItems.tags.network
            : 0;
          this.checkTags.running_state.total = checkItems.tags.running_state
            ? checkItems.tags.running_state
            : 0;
          this.checkTags.others.total = checkItems.tags.others
            ? checkItems.tags.others
            : 0;
        })
        .catch((err) => {
          console.log(err);
          // this.$router.push({ name: "cluster" });
        });
    },
    confirmRunCheck() {
      if (this.total == 0) {
        this.visible = false;
        this.$notification["info"]({
          title: this.$t("check.execute.notification.no_probe.title"),
          message: this.$t("check.execute.notification.no_probe"),
        });
        return;
      }

      this.visible = true;
    },
    runCheck() {
      this.visible = false;
      this.current = 0;
      this.checkTags.cluster.current = 0;
      this.checkTags.network.current = 0;
      this.checkTags.running_state.current = 0;
      this.checkTags.others.current = 0;
      this.isCompleted = false;
      this.isRunning = true;
      this.progressStatus = "active";

      const webSocket = runExecute(this.clusterID);

      webSocket.onmessage = (event) => {
        const result = JSON.parse(event.data);
        console.log(result);

        if (result.is_finished == true) {
          this.isRunning = false;
          this.isCompleted = true;
          this.reportID = result.check_id;
          this.progressStatus = "success";
          this.$notification["success"]({
            message: "Success",
            description: this.$t("check.execute.notification.success"),
            duration: 3,
          });
          webSocket.close();
          return;
        }

        if (result.is_timeout == true) {
          this.isRunning = false;
          this.progressStatus = "exception";
          this.$notification["error"]({
            message: "Timeout",
            description: this.$t("check.execute.notification.timeout"),
            duration: 3,
          });
          webSocket.close();
          return;
        }

        if (result.err != null) {
          this.$notification["warning"]({
            message: "Warning",
            description: result.err,
            duration: 3,
          });
        }

        console.log(result.data);
        if (result.data == null || result.data.length < 1) {
          return;
        }

        this.data = result.data.concat(this.data);
        this.current = this.current + 1;
        const probeTag = result.data[0].check_tag;
        if (probeTag == "cluster") {
          this.checkTags.cluster.current += 1;
        } else if (probeTag == "network") {
          this.checkTags.network.current += 1;
        } else if (probeTag == "running_state") {
          this.checkTags.running_state.current += 1;
        } else {
          this.checkTags.others.current += 1;
        }
      };

      webSocket.onclose = (event) => {
        this.isRunning = false;
        if (!this.isCompleted) {
          this.isRunning = false;
          this.progressStatus = "exception";
          this.$notification["error"]({
            message: "Error",
            description: this.$t("check.execute.notification.error.disconnect"),
            duration: 3,
          });
        }
      };

      webSocket.onerror = (event) => {
        this.isRunning = false;
        this.progressStatus = "exception";
        this.$notification["error"]({
          message: "Error",
          description: this.$t("check.execute.notification.error"),
          duration: 3,
        });
      };
    },
    editProbe() {
      this.$router.push({
        name: "ProbeList",
        params: { id: this.clusterID },
      });
    },
    getResultReport() {
      this.$router.push({
        name: "ReportDetail",
        params: { id: this.reportID },
      });
    },
    getTotalPercent() {
      return this.total === 0
        ? 100
        : Math.floor((this.current / this.total) * 100);
    },
    getClusterPercent() {
      return this.checkTags.cluster.total === 0
        ? 100
        : Math.floor(
            (this.checkTags.cluster.current / this.checkTags.cluster.total) *
              100
          );
    },
    getNetworkPercent() {
      return this.checkTags.network.total === 0
        ? 100
        : Math.floor(
            (this.checkTags.network.current / this.checkTags.network.total) *
              100
          );
    },
    getStatePercent() {
      return this.checkTags.running_state.total === 0
        ? 100
        : Math.floor(
            (this.checkTags.running_state.current /
              this.checkTags.running_state.total) *
              100
          );
    },
    getOthersPercent() {
      return this.checkTags.others.total === 0
        ? 100
        : Math.floor(
            (this.checkTags.others.current / this.checkTags.others.total) * 100
          );
    },
    mapOperatorValue(operator) {
      return mapOperatorValue(operator);
    },
    mapTagText(tag) {
      return mapTagText(tag);
    },
    mapStatusIconType(status) {
      switch (status) {
        case -1:
          return "close-circle";
        case 0:
          return "check-circle";
        case 1:
          return "exclamation-circle"
        case 2:
          return "exclamation-circle";
        default:
          return "exclamation-circle";
      }
    },
    mapStatusIconStyle(status) {
      switch (status) {
        case -1:
          return "color: #FF0000; fontSize: 25px";
        case 0:
          return "color: #00FF00; fontSize: 25px";
        case 1:
          return "color: #FF8000; fontSize: 25px";
        case 2:
          return "color: #FF8000; fontSize: 25px";
        default:
          return "color: #FF8000; fontSize: 25px";
      }
    },
  },
};
</script>

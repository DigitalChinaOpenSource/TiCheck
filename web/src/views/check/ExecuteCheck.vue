<template>
  <div>
    <a-page-header
      style="border: 1px solid rgb(235, 237, 240); margin-bottom: 20px"
      :title="$t('check.execute.title')"
    />
    <div style="float: left">
      <span> {{ $t('check.execute.times') }}: {{ checkTimes }} </span>
      <a-divider type="vertical" />
      <span> {{ $t('check.execute.lastTime') }}: {{ lastTime }}</span>
    </div>
    <div style="float: right">
      <a-button class="button" type="primary" :loading="isRunning" @click="runCheck">
        Run Check
      </a-button>
      <a-button class="button" type="primary" @click="beginX">Edit Probe</a-button>
      <a-button class="button" type="primary" :disabled="!isCompleted"> Check Result </a-button>
    </div>
    <div>
      <a-progress
        style="margin-top: 20px"
        :percent="((current / total) * 100).toFixed(2)"
        :status="progressStatus"
        strokeWidth="12"
        :stroke-color="{
          from: '#108ee9',
          to: '#87d068',
        }"
      />
    </div>
    <!-- just display when there is no task -->
    <div v-show="!isRunning&&!isCompleted">
      <div class="progress">
        <a-progress
          type="dashboard"
          :percent="getClusterPercent()"
          :width="200"
        >
          <template #format="">
            <div class="large">
              <span>{{ checkTags.cluster.current }}/{{ checkTags.cluster.total }} </span>
              <br /><br />
              <span> {{ $t('check.probe.tag.cluster') }} </span>
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
              <span>{{ checkTags.network.current }}/{{ checkTags.network.total }} </span>
              <br /><br />
              <span> {{ $t('check.probe.tag.network') }} </span>
            </div>
          </template>
        </a-progress>
      </div>
      <div class="progress">
        <a-progress
          type="dashboard"
          :percent="getStatePercent()"
          :width="200"
        >
          <template #format="">
            <div class="large">
              <span>{{ checkTags.state.current }}/{{ checkTags.state.total }} </span>
              <br /><br />
              <span> {{ $t('check.probe.tag.state') }} </span>
            </div>
          </template>
        </a-progress>
      </div>
    </div>
    <!-- just display when task is running-->
    <div style="float: left; margin-top: 20px" v-show="isRunning||isCompleted">
      <a-progress type="circle" :percent="getClusterPercent()" :width="100">
        <template #format="">
          <div class="small">
            <span> {{ checkTags.cluster.current }}/{{ checkTags.cluster.total }} </span>
            <br /><br />
            <span> {{ $t('check.probe.tag.cluster') }} </span>
          </div>
        </template>
      </a-progress>
      <br />
      <br />
      <a-progress
        type="circle"
        :percent="getNetworkPercent()"
        :width="100"
      >
        <template #format="">
          <div class="small">
            <span> {{ checkTags.network.current }}/{{ checkTags.network.total }} </span>
            <br /><br />
            <span> {{ $t('check.probe.tag.network') }} </span>
          </div>
        </template>
      </a-progress>
      <br />
      <br />
      <a-progress
        type="circle"
        :percent="getStatePercent()"
        :width="100"
      >
        <template #format="">
          <div class="small">
            <span>{{ checkTags.state.current }}/{{ checkTags.state.total }} </span>
            <br /><br />
            <span> {{ $t('check.probe.tag.state') }} </span>
          </div>
        </template>
      </a-progress>
    </div>
    <div style="padding-left: 140px; margin-top: 20px" v-show="isRunning">
      <a-list item-layout="horizontal" :data-source="data">
        <a-list-item slot="renderItem" slot-scope="item">
          <a-list-item-meta
            description="Ant Design, a design language for background applications, is refined by Ant UED Team"
          >
            <a slot="title" href="https://www.antdv.com/">{{ item.title }}</a>
            <a-avatar
              slot="avatar"
              src="https://zos.alipayobjects.com/rmsportal/ODTLcjxAfvqbxHnVXCYX.png"
            />
          </a-list-item-meta>
        </a-list-item>
      </a-list>
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
const total = 0;
const current = 0;
const checkTimes = 0;
const lastTime = "xxxx-xx-xx";
const progressStatus = 'active';
const checkTags = {
  cluster: {
    total: 0,
    current: 0,
  },
  network: {
    total: 0,
    current: 0,
  },
  state: {
    total: 0,
    current: 0,
  },
};
const data = [];
export default {
  data() {
    return {
      current,
      total,
      checkTimes,
      lastTime,
      progressStatus,
      checkTags,
      data,
      isRunning: false,
      isCompleted: false,
    };
  },
  methods: {
    getClusterData() {
      this.current += 1;
      if (this.current === this.total) {
        this.current = this.total;
      }
      this.data.unshift({
        title: "Ant Design Title " + this.current,
      });
    },
    getTotalPercent() {
      return ((this.current / this.total) * 100).toFixed(2);
    },
    getClusterPercent() {
      return (
        (this.checkTags.cluster.current / this.checkTags.cluster.total) *
        100
      ).toFixed(2);
    },
    getNetworkPercent() {
      return (
        (this.checkTags.network.current / this.checkTags.network.total) *
        100
      ).toFixed(2);
    },
    getStatePercent() {
      return (
        (this.checkTags.state.current / this.checkTags.state.total) *
        100
      ).toFixed(2);
    },
    runCheck() {
      this.data = [];
      this.current = 0;
      this.checkTags.cluster.current = 0;
      this.checkTags.network.current = 0;
      this.checkTags.state.current = 0;
      this.isCompleted = false;
      this.isRunning = true;
      this.progressStatus = 'active';
      setTimeout(() => {
        this.progressStatus = 'success';
        this.isRunning = false;
        this.isCompleted = true;
      }, 3000);
    },
    executeCheck() {
      const id = this.$route.query.id;
      console.log('id=>', id)
      if (id != null) {
        this.runCheck()
        console.log('id is not null', id)
      }
    }
  },
  created() {
    window.that = this
    this.executeCheck()
  }
};
</script>

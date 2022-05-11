<template>
  <div class="page-header-index-wide">
    <a-card :bordered="false" :bodyStyle="{ padding: '16px 0', height: '100%' }" :style="{ height: '100%' }">
      <div class="account-settings-info-main" :class="{ 'mobile': isMobile }" :style="{ height: '100%' }">
        <div class="account-settings-info-left">
          <a-menu
            :mode="isMobile ? 'horizontal' : 'inline'"
            :style="{ border: '0', width: isMobile ? '560px' : 'auto'}"
            :selectedKeys="selectedKeys"
            :v-bind="clusrerID"
            type="inner"
            @openChange="onOpenChange"
          >
            <a-menu-item key="/cluster/info/detail/:id">
              <router-link :to="{ name: 'ClusterInfo', params: {id: clusrerID}}">
                <a-icon type="container" />
                {{ $t('menu.cluster.info') }}
              </router-link>
            </a-menu-item>
            <a-menu-item key="/cluster/check/history/:id">
              <router-link :to="{ name: 'CheckHistory', params: {id: clusrerID}}">
                <a-icon type="profile" />
                {{ $t('menu.cluster.check.history') }}
              </router-link>
            </a-menu-item>
            <a-menu-item key="/cluster/check/probe/:id">
              <router-link :to="{ name: 'ProbeList', params: {id: clusrerID}}">
                <a-icon type="block" />
                {{ $t('menu.cluster.check.probe') }}
              </router-link>
            </a-menu-item>
            <a-menu-item key="/cluster/check/execute/:id">
              <router-link :to="{ name: 'ExecuteCheck', params: {id: clusrerID}}">
                <a-icon type="rocket" />
                {{ $t('menu.cluster.check.execute') }}
              </router-link>
            </a-menu-item>
            <a-menu-item key="/cluster/info/scheduler/:id">
              <router-link :to="{ name: 'ClusterScheduler', params: {id: clusrerID}}">
                <a-icon type="schedule" />
                {{ $t('menu.cluster.scheduler') }}
              </router-link>
            </a-menu-item>
            <a-menu-item key="/cluster/info/setting/:id">
              <router-link :to="{ name: 'ClusterSetting', params: {id: clusrerID}}">
                <a-icon type="setting" />
                {{ $t('menu.cluster.settings') }}
              </router-link>
            </a-menu-item>
            <!-- <a-menu-item key="/cluster/user">
              <router-link :to="{ name: 'ClusterUser', params: {id: clusrerID}}">
                <a-icon type="folder-open" />
                {{ $t('menu.cluster.user') }}
              </router-link>
            </a-menu-item> -->
          </a-menu>
        </div>
        <div class="account-settings-info-right">
          <!-- <div class="account-settings-info-title">
            <span>{{ $t($route.meta.title) }}</span>
          </div> -->
          <route-view></route-view>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script>
import { RouteView } from '@/layouts'
import { baseMixin } from '@/store/app-mixin'

export default {
  components: {
    RouteView
  },
  mixins: [baseMixin],
  data () {
    return {
      openKeys: [],
      selectedKeys: [],
      clusrerID: this.$route.params.id
    }
  },
  mounted () {
    this.updateMenu()
  },
  methods: {
    onOpenChange (openKeys) {
      this.openKeys = openKeys
    },
    updateMenu () {
      const routes = this.$route.matched.concat()
      this.selectedKeys = [ routes.pop().path ]
      console.log(this.selectedKeys)
    },
    getCluterId () {
      return this.clusrerID
    }
  },
  watch: {
    '$route' (val) {
      this.updateMenu()
    }
  }
}
</script>

<style lang="less" scoped>
  .account-settings-info-main {
    width: 100%;
    display: flex;
    height: 100%;
    overflow: auto;

    &.mobile {
      display: block;

      .account-settings-info-left {
        border-right: unset;
        border-bottom: 1px solid #e8e8e8;
        width: 100%;
        height: 50px;
        overflow-x: auto;
        overflow-y: scroll;
      }
      .account-settings-info-right {
        padding: 20px 40px;
      }
    }

    .account-settings-info-left {
      border-right: 1px solid #f0f2f5;
      width: 224px;
      position: fixed;
      display: block;
      height: 88%;
    }

    .account-settings-info-right {
      flex: 1 1;
      padding: 8px 40px;
      display: block;
      margin-left: 225px;

      .account-settings-info-title {
        color: rgba(0,0,0,.85);
        font-size: 20px;
        font-weight: 500;
        line-height: 28px;
        margin-bottom: 12px;
      }
      .account-settings-info-view {
        padding-top: 12px;
      }
    }
  }

</style>

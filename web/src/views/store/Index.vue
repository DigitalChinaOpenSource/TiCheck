<template>
  <div class="page-header-index-wide">
    <a-card :bordered="false" :bodyStyle="{ padding: '16px 0', height: '100%' }" :style="{ height: '100%' }">
      <div class="menus-page-info-main" :class="{ 'mobile': isMobile }">
        <div class="menus-page-info-left">
          <a-menu
            :mode="isMobile ? 'horizontal' : 'inline'"
            :style="{ border: '0', width: isMobile ? '560px' : 'auto'}"
            :selectedKeys="selectedKeys"
            type="inner"
            @openChange="onOpenChange"
          >
            <a-menu-item key="/store/local">
              <router-link :to="{ name: 'StoreLocal' }">
                  <a-icon type="folder-open" />
                {{ $t('menu.store.local') }}
              </router-link>
            </a-menu-item>
            <a-menu-item key="/store/remote">
              <router-link :to="{ name: 'StoreRemote' }">
                  <a-icon type="cloud-server" />
                {{ $t('menu.store.remote') }}
              </router-link>
            </a-menu-item>
            <a-menu-item key="/store/custom">
              <router-link :to="{ name: 'StoreCustom' }">
                  <a-icon type="api" />
                {{ $t('menu.store.custom') }}
              </router-link>
            </a-menu-item>
          </a-menu>
        </div>
        <div class="menus-page-info-right">
          <!-- <div class="menus-page-info-title">
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
      selectedKeys: []
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
  .menus-page-info-main {
    width: 100%;
    display: flex;
    height: 100%;
    overflow: auto;

    &.mobile {
      display: block;

      .menus-page-info-left {
        border-right: unset;
        border-bottom: 1px solid #e8e8e8;
        width: 100%;
        height: 50px;
        overflow-x: auto;
        overflow-y: scroll;
      }
      .menus-page-info-right {
        padding: 20px 40px;
      }
    }

    .menus-page-info-left {
      border-right: 1px solid #f0f2f5;
      width: 224px;
      position: fixed;
      display: block;
      height: 88%;
    }

    .menus-page-info-right {
      flex: 1 1;
      padding: 8px 40px;
      display: block;
      margin-left: 225px;

      .menus-page-info-title {
        color: rgba(0,0,0,.85);
        font-size: 20px;
        font-weight: 500;
        line-height: 28px;
        margin-bottom: 12px;
        border-bottom: 1px solid #e8e8e8;
    border-radius: 2px 2px 0 0;
      }
      .menus-page-info-view {
        padding-top: 12px;
      }
    }
  }

</style>

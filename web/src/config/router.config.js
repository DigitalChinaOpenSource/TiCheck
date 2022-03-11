// eslint-disable-next-line
import { UserLayout, BasicLayout, MenuLayout } from '@/layouts'
// import { bxAnaalyse } from '@/core/icons'

const RouteView = {
  name: 'RouteView',
  render: h => h('router-view')
}

export const asyncRouterMap = [
  {
    path: '/',
    name: 'index',
    component: BasicLayout,
    meta: { title: 'menu.home' },
    redirect: '/cluster',
    children: [
      {
        path: '/cluster',
        name: 'cluster',
        redirect: '/cluster/list',
        component: RouteView,
        hideChildrenInMenu: true,
        meta: { title: 'menu.cluster', keepAlive: true, icon: 'cluster', permission: ['sham'] },
        children: [
          {
            path: '/cluster/list',
            name: 'clusterList',
            component: () => import('@/views/cluster/List'),
            meta: { title: 'menu.cluster.list', keepAlive: true, permission: ['sham'] }
          },
          {
            path: '/cluster/info/:id',
            name: 'clusterInfo',
            redirect: '/cluster/info/detail/:id',
            component: () => import('@/views/cluster/Index'),
            meta: { title: 'menu.cluster.info', keepAlive: true, permission: ['sham'] },
            children: [
              {
                path: '/cluster/info/detail/:id',
                name: 'ClusterInfo',
                component: () => import('@/views/cluster/Info'),
                meta: { title: 'menu.cluster.info.detail', keepAlive: true, permission: ['sham'] }
              },
              {
                path: '/cluster/info/setting/:id',
                name: 'ClusterSet',
                component: () => import('@/views/cluster/Setting'),
                meta: { title: 'menu.cluster.info.setting', keepAlive: true, permission: ['sham'] }
              },
              // check
              {
                path: '/cluster/check/history/:id',
                name: 'CheckHistory',
                component: () => import('@/views/check/History'),
                meta: { title: 'menu.cluster.check.history', keepAlive: true, permission: ['sham'] }
              },
              {
                path: '/cluster/check/report/:id',
                name: 'ReportDetail',
                component: () => import('@/views/check/ReportDetail'),
                meta: { title: 'menu.cluster.check.history.report', keepAlive: true, permission: ['sham'] }
              }
            ]
          }
        ]
      },
      // store
      {
        path: '/store',
        name: 'Store',
        redirect: '/store/local',
        component: () => import('@/views/store/Index'),
        // component: MenuLayout,
        meta: { title: 'menu.store', icon: 'appstore', permission: ['sham'] },
        hideChildrenInMenu: true,
        children: [
          {
            path: '/store/local',
            name: 'StoreLocal',
            component: () => import('@/views/store/Local'),
            meta: { title: 'menu.store.local', keepAlive: true, permission: ['sham'] }
          },
          {
            path: '/store/remote',
            name: 'StoreRemote',
            component: () => import('@/views/store/Remote'),
            meta: { title: 'menu.store.remote', keepAlive: true, permission: ['sham'] }
          },
          {
            path: '/store/custom',
            name: 'StoreCustom',
            component: () => import('@/views/store/Custom'),
            meta: { title: 'menu.store.custom', keepAlive: true, permission: ['sham'] }
          }
        ]
      },
      // account
      {
        path: '/account',
        component: RouteView,
        redirect: '/account/center',
        name: 'account',
        meta: { title: 'menu.account', icon: 'user', keepAlive: true, permission: ['sham'] },
        children: [
        ]
      },
      // list
      {
        path: '/list',
        name: 'list',
        component: RouteView,
        redirect: '/list/table-list',
        meta: { title: 'menu.setting', icon: 'setting', permission: ['sham'] },
        children: [
          {
            path: '/list/table-list/:pageNo([1-9]\\d*)?',
            name: 'TableListWrapper',
            hideChildrenInMenu: true, // 强制显示 MenuItem 而不是 SubMenu
            // component: () => import('@/views/list/TableList'),
            meta: { title: 'menu.setting.store', keepAlive: true, permission: ['sham'] }
          }
        ]
      }
    ]
  },
  {
    path: '*',
    redirect: '/404',
    hidden: true
  }
]

/**
 * 基础路由
 * @type { *[] }
 */
export const constantRouterMap = [
  {
    path: '/user',
    component: UserLayout,
    redirect: '/user/login',
    hidden: true,
    children: [
      {
        path: 'login',
        name: 'login',
        component: () => import(/* webpackChunkName: "user" */ '@/views/user/Login')
      },
      {
        path: 'register',
        name: 'register',
        component: () => import(/* webpackChunkName: "user" */ '@/views/user/Register')
      },
      {
        path: 'register-result',
        name: 'registerResult',
        component: () => import(/* webpackChunkName: "user" */ '@/views/user/RegisterResult')
      },
      {
        path: 'recover',
        name: 'recover',
        component: undefined
      }
    ]
  },

  {
    path: '/404',
    component: () => import(/* webpackChunkName: "fail" */ '@/views/exception/404')
  }
]

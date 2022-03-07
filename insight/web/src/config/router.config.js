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
            path: '/cluster/detail',
            name: 'clusterDetail',
            redirect: '/cluster/check/history',
            component: () => import('@/views/cluster/Index'),
            meta: { title: 'menu.cluster.detail', keepAlive: true, permission: ['sham'] },
            children: [
              // check
              {
                path: '/cluster/check/history',
                name: 'CheckHistory',
                component: () => import('@/views/check/History'),
                meta: { title: 'menu.cluster.check.history', keepAlive: true, permission: ['sham'] }
              },
              {
                path: '/cluster/check/history/detail',
                name: 'HistoryDetail',
                component: () => import('@/views/check/HistoryDetail'),
                meta: { title: 'menu.cluster.check.history.detail', keepAlive: true, permission: ['sham'] }
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

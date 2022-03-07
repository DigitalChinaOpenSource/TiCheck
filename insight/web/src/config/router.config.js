// eslint-disable-next-line
import { UserLayout, BasicLayout, BlankLayout } from '@/layouts'
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
        meta: { title: 'menu.dashboard', keepAlive: true, icon: 'cluster', permission: ['sham'] },
        children: [
          {
            path: '/cluster/list',
            name: 'cluster',
            component: () => import('@/views/cluster/List'),
            meta: { title: 'menu.dashboard.workplace', keepAlive: true, permission: ['sham'] }
          }
        ]
      },
      // store
      {
        path: '/store',
        redirect: '/store/local',
        component: RouteView,
        meta: { title: 'menu.form', icon: 'appstore', permission: ['sham'] },
        children: [
          {
            path: '/store/local',
            name: 'BaseForm',
            component: () => import('@/views/store/Local'),
            meta: { title: 'menu.form.basic-form', keepAlive: true, permission: ['sham'] }
          },
          {
            path: '/form/remote',
            name: 'StepForm',
            component: () => import('@/views/store/Remote'),
            meta: { title: 'menu.form.step-form', keepAlive: true, permission: ['sham'] }
          },
          {
            path: '/form/custome',
            name: 'AdvanceForm',
            component: () => import('@/views/store/Custom'),
            meta: { title: 'menu.form.advanced-form', keepAlive: true, permission: ['sham'] }
          }
        ]
      },
      {
      // check
        path: '/check',
        name: 'check',
        hideChildrenInMenu: true,
        redirect: '/check/history',
        component: RouteView,
        meta: { title: 'menu.check', icon: 'check', permission: ['sham'] },
        children: [
          {
            path: '/check/history',
            name: 'History',
            hideChildrenInMenu: true,
            component: () => import('@/views/check/History'),
            meta: { title: 'menu.check.history', keepAlive: true, permission: ['sham'] },
            children: [
              {
                path: '/check/history/detail',
                name: 'HistoryDetail',
                component: () => import('@/views/check/HistoryDetail'),
                meta: { title: 'menu.check.history.detail', keepAlive: true, permission: ['sham'] }
              }
            ]
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

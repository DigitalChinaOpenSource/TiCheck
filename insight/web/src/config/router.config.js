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
    redirect: '/dashboard',
    children: [
      // dashboard
      {
        path: '/dashboard',
        name: 'dashboard',
        redirect: '/cluster/list',
        component: RouteView,
            hideChildrenInMenu: true,
        meta: { title: 'menu.dashboard', keepAlive: true, icon: 'cluster', permission: ['dashboard'] },
        children: [
          {
            path: '/cluster/list',
            name: 'Workplace',
            component: () => import('@/views/cluster/List'),
            meta: { title: 'menu.dashboard.workplace', keepAlive: true, permission: ['dashboard'] }
          }
        ]
      },
      // store
      {
        path: '/store',
        redirect: '/store/local',
        component: RouteView,
        meta: { title: 'menu.form', icon: 'appstore', permission: ['form'] },
        children: [
          {
            path: '/store/local',
            name: 'BaseForm',
            component: () => import('@/views/store/Local'),
            meta: { title: 'menu.form.basic-form', keepAlive: true, permission: ['form'] }
          },
          {
            path: '/form/remote',
            name: 'StepForm',
            component: () => import('@/views/store/Remote'),
            meta: { title: 'menu.form.step-form', keepAlive: true, permission: ['form'] }
          },
          {
            path: '/form/custome',
            name: 'AdvanceForm',
            component: () => import('@/views/store/Custom'),
            meta: { title: 'menu.form.advanced-form', keepAlive: true, permission: ['form'] }
          }
        ]
      },
      // account
      {
        path: '/account',
        component: RouteView,
        redirect: '/account/center',
        name: 'account',
        meta: { title: 'menu.account', icon: 'user', keepAlive: true, permission: ['user'] },
        children: [
        ]
      },
      // list
      {
        path: '/list',
        name: 'list',
        component: RouteView,
        redirect: '/list/table-list',
        meta: { title: 'menu.setting', icon: 'setting', permission: ['table'] },
        children: [
          {
            path: '/list/table-list/:pageNo([1-9]\\d*)?',
            name: 'TableListWrapper',
            hideChildrenInMenu: true, // 强制显示 MenuItem 而不是 SubMenu
            // component: () => import('@/views/list/TableList'),
            meta: { title: 'menu.setting.store', keepAlive: true, permission: ['table'] }
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

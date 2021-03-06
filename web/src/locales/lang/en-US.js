import antdEnUS from 'ant-design-vue/es/locale-provider/en_US'
import momentEU from 'moment/locale/eu'
import global from './en-US/global'
import dictData from './en-US/dict'

import menu from './en-US/menu'
import setting from './en-US/setting'
import user from './en-US/user'

import dashboard from './en-US/dashboard'
import form from './en-US/form'
import result from './en-US/result'
import account from './en-US/account'
import cluster from './en-US/cluster'
import check from './en-US/check'

import store from './en-US/store'

const components = {
  antLocale: antdEnUS,
  momentName: 'eu',
  momentLocale: momentEU
}

export default {
  message: '-',

  'layouts.usermenu.dialog.title': 'Message',
  'layouts.usermenu.dialog.content': 'Are you sure you would like to logout?',
  'layouts.userLayout.title': 'TiDB automated checklist for hackathon 2021.',
  'layouts.list.load-more': 'Load More',
  'layouts.list.no-more-data': 'no more data',
  ...components,
  ...global,
  ...dictData,
  ...menu,
  ...setting,
  ...user,
  ...dashboard,
  ...form,
  ...result,
  ...account,
  ...store,
  ...cluster,
  ...check
}

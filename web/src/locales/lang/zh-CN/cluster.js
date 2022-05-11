import list from './cluster/list'
import info from './cluster/info'
import scheduler from './cluster/scheduler'

export default {
  ...list,
  ...info,
  ...scheduler
}

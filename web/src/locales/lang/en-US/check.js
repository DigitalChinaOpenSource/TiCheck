import history from './check/history'
import probe from './check/probe'
import execute from './check/execute'

export default {
  ...history,
  ...probe,
  ...execute
}

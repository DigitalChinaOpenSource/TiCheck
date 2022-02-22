<template>
  <global-footer class="footer custom-render">
    <template v-slot:links>
      <a href="https://github.com/DigitalChinaOpenSource/TiDB-for-PostgreSQL" target="_blank">TiDB For PostgreSQL</a>
      <a href="https://github.com/DigitalChinaOpenSource/TiCheck" target="_blank">TiCheck</a>
      <a href="https://www.zhihu.com/org/shen-zhou-shu-ma-yun-ji-di" target="_blank">Contact Us</a>
    </template>
    <template v-slot:copyright>
      <a href="https://www.zhihu.com/org/shen-zhou-shu-ma-yun-ji-di" target="_blank">神州数码武汉云基地</a>
    </template>
  </global-footer>
</template>

<script>
import { GlobalFooter } from '@ant-design-vue/pro-layout'

export default {
  name: 'ProGlobalFooter',
  components: {
    GlobalFooter
  }
}
</script>

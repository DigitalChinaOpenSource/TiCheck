<template>
  <div>
    <a-card :bordered="false" class="ant-pro-components-tag-select" :title="$t('store.page.remote.title')">
         <div slot="extra">
        <a-radio-group default-value="all" @change="onSearch">
          <a-radio-button value="all">全部分类</a-radio-button>
          <a-radio-button value="processing">集群</a-radio-button>
          <a-radio-button value="waiting">网络</a-radio-button>
          <a-radio-button value="run">运行状态</a-radio-button>
        </a-radio-group>
        <a-input-search style="margin-left: 16px; width: 272px;" @search="onSearch" />
      </div>
    </a-card>

    <a-card style="margin-top: 10px;" :bordered="false">
      <a-list
        size="large"
        rowKey="id"
        :loading="loading"
        itemLayout="vertical"
        :dataSource="data"
      >
        <a-list-item :key="item.id" slot="renderItem" slot-scope="item, index">
          <template slot="actions">
            <span> <a-icon type="star-o" style="margin-right: 8px" /> {{ item.star }} </span>
            <span> <a-icon type="like-o" style="margin-right: 8px" /> {{ item.like }} </span>
            <span> <a-icon type="message" style="margin-right: 8px" /> {{ item.message }} </span>
          </template>
          <a-list-item-meta>
            <a slot="title" target="_blank" href="https://github.com/DigitalChinaOpenSource/TiCheck_ScriptWarehouse/tree/main/scripts">
            <a-avatar size="large" :src=" index%2==0?python:shell " />
            {{ item.name }}</a>
            <template slot="description">
              <span>
                <a-tag color="blue">集群</a-tag>
                <a-tag color="blue">网络</a-tag>
              </span>
            </template>
          </a-list-item-meta>
          <article-list-content :description="item.description" :owner="item.owner" :avatar="item.avatar" :href="item.href" :updateAt="item.updatedAt" />
        </a-list-item>
        <div slot="footer" v-if="data.length > 0" style="text-align: center; margin-top: 16px;">
          <a-button @click="loadMore" :loading="loadingMore">加载更多...</a-button>
        </div>
      </a-list>
    </a-card>
  </div>
</template>

<script>
import { TagSelect, ArticleListContent } from '@/components'
// import IconText from './components/IconText'
const TagSelectOption = TagSelect.Option

export default {
  components: {
    TagSelect,
    TagSelectOption,
    ArticleListContent
  },
  data () {
    return {
      python: require('@/assets/icons/python.png'),
      shell: require('@/assets/icons/shell.png'),
      loading: true,
      loadingMore: false,
      data: []
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    onSearch (value) {
      alert(`selected ${value}`)
    },
    getList () {
      this.$http.get('/store/remote').then(res => {
        console.log('res', res)
        this.data = res.script_list
        this.loading = false
      })
    },
    loadMore () {
      this.loadingMore = true
      this.$http.get('/list/article').then(res => {
        this.data = this.data.concat(res.result)
      }).finally(() => {
        this.loadingMore = false
      })
    }
  }
}
</script>

<style lang="less" scoped>

</style>

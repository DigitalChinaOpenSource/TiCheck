<template>
  <div>
    <a-card :bordered="false" class="ant-pro-components-tag-select" :title="$t('store.page.local.title')">
         <div slot="extra">
        <a-radio-group default-value="all" @change="onSearch">
          <a-radio-button value="all">全部分类</a-radio-button>
          <a-radio-button value="集群">集群</a-radio-button>
          <a-radio-button value="网络">网络</a-radio-button>
          <a-radio-button value="运行状态">运行状态</a-radio-button>
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
          <a-list-item-meta>
            <a slot="title" v-on:click="showReadme(item.id)" >
            <a-avatar size="large" v-if="item.file_name.split('.')[1]=='sh'" :src="icon.shell" />
            <a-avatar size="large" v-if="item.file_name.split('.')[1]=='py'" :src="icon.python" />
            {{ item.script_name }}</a>
            <template slot="description">
              <span>
                <a-tag color="blue">{{ item.tag }}</a-tag>
              </span>
            </template>
          </a-list-item-meta>
          <article-list-content :description="item.description" :owner="item.creator" avatar="" href="" :updateAt="item.update_time" />
        </a-list-item>
        <div slot="footer" v-if="data.length > 0" style="text-align: center; margin-top: 16px;">
          <a-button @click="loadMore" v-if="showMore" :loading="loadingMore">加载更多</a-button>
          <p v-if="!showMore" class="ant-result-subtitle">----没有更多数据了----</p>
        </div>
      </a-list>
    </a-card>

    <a-modal
      :title="$t('store.page.local.modal.title')"
      width="880px"
      :visible="visible"
      @cancel="handleCancel"
    >
      <div class="markdown-body">
        <vue-markdown :source="readmeText" v-highlight></vue-markdown>
      </div>
    </a-modal>
   
  </div>
</template>

<script>
import { ArticleListContent } from '@/components'
import VueMarkdown from 'vue-markdown'
import 'github-markdown-css/github-markdown.css'

export default {
  components: {
    ArticleListContent,
    VueMarkdown
  },
  data () {
    return {
      icon: {
          python: require('@/assets/icons/python.png'),
          shell: require('@/assets/icons/shell.png')
      },
      msg:'',
      key: 0,
      readmeText: '',
      visible: false,
      loading: true,
      showMore: true,
      loadingMore: false,
      data: [],
      page:0,
      search:{
        tag:'all',
        name:''
      }
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    showReadme (id) {
      this.readmeText = ''
      this.visible = true
      this.$http.get('/store/local/readme?name='+id).then((response) => {
         　　this.readmeText = response;
     　　});
    },
    handleCancel (e) {
      this.visible = false
    },
    onSearch (value) {
      this.data =[]
      this.page=0
      if (value.target){
        this.search.tag=value.target.value
      }else{
        this.search.name=value
      }
      this.getList()
    },
    getList () {
      this.page++
      this.loading = true
      this.showMore = true
      this.$http.get(
        '/store/local?page='+this.page+'&page_size=10&tag='+this.search.tag+"&name="+this.search.name).then(res => {
        this.data = this.data.concat(res.rows)
        this.loading = false
        this.loadingMore = false
        if (res.total <= this.data.length) {
          this.showMore = false
        }
      })
    },
    loadMore () {
      this.loadingMore = true
      this.getList()
    }
  }
}
</script>

<style lang="less" scoped>
.ant-pro-components-tag-select {
  /deep/ .ant-pro-tag-select .ant-tag {
    margin-right: 24px;
    padding: 0 8px;
    font-size: 14px;
  }
}

.list-articles-trigger {
  margin-left: 12px;
}
</style>

<template>
  <div>
    <a-card :bordered="false" class="ant-pro-components-tag-select" title="Remote Store">
         <div slot="extra">
        <a-radio-group default-value="all" @change="onSearch">
          <a-radio-button value="all">全部分类</a-radio-button>
          <a-radio-button value="processing">集群</a-radio-button>
          <a-radio-button value="waiting">网络</a-radio-button>
          <a-radio-button value="run">运行状态</a-radio-button>
        </a-radio-group>
        <a-input-search style="margin:0 16px; width: 272px;" @search="onSearch" />
        <a-button type="primary" icon="upload" @click="showModal"> Upload </a-button>
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
            <a slot="title" >
            <a-avatar size="large" :src=" index%2==0?python:shell " />
            {{ item.title }}</a>
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

    <a-modal
      title="Upload Custom Script"
      width="680px"
      :visible="visible"
      :confirm-loading="confirmLoading"
      @ok="handleOk"
      @cancel="handleCancel"
    >
      <p>{{ ModalText }}</p>
    </a-modal>
  </div>
</template>

<script>
import { TagSelect, StandardFormRow, ArticleListContent } from '@/components'
// import IconText from './components/IconText'
const TagSelectOption = TagSelect.Option

export default {
  components: {
    TagSelect,
    TagSelectOption,
    StandardFormRow,
    ArticleListContent
  },
  data () {
    return {
      python: require('@/assets/icons/python.png'),
      shell: require('@/assets/icons/shell.png'),
      visible: false,
      confirmLoading: false,
      loading: true,
      loadingMore: false,
      data: [],
      form: this.$form.createForm(this)
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
      this.$http.get('/list/article').then(res => {
        console.log('res', res)
        this.data = res.result
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
    },
    showModal () {
      this.visible = true
    },
    handleOk (e) {
      this.ModalText = 'The modal will be closed after two seconds'
      this.confirmLoading = true
      setTimeout(() => {
        this.visible = false
        this.confirmLoading = false
      }, 2000)
    },
    handleCancel (e) {
      console.log('Clicked cancel button')
      this.visible = false
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

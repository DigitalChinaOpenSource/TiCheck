<template>
<page-header-wrapper>
  <div  class="page-header-wrapper">
    <div class="ant-pro-pages-list-applications-filterCardList">
      <a-list
        :loading="loading"
        :data-source="data"
        :grid="{ gutter: 24, xl: 4, lg: 3, md: 3, sm: 2, xs: 1 }"
        style="margin-top: 24px;"
        :pagination="{showSizeChanger: true, pageSize: 5, total: 50, showTotal: (total, range) => `total ${total} items`}"
      >
        <a-list-item slot="renderItem" slot-scope="item, index">
          <a-card :body-style="{ paddingBottom: 20 }" hoverable>
            <a-card-meta :title="item.title">
              <template slot="avatar">
                <a-avatar size="large" :src=" index%2==0?python:shell " />
              </template>
            </a-card-meta>
                <a-tag color="blue">
        节点状态
      </a-tag>
            <template slot="actions">
              <a-tooltip title="下载">
                <a-icon type="download" />
              </a-tooltip>
              <a-tooltip title="编辑">
                <a-icon type="edit" />
              </a-tooltip>
              <a-tooltip title="分享">
                <a-icon type="share-alt" />
              </a-tooltip>
              <a-dropdown>
                <a class="ant-dropdown-link">
                  <a-icon type="ellipsis" />
                </a>
                <a-menu slot="overlay">
                  <a-menu-item>
                    <a href="javascript:;">1st menu item</a>
                  </a-menu-item>
                  <a-menu-item>
                    <a href="javascript:;">2nd menu item</a>
                  </a-menu-item>
                  <a-menu-item>
                    <a href="javascript:;">3rd menu item</a>
                  </a-menu-item>
                </a-menu>
              </a-dropdown>
            </template>
            <div class="cardInfo">
    <div>
      <p>脚本描述脚本描述脚本描述脚本描述脚本描述脚本描述脚本描述脚本描述脚本描述</p>
    </div>
            </div>
          </a-card>
        </a-list-item>
      </a-list>
    </div>
  </div>
   </page-header-wrapper>
</template>

<script>
import moment from 'moment'
import { TagSelect, StandardFormRow, Ellipsis, AvatarList } from '@/components'
// import CardInfo from './components/CardInfo'
const TagSelectOption = TagSelect.Option
const AvatarListItem = AvatarList.Item

export default {
  components: {
    AvatarList,
    AvatarListItem,
    Ellipsis,
    TagSelect,
    TagSelectOption,
    StandardFormRow
  },
  data () {
    return {
      python: require('@/assets/icons/python.png'),
      shell: require('@/assets/icons/shell.png'),
      data: [],
      form: this.$form.createForm(this),
      loading: true
    }
  },
  filters: {
    fromNow (date) {
      return moment(date).fromNow()
    }
  },
  mounted () {
    this.getList()
  },
  methods: {
    handleChange (value) {
      console.log(`selected ${value}`)
    },
    getList () {
      this.$http.get('/list/article', { params: { count: 8 } }).then(res => {
        console.log('res', res)
        this.data = res.result
        this.loading = false
      })
    }
  }
}
</script>

<style lang="less" scoped>
@import "~@/components/index.less";
@import "~@/utils/utils.less";
.ant-pro-components-tag-select {
  /deep/ .ant-pro-tag-select .ant-tag {
    margin-right: 24px;
    padding: 0 8px;
    font-size: 14px;
  }
}
.ant-pro-pages-list-projects-cardList {
  margin-top: 24px;

  /deep/ .ant-card-meta-title {
    margin-bottom: 4px;
  }

  /deep/ .ant-card-meta-description {
    height: 44px;
    overflow: hidden;
    line-height: 22px;
  }

  .cardItemContent {
    display: flex;
    height: 20px;
    margin-top: 16px;
    margin-bottom: -4px;
    line-height: 20px;

    > span {
      flex: 1 1;
      color: rgba(0, 0, 0, 0.45);
      font-size: 12px;
    }

    /deep/ .ant-pro-avatar-list {
      flex: 0 1 auto;
    }
  }
}
.cardInfo {

  .clearfix();

  margin-top: 16px;
  margin-left: 55px;
  & > div {
    position: relative;
    float: left;
    width: 100%;
    text-align: left;
    p:first-child {
      margin-bottom: 4px;
      color: @text-color-secondary;
      font-size: 12px;
      line-height: 20px;
    }
  }
}
</style>

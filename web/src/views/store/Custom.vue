<template>
  <div>
    <a-card
      :bordered="false"
      class="ant-pro-components-tag-select"
      :title="$t('store.page.custom.title')"
    >
      <div slot="extra">
        <a-radio-group default-value="all" @change="onSearch">
          <a-radio-button value="all">全部分类</a-radio-button>
          <a-radio-button value="集群">集群</a-radio-button>
          <a-radio-button value="网络">网络</a-radio-button>
          <a-radio-button value="运行状态">运行状态</a-radio-button>
        </a-radio-group>
        <a-input-search
          style="margin: 0 16px; width: 272px"
          @search="onSearch"
        />
        <a-button type="primary" icon="file-add" @click="showModal">
          {{ $t("store.page.custom.upload") }}
        </a-button>
      </div>
    </a-card>

    <a-card style="margin-top: 10px" :bordered="false">
      <a-list
        size="large"
        rowKey="id"
        :loading="loading"
        itemLayout="vertical"
        :dataSource="data"
      >
        <a-list-item :key="item.id" slot="renderItem" slot-scope="item, index">
          <a-list-item-meta>
            <a slot="title">
              <a-avatar
                size="large"
                v-if="item.file_name.split('.')[1] == 'sh'"
                :src="icon.shell"
              />
              <a-avatar
                size="large"
                v-if="item.file_name.split('.')[1] == 'py'"
                :src="icon.python"
              />
              {{ item.script_name }}</a
            >
            <template slot="description">
              <span>
                <a-tag color="blue">{{ item.tag }}</a-tag>
              </span>
            </template>
          </a-list-item-meta>
          <article-list-content
            :description="item.description"
            :owner="item.creator"
            avatar=""
            href=""
            :updateAt="item.update_time"
          />
        </a-list-item>
        <div
          slot="footer"
          v-if="data.length > 0"
          style="text-align: center; margin-top: 16px"
        >
          <a-button @click="loadMore" v-if="showMore" :loading="loadingMore"
            >加载更多</a-button
          >
          <p v-if="!showMore" class="ant-result-subtitle">
            ----没有更多数据了----
          </p>
        </div>
      </a-list>
    </a-card>

    <a-modal
      :title="$t('store.page.custom.modal.title')"
      width="680px"
      :visible="visible"
      :confirm-loading="confirmLoading"
      @ok="handleOk"
      @cancel="handleCancel"
    >
      <a-upload
        accept=".zip"
        :disabled="disabled"
        :multiple="false"
        :file-list="fileList"
        :remove="removeFile"
        :before-upload="
          (file) => {
            if (this.fileList.length > 0) {
              this.$message.error('只能上传一个文件');
              return false;
            }
            if (file.size > 1024 * 1024 * 10) {
              this.$message.error('文件大小不能超过10M');
              return false;
            }
            const isZip = file.type.includes('zip');
            if (!isZip) {
              this.$message.error('请上传zip文件');
              return false;
            }
            this.fileList = [...this.fileList, file];
            return false;
          }
        "
      >
        <a-button>
          <a-icon type="upload" /> {{ $t("store.page.custom.modal.upload") }}
        </a-button>
      </a-upload>
      <div class="desc">
        <div
          style="
            font-size: 16px;
            color: rgba(0, 0, 0, 0.85);
            font-weight: 500;
            margin-bottom: 16px;
          "
        >
          {{ $t("result.fail.error.hint-title") }}
        </div>
        <div style="margin-bottom: 16px">
          <a-icon
            type="close-circle-o"
            style="color: #f5222d; margin-right: 8px"
          />
          {{ $t("result.fail.error.hint-text1") }}
          <a style="margin-left: 16px"
            >{{ $t("result.fail.error.hint-btn1") }} <a-icon type="right"
          /></a>
        </div>
        <div>
          <a-icon
            type="close-circle-o"
            style="color: #f5222d; margin-right: 8px"
          />
          {{ $t("result.fail.error.hint-text2") }}
          <a style="margin-left: 16px"
            >{{ $t("result.fail.error.hint-btn2") }} <a-icon type="right"
          /></a>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script>
import { TagSelect, ArticleListContent } from "@/components";
import request from "@/utils/request";
// import IconText from './components/IconText'
const TagSelectOption = TagSelect.Option;

export default {
  components: {
    TagSelect,
    TagSelectOption,
    ArticleListContent,
  },
  data() {
    return {
      icon: {
        python: require("@/assets/icons/python.png"),
        shell: require("@/assets/icons/shell.png"),
      },
      visible: false,
      confirmLoading: false,
      loading: true,
      loadingMore: false,
      data: [],
      page: 0,
      search: {
        tag: "all",
        name: "",
      },
      fileList: [],
    };
  },
  mounted() {
    this.getList();
  },
  methods: {
    onSearch(value) {
      this.data = [];
      this.page = 0;
      if (value.target) {
        this.search.tag = value.target.value;
      } else {
        this.search.name = value;
      }
      this.getList();
    },
    getList() {
      this.page++;
      this.loading = true;
      this.showMore = true;
      this.$http
        .get(
          "/store/custom?page=" +
            this.page +
            "&page_size=10&tag=" +
            this.search.tag +
            "&name=" +
            this.search.name
        )
        .then((res) => {
          this.data = this.data.concat(res.rows);
          this.loadingMore = false;
          if (res.total <= this.data.length) {
            this.showMore = false;
          }
        })
        .finally(() => {
          this.loading = false;
        });
    },
    loadMore() {
      this.loadingMore = true;
      this.getList();
    },
    showModal() {
      this.visible = true;
    },
    handleOk(e) {
      this.confirmLoading = true;

      const { fileList } = this;
      const formData = new FormData();
      formData.append("file", fileList[0]);

      request({
        url: "/store/custom",
        method: "post",
        processData: false,
        data: formData,
      })
        .then((res) => {
          this.visible = false;
          this.$message.success(`ile uploaded successfully`);
        })
        .catch((err) => {
          this.$message.error(`upload failed.`);
        })
        .finally(() => {
          this.confirmLoading = false;
        });
    },
    handleCancel(e) {
      this.visible = false;
      this.fileList =[]; 
    },
    removeFile(file) {
      const index = this.fileList.indexOf(file);
      const newFileList = this.fileList.slice();
      newFileList.splice(index, 1);
      this.fileList = newFileList;
    },
  },
};
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

.desc {
  margin-top: 24px;
  padding: 24px 24px;
  background-color: #fafafa;
}
</style>

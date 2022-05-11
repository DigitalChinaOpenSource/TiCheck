<template>
  <div>
    <a-card
      :bordered="false"
      class="ant-pro-components-tag-select"
      :title="$t('store.page.custom.title')"
    >
      <div slot="extra">
        <a-radio-group default-value="all" @change="onSearch">
          <a-radio-button value="all">{{$t("check.probe.tag.all")}}</a-radio-button>
          <a-radio-button value="cluster">{{$t("check.probe.tag.cluster")}}</a-radio-button>
          <a-radio-button value="network">{{$t("check.probe.tag.network")}}</a-radio-button>
          <a-radio-button value="running_state">{{$t("check.probe.tag.running_state")}}</a-radio-button>
          <a-radio-button value="others">{{$t("check.probe.tag.others")}}</a-radio-button>
        </a-radio-group>
        <a-input-search
          style="margin: 0 16px; width: 272px"
          @search="onSearch"
        />
        <a-button type="primary" ghost icon="file-add" @click="showUploadModal">
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
            <a slot="title" v-on:click="showReadme(item.id)">
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
                <a-tag color="blue">{{ mapTag(item.tag) }}</a-tag>
              </span>
            </template>
          </a-list-item-meta>
          <!-- <article-list-content
            :description="item.description"
            :owner="item.creator"
            avatar=""
            href=""
            :updateAt="item.update_time"
          /> -->
          <div
            class="antd-pro-components-article-list-content-index-listContent"
          >
            <div class="description">
              <slot>
                {{ item.description }}
              </slot>
            </div>
            <div class="extra">
              <!-- <a-avatar :src="avatar" size="small" /> -->
              <a>{{ item.creator }}</a> {{ $t("store.page.update-time") }}
              <!-- <a :href="href">{{ href }}</a> -->
              <em>{{ item.update_time | moment }}</em>
              <div class="actions">
                <a-popconfirm
                  :title="$t('store.page.custom.delete-tip')"
                  @confirm="deleteOne(item.id)"
                >
                  <a-button type="link" icon="delete"> Delete </a-button>
                </a-popconfirm>
              </div>
            </div>
          </div>
        </a-list-item>
        <div
          slot="footer"
          v-if="data.length > 0"
          style="text-align: center; margin-top: 16px"
        >
          <a-button @click="loadMore" v-if="showMore" :loading="loadingMore">{{
            $t("layouts.list.load-more")
          }}</a-button>
          <p v-if="!showMore" class="ant-result-subtitle">
            ---- {{ $t("layouts.list.no-more-data") }} ----
          </p>
        </div>
      </a-list>
    </a-card>
    <a-modal
      :title="$t('store.page.custom.readme.title')"
      width="880px"
      :visible="readme.visible"
      :footer="null"
      @cancel="closeReadme"
    >
      <div class="markdown-body">
        <vue-markdown :source="readme.text" v-highlight></vue-markdown>
      </div>
    </a-modal>
    <a-modal
      :title="$t('store.page.custom.upload.modal-title')"
      width="680px"
      :visible="upload.visible"
      :confirm-loading="upload.confirmLoading"
      @ok="handleUploadOk"
      @cancel="handleUploadCancel"
    >
      <a-upload
        accept=".zip"
        :multiple="false"
        :file-list="upload.fileList"
        :remove="removeFile"
        :before-upload="
          (file) => {
            if (this.upload.fileList.length > 0) {
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
            this.upload.fileList = [...this.upload.fileList, file];
            return false;
          }
        "
      >
        <a-button>
          <a-icon type="upload" />
          {{ $t("store.page.custom.upload.select-file") }}
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
          {{ $t("store.page.custom.upload.tips-title") }}:
        </div>
        <div style="margin-bottom: 16px">
          <a-icon
            type="exclamation-circle-o"
            style="color: #f5222d; margin-right: 8px"
          />
          {{ $t("store.page.custom.upload.tips-text1") }}
        </div>
        <div style="margin-bottom: 16px">
          <a-icon
            type="exclamation-circle-o"
            style="color: #f5222d; margin-right: 8px"
          />
          {{ $t("store.page.custom.upload.tips-text2") }}
        </div>
        <div>
          <a-icon
            type="exclamation-circle-o"
            style="color: #f5222d; margin-right: 8px"
          />
          {{ $t("store.page.custom.upload.tips-text3") }}
          <a
            style="margin-left: 16px"
            href="https://github.com/DigitalChinaOpenSource/TiCheck"
            target="_blank"
            >{{ $t("store.page.custom.upload.tips-more") }}
            <a-icon type="right"
          /></a>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script>
import { ArticleListContent } from "@/components";
import {  mapTagText } from "@/api/check";
import request from "@/utils/request";
import VueMarkdown from "vue-markdown";
import "github-markdown-css/github-markdown.css";

export default {
  components: {
    VueMarkdown,
    ArticleListContent,
  },
  data() {
    return {
      icon: {
        python: require("@/assets/icons/python.png"),
        shell: require("@/assets/icons/shell.png"),
      },
      upload: {
        visible: false,
        confirmLoading: false,
        fileList: [],
      },
      readme: {
        text: "",
        visible: false,
      },
      loading: true,
      loadingMore: false,
      data: [],
      page: 0,
      search: {
        tag: "all",
        name: "",
      },
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
          this.data = this.data.concat(res.data.rows);
          this.loadingMore = false;
          if (res.data.total <= this.data.length) {
            this.showMore = false;
          }
        })
        .finally(() => {
          this.loading = false;
        });
    },
    mapTag(tag) {
      return mapTagText(tag);
    },
    loadMore() {
      this.loadingMore = true;
      this.getList();
    },
    deleteOne(id) {
        request({
        url: "/store/custom/" + id,
        method: "delete",
        }).then((res) => {
          if (res.success) {
            this.data = this.data.filter((item) => item.id !== id);
          }
        })
        .catch((err) => {});
    },
    showUploadModal() {
      this.upload.visible = true;
    },
    handleUploadOk(e) {
      this.upload.confirmLoading = true;

      const formData = new FormData();
      formData.append("file", this.upload.fileList[0]);

      request({
        url: "/store/custom",
        method: "post",
        processData: false,
        data: formData,
      })
        .then((res) => {
          if (res.success) {
            this.upload.visible = false;
            this.$message.success(res.msg);
            this.onSearch("");
          }
        })
        .catch((err) => {})
        .finally(() => {
          this.upload.confirmLoading = false;
        });
    },
    handleUploadCancel(e) {
      this.upload.visible = false;
      this.upload.fileList = [];
    },
    removeFile(file) {
      const index = this.upload.fileList.indexOf(file);
      const newFileList = this.upload.fileList.slice();
      newFileList.splice(index, 1);
      this.upload.fileList = newFileList;
    },
    showReadme(id) {
      this.readme.text = "";
      this.readme.visible = true;
      this.$http.get("/store/custom/readme/" + id).then((response) => {
        this.readme.text = response.data;
      });
    },
    closeReadme(e) {
      this.readme.visible = false;
    },
  },
};
</script>

<style lang="less" scoped>
@import "../../components/index.less";

.antd-pro-components-article-list-content-index-listContent {
  .description {
    // max-width: 720px;
    line-height: 22px;
  }
  .extra {
    margin-top: 16px;
    color: @text-color-secondary;
    line-height: 22px;

    & /deep/ .ant-avatar {
      position: relative;
      top: 1px;
      width: 20px;
      height: 20px;
      margin-right: 8px;
      vertical-align: top;
    }

    & > em {
      margin-left: 16px;
      color: @disabled-color;
      font-style: normal;
    }

    .actions {
      display: inline;
      float: right;
    }
  }
}

@media screen and (max-width: @screen-xs) {
  .antd-pro-components-article-list-content-index-listContent {
    .extra {
      & > em {
        display: block;
        margin-top: 8px;
        margin-left: 0;
      }
    }
  }
}

.desc {
  margin-top: 24px;
  padding: 24px 24px;
  background-color: #fafafa;
}
</style>

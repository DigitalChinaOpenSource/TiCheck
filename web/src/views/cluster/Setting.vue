<template>
  <a-page-header
    :ghost="false"
    :style="{ marginTop: '24px'}"
    :title="$t('cluster.setting.title')">
    <div>
      <a-form :form="clusterForm">
        <a-form-item
          :label="$t('cluster.list.name')"
          :labelCol="{lg: {span: 2}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
          <a-input
            v-decorator="['name',{rules: [{ required: true }]}]"
            :placeholder="$t('cluster.list.input.name')"
            name="name" />
        </a-form-item>
        <a-form-item
          :label="$t('cluster.list.prometheus')"
          :labelCol="{lg: {span: 2}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
          <a-input
            name="url"
            addon-before="http://"
            :placeholder="$t('cluster.list.input.prometheus')"
            v-decorator="[
              'url',
              {
                rules: [
                  {
                    required: true,
                    message: ''
                  }]
              }]" />
        </a-form-item>
        <a-form-item
          :label="$t('cluster.list.user')"
          :labelCol="{lg: {span: 2}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
          <a-input
            name="user"
            :placeholder="$t('cluster.list.input.user')"
            v-decorator="['user',{rules: [{ required: true }]}]" />
        </a-form-item>
        <a-form-item
          :label="$t('cluster.list.passwd')"
          :labelCol="{lg: {span: 2}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
          <a-input-password
            name="passwd"
            :placeholder="$t('cluster.list.input.passwd')"
            v-decorator="['passwd',{rules: [{ required: true }]}]" />
        </a-form-item>
        <a-form-item
          :label="$t('cluster.list.description')"
          :labelCol="{lg: {span: 2}, sm: {span: 7}}"
          :wrapperCol="{lg: {span: 10}, sm: {span: 17} }">
          <a-textarea
            :auto-size="{ minRows: 4, maxRows: 6 }"
            name="description"
            :placeholder="$t('cluster.list.input.description')"
            v-decorator="['description']" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleOk">update</a-button>
        </a-form-item>
      </a-form>
    </div>
  </a-page-header>
</template>

<script>

export default {
  name: 'ClusterSet'
}

</script>

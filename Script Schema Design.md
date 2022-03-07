### 组件设计规范

每个组件由以下3种文件构成：

- *元信息文件(package.json)，描述了组件的基本信息

- *脚本文件(python或者shell)，至少要有一个主文件

- 说明文档(readme.md)，产品使用手册、原理解析、参数、安全风险等



#### package.json
```
{
  "_id": "alive_pd_number", //id要和打包名称一致
  "name": "存活的PD节点数量",
  "author": {
    "name": "DigitalChina",
    "email": "heao@digitalchina.com"
  },
  "description": "this script description.",
  "files": [
    "index.js",
    "lib/"
  ],
  "main": "index.js", //入口文件
  "tags": [
    "集群",
    "网络",
    "运行状态"
  ],
  "rules": [
      {
          "operator": 3, //0-无，1-等于，2-大于，3-大于等于，4-小于，5-小于等于
          "threshold": "3",
          "args": []
      }
  ],
  "homepage": "",
  "version": "1.0.1",
  "createTime": "",
  "updateTime": ""
}

```
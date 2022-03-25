### 组件设计规范

每个组件由以下3种文件构成：

- *元信息文件(package.json)，描述了组件的基本信息

- *脚本文件(python或者shell)，至少要有一个主文件

- 说明文档(readme.md)，产品使用手册、原理解析、参数、安全风险等



#### package.json格式
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
    "cluster",
    "network",
    "running_state",
    "others"
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


### 输入输出规范

我们为每一个Probe制定了统一的输入输出格式，会把一些全局参数传给要运行的脚本，方然也支持给脚本自定义启动参数。与此同时，脚本输出信息也要符合一定要求才能被TiCheck捕获到，这和是否能正确判断阈值至关重要。

#### 默认输入参数

针对`shell`和`python`脚本，我们会统一传入三个参数，他们依次是：

- BasePath，程序运行的主目录
- MysqlLoginPath，tidb集群的登录连接信息，可以使用login path的方式登录执行sql
- PremetheusPath，tidb集群的Prometheus地址，可以使用Psql查询集群监控指标

#### 自定义输入参数

#### 巡检值输出格式
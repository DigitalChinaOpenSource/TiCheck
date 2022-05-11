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

- BasePath，比如`/data/tidb/ticheck`，程序运行的主目录
- MysqlLoginPath，比如`tidb-login`，tidb集群的登录连接信息，可以使用login path的方式登录执行sql
- PremetheusPath，比如`http://10.0.0.1:9090`，tidb集群的Prometheus地址，可以使用Psql查询集群监控指标

> **注意:**
> 这三个参数需要按传入顺序获取，比如在shell中依次是`$1`、`$2`、`$3`，在python中依次是`sys.argv[1]`、`sys.argv[2]`、`sys.argv[3]`。

#### 自定义输入参数

脚本自定义参数在`package.json`的`rules.args`中设置，可以设置多个自定义参数，目前不支持使用`--参数名=参数值`或者`--参数名 参数值`的形式，**每个参数只填写参数值即可，这些参数会追加到默认参数后面传入，也就是说它们的顺序在3以后，获取方式参考前面的描述**。

#### 巡检值输出格式

TiCheck能够识别两种格式的输出信息。

第一种是本次巡检的实际值，这个结果会用来和设置的阈值做对比，得出巡检结果。它的格式是`[tck_result:]xxxx=yy`。参考示例：
```
// shell
echo "[tck_result:] TiDB节点数量=2"
echo "$i" | awk '{print [tck_result:] $1"."$2=无主键}'

// python
print ("[tck_result:] TiDB节点数量=2")
```
每一行仅输出一个巡检项实际值，有多个巡检项的情况下请分多行输出，例如需要检查每一个节点的网络流量的时候每一个节点都应该是一行单独的输出。

第二种是脚本需要记录的日志，这些信息通常用于跟踪脚本执行逻辑是否符合预期，它的格式是`[tck_log:]xxxxx`，目前只能在go控制台中打印出来，后面会考虑持久化到数据库中，给巡检结果提供查看入口。
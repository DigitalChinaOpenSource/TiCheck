# TiCheck

TiDB automated checklist for hackathon 2021.

## 团队

项目成员都来自神州数码 TiDB 团队，TiDB For PostgreSQL 项目核心开发者。

- [hoho](https://github.com/hey-hoho)
- [david](https://github.com/AmoebaProtozoa)
- [pupillord](https://github.com/pupillord)

## 项目介绍

TiCheck 是为 TiDB 设计的自动化，可扩展的检查工具，核心特点有：

- 自动，自动生成检查报告，定时执行检查，自动提交检查报告
- 可视，检查报告可视化展示，检查报告可视化查看
- 可定制，在网页界面上添加检查项目，更改检查项目的配置
- 可扩展，可以自定义检查规则，支持多种语言

## 背景

每个 DBA 都有他/她做日常巡检的一套“脚本武器库”，作为 TiDB DBA 的我们当然也不例外。但在实际的使用时，我们发现手动的使用这些脚本总会有这样那样的缺点：

1. 监控数据的多种来源造成的数据获取问题：
   这些脚本有的是针对机器状态的 shell 脚本，有的是针对数据健康的 sql 语句，有的是对集群 metric 进行收集的 prometheusQL。完成一次巡检通常意味着要到多个数据源执行各个脚本，而我们希望只需要一个统一的平台，可以自动化执行所有的脚本。而且对于这些重复性高的工作，我们希望能够定时自动执行。在隔离要求较高的场景下，我们还希望它能自动上传至指定 ftp 服务器，避免每次操作机器的授权问题。

2. 缺少统一的可视化界面
   在这些脚本中，我们可以看到每个脚本的执行结果，但基本时一个脚本一个结果，而且展示的界面还不同。Shell 脚本在终端上执行，得到的结果可能只是一个数字 sql 脚本在 mysql 上执行，得到的结果是一个表，prometheusQL 脚本在 prometheus 上执行，得到的结果是一个图。我们希望能够把这些结果统一的展示出来，更直观的感受集群的状态，以便更好的管理。

3. 定制与扩展的局限
   因为 Prometheus 与 Grafana 的存在，我们同样想过将所有巡检的结果都集中到一个 Grafana 面板里，但这带来了极大的局限性。首先就是数据来源的局限性。Grafana 所有的数据来源都是 Prometheus，而总有一些巡检项目是需要实际跑 SQL/Shell 的，比如说业务逻辑相关的检查，或是针对非 TiDB 但相关服务的状态检查。如果要将这些也集成至 Prometheus，则需要修改甚至重写一些 prometheus exporter，带来极大的业务量与极低的便携性与扩展性。而且我们要考虑到甲方对于这些脚本的审核，所以我们希望这些脚本以及配置项能够以明文单独展示。

## 设计

架构设计
![QQ截图20211229174657.png](https://s2.loli.net/2021/12/29/uz6mrfDYxyAp921.png)

功能原型
![image.png](https://s2.loli.net/2021/12/30/BD6C8Yvy49RqWaG.png)

## 进度

本期 Hackathon 计划完成的功能点：

- [x] 安装逻辑
- [x] 脚本运行逻辑
- [x] 自定义巡检配置
- [ ] 可视化 UI 后端
- [x] 可视化 UI 前端
- [ ] 巡检汇总统计
- [x] 数据库用户联名登录

后期规划功能点：

- [ ] 通过 UI 做巡检配置
- [ ] 通过 UI 做脚本扩展


## Quick Start

### Use Docker

```bash
# clone repo
git clone -b dev https://github.com/DigitalChinaOpenSource/TiCheck.git
cd TiCheck

# use docker to build image
docker build -t ticheck:latest .

# use docker to run this image 
cker run --name ticheck -p 8081:8081 -d ticheck:latest

# and you can access ticheck through port 801
```

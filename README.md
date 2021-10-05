# ProxyPanel

通过类似salt stack的方式提供管理多个代理服务器节点的面板，不做流控、审计、收费

## 主要功能

- 通过nginx反向代理中增加lua脚本，基于本地redis检查用户订阅状态来支持多用户
- 自动生成V2rayN、shadowrocket兼容的订阅

## 架构与实现细节

- Panel Master
直连延迟最低的某节点，便于用户和admin操作

- Panel Worker
其他节点，通过心跳报活，同时拉取新配置、同步用户订阅状态

### 用户数据

- 数据在Panel master上，其他Panel Instance拉配置并apply，【可能的单点故障】

### nginx鉴权

- v2ray单端口部署ws+tls，通过path里的query参数进redis里查订阅过期时间

- 不做流控、审计

### 管理面板整个上云Pros & Cons

#### Cons

- 额外成本?

- 至少需要两套key：一套只读数据库给worker用，一套读写给master用，配置和运维更复杂?

#### Pros

- 额外成本可能非常小

- master不会有CPU压力，由于无状态，也不再需要master，不再存在单点故障

- 配合部署脚本使用的话会更方便，节点部署后可以自动向面板注册

#### Questions

- 起一个定时脚本拉订阅清单直接查读文件/查redis/查云上数据库不确定哪个开销更小 <- redis会多一个进程（服务）运行在单核机器上（都跑了v2ray和nginx了，还差个redis？），文件是阻塞操作会严重影响性能，直接访问云上数据库不知道是否会对体验有大影响，反正只有第一次建立连接的时候需要

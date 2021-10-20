# Lazarus

提供管理多个服务器节点的面板与用户端

## 主要功能

### 一阶段

- 通过nginx反向代理中增加lua脚本，基于本地redis检查用户订阅状态来支持多用户
- 自动生成V2rayN、shadowrocket兼容的订阅
- 用户面板中显示全部metadata：服务端在线数、总用量、并发数……

### 二阶段

- 配置其他机器key后自动部署v2ray和worker实例
- terraform自动买机器扩容

## 架构与实现细节

- Panel Master
跑在firebase/其他云/最低延迟服上，中控节点 (单点故障)

- Panel Worker
所有普通节点，通过心跳向master报活，同时拉取新配置（同步用户订阅状态）  

如果拉不到状态则默认deny all 

### nginx鉴权

- v2ray单端口部署ws+tls，使用access_by_lua做实时鉴权，用path里的query参数向本地Panel worker查询过期时间

- 不做流控、审计

## Proxy Panel Instances设计

### Panel Worker

- golang写，二进制可直接部署
- 所有信息只留在内存中，重启后需要向panel master重新拉取
- 向内网lua提供http(tcp?)接口用于鉴权（只使用状态码）
- 向外网panel master发心跳，多次超时则认为master down，deny all access，心跳还包含当前资源使用情况: CPU, MEM，带宽，并发数
- 需要部署时无感支持ipv4和v6
- 配置文件：panel master ip和端口，心跳频率，masterDownTimeoutCountThres

### Panel Master

- 接受proxy worker上报数据，持久化，可生成监控图，同时响应中增量更新配置信息（至少实现配置是否有更新，如果有则附上配置）
- 业务配置信息存储：redis，用key的expire实现过期；
- 注册信息同样redis即可。需要在redis上做事务

- 如果用云原生，如firebase：需要重新设计订阅过期：加个expireAt字段，服务端下发时进行计算（数据量大了之后可能会有性能问题）
- 本地存的话sqlite就够了，低写较低读。

### RPC

- gRPC+SSL

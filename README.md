# Lazarus

"""Call forth Lazurus from his cave, and resurrect the dead."""

Proxy management server, with user interaction panel.

提供管理多个服务器节点的面板与用户端

## Main Functions 主要功能

### Stage 1 一阶段

- Lua Scripts with Open-Resty to support user authentication on workers. 通过nginx反向代理中增加lua脚本，基于本地worker检查用户订阅状态来支持多用户
- V2rayN & shadowrocket compatible subscription links generator. 自动生成V2rayN、shadowrocket兼容的订阅
- A user panel to display statistics: concurrency on port 443, total online users, total data used... 用户面板中显示全部metadata：服务端在线数、总用量、443并发数……

### Stage 2 二阶段

- Auto deployment of v2ray, nginx and workers with ssh keys on designated servers. 配置其他机器key后自动部署v2ray,nginx和worker实例
- Auto scaling with terraform. terraform自动买机器扩容
- Dynamic naming service with Cloudflare DNS API. 所有worker的域名为四级域名([WORKERID].xxx.xxx.xxx)，由master调cloudflare api实现动态解析，worker只需master的RPC地址即可工作，避免复杂配置  
- dockerize for easy deployments.

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
- 需要部署时支持ipv4和v6
- 配置文件：panel master ip和端口，心跳频率，masterDownTimeoutCountThres

#### worker 时间轴

- master收到建机器请求后去cf加dns，同时通过ssh登录worker机器安装docker-compose并启动worker,nginx和v2ray镜像，数据盘挂到宿主机/opt上
- 初始化：worker启动后如果发现已经完成配置则goto注册；首先向master请求装需要的配置信息：add/host（initializeRequest）
- master回复initializeResponse，下发配置
- worker应用配置，向lets encrypt申请证书.任何一步失败则退出，由master检测超时，成功则在cwd下写一个文件表明环境配置正确，以便重启时goto注册
- 注册：worker配置完成后向master进行服务注册，若注册失败则退出
- 心跳：worker启动后向master定期发带自身负载数据的心跳包
- master发生故障（worker发现多次心跳rpc超时）后worker直接退出进程
- worker进程退出后由进程管理器（一阶段systemd，二阶段docker）重启

### Panel Master

- 接受proxy worker上报数据，刷盘？，可生成监控图，同时响应中增量更新配置信息（至少实现配置是否有更新，如果有则附上配置）

- 设计订阅过期：expireAt字段，服务端每次下发时进行计算（数据量大了之后可能会有性能问题）
- sqlite持久化用户数据(低写较低读)。

### RPC

- thrift+SSL
- `thrift -r  --gen go idl/service.thrift`

## Deploy

```bash
cd lazarus
cp lazarus.service /usr/lib/systemd/system/lazarus.service
systemctl daemon-reload
systemctl start lazarus.service
```

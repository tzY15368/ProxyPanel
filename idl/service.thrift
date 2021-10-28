// thrift -r  --gen go idl/service.thrift
namespace go RPCService

typedef set<string> UserData

struct InitializeRequest {
    1:required string mac;
}

struct InitializeResponse{
    1:required string add;
    2:required string host;
}

struct RegisterRequest {
    // 默认所有实例都跑在/index.php 443端口
    1:required string mac;
}

struct RegisterResponse {
    1:required i32 heartBeatRateIntervalSec;
    2:required i32 HeartBeatErrorThres;
}

struct HeartbeatRequest {
    1:required string mac;

    // cpu占用百分比
    4:optional i32 cpu;

    // mem占用百分比
    5:optional i32 mem;

    // 并发数
    6:optional i32 activeConn;

    // 总流量(mb)
    7:optional i32 currentData;

    // 可用流量(mb)
    8:optional i32 totalData;
}

struct HeartbeatResponse {
    1:required bool hasUpdate;

    // 如果hasupdate==true，应用data中的filter
    2:optional UserData data;
}

service LazarusService {
    InitializeResponse DoInitialize(1:InitializeRequest ir)
    RegisterResponse DoRegister(1:RegisterRequest rr)
    HeartbeatResponse DoHeartBeat(1:HeartbeatRequest hbr)
}

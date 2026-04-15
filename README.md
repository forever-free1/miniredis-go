# miniredis-go

一个 Go 语言实现的简化版 Redis 服务器。

## 功能特性

### 已实现命令

#### String 命令
- `PING` - 心跳检测
- `GET` - 获取值
- `SET key value [EX seconds]` - 设置值（可选过期时间）
- `DEL key [key ...]` - 删除键
- `EXISTS key [key ...]` - 检查键是否存在
- `INCR key` - 递增
- `DECR key` - 递减
- `APPEND key value` - 追加值
- `STRLEN key` - 获取字符串长度

#### List 命令
- `LPUSH key value [value ...]` - 从左侧插入
- `RPUSH key value [value ...]` - 从右侧插入
- `LRANGE key start stop` - 获取列表范围
- `LLEN key` - 获取列表长度
- `LINDEX key index` - 按索引获取元素

#### Hash 命令
- `HSET key field value` - 设置 Hash 字段
- `HGET key field` - 获取 Hash 字段
- `HGETALL key` - 获取所有字段
- `HDEL key field [field ...]` - 删除字段
- `HEXISTS key field` - 检查字段是否存在
- `HLEN key` - 获取字段数量

#### Set 命令
- `SADD key member [member ...]` - 添加成员
- `SMEMBERS key` - 获取所有成员
- `SISMEMBER key member` - 检查成员是否存在
- `SCARD key` - 获取成员数量
- `SREM key member [member ...]` - 删除成员

#### 过期时间命令
- `EXPIRE key seconds` - 设置过期时间
- `TTL key` - 获取剩余生存时间

#### Pub/Sub 命令
- `PUBLISH channel message` - 发布消息
- `SUBSCRIBE channel [channel ...]` - 订阅频道
- `UNSUBSCRIBE [channel ...]` - 取消订阅
- `PSUBSCRIBE pattern [pattern ...]` - 按模式订阅
- `PUNSUBSCRIBE [pattern ...]` - 取消模式订阅

#### 事务命令
- `MULTI` - 开始事务
- `EXEC` - 执行事务
- `DISCARD` - 丢弃事务

## 快速开始

### 构建
```bash
go build -o miniredis.exe .
```

### 运行
```bash
./miniredis.exe
# 默认监听 6379 端口

# 指定端口
./miniredis.exe :6380
```

### 测试
```bash
go run ./cmd/test.go
```

## 项目结构

```
miniredis-go/
├── main.go              # 程序入口
├── README.md            # 本文件
├── task_plan.md         # 开发计划
├── progress.md          # 进度记录
├── findings.md          # 研究发现
├── server/
│   ├── server.go        # TCP 服务器
│   ├── handler.go       # 命令处理器
│   ├── database.go      # 数据存储
│   └── resp.go          # RESP 协议解析
└── cmd/
    └── test.go          # 测试客户端
```

## 技术实现

- **网络**: Go 原生 `net` 包，goroutine per connection
- **协议**: RESP (Redis Serialization Protocol)
- **存储**: 内存存储，无持久化
- **并发**: sync.RWMutex 保护共享数据

## 限制

- 不支持数据持久化
- 不支持集群
- 不支持 SSL/TLS
- Pub/Sub 为简化实现，消息不实际推送给订阅者

## License

MIT

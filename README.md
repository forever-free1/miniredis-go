# miniredis-go

一个 Go 语言实现的轻量级 Redis 服务器，实现了 Redis 最核心最常用的功能。

## 目录

- [功能特性](#功能特性)
- [快速开始](#快速开始)
- [命令参考](#命令参考)
- [项目结构](#项目结构)
- [技术实现](#技术实现)
- [测试](#测试)
- [限制与已知问题](#限制与已知问题)
- [贡献](#贡献)
- [许可](#许可)

---

## 功能特性

### 已实现的数据结构

| 数据类型 | 支持命令 | 说明 |
|---------|---------|------|
| String | PING, GET, SET, DEL, EXISTS, INCR, DECR, APPEND, STRLEN | 字符串类型，包含原子递增/递减 |
| List | LPUSH, RPUSH, LRANGE, LLEN, LINDEX | 双端列表，支持从任意一端操作 |
| Hash | HSET, HGET, HGETALL, HDEL, HEXISTS, HLEN | 哈希表，适合存储对象 |
| Set | SADD, SMEMBERS, SISMEMBER, SCARD, SREM | 无序集合，支持成员关系判断 |

### 已实现的高级功能

| 功能 | 命令 | 说明 |
|------|------|------|
| 过期时间 | EXPIRE, TTL, SET EX | 支持键的自动过期 |
| 事务 | MULTI, EXEC, DISCARD | 支持命令批量执行 |
| 发布/订阅 | PUBLISH, SUBSCRIBE, UNSUBSCRIBE, PSUBSCRIBE, PUNSUBSCRIBE | 支持频道和模式匹配 |

---

## 快速开始

### 环境要求

- Go 1.16 或更高版本

### 构建

```bash
# 克隆仓库
git clone https://github.com/forever-free1/miniredis-go.git
cd miniredis-go

# 构建
go build -o miniredis.exe .
```

### 运行

```bash
# 默认监听 6379 端口
./miniredis.exe

# 自定义端口
./miniredis.exe :6380

# 自定义地址
./miniredis.exe localhost:6379
```

### 测试

```bash
# 运行综合测试套件
go run ./cmd/test.go
```

---

## 命令参考

### String 命令

| 命令 | 语法 | 说明 | 示例 |
|------|------|------|------|
| PING | `PING [message]` | 心跳检测 | `PING` → +PONG |
| GET | `GET key` | 获取值 | `GET name` → $5\r\nhello |
| SET | `SET key value [EX seconds]` | 设置值 | `SET key value EX 60` |
| DEL | `DEL key [key ...]` | 删除键 | `DEL key1 key2` |
| EXISTS | `EXISTS key [key ...]` | 检查键是否存在 | `EXISTS key` → :1 |
| INCR | `INCR key` | 原子递增 | `INCR counter` → :1 |
| DECR | `DECR key` | 原子递减 | `DECR counter` → :0 |
| APPEND | `APPEND key value` | 追加字符串 | `APPEND key abc` → :8 |
| STRLEN | `STRLEN key` | 获取长度 | `STRLEN key` → :5 |

### List 命令

| 命令 | 语法 | 说明 | 示例 |
|------|------|------|------|
| LPUSH | `LPUSH key value [value ...]` | 从左侧插入 | `LPUSH list a b c` |
| RPUSH | `RPUSH key value [value ...]` | 从右侧插入 | `RPUSH list d e` |
| LRANGE | `LRANGE key start stop` | 获取范围元素 | `LRANGE list 0 -1` |
| LLEN | `LLEN key` | 获取长度 | `LLEN list` → :5 |
| LINDEX | `LINDEX key index` | 按索引获取 | `LINDEX list 0` |

### Hash 命令

| 命令 | 语法 | 说明 | 示例 |
|------|------|------|------|
| HSET | `HSET key field value` | 设置字段 | `HSET user name Alice` → :1 |
| HGET | `HGET key field` | 获取字段值 | `HGET user name` → $5\r\nAlice |
| HGETALL | `HGETALL key` | 获取所有字段值 | `HGETALL user` |
| HDEL | `HDEL key field [field ...]` | 删除字段 | `HDEL user age` |
| HEXISTS | `HEXISTS key field` | 检查字段存在 | `HEXISTS user name` → :1 |
| HLEN | `HLEN key` | 获取字段数量 | `HLEN user` → :2 |

### Set 命令

| 命令 | 语法 | 说明 | 示例 |
|------|------|------|------|
| SADD | `SADD key member [member ...]` | 添加成员 | `SADD tags go redis` → :3 |
| SMEMBERS | `SMEMBERS key` | 获取所有成员 | `SMEMBERS tags` |
| SISMEMBER | `SISMEMBER key member` | 检查成员存在 | `SISMEMBER tags go` → :1 |
| SCARD | `SCARD key` | 获取成员数量 | `SCARD tags` → :3 |
| SREM | `SREM key member [member ...]` | 删除成员 | `SREM tags redis` |

### 过期时间命令

| 命令 | 语法 | 说明 | 示例 |
|------|------|------|------|
| EXPIRE | `EXPIRE key seconds` | 设置过期秒数 | `EXPIRE key 60` → :1 |
| TTL | `TTL key` | 获取剩余生存时间 | `TTL key` → :59 |

**TTL 返回值说明：**
- `-2`: 键不存在
- `-1`: 键没有过期设置
- `>=0`: 剩余秒数

### Pub/Sub 命令

| 命令 | 语法 | 说明 | 示例 |
|------|------|------|------|
| PUBLISH | `PUBLISH channel message` | 发布消息 | `PUBLISH news "hello"` → :0 |
| SUBSCRIBE | `SUBSCRIBE channel [channel ...]` | 订阅频道 | `SUBSCRIBE news` |
| UNSUBSCRIBE | `UNSUBSCRIBE [channel ...]` | 取消订阅 | `UNSUBSCRIBE news` |
| PSUBSCRIBE | `PSUBSCRIBE pattern [pattern ...]` | 按模式订阅 | `PSUBSCRIBE news.*` |
| PUNSUBSCRIBE | `PUNSUBSCRIBE [pattern ...]` | 取消模式订阅 | `PUNSUBSCRIBE news.*` |

**模式匹配规则：**
- `*` 匹配任意字符序列
- `?` 匹配任意单个字符

### 事务命令

| 命令 | 语法 | 说明 | 示例 |
|------|------|------|------|
| MULTI | `MULTI` | 开始事务 | `MULTI` → +OK |
| EXEC | `EXEC` | 执行事务 | `EXEC` → *3\r\n... |
| DISCARD | `DISCARD` | 丢弃事务 | `DISCARD` → +OK |

**事务使用示例：**
```
MULTI
SET key1 value1
SET key2 value2
INCR counter
EXEC
```

---

## 项目结构

```
miniredis-go/
├── main.go                 # 程序入口
├── README.md               # 项目文档
├── task_plan.md            # 开发计划
├── progress.md             # 开发进度记录
├── findings.md             # 研究发现
│
├── server/
│   ├── server.go          # TCP 服务器实现
│   │                       #   - 监听端口接受连接
│   │                       #   - 为每个连接创建 Handler
│   │                       #   - 读取请求并返回响应
│   │                       #
│   ├── handler.go         # 命令处理器
│   │                       #   - 所有 Redis 命令的实现
│   │                       #   - Handler 结构维护事务状态
│   │                       #
│   ├── database.go        # 数据存储层
│   │                       #   - 内存存储 (store, listStore, hashStore, setStore)
│   │                       #   - Pub/Sub 状态管理
│   │                       #   - 线程安全 (sync.RWMutex)
│   │                       #
│   └── resp.go            # RESP 协议解析
│                             #   - ParseCommand: 解析命令
│                             #   - Encode*: 编码响应
│
└── cmd/
    └── test.go            # 测试客户端
                              #   - 综合测试所有命令
                              #   - 支持持久连接的测试
```

---

## 技术实现

### 网络层

- 使用 Go 标准库 `net` 包
- 采用 **goroutine per connection** 模式
- 每个连接独立处理，并发安全

### 协议

实现 Redis Serialization Protocol (RESP)：

- **请求格式**：Array (以 `*` 开头)
- **响应格式**：
  - Simple String: `+OK\r\n`
  - Error: `-ERR message\r\n`
  - Integer: `:123\r\n`
  - Bulk String: `$5\r\nhello\r\n`
  - Array: `*3\r\n...`

### 数据存储

- **内存存储**：所有数据存储在内存中
- **线程安全**：使用 `sync.RWMutex` 保护共享数据
- **过期处理**：过期时间存储在 `ExpireAt` 字段，GET/EXISTS 时检查

### 并发模型

```
Client 1 ──┐
Client 2 ──┼── TCP Listener ── goroutine per connection
Client N ──┘
              │
              ▼
         Handler (per connection)
              │
              ▼
         ExecuteCommand()
              │
              ▼
         Data Store (protected by RWMutex)
```

---

## 测试

### 运行测试

```bash
# 启动服务器
./miniredis.exe

# 在另一个终端运行测试
go run ./cmd/test.go
```

### 测试覆盖

测试套件验证以下功能：

- String: PING, SET, GET, INCR, DECR, APPEND, STRLEN, EXISTS, DEL
- List: LPUSH, RPUSH, LRANGE, LLEN, LINDEX
- Hash: HSET, HGET, HGETALL, HDEL, HEXISTS, HLEN
- Set: SADD, SMEMBERS, SISMEMBER, SCARD, SREM
- Expire: EXPIRE, TTL, SET EX
- Pub/Sub: PUBLISH, SUBSCRIBE, UNSUBSCRIBE, PSUBSCRIBE, PUNSUBSCRIBE
- Transaction: MULTI, EXEC, DISCARD

### 使用 redis-cli 测试

```bash
# 连接服务器
redis-cli

# 测试命令
> PING
PONG

> SET mykey "hello"
OK

> GET mykey
"hello"

> INCR counter
(integer) 1

> MULTI
OK

> SET key1 value1
QUEUED

> INCR counter
QUEUED

> EXEC
1) OK
2) (integer) 2
```

---

## 限制与已知问题

### 已知限制

1. **无持久化**：数据仅存储在内存，服务器关闭后丢失
2. **无集群支持**：不支持 Redis Cluster
3. **无 SSL/TLS**：不支持加密连接
4. **无认证**：不支持密码保护
5. **Pub/Sub 简化实现**：订阅者消息推送未完全实现

### 与 Redis 的差异

1. 不支持事务内的 WATCH 命令
2. 不支持持久化的 AOF 和 RDB
3. 不支持发布消息的实际推送（仅计数）
4. 不支持 Lua 脚本

---

## 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

---

## 许可

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

---

## 致谢

- [Redis 官方文档](https://redis.io/documentation) - 命令参考和协议说明
- [Go 语言](https://golang.org/) - 项目实现语言

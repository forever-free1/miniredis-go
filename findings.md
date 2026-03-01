# Findings & Decisions

## Requirements
<!-- 用户需求：实现 Go 语言简化版 Redis -->
- 实现 TCP 服务器，监听 6379 端口
- 实现 RESP (Redis Serialization Protocol) 协议解析
- 实现 PING 命令 - 心跳检测
- 实现核心 String 命令：GET, SET, DEL, EXISTS, INCR, DECR, APPEND, STRLEN
- 实现 List 命令：LPUSH, RPUSH, LRANGE, LLEN, LINDEX
- 实现 Hash 命令：HSET, HGET, HGETALL, HDEL, HEXISTS
- 实现 Set 命令：SADD, SMEMBERS, SISMEMBER, SCARD, SREM
- 实现 EXPIRE key seconds - 设置过期时间
- 实现 TTL key - 查看剩余时间
- 保持与 Redis 协议兼容
- 每次完成功能后推送到 GitHub，必须经过测试

## Research Findings

### Redis 核心命令优先级
根据使用频率，优先级排序：
1. **String** - GET, SET (最常用)
2. **Key** - DEL, EXISTS, EXPIRE, TTL
3. **Hash** - HSET, HGET, HGETALL
4. **List** - LPUSH, RPUSH, LRANGE
5. **Set** - SADD, SMEMBERS, SISMEMBER

### RESP 协议要点
- 简单字符串: +OK\r\n
- 错误: -ERR message\r\n
- 整数: :1000\r\n
- 批量字符串: $6\r\nfoobar\r\n
- 数组: *2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n

### 项目当前状态
- 空白模板项目
- 仅有 go.mod 和占位符 main.go
- 无任何 Redis 功能实现

## Technical Decisions
| Decision | Rationale |
|----------|-----------|
| RESP2 协议 | Redis 标准协议，兼容 redis-cli |
| 内存 Map 存储 | 简化实现，快速原型 |
| Goroutine per connection | Go 并发模型天然支持 |

## Issues Encountered
| Issue | Resolution |
|-------|------------|
|       |            |

## Resources
- Redis Protocol: https://redis.io/topics/protocol
- Redis Commands: https://redis.io/commands

## Visual/Browser Findings
- 项目目录结构简单，无现有实现

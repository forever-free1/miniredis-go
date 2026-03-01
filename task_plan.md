# Task Plan: miniredis-go 实现计划

## Goal
实现一个 Go 语言的简化版 Redis 服务器，实现 Redis 最核心最常用的功能。

## 版本控制
- 每次完成一部分计划时，代码必须推送到 https://github.com/forever-free1/miniredis-go.git
- 代码必须经过测试才能推送
- 使用 Git 进行版本控制

## Current Phase
Phase 2

## Phases

### Phase 1: 需求分析与架构设计
- [x] 探索项目当前状态
- [x] 确定要实现的 Redis 核心命令
- [x] 设计系统架构
- [x] 制定详细实现计划
- **Status:** complete

### Phase 2: 基础框架搭建
- [x] 实现 TCP 服务器 (监听 6379 端口)
- [x] 实现 RESP 协议解析器
- [x] 实现 PING 命令 - 心跳检测
- [x] 实现基础命令路由 (GET, SET, DEL, EXISTS)
- **Status:** complete

### Phase 3: 核心 String 命令实现
- [x] 实现 GET 命令
- [x] 实现 SET 命令
- [x] 实现 DEL 命令
- [x] 实现 EXISTS 命令
- [ ] 实现其他 String 相关命令 (INCR, DECR, APPEND, STRLEN)
- **Status:** in_progress

### Phase 4: 常用数据结构实现
- [ ] 实现 List 命令 (LPUSH, RPUSH, LRANGE, etc.)
- [ ] 实现 Hash 命令 (HSET, HGET, HGETALL, etc.)
- [ ] 实现 Set 命令 (SADD, SMEMBERS, SISMEMBER, etc.)
- **Status:** pending

### Phase 5: 高级功能实现
- [ ] 实现 EXPIRE key seconds - 设置过期时间
- [ ] 实现 TTL key - 查看剩余时间
- [ ] 实现 Pub/Sub 功能
- [ ] 实现事务支持 (MULTI, EXEC, DISCARD)
- **Status:** pending

### Phase 6: 测试与验证
- [ ] 编写单元测试
- [ ] 使用 redis-cli 连接测试
- [ ] 性能测试
- **Status:** pending

## Key Questions
1. 需要支持哪些 Redis 数据类型？ (String, List, Hash, Set)
2. 是否需要支持持久化？ (当前版本不考虑)
3. 并发处理策略？ (goroutine per connection)

## Decisions Made
| Decision | Rationale |
|----------|-----------|
| 使用 RESP2 协议 | Redis 标准协议，兼容 redis-cli |
| 内存存储 | 简化实现，不考虑持久化 |
| Goroutine per connection | Go 语言天然支持，简化并发处理 |

## Errors Encountered
| Error | Attempt | Resolution |
|-------|---------|------------|
|       |         |            |

## Notes
- 优先级：先实现最常用的 String 命令
- 参考 Redis 官方文档实现命令
- 保持与 Redis 兼容的响应格式

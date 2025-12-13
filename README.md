# gochat

A simple chat application written in Go.

## 日志服务

日志服务使用 [zap](go.uber.org/zap) 实现
使用教程参考：[【golang】最受欢迎的开源日志库——zap零基础入门使用](https://www.bilibili.com/video/BV1Rk99YHEM6)

## Swagger

Swag 使用参考：<https://github.com/swaggo/swag/blob/master/README_zh-CN.md>

1）安装 swag

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

2）swag 格式化

```shell
swag fmt
```

3）swag 生成文档

```shell
swag init -g ./internal/router/router.go
```

4）开发环境访问：<http://localhost:8083/swagger/index.html>

## 错误码

### 错误码规范

错误码（code）结构：DDCCC (5 位数字)

- DD (前两位)：领域/模块标识符（Domain）
- CCC (后三位)：该领域内的具体错误编号（Code）

1）DD（领域/模块标识符）
这部分标识了错误是发生在哪个核心业务模块或系统组件。

| DD 编号 | 领域/模块   | 描述                                                                        |
| ------- | ----------- | --------------------------------------------------------------------------- |
| 10      | 系统级/通用 | 数据库、缓存、内部 RPC 错误、通用服务器错误 (通常对应 5xx HTTP 状态码)。    |
| 20      | 用户/认证   | 用户注册、登录、Token 验证、权限问题 (通常对应 401/403 HTTP 状态码)。       |
| 30      | 参数/校验   | 所有客户端输入校验失败、格式错误、必填字段缺失 (通常对应 400 HTTP 状态码)。 |
| 40      | 资源/数据   | 资源找不到、资源已被删除、资源唯一性冲突 (通常对应 404/409 HTTP 状态码)。   |
| 50      | 订单/交易   | 订单创建、支付、退款等业务流程错误。                                        |
| 60      | 商品/库存   | 商品上下架、库存不足等错误。                                                |

2）CCC（具体错误编号）
这部分用于区分领域内具体的错误类型，从 001 开始递增。

### 分层

| 层级                | 职责             | 错误载体       | 核心作用                                            |
| ------------------- | ---------------- | -------------- | --------------------------------------------------- |
| Repository (DB/DAO) | 基础设施错误处理 | Go error       | 仅传递数据库或外部服务错误                          |
| Service (业务层)    | 核心错误映射     | \*ServiceError | 封装 业务码 (Code) 和 HTTP 状态码 (HTTPStatus)      |
| Handler (控制层)    | 最终 HTTP 响应   | HTTP 响应      | 读取 \*ServiceError，设置 HTTP 状态码和 JSON 错误体 |

> DAO（Data Access Object）：数据访问对象层，DAO 层的主要目标是实现业务逻辑层与持久化机制的解耦。

## TODO

- [x] 完善日志
- [x] 验证 token 时效
- [x] 零值陷阱
- [ ] 参数区分大小写

# gochat

A simple chat application written in Go.

## 日志服务

日志服务使用 [zap](go.uber.org/zap) 实现
使用教程参考：[【golang】最受欢迎的开源日志库——zap零基础入门使用](https://www.bilibili.com/video/BV1Rk99YHEM6)

## Swagger

Swag 使用参考：<https://github.com/swaggo/swag/blob/master/README_zh-CN.md>

1. 安装 swag

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

2. swag 格式化

```shell
swag fmt
```

3. swag 生成文档

```shell
swag init -g ./internal/router/router.go
```

1. 开发环境访问：<http://localhost:8083/swagger/index.html>

## TODO

- [x] 完善日志
- [x] 验证 token 时效

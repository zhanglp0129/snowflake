# snowflake
雪花算法，用于生成分布式数据库中的主键。能保证全局唯一性、趋势递增和信息安全。

## 原理

## 使用
1. 安装依赖
```shell
go get -u github.com/zhanglp0129/snowflake
```

2. 创建雪花算法配置
```go
cfg := snowflake.DefaultConfig
cfg.SetStartTime(time.Parse("2006-01-02 15:04:05.000", "2024-8-14 00:00:00.000"))
```

3. 创建工作节点
```go
worker, err := snowflake.NewWorker(cfg, 0)
```

4. 生成id
```go
id, err := worker.GenerateId()
```

## LICENSE

MIT
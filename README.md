# Go-Web-Template

## 简介

Go Web 项目模版。集成了常用的组件，方便开发。

目前实现了如下组件：

|类型|组件名字|组件描述|
|::|::|::|
|配置| config   |解析配置文件|
|日志|logger|日志打印|
|数据库|mysql|MySQL 数据库连接|
|缓存|redis|Redis 连接|
|消息队列|kafka|Kafka 生产者、消费者|
|中心化服务|nacos|Nacos是一个分布式配置中心、服务注册和发现中心|

### 为什么这么做呢？

有两个组件是所有其它项目都可能依赖的组件，一个是 config，另一个是 logger。

配置和日志是所有项目必须的内容。

使用配置文件定义组件的初始化参数，可以让组件的初始化变得简单，把配置信息添加加到配置中即可。

使用日志插件，可以在项目中方便的打印运行信息，方便问题排查和Debug。

### 怎么使用？
组件的使用方法都类似，首先需要在配置文件中加入对应组件的配置，然后把配置传给组件的初始化方法，即可初始化组件。

**配置文件：**

```go
config.Init(confFile, global.Configs)
```

**日志插件**

```go
logger.Init(global.Configs.App.Name, global.Configs.LoggerConfig)
```

**MySQL**

```go
# 初始化
mysql.Init(background, global.Configs.Mysql)
# 使用
db := mysql.GetDB(background, "test")
```


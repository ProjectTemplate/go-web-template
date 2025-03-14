## Go-Web-Sandbox 

### 说明

该模用 Docker 块运行 go-web-template 项目运行依赖的组件。方便用户启动简单的测试。

目前支持的组件：

|组件名|是否已支持|描述|访问地址|
|:-:|:--:|:-:|---|
|Kafka|否❎|消息队列||
|Redis|否❎|Key/Value 缓存||
|MySQL|否❎|数据库||
|Jaeger|是✅|Trace 可视化后端|http://127.0.0.1:16686/jaeger/ui/search|
|otel-collector|是✅|OpenTelemetry 收集器|HTTP上报：http://127.0.0.1:4318 <br>GRPC上报：http://127.0.0.1:4317|
|Prometheus|是✅|数据采集|http://localhost:9090/|
|Grafana|是✅|可视化面板|http://127.0.0.1:3000/grafana/dashboards|
|opensearch|是✅|数据存储|http://localhost:9200|

**如何使用：** 在 `go-web-sandbox` 目录执行如下命令即可启动。


```shell
docker compose up --force-recreate --remove-orphans --detach
```



> 备注：改项目改进自 [https://github.com/open-telemetry/opentelemetry-demo](https://github.com/open-telemetry/opentelemetry-demo)
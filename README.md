# 环境搭建

### 1. golang下载地址

- https://go.dev/dl

```shell
wget https://go.dev/dl/go1.20.11.linux-amd64.tar.gz
tar zxvf go1.20.11.linux-amd64.tar.gz
mv go go1.20.11
export GOROOT=~/go1.20.11
mkdir go
export GOPATH=~/go
go build -tags=jsoniter .
```

### 2. 初始化项目

- 设置代理

```shell
go env -w GOPROXY=https://goproxy.cn,direct
#go env -w GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```

- go.mod不存在，初始化项目

```shell
go mod init wsd-athena-go
go get -u github.com/gin-gonic/gin
```

- go.mod已存在，加入项目

```shell
go mod tidy
```

# 教程

### 推荐个实战的教程

- https://juejin.cn/post/7016742808560074783

### gin

```shell
go get -u github.com/gin-gonic/gin
```

- https://gin-gonic.com/zh-cn/docs/

### log

```shell
go get -u go.uber.org/zap
go get -u gopkg.in/natefinch/lumberjack.v2
```

- 教程: https://github.com/uber-go/zap

### mysql gorm

```shell
go get -u gorm.io/gorm
# GORM 官方支持 sqlite、mysql、postgres、sqlserver
go get -u gorm.io/driver/mysql
```

- gorm教程: https://topgoer.com/%E6%95%B0%E6%8D%AE%E5%BA%93%E6%93%8D%E4%BD%9C/gorm/gorm%E6%9F%A5%E8%AF%A2.html

### redis

```shell
go get -u github.com/go-redis/redis/v8
```

- 官网教程: https://redis.uptrace.dev/zh

### rocketmq

```shell
go get -u github.com/apache/rocketmq-client-go/v2
export ROCKETMQ_GO_LOG_LEVEL=error
```

- 官网教程: https://rocketmq.apache.org/zh/docs/quickStart/01quickstart/

### swagger

```shell
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

- 官网教程: https://github.com/swaggo/gin-swagger

```shell
swag init
```

- 访问地址: http://127.0.0.1:48085/swagger/index.html

### xxl-job

```shell
go get -u github.com/xxl-job/xxl-job-executor-go
```

- 官网教程: https://github.com/gin-middleware/xxl-job-executor

### nacos

```shell
go get -u github.com/nacos-group/nacos-sdk-go/v2
```

- 官网教程: https://github.com/nacos-group/nacos-sdk-go/blob/master/README_CN.md

### grpc

```shell
go get -u google.golang.org/grpc
```

- 官网教程: https://grpc.io/docs/languages/go/quickstart/

```shell
protoc -I=. --go_out=./pb --go-grpc_out=./pb ./proto/*
```

### jwt

```shell
go get -u github.com/dgrijalva/jwt-go
```

### validator

```shell
go get -u github.com/go-playground/validator/v10
```

### cmux

```shell
go get -u github.com/soheilhy/cmux
```

- 官网教程: https://github.com/soheilhy/cmux

### pprof

```shell
go get github.com/gin-contrib/pprof
```

- 查看地址: http://localhost:48085/debug/pprof/

### 中间件

- https://github.com/orgs/gin-contrib/repositories?type=all

# 代码片段

# 好用的小工具

### SQL生成go结构体、SQL转Golang Struct、SQL转Struct、SQL转Go

- https://wetools.cc/sql2go

# 其他框架

- 字节go框架: https://www.cloudwego.io/zh/docs/hertz/getting-started/
- echo框架: https://echo.laily.net/
- fasthttp: https://github.com/valyala/fasthttp
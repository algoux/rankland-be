# 架构设计
架构层级分为三层，接口层、逻辑层、数据访问层。和一些纵向的工具层

## 1. 服务初始化
> Go 程序的初始化需要在 main package 中，为了简化 main 函数的逻辑，只进行最必要的 init 函数调用。所有需要初始化的组件放在 init 目录下构建

当前需要构建的组件有：获取 config 目录下的静态配置参数、初始化 gin 框架、使用 gorm 框架初始化 postgre 数据库、使用 go-redis 框架初始化 redis

## 2. 接口层
> 用来提供 API 接口和处理一些通用参数的校验。如验证入参是否必传、是否有效、是否在合适范围等，不处理业务逻辑的参数校验

映射到项目目录中，主要是：api 目录和 param 目录

## 3. 逻辑层
> 用于处理业务逻辑和封装底层数据层的业务原子化行为

映射到项目目录中，主要是：service 目录（业务逻辑层）和 logic 目录（原子化行为接口）

- service: 主要处理业务逻辑，尽可能避免直接操作缓存、数据等行为
- logic: 为 service 层服务，主要是对处理缓存和数据等相关逻辑，为 service 提供所必须得原子化行为接口。如：获取某项数据，支持缓存获取等


## 4. 数据访问层
> 用来和数据库进行最简单的 CRUD 操作和数据模型的定义，不和业务绑定

映射到项目目录中，主要是：model 目录

## 5. 其他
> 其他一些关键目录的主要含义

### error
> 用于封装业务的错误码，定义业务内部和对外的 http 交互错误码行为

### middleware 
> 定义业务的中间件，用于集中处理拦截、校验等逻辑

### util
> 业务中需要用到的一些通用转换工具

## script
> 项目运行的脚本

## 

## 6. 使用Docker进行部署
1. 在宿主机创建新的配置文件`config.yaml`（模板配置文件位于源码的`/config/config.template.yaml`），注意`application.host`需要设置为`0.0.0.0`才可以在docker进行端口转发。

2. 从腾讯云拉取镜像
    ```shell
    docker pull ccr.ccs.tencentyun.com/algoux/rankland-be
    ```
3. 运行镜像
    ```
    docker run -p 80:8000 -v /path/config.yaml:/app/rankland/config/config.yaml -v /path/temp/logs:/app/rankland/temp/logs -d ccr.ccs.tencentyun.com/algoux/rankland-be
    ```
# sim
基于grpc+etcd的微服务即时通讯项目

### 项目依赖
- gin
- grpc
- etcd
- gorm
- websocket
- jwt-go
- zap
- rabbitmq
- viper
- protobuf
- validator

### 项目结构
```text
sim /
|—— app                 // 项目应用
|  |—— core              // 框架内核
|  |—— gateway           // 网关服务
|     |—— http            // 包含控制器、中间件、表单验证登，应用请求处理逻辑都放这里
|     └── rpc             // 注册请求处理需要的rpc服务
|  |—— global            // 包含框架全局的常量、变量
|  |—— listeners         // 存放事件监听handle
|  |—— model             // 数据模型
|  |—— providers         // 服务注册
|  |—— services          // rpc服务
|  └── util              // 应用内自定义的一些包
|—— bootstrap           // 启动项目运行所需要的所有包
|—— cmd                 // 服务启动
|   |—— gateway         // 网关服务启动
|   |—— im              // im rpc服务启动
|   └── user            // user rpc服务启动
|—— config              // 配置文件
|—— database            // 数据库文件
|—— idl                 // protobuf文件
|—— public              // 存放共用静态文件
|—— routes              // 路由文件目录
└── storage             // 存储运行日志、上传文件等
```

### 项目说明
目前项目完成框架环境的搭建，实现了基本的聊天逻辑，如登录注册、客户端管理、消息收发等功能，后续会
完善聊天内容的存储、离线消息处理机制、以及消息的多端同步等功能

### 配置文件
将`config/config.yaml.example`拷贝至`config.yaml`，并且完善对应的配置修改即可
> 注意：首次部署需访问rabbitmq管理端创建virtual host
#### rabbitmq配置
访问`http://127.0.0.1:15672`，找到`Virtual Hosts`创建`sim`虚拟主机，可参考[官方文档](https://www.rabbitmq.com/vhosts.html)


### 项目启动
#### 启动命令
```shell
make env-up     // 启动docker容器环境
make im         // 启动im服务
make user       // 启动user服务
make gateway    // 启动网关服务
make env-down   // 关闭容器环境
make proto      // 生成rpc服务所需的protoc文件
```

### 启动项目
#### step1 启动项目环境
`make env-up`
#### step2 启动各服务
`make im && make user && make gateway`

最后访问`http://localhost:1024`即可
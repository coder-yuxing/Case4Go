# Bookstore 
一个图书管理API服务, 提供针对图书的CRUD（创建、检索、更新、删除）的基于http的API。
API采用RESTful风格设计，服务提供的API如下：
![api list](https://github.com/coder-yuxing/Case4Go/blob/main/projects/bookstore/docs/img/api_list.png?raw=true)

## 项目结构
```azure
├── cmd/
│   └── bookstore/         // 放置bookstore main包源码
│       └── main.go
├── go.mod                 // module bookstore的go.mod
├── go.sum
├── internal/              // 存放项目内部包的目录
│   └── store/
│       └── memstore.go     
├── server/                // HTTP服务器模块
│   ├── middleware/
│   │   └── middleware.go
│   └── server.go          
└── store/                 // 图书数据存储模块
    ├── factory/
    │   └── factory.go
    └── store.go
```
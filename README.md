
go-micro-微服务框架

### 启动etcdv3集群

这里我用两台机器三个节点([192.168.3.45:2379, 192.168.3.45:2479],  [192.168.3.118:2379])
#### node1
```bash
etcd --name cd0 \
--initial-advertise-peer-urls http://192.168.3.45:2380 \
--listen-peer-urls http://192.168.3.45:2380 \
--listen-client-urls http://192.168.3.45:2379 \
--advertise-client-urls http://192.168.3.45:2379 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster cd0=http://192.168.3.45:2380,cd1=http://192.168.3.45:2480,cd2=http://192.168.3.118:2380 \
--initial-cluster-state new
```
#### node2
```
etcd --name cd1 \
--initial-advertise-peer-urls http://192.168.3.45:2480 \
--listen-peer-urls http://192.168.3.45:2480 \
--listen-client-urls http://192.168.3.45:2479 \
--advertise-client-urls http://192.168.3.45:2479 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster cd0=http://192.168.3.45:2380,cd1=http://192.168.3.45:2480,cd2=http://192.168.3.118:2380 \
--initial-cluster-state new

```
#### node3
```
etcd --name cd2 \
--initial-advertise-peer-urls http://192.168.3.118:2380 \
--listen-peer-urls http://192.168.3.118:2380 \
--listen-client-urls http://192.168.3.118:2379 \
--advertise-client-urls http://192.168.3.118:2379 \
--initial-cluster-token etcd-cluster-1 \
--initial-cluster cd0=http://192.168.3.45:2380,cd1=http://192.168.3.45:2480,cd2=http://192.168.3.118:2380 \
--initial-cluster-state new
```
#### 测试etcd集群是否正常
[![F6CCEAF6-4C23-4599-97E4-CBD35CB029FD.png](https://i.loli.net/2019/06/20/5d0b1c214bfab55031.png)](https://i.loli.net/2019/06/20/5d0b1c214bfab55031.png)

### micro安装etcd插件

在$GOPATH/src/github.com/micro/micro 下新建plugins.go

```
package main

import (
    _ "github.com/micro/go-plugins/registry/etcdv3"
)
```
```
$ go build -i -o micro ./main.go ./plugins.go
$ mv micro $GOBIN/
$ micro --registry=etcdv3 list services
```

### 启动服务

#### clone代码

git clone git@github.com:weiwenwang/go-mcro-demo.git
cd go-micro-demo


#### 启动Micro Api
```
micro --registry=etcdv3 --registry_address=http://192.168.3.45:2479 api --handler=api
```
HTTP API Listening on [::]:8080  监听到本地8080端口

#### 启动greeter Api并注册到etcd

cd go-micro-demo/api
```
go run main.go --registry=etcdv3 --registry_address=http://192.168.3.45:2479
```

#### 启动greeter srv并注册到etcd
cd go-micro-demo/srv
```
go run main.go --registry=etcdv3 --registry_address=http://192.168.3.118:2379
```


#### 逻辑图

这个图是我借用[别处](http://btfak.com/%E5%BE%AE%E6%9C%8D%E5%8A%A1/2016/03/20/micro/)的, 图上的ip和我这不一致, 服务之间的调用完全一样

[![greeter.png](https://i.loli.net/2019/06/19/5d09fbb10b5c556537.png)](https://i.loli.net/2019/06/19/5d09fbb10b5c556537.png)

### 查看etcdv3所有的key

```
$ etcdctl --endpoints="http://192.168.3.45:2379" --prefix --keys-only=true get /
/micro-registry/go.micro.api.greeter/go.micro.api.greeter-749681bd-d47b-4f3e-be28-92014bf559b3

/micro-registry/go.micro.api/go.micro.api-babd3725-c520-4d3e-8abd-fdb95b93d663

/micro-registry/go.micro.srv.greeter/go.micro.srv.greeter-c96c6792-2d46-4932-8b00-2b4f3db7624e
```

### 浏览器访问Micro Api
http://localhost:8080/greeter/say/hello?name=John

```
{
    "id": "go.micro.client",
    "code": 500,
    "detail": "error selecting go.micro.srv.greeter node: not found",
    "status": "Internal Server Error"
}
```

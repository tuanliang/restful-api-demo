# BOOK管理系统样例

微服务开发的时候，A，B,C
都需要做分页,都需要更新模式，把这些公共功能，抽象公共库

## 公共库

import把需要功能导入进行

## 功能的protobuf定义

import protobuf公共库的时候，需要版本对应上

比如我们公开库 demo1 v1，生成的代码就是v1，依赖的v1，我们依赖的protobuf 也是v1

这些库的位置:
  Page公共protobuf: ${GOMODCACHE}/github.com/infraboard/mcube@version /pb/page
  Request protobuf: UpdateMode ${GOMODCACHE}/github.com/infraboard/mcube@version /pb/request

## 第一种解决方案

把依赖的protobuf的对应的版本存放本地的 /usr/local/include 下面，通过-I=/usr/local/include 指定外部依赖库的位置

类似于你把这个公共的protobuf 作为了一个全局的版本，如果你本地有2个项口，依赖2个版本，这个时候放在全局就法解决

## 把依赖放入项目里面

把依赖的protobuf copy到项目内部，然后通过-I=common(自己建的文件夹)/pb，不同项П就可以依赖指定的版木，互补干扰

项目A: mcube@v1.1
项目B: mcube@v1.2

有没有这个的工具，能把对应版本的外部库的protobuf定义 copy项目里面？应该是没有

## protobuf外部依赖copy

项目依赖的protobuf的版木是v1.9.7 (go.mod)，需要copy就是v1.9.7 protobuf1．找到项目依赖的外部公共库的版本

1.找到项目依赖的外部公共库的版本
  得到版本号：go list -m module_name ,eg(go list -m github.com/infraboard/mcube) v1.9.7
    ```shell
      go list -m github.com/infraboard/mcube
    ```
  或者看go.mod
2.找到该包的本地地址
+ go mod存在的位置
+ 包名称
+ 版木号

```proto
import "github.com/infraboard/mcube/pt/request/request.proto
```

3.确定目标copy位置
copy的目的地址:common/pb/ github.com/infraboard/mcube/pb
+ 搜索位置:-I=common/pb
+ 包前缀:github.com/infraboard/mcube/pb

4.copy对应版本的protobuf文件到当前目录

5.清理多余的Go文件


1-5小总结：简单来说：执行以下步骤即可
```shell
mkdir -pv common/pb/github.com/infraboard/mcube/pb
```
go env --> 找到自己的GOMODCACHE目录（E:\go_path\pkg\mod）
```shell
go env
```
```shell
cp -r /e/go_path/pkg/mod/github.com/infraboard/mcube\@v1.9.7/pb/* common/pb/github.com/infraboard/mcube/pb/
```
```shell
rm -rf common/pb/github.com/infraboard/mcube/pb/*/*.go
```

6．使川这些依赖protobuf来生成代码:

```shell
protoc -I=. -I=common/pb --go_out=. --go_opt=module="github.com/tuanliang/restful-api-demo" 
--go-grpc_out=. --go-grpc_opt=module="github.com/tuanliang/restful-api-demo" apps/*/pb/*.proto
```

```shell
go fmt ./...
```

```shell
protoc-go-inject-tag -input=apps/*/*.pb.go
```

```shell

```
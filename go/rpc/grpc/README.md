
# grpc 

## 文档

>- https://github.com/grpc/grpc-go
>- https://www.grpc.io/docs/protoc-installation/
>- https://developers.google.com/protocol-buffers/docs/gotutorial
>- https://www.grpc.io/docs/languages/go/quickstart/


## 安装

>- https://github.com/protocolbuffers/protobuf/releases 下载对应的系统版本
>- go get -v -u google.golang.org/protobuf/cmd/protoc-gen-go
>- go get -v -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

### 未使用 go modules
>- go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
>- go get -u google.golang.org/grpc
>- protoc --go_out=plugins=grpc:. *.proto

## 生成文件

>- protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto
>- protoc --gou_out=. --go-grpc_out=. *.proto


## grpc 格式定义

```
service 服务名{
	rpc 函数名(入参:消息体) returns (出参:消息体){}
}
```

## proto 文件格式demo

```
//设置版本，默认是 proto2
syntax = "proto3";

option go_package="/pb";

//包名
package pb;


//枚举类型, 枚举值必须从 0 开始
enum Gender {
	Man = 0;
	Woman = 1;
}

//消息体
message Student {
	int32 id = 1;   //后面数字不可以重复
	string name = 2;
	Gender gender = 3;
	repeated int32 scores = 4; //数组->切片
	Company compay = 5; //嵌套
	
	//联合体
	oneof data {
		string teacher = 6;
		string class = 7;
	}
}

//消息体
message Company {
	string Name = 1;
}

//添加rpc定义
service hello{
	rpc Say(Student) returns (Student){}
}
```
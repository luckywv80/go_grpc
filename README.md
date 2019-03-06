# go_grpc
##protouf

**编译安装protobuf**
  
  1、从https://github.com/google/protobuf/releases 获取 Protobuf 编译器 protoc  
  tar zxvf protobuf-all-3.7.0.tar.gz
  
  2、cd protobuf-3.7.0
  
  3、./confiure --prefix=/usr/local/protobuf
  
  4、make
  
  5、make install
  
**修改环境变量**

    1、vim ~/.bash_profile文件
  
    2、输入如下 export PROTOBUF=/usr/local/protobuf export PATH=$PATH:$PROTOBUF/bin
  
    3、source ~/.bash_profile 使环境变量生效
  
## protoc-gen-go

    该插件用于编译 .proto 文件为 Golang 源文件
   
**安装protoc-gen-go**
   
    1、执行 go get github.com/golang/protobuf/protoc-gen-go 命令
   
    2、cd $GOPATH/src/github.com/golang/protobuf/protoc-gen-go
   
    3、go build
   
    4、go install

**修改protoc-gen-go环境变量**
        
    1、vim ~/.bash_profile文件   
    
    2、输入如下export GOPATH=/Users/lilonggen/dev/go export PATH=$PATH:$GOPATH/bin 其中GOPATH为go安装时的path    
   
    3、source ~/.bash_profile 使环境变量生效

## goprotobuf 提供的支持库

    该库主要包含编码、解码等功能
    
**安装proto**    
    
    1、go get github.com/golang/protobuf/proto
    
    2、cd $GOPATH/src/github.com/golang/protobuf/proto
    
    3、go build
    
    4、go install
    
## 安装grpc

    1、git clone https://github.com/grpc/grpc-go.git $GOPATH/src/google.golang.org/grpc

    2、git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
    
    3、git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x/text
    
    4、git clone https://github.com/google/go-genproto.git $GOPATH/src/google.golang.org/genproto
    
    5、cd $GOPATH/src/
    
    6、go install google.golang.org/grpc
  
## go使用grpc
    
    1、到blog目录中执行如下命令
    2、protoc -I blogpb/  blogpb/blog.proto --go_out=plugins=grpc:blogpb 会在blogpb文件夹下生成blog.pb.go文件

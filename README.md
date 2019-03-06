# go_grpc
# 环境搭建
-- 编译安装protobuf

1.从 https://github.com/google/protobuf/releases 获取 Protobuf 编译器 protoc


go 使用grpc例子
到blog目录中执行如下命令
   protoc -I blogpb/  blogpb/blog.proto --go_out=plugins=grpc:blogpb 
会在blogpb文件夹下生成blog.pb.go文件

# go_grpc
go 使用grpc例子
到blog目录中执行如下命令
   protoc -I blogpb/  blogpb/blog.proto --go_out=plugins=grpc:blogpb 
会在blogpb文件夹下生成blog.pb.go文件

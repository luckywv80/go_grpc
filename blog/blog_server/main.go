package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"grpc_test/blog/blogpb"
	"log"
	"net"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"os/signal"
)

var Conns *sql.DB

type server struct {
}

type blogItem struct {
	Id 	int64 `json:"id" form:"id"`
	AuthorId int64 `json:"author_id" form:"author_id"`
	Title string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
}

func (*server) CreateBlog(cxt context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Create blog request")

	fmt.Println(Conns)

	blog := req.GetBlog()
	data := blogItem{
		AuthorId: 	blog.GetAuthorId(),
		Title:		blog.GetTitle(),
		Content:	blog.GetContent(),
	}
	res, err := Conns.Exec("INSERT INTO blog (author_id, title, content) VALUES (?, ?, ?)", data.AuthorId, data.Title, data.Content)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       id,
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}

func (*server) ReadBlog(cxt context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("Read blog request")
	blog := blogItem{}
	blogId := req.GetBlogId()
	err := Conns.QueryRow("SELECT * FROM blog WHERE id=? LIMIT 1", blogId).Scan(&blog.Id, &blog.AuthorId, &blog.Title, &blog.Content)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}
	return &blogpb.ReadBlogResponse{
		Blog:dataToBlogPb(&blog),
	}, nil
}

func (*server) UpdateBlog(cxt context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	fmt.Println("update blog request")
	blog := req.GetBlog()

	data := &blogItem{}
	data.Id = blog.GetId()
	data.AuthorId = blog.GetAuthorId()
	data.Title = blog.GetTitle()
	data.Content = blog.GetContent()

	res, err := Conns.Prepare("UPDATE blog SET author_id=?,title=?,content=? WHERE id=?")
	defer res.Close()
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	rs, err := res.Exec(data.Id, data.Title, data.Content, data.Id)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}
	_, err = rs.RowsAffected()
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}
	return &blogpb.UpdateBlogResponse{
		Blog: dataToBlogPb(data),
	}, nil
}

func (*server) DeleteBlog(cxt context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	fmt.Println("Delete blog request")
	blogId := req.GetBlogId()
	rs, err := Conns.Exec("DELETE FROM blog WHERE id=?", blogId)
	if err != nil {
		log.Fatalln(err)
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}
	_, err = rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}
	return &blogpb.DeleteBlogResponse{BlogId: req.GetBlogId()}, nil
}

func (*server) ListBlog(req *blogpb.ListBlogRequest, stream blogpb.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")

	rows, err := Conns.Query("SELECT * FROM blog")
	defer rows.Close()
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	for rows.Next() {
		data := &blogItem{}
		_ = rows.Scan(&data.Id, &data.AuthorId, &data.Title, &data.Content)

		_ = stream.Send(&blogpb.ListBlogResponse{Blog: dataToBlogPb(data)})
	}
	return nil
}

func dataToBlogPb(data *blogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.Id,
		AuthorId: data.AuthorId,
		Content:  data.Content,
		Title:    data.Title,
	}
}

func main() {
	fmt.Print("Connecting to mysql")

	Conns, _ = sql.Open("mysql", "root:123456@tcp(192.168.33.11)/test")
	fmt.Println(Conns)

	err := Conns.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Blog Service Started")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	blogpb.RegisterBlogServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	Conns.Close()
	fmt.Println("End of Program")
}

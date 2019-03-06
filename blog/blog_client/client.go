package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_test/blog/blogpb"
	"io"
	"log"
)

func main() {
	fmt.Println("Blog Client")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	// create Blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: 1,
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()
	fmt.Println(blogID)

	// read Blog
	fmt.Println("Reading the blog")
	readBlog, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: 2})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v \n", err2)
	}
	fmt.Println(readBlog)

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, readBlogErr := c.ReadBlog(context.Background(), readBlogReq)
	if readBlogErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readBlogErr)
	}
	fmt.Printf("Blog was read: %v \n", readBlogRes)

	//update blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: 2,
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}
	updateRes, updateErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf("Blog was updated: %v\n", updateRes)

	//delete blog
	deleteRes, deleteErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	//list blog
	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}

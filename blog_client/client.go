package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/flaco-segura/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// Create blog
	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Flaco",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}
	createBlogRes, errCrt := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if errCrt != nil {
		log.Fatalf("Unexpected error: %v", errCrt)
	}
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// Read Blog
	fmt.Println("Reading the blog")

	// Checking NOT FOUND
	_, errRd := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "s0m3f4k31D"})
	if errRd != nil {
		fmt.Printf("Error happened while reading: %v", errRd)
	}

	readBlogRes, errRd := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogID})
	if errRd != nil {
		fmt.Printf("Error happened while reading: %v", errRd)
	}
	fmt.Printf("Blog %v is: %v\n", blogID, readBlogRes)

	// Update Blog
	fmt.Println("Update the blog")
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Flaco",
		Title:    "My n blog",
		Content:  "Content of the n blog",
	}

	updRes, errUpd := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if errUpd != nil {
		log.Fatalf("Unexpected error: %v", errUpd)
	}
	fmt.Printf("Blog %v was updated: %v\n", blogID, updRes)

	// Delete Blog
	deleteRes, errDlt := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if errDlt != nil {
		log.Fatalf("Unexpected error: %v", errDlt)
	}
	fmt.Printf("Blog %v was deleted\n", deleteRes)

	// List Blog
	stream, errLst := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if errLst != nil {
		log.Fatalf("Error while calling ListBlog: %v", errLst)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Println(res.GetBlog())
	}
}

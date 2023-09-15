package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "grpc-mongo/proto"
)

func CallListGroup(ctx context.Context, client pb.Blog_ServiceClient) {

	stream, err := client.ListBlog(ctx, &pb.ListBlogReq{})

	if err != nil {
		log.Printf("Error %v", err)
		return
	}

	for {
		message, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error While Receving the Stream %v", err)
		}

		newBlog := message.Blog

		log.Printf("%s", newBlog)
	}

	println("Streaming Finished")
}

func main() {

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to Connect to the Server %v", err)
	}

	Client := pb.NewBlog_ServiceClient(conn)

	ctx := context.Background()

	// Create a Blog message to send
	// blog := &pb.CreateBlogReq{Blog: &pb.Blog{
	// 	Id:       "3",
	// 	AuthorId: "author123",
	// 	Title:    "Sample Blog",
	// 	Content:  "This is a sample blog post.",
	// },
	// }

	//client.CreateBlog(ctx, blog)

	//CallListGroup(ctx, Client)

	res, err := Client.ReadBlog(ctx, &pb.ReadBlogReq{Id: "1"})
	if err != nil {
		fmt.Println("Error")
		return
	}
	bon := res.Blog

	fmt.Println(bon)

	riz, viz := Client.UpdateBlog(ctx, &pb.UpdateBlogReq{Blog: &pb.Blog{
		Id:       "2",
		AuthorId: "author123",
		Content:  "Rizzi G",
		Title:    "King of FaisalAbad",
	}})

	if viz != nil {
		fmt.Println("Error")
		return
	}
	rizviz := riz.Blog

	fmt.Println(rizviz)

}

package main

import (
	"context"
	pb "grpc-mongo/proto"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Blog struct {
	ID       string `bson:"id,omitempty"`
	AuthorID string `bson:"author_id"`
	Content  string `bson:"content"`
	Title    string `bson:"title"`
}

func (s *helloServer) CreateBlog(ctx context.Context, req *pb.CreateBlogReq) (*pb.CreateBlogRes, error) {

	blog := req.Blog

	newBlog := &Blog{
		ID:       blog.Id,
		AuthorID: blog.AuthorId,
		Content:  blog.Content,
		Title:    blog.Title,
	}

	if _, err := Blogdb.InsertOne(ctx, newBlog); err != nil {
		return nil, err
	} else {
		log.Printf("Succesfully Created a Blog")
	}

	return &pb.CreateBlogRes{Blog: blog}, nil
}

func (s *helloServer) CreateBlogs(ctx context.Context, req *pb.CreateBlogsReq) (*pb.CreateBlogsRes, error) {

	blog := req.Blog

	var blogs []interface{}

	for _, newBlog := range blog {

		tempBlog := Blog{
			ID:       newBlog.Id,
			AuthorID: newBlog.AuthorId,
			Title:    newBlog.Title,
			Content:  newBlog.Content,
		}
		blogs = append(blogs, tempBlog)
	}

	if _, err := Blogdb.InsertMany(ctx, blogs); err != nil {
		return nil, err
	} else {
		log.Printf("Succesfully Created Blogs")
	}
	return &pb.CreateBlogsRes{Blog: blog}, nil
}

func (s *helloServer) ReadBlog(ctx context.Context, req *pb.ReadBlogReq) (*pb.ReadBlogRes, error) {

	id := req.Id
	var newBlog Blog
	filter := bson.M{"id": id}

	if err := Blogdb.FindOne(ctx, filter).Decode(&newBlog); err != nil {
		return nil, err
	} else {
		log.Printf("Found the Searched Blog")
	}
	return &pb.ReadBlogRes{
		Blog: &pb.Blog{
			Id:       newBlog.ID,
			AuthorId: newBlog.AuthorID,
			Content:  newBlog.Content,
			Title:    newBlog.Title,
		},
	}, nil
}

func (s *helloServer) UpdateBlog(ctx context.Context, req *pb.UpdateBlogReq) (*pb.UpdateBlogRes, error) {

	newBlog := req.GetBlog()

	filter := bson.M{"id": newBlog.GetId()}
	update := bson.M{
		"$set": bson.M{
			"id":        newBlog.GetId(),
			"author_id": newBlog.GetAuthorId(),
			"content":   newBlog.GetContent(),
			"title":     newBlog.GetTitle(),
		},
	}

	result := Blogdb.FindOneAndUpdate(ctx, filter, update)
	//res, err := Blogdb.UpdateOne(ctx, filter, update)
	var decode Blog

	err := result.Decode(&decode)

	if err != nil {
		return nil, err
	}

	log.Printf("Successfully Updated Blog")

	return &pb.UpdateBlogRes{Blog: &pb.Blog{
		Id:       decode.ID,
		AuthorId: decode.AuthorID,
		Content:  decode.Content,
		Title:    decode.Title,
	}}, nil
}

func (s *helloServer) DeleteBlog(ctx context.Context, req *pb.DeleteBlogReq) (*pb.DeleteBlogRes, error) {

	id := req.Id
	filter := bson.M{"id": id}

	res, err := Blogdb.DeleteOne(ctx, filter)

	if res.DeletedCount != 1 {
		return nil, err
	}

	log.Printf("Successfully Delted a Blog")

	return &pb.DeleteBlogRes{Success: true}, nil

}

func (s *helloServer) ListBlog(req *pb.ListBlogReq, stream pb.Blog_Service_ListBlogServer) error {

	filter := bson.M{}

	cursor, err := Blogdb.Find(context.Background(), filter)

	if err != nil {
		log.Printf("Error in Opening the Database: %v", err)
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var blog Blog // Declare the blog struct inside the loop

		if err := cursor.Decode(&blog); err != nil {
			log.Printf("Error in Decoding Cursor: %v", err)
			return err
		}

		if err := stream.Send(&pb.ListBlogRes{
			Blog: &pb.Blog{
				Id:       blog.ID,
				AuthorId: blog.AuthorID,
				Content:  blog.Content,
				Title:    blog.Title,
			},
		}); err != nil {
			log.Printf("Error sending data to the client: %v", err)
			return err
		}
		time.Sleep(2 * time.Second)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Error in Parsing the Database: %v", err)
		return err
	}

	return nil
}

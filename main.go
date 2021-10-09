package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type User struct {
	ID       int    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

type Post struct {
	ID        int       `json:"id,omitempty" bson:"id,omitempty"`
	Caption   string    `json:"caption,omitempty" bson:"caption,omitempty"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty"`
	Timestamp time.Time `json:"time,omitempty" bson:"time,omitempty"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var client *mongo.Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	h := sha256.New()

	collection := client.Database("instagram").Collection("users")
	user1 := User{1, "Girish", "girish@gmail.com", "pass1"}
	user2 := User{2, "Balaji", "girish1@gmail.com", "pass2"}
	user3 := User{3, "Tharun", "girish2@gmail.com", "pass3"}

	users := []interface{}{user1, user2, user3}

	insertManyResult, err := collection.InsertMany(context.TODO(), users)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.NewRequest("POST", "http://localhost:3000/users", bytes.NewBuffer((h.Sum([]byte(user1.Password)))))

	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))

	fmt.Println("Inserted multiple users: ", insertManyResult.InsertedIDs)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var result User
	var client *mongo.Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.NewRequest("GET", "localhost:3000/user/", nil)
	if err != nil {
		print(err)
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /user/<id> in task handler", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := resp.URL.Query()
	q.Add("userid", "1")

	collection := client.Database("instagram").Collection("users")
	filter := bson.D{{"_id", id}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single user: %+v\n", result)

}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var client *mongo.Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("instagram").Collection("posts")
	tm := time.Now()
	post1 := Post{1, "this is post 1", "https://images.unsplash.com/photo-1515091943-9d5c0ad475af?ixid=MnwxMjA3fDB8MHxzZWFyY2h8OHx8aW5kaWF8ZW58MHx8MHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60", tm}
	post2 := Post{2, "this is post 1", "https://images.unsplash.com/photo-1547707188-cdbffa0bc270?ixid=MnwxMjA3fDB8MHxzZWFyY2h8NXx8aW5kaWF8ZW58MHx8MHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60", tm}
	post3 := Post{3, "this is post 1", "https://images.unsplash.com/photo-1532375810709-75b1da00537c?ixid=MnwxMjA3fDB8MHxzZWFyY2h8Mnx8aW5kaWF8ZW58MHx8MHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60", tm}
	post4 := Post{3, "this is post 2", "https://images.unsplash.com/photo-1524492412937-b28074a5d7da?ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8aW5kaWF8ZW58MHx8MHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60", tm}
	post5 := Post{3, "this is post 3", "https://images.unsplash.com/photo-1524492412937-b28074a5d7da?ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8aW5kaWF8ZW58MHx8MHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60", tm}
	post6 := Post{4, "this is post 4", "https://images.unsplash.com/photo-1524492412937-b28074a5d7da?ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8aW5kaWF8ZW58MHx8MHx8&ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60", tm}
	posts := []interface{}{post1, post2, post3, post4, post5, post6}

	insertManyResult, err := collection.InsertMany(context.TODO(), posts)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.NewRequest("POST", "localhost:3000/posts", bytes.NewBuffer([]byte(post2.Caption)))

	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}
	fmt.Println(string(body))

	fmt.Println("Inserted multiple posts: ", insertManyResult.InsertedIDs)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	var result Post
	var client *mongo.Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.NewRequest("GET", "localhost:3000/post/", nil)
	if err != nil {
		print(err)
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 2 {
		http.Error(w, "expect /post/<id> in task handler", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := resp.URL.Query()
	q.Add("postid", "3")
	resp.URL.RawQuery = q.Encode()

	collection := client.Database("instagram").Collection("posts")
	filter := bson.D{{"id", id}}
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found post: %+v\n", result)
}

func ListAllPost(w http.ResponseWriter, r *http.Request) {
	var client *mongo.Client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.NewRequest("GET", "localhost:3000/posts/users/", nil)
	if err != nil {
		print(err)
	}

	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "expect /posts/users/<id> in task handler", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := resp.URL.Query()
	q.Add("postid", "3")
	resp.URL.RawQuery = q.Encode()

	collection := client.Database("instagram").Collection("posts")
	filter := bson.D{{"id", id}}

	findOptions := options.Find()
	findOptions.SetLimit(3)

	var results []Post
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found multiple posts corresponding to the particular user! : %+v\n", results)
}

func handleRequests() {
	http.HandleFunc("/users", CreateUser)
	http.HandleFunc("/user/", GetUser)
	http.HandleFunc("/posts", CreatePost)
	http.HandleFunc("/post/", GetPost)
	http.HandleFunc("/posts/users/", ListAllPost)
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {

	handleRequests()
	fmt.Println("Appointy is my dream company")

}

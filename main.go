package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Record contains data for urls in database
type Record struct {
	ShortURL string `bson:"_id,omitempty"`
	LongURL  string `bson:"longurl"`
}

var (
	client *mongo.Client
)

func Index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Welcome!")
}

func Hello(ctx *fasthttp.RequestCtx) {
	var (
		contex context.Context
		cancel context.CancelFunc
	)

	collection := client.Database("Pretiffy").Collection("urls")
	var record bson.M
	err := collection.FindOne(contex, bson.M{"_id": ctx.UserValue("short")}).Decode(&record)
	if err == nil {
		contex, cancel = context.WithTimeout(context.Background(), 15*time.Second)
	} else {
		contex, cancel = context.WithCancel(context.Background())
	}
	ctx.Redirect(record["longurl"].(string), http.StatusSeeOther)
	defer cancel()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	// Set client options
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_DB"))

	var err error
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err = mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	hash := generateRandomHash(10)
	insertUrl(hash, "https://www.facebook.com/")

	router := router.New()
	router.GET("/", Index)
	router.GET("/{short}", Hello)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

}

// generate random base 64 with length n
func generateRandomHash(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// insert record {short, long} in mongodb client
func insertUrl(short string, long string) error {
	record := Record{short, long}

	collection := client.Database("Pretiffy").Collection("urls")

	instertResult, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		return fmt.Errorf("failed to insert short and long url to prettify collection")
	}

	fmt.Println("Inserted post with ID :", instertResult.InsertedID)
	return nil
}

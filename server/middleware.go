package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Record contains data for urls in database
type Record struct {
	ShortURL string `bson:"_id,omitempty"`
	LongURL  string `bson:"longurl"`
}

const (
	dbName         = "Pretiffy"
	collectionName = "urls"
	hostURL        = "localhost:8000/"
)

var (
	collection *mongo.Collection
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_DB"))

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)

	}

	fmt.Println("MongoDB connected")

	collection = client.Database(dbName).Collection(collectionName)
}

// Encode get request from user with long url
// and gives user short url
func Encode(w http.ResponseWriter, r *http.Request) {
	hash := generateRandomHash(10)
	err := json.NewEncoder(w).Encode(hostURL + hash)
	if err != nil {
		log.Printf("could not encode response to output: %v", err)
	}
	var input struct {
		URL string `json:"url"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &input)
	if err != nil {
		log.Fatal("cant unmarshall body respones : %w", err)
	}

	insertURL(hash, input.URL)

}

// Short url get long url from db
func Short(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var record bson.M
	collection.FindOne(context.Background(), bson.M{"_id": vars["short"]}).Decode(&record)

	if record != nil {
		http.Redirect(w, r, record["longurl"].(string), http.StatusSeeOther)
	}

}

// insert record {short, long} in mongodb client
func insertURL(short string, long string) error {
	record := Record{short, long}

	instertResult, err := collection.InsertOne(context.Background(), record)
	if err != nil {
		return fmt.Errorf("failed to insert short and long url to prettify collection")
	}
	fmt.Println("Inserted post with ID :", instertResult.InsertedID)
	return nil
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

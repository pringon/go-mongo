package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

type TodoDocument struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
}

type Todo struct {
	Title       string
	Description string
}

func Connect(options *options.ClientOptions) (*mongo.Client, error) {

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), options)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func ToDoc(v interface{}) (doc *bson.D, err error) {
	data, err := bson.Marshal(v)
	if err != nil {
		return
	}

	err = bson.Unmarshal(data, &doc)
	return
}

func GetAll(client *mongo.Collection) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		cur, err := client.Find(context.Background(), bson.D{})
		if err != nil {
			log.Print(err)
			return
		}
		defer cur.Close(context.Background())

		var todoItems []TodoDocument

		for cur.Next(context.Background()) {
			var todoItem TodoDocument
			err := cur.Decode(&todoItem)
			if err != nil {
				log.Print(err)
				return
			}

			todoItems = append(todoItems, todoItem)
		}
		js, err := json.Marshal(todoItems)
		if err != nil {
			log.Print(err)
			return
		}

		w.Write(js)
	}
}

func GetById(client *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		todoId, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			log.Print(err)
			return
		}

		var todoItem TodoDocument
		filter := bson.D{{"_id", todoId}}
		err = client.FindOne(context.Background(), filter).Decode(&todoItem)
		if err != nil {
			log.Print(err)
			return
		}
		js, err := json.Marshal(todoItem)
		if err != nil {
			log.Print(err)
			return
		}
		w.Write(js)
	}
}

func CreateTodo(client *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var todoItem TodoDocument
		_ = json.NewDecoder(r.Body).Decode(&todoItem)
		todoItem.ID = primitive.NewObjectID()

		_, err := client.InsertOne(context.Background(), todoItem)
		if err != nil {
			log.Print(err)
			return
		}
		js, err := json.Marshal(todoItem)
		if err != nil {
			log.Print(err)
			return
		}
		w.Write(js)
	}
}

func UpdateTodo(client *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)

		var todoItem TodoDocument
		_ = json.NewDecoder(r.Body).Decode(&todoItem)
		var err error
		todoItem.ID, err = primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			log.Print(err)
			return
		}

		todoDoc, err := ToDoc(todoItem)
		if err != nil {
			log.Print(err)
			return
		}

		filter := bson.D{{"_id", todoItem.ID}}
		update := bson.D{{"$set", todoDoc}}

		_, err = client.UpdateOne(context.Background(), filter, update)
		if err != nil {
			log.Print(err)
			return
		}

		js, err := json.Marshal(todoItem)
		if err != nil {
			log.Print(err)
			return
		}
		w.Write(js)
	}
}

func DeleteTodo(client *mongo.Collection) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		itemId, err := primitive.ObjectIDFromHex(params["id"])
		if err != nil {
			log.Print(err)
			return
		}

		filter := bson.D{{"_id", itemId}}

		_, err = client.DeleteOne(context.Background(), filter)
		if err != nil {
			log.Print(err)
			return
		}

		w.Write([]byte("{ code: \"DELETE_COMPLETE\"}"))
	}
}

function buildMongoURI(user string, password string) string {
        return fmt.Sprintf("mongodb+srv://%s:%s@cluster0-eecpk.gcp.mongodb.net/test?retryWrites=true&w=majority", mongoUser, mongoPass)
}

func main() {
	mongoUser := os.Getenv("MONGO_USER")
	mongoPass := os.Getenv("MONGO_PASSWORD")
	client, err := mongo.NewClient(options.Client().ApplyURI(buildMongoURI(mongoUser, mongoPass)))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
        if err := client.Connet(ctx); err != nil {
		log.Fatal(err)
	}

	todo := client.Database("test").Collection("Todo")

	// Router initialisation
	r := mux.NewRouter()

	// Routes config
	r.HandleFunc("/{id}", GetById(todo)).Methods("GET")
	r.HandleFunc("/", GetAll(todo)).Methods("GET")
	r.HandleFunc("/", CreateTodo(todo)).Methods("POST")
	r.HandleFunc("/{id}", UpdateTodo(todo)).Methods("PUT")
	r.HandleFunc("/{id}", DeleteTodo(todo)).Methods("DELETE")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

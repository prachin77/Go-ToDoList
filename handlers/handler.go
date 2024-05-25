package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var id primitive.ObjectID

type Task struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	TaskValue string             `json:"taskvalue"`
	TaskDate  string             `json:"taskdate"`
	// Add other fields as needed
}

const (
	connectionString = "mongodb://localhost:27017"
	dbName           = "timepass"
	collName         = "todolist"
)

// this is a pointer(reference) to collection in mongo db
var collection *mongo.Collection

func init() {
	clientOpt := options.Client().ApplyURI(connectionString)

	// connect to mongodb
	client, err := mongo.Connect(context.TODO(), clientOpt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection to mongo db successfull ✌️✌️")

	collection = client.Database(dbName).Collection(collName)

	// collection instance
	fmt.Println("collection instance is ready")
}

func DefaultRoute(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/main.html"))
	tmpl.Execute(w, nil)
}

func GetApp(w http.ResponseWriter, r *http.Request) {
	// get all task from mongo db
	tasks, err := GetAllTasks(w, r)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("tasks from app page : ", tasks)
		tmpl := template.Must(template.ParseFiles("./templates/todo.html", "./templates/todocompo.html"))
		tmpl.Execute(w, tasks)
		// tmpl.Execute(w, nil)
	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) ([]Task, error) {
	// Create a slice to hold the tasks
	var tasks []Task

	// Define a context for the operation
	ctx := context.TODO()

	// Define options to customize the query
	findOptions := options.Find()

	// Find all documents in the collection with specified options
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into a Task
	for cursor.Next(ctx) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Return the slice of tasks
	return tasks, nil
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	id = primitive.NewObjectID()
	task := Task{
		ID:        id,
		TaskValue: r.PostFormValue("addtask"),
		TaskDate:  time.Now().Format("2006-01-02"),
	}
	// Insert the task into the MongoDB collection
	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		// http.Error(w, "Failed to add task to database", http.StatusInternalServerError)
		log.Fatal(err)
		// return
	} else {
		tmpl := template.Must(template.ParseFiles("./templates/todocompo.html"))
		tmpl.Execute(w, task)
	}

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract task ID from request URL
	taskId := mux.Vars(r)["id"]
	fmt.Println("task id : ", taskId)

	// Convert task ID string to ObjectID
	objId, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("task id with hex converted value = ", objId)
	filter := bson.M{"_id": objId}

	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

}

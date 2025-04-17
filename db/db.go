package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DB is the global MongoDB database instance
var DB *mongo.Database
// EmployeeCollection is the MongoDB collection for employee records
var EmployeeCollection *mongo.Collection

// ConnectDB establishes a connection to MongoDB and initializes DB and collection variables
func ConnectDB(){
	// Create a context with a 2-second timeout for database operations
	ctx,cancel:=context.WithTimeout(context.Background(),2*time.Second)
	defer cancel()
	
	// Load environment variables from .env file
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal("error in loading .env file")
	}

	// Get MongoDB connection string from environment variables
	mongoURI:=os.Getenv("MONGO_URI")
	if mongoURI==""{
		log.Fatal("no mongoURI available!")
	}

	// Configure MongoDB client options with the connection string
	clientOptions:=options.Client().ApplyURI(mongoURI)

	// Establish connection to MongoDB
	client,err:=mongo.Connect(ctx,clientOptions)
	if err!=nil{
		log.Fatal("error in connecting to mongoDB")
	}

	// Verify connection with a ping to the primary node
	err=client.Ping(ctx,readpref.Primary())
	if err!=nil{
		log.Fatal("Attempt to connected but ping failed")
	}

	// Initialize the database and collection
	DB=client.Database("employee_db")
	EmployeeCollection=DB.Collection("employee_collection")
	log.Println("Connected to mongodb",EmployeeCollection.Name())
}
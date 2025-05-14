package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// DB is the global MongoDB database instance
var DB *mongo.Database
// EmployeeCollection is the MongoDB collection for employees
var EmployeeCollection *mongo.Collection

// ConnectDB establishes a connection to MongoDB
func ConnectDB(){
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	
	if os.Getenv("DOCKER_ENV") != "true" {
        err := godotenv.Load()
        if err != nil {
            log.Println("No .env file found or failed to load .env (not in Docker).")
        }
    }

	mongoURI:=os.Getenv("MONGO_URI")
	if mongoURI==""{
		log.Fatal("no mongoURI available!")
	}
	//log.Print("Attempting to connect to MongoDB with URI")

	clientOptions:=options.Client().ApplyURI(mongoURI)

	client,err:=mongo.Connect(ctx,clientOptions)
	if err!=nil{
		log.Fatal("error in connecting to mongoDB")
	}

	err=client.Ping(ctx,readpref.Primary())
	if err!=nil{
		log.Fatalf("Attempt to connected but ping failed: %v",err)
	}

	DB=client.Database(os.Getenv("DATABASE_NAME"))
	EmployeeCollection=DB.Collection(os.Getenv("COLLECTION_NAME"))
	log.Println("Connected to mongodb",EmployeeCollection.Name())

	indexModel:=mongo.IndexModel{
		Keys: bson.D{
			{Key:"email", Value: "text"},
			{Key:"name", Value: "text"},
		},
	}
	
	_,err=EmployeeCollection.Indexes().CreateOne(ctx,indexModel)
	if err!=nil{
		log.Println("error in creating index for employee collection",err)
	} else{
		log.Println("Index created for employee collection")
	}

	
}

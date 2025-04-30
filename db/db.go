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
// EmployeeCollection is the MongoDB collection for employees
var EmployeeCollection *mongo.Collection

// ConnectDB establishes a connection to MongoDB
func ConnectDB(){
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal("error in loading .env file")
	}

	mongoURI:=os.Getenv("MONGO_URI")
	if mongoURI==""{
		log.Fatal("no mongoURI available!")
	}
	log.Printf("Attempting to connect to MongoDB with URI: %s", mongoURI)

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
}

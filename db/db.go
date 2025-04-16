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

var DB *mongo.Database
var EmployeeCollection *mongo.Collection

func ConnectDB(){
	ctx,cancel:=context.WithTimeout(context.Background(),2*time.Second)
    defer cancel()
    
	err:=godotenv.Load()
	if err!=nil{
		log.Fatal("error in loading .env file")
	}

	mongoURI:=os.Getenv("MONGO_URI")
	if mongoURI==""{
		log.Fatal("no mongoURI available!")
	}

	clientOptions:=options.Client().ApplyURI(mongoURI)

	client,err:=mongo.Connect(ctx,clientOptions)
	if err!=nil{
		log.Fatal("error in connecting to mongoDB")
	}

	err=client.Ping(ctx,readpref.Primary())
	if err!=nil{
		log.Fatal("Attempt to connected but ping failed")
	}


	
	DB=client.Database("employee_db")
	EmployeeCollection=DB.Collection("employee_collection")
	log.Println("Connected to mongodb",EmployeeCollection.Name())
}
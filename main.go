package main

import (
	"employee/db"
	"employee/routes"
	"log"
	"net/http"
	"os"
)

func main(){
  db.ConnectDB()
  r:=routes.RegisterRoutes()
  log.Println("Server running on port 8080...")
  dockerEnv:=os.Getenv("DOCKER_ENV")
  if dockerEnv=="true"{
	log.Println("Docker environment detected")
  } else{
	log.Println("Not running in Docker")
  }
  log.Fatal(http.ListenAndServe(":8080",r))

  
}
package app

import (
	"employee/db"
	"employee/routes"
	"log"
	"net/http"
	"os"
)

func Start(){
	// Initialize the database connection
	db.ConnectDB()

	// Initialize the router and register routes
	r:=routes.RegisterRoutes()

	// Check if running in Docker
	dockerEnv:=os.Getenv("DOCKER_ENV")
    if dockerEnv=="true"{
	   log.Println("Docker environment detected")
      } else{
    	log.Println("Not running in Docker")
     }

	log.Println("Server running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080",r))
	



}
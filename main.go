package main

import (
	"employee/db"
	"employee/routes"
	"log"
	"net/http"
)

func main(){
  db.ConnectDB()
  r:=routes.RegisterRoutes()
  log.Println("Server running on port 8080...")
  log.Fatal(http.ListenAndServe(":8080",r))
}
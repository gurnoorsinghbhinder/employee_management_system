package app

import (
	"employee/db"
	"employee/routes"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.uber.org/dig"
)

func Start(){
	container:=dig.New()

	//connect db
	container.Provide(db.InitDB)

	//Initialized routes 
	container.Provide(routes.RegisterRoutes)

	

	// Check if running in Docker
   err := container.Invoke(func( r *mux.Router) {
    dockerEnv := os.Getenv("DOCKER_ENV")
    if dockerEnv == "true" {
        log.Println("Docker environment detected")
    } else {
        log.Println("Not running in Docker")
    }
    log.Println("Server running on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", r))
})
    if err != nil {
        log.Fatalf("Failed to start: %v", err)
    }		
}

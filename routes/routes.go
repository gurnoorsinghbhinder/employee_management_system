package routes

import (
	"employee/controller"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router{
	r:=mux.NewRouter()
	r.HandleFunc("/employees",controller.CreateEmployee).Methods("POST")
	r.HandleFunc("/employees",controller.GetAllEmployee).Methods("GET")
	r.HandleFunc("/employees",controller.DeleteAllEmployees).Methods("DELETE")
	r.HandleFunc("/employees/{id}",controller.DeleteEmployee).Methods("DELETE")
	r.HandleFunc("/employees/{id}",controller.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/employees/{id}",controller.GetEmployeeByID).Methods("GET")
	return r
}
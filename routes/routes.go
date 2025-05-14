package routes

import (
	"employee/controller"
	"employee/repository"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router{
	
	repo := repository.NewEmployeeRepository()
	controller := controller.NewEmployeeController(repo)

	router := mux.NewRouter()

	router.HandleFunc("/employees", controller.CreateEmployee).Methods("POST")
	router.HandleFunc("/employees", controller.GetAllEmployee).Methods("GET")
	router.HandleFunc("/employees/search", controller.SearchEmployees).Methods("GET")
	router.HandleFunc("/employees/{id}", controller.UpdateEmployee).Methods("PUT")
	router.HandleFunc("/employees/{id}", controller.DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/employees", controller.DeleteAllEmployees).Methods("DELETE")
	router.HandleFunc("/employees/{id}", controller.GetEmployeeByID).Methods("GET")
	router.HandleFunc("/paginated", controller.GetPaginatedEmployees).Methods("GET")
	router.HandleFunc("/employees/bulk", controller.CreateEmployeesBulk).Methods("POST")
	
	return router
}



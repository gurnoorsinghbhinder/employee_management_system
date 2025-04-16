package controller

import (
	"employee/models"
	"employee/repository"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEmployee(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var employee models.Employee
	json.NewDecoder(r.Body).Decode(&employee)
	err:=repository.CreateEmployee(employee)
	if err!=nil{
		log.Print("error creating employee")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

func GetAllEmployee(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	employee,err:=repository.GetAllEmployees()
	if err!=nil{
		log.Print("error getting all employees")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

func UpdateEmployee(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	idParam:=mux.Vars(r)["id"]
	id,_:=primitive.ObjectIDFromHex(idParam)

	var employee models.Employee
	json.NewDecoder(r.Body).Decode(&employee)

	err:=repository.UpdateEmployees(id,employee)
	if err!=nil{
		log.Print("error updating employees")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)


}

func DeleteEmployee(w http.ResponseWriter,r *http.Request){
	idParam:=mux.Vars(r)["id"]
	id,_:=primitive.ObjectIDFromHex(idParam)

	err:=repository.DeleteEmployee(id)
	if err!=nil{
        log.Print("error deleting employee")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllEmployees(w http.ResponseWriter,r *http.Request){
	err:=repository.DeleteAllEmployees()
	if err!=nil {
		log.Print("error deleting all employees")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func GetEmployeeByID(w http.ResponseWriter,r *http.Request){
	idParam:=mux.Vars(r)["id"]
	id,err:=primitive.ObjectIDFromHex(idParam)
    if err!=nil{
		http.Error(w,"Invalid ID",http.StatusBadRequest)
		return 
	}

	employee,err:=repository.GetEmployeeByID(id)
	if err!=nil{
		log.Print("error getting employee by ID")
		http.Error(w,err.Error(),http.StatusNotFound)
		return 
	}
	json.NewEncoder(w).Encode(employee)
}
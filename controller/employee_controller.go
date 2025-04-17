// Package controller handles HTTP requests for the employee management system.
package controller

import (
	"employee/models"
	"employee/repository"
	"employee/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateEmployee handles POST requests to create a new employee.
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employee models.Employee
	json.NewDecoder(r.Body).Decode(&employee)
	
	if !utils.IsValidEmail(*employee.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	
	taken, err := repository.IsEmailTaken(*employee.Email)
	if err != nil {
		http.Error(w, "Error while checking Email", http.StatusInternalServerError)
		return
	}
	if taken {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	err = repository.CreateEmployee(employee)
	if err != nil {
		log.Print("error creating employee")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// GetAllEmployee retrieves all employees
func GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	employee, err := repository.GetAllEmployees()
	if err != nil {
		log.Print("error getting all employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// UpdateEmployee updates an existing employee by ID
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idParam := mux.Vars(r)["id"]
	id, _ := primitive.ObjectIDFromHex(idParam)

	var employee models.Employee
	if err:=json.NewDecoder(r.Body).Decode(&employee);err!=nil{
		http.Error(w,"invalid inputs",http.StatusBadRequest)
	}

	if employee.Name==nil || strings.TrimSpace(*employee.Name)==""{
        http.Error(w,"name is required",http.StatusBadRequest)
	}

	if employee.Email==nil || strings.TrimSpace(*employee.Email)==""{
		http.Error(w,"email is required",http.StatusBadRequest)
	}



	if !utils.IsValidEmail(*employee.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	taken, err := repository.IsEmailTakenByOther(id,*employee.Email)
	if err != nil {
		http.Error(w, "Error while checking Email", http.StatusInternalServerError)
		return
	}
	if taken {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	err = repository.UpdateEmployees(id, employee)
	if err != nil {
		log.Print("error updating employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// DeleteEmployee removes an employee by ID
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, _ := primitive.ObjectIDFromHex(idParam)

	err := repository.DeleteEmployee(id)
	if err != nil {
		log.Print("error deleting employee")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteAllEmployees removes all employees from the database
func DeleteAllEmployees(w http.ResponseWriter, r *http.Request) {
	err := repository.DeleteAllEmployees()
	if err != nil {
		log.Print("error deleting all employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetEmployeeByID retrieves a specific employee by ID
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 
	}

	employee, err := repository.GetEmployeeByID(id)
	if err != nil {
		log.Print("error getting employee by ID")
		http.Error(w, err.Error(), http.StatusNotFound)
		return 
	}
	json.NewEncoder(w).Encode(employee)
}

// GetPaginatedEmployees retrieves employees with pagination
func GetPaginatedEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		log.Printf("Invalid 'page' query: %v. Falling back to 1.", pageStr)
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		log.Printf("Invalid 'limit' query: %v. Falling back to 10.", limitStr)
		limit = 10
	}

	employees, err := repository.GetPaginatedEmployees(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employees)
}

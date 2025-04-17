// Package controller handles HTTP request processing for the employee management system.
// It contains handlers for CRUD operations on employee resources.
package controller

import (
	"employee/models"
	"employee/repository"
	"employee/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateEmployee handles POST requests to create a new employee.
// It validates the email format and checks if the email is already in use
// before storing the employee data in the database.
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employee models.Employee
	json.NewDecoder(r.Body).Decode(&employee)
	
	// Validate email format
	if !utils.IsValidEmail(employee.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	
	// Check if email is already taken
	taken, err := repository.IsEmailTaken(employee.Email)
	if err != nil {
		http.Error(w, "Error while checking Email", http.StatusInternalServerError)
		return
	}
	if taken {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	// Persist employee to database
	err = repository.CreateEmployee(employee)
	if err != nil {
		log.Print("error creating employee")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// GetAllEmployee handles GET requests to retrieve all employees from the database.
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

// UpdateEmployee handles PUT requests to update an existing employee.
// It validates the request data and checks for email conflicts before updating.
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Extract employee ID from URL
	idParam := mux.Vars(r)["id"]
	id, _ := primitive.ObjectIDFromHex(idParam)

	var employee models.Employee
	json.NewDecoder(r.Body).Decode(&employee)

	// Validate email format
	if !utils.IsValidEmail(employee.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Check if email is already taken by another employee
	taken, err := repository.IsEmailTaken(employee.Email)
	if err != nil {
		http.Error(w, "Error while checking Email", http.StatusInternalServerError)
		return
	}
	if taken {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	// Update employee in database
	err = repository.UpdateEmployees(id, employee)
	if err != nil {
		log.Print("error updating employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// DeleteEmployee handles DELETE requests to remove a specific employee by ID.
func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract employee ID from URL
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

// DeleteAllEmployees handles DELETE requests to remove all employees from the database.
func DeleteAllEmployees(w http.ResponseWriter, r *http.Request) {
	err := repository.DeleteAllEmployees()
	if err != nil {
		log.Print("error deleting all employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetEmployeeByID handles GET requests to retrieve a specific employee by their ID.
func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	// Extract employee ID from URL
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

// GetPaginatedEmployees handles GET requests with pagination parameters.
// It accepts page number and limit as query parameters.
func GetPaginatedEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query()

	// Extract and validate pagination parameters
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

	// Fetch paginated results
	employees, err := repository.GetPaginatedEmployees(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employees)
}
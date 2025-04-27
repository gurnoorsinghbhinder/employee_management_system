package controller

import (
	"encoding/json"
	"employee/models"
	"employee/repository"
	"employee/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmployeeController struct {
	Repo repository.EmployeeRepository
}

func NewEmployeeController(repo repository.EmployeeRepository) *EmployeeController {
	return &EmployeeController{Repo: repo}
}

// CreateEmployee handles POST requests to create a new employee.
func (ec *EmployeeController) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var employee models.Employee
	json.NewDecoder(r.Body).Decode(&employee)

	var validationError []string

	// Validation checks
	if employee.Name == nil || strings.TrimSpace(*employee.Name) == "" {
		validationError = append(validationError, "Name is required")
	}
	if employee.Email == nil || strings.TrimSpace(*employee.Email) == "" {
		validationError = append(validationError, "Email is required")
	}
	if employee.Position == nil {
		validationError = append(validationError, "Position is required")
	}
	if employee.Salary == nil || *employee.Salary < 0 {
		validationError = append(validationError, "Salary must be a positive number")
	}
	if len(validationError) > 0 {
		http.Error(w, strings.Join(validationError, ","), http.StatusBadRequest)
		return
	}

	// Email validation
	if !utils.IsValidEmail(*employee.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Email already taken check
	taken, err := ec.Repo.IsEmailTaken(*employee.Email)
	if err != nil {
		http.Error(w, "Error while checking Email", http.StatusInternalServerError)
		return
	}
	if taken {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	// Create Employee
	if err := ec.Repo.CreateEmployee(employee); err != nil {
		log.Print("error creating employee")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// GetAllEmployee retrieves all employees.
func (ec *EmployeeController) GetAllEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	employees, err := ec.Repo.GetAllEmployees()
	if err != nil {
		log.Print("error getting all employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employees)
}

// UpdateEmployee updates an existing employee by ID.
func (ec *EmployeeController) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idParam := mux.Vars(r)["id"]
	id, _ := primitive.ObjectIDFromHex(idParam)

	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		http.Error(w, "invalid inputs", http.StatusBadRequest)
		return
	}

	// Validation checks
	if employee.Name == nil || strings.TrimSpace(*employee.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if employee.Email == nil || strings.TrimSpace(*employee.Email) == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	if !utils.IsValidEmail(*employee.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Prepare update fields
	update := bson.M{}
	if employee.Name != nil && strings.TrimSpace(*employee.Name) != "" {
		update["name"] = *employee.Name
	}
	if employee.Email != nil && strings.TrimSpace(*employee.Email) != "" {
		update["email"] = *employee.Email
	}
	if employee.Position != nil && strings.TrimSpace(*employee.Position) != "" {
		update["position"] = *employee.Position
	}
	if employee.Salary != nil {
		update["salary"] = *employee.Salary
	}
	if employee.Joining != nil {
		update["joining"] = *employee.Joining
	}
	if employee.PreviousOrgs != nil {
		update["previousOrgs"] = *employee.PreviousOrgs
	}

	// Check if email is taken by someone else
	taken, err := ec.Repo.IsEmailTakenByOther(id, *employee.Email)
	if err != nil {
		http.Error(w, "Error while checking Email", http.StatusInternalServerError)
		return
	}
	if taken {
		http.Error(w, "Email already in use", http.StatusBadRequest)
		return
	}

	// Update employee
	if err := ec.Repo.UpdateEmployees(id, update); err != nil {
		log.Print("error updating employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// DeleteEmployee removes an employee by ID.
func (ec *EmployeeController) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, _ := primitive.ObjectIDFromHex(idParam)

	if err := ec.Repo.DeleteEmployee(id); err != nil {
		log.Print("error deleting employee")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteAllEmployees removes all employees from the database.
func (ec *EmployeeController) DeleteAllEmployees(w http.ResponseWriter, r *http.Request) {
	if err := ec.Repo.DeleteAllEmployees(); err != nil {
		log.Print("error deleting all employees")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetEmployeeByID retrieves a specific employee by ID.
func (ec *EmployeeController) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	employee, err := ec.Repo.GetEmployeeByID(id)
	if err != nil {
		log.Print("error getting employee by ID")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

// GetPaginatedEmployees retrieves employees with pagination.
func (ec *EmployeeController) GetPaginatedEmployees(w http.ResponseWriter, r *http.Request) {
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

	employees, err := ec.Repo.GetPaginatedEmployees(page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employees)
}

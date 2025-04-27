package controller

import (
	"bytes"
	"encoding/json"
	"employee/models"
	"employee/repository"
	
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
)

type mockEmployeeRepo struct {
	repository.EmployeeRepository
	MockCreateEmployee     func(models.Employee) error
	MockGetAllEmployees    func() ([]models.Employee, error)
	MockGetEmployeeByID    func(primitive.ObjectID) (models.Employee, error)
	MockUpdateEmployees    func(primitive.ObjectID, bson.M) error // <-- Update signature here
	MockDeleteEmployee     func(primitive.ObjectID) error
	MockDeleteAllEmployees func() error
	MockIsEmailTaken       func(string) (bool, error)
	MockIsEmailTakenByOther func(primitive.ObjectID, string) (bool, error)
	MockGetPaginatedEmployees func(int, int) ([]models.Employee, error)
}

func (m *mockEmployeeRepo) CreateEmployee(emp models.Employee) error {
	return m.MockCreateEmployee(emp)
}
func (m *mockEmployeeRepo) GetAllEmployees() ([]models.Employee, error) {
	return m.MockGetAllEmployees()
}
func (m *mockEmployeeRepo) GetEmployeeByID(id primitive.ObjectID) (models.Employee, error) {
	return m.MockGetEmployeeByID(id)
}
func (m *mockEmployeeRepo) UpdateEmployees(id primitive.ObjectID, update bson.M) error { // <-- Here
	return m.MockUpdateEmployees(id, update)
}
func (m *mockEmployeeRepo) DeleteEmployee(id primitive.ObjectID) error {
	return m.MockDeleteEmployee(id)
}
func (m *mockEmployeeRepo) DeleteAllEmployees() error {
	return m.MockDeleteAllEmployees()
}
func (m *mockEmployeeRepo) IsEmailTaken(email string) (bool, error) {
	return m.MockIsEmailTaken(email)
}
func (m *mockEmployeeRepo) IsEmailTakenByOther(id primitive.ObjectID, email string) (bool, error) {
	return m.MockIsEmailTakenByOther(id, email)
}
func (m *mockEmployeeRepo) GetPaginatedEmployees(page int, limit int) ([]models.Employee, error) {
	return m.MockGetPaginatedEmployees(page, limit)
}

// --- Test helpers ---
func stringPointer(s string) *string {
	return &s
}
func floatPointer(f float64) *float64 {
	return &f
}

// ------------------------- Tests start -----------------------------

func TestCreateEmployee(t *testing.T) {
	mockRepo := &mockEmployeeRepo{
		MockIsEmailTaken: func(email string) (bool, error) { return false, nil },
		MockCreateEmployee: func(emp models.Employee) error { return nil },
	}
	controller := NewEmployeeController(mockRepo)

	payload := models.Employee{
		Name:     stringPointer("John Doe"),
		Email:    stringPointer("john@example.com"),
		Position: stringPointer("Developer"),
		Salary:   floatPointer(50000),
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/employees", bytes.NewReader(body))
	w := httptest.NewRecorder()

	controller.CreateEmployee(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllEmployee(t *testing.T) {
	mockRepo := &mockEmployeeRepo{
		MockGetAllEmployees: func() ([]models.Employee, error) {
			return []models.Employee{
				{Name: stringPointer("Alice")},
			}, nil
		},
	}
	controller := NewEmployeeController(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/employees", nil)
	w := httptest.NewRecorder()

	controller.GetAllEmployee(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetEmployeeByID(t *testing.T) {
	mockID := primitive.NewObjectID()
	mockRepo := &mockEmployeeRepo{
		MockGetEmployeeByID: func(id primitive.ObjectID) (models.Employee, error) {
			return models.Employee{Name: stringPointer("Bob")}, nil
		},
	}
	controller := NewEmployeeController(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/employees/"+mockID.Hex(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": mockID.Hex()})
	w := httptest.NewRecorder()

	controller.GetEmployeeByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateEmployee(t *testing.T) {
	mockID := primitive.NewObjectID()
	mockRepo := &mockEmployeeRepo{
		MockUpdateEmployees: func(id primitive.ObjectID, update bson.M) error { return nil },
		MockIsEmailTakenByOther: func(id primitive.ObjectID, email string) (bool, error) { return false, nil },
	}
	controller := NewEmployeeController(mockRepo)

	payload := models.Employee{
		Name:  stringPointer("Charlie"),
		Email: stringPointer("charlie@example.com"),
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/employees/"+mockID.Hex(), bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": mockID.Hex()})
	w := httptest.NewRecorder()

	controller.UpdateEmployee(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteEmployee(t *testing.T) {
	mockID := primitive.NewObjectID()
	mockRepo := &mockEmployeeRepo{
		MockDeleteEmployee: func(id primitive.ObjectID) error { return nil },
	}
	controller := NewEmployeeController(mockRepo)

	req := httptest.NewRequest(http.MethodDelete, "/employees/"+mockID.Hex(), nil)
	req = mux.SetURLVars(req, map[string]string{"id": mockID.Hex()})
	w := httptest.NewRecorder()

	controller.DeleteEmployee(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteAllEmployees(t *testing.T) {
	mockRepo := &mockEmployeeRepo{
		MockDeleteAllEmployees: func() error { return nil },
	}
	controller := NewEmployeeController(mockRepo)

	req := httptest.NewRequest(http.MethodDelete, "/employees", nil)
	w := httptest.NewRecorder()

	controller.DeleteAllEmployees(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestGetPaginatedEmployees(t *testing.T) {
	mockRepo := &mockEmployeeRepo{
		MockGetPaginatedEmployees: func(page int, limit int) ([]models.Employee, error) {
			return []models.Employee{
				{Name: stringPointer("David")},
			}, nil
		},
	}
	controller := NewEmployeeController(mockRepo)

	req := httptest.NewRequest(http.MethodGet, "/employees?page=1&limit=2", nil)
	w := httptest.NewRecorder()

	controller.GetPaginatedEmployees(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
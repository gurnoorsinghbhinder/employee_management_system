package mocks

import (
	"employee/models"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmployeeRepositoryMock is a mock of the EmployeeRepository interface
type EmployeeRepositoryMock struct {
	mock.Mock
}

func (m *EmployeeRepositoryMock) CreateEmployeesBulk(employees []models.Employee) error {
    args := m.Called(employees)
    return args.Error(0)
}

func (m *EmployeeRepositoryMock) SearchEmployees(query string) ([]models.Employee, error) {
    args := m.Called(query)
    return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *EmployeeRepositoryMock) CreateEmployee(emp models.Employee) error {
	args := m.Called(emp)
	return args.Error(0)
}

func (m *EmployeeRepositoryMock) GetAllEmployees() ([]models.Employee, error) {
	args := m.Called()
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *EmployeeRepositoryMock) GetEmployeeByID(id primitive.ObjectID) (models.Employee, error) {
	args := m.Called(id)
	return args.Get(0).(models.Employee), args.Error(1)
}

func (m *EmployeeRepositoryMock) UpdateEmployees(id primitive.ObjectID, update bson.M) error {
	args := m.Called(id, update) 
	return args.Error(0)
}

func (m *EmployeeRepositoryMock) DeleteEmployee(id primitive.ObjectID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *EmployeeRepositoryMock) DeleteAllEmployees() error {
	args := m.Called()
	return args.Error(0)
}

func (m *EmployeeRepositoryMock) IsEmailTaken(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *EmployeeRepositoryMock) IsEmailTakenByOther(id primitive.ObjectID, email string) (bool, error) {
	args := m.Called(id, email)
	return args.Bool(0), args.Error(1)
}

func (m *EmployeeRepositoryMock) GetPaginatedEmployees(page int, limit int) ([]models.Employee, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Employee), args.Error(1)
}
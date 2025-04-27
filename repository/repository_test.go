package repository

import (
	"testing"
	"employee/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/stretchr/testify/suite"
	"employee/repository/mocks"
)

type EmployeeRepositoryTestSuite struct {
	suite.Suite
	mockRepo *mocks.EmployeeRepositoryMock
	repo     EmployeeRepository
}

func stringPointer(s string) *string {
	return &s 
}

func (suite *EmployeeRepositoryTestSuite) SetupTest() {
	suite.mockRepo = new(mocks.EmployeeRepositoryMock)
	suite.repo = suite.mockRepo
}

func (suite *EmployeeRepositoryTestSuite) TestCreateEmployee() {
	objectID := primitive.NewObjectID()
	emp := models.Employee{
		ID:    &objectID,  
		Email: stringPointer("test@company.com"),  
		Name:  stringPointer("John Doe"),
	}

	suite.mockRepo.On("CreateEmployee", emp).Return(nil)
	err := suite.repo.CreateEmployee(emp)

	suite.NoError(err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *EmployeeRepositoryTestSuite) TestGetAllEmployees() {
	objectID := primitive.NewObjectID()
	emp := models.Employee{
		ID:    &objectID,  
		Email: stringPointer("test@company.com"),  
		Name:  stringPointer("John Doe"),
	}

	suite.mockRepo.On("GetAllEmployees").Return([]models.Employee{emp}, nil)
	employees, err := suite.repo.GetAllEmployees()

	suite.NoError(err)
	suite.Equal(1, len(employees))
	suite.Equal("test@company.com", *employees[0].Email)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *EmployeeRepositoryTestSuite) TestGetEmployeeByID() {
	objectID := primitive.NewObjectID()
	emp := models.Employee{
		ID:    &objectID,  
		Email: stringPointer("test@company.com"),  
		Name:  stringPointer("John Doe"),
	}

	suite.mockRepo.On("GetEmployeeByID", objectID).Return(emp, nil)
	employee, err := suite.repo.GetEmployeeByID(objectID)

	suite.NoError(err)
	suite.Equal(objectID, *employee.ID)
	suite.Equal("test@company.com", *employee.Email)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *EmployeeRepositoryTestSuite) TestIsEmailTaken() {
	email := "test@company.com"

	suite.mockRepo.On("IsEmailTaken", email).Return(true, nil)
	taken, err := suite.repo.IsEmailTaken(email)

	suite.NoError(err)
	suite.True(taken)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestEmployeeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EmployeeRepositoryTestSuite))
}

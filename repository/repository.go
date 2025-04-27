package repository

import (
	"employee/db"
	"employee/models"
	"employee/utils"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type EmployeeRepository interface {
	CreateEmployee(employee models.Employee) error
	IsEmailTaken(email string) (bool, error)
	IsEmailTakenByOther(id primitive.ObjectID, email string) (bool, error)
	GetAllEmployees() ([]models.Employee, error)
	UpdateEmployees(id primitive.ObjectID, update bson.M) error
	DeleteEmployee(id primitive.ObjectID) error
	DeleteAllEmployees() error
	GetEmployeeByID(id primitive.ObjectID) (models.Employee, error)
	GetPaginatedEmployees(page int, limit int) ([]models.Employee, error)
}


type employeeRepository struct {}

func NewEmployeeRepository() EmployeeRepository {
	return &employeeRepository{}
}

// Now attach all methods to *employeeRepository

func (r *employeeRepository) CreateEmployee(emp models.Employee) error {
	ctx, cancel := utils.GetContext()
	defer cancel()

	_, err := db.EmployeeCollection.InsertOne(ctx, emp)
	return err
}

func (r *employeeRepository) GetAllEmployees() ([]models.Employee, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	var allEmployees []models.Employee
	cursor, err := db.EmployeeCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var employee models.Employee
		if err := cursor.Decode(&employee); err != nil {
			return nil, err
		}
		allEmployees = append(allEmployees, employee)
	}
	return allEmployees, nil
}

func (r *employeeRepository) UpdateEmployees(id primitive.ObjectID, update bson.M) error {
	ctx, cancel := utils.GetContext()
	defer cancel()

	if len(update) == 0 {
		return errors.New("nothing to update")
	}

	_, err := db.EmployeeCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *employeeRepository) DeleteEmployee(id primitive.ObjectID) error {
	ctx, cancel := utils.GetContext()
	defer cancel()

	_, err := db.EmployeeCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *employeeRepository) DeleteAllEmployees() error {
	ctx, cancel := utils.GetContext()
	defer cancel()

	_, err := db.EmployeeCollection.DeleteMany(ctx, bson.M{})
	return err
}

func (r *employeeRepository) GetEmployeeByID(id primitive.ObjectID) (models.Employee, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	var employee models.Employee
	err := db.EmployeeCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&employee)
	return employee, err
}

func (r *employeeRepository) GetPaginatedEmployees(page, limit int) ([]models.Employee, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := db.EmployeeCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	var employees []models.Employee
	if err := cursor.All(ctx, &employees); err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepository) IsEmailTaken(email string) (bool, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	count, err := db.EmployeeCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *employeeRepository) IsEmailTakenByOther(id primitive.ObjectID, email string) (bool, error) {
	ctx, cancel := utils.GetContext()
	defer cancel()

	filter := bson.M{
		"email": email,
		"_id":   bson.M{"$ne": id},
	}

	count, err := db.EmployeeCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

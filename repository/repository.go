package repository

import (
	"employee/db"
	"employee/models"
	"employee/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateEmployee adds a new employee to the database
// Parameters:
//   - emp: The employee model to be inserted
// Returns:
//   - error: Any error encountered during insertion
func CreateEmployee(emp models.Employee) error{
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.InsertOne(ctx,emp)
	return err
}

// GetAllEmployees retrieves all employees from the database
// Returns:
//   - []models.Employee: Slice containing all employees
//   - error: Any error encountered during retrieval
func GetAllEmployees()([]models.Employee,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()
	
	var allEmployees []models.Employee
	cursor,err:=db.EmployeeCollection.Find(ctx,bson.M{})
	if err!=nil{
		return nil,err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx){
		var employee models.Employee
		if err:=cursor.Decode(&employee);err!=nil{
			return nil,err
		}
		allEmployees=append(allEmployees, employee)
	}
	return allEmployees,nil
	
}

// UpdateEmployees updates an employee's information in the database
// Parameters:
//   - id: ObjectID of the employee to update
//   - emp: Updated employee data
// Returns:
//   - error: Any error encountered during update
func UpdateEmployees(id primitive.ObjectID,emp models.Employee) error {
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.UpdateOne(ctx,bson.M{"_id":id},bson.M{"$set":emp})
	return err
}

// DeleteEmployee removes a specific employee from the database
// Parameters:
//   - id: ObjectID of the employee to delete
// Returns:
//   - error: Any error encountered during deletion
func DeleteEmployee(id primitive.ObjectID) error {
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.DeleteOne(ctx,bson.M{"_id":id})
	return err
}

// DeleteAllEmployees removes all employees from the database
// Returns:
//   - error: Any error encountered during deletion
func DeleteAllEmployees()error{
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.DeleteMany(ctx,bson.M{})
	return err
}

// GetEmployeeByID retrieves a specific employee by their ID
// Parameters:
//   - id: ObjectID of the employee to retrieve
// Returns:
//   - models.Employee: The found employee
//   - error: Any error encountered during retrieval
func GetEmployeeByID(id primitive.ObjectID)(models.Employee,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()

	var employee models.Employee
	err:=db.EmployeeCollection.FindOne(ctx,bson.M{"_id":id}).Decode(&employee)
	return employee,err
}

// GetPaginatedEmployees retrieves employees with pagination
// Parameters:
//   - page: Current page number (1-based indexing)
//   - limit: Number of records per page
// Returns:
//   - []models.Employee: Slice of employees for the requested page
//   - error: Any error encountered during retrieval
func GetPaginatedEmployees(page,limit int)([]models.Employee,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()

	skip:=(page-1)*limit
	opts:=options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor,err:=db.EmployeeCollection.Find(ctx,bson.M{},opts)
	if err!=nil{
		return nil,err
	}
	
	var employees []models.Employee
	if err:=cursor.All(ctx,&employees);err!=nil{
		return nil,err
	}
	return employees,nil
}

// IsEmailTaken checks if an email is already registered in the database
// Parameters:
//   - email: Email address to check
// Returns:
//   - bool: true if email exists, false otherwise
//   - error: Any error encountered during the check
func IsEmailTaken(email string)(bool,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()

	count,err:=db.EmployeeCollection.CountDocuments(ctx,bson.M{"email":email})
	if err!=nil{
		return false,err
	}
	return count>0,nil
}

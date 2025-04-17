package repository

import (
	"employee/db"
	"employee/models"
	"employee/utils"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Creates a new employee in the database
func CreateEmployee(emp models.Employee) error{
	ctx,cancel:=utils.GetContext()
	defer cancel()

	if emp.Name==nil || emp.Email==nil{
		return errors.New("name and email are required")
	}

	_,err:=db.EmployeeCollection.InsertOne(ctx,emp)
	return err
}

// Retrieves all employees
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

// Updates an employee by ID
func UpdateEmployees(id primitive.ObjectID,emp models.Employee) error {
	ctx,cancel:=utils.GetContext()
	defer cancel()

	update:=bson.M{}
	
	if emp.Name!=nil && strings.TrimSpace(*emp.Name)!=""{
		update["name"]=*emp.Name
	}
	
	if emp.Email!=nil && strings.TrimSpace(*emp.Email)!=""{
		update["email"]=*emp.Email
	}

	if emp.Position!=nil && strings.TrimSpace(*emp.Position)!=""{
		update["position"]=*emp.Position
	}

	if emp.Salary!=nil {
		update["salary"]=*emp.Salary
	}

	if emp.Joining!=nil {
		update["joining"]= *emp.Joining
	}

	if emp.PreviousOrgs!=nil {
		update["previousOrgs"]=*emp.Joining
	}

	if len(update)==0{
		return errors.New("nothing to update")
	}

	_,err:=db.EmployeeCollection.UpdateOne(ctx,bson.M{"_id":id},bson.M{"$set":emp})
	return err
}

// Deletes an employee by ID
func DeleteEmployee(id primitive.ObjectID) error {
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.DeleteOne(ctx,bson.M{"_id":id})
	return err
}

// Deletes all employees
func DeleteAllEmployees()error{
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.DeleteMany(ctx,bson.M{})
	return err
}

// Gets employee by ID
func GetEmployeeByID(id primitive.ObjectID)(models.Employee,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()

	var employee models.Employee
	err:=db.EmployeeCollection.FindOne(ctx,bson.M{"_id":id}).Decode(&employee)
	return employee,err
}

// Retrieves employees with pagination
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

// Checks if email already exists
func IsEmailTaken(email string)(bool,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()

	count,err:=db.EmployeeCollection.CountDocuments(ctx,bson.M{"email":email})
	if err!=nil{
		return false,err
	}
	return count>0,nil
}

func IsEmailTakenByOther(id primitive.ObjectID,email string)(bool,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()

	filter:=bson.M{
		"email":email,
		"_id":bson.M{"$ne":id},
	}

	count,err:=db.EmployeeCollection.CountDocuments(ctx,filter)
	if err!=nil{
		return false,err
	}
	return count>0,nil
}

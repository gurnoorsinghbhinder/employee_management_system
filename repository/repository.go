package repository

import (
	"employee/db"
	"employee/models"
	"employee/utils"
	
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func CreateEmployee(emp models.Employee) error{
    ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.InsertOne(ctx,emp)
	return err
}

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

func UpdateEmployees(id primitive.ObjectID,emp models.Employee) error {
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.UpdateOne(ctx,bson.M{"_id":id},bson.M{"$set":emp})
    return err
}

func DeleteEmployee(id primitive.ObjectID) error {
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.DeleteOne(ctx,bson.M{"_id":id})
	return err
}

func DeleteAllEmployees()error{
	ctx,cancel:=utils.GetContext()
	defer cancel()

	_,err:=db.EmployeeCollection.DeleteMany(ctx,bson.M{})
	return err
}

func GetEmployeeByID(id primitive.ObjectID)(models.Employee,error){
	ctx,cancel:=utils.GetContext()
	defer cancel()

    var employee models.Employee
	err:=db.EmployeeCollection.FindOne(ctx,bson.M{"_id":id}).Decode(&employee)
	return employee,err
}
package utils

import (
	"context"
	"encoding/json"
	"net/mail"
	"time"
)

type CustomTime struct{
	time.Time
}

func (ct *CustomTime) UnmarshlJson(b []byte)error{
	var str string
	if err:=json.Unmarshal(b,&str);err!=nil{
		return err
	}
	if str==""{
		ct.Time=time.Time{}  //means creating a new value in struct 
		return nil
	}
	parsed,err:=time.Parse("2006-01-02",str)
	if err!=nil{
		return err
	}
	ct.Time=parsed
	return nil
	
}

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)

}

func IsValidEmail(email string)bool{

	_,err:=mail.ParseAddress(email)
	return err==nil
}
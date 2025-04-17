package utils

import (
	"context"
	"time"
	"net/mail"
)

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)

}

func IsValidEmail(email string)bool{

	_,err:=mail.ParseAddress(email)
	return err==nil
}
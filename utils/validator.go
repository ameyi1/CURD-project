package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
func Getvalimsg(err error) string {
	var errlist []string
	validationErr := err.(validator.ValidationErrors)
	for _, item := range validationErr {
		field := item.Field()
		tag := item.Tag()
		param := item.Param()

		switch tag {
		case "required":
			errlist = append(errlist, field+"不能为空")
		case "min":
			if item.Kind().String() == "string" {
				errlist = append(errlist, field+"长度不能小于"+param+"位")
			} else {
				errlist = append(errlist, field+"长度不能小于"+param)
			}
		case "max":
			if item.Kind().String() == "string" {
				errlist = append(errlist, field+"长度不能大于"+param+"位")
			} else {
				errlist = append(errlist, field+"长度不能大于"+param)
			}
		default:
			errlist = append(errlist, field+"不合法")
		}
	}
	return strings.Join(errlist, ";")
}

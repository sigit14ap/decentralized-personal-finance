package helpers

import (
	"fmt"

	"github.com/go-playground/validator"
	pb "github.com/sigit14ap/personal-finance/auth-service/internal/proto"
)

func ValidateRequest(request interface{}) []*pb.ResponseError {
	var response []*pb.ResponseError
	validate := validator.New()

	err := validate.Struct(request)
	if err == nil {
		return nil
	}

	for _, err := range err.(validator.ValidationErrors) {
		var message error
		switch err.Tag() {
		case "required":
			message = fmt.Errorf("field '%s' cannot be blank", err.Field())
		case "email":
			message = fmt.Errorf("field '%s' must be a valid email address", err.Field())
		case "len":
			message = fmt.Errorf("field '%s' must be exactly %v characters long", err.Field(), err.Param())
		default:
			message = fmt.Errorf("field '%s': '%v' must satisfy '%s' '%v' criteria", err.Field(), err.Value(), err.Tag(), err.Param())
		}

		errorField := &pb.ResponseError{
			Field:   err.Field(),
			Message: message.Error(),
		}

		response = append(response, errorField)
	}

	return response
}

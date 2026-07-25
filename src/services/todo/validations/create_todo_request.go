package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/furkancmn57/go-architecture-template/src/common/apperr"
	"github.com/furkancmn57/go-architecture-template/src/models/requests"
)

// CreateTodoRequest validates a create todo payload.
func CreateTodoRequest(req requests.CreateTodoRequest) *apperr.Error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.Description, validation.Length(0, 2000)),
	)
	return apperr.FromValidation(err)
}

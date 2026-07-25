package validations

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/furkancmn57/go-architecture-template/src/common/apperr"
	"github.com/furkancmn57/go-architecture-template/src/models/requests"
)

// UpdateTodoRequest validates an update todo payload.
func UpdateTodoRequest(req requests.UpdateTodoRequest) *apperr.Error {
	err := validation.ValidateStruct(&req,
		validation.Field(&req.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&req.Description, validation.Length(0, 2000)),
	)
	return apperr.FromValidation(err)
}

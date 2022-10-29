package validators

import (
	"pirs.io/process/validation/models"
)

type SchemaValidator struct {
	next models.Validator
}

func (sv *SchemaValidator) Validate(data models.ImportProcessValidationData) {
	// todo
	panic("not implemented")
}

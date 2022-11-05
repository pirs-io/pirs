package models

type Validator interface {
	Validate(*ImportProcessValidationData)
	SetNext(Validator)
	ExecuteNextIfPresent(*ImportProcessValidationData)
}

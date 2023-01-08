package models

type Validator interface {
	Validate(Validable)
	SetNext(Validator)
	ExecuteNextIfPresent(Validable)
}

type Validable interface {
	IsValidable()
}

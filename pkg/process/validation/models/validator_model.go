package models

type Validator interface {
	Validate(interface{})
	SetNext(Validator)
	ExecuteNextIfPresent(interface{})
}

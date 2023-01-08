package validators

import (
	"pirs.io/process/validation/models"
)

// An SchemaValidator contains next validator instance. It's implementation of models.Validator
type SchemaValidator struct {
	next models.Validator
}

// Validate takes models.ImportValidationData and validates it in schema context. If data is valid, it sets
// field IsSchemaValid of models.ValidationFlags to true. Otherwise it sets to false.
func (sv *SchemaValidator) Validate(data models.Validable) {
	typedData := data.(*models.ImportValidationData)
	var isValid bool
	defer sv.ExecuteNextIfPresent(typedData)
	defer func() { typedData.ValidationFlags.IsSchemaValid = isValid }()

	// todo
	// neexistuje standard libka, len github projekty. Vacsina to su nieco ako skolske projekty, ktore nie su podporovane
	// neviem ci take nieco chceme do projektu
	// napriklad https://github.com/krolaw/xsd, 4 rocna libka, nadstavba C kodu, obsahuje todos, pravdepodobne nepodporovana

	isValid = true
}

func (sv *SchemaValidator) ExecuteNextIfPresent(data models.Validable) {
	if sv.next != nil {
		sv.next.Validate(data)
	}
}

func (sv *SchemaValidator) SetNext(validator models.Validator) {
	sv.next = validator
}

func isSchemaValid() {
}

func findXsdInXml() {
}

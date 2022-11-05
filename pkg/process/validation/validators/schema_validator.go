package validators

import (
	"pirs.io/process/validation/models"
)

type SchemaValidator struct {
	next models.Validator
}

func (sv *SchemaValidator) Validate(data *models.ImportProcessValidationData) {
	var isValid bool
	defer sv.ExecuteNextIfPresent(data)
	defer func() { data.ValidationFlags.IsSchemaValid = isValid }()

	// todo
	// neexistuje standard libka, len github projekty. Vacsina to su nieco ako skolske projekty, ktore nie su podporovane
	// neviem ci take nieco chceme do projektu
	// napriklad https://github.com/krolaw/xsd, 4 rocna libka, nadstavba C kodu, obsahuje todos, pravdepodobne nepodporovana

	isValid = true
}

func (sv *SchemaValidator) ExecuteNextIfPresent(data *models.ImportProcessValidationData) {
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

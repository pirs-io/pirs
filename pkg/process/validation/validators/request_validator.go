package validators

import (
	"pirs.io/process/validation/models"
	"regexp"
	"strings"
)

const (
	FILENAME_REGEX   = "^[\\w\\-. ]+$"
	PARTIALURI_REGEX = "^[\\w]+$"
)

// An ImportProcessRequestValidator contains next validator instance. It's implementation of models.Validator
type ImportProcessRequestValidator struct {
	next models.Validator
}

type ImportPackageRequestValidator struct {
	next models.Validator
}

// Validate takes models.ImportProcessValidationData and validates it in request context. If data is valid, it sets
// field IsRequestValid of models.ValidationFlags to true. Otherwise it sets to false.
func (rv *ImportProcessRequestValidator) Validate(data *models.ImportProcessValidationData) {
	var isValid bool
	defer rv.ExecuteNextIfPresent(data)
	defer func() { data.ValidationFlags.IsRequestValid = isValid }()

	if isValidDataLength(data.ReqData.ProcessData.Len(), data.ReqData.ProcessSize) &&
		isValidFileName(data.ReqData.ProcessFileName) &&
		isValidPartialUri(data.ReqData.PartialUri) {
		isValid = true
	} else {
		isValid = false
	}
}

func (rv *ImportProcessRequestValidator) ExecuteNextIfPresent(data *models.ImportProcessValidationData) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}

func (rv *ImportProcessRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

func (rv *ImportPackageRequestValidator) Validate(data *models.ImportPackageValidationData) {
	// todo
	panic("not implemented")
}

func (rv *ImportPackageRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

func isValidFileName(fileName string) bool {
	isValid, err := regexp.MatchString(FILENAME_REGEX, fileName)
	if err != nil {
		return false
	}
	return isValid
}
func isValidDataLength(dataLen int, dataLenFromReq int) bool {
	return dataLen != 0 && dataLen == dataLenFromReq
}
func isValidPartialUri(partialUri string) bool {
	if strings.Contains(partialUri, ".") {
		splitPartialUri := strings.Split(partialUri, ".")
		if len(splitPartialUri) == 3 {
			for _, part := range splitPartialUri {
				if len(part) == 0 {
					return false
				}
				isValid, err := regexp.MatchString(PARTIALURI_REGEX, part)
				if err != nil || !isValid {
					return false
				}
			}
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

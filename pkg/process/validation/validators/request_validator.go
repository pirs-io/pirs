package validators

import (
	"pirs.io/process/validation/models"
	"regexp"
)

const (
	FILENAME_REGEX = "^[\\w\\-. ]+$"
)

type ImportProcessRequestValidator struct {
	next models.Validator
}

type ImportPackageRequestValidator struct {
	next models.Validator
}

func (rv *ImportProcessRequestValidator) Validate(data *models.ImportProcessValidationData) {
	var isValid bool
	defer rv.ExecuteNextIfPresent(data)
	defer func() { data.ValidationFlags.IsRequestValid = isValid }()

	if !isValidDataLength(data.ReqData.ProcessData.Len(), data.ReqData.ProcessSize) {
	} else if !isValidFileName(data.ReqData.ProcessFileName) {
	} else {
		isValid = true
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

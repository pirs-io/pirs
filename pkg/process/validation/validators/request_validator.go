package validators

import (
	"pirs.io/process/validation/models"
	"regexp"
	"sync"
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

func (rq *ImportProcessRequestValidator) Validate(data *models.ImportProcessValidationData) {
	var wg sync.WaitGroup
	wg.Add(2)
	data.ValidationFlags.IsRequestValid = true

	go func(dataLen int, dataLenFromReq int) {
		defer wg.Done()
		if dataLen == 0 {
			data.ValidationFlags.IsRequestValid = false
		} else if dataLen != dataLenFromReq {
			data.ValidationFlags.IsRequestValid = false
		}
	}(data.ReqData.ProcessData.Len(), data.ReqData.ProcessSize)

	go func(fileName string) {
		defer wg.Done()
		if !isValidFileName(fileName) {
			data.ValidationFlags.IsRequestValid = false
		}
	}(data.ReqData.ProcessFileName)

	wg.Wait()
	if rq.next != nil {
		rq.next.Validate(data)
	}
}

func (rq *ImportProcessRequestValidator) SetNext(validator models.Validator) {
	rq.next = validator
}

func (rq *ImportPackageRequestValidator) Validate(data *models.ImportPackageValidationData) {
	panic("not implemented")
}

func (rq *ImportPackageRequestValidator) SetNext(validator models.Validator) {
	rq.next = validator
}

func isValidFileName(fileName string) bool {
	isValid, err := regexp.MatchString(FILENAME_REGEX, fileName)
	if err != nil {
		return false
	}
	return isValid
}

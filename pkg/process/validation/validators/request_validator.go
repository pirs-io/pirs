package validators

import (
	"pirs.io/process/validation/models"
	"regexp"
	"strings"
)

const (
	FILENAME_REGEX = "^[\\w\\-. ]+$"
	URI_REGEX      = "^[\\w]+$"
	VERSION_REGEX  = "^[\\d]+$"
)

// An ImportProcessRequestValidator contains next validator instance. It's implementation of models.Validator
type ImportProcessRequestValidator struct {
	next models.Validator
}

type ImportPackageRequestValidator struct {
	next models.Validator
}

type DownloadProcessRequestValidator struct {
	next models.Validator
}

type DownloadPackageRequestValidator struct {
	next models.Validator
}

// Validate takes models.ImportProcessValidationData and validates it in request context. If data is valid, it sets
// field IsRequestValid of models.ValidationFlags to true. Otherwise it sets to false.
func (rv *ImportProcessRequestValidator) Validate(data interface{}) {
	typedData := data.(*models.ImportProcessValidationData)
	var isValid bool
	defer rv.ExecuteNextIfPresent(typedData)
	defer func() { typedData.ValidationFlags.IsRequestValid = isValid }()

	if isValidDataLength(typedData.ReqData.ProcessData.Len(), typedData.ReqData.ProcessSize) &&
		isValidFileName(typedData.ReqData.ProcessFileName) &&
		isValidPartialUri(typedData.ReqData.PartialUri) {
		isValid = true
	} else {
		isValid = false
	}
}

func (rv *ImportPackageRequestValidator) Validate(data *models.ImportPackageValidationData) {
	// todo
	panic("not implemented")
}

// Validate todo
func (rv *DownloadProcessRequestValidator) Validate(data interface{}) {
	typedData := data.(*models.DownloadValidationData)
	var isValid bool
	defer rv.ExecuteNextIfPresent(typedData)
	defer func() { typedData.ValidationFlags.IsRequestValid = isValid }()

	if isValidUri(typedData.ReqData.TargetUri) {
		isValid = true
	} else {
		isValid = false
	}
}

func (rv *DownloadPackageRequestValidator) Validate(data interface{}) {
	typedData := data.(*models.DownloadValidationData)
	var isValid bool
	defer rv.ExecuteNextIfPresent(typedData)
	defer func() { typedData.ValidationFlags.IsRequestValid = isValid }()

	if isValidPackageUri(typedData.ReqData.TargetUri) {
		isValid = true
	} else {
		isValid = false
	}
}

func (rv *ImportProcessRequestValidator) ExecuteNextIfPresent(data interface{}) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}

func (rv *DownloadProcessRequestValidator) ExecuteNextIfPresent(data interface{}) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}

func (rv *DownloadPackageRequestValidator) ExecuteNextIfPresent(data interface{}) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}

func (rv *ImportProcessRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

func (rv *ImportPackageRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

func (rv *DownloadProcessRequestValidator) SetNext(validator models.Validator) {
	rv.next = validator
}

func (rv *DownloadPackageRequestValidator) SetNext(validator models.Validator) {
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
		splitUri := strings.Split(partialUri, ".")
		if len(splitUri) == 3 {
			for _, part := range splitUri {
				if len(part) == 0 {
					return false
				}
				isValid, err := regexp.MatchString(URI_REGEX, part)
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
func isValidUri(uri string) bool {
	if strings.Contains(uri, ".") && strings.Contains(uri, ":") {
		splitUri := strings.Split(uri, ".")
		if len(splitUri) == 4 && strings.Contains(splitUri[len(splitUri)-1], ":") {
			splitLastElem := strings.Split(splitUri[len(splitUri)-1], ":")
			if len(splitLastElem) == 2 {
				splitUri[3] = splitLastElem[0]
				splitUri = append(splitUri, splitLastElem[1])
			} else {
				return false
			}

			for idx, part := range splitUri {
				if len(part) == 0 {
					return false
				}
				if idx == 4 {
					isValid, err := regexp.MatchString(VERSION_REGEX, part)
					if err != nil || !isValid {
						return false
					}
				} else {
					isValid, err := regexp.MatchString(URI_REGEX, part)
					if err != nil || !isValid {
						return false
					}
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
func isValidPackageUri(uri string) bool {
	if strings.Contains(uri, ".") {
		splitUri := strings.Split(uri, ".")
		if len(splitUri) == 3 {
			for _, part := range splitUri {
				if len(part) == 0 {
					return false
				}
				isValid, err := regexp.MatchString(URI_REGEX, part)
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

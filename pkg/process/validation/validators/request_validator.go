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

// An ImportRequestValidator contains next validator instance. It's implementation of models.Validator
type ImportRequestValidator struct {
	next models.Validator
}

type DownloadProcessRequestValidator struct {
	next models.Validator
}

type DownloadPackageRequestValidator struct {
	next models.Validator
}

// Validate takes models.ImportValidationData and validates it in request context. If data is valid, it sets
// field IsRequestValid of models.ImportProcessValidationFlags to true. Otherwise it sets to false.
func (rv *ImportRequestValidator) Validate(data models.Validable) {
	typedData := data.(*models.ImportValidationData)
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

// Validate takes models.DownloadValidationData and validates it in request context. If data is valid, it sets
// field IsRequestValid of models.DownloadValidationFlags to true. Otherwise it sets to false. It validates URI.
func (rv *DownloadProcessRequestValidator) Validate(data models.Validable) {
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

// Validate takes models.DownloadValidationData and validates it in request context. If data is valid, it sets
// field IsRequestValid of models.DownloadValidationFlags to true. Otherwise it sets to false. It validates package URI.
func (rv *DownloadPackageRequestValidator) Validate(data models.Validable) {
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

func (rv *ImportRequestValidator) ExecuteNextIfPresent(data models.Validable) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}

func (rv *DownloadProcessRequestValidator) ExecuteNextIfPresent(data models.Validable) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}

func (rv *DownloadPackageRequestValidator) ExecuteNextIfPresent(data models.Validable) {
	if rv.next != nil {
		rv.next.Validate(data)
	}
}

func (rv *ImportRequestValidator) SetNext(validator models.Validator) {
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

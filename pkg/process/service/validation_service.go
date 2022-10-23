package service

import (
	"errors"
	"strings"
)

type ValidationService struct {
	allowedExtensions []string
}

func NewValidationService(rawExtensions string) *ValidationService {
	parsedExtensions := strings.Split(rawExtensions, ",")
	for i, extension := range parsedExtensions {
		parsedExtensions[i] = strings.TrimSpace(extension)
	}
	return &ValidationService{allowedExtensions: parsedExtensions}
}

func (vs *ValidationService) ValidateProcessData(data ImportProcessRequestData) (bool, error) {
	var isValid bool
	// check if request is correctly initialized
	// validate extension
	extension, err := vs.parseExtensionFromFileName(data.ProcessFileName)
	if err != nil {
		// todo log err
		return false, err
	}
	isValid = vs.validateExtension(extension)
	if !isValid {
		return false, nil
	}
	// get schema
	// validate schema
	return true, nil
}

func (vs *ValidationService) parseExtensionFromFileName(filename string) (string, error) {
	if !strings.Contains(filename, ".") {
		return "", errors.New("given filename does not contain type extension")
	}
	splitFileName := strings.Split(filename, ".")
	return splitFileName[len(splitFileName)-1], nil
}

func (vs *ValidationService) validateExtension(extension string) bool {
	return contains(vs.allowedExtensions, extension)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

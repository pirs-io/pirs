package validators

import (
	"errors"
	"pirs.io/process/validation/models"
	"strings"
)

type ExtensionValidator struct {
	allowedExtensions []string
	next              models.Validator
}

func NewExtensionValidator(rawExtensions string) *ExtensionValidator {
	parsedExtensions := strings.Split(rawExtensions, ",")
	for i, extension := range parsedExtensions {
		parsedExtensions[i] = strings.TrimSpace(extension)
	}
	return &ExtensionValidator{allowedExtensions: parsedExtensions}
}

func (ev *ExtensionValidator) Validate(data *models.ImportProcessValidationData) {
	extension, err := ev.parseExtensionFromFileName(data.ReqData.ProcessFileName)
	if err != nil {
		data.ValidationFlags.IsExtensionValid = false
		// todo log err
	}
	isValid := ev.validateExtension(extension)
	if !isValid {
		data.ValidationFlags.IsExtensionValid = false
		// todo log/err non valid
	} else {
		data.ValidationFlags.IsExtensionValid = true
	}
	if ev.next != nil {
		ev.next.Validate(data)
	}
}

func (ev *ExtensionValidator) SetNext(validator models.Validator) {
	ev.next = validator
}

func (ev *ExtensionValidator) parseExtensionFromFileName(filename string) (string, error) {
	if !strings.Contains(filename, ".") {
		return "", errors.New("given filename does not contain type extension")
	}
	splitFileName := strings.Split(filename, ".")
	return splitFileName[len(splitFileName)-1], nil
}

func (ev *ExtensionValidator) validateExtension(extension string) bool {
	return contains(ev.allowedExtensions, extension)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

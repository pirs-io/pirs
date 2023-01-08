package validators

import (
	"bytes"
	"github.com/antchfx/xmlquery"
	"github.com/gabriel-vasile/mimetype"
	"pirs.io/commons"
	"pirs.io/process/enums"
	"pirs.io/process/metadata/determiner"
	"pirs.io/process/validation/models"
	"strings"
)

var (
	log = commons.GetLoggerFor("FileTypeValidator")
)

// A FileTypeValidator contains configurations to validate content type and next validator instance. It's implementation
// of models.Validator
type FileTypeValidator struct {
	allowedExtensions    []string
	ignoreWrongExtension bool
	next                 models.Validator
}

// NewFileTypeValidator creates FileTypeValidator instance. Parameter rawExtensions is initialized in config package and
// represents all the file type extensions, that can be accepted. If parameter ignoreWrongExtension is set to true, then
// method FileTypeValidator.Validate doesn't compare extension from the filename with the extension found out by system.
func NewFileTypeValidator(rawExtensions string, ignoreWrongExtension bool) *FileTypeValidator {
	parsedExtensions := strings.Split(rawExtensions, ",")
	for i, extension := range parsedExtensions {
		parsedExtensions[i] = strings.TrimSpace(extension)
	}
	return &FileTypeValidator{
		allowedExtensions:    parsedExtensions,
		ignoreWrongExtension: ignoreWrongExtension,
	}
}

// Validate takes models.ImportValidationData and validates it in file-type context. If data is valid, it sets
// field IsFileTypeValid of models.ValidationFlags to true. Otherwise, it sets to false.
func (ft *FileTypeValidator) Validate(data models.Validable) {
	typedData := data.(*models.ImportValidationData)
	var isValid bool
	defer ft.ExecuteNextIfPresent(typedData)
	defer func() { typedData.ValidationFlags.IsFileTypeValid = isValid }()

	fileType := mimetype.Detect(typedData.ReqData.ProcessData.Bytes())

	if ft.isAllowedType(fileType.Extension()) &&
		(ft.ignoreWrongExtension || fileType.Extension() == ft.parseExtensionFromFileName(typedData.ReqData.ProcessFileName)) {
		// here we can call xmlquery.Parse
		doc, err := xmlquery.Parse(bytes.NewReader(typedData.ReqData.ProcessData.Bytes()))
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		processType := determiner.DetermineProcessType(doc)
		if processType != enums.UNKNOWN {
			isValid = true
		} else {
			isValid = false
		}
	} else {
		isValid = false
	}
}

func (ft *FileTypeValidator) ExecuteNextIfPresent(data models.Validable) {
	if ft.next != nil {
		ft.next.Validate(data)
	}
}

func (ft *FileTypeValidator) SetNext(validator models.Validator) {
	ft.next = validator
}

func (ft *FileTypeValidator) parseExtensionFromFileName(filename string) string {
	if !strings.Contains(filename, ".") {
		return ""
	}
	splitFileName := strings.Split(filename, ".")
	return "." + splitFileName[len(splitFileName)-1]
}

func (ft *FileTypeValidator) isAllowedType(extension string) bool {
	return contains(ft.allowedExtensions, extension)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

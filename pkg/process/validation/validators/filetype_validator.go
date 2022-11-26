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

type FileTypeValidator struct {
	allowedExtensions    []string
	ignoreWrongExtension bool
	next                 models.Validator
}

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

func (ft *FileTypeValidator) Validate(data *models.ImportProcessValidationData) {
	var isValid bool
	defer ft.ExecuteNextIfPresent(data)
	defer func() { data.ValidationFlags.IsFileTypeValid = isValid }()

	fileType := mimetype.Detect(data.ReqData.ProcessData.Bytes())
	doc, err := xmlquery.Parse(bytes.NewReader(data.ReqData.ProcessData.Bytes()))
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	processType := determiner.DetermineProcessType(doc)
	if !ft.isAllowedType(fileType.Extension()) || processType == enums.UNKNOWN {
	} else if !ft.ignoreWrongExtension && fileType.Extension() != ft.parseExtensionFromFileName(data.ReqData.ProcessFileName) {
	} else {
		isValid = true
	}
}

func (ft *FileTypeValidator) ExecuteNextIfPresent(data *models.ImportProcessValidationData) {
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

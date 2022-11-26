package validators

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"io"
	"os"
	"pirs.io/commons"
	"pirs.io/process/service/models"
	valModels "pirs.io/process/validation/models"
	"testing"
)

var (
	log_test = commons.GetLoggerFor("FileTypeValidatorTest")
)

const (
	RESOURCES   = "../../resources/testing/"
	PF1         = "empty_petriflow.xml"
	PF2         = "v5_petriflow.xml"
	PF3         = "v6_petriflow.xml"
	BPMN1       = "empty.bpmn"
	WRONG1      = "microservices.drawio"
	WRONG2      = "Spresnenie_harmonogramu_zimneho_semestra_2022.pdf"
	PARTIAL_URI = "io.pirs.testing"
)

func TestFileTypeValidator_Validate(t *testing.T) {
	// todo use prod config
	conf, err := GetTestingConfig("../../dev.env", &TestingConfig{})
	if err != nil {
		log_test.Fatal().Msgf("Unable to load application config for FileTypeValidatorTest! %s", err)
	}
	validator := NewFileTypeValidator(conf.AllowedFileExtensions, conf.IgnoreWrongExtension)

	valData := buildValidationData(RESOURCES+PF3, PF3)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+PF2, PF2)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+BPMN1, BPMN1)
	validator.Validate(&valData)
	assert.Equal(t, true, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+WRONG1, WRONG1)
	validator.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsFileTypeValid)

	valData = buildValidationData(RESOURCES+WRONG2, WRONG2)
	validator.Validate(&valData)
	assert.Equal(t, false, valData.ValidationFlags.IsFileTypeValid)
}

func buildValidationData(path string, filename string) valModels.ImportProcessValidationData {
	file, err := os.Open(path)
	if err != nil {
		log_test.Fatal().Msg(err.Error())
	}
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024)
	var data []byte
	var totalSize = 0
	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log_test.Fatal().Msg(err.Error())
			}
			break
		}
		data = append(data, buf[:n]...)
		totalSize = totalSize + n
	}

	err = file.Close()
	if err != nil {
		log_test.Fatal().Msg(err.Error())
	}

	return valModels.ImportProcessValidationData{
		ReqData: models.ImportProcessRequestData{
			Ctx:             context.Background(),
			PartialUri:      PARTIAL_URI,
			ProcessFileName: filename,
			ProcessData:     *bytes.NewBuffer(data),
			ProcessSize:     totalSize,
		},
		ValidationFlags: valModels.ValidationFlags{},
	}
}

type TestingConfig struct {
	AllowedFileExtensions string `mapstructure:"ALLOWED_FILE_EXTENSIONS"`
	IgnoreWrongExtension  bool   `mapstructure:"IGNORE_WRONG_EXTENSION"`
}

func GetTestingConfig(configFilePath string, c *TestingConfig) (res *TestingConfig, err error) {
	viper.SetConfigType("env")
	viper.SetConfigFile(configFilePath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
		return nil, err
	}
	return c, nil
}

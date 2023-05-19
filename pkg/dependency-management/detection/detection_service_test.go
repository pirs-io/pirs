package detection

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"pirs.io/commons"
	"pirs.io/commons/enums"
	"pirs.io/dependency-management/detection/models"
	"pirs.io/dependency-management/mocks"
	"testing"
)

const (
	mockProcessData = "</dataRef>\n\t\t\t<dataRef>\n\t\t\t\t<id>mileage</id>\n\t\t\t\t<logic>\n\t\t\t\t\t<behavior>visible</behavior>\n\t\t\t\t</logic>\n\t\t\t\t<layout>\n\t\t\t\t\t<x>3</x>\n\t\t\t\t\t<y>1</y>\n\t\t\t\t\t<rows>1</rows>\n\t\t\t\t\t<cols>1</cols>\n\t\t\t\t\t<offset>0</offset>\n\t\t\t\t\t<template>material</template>\n\t\t\t\t\t<appearance>outline</appearance>\n\t\t\t\t</layout>\n\t\t\t</dataRef>\n\t\t\t<dataRef>\n\t\t\t\t<id>vin</id>\n\t\t\t\t<logic>\n\t\t\t\t\t<behavior>visible</behavior>\n\t\t\t\t</logic>\n\t\t\t\t<layout>\n\t\t\t\t\t<x>0</x>\n\t\t\t\t\t<y>2</y>\n\t\t\t\t\t<rows>1</rows>\n\t\t\t\t\t<cols>1</cols>\n\t\t\t\t\t<offset>0</offset>\n\t\t\t\t\t<template>material</template>\n\t\t\t\t\t<appearance>outline</appearance>\n\t\t\t\t</layout>\n\t\t\t</dataRef>"
)

func TestDetectionService_Detect(t *testing.T) {
	data := []byte(mockProcessData)
	rawHash := commons.HashBytesToSHA256(data)
	checksum := commons.ConvertBytesToString(rawHash)

	// case 1
	service := DetectionService{
		detectorChain: &mocks.Detector{
			SizeOfReturnedArray: 0,
		},
	}
	request := models.DetectRequestData{
		CheckSum:    checksum,
		ProcessType: enums.Petriflow,
		ProcessData: *bytes.NewBuffer(data),
	}
	response := service.Detect(request)
	assert.Equal(t, codes.NotFound, response.Status)
	assert.Equal(t, primitive.NilObjectID, response.Metadata[len(response.Metadata)-1].ID)

	// case 2
	service = DetectionService{
		detectorChain: &mocks.Detector{
			SizeOfReturnedArray: 2,
		},
	}
	response = service.Detect(request)
	assert.Equal(t, codes.OK, response.Status)
	assert.Equal(t, primitive.NilObjectID, response.Metadata[len(response.Metadata)-1].ID)

	// case 3
	service = DetectionService{
		detectorChain: &mocks.Detector{
			SizeOfReturnedArray: 2,
		},
	}
	request = models.DetectRequestData{
		CheckSum:    "awdawdawdawd",
		ProcessType: enums.Petriflow,
		ProcessData: *bytes.NewBuffer(data),
	}
	response = service.Detect(request)
	assert.Equal(t, codes.InvalidArgument, response.Status)

	// case 4
	service = DetectionService{
		detectorChain: &mocks.Detector{
			SizeOfReturnedArray: 2,
		},
	}
	request = models.DetectRequestData{
		CheckSum:    checksum,
		ProcessType: enums.Petriflow,
		ProcessData: *bytes.NewBuffer([]byte("awdawdawd")),
	}

	response = service.Detect(request)
	assert.Equal(t, codes.InvalidArgument, response.Status)
}

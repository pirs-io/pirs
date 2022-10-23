package mock

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

func CheckAuthorization(userRoles string, desiredRoles []string) bool {
	for _, role := range desiredRoles {
		if strings.Contains(userRoles, role) {
			return true
		}
	}
	return false
}

func FindOrCreateFile(filePath string) (*os.File, error) {
	return os.Create(filePath)
}

// process files
type DiskProcessStore struct {
	processFolder string
}

func NewDiskProcessStore(processFolder string) *DiskProcessStore {
	return &DiskProcessStore{
		processFolder: processFolder,
	}
}

func (store *DiskProcessStore) SaveProcessFile(processData bytes.Buffer) (bool, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return false, fmt.Errorf("cannot generate process uuid: %w", err)
	}

	processPath := fmt.Sprintf("%s/%s%s", store.processFolder, newUUID, ".pf")

	file, err := os.Create(processPath)
	if err != nil {
		return false, fmt.Errorf("cannot create process file: %w", err)
	}

	_, err = processData.WriteTo(file)
	if err != nil {
		return false, fmt.Errorf("cannot write process to file: %w", err)
	}

	return true, nil
}

package mock

import (
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

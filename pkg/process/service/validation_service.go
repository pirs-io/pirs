package service

import "os"

type ValidationService struct{}

func (vs *ValidationService) ValidateProcessFile(processFile os.File) (bool, error) {
	// todo (validation package files, structure)
	return true, nil
}

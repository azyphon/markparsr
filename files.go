package markparsr

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileValidator checks that required files for a Terraform module exist and aren't empty.
type FileValidator struct {
	rootDir         string
	requiredFiles   []string
	additionalFiles []string
}

// NewFileValidator creates a validator that checks for certain files in module directory.
// It accepts both the README path and module path to check files in both locations.
// additionalFiles specifies extra files to validate beyond the standard required ones.
func NewFileValidator(readmePath string, modulePath string, additionalFiles []string) *FileValidator {
	// Standard required files that should always be present
	requiredFiles := []string{
		readmePath,
		filepath.Join(modulePath, "outputs.tf"),
		filepath.Join(modulePath, "variables.tf"),
		filepath.Join(modulePath, "terraform.tf"),
	}

	// Convert relative additional file paths to absolute paths
	var absAdditionalFiles []string
	for _, file := range additionalFiles {
		if !filepath.IsAbs(file) {
			absAdditionalFiles = append(absAdditionalFiles, filepath.Join(modulePath, file))
		} else {
			absAdditionalFiles = append(absAdditionalFiles, file)
		}
	}

	return &FileValidator{
		rootDir:         modulePath,
		requiredFiles:   requiredFiles,
		additionalFiles: absAdditionalFiles,
	}
}

// Validate checks that all required and additional files exist and are not empty.
func (fv *FileValidator) Validate() []error {
	var allErrors []error

	// Validate required files
	for _, filePath := range fv.requiredFiles {
		if err := validateFile(filePath); err != nil {
			allErrors = append(allErrors, fmt.Errorf("required %v", err))
		}
	}

	// Validate additional files
	for _, filePath := range fv.additionalFiles {
		if err := validateFile(filePath); err != nil {
			allErrors = append(allErrors, fmt.Errorf("additional %v", err))
		}
	}

	return allErrors
}

// validateFile checks if a file exists and is not empty.
func validateFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", filepath.Base(filePath))
		}
		return fmt.Errorf("error accessing file: %s: %w", filepath.Base(filePath), err)
	}
	if fileInfo.Size() == 0 {
		return fmt.Errorf("file is empty: %s", filepath.Base(filePath))
	}
	return nil
}

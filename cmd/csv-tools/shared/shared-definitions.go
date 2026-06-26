package shared

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func PrepareFileHandles(originalCsvFile string) (originalCsv *os.File, tempCsv *os.File, prepErr error) {
	originalCsv, readErr := os.Open(originalCsvFile)
	if originalCsv == nil {
		prepErr = fmt.Errorf("failed to read original CSV file")
		return
	}
	if readErr != nil {
		prepErr = errors.WithMessage(readErr, "failed to read original CSV file")
		return
	}
	tempCsv, tempErr := os.CreateTemp("", fmt.Sprintf("*_%s", filepath.Base(originalCsvFile)))
	if tempCsv == nil {
		prepErr = fmt.Errorf("failed to create temp file")
		return
	}
	if tempErr != nil {
		prepErr = errors.WithMessage(tempErr, "failed to create temp file")
		return
	}
	return
}

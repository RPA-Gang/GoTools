package shared

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"

	"GoTools/pkg/helpers"
)

func ResolveInputPath(filePathPtr *string) (string, helpers.ErrMsg) {
	pipeInput, err := os.Stdin.Stat()
	if err != nil {
		return "", helpers.ErrMsg{Err: err, Code: helpers.ErrStdin}
	}
	if pipeInput.Mode()&os.ModeNamedPipe != 0 {
		reader := bufio.NewReader(os.Stdin)
		input, inputErr := reader.ReadString('\n')
		if inputErr != nil && !errors.Is(inputErr, io.EOF) {
			return "", helpers.ErrMsg{Err: inputErr, Code: helpers.ErrStdin}
		}
		return strings.TrimSpace(input), helpers.ErrMsg{Code: helpers.Success}
	}
	if *filePathPtr != "" {
		return *filePathPtr, helpers.ErrMsg{Code: helpers.Success}
	}
	return "", helpers.ErrMsg{
		Err:  fmt.Errorf("no CSV path provided from pipe nor --path flag"),
		Code: helpers.ErrNoInput,
	}
}

func ProcessCSV(path string, transform func([]string, int) ([]string, error)) helpers.ErrMsg {
	if exists, _ := helpers.PathExists(path); !exists {
		return helpers.ErrMsg{Err: fmt.Errorf("file '%s' does not exist", path), Code: helpers.ErrNoFile}
	}
	if !helpers.CheckExtension(path, ".csv") {
		return helpers.ErrMsg{
			Err:  fmt.Errorf("file '%s' is not a CSV file", path),
			Code: helpers.ErrInvalidFileType,
		}
	}
	tempFile, ioErr := readWriteCSV(path, transform)
	if ioErr != nil {
		return helpers.ErrMsg{Err: ioErr, Code: helpers.ErrReadWrite}
	}
	if err := os.Remove(path); err != nil {
		return helpers.ErrMsg{Err: err, Code: helpers.ErrWriteFile}
	}
	if err := helpers.MoveFile(tempFile, path); err != nil {
		return helpers.ErrMsg{Err: err, Code: helpers.ErrMoveFile}
	}
	log.Info("Successfully amended file", "original", filepath.Base(path), "amended", filepath.Base(tempFile))
	return helpers.ErrMsg{Code: helpers.Success}
}

func readWriteCSV(path string, transform func([]string, int) ([]string, error)) (string, error) {
	originalCsv, tempCsv, err := PrepareFileHandles(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = originalCsv.Close() }()
	defer func() { _ = tempCsv.Close() }()
	reader := csv.NewReader(originalCsv)
	writer := csv.NewWriter(tempCsv)
	defer writer.Flush()
	line := 0
	for {
		record, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return tempCsv.Name(), readErr
		}
		out, txErr := transform(record, line)
		if txErr != nil {
			return tempCsv.Name(), txErr
		}
		if err := writer.Write(out); err != nil {
			return tempCsv.Name(), err
		}
		line++
	}
	if err := writer.Error(); err != nil {
		return tempCsv.Name(), err
	}
	return tempCsv.Name(), nil
}

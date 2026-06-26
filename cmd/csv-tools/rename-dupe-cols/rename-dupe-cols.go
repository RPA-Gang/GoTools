package main

import (
	"flag"
	"time"

	"github.com/charmbracelet/log"

	"GoTools/cmd/csv-tools/shared"
	. "GoTools/pkg/helpers"
)

// main is the entry point of the program.
func main() {
	log.SetLevel(log.DebugLevel)
	startTime := time.Now()

	processingErr := ErrMsg{Code: Success}
	defer func() {
		log.Debug(
			"DONE!",
			"time", time.Since(startTime),
		)
		processingErr.Exit()
	}()
	filePathPtr := flag.String("path", "", "CSV file path")
	flag.Parse()
	inputPath, inputErr := shared.ResolveInputPath(filePathPtr)
	if inputErr.Code != Success {
		processingErr = inputErr
		return
	}
	processingErr = shared.ProcessCSV(inputPath, func(record []string, line int) ([]string, error) {
		if line == 0 {
			return RenameDuplicates(record, true), nil
		}
		return record, nil
	})
}

package main

import (
	"flag"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"GoTools/cmd/csv-tools/shared"
	. "GoTools/pkg/helpers"
)

//goland:noinspection DuplicatedCode
func main() {
	log.SetLevel(log.DebugLevel)
	startTime := time.Now()
	processingErr := ErrMsg{Code: Success}
	defer func() {
		log.Debug(
			"DONE!",
			"time",
			time.Since(startTime),
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
	processingErr = shared.ProcessCSV(inputPath, func(record []string, _ int) ([]string, error) {
		out := make([]string, len(record))
		for i, field := range record {
			out[i] = strings.TrimSpace(field)
		}
		return out, nil
	})
}

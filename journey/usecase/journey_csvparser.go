package usecase

import (
	"encoding/csv"
	"io"
	"os"
	"sync"

	"go.uber.org/zap"
)

func Parse(logger *zap.SugaredLogger, f *os.File) (int64, error) {
	logger.Debugw("Start parsing csv file",
		"file", f.Name(),
	)

	csvReader := csv.NewReader(f)

	logger.Debug("Reading headers")
	if _, err := csvReader.Read(); err != nil {
		logger.Error("Error reading headers",
			"error", err,
		)
		return 0, err
	}

	numWorkers := 10
	jobs := make(chan []string, numWorkers)
	res := make(chan []string)

	var workerGroup sync.WaitGroup

	worker := func(jobs <-chan []string, results chan<- []string) {
		logger.Debug("Worker started")
		for {
			select {
			case job, ok := <-jobs:
				if !ok {
					logger.Debug("Worker ended")
					return
				}
				logger.Debugw("Line received",
					"csvLine", job,
				)
				results <- job
			}
		}
	}

	logger.Debug("Workers Initializing")
	for w := 0; w < numWorkers; w++ {
		workerGroup.Add(1)
		go func() {
			defer workerGroup.Done()
			worker(jobs, res)
		}()
	}

	go func() {
		for {
			line, err := csvReader.Read()
			if err == io.EOF {
				logger.Debug("End of file reached")
				break
			}

			if err != nil {
				logger.Error("Error reading csv file:", err.Error())
				break
			}
			jobs <- line
		}
		close(jobs)
	}()

	go func() {
		workerGroup.Wait()
		close(res)
	}()

	nbReadLine := 0
	for r := range res {
		logger.Infow("New line read",
			"resLine", r,
		)
		nbReadLine++
	}

	return int64(nbReadLine), nil
}

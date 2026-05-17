package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"runtime"
	"sync"
)

func ProcessPipelineFromReader(reader io.Reader, tableName string, headers []string) (int, error) {
	numWorkers := MinInt(runtime.NumCPU(), 8)

	fmt.Printf("Using %d workers (CPU cores: %d)\n", numWorkers, runtime.NumCPU())

	recordChan := make(chan []string, 10000)
	resultChan := make(chan int, numWorkers)

	var wg sync.WaitGroup

	wg.Add(1)
	go streamFromReader(reader, recordChan, &wg)

	for i := range numWorkers {
		wg.Add(1)
		go streamWorker(i, tableName, headers, recordChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	totalLines := 0
	for count := range resultChan {
		totalLines += count
	}

	return totalLines, nil
}

func streamFromReader(reader io.Reader, recordChan chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(recordChan)

	bufferedReader := bufio.NewReaderSize(reader, 4*1024*1024)
	csvReader := csv.NewReader(bufferedReader)
	csvReader.Comma = ','
	csvReader.ReuseRecord = true
	csvReader.LazyQuotes = true

	_, _ = csvReader.Read()

	lineCount := 0
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		recordCopy := make([]string, len(record))
		copy(recordCopy, record)

		recordChan <- recordCopy
		lineCount++

		if lineCount%2000000 == 0 {
			fmt.Printf("Streamed: %.1fM lines\n", float64(lineCount)/1000000)
		}
	}

	fmt.Printf("Reader finished: %.1fM lines\n", float64(lineCount)/1000000)
}

func streamWorker(workerID int, tableName string, headers []string, recordChan <-chan []string, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	lineCount := 0
	batchSize := 200000
	batch := make([][]string, 0, batchSize)

	fmt.Printf("Worker %d started\n", workerID)

	for record := range recordChan {
		batch = append(batch, record)

		if len(batch) >= batchSize {
			if err := InsertBatch(tableName, headers, batch); err != nil {
				fmt.Printf("Worker %d batch error: %v\n", workerID, err)
				continue
			}
			lineCount += len(batch)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := InsertBatch(tableName, headers, batch); err == nil {
			lineCount += len(batch)
		}
	}

	resultChan <- lineCount
	fmt.Printf("Worker %d: %.1fM lines\n", workerID, float64(lineCount)/1000000)
}

package csv

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
)

func ProcessPipelineParallel(db *sql.DB, csvPath, tableName string, headers []string) (int, error) {
	numWorkers := MinInt(runtime.NumCPU(), 8)

	fmt.Printf("🚀 Using %d workers (CPU cores: %d)\n", numWorkers, runtime.NumCPU())

	recordChan := make(chan []string, 100000)
	resultChan := make(chan int, numWorkers)

	var wg sync.WaitGroup

	wg.Add(1)
	go streamCSVReader(csvPath, recordChan, &wg)

	for i := range numWorkers {
		wg.Add(1)
		go streamWorker(i, db, tableName, headers, recordChan, resultChan, &wg)
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

func streamCSVReader(csvPath string, recordChan chan<- []string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(recordChan)

	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Printf("❌ Reader error: %v\n", err)
		return
	}
	defer file.Close()

	bufferedReader := bufio.NewReaderSize(file, 4*1024*1024)
	reader := csv.NewReader(bufferedReader)
	reader.Comma = ','
	reader.ReuseRecord = true

	reader.Read()

	lineCount := 0
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		recordCopy := make([]string, len(record))
		copy(recordCopy, record)

		select {
		case recordChan <- recordCopy:
			lineCount++
		default:
			recordChan <- recordCopy
			lineCount++
		}

		if lineCount%2000000 == 0 {
			fmt.Printf("📖 Streamed: %.1fM lines\n", float64(lineCount)/1000000)
		}
	}

	fmt.Printf("📖 Reader finished: %.1fM lines\n", float64(lineCount)/1000000)
}

func streamWorker(workerID int, db *sql.DB, tableName string, headers []string, recordChan <-chan []string, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	lineCount := 0
	batchSize := 200000
	batch := make([][]string, 0, batchSize)

	fmt.Printf("⚡ Ultra Worker %d started\n", workerID)

	for record := range recordChan {
		batch = append(batch, record)

		if len(batch) >= batchSize {
			if err := InsertBatch(db, tableName, headers, batch); err != nil {
				fmt.Printf("❌ Worker %d batch error: %v\n", workerID, err)
				continue
			}
			lineCount += len(batch)
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := InsertBatch(db, tableName, headers, batch); err == nil {
			lineCount += len(batch)
		}
	}

	resultChan <- lineCount
	fmt.Printf("🏁 Ultra Worker %d: %.1fM lines\n", workerID, float64(lineCount)/1000000)
}

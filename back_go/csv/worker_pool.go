package csv

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/lib/pq"
)

type WorkerPool struct {
	db        *sql.DB
	tableName string
	headers   []string
	numWorkers int
}

type ChunkResult struct {
	ChunkID   int
	LineCount int
	Error     error
	Duration  time.Duration
}

func NewWorkerPool(db *sql.DB, tableName string, headers []string, numWorkers int) *WorkerPool {
	return &WorkerPool{
		db:         db,
		tableName:  tableName,
		headers:    headers,
		numWorkers: numWorkers,
	}
}

func (wp *WorkerPool) ProcessChunks(chunks []CSVChunk) (int, error) {
	fmt.Printf("🔥 Processing %d chunks with %d workers\n", len(chunks), wp.numWorkers)

	chunkChan := make(chan CSVChunk, len(chunks))
	resultChan := make(chan ChunkResult, len(chunks))

	for i, chunk := range chunks {
		chunk.ID = i
		chunkChan <- chunk
	}
	close(chunkChan)

	var wg sync.WaitGroup
	for i := 0; i < wp.numWorkers; i++ {
		wg.Add(1)
		go wp.worker(i, chunkChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return wp.collectResults(resultChan)
}

func (wp *WorkerPool) worker(workerID int, chunkChan <-chan CSVChunk, resultChan chan<- ChunkResult, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("⚡ Worker %d started\n", workerID)

	for chunk := range chunkChan {
		start := time.Now()
		lineCount, err := wp.processChunk(chunk)
		duration := time.Since(start)

		result := ChunkResult{
			ChunkID:   chunk.ID,
			LineCount: lineCount,
			Error:     err,
			Duration:  duration,
		}

		resultChan <- result

		if err != nil {
			fmt.Printf("❌ Worker %d failed chunk %d: %v\n", workerID, chunk.ID, err)
		} else {
			fmt.Printf("✅ Worker %d: chunk %d completed (%d lines in %.1fs)\n",
				workerID, chunk.ID, lineCount, duration.Seconds())
		}
	}

	fmt.Printf("🏁 Worker %d finished\n", workerID)
}

func (wp *WorkerPool) processChunk(chunk CSVChunk) (int, error) {
	tx, err := wp.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("cannot start transaction: %v", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(pq.CopyIn(wp.tableName, wp.headers...))
	if err != nil {
		return 0, fmt.Errorf("cannot prepare COPY: %v", err)
	}

	for _, record := range chunk.Lines {
		values := make([]any, len(record))
		for i, v := range record {
			values[i] = v
		}

		if _, err := stmt.Exec(values...); err != nil {
			return 0, fmt.Errorf("COPY error: %v", err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		return 0, fmt.Errorf("cannot finalize COPY: %v", err)
	}

	if err := stmt.Close(); err != nil {
		return 0, fmt.Errorf("cannot close COPY: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("cannot commit transaction: %v", err)
	}

	return len(chunk.Lines), nil
}

func (wp *WorkerPool) collectResults(resultChan <-chan ChunkResult) (int, error) {
	totalLines := 0

	for result := range resultChan {
		if result.Error != nil {
			return 0, fmt.Errorf("chunk %d failed: %v", result.ChunkID, result.Error)
		}

		totalLines += result.LineCount
		linesPerSec := float64(result.LineCount) / result.Duration.Seconds()

		fmt.Printf("📊 Chunk %d: %d lines (%.0f lines/sec)\n",
			result.ChunkID, result.LineCount, linesPerSec)
	}

	return totalLines, nil
}

package csv

import (
	"context"
	"csv-importer/config"
	"csv-importer/database"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

type WorkerPool struct {
	cfg        *config.Config
	tableName  string
	headers    []string
	numWorkers int
}

type ChunkResult struct {
	ChunkID   int
	LineCount int
	Error     error
	Duration  time.Duration
}

func NewWorkerPool(tableName string, headers []string, numWorkers int) *WorkerPool {
	return &WorkerPool{
		cfg:        config.Load(),
		tableName:  tableName,
		headers:    headers,
		numWorkers: numWorkers,
	}
}

func (wp *WorkerPool) ProcessChunks(chunks []CSVChunk) (int, error) {
	fmt.Printf("🔥 Processing %d chunks with %d workers (pgx)\n", len(chunks), wp.numWorkers)

	chunkChan := make(chan CSVChunk, len(chunks))
	resultChan := make(chan ChunkResult, len(chunks))

	for i, chunk := range chunks {
		chunk.ID = i
		chunkChan <- chunk
	}
	close(chunkChan)

	var wg sync.WaitGroup
	for i := range wp.numWorkers {
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

	fmt.Printf("⚡ pgx Worker %d started\n", workerID)

	conn, err := database.ConnectPgxNative(wp.cfg)
	if err != nil {
		fmt.Printf("❌ Worker %d connection failed: %v\n", workerID, err)
		return
	}
	defer conn.Close(context.Background())

	for chunk := range chunkChan {
		start := time.Now()
		lineCount, err := wp.processChunk(conn, chunk)
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

	fmt.Printf("🏁 pgx Worker %d finished\n", workerID)
}

func (wp *WorkerPool) processChunk(conn *pgx.Conn, chunk CSVChunk) (int, error) {
	rows := make([][]any, len(chunk.Lines))
	for i, record := range chunk.Lines {
		row := make([]any, len(record))
		for j, v := range record {
			row[j] = v
		}
		rows[i] = row
	}

	rowsAffected, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{wp.tableName},
		wp.headers,
		pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
			return rows[i], nil
		}),
	)

	if err != nil {
		return 0, fmt.Errorf("pgx copy from failed: %w", err)
	}

	return int(rowsAffected), nil
}

func (wp *WorkerPool) collectResults(resultChan <-chan ChunkResult) (int, error) {
	totalLines := 0

	for result := range resultChan {
		if result.Error != nil {
			return 0, fmt.Errorf("chunk %d failed: %v", result.ChunkID, result.Error)
		}

		totalLines += result.LineCount
		linesPerSec := float64(result.LineCount) / result.Duration.Seconds()

		fmt.Printf("📊 Chunk %d: %d lines (%.0f lines/sec) [pgx]\n",
			result.ChunkID, result.LineCount, linesPerSec)
	}

	return totalLines, nil
}

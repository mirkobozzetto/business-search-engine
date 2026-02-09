package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type CSVChunk struct {
	ID    int
	Lines [][]string
}

func CreateChunks(csvPath string, chunkSize int) ([]CSVChunk, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	// Skip header
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("cannot read header: %v", err)
	}

	var chunks []CSVChunk
	var currentChunk [][]string
	lineCount := 0

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("error reading line %d: %v", lineCount, err)
		}

		currentChunk = append(currentChunk, record)
		lineCount++

		if len(currentChunk) >= chunkSize {
			chunks = append(chunks, CSVChunk{Lines: currentChunk})
			currentChunk = nil
		}
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, CSVChunk{Lines: currentChunk})
	}

	fmt.Printf("📦 Created %d chunks from %d total lines\n", len(chunks), lineCount)
	return chunks, nil
}

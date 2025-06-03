package cache

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log/slog"
)

type CompressionResult struct {
	Data             []byte
	Key              string
	IsCompressed     bool
	OriginalSize     int
	FinalSize        int
	CompressionRatio float64
}

func (r *RedisCache) smartCompress(key string, jsonData []byte) (*CompressionResult, error) {
	originalSize := len(jsonData)

	if originalSize <= COMPRESSION_THRESHOLD {
		return &CompressionResult{
			Data:         jsonData,
			Key:          key,
			IsCompressed: false,
			OriginalSize: originalSize,
			FinalSize:    originalSize,
		}, nil
	}

	if originalSize > MAX_UNCOMPRESSED_SIZE {
		slog.Info("Large dataset detected, attempting compression",
			"key", key,
			"size_mb", originalSize/1024/1024)
	}

	compressed, err := r.compressData(jsonData)
	if err != nil {
		slog.Warn("Compression failed, using uncompressed",
			"key", key,
			"original_mb", originalSize/1024/1024,
			"error", err)

		if originalSize > MAX_UNCOMPRESSED_SIZE {
			return nil, fmt.Errorf("dataset too large and compression failed: %d MB", originalSize/1024/1024)
		}

		return &CompressionResult{
			Data:         jsonData,
			Key:          key,
			IsCompressed: false,
			OriginalSize: originalSize,
			FinalSize:    originalSize,
		}, nil
	}

	compressedSize := len(compressed)
	compressionRatio := float64(compressedSize) / float64(originalSize) * 100

	if compressedSize > MAX_COMPRESSED_SIZE {
		return nil, fmt.Errorf("dataset too large even after compression: %d MB compressed", compressedSize/1024/1024)
	}

	slog.Info("Smart compression successful",
		"key", key,
		"original_mb", originalSize/1024/1024,
		"compressed_mb", compressedSize/1024/1024,
		"compression_ratio", fmt.Sprintf("%.1f%%", compressionRatio),
		"saved_mb", (originalSize-compressedSize)/1024/1024)

	return &CompressionResult{
		Data:             compressed,
		Key:              COMPRESSION_PREFIX + key,
		IsCompressed:     true,
		OriginalSize:     originalSize,
		FinalSize:        compressedSize,
		CompressionRatio: compressionRatio,
	}, nil
}

func (r *RedisCache) compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		writer.Close()
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *RedisCache) decompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	return io.ReadAll(reader)
}

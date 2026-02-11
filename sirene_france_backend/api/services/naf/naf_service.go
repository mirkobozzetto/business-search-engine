package naf

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

type NafCode struct {
	Code         string `json:"code"`
	Label        string `json:"label"`
	SectionCode  string `json:"section_code"`
	SectionLabel string `json:"section_label"`
}

type NafSection struct {
	Code  string `json:"code"`
	Label string `json:"label"`
	Count int    `json:"count"`
}

type nafService struct {
	db *sql.DB
}

func NewNafService(db *sql.DB) *nafService {
	return &nafService{db: db}
}

func (s *nafService) SearchByLabel(ctx context.Context, query string, limit, offset int) ([]NafCode, int, error) {
	pattern := "%" + strings.ToLower(query) + "%"

	var total int
	var countErr error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		countErr = s.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM naf_reference WHERE label ILIKE $1`, pattern).Scan(&total)
	}()

	rows, err := s.db.QueryContext(ctx,
		`SELECT code, label, section_code, section_label FROM naf_reference WHERE label ILIKE $1 ORDER BY code LIMIT $2 OFFSET $3`,
		pattern, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("naf search failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	codes := make([]NafCode, 0, limit)
	for rows.Next() {
		var n NafCode
		if err := rows.Scan(&n.Code, &n.Label, &n.SectionCode, &n.SectionLabel); err != nil {
			continue
		}
		codes = append(codes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("naf rows error: %w", err)
	}

	wg.Wait()
	if countErr != nil {
		total = len(codes)
	}

	return codes, total, nil
}

func (s *nafService) ListSections(ctx context.Context) ([]NafSection, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT section_code, section_label, COUNT(*) as count FROM naf_reference GROUP BY section_code, section_label ORDER BY section_code`)
	if err != nil {
		return nil, fmt.Errorf("naf sections query failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var sections []NafSection
	for rows.Next() {
		var s NafSection
		if err := rows.Scan(&s.Code, &s.Label, &s.Count); err != nil {
			continue
		}
		sections = append(sections, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("naf sections rows error: %w", err)
	}

	return sections, nil
}

func (s *nafService) GetByCode(ctx context.Context, code string) (*NafCode, error) {
	var n NafCode
	err := s.db.QueryRowContext(ctx,
		`SELECT code, label, section_code, section_label FROM naf_reference WHERE code = $1`, code).
		Scan(&n.Code, &n.Label, &n.SectionCode, &n.SectionLabel)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("naf get by code failed: %w", err)
	}
	return &n, nil
}

func (s *nafService) GetBySection(ctx context.Context, sectionCode string) ([]NafCode, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT code, label, section_code, section_label FROM naf_reference WHERE section_code = $1 ORDER BY code`, sectionCode)
	if err != nil {
		return nil, fmt.Errorf("naf get by section failed: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var codes []NafCode
	for rows.Next() {
		var n NafCode
		if err := rows.Scan(&n.Code, &n.Label, &n.SectionCode, &n.SectionLabel); err != nil {
			continue
		}
		codes = append(codes, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("naf section rows error: %w", err)
	}

	return codes, nil
}

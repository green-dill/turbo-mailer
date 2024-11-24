package postfix

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/samber/lo"
)

type NopeDnsServer struct{}

func NewNopeDnsServer() (*NopeDnsServer, error) {
	return &NopeDnsServer{}, nil
}

func (s *NopeDnsServer) GetDNSRecords(domain, recordType string) ([]cloudflare.DNSRecord, error) {
	return nil, nil
}

func (s *NopeDnsServer) UpdateDNSRecord(domain, recordType, content string) error {
	const recordFile = "./storage/record.txt"
	record := fmt.Sprintf("%s\t%s\t%s", domain, recordType, content)

	// Create directory if not exists
	if err := os.MkdirAll("./storage", 0755); err != nil {
		return fmt.Errorf("create directory failed: %w", err)
	}

	// Read existing records
	records := make([]string, 0, 8)
	if data, err := os.ReadFile(recordFile); err == nil {
		records = lo.Filter(strings.Split(string(data), "\n"), func(r string, _ int) bool {
			return len(strings.TrimSpace(r)) > 0
		})
	}

	// Find and update existing record
	found := false
	for i, r := range records {
		if strings.HasPrefix(r, domain+"\t"+recordType) {
			records[i] = record
			found = true
			break
		}
	}

	// Append new record if not found
	if !found {
		records = append(records, record)
	}

	// Write back to file
	if err := os.WriteFile(recordFile, []byte(strings.Join(records, "\n\n\n")), 0644); err != nil {
		return fmt.Errorf("write record file failed: %w", err)
	}
	return nil
}

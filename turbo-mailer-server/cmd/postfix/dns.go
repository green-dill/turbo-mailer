package postfix

import (
	"github.com/cloudflare/cloudflare-go"
)

type DnsServer interface {
	GetDNSRecords(domain, recordType string) ([]cloudflare.DNSRecord, error)
	UpdateDNSRecord(domain, recordType, content string) error
}

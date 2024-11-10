package postfix

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type DnsServer interface {
	GetDNSRecords(domain, recordType string) ([]cloudflare.DNSRecord, error)
	UpdateDNSRecord(domain, recordType, content string) error
}

type CloudflareDnsServer struct {
	api  *cloudflare.API
	zone *cloudflare.ResourceContainer
}

func NewCloudflareDnsServer(apiKey, zoneID string) (*CloudflareDnsServer, error) {
	api, err := cloudflare.NewWithAPIToken(apiKey)
	if err != nil {
		return nil, err
	}

	zone := cloudflare.ZoneIdentifier(zoneID)

	return &CloudflareDnsServer{
		api:  api,
		zone: zone,
	}, nil
}

func NewDnsServer() (DnsServer, error) {
	return NewCloudflareDnsServer(viper.GetString("cloudflare.api_key"), viper.GetString("cloudflare.zone_id"))
}

func (s *CloudflareDnsServer) GetDNSRecords(domain, recordType string) ([]cloudflare.DNSRecord, error) {
	records, _, err := s.api.ListDNSRecords(context.Background(), s.zone, cloudflare.ListDNSRecordsParams{
		Type: recordType,
		Name: domain,
	})
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *CloudflareDnsServer) UpdateDNSRecord(domain, recordType, content string) error {
	records, err := s.GetDNSRecords(domain, recordType)
	if err != nil {
		return err
	}

	record, ok := lo.Find(records, func(record cloudflare.DNSRecord) bool {
		return record.Name == domain && record.Type == recordType
	})

	if !ok {
		// create new record
		createParams := cloudflare.CreateDNSRecordParams{
			Type:    recordType,
			Name:    domain,
			Content: content,
		}
		if recordType == "TXT" {
			createParams.TTL = 1
		}
		if recordType == "MX" {
			createParams.Priority = lo.ToPtr(uint16(10))
		}
		_, err = s.api.CreateDNSRecord(context.Background(), s.zone, createParams)
	} else {
		if record.Content == content {
			return nil
		}

		_, err = s.api.UpdateDNSRecord(context.Background(), s.zone, cloudflare.UpdateDNSRecordParams{
			ID:      record.ID,
			Type:    recordType,
			Name:    domain,
			Content: content,
		})

	}

	return err
}

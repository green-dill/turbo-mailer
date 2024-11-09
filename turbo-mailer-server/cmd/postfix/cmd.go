package postfix

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
	"turbo-mailer-server/internal/initialize"

	"github.com/imroc/req/v3"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	keepRunning = true
)

var Cmd = &cobra.Command{
	Use:   "postfix",
	Short: "Postfix dns record auto setup",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.Do(cmd.Context())
		boot()
	},
}

func init() {
	Cmd.Flags().BoolVarP(&keepRunning, "keep-running", "k", false, "keep running")
}

var (
	dnsServer DnsServer
)

func boot() {
	defer func() {
		if err := recover(); err != nil {
			log.Error().Interface("error", err).Msg("postfix server sidecar panic")
			time.Sleep(10 * time.Second)
			boot()
		}
	}()

	var err error
	dnsServer, err = NewDnsServer()
	if err != nil {
		log.Error().Err(err).Msg("create dns server failed")
		return
	}

	for {
		if err = checkAndSetDNSRecord(); err != nil {
			log.Error().Err(err).Msg("check and set dns record failed")
		}

		if !keepRunning {
			break
		}

		time.Sleep(1 * time.Minute)
	}
}

func checkAndSetDNSRecord() error {
	domain := viper.GetString("domain")
	if domain == "" {
		return fmt.Errorf("domain is not set")
	}

	externalIP, err := getExternalIP()
	if err != nil {
		return fmt.Errorf("get external ip failed: %w", err)
	}

	if err := checkAndSetSPFRecord(domain, externalIP); err != nil {
		return fmt.Errorf("check and set spf record failed: %w", err)
	}

	if err := checkAndSetDKIMRecord(domain); err != nil {
		return fmt.Errorf("check and set dkim record failed: %w", err)
	}

	return nil
}

func getExternalIP() (string, error) {
	r, err := req.R().Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	ip := strings.TrimSpace(r.String())

	// check ip is valid
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid ip: %s", ip)
	}

	return ip, nil
}

func checkAndSetSPFRecord(domain, externalIP string) error {
	txtRecords, err := dnsServer.GetDNSRecords(domain, "TXT")
	if err != nil {
		return fmt.Errorf("get dns records failed: %w", err)
	}

	var spfRecord string
	for _, record := range txtRecords {
		if strings.Contains(record.Content, "v=spf1") {
			spfRecord = record.Content
			break
		}
	}

	ipEntry := "ip4:" + externalIP
	if spfRecord != "" {
		if strings.Contains(spfRecord, ipEntry) {
			log.Info().Str("domain", domain).Str("ip", externalIP).Msg("SPF record already contains IP")
			return nil
		}
	}

	var newSPFRecord string
	if spfRecord == "" {
		newSPFRecord = fmt.Sprintf("v=spf1 %s -all", ipEntry)
	} else {
		parts := strings.Fields(spfRecord)
		lastPart := parts[len(parts)-1]
		if lastPart == "-all" || lastPart == "~all" {
			parts = parts[:len(parts)-1]
		}
		parts = append(parts, ipEntry, "-all")
		newSPFRecord = strings.Join(parts, " ")
	}

	log.Info().
		Str("domain", domain).
		Str("old_record", spfRecord).
		Str("new_record", newSPFRecord).
		Msg("Updating SPF record")

	return updateDNSRecord(domain, "TXT", newSPFRecord)
}

func generatePrivateKey(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create directory failed: %w", err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generate private key failed: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create private key file failed: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, privateKeyPEM); err != nil {
		return fmt.Errorf("write private key file failed: %w", err)
	}

	log.Info().Str("path", path).Msg("Generated new private key")
	return nil
}

func checkAndSetDKIMRecord(domain string) error {
	privateKeyPath := "./storage/mail.pem"

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		log.Info().Str("path", privateKeyPath).Msg("Private key file not found, generating new one")
		if err := generatePrivateKey(privateKeyPath); err != nil {
			return fmt.Errorf("generate private key failed: %w", err)
		}
	}

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("read private key failed: %w", err)
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("parse private key failed: %w", err)
	}

	publicKey := &privateKey.PublicKey
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKeyBytes)

	dkimRecord := fmt.Sprintf("v=DKIM1; k=rsa; h=sha256; s=mail; p=%s", publicKeyBase64)

	selector := "mail"
	dkimDomain := fmt.Sprintf("%s._domainkey.%s", selector, domain)

	txtRecords, err := dnsServer.GetDNSRecords(dkimDomain, "TXT")
	if err != nil {
		return fmt.Errorf("get dns records failed: %w", err)
	}

	for _, record := range txtRecords {
		if strings.HasSuffix(strings.Trim(record.Content, `"`), dkimRecord[len(dkimRecord)-16:]) {
			log.Info().Str("domain", dkimDomain).Msg("DKIM record is up to date")
			return nil
		}
	}

	log.Info().
		Str("domain", dkimDomain).
		Str("record", dkimRecord).
		Msg("Updating DKIM record")

	return updateDNSRecord(dkimDomain, "TXT", dkimRecord)
}

func updateDNSRecord(domain, recordType, content string) error {
	log.Info().Str("domain", domain).Str("record_type", recordType).Str("content", content).Msg("Updating DNS record")

	content = fmt.Sprintf(`"%s"`, content)

	if err := dnsServer.UpdateDNSRecord(domain, recordType, content); err != nil {
		return fmt.Errorf("update dns record failed: %w", err)
	}

	return nil
}

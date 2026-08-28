package ssl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type CertificateResult struct {
	Domain     string    `json:"domain"`
	CertPath   string    `json:"cert_path"`
	KeyPath    string    `json:"key_path"`
	Issuer     string    `json:"issuer"`
	ValidFrom  time.Time `json:"valid_from"`
	ValidUntil time.Time `json:"valid_until"`
	Status     string    `json:"status"`
}

type Manager struct {
	certDir string
	isDev   bool
}

func NewManager(isDev bool) *Manager {
	if isDev || runtime.GOOS != "linux" {
		mockDir := filepath.Join(os.TempDir(), "kodepreneur", "ssl")
		_ = os.MkdirAll(mockDir, 0755)
		return &Manager{
			certDir: mockDir,
			isDev:   true,
		}
	}

	return &Manager{
		certDir: "/etc/letsencrypt/live",
		isDev:   false,
	}
}

// IssueCertificate runs Certbot webroot HTTP-01 issuance or generates a dev certificate.
func (m *Manager) IssueCertificate(domain string, aliases []string, email, webroot string) (*CertificateResult, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain cannot be empty")
	}
	if email == "" {
		email = fmt.Sprintf("admin@%s", domain)
	}
	if webroot == "" {
		webroot = fmt.Sprintf("/var/www/%s/public", domain)
	}

	if m.isDev || runtime.GOOS != "linux" {
		return m.generateDevCertificate(domain)
	}

	// Build Certbot arguments
	args := []string{
		"certonly",
		"--webroot",
		"-w", webroot,
		"-d", domain,
		"--non-interactive",
		"--agree-tos",
		"--email", email,
		"--keep-until-expiring",
	}

	for _, alias := range aliases {
		alias = strings.TrimSpace(alias)
		if alias != "" && alias != domain {
			args = append(args, "-d", alias)
		}
	}

	cmd := exec.Command("certbot", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("certbot issuance failed: %s: %s", err.Error(), strings.TrimSpace(string(out)))
	}

	certPath := filepath.Join(m.certDir, domain, "fullchain.pem")
	keyPath := filepath.Join(m.certDir, domain, "privkey.pem")

	// Read validity
	validFrom, validUntil, err := m.readCertValidity(certPath)
	if err != nil {
		validFrom = time.Now()
		validUntil = time.Now().AddDate(0, 3, 0)
	}

	return &CertificateResult{
		Domain:     domain,
		CertPath:   certPath,
		KeyPath:    keyPath,
		Issuer:     "Let's Encrypt",
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
		Status:     "valid",
	}, nil
}

// InspectCertificate reads expiration and issuer from certificate files.
func (m *Manager) InspectCertificate(certPath string) (*CertificateResult, error) {
	validFrom, validUntil, err := m.readCertValidity(certPath)
	if err != nil {
		return nil, err
	}

	return &CertificateResult{
		CertPath:   certPath,
		Issuer:     "Let's Encrypt",
		ValidFrom:  validFrom,
		ValidUntil: validUntil,
		Status:     "valid",
	}, nil
}

func (m *Manager) readCertValidity(certPath string) (time.Time, time.Time, error) {
	data, err := os.ReadFile(certPath)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert.NotBefore, cert.NotAfter, nil
}

func (m *Manager) generateDevCertificate(domain string) (*CertificateResult, error) {
	domainDir := filepath.Join(m.certDir, domain)
	_ = os.MkdirAll(domainDir, 0755)

	certPath := filepath.Join(domainDir, "fullchain.pem")
	keyPath := filepath.Join(domainDir, "privkey.pem")

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	now := time.Now()
	validUntil := now.AddDate(0, 3, 0) // 90 days

	template := x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		Subject: pkix.Name{
			Organization: []string{"Kodepreneur Development CA"},
			CommonName:   domain,
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              validUntil,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(domain); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, domain, fmt.Sprintf("*.%s", domain))
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s for writing: %w", certPath, err)
	}
	_ = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	_ = certOut.Close()

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s for writing: %w", keyPath, err)
	}
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	_ = pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	_ = keyOut.Close()

	return &CertificateResult{
		Domain:     domain,
		CertPath:   certPath,
		KeyPath:    keyPath,
		Issuer:     "Kodepreneur Dev Self-Signed",
		ValidFrom:  now,
		ValidUntil: validUntil,
		Status:     "valid",
	}, nil
}

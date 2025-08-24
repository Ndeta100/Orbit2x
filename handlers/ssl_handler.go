// handlers/ssl_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Ndeta100/orbit2x/internal"

	"github.com/Ndeta100/orbit2x/views/ssl"
)

// HandleSSLIndex renders the SSL Certificate checker page
func HandleSSLIndex(w http.ResponseWriter, r *http.Request) error {
	return ssl.SSLChecker().Render(r.Context(), w)
}

// HandleSSLCheck analyzes SSL certificates for a given hostname
func HandleSSLCheck(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return ssl.SSLCertificateResult(ssl.CertificateInfo{
			Error: "Failed to parse form data",
		}).Render(r.Context(), w)
	}

	// Get hostname from form
	hostname := r.FormValue("hostname")
	if hostname == "" {
		return ssl.SSLCertificateResult(ssl.CertificateInfo{
			Error: "Hostname is required",
		}).Render(r.Context(), w)
	}

	// Clean the hostname (remove protocols)
	hostname = strings.TrimPrefix(hostname, "http://")
	hostname = strings.TrimPrefix(hostname, "https://")
	hostname = strings.Split(hostname, "/")[0] // Remove path if any

	// Set a timeout for the certificate check (10 seconds)
	timeout := 10

	// Get certificate details
	certDetails, err := internal.GetCertificateDetails(hostname, timeout)
	if err != nil {
		return ssl.SSLCertificateResult(ssl.CertificateInfo{
			Error: fmt.Sprintf("Failed to check certificate: %v", err),
		}).Render(r.Context(), w)
	}

	// Check expiration status (30-day threshold for "expiring soon")
	internal.CheckExpirationStatus(&certDetails, 30)

	// Convert certificate details to our template format
	certInfo := ssl.CertificateInfo{
		DaysUntilExpiration: certDetails.DaysUntilExpiration,
		IssuerName:          certDetails.IssuerName,
		SubjectName:         certDetails.SubjectName,
		SerialNumber:        certDetails.SerialNumber,
		ExpiringSoon:        certDetails.ExpiringSoon,
		Expired:             certDetails.Expired,
		Hostname:            hostname,
		TimeTaken:           certDetails.TimeTaken,
		ExpirationDate:      certDetails.ExpirationDate,
	}

	// Render the result
	return ssl.SSLCertificateResult(certInfo).Render(r.Context(), w)
}

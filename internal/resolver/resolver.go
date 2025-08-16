package resolver

import (
	"fmt"
	"github.com/domainr/whois"
	"github.com/miekg/dns"
	"net"
)

type Resolver interface {
	LookupAddr(ip string) ([]string, error)
	LookupIP(host string) ([]net.IP, error)
	LookupCNAME(host string) (string, error)
	LookupMX(host string) ([]*net.MX, error)
	LookupNS(host string) ([]*net.NS, error)
	LookupTXT(host string) ([]string, error)
	LookupSOA(host string) ([]string, error)
	LookupWHOIS(host string) (string, error)
}

// DefaultResolver uses Go's net package for lookups
type DefaultResolver struct{}

func (r *DefaultResolver) LookupAddr(ip string) ([]string, error) {
	return net.LookupAddr(ip)
}

func (r *DefaultResolver) LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}

func (r *DefaultResolver) LookupCNAME(host string) (string, error) {
	return net.LookupCNAME(host)
}

func (r *DefaultResolver) LookupMX(host string) ([]*net.MX, error) {
	return net.LookupMX(host)
}

func (r *DefaultResolver) LookupNS(host string) ([]*net.NS, error) {
	return net.LookupNS(host)
}

func (r *DefaultResolver) LookupTXT(host string) ([]string, error) {
	return net.LookupTXT(host)
}

// LookupSOA performs SOA record lookup using miekg/dns
func (r *DefaultResolver) LookupSOA(host string) ([]string, error) {
	client := &dns.Client{}
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(host), dns.TypeSOA)

	resp, _, err := client.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return nil, fmt.Errorf("SOA Lookup failed: %v", err)
	}

	var soa []string
	for _, ans := range resp.Answer {
		if soaRecord, ok := ans.(*dns.SOA); ok {
			soa = append(soa, fmt.Sprintf("Master Name Server: %s", soaRecord.Ns))
			soa = append(soa, fmt.Sprintf("Responsible Email: %s", soaRecord.Mbox))
			soa = append(soa, fmt.Sprintf("Serial: %d", soaRecord.Serial))
			soa = append(soa, fmt.Sprintf("Refresh: %d seconds", soaRecord.Refresh))
			soa = append(soa, fmt.Sprintf("Retry: %d seconds", soaRecord.Retry))
			soa = append(soa, fmt.Sprintf("Expire: %d seconds", soaRecord.Expire))
			soa = append(soa, fmt.Sprintf("Minimum TTL: %d seconds", soaRecord.Minttl))
		}
	}
	return soa, nil
}

// LookupWHOIS performs a WHOIS lookup using domain/whois
func (r *DefaultResolver) LookupWHOIS(domain string) (string, error) {
	// create WHOIS request for a given domain
	request, err := whois.NewRequest(domain)
	if err != nil {
		return "", fmt.Errorf("failed to create WHOIS request: %v", err)
	}

	// fetch WHOIS response
	response, err := whois.DefaultClient.Fetch(request)
	if err != nil {
		return "", fmt.Errorf("WHOIS lookup failed: %v", err)
	}

	// return WHOIS response as string
	return response.String(), nil
}

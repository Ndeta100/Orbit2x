package resolver

import "fmt"

// DNSResult represents the results of a specific DNS record type lookup
type DNSResult struct {
	Type    string
	Records []string
	Error   string
}

// DNSLookupResults contains all DNS lookup results for a domain
type DNSLookupResults struct {
	Domain  string
	Results map[string]DNSResult
}

// PerformAllLookups executes all DNS lookups for a domain and returns structured results
func PerformAllLookups(r Resolver, domain string) DNSLookupResults {
	results := DNSLookupResults{
		Domain:  domain,
		Results: make(map[string]DNSResult),
	}

	// A Records
	aRecords, aErr := getARecords(r, domain)
	if aErr != nil {
		results.Results["A"] = DNSResult{Type: "A", Error: aErr.Error()}
	} else {
		results.Results["A"] = DNSResult{Type: "A", Records: aRecords}
	}

	// AAAA Records
	aaaaRecords, aaaaErr := getAAAARecords(r, domain)
	if aaaaErr != nil {
		results.Results["AAAA"] = DNSResult{Type: "AAAA", Error: aaaaErr.Error()}
	} else {
		results.Results["AAAA"] = DNSResult{Type: "AAAA", Records: aaaaRecords}
	}

	// CNAME Record
	cname, cnameErr := r.LookupCNAME(domain)
	if cnameErr != nil {
		results.Results["CNAME"] = DNSResult{Type: "CNAME", Error: cnameErr.Error()}
	} else {
		results.Results["CNAME"] = DNSResult{Type: "CNAME", Records: []string{cname}}
	}

	// MX Records
	mxRecords, mxErr := getMXRecords(r, domain)
	if mxErr != nil {
		results.Results["MX"] = DNSResult{Type: "MX", Error: mxErr.Error()}
	} else {
		results.Results["MX"] = DNSResult{Type: "MX", Records: mxRecords}
	}

	// NS Records
	nsRecords, nsErr := getNSRecords(r, domain)
	if nsErr != nil {
		results.Results["NS"] = DNSResult{Type: "NS", Error: nsErr.Error()}
	} else {
		results.Results["NS"] = DNSResult{Type: "NS", Records: nsRecords}
	}

	// TXT Records
	txtRecords, txtErr := r.LookupTXT(domain)
	if txtErr != nil {
		results.Results["TXT"] = DNSResult{Type: "TXT", Error: txtErr.Error()}
	} else {
		results.Results["TXT"] = DNSResult{Type: "TXT", Records: txtRecords}
	}

	// SOA Record
	soaRecords, soaErr := r.LookupSOA(domain)
	if soaErr != nil {
		results.Results["SOA"] = DNSResult{Type: "SOA", Error: soaErr.Error()}
	} else {
		results.Results["SOA"] = DNSResult{Type: "SOA", Records: soaRecords}
	}

	// WHOIS Info
	whoisInfo, whoisErr := r.LookupWHOIS(domain)
	if whoisErr != nil {
		results.Results["WHOIS"] = DNSResult{Type: "WHOIS", Error: whoisErr.Error()}
	} else {
		// Split WHOIS info into lines for better display
		results.Results["WHOIS"] = DNSResult{Type: "WHOIS", Records: []string{whoisInfo}}
	}

	return results
}

// Helper functions to convert your existing functions into returnable values

func getARecords(r Resolver, domain string) ([]string, error) {
	ips, err := r.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	var records []string
	for _, ip := range ips {
		if ip.To4() != nil {
			records = append(records, ip.String())
		}
	}

	return records, nil
}

func getAAAARecords(r Resolver, domain string) ([]string, error) {
	ips, err := r.LookupIP(domain)
	if err != nil {
		return nil, err
	}

	var records []string
	for _, ip := range ips {
		if ip.To16() != nil && ip.To4() == nil {
			records = append(records, ip.String())
		}
	}

	return records, nil
}

func getMXRecords(r Resolver, domain string) ([]string, error) {
	mxRecords, err := r.LookupMX(domain)
	if err != nil {
		return nil, err
	}

	var records []string
	for _, mx := range mxRecords {
		records = append(records, fmt.Sprintf("%s (Priority: %d)", mx.Host, mx.Pref))
	}

	return records, nil
}

func getNSRecords(r Resolver, domain string) ([]string, error) {
	nsRecords, err := r.LookupNS(domain)
	if err != nil {
		return nil, err
	}

	var records []string
	for _, ns := range nsRecords {
		records = append(records, ns.Host)
	}

	return records, nil
}

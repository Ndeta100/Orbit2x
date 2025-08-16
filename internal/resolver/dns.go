package resolver

import (
	"fmt"
	"log"
	"strings"
)

func PerformDNSLookups(resolver Resolver, domain string) {
	header := fmt.Sprintf("DNS Lookup Results for: %s", domain)
	fmt.Println("\n" + header)
	fmt.Println(strings.Repeat("=", len(header)))

	lookupA(resolver, domain)
	lookupAAAA(resolver, domain)
	lookupCNAME(resolver, domain)
	lookupMX(resolver, domain)
	lookupNS(resolver, domain)
	lookupTXT(resolver, domain)
	lookupSOA(resolver, domain)
	lookupWHOIS(resolver, domain)
}

func lookupA(resolver Resolver, domain string) {
	ips, err := resolver.LookupIP(domain)
	if err != nil {
		fmt.Printf("A Record: Failed (%v)\n", err)
		return
	}

	fmt.Println("\nA Records (IPv4):")
	for _, ip := range ips {
		if ip.To4() != nil {
			fmt.Printf("- %s\n", ip)
		}
	}
}

func lookupAAAA(resolver Resolver, domain string) {
	ips, err := resolver.LookupIP(domain)
	if err != nil {
		fmt.Printf("AAAA Record: Failed (%v)\n", err)
		return
	}

	fmt.Println("\nAAAA Records (IPv6):")
	for _, ip := range ips {
		if ip.To16() != nil && ip.To4() == nil {
			fmt.Printf("- %s\n", ip)
		}
	}
}

func lookupCNAME(resolver Resolver, domain string) {
	cname, err := resolver.LookupCNAME(domain)
	if err != nil {
		fmt.Printf("CNAME Record: Failed (%v)\n", err)
		return
	}

	fmt.Println("\nCNAME Record:")
	fmt.Printf("- %s\n", cname)
}

func lookupMX(resolver Resolver, domain string) {
	mxRecords, err := resolver.LookupMX(domain)
	if err != nil {
		fmt.Printf("MX Record: Failed (%v)\n", err)
		return
	}

	fmt.Println("\nMX Records (Mail Servers):")
	for _, mx := range mxRecords {
		fmt.Printf("- %s (Priority: %d)\n", mx.Host, mx.Pref)
	}
}

func lookupNS(resolver Resolver, domain string) {
	nsRecords, err := resolver.LookupNS(domain)
	if err != nil {
		fmt.Printf("NS Record: Failed (%v)\n", err)
		return
	}

	fmt.Println("\nNS Records (Name Servers):")
	for _, ns := range nsRecords {
		fmt.Printf("- %s\n", ns.Host)
	}
}

func lookupTXT(resolver Resolver, domain string) {
	txtRecords, err := resolver.LookupTXT(domain)
	if err != nil {
		fmt.Printf("TXT Record: Failed (%v)\n", err)
		return
	}

	fmt.Println("\nTXT Records:")
	for _, txt := range txtRecords {
		fmt.Printf("- %s\n", txt)
	}
}

func lookupSOA(r Resolver, domain string) {
	soaRecords, err := r.LookupSOA(domain)
	if err != nil {
		log.Printf("SOA Record lookup failed: %v\n", err)
		return
	}
	fmt.Println("\nSOA Record:")
	for _, line := range soaRecords {
		fmt.Println(line)
	}
}

func lookupWHOIS(r Resolver, domain string) string {
	whoisInfo, err := r.LookupWHOIS(domain)
	if err != nil {
		log.Printf("WHOIS lookup failed: %v\n", err)
		return fmt.Sprintf("WHOIS lookup failed: %v", err)
	}
	fmt.Println("\nWHOIS Information:")
	fmt.Println(whoisInfo)
	return whoisInfo
}

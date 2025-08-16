// handlers/subnet_handler.go
package handlers

import (
	"fmt"
	"math"
	"net"
	"net/http"

	"github.com/Ndeta100/orbit2x/views/subnet" // Adjust to your actual path
)

// HandleSubnetIndex renders the subnet calculator page
func HandleSubnetIndex(w http.ResponseWriter, r *http.Request) error {
	return subnet.Index().Render(r.Context(), w)
}

// HandleSubnetCalculateCIDR calculates subnet details from CIDR notation
func HandleSubnetCalculateCIDR(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: "Failed to parse form data",
		}).Render(r.Context(), w)
	}

	// Get CIDR from form
	cidrNotation := r.FormValue("cidr")
	if cidrNotation == "" {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: "CIDR notation is required",
		}).Render(r.Context(), w)
	}

	// Parse CIDR
	ip, ipNet, err := net.ParseCIDR(cidrNotation)
	if err != nil {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: fmt.Sprintf("Invalid CIDR notation: %v", err),
		}).Render(r.Context(), w)
	}

	// Calculate subnet information
	result, err := calculateSubnet(ip, ipNet)
	if err != nil {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: fmt.Sprintf("Calculation error: %v", err),
		}).Render(r.Context(), w)
	}

	// Render the result
	return subnet.SubnetResults(result).Render(r.Context(), w)
}

// HandleSubnetCalculateMask calculates subnet details from IP and subnet mask
func HandleSubnetCalculateMask(w http.ResponseWriter, r *http.Request) error {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: "Failed to parse form data",
		}).Render(r.Context(), w)
	}

	// Get IP and mask from form
	ipAddress := r.FormValue("ip")
	subnetMask := r.FormValue("mask")

	if ipAddress == "" || subnetMask == "" {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: "IP address and subnet mask are required",
		}).Render(r.Context(), w)
	}

	// Parse IP
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: "Invalid IP address format",
		}).Render(r.Context(), w)
	}

	// Parse subnet mask
	mask := net.ParseIP(subnetMask)
	if mask == nil {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: "Invalid subnet mask format",
		}).Render(r.Context(), w)
	}

	// Convert mask to IPNet
	ipv4Mask := net.IPv4Mask(mask[12], mask[13], mask[14], mask[15])
	//cidr, _ := net.IPv4Mask(mask[12], mask[13], mask[14], mask[15]).Size()
	ipNet := &net.IPNet{
		IP:   ip.Mask(ipv4Mask),
		Mask: ipv4Mask,
	}

	// Calculate subnet information
	result, err := calculateSubnet(ip, ipNet)
	if err != nil {
		return subnet.SubnetResults(subnet.SubnetResult{
			Error: fmt.Sprintf("Calculation error: %v", err),
		}).Render(r.Context(), w)
	}

	// Render the result
	return subnet.SubnetResults(result).Render(r.Context(), w)
}

// calculateSubnet performs the subnet calculations
func calculateSubnet(ip net.IP, ipNet *net.IPNet) (subnet.SubnetResult, error) {
	// Ensure we're working with IPv4
	ip = ip.To4()
	if ip == nil {
		return subnet.SubnetResult{}, fmt.Errorf("only IPv4 addresses are supported")
	}

	// Get network and broadcast addresses
	networkIP := ipNet.IP.To4()

	// Calculate broadcast address (network address OR NOT subnet mask)
	mask := ipNet.Mask
	broadcastIP := make(net.IP, len(networkIP))
	for i := 0; i < len(networkIP); i++ {
		broadcastIP[i] = networkIP[i] | ^mask[i]
	}

	// Calculate first and last usable IP addresses
	firstUsableIP := make(net.IP, len(networkIP))
	lastUsableIP := make(net.IP, len(broadcastIP))

	copy(firstUsableIP, networkIP)
	copy(lastUsableIP, broadcastIP)

	// For most subnets, first usable IP is network address + 1,
	// except for /31 and /32 subnets which have special rules
	maskSize, _ := mask.Size()
	if maskSize < 31 {
		firstUsableIP[3]++
		lastUsableIP[3]--
	}

	// Calculate total hosts and usable hosts
	totalHosts := math.Pow(2, float64(32-maskSize))
	usableHosts := totalHosts - 2

	// Special cases for /31 (point-to-point) and /32 (single host) networks
	if maskSize == 31 {
		usableHosts = 2 // RFC 3021 - No network or broadcast addresses in /31
	} else if maskSize == 32 {
		usableHosts = 1 // Single host
	}

	if usableHosts < 0 {
		usableHosts = 0
	}

	// Convert subnet mask to binary
	maskBinary := ""
	for i := 0; i < len(mask); i++ {
		maskBinary += fmt.Sprintf("%08b.", mask[i])
	}
	maskBinary = maskBinary[:len(maskBinary)-1] // Remove trailing dot

	// Create result
	result := subnet.SubnetResult{
		NetworkAddress:   networkIP.String(),
		BroadcastAddress: broadcastIP.String(),
		FirstUsableIP:    firstUsableIP.String(),
		LastUsableIP:     lastUsableIP.String(),
		TotalHosts:       int(totalHosts),
		UsableHosts:      int(usableHosts),
		SubnetMaskDec:    net.IP(mask).String(),
		SubnetMaskBin:    maskBinary,
		CIDR:             maskSize,
	}

	return result, nil
}

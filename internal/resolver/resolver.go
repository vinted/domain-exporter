package resolver

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

func QueryNameserver(nameserver, domain string, recordType uint16) (*dns.Msg, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), recordType)
	msg.RecursionDesired = false

	client := &dns.Client{Timeout: 5 * time.Second}

	resp, _, err := client.Exchange(msg, net.JoinHostPort(nameserver, "53"))
	if err != nil {
		return nil, fmt.Errorf("failed to query nameserver %s: %w", nameserver, err)
	}

	return resp, nil
}

func GetAuthoritativeNameservers(tld string) ([]string, error) {
	nsRecords, err := net.LookupNS(tld)
	if err != nil {
		return nil, fmt.Errorf("NS lookup failed: %w", err)
	}

	if len(nsRecords) == 0 {
		return nil, fmt.Errorf("no nameservers found for TLD %s", tld)
	}

	nameservers := make([]string, len(nsRecords))
	for i, ns := range nsRecords {
		nameservers[i] = ns.Host
	}

	return nameservers, nil
}

package collector

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vinted/domain-exporter/config"
	"github.com/vinted/domain-exporter/internal/resolver"
)

type DomainCollector struct {
	domains []config.Domain
	status  *prometheus.Desc
}

func NewDomainCollector(domains []config.Domain) *DomainCollector {
	return &DomainCollector{
		domains: domains,
		status: prometheus.NewDesc(
			"domain_rcode_status",
			"DNS RCODE of the domain query (0=NOERROR, 2=SERVFAIL, 3=NXDOMAIN, etc.). -1 indicates a network/query error.",
			[]string{"domain"},
			nil,
		),
	}
}

func (c *DomainCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.status
}

func (c *DomainCollector) Collect(ch chan<- prometheus.Metric) {
	for _, domain := range c.domains {
		rcode, err := c.getDomainRcode(domain)

		value := float64(rcode)
		if err != nil {
			slog.Error("error querying domain", "domain", domain.Name, "error", err)
			value = -1
		} else if rcode != dns.RcodeSuccess {
			slog.Warn("domain returned non-zero RCODE", "domain", domain.Name, "rcode", rcode, "rcode_name", dns.RcodeToString[rcode])
		}

		ch <- prometheus.MustNewConstMetric(c.status, prometheus.GaugeValue, value, domain.Name)
	}
}

func (c *DomainCollector) getDomainRcode(domain config.Domain) (int, error) {
	rcode := -1

	nameservers, err := resolver.GetAuthoritativeNameservers(domain.TLD)
	if err != nil {
		return rcode, fmt.Errorf("failed to fetch authoritative nameservers: %w", err)
	}

	var queryErr error
	for _, ns := range nameservers {
		resp, err := resolver.QueryNameserver(ns, domain.Name, dns.TypeNS)
		if err == nil {
			return resp.Rcode, nil
		}
		queryErr = err
	}

	return rcode, fmt.Errorf("failed to query domain: %w", queryErr)
}

func staticPage(w http.ResponseWriter, r *http.Request) {
	page := `<html>
    <head><title>Domain Exporter</title></head>
    <body>
    <h1>Domain Exporter</h1>
    <p><a href='metrics'>Metrics</a></p>
    </body>
    </html>`
	fmt.Fprintln(w, page)
}

func Start(httpListenAddress string, domains []config.Domain) error {
	collector := NewDomainCollector(domains)
	prometheus.MustRegister(collector)

	mux := http.NewServeMux()
	mux.HandleFunc("/", staticPage)
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    httpListenAddress,
		Handler: mux,
	}

	slog.Info("Listening on " + httpListenAddress)
	return server.ListenAndServe()
}

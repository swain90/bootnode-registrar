package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

const (
	defaultTimeout time.Duration = 5 * time.Second
)

var (
	conf *dns.ClientConfig
)

func query(fqdn string, qtype uint16) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(fqdn, qtype)
	m.RecursionDesired = true

	c := new(dns.Client)
	c.ReadTimeout = defaultTimeout

	for _, server := range conf.Servers {
		r, _, err := c.Exchange(m, server+":"+conf.Port)
		if err != nil {
			return nil, fmt.Errorf("failed to exchange with server %s: %w", server, err)
		}
		if r == nil {
			return nil, errors.New("received nil response from server")
		}
		if r.Rcode == dns.RcodeNameError || r.Rcode == dns.RcodeSuccess {
			return r, nil
		}
	}

	return nil, errors.New("no name server to answer the question")
}

func ResolveAddressRecord(dnsName string) ([]string, error) {
	var err error
	conf, err = dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil || conf == nil {
		return nil, fmt.Errorf("failed to initialize local resolver. %s", err)
	}

	r, err := query(dns.Fqdn(dnsName), dns.TypeA)
	if err != nil || r == nil {
		return nil, err
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("error getting address. %s", dns.RcodeToString[r.Rcode])
	}

	var ips []string
	for _, ans := range r.Answer {
		switch ans := ans.(type) {
		case *dns.A:
			ips = append(ips, ans.A.String())
		}
	}

	return ips, nil
}

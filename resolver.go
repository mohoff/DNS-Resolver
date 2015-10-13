package main

import (
	"errors"
	"log"
	"strings"

	"github.com/miekg/dns"
)

var (
	DNSPORT = "53"
)

// Resolver struct holds essential parameters which are used in the context
// of github.com/miekg/dns
type Resolver struct {
	dnsServers  []string
	rrTypes     []string
	edns        bool
	dnsQueryMsg *dns.Msg
	dnsClient   *dns.Client
	Result      *Result
}

// NewResolver returns an initialized Resolver struct.
func NewResolver(config *Config) *Resolver {
	msg := &dns.Msg{}
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.SetQuestion("", dns.TypeANY)

	if config.edns {
		msg = handleEDNS(msg)
	}

	return &Resolver{
		config.dnsServers,
		config.rrTypes,
		config.edns,
		msg,
		&dns.Client{},
		&Result{
			res: make(map[interface{}][]string),
		},
	}
}

// Use DNS extension EDNS which can be used for DNS over UDP
func handleEDNS(msg *dns.Msg) *dns.Msg {
	opt := &dns.OPT{
		Hdr: dns.RR_Header{
			Name:   ".",
			Rrtype: dns.TypeOPT,
		},
	}
	opt.SetUDPSize(dns.DefaultMsgSize)
	msg.Extra = append(msg.Extra, opt)
	return msg
}

// Resolve queries every DNS-server with every specififed RR-type
func (r *Resolver) Resolve(domainname string) (*Result, error) {
	result := &Result{
		input: domainname,
	}

	for _, dnsServer := range r.dnsServers {
		log.Println("DNS-SERVER:", dnsServer)
		for _, rrType := range r.rrTypes {
			dnsResponseMsg, err := r.queryDNSServer(dnsServer, domainname, rrType, r.edns)
			if err != nil {
				continue
			}
			result.addMsg(dnsResponseMsg)
		}
	}

	return r.Result, nil
}

// Uses github.com/miekg/dns package, in order to invoke one query.
func (r *Resolver) queryDNSServer(dnsServer, domainname, rrType string, edns bool) (*dns.Msg, error) {
	fqdn := dns.Fqdn(domainname)
	r.dnsQueryMsg.Id = dns.Id()
	r.dnsQueryMsg.SetQuestion(fqdn, dns.StringToType[rrType])
	dnsServerSocket := dnsServer + ":" + DNSPORT
	dnsResponseMsg, err := dns.Exchange(r.dnsQueryMsg, dnsServerSocket)

	if err != nil {
		return nil, errors.New("dns.Exchange() failed")
	}

	if r.dnsQueryMsg.Id != dnsResponseMsg.Id {
		log.Printf("DNS msgID mismatch: Request-ID(%d), Response-ID(%d)", r.dnsQueryMsg.Id, dnsResponseMsg.Id)
		return nil, errors.New("DNS Msg-ID mismatch")
	}

	if dnsResponseMsg.MsgHdr.Truncated {
		if r.dnsClient.Net == "tcp" {
			return nil, errors.New("Received invalid truncated Msg over TCP") //fmt.Errorf("Got truncated message on tcp")
		}
		if edns {
			r.dnsClient.Net = "tcp"
		}
		return r.queryDNSServer(dnsServer, domainname, rrType, !edns)
	}

	return dnsResponseMsg, nil
}

// Result struct contains a map which stores DNS-responses.
type Result struct {
	input string
	res   map[interface{}][]string
}

// Adds DNS responses to a map. Most common RR-types are applied to retrieve
// proper information.
func (r *Result) addMsg(msg *dns.Msg) {
	if r.res == nil {
		r.res = make(map[interface{}][]string)
	}

	for _, rr := range msg.Answer {
		record := strings.Fields(rr.String())
		log.Println("Record:", record)
		switch rrType := rr.(type) {
		case *dns.A:
			r.res[rrType] = append(r.res[rrType], rr.(*dns.A).A.String())
		case *dns.AAAA:
			r.res[rrType] = append(r.res[rrType], rr.(*dns.AAAA).AAAA.String())
		case *dns.CNAME:
			r.res[rrType] = append(r.res[rrType], rr.(*dns.CNAME).Target)
		default:
			r.res[rrType] = append(r.res[rrType], rr.String())
		}
	}
}

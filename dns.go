package main

import (
	"github.com/miekg/dns"
	"log"
	"strconv"
	"strings"
)

func NewRR(s string) dns.RR { r, _ := dns.NewRR(s); return r }

func initNameServer(port int) {
	dns.HandleFunc(config.nsDomain, func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Authoritative = true
		name := r.Question[0].Name
		if strings.ToLower(name) == config.nsDomain {
			m.Answer = []dns.RR{}
			m.Ns = []dns.RR{
				&dns.SOA{
					Hdr:     dns.RR_Header{Name: config.nsDomain, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: config.nsRecordTTL},
					Ns:      config.nsRecordValue,
					Mbox:    config.nsRecordMbox,
					Serial:  2014042809,
					Refresh: 3600,
					Retry:   3600,
					Expire:  3600,
					Minttl:  config.nsRecordTTL,
				},
			}
		} else {
			m.Answer = []dns.RR{
				&dns.CNAME{
					Hdr:    dns.RR_Header{Name: name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: config.nsRecordTTL},
					Target: config.nsRecordCNAME,
				},
			}
			m.Ns = []dns.RR{
				&dns.NS{
					Hdr: dns.RR_Header{Name: config.nsDomain, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: config.nsRecordTTL},
					Ns:  config.nsRecordValue,
				},
			}
		}
		err := w.WriteMsg(m)
		if err != nil {
			log.Printf("%v while handling %s", err, name)
		}
	})
	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to listen UDP port %d for DNS server: %s.\n", port, err.Error())
		}
	}()
	go func() {
		srv := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "tcp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to listen TCP port %d for DNS server: %s.\n", port, err.Error())
		}
	}()
}

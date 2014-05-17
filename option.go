package main

import (
	"flag"
	"github.com/miekg/dns"
	"log"
	"os"
)

type Config struct {
	httpServerPort int
	nameServerPort int
	nsRecordCNAME  string
	nsRecordValue  string
	nsRecordTTL    uint32
	nsRecordMbox   string
	nsDomain       string
}

var (
	config Config
)

func getOptions() {
	var err error

	config = Config{}

	flag.IntVar(&config.httpServerPort, "http-port", 80, "http server port")
	flag.IntVar(&config.nameServerPort, "dns-port", 53, "dns server port")
	flag.StringVar(&config.nsRecordCNAME, "cname", "", "probe record cname to")
	flag.StringVar(&config.nsRecordValue, "ns", "ns-bench.edgeaccel.net.", "primary nameserver")
	flag.StringVar(&config.nsRecordMbox, "email", "dns-test.edgens.com.", "email")
	flag.StringVar(&config.nsDomain, "domain", "probe.edgeaccel.net.", "domain of probe record")

	flag.Parse()

	config.nsRecordTTL = 600

	if config.nsRecordCNAME == "" {
		config.nsRecordCNAME, err = os.Hostname()
		if err != nil {
			log.Fatalf("Unable to get hostname: %s.\n", err.Error())
		} else {
			if config.nsRecordCNAME == "" {
				log.Fatalf("Unable to get hostname: unknown error\n")
			} else {
				config.nsRecordCNAME = dns.Fqdn(config.nsRecordCNAME)
			}
			log.Printf("No CNAME option given, using hostname \"%s\".\n", config.nsRecordCNAME)
		}
	}

	checkFqdn := func(data string, message string) {
		if dns.IsFqdn(data) == false {
			log.Fatalf(message, data)
		}
	}

	checkFqdn(config.nsRecordCNAME, "Probe record CNAME \"%s\" isn't not a fully qualified domain name.\n")
	checkFqdn(config.nsRecordValue, "Primary Nameserver \"%s\" isn't not a fully qualified domain name.\n")
	checkFqdn(config.nsRecordMbox, "SOA Mailbox \"%s\" isn't not a fully qualified domain name.\n")
	checkFqdn(config.nsDomain, "Domain of probe record \"%s\" isn't not a fully qualified domain name.\n")
}

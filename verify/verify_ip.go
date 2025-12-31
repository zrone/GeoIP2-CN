package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

const defaultDataFile = "Country.mmdb"

func main() {
	var dbPath string
	var ipList string
	var expect string
	var allowMiss bool
	flag.StringVar(&dbPath, "db", defaultDataFile, "path to mmdb file")
	flag.StringVar(&ipList, "ip", "", "comma-separated list of IPs to verify")
	flag.StringVar(&expect, "expect", "CN", "expected ISO code (empty to skip assertion)")
	flag.BoolVar(&allowMiss, "allow-miss", false, "allow lookup miss (no record found)")
	flag.Parse()

	if strings.TrimSpace(ipList) == "" {
		log.Fatal("missing required flag: -ip")
	}

	db, err := geoip2.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, ipTxt := range strings.Split(ipList, ",") {
		ipTxt = strings.TrimSpace(ipTxt)
		if ipTxt == "" {
			continue
		}
		ip := net.ParseIP(ipTxt)
		if ip == nil {
			log.Fatalf("invalid ip: %s", ipTxt)
		}
		record, err := db.Country(ip)
		if err != nil || record == nil {
			if allowMiss {
				fmt.Printf("IP:%s-MISS\n", ipTxt)
				continue
			}
			log.Fatal(err)
		}

		iso := record.Country.IsoCode
		fmt.Printf("IP:%s-IsoCode:%s\n", ipTxt, iso)
		if expect != "" && iso != expect {
			log.Fatalf("unexpected ISO code for %s: got %q expect %q", ipTxt, iso, expect)
		}
	}
}

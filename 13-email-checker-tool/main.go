package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("domain, hasMx, hasSPF, sprRecord, hasDMARC, dmarcRecord")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func checkDomain(domain string) {
	var hasMx, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(mxRecords) > 0 {
		hasMx = true
	}

	spfRecords, err := net.LookupTXT(domain)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, record := range spfRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%s, %t, %t, %s, %t, %s\n", domain, hasMx, hasSPF, spfRecord, hasDMARC, dmarcRecord)

}

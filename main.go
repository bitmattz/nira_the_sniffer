// scan_ping_arp.go
package main

import (
	"fmt"
	port "nira_the_sniffer/services"
)

func main() {
	fmt.Println("PortScanning in Go")
	results := port.ScanPorts("localhost")

	fmt.Printf("Port scan results: %v\n", results)

}

// scan_ping_arp.go
package main

import (
	"fmt"

	portHandler "github.com/bitmattz/nira_the_sniffer/services"
	presenter "github.com/bitmattz/nira_the_sniffer/services"
)

func main() {
	fmt.Println("PortScanning in Go")
	presenter.StartApplicationPresenter()
	portHandler.ScanPorts("localhost")

}

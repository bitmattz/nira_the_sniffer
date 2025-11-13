// scan_ping_arp.go
package main

import (
	"fmt"

	// presenter "github.com/bitmattz/nira_the_sniffer/services/view"
	portHandler "github.com/bitmattz/nira_the_sniffer/services/ports"
)

func main() {
	fmt.Println("PortScanning in Go")
	// presenter.StartApplicationPresenter()
	portHandler.ScanPorts("192.168.1.49")
	//portHandler.ScanPorts("localhost")
}

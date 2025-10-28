package port

import (
	"net"
	"nira_the_sniffer/models"
	"strconv"
	"time"
)

func ScanPort(protocol, hostname string, port int) models.PortScan {
	result := models.PortScan{
		Port: port,
	}
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 60*time.Second)

	if err != nil {
		result.State = "closed"
		return result
	}

	defer conn.Close()
	result.State = "open"
	return result
}

func ScanPorts(hostname string) []models.PortScan {
	var results []models.PortScan

	for i := 1; i < 1024; i++ {
		portScanned := ScanPort("tcp", hostname, i)
		if portScanned.State == "open" {
			results = append(results, portScanned)
		}
	}

	return results
}

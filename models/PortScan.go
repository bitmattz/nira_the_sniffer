package models

type PortScan struct {
	Port   int    `json:"port"`
	State  string `json:"state"`
	Banner string `json:"service,omitempty"`
	PID    string `json:"pid,omitempty"`
	Status string `json:"status,omitempty"`
}

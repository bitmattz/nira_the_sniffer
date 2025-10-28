package models

import "time"

type Device struct {
	IP        string
	MAC       string
	Hostname  string
	Vendor    string
	FirstSeen time.Time
	LastSeen  time.Time
	SeenCount int
}

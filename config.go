package main

import (
	"strconv"
)

type Config struct {
	port        int
	ServingPath string
	dnsServers  []string
	rrTypes     []string
	edns        bool
	MIMEType    string
}

func NewConfig(port int, servingPath string, dnsServers []string, rrTypes []string, edns bool, MIMEType string) *Config {
	return &Config{
		port:        port,
		ServingPath: servingPath,
		dnsServers:  dnsServers,
		rrTypes:     rrTypes,
		edns:        edns,
		MIMEType:    MIMEType,
	}
}

// Converts numeric port into string format with prepended colon.
func (c *Config) GetPortString() string {
	if c.port != 0 {
		return ":" + strconv.Itoa(c.port)
	}
	return ""
}

package main

import (
	"fmt"
	"log"
	"net/http"
)

// Webserver struct holds basic parameters in order to work as DNS Resolver
type Webserver struct {
	Port     int
	Config   *Config
	Resolver *Resolver
}

// NewWebserver takes a *Config and returns an initialized Webserver struct
func NewWebserver(c *Config) *Webserver {
	return &Webserver{
		c.port,
		c,
		NewResolver(c),
	}
}

// Defines handling function and starts listening on port.
func (w *Webserver) startWebserver() {
	// TODO: enable or shut down TLS explicitly
	fmt.Println("Webserver listening on", w.Config.GetPortString())
	http.Handle(w.Config.ServingPath, w)
	err := http.ListenAndServe(w.Config.GetPortString(), nil)
	if err != nil {
		log.Println("Could not start webserver.")
	}
}

// ServeHTTP implements http.Handler interface.
// Is called when client requests page for s.config.ServingPath
func (w *Webserver) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	log.Println("Request received.")

	domainname := getDomainName(r)

	log.Println("all:", r.URL.String())
	log.Println("ARG:", domainname)

	// handle request
	_, err := w.Resolver.Resolve(domainname)
	if err != nil {
		fmt.Errorf("Resolve failed")
	}

	// set correct MIME-type in response to match JSON data
	wr.Header().Set("Content-Type", w.Config.MIMEType)

	log.Println("Response sent.\n-----")
}

// getDomainName is a helper function to extract the domain name from *http.Request
func getDomainName(r *http.Request) string {
	return r.URL.RawQuery
}

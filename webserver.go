package main

import (
	"fmt"
	"log"
	"net/http"
)

type Webserver struct {
	Port     int
	Config   *Config
	Resolver *Resolver
	// Address to listen on, ":dns" if empty.
	//Addr string
	// if "tcp" it will invoke a TCP listener, otherwise an UDP one.
	//Net string
	// TCP Listener to use, this is to aid in systemd's socket activation.
	//Listener net.Listener
	// UDP "Listener" to use, this is to aid in systemd's socket activation.
	//PacketConn net.PacketConn
	// Handler to invoke, dns.DefaultServeMux if nil.
	//Handler Handler
	//Server *http.Server
	// ...
}

/*type DNSHandler struct{}

func (*DNSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := getDomainName(r)

}*/

func getDomainName(r *http.Request) string {
	return r.URL.RawQuery
}

func NewWebserver(c *Config) *Webserver {
	return &Webserver{
		c.port,
		c,
		NewResolver(c),
		/*&http.Server{
			Addr:    c.GetPortString(),
			Handler: &DNSHandler{},
			//ReadTimeout:    10 * time.Second,
			//WriteTimeout:   10 * time.Second,
			//MaxHeaderBytes: 1 << 20,
		},*/
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
// Gets called when client requests page for s.config.ServingPath
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

	//byteResponse := []byte(response)

	// set correct MIME-type in response to match JSON data
	wr.Header().Set("Content-Type", w.Config.MIMEType)

	// send response
	//wr.Write(byteResponse)

	log.Println("Response sent.\n-----")
}

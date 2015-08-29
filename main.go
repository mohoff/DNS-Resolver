package main

import (
//_ "fmt"
)

func main() {

	c := NewConfig(
		8080,
		"/lookup",
		[]string{"8.8.8.8", "127.0.1.1"},
		[]string{"A", "AAAA", "MX", "CNAME", "PTR"}, //"
		false,
		"application/json",
	)

	w := NewWebserver(c)
	w.startWebserver()
}

package main

import (
	"log"

	"github.com/na4ma4/go-hostsfile"
)

func main() {
	cb := func(ip, h string) {
		log.Printf("'%s': '%s'", ip, h)
	}

	if err := hostsfile.ParseHostsFile("test/hostfile", cb); err != nil {
		log.Fatal(err)
	}
}

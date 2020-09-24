package main

import (
	"log"

	"github.com/na4ma4/go-hostsfile"
)

func main() {
	hosts := map[string][]string{}
	cb := func(ip, h string) {
		hosts[ip] = append(hosts[ip], h)
	}

	if err := hostsfile.ParseHostsFile("test/hostfile", cb); err != nil {
		log.Fatal(err)
	}

	for ip, hl := range hosts {
		for _, h := range hl {
			log.Printf("'%s': '%s'", ip, h)
		}
	}
}

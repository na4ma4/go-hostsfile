# go-hostsfile

[![PkgGoDev](https://pkg.go.dev/badge/github.com/na4ma4/go-hostsfile)](https://pkg.go.dev/github.com/na4ma4/go-hostsfile)
[![GoDoc](https://godoc.org/github.com/na4ma4/go-hostsfile/src/jwt?status.svg)](https://godoc.org/github.com/na4ma4/go-hostsfile)

(yet another) golang hostsfile parser this time using a callback function and hopefully very fast and memory efficient.

## Installation

```shell
go get -u github.com/na4ma4/go-hostsfile
```

## Example

```golang
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
```

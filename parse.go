package hostsfile

import (
	"bytes"
	"io"
	"os"
)

const (
	byteNewLine = byte('\n')
	byteReturn  = byte('\r')
	byteSpace   = byte(' ')
	byteTab     = byte('\t')
	byteComment = byte('#')
)

// CallbackFunc is used to process the data discovered in the hosts file.
type CallbackFunc func(ipAddr, host string)

// ParseHostsFile takes a host file path and a CallbackFunc and returns an error
// if the file can not be opened.
func ParseHostsFile(hostsFileName string, cb CallbackFunc) error {
	f, err := os.Open(hostsFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	return ParseHostsReader(f, cb)
}

// ParseHosts accepts a byte slice and sends all host/ip matches found to the CallbackFunc.
func ParseHosts(hostsFile []byte, cb CallbackFunc) error {
	return ParseHostsReader(bytes.NewBuffer(hostsFile), cb)
}

// ParseHostsReader parses host entries from an io.Reader.
func ParseHostsReader(r io.Reader, cb CallbackFunc) error {
	c := make([]byte, 1)
	buf := bytes.NewBuffer(nil)
	ipAddr := bytes.NewBuffer(nil)
	skipline := false

	for {
		if _, err := r.Read(c); err != nil {
			break
		}

		switch c[0] {
		case byteNewLine, byteReturn:
			if !skipline && ipAddr.Len() > 0 && buf.Len() > 0 {
				// ipAddr is set, process data before newline
				cb(ipAddr.String(), buf.String())
			}

			skipline = false

			ipAddr.Reset()
			buf.Reset()
		case byteComment:
			if !skipline && ipAddr.Len() > 0 && buf.Len() > 0 {
				// ipAddr is set, process data before newline
				cb(ipAddr.String(), buf.String())
			}

			skipline = true

			ipAddr.Reset()
			buf.Reset()
		case byteSpace, byteTab:
			switch {
			case skipline:
				continue
			case ipAddr.Len() == 0:
				// ipAddr is not set, set the ipAddr to the data
				ipAddr.Write(buf.Bytes())
			case buf.Len() > 0:
				// if lastpos and i are equal, it's zero length data, probably a series of spaces
				cb(ipAddr.String(), buf.String())
			}

			buf.Reset()
		default:
			buf.WriteByte(c[0])
		}
	}

	if buf.Len() > 0 && !skipline {
		// Last line wasn't a comment, so process remaining data
		cb(ipAddr.String(), buf.String())
	}

	return nil
}

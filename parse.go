package hostsfile

import (
	"io/ioutil"
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
	bs, err := ioutil.ReadFile(hostsFileName)
	if err != nil {
		return err
	}

	return ParseHosts(bs, cb)
}

// ParseHosts parses the bytes from a hosts file.
func ParseHosts(hostsFile []byte, cb CallbackFunc) error {
	lastpos := 0
	ipAddr := ""
	skipline := false

	for i := range hostsFile {
		switch hostsFile[i] {
		case byteNewLine, byteReturn:
			if skipline {
				skipline = false
				lastpos = i + 1

				continue
			}

			if ipAddr != "" && lastpos < i {
				// ipAddr is set, process data before newline
				cb(ipAddr, string(hostsFile[lastpos:i]))
			}

			ipAddr = ""
			lastpos = i + 1
		case byteComment:
			if skipline {
				continue
			} else if ipAddr != "" && lastpos < i {
				// ipAddr is set, process data before newline
				cb(ipAddr, string(hostsFile[lastpos:i]))
			}

			skipline = true
		case byteSpace, byteTab:
			switch {
			case skipline:
				continue
			case ipAddr == "":
				// ipAddr is not set, set the ipAddr to the data
				ipAddr = string(hostsFile[lastpos:i])
			case lastpos < i:
				// if lastpos and i are equal, it's zero length data, probably a series of spaces
				cb(ipAddr, string(hostsFile[lastpos:i]))
			}

			lastpos = i + 1
		}
	}

	if lastpos < len(hostsFile) && !skipline {
		// Last line wasn't a comment, so process remaining data
		cb(ipAddr, string(hostsFile[lastpos:]))
	}

	return nil
}

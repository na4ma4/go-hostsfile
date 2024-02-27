package hostsfile_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/na4ma4/go-hostsfile"
)

func getCallbackFunc(ip, hn *string) hostsfile.CallbackFunc {
	return func(ipAddr, hostName string) {
		*ip = ipAddr
		*hn = hostName
	}
}

func TestParser_File_Read(t *testing.T) {
	err := hostsfile.ParseHostsFile("testdata/hostfile", func(_, _ string) {})
	if err != nil {
		t.Errorf("hostsfile.ParseHostsFile(): error got '%s', want nil", err)
	}
}

func TestParser_ByteSlice_Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     []byte
		ip       string
		hostname string
	}{
		{"can handle data on last line", []byte(`127.0.0.3	localhostwithoutnewline`), "127.0.0.3", "localhostwithoutnewline"},
		{"can handle comments without spaces", []byte(`127.1.2.2	hostname2# comment on same line`), "127.1.2.2", "hostname2"},
		{"ignores comments", []byte("#127.1.2.2	hostname2"), "", ""},
		{"ignores multiple spaces instead of tabs", []byte("127.1.2.2    hostname2   "), "127.1.2.2", "hostname2"},
		{"can handle comments on same line", []byte(`127.1.2.1	hostname1 # comment on same line`), "127.1.2.1", "hostname1"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var ip, hn string
			cb := getCallbackFunc(&ip, &hn)
			err := hostsfile.ParseHosts(tt.data, cb)
			if err != nil {
				t.Errorf("hostsfile.ParseHosts(): error got '%s', want nil", err)
			}

			if ip != tt.ip {
				t.Errorf("hostsfile.ParseHosts(): IP Address got '%s', want '%s'", ip, tt.ip)
			}
			if hn != tt.hostname {
				t.Errorf("hostsfile.ParseHosts(): Hostname got '%s', want '%s'", hn, tt.hostname)
			}
		})
	}
}

func TestParser_ByteSlice_IPv6(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan string)
	cb := func(i, h string) {
		c <- fmt.Sprintf("%s:%s", i, h)
	}
	data := []byte(`::1     localhost ip6-localhost ip6-loopback`)

	errChan := make(chan error)

	go func() {
		errChan <- hostsfile.ParseHosts(data, cb)
	}()

	go func() {
		if line := <-c; line != "::1:localhost" {
			t.Errorf("hostsfile.ParseHosts(): IPv6 got '%s', want '%s'", line, "::1:localhost")
		}
		if line := <-c; line != "::1:ip6-localhost" {
			t.Errorf("hostsfile.ParseHosts(): IPv6 got '%s', want '%s'", line, "::1:ip6-localhost")
		}
		if line := <-c; line != "::1:ip6-loopback" {
			t.Errorf("hostsfile.ParseHosts(): IPv6 got '%s', want '%s'", line, "::1:ip6-loopback")
		}
		cancel()
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("hostsfile.ParseHosts(): parser error: %s", err)
				return
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
				t.Errorf("hostsfile.ParseHosts(): context error: %s", err)
			}
			return
		}
	}

	// It("can handle multiple names", func(done Done) {
	// 	c := make(chan string)
	// 	cb := func(i, h string) {
	// 		c <- fmt.Sprintf("%s:%s", i, h)
	// 	}
	// 	data := []byte("127.1.1.1	multiname1 multiname2 multiname3.withdomain.org")

	// 	go func() {
	// 		err := hostsfile.ParseHosts(data, cb)
	// 		Expect(err).NotTo(HaveOccurred())
	// 	}()

	// 	Expect(<-c).To(Equal("127.1.1.1:multiname1"))
	// 	Expect(<-c).To(Equal("127.1.1.1:multiname2"))
	// 	Expect(<-c).To(Equal("127.1.1.1:multiname3.withdomain.org"))
	// 	close(done)
	// }, 1)
}

func TestParser_ByteSlice_MultipleNames(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	c := make(chan string)
	cb := func(i, h string) {
		c <- fmt.Sprintf("%s:%s", i, h)
	}
	data := []byte(`127.1.1.1	multiname1 multiname2 multiname3.withdomain.org`)

	errChan := make(chan error)

	go func() {
		errChan <- hostsfile.ParseHosts(data, cb)
	}()

	go func() {
		if line := <-c; line != "127.1.1.1:multiname1" {
			t.Errorf("hostsfile.ParseHosts(): IPv6 got '%s', want '%s'", line, "127.1.1.1:multiname1")
		}
		if line := <-c; line != "127.1.1.1:multiname2" {
			t.Errorf("hostsfile.ParseHosts(): IPv6 got '%s', want '%s'", line, "127.1.1.1:multiname2")
		}
		if line := <-c; line != "127.1.1.1:multiname3.withdomain.org" {
			t.Errorf("hostsfile.ParseHosts(): IPv6 got '%s', want '%s'", line, "127.1.1.1:multiname3.withdomain.org")
		}
		cancel()
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				t.Errorf("hostsfile.ParseHosts(): parser error: %s", err)
				return
			}
		case <-ctx.Done():
			if err := ctx.Err(); err != nil && !errors.Is(err, context.Canceled) {
				t.Errorf("hostsfile.ParseHosts(): context error: %s", err)
			}
			return
		}
	}

	// It("can handle multiple names", func(done Done) {
	// 	c := make(chan string)
	// 	cb := func(i, h string) {
	// 		c <- fmt.Sprintf("%s:%s", i, h)
	// 	}
	// 	data := []byte("127.1.1.1	multiname1 multiname2 multiname3.withdomain.org")

	// 	go func() {
	// 		err := hostsfile.ParseHosts(data, cb)
	// 		Expect(err).NotTo(HaveOccurred())
	// 	}()

	// 	Expect(<-c).To(Equal("127.1.1.1:multiname1"))
	// 	Expect(<-c).To(Equal("127.1.1.1:multiname2"))
	// 	Expect(<-c).To(Equal("127.1.1.1:multiname3.withdomain.org"))
	// 	close(done)
	// }, 1)
}

// var _ = Describe("Parser", func() {
// 	// Describe("source=file", func() {
// 	// 	Describe("reading hosts file", func() {
// 	// 		It("is successful", func() {
// 	// 			err := hostsfile.ParseHostsFile("testdata/hostfile", func(h, i string) {})
// 	// 			Expect(err).NotTo(HaveOccurred())
// 	// 		})
// 	// 	})
// 	// })

// 	Describe("source=byteslice", func() {
// 		Describe("hosts file", func() {
// 			// It("can handle data on last line", func() {
// 			// 	var ip, hn string
// 			// 	cb := getCallbackFunc(&ip, &hn)
// 			// 	data := []byte(`127.0.0.3	localhostwithoutnewline`)

// 			// 	err := hostsfile.ParseHosts(data, cb)
// 			// 	Expect(err).NotTo(HaveOccurred())
// 			// 	Expect(ip).To(Equal("127.0.0.3"))
// 			// 	Expect(hn).To(Equal("localhostwithoutnewline"))
// 			// })

// 			// It("can handle comments without spaces", func() {
// 			// 	var ip, hn string
// 			// 	cb := getCallbackFunc(&ip, &hn)
// 			// 	data := []byte(`127.1.2.2	hostname2# comment on same line`)

// 			// 	err := hostsfile.ParseHosts(data, cb)
// 			// 	Expect(err).NotTo(HaveOccurred())
// 			// 	Expect(ip).To(Equal("127.1.2.2"))
// 			// 	Expect(hn).To(Equal("hostname2"))
// 			// })

// 			// It("ignores comments", func() {
// 			// 	var ip, hn string
// 			// 	cb := getCallbackFunc(&ip, &hn)
// 			// 	data := []byte("#127.1.2.2	hostname2")

// 			// 	err := hostsfile.ParseHosts(data, cb)
// 			// 	Expect(err).NotTo(HaveOccurred())
// 			// 	Expect(ip).To(Equal(""))
// 			// 	Expect(hn).To(Equal(""))
// 			// })

// 			// It("ignores multiple spaces instead of tabs", func() {
// 			// 	var ip, hn string
// 			// 	cb := getCallbackFunc(&ip, &hn)
// 			// 	data := []byte("127.1.2.2    hostname2   ")

// 			// 	err := hostsfile.ParseHosts(data, cb)
// 			// 	Expect(err).NotTo(HaveOccurred())
// 			// 	Expect(ip).To(Equal("127.1.2.2"))
// 			// 	Expect(hn).To(Equal("hostname2"))
// 			// })

// 			// It("can handle comments on same line", func() {
// 			// 	var ip, hn string
// 			// 	cb := getCallbackFunc(&ip, &hn)
// 			// 	data := []byte(`127.1.2.1	hostname1 # comment on same line`)

// 			// 	err := hostsfile.ParseHosts(data, cb)
// 			// 	Expect(err).NotTo(HaveOccurred())
// 			// 	Expect(ip).To(Equal("127.1.2.1"))
// 			// 	Expect(hn).To(Equal("hostname1"))
// 			// })

// 			It("can handle IPv6 hosts", func(done Done) {
// 				c := make(chan string)
// 				cb := func(i, h string) {
// 					c <- fmt.Sprintf("%s:%s", i, h)
// 				}
// 				data := []byte(`::1     localhost ip6-localhost ip6-loopback`)

// 				go func() {
// 					err := hostsfile.ParseHosts(data, cb)
// 					Expect(err).NotTo(HaveOccurred())
// 				}()

// 				Expect(<-c).To(Equal("::1:localhost"))
// 				Expect(<-c).To(Equal("::1:ip6-localhost"))
// 				Expect(<-c).To(Equal("::1:ip6-loopback"))
// 				close(done)
// 			}, 1)

// 			It("can handle multiple names", func(done Done) {
// 				c := make(chan string)
// 				cb := func(i, h string) {
// 					c <- fmt.Sprintf("%s:%s", i, h)
// 				}
// 				data := []byte("127.1.1.1	multiname1 multiname2 multiname3.withdomain.org")

// 				go func() {
// 					err := hostsfile.ParseHosts(data, cb)
// 					Expect(err).NotTo(HaveOccurred())
// 				}()

// 				Expect(<-c).To(Equal("127.1.1.1:multiname1"))
// 				Expect(<-c).To(Equal("127.1.1.1:multiname2"))
// 				Expect(<-c).To(Equal("127.1.1.1:multiname3.withdomain.org"))
// 				close(done)
// 			}, 1)
// 		})
// 	})
// })

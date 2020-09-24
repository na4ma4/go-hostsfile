package hostsfile_test

import (
	"fmt"

	"github.com/na4ma4/go-hostsfile"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func getCallbackFunc(ip, hn *string) hostsfile.CallbackFunc {
	return func(ipAddr, hostName string) {
		*ip = ipAddr
		*hn = hostName
	}
}

var _ = Describe("Parser", func() {
	Describe("source=file", func() {
		Describe("reading hosts file", func() {
			It("is successful", func() {
				err := hostsfile.ParseHostsFile("test/hostfile", func(h, i string) {})
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("source=byteslice", func() {
		Describe("hosts file", func() {
			It("can handle data on last line", func() {
				var ip, hn string
				cb := getCallbackFunc(&ip, &hn)
				data := []byte(`127.0.0.3	localhostwithoutnewline`)

				err := hostsfile.ParseHosts(data, cb)
				Expect(err).NotTo(HaveOccurred())
				Expect(ip).To(Equal("127.0.0.3"))
				Expect(hn).To(Equal("localhostwithoutnewline"))
			})

			It("can handle comments without spaces", func() {
				var ip, hn string
				cb := getCallbackFunc(&ip, &hn)
				data := []byte(`127.1.2.2	hostname2# comment on same line`)

				err := hostsfile.ParseHosts(data, cb)
				Expect(err).NotTo(HaveOccurred())
				Expect(ip).To(Equal("127.1.2.2"))
				Expect(hn).To(Equal("hostname2"))
			})

			It("ignores comments", func() {
				var ip, hn string
				cb := getCallbackFunc(&ip, &hn)
				data := []byte("#127.1.2.2	hostname2")

				err := hostsfile.ParseHosts(data, cb)
				Expect(err).NotTo(HaveOccurred())
				Expect(ip).To(Equal(""))
				Expect(hn).To(Equal(""))
			})

			It("ignores multiple spaces instead of tabs", func() {
				var ip, hn string
				cb := getCallbackFunc(&ip, &hn)
				data := []byte("127.1.2.2    hostname2   ")

				err := hostsfile.ParseHosts(data, cb)
				Expect(err).NotTo(HaveOccurred())
				Expect(ip).To(Equal("127.1.2.2"))
				Expect(hn).To(Equal("hostname2"))
			})

			It("can handle comments on same line", func() {
				var ip, hn string
				cb := getCallbackFunc(&ip, &hn)
				data := []byte(`127.1.2.1	hostname1 # comment on same line`)

				err := hostsfile.ParseHosts(data, cb)
				Expect(err).NotTo(HaveOccurred())
				Expect(ip).To(Equal("127.1.2.1"))
				Expect(hn).To(Equal("hostname1"))
			})

			It("can handle IPv6 hosts", func(done Done) {
				c := make(chan string)
				cb := func(i, h string) {
					c <- fmt.Sprintf("%s:%s", i, h)
				}
				data := []byte(`::1     localhost ip6-localhost ip6-loopback`)

				go func() {
					err := hostsfile.ParseHosts(data, cb)
					Expect(err).NotTo(HaveOccurred())
				}()

				Expect(<-c).To(Equal("::1:localhost"))
				Expect(<-c).To(Equal("::1:ip6-localhost"))
				Expect(<-c).To(Equal("::1:ip6-loopback"))
				close(done)
			}, 1)

			It("can handle multiple names", func(done Done) {
				c := make(chan string)
				cb := func(i, h string) {
					c <- fmt.Sprintf("%s:%s", i, h)
				}
				data := []byte("127.1.1.1	multiname1 multiname2 multiname3.withdomain.org")

				go func() {
					err := hostsfile.ParseHosts(data, cb)
					Expect(err).NotTo(HaveOccurred())
				}()

				Expect(<-c).To(Equal("127.1.1.1:multiname1"))
				Expect(<-c).To(Equal("127.1.1.1:multiname2"))
				Expect(<-c).To(Equal("127.1.1.1:multiname3.withdomain.org"))
				close(done)
			}, 1)
		})
	})
})

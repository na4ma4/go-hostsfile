package hostsfile_test

import (
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
	Describe("reading hosts file", func() {
		It("is successful", func() {
			err := hostsfile.ParseHostsFile("test/hostfile", func(h, i string) {})
			Expect(err).NotTo(HaveOccurred())
		})
	})

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

		It("can comments without spaces", func() {
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
	})
})

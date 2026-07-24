package dns

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Context("snapshotUpstreams", func() {
	writeResolvConf := func(text string) string {
		tmpFile, err := os.CreateTemp("", "")
		if err != nil {
			panic(err)
		}

		_, err = tmpFile.WriteString(text)
		if err != nil {
			panic(err)
		}

		return tmpFile.Name()
	}

	It("should return nameservers excluding opctl-managed entries & the node's own address", func() {
		/* arrange */
		path := writeResolvConf(`nameserver 127.0.0.1 # do not edit; managed by opctl
# comment line
search example.com
nameserver 127.0.0.53
nameserver 10.0.0.2
nameserver notanip
options edns0
`)

		/* act */
		upstreams := snapshotUpstreams(path, "127.0.0.1", "53")

		/* assert */
		Expect(upstreams).To(Equal([]string{"127.0.0.53:53", "10.0.0.2:53"}))
	})

	It("should exclude the node's own listen address when untagged", func() {
		/* arrange */
		path := writeResolvConf(`nameserver 127.0.0.1
nameserver 10.0.0.2
`)

		/* act */
		upstreams := snapshotUpstreams(path, "127.0.0.1", "53")

		/* assert */
		Expect(upstreams).To(Equal([]string{"10.0.0.2:53"}))
	})

	It("should exclude loopback nameservers given a wildcard listen address", func() {
		/* arrange */
		path := writeResolvConf(`nameserver 127.0.0.53
nameserver 10.0.0.2
`)

		/* act */
		upstreams := snapshotUpstreams(path, "0.0.0.0", "53")

		/* assert */
		Expect(upstreams).To(Equal([]string{"10.0.0.2:53"}))
	})

	It("should keep the node's own IP as upstream given the node listens on another port", func() {
		/* arrange */
		path := writeResolvConf(`nameserver 127.0.0.1
`)

		/* act */
		upstreams := snapshotUpstreams(path, "127.0.0.1", "5353")

		/* assert */
		Expect(upstreams).To(Equal([]string{"127.0.0.1:53"}))
	})

	It("should return no upstreams given resolv.conf isn't readable", func() {
		/* act */
		upstreams := snapshotUpstreams("/does/not/exist", "127.0.0.1", "53")

		/* assert */
		Expect(upstreams).To(BeEmpty())
	})
})

package client_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/saad-karim/saad-karim-parse-and-post/client"
)

var _ = Describe("client", func() {
	Context("push metadata", func() {
		var (
			c *client.Client
		)

		BeforeEach(func() {
			c = &client.Client{
				HTTPClient: &http.Client{},
				RetryLimit: 10,
			}
		})

		It("reads in csv file", func() {
			err := c.PushMetadata("../backup.csv")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

package api_e2e

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/flattrack/flattrack/src/backend/common"
)

var authToken string

var _ = Describe("API e2e tests", func() {
	apiServer := common.GetEnvOrDefault("API_SERVER", "http://localhost:8080")

	It("should reach root of API endpoint", func() {
		By("fetching from API")
		apiEndpoint := "api"
		resp, err := http.Get(fmt.Sprintf("%v/%v", apiServer, apiEndpoint))
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		Expect(common.GetHTTPresponseBodyContents(resp).Metadata.URL).To(Equal(fmt.Sprintf("/%v", apiEndpoint)))
	})

	It("should return hello message on hello endpoint", func() {
		By("fetching from API")
		apiEndpoint := "api/hello"
		resp, err := http.Get(fmt.Sprintf("%v/%v", apiServer, apiEndpoint))
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		response := common.GetHTTPresponseBodyContents(resp)
		Expect(response.Metadata.URL).To(Equal(fmt.Sprintf("/%v", apiEndpoint)))
		Expect(response.Metadata.Response).To(Equal("Hello"))
	})

})

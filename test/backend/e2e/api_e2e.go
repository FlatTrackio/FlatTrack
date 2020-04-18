package api_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/routes"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

/*
  note:
    due to the current lack of a client library and/or testing frame work, e2e test currently require an already set up instance
*/

var jwtToken string

var _ = Describe("API e2e tests", func() {
	apiServer := common.GetEnvOrDefault("APP_HOST", "http://localhost:8080")
	jwtToken = os.Getenv("APP_TEST_JWT")

	It("should reach root of API endpoint", func() {
		By("fetching from API's root")
		apiEndpoint := "api"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		Expect(routes.GetHTTPresponseBodyContents(resp).Metadata.URL).To(Equal(fmt.Sprintf("/%v", apiEndpoint)))
	})

	It("should be initialized", func() {
		By("querying the API")
		apiEndpoint := "api/system/initialized"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		Expect(routes.GetHTTPresponseBodyContents(resp).Data.(bool)).To(Equal(true), "instance should be initialized")
	})

	It("should have a flat name", func() {
		By("querying the API")
		apiEndpoint := "api/system/flatName"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).ToNot(Equal(""), "flatName should not be empty")
	})

	// TODO /api/admin/register - not yet possible in the same manner

	It("should have at least one user account", func() {
		By("querying the API")
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		response := routes.GetHTTPresponseBodyContents(resp)
		Expect(len(response.List.([]interface{})) > 0).To(Equal(true), "should at least one user account")
	})

	It("should return properties of a single user account", func() {
		By("listing all user accounts")
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		allUserAccountsResponse := routes.GetHTTPresponseBodyContents(resp).List //.([]interface{})//.([]types.UserSpec)
		allUserAccountsJSON, err := json.Marshal(allUserAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var allUserAccounts []types.UserSpec
		json.Unmarshal(allUserAccountsJSON, &allUserAccounts)
		firstUserAccount := allUserAccounts[0]

		By("listing all user accounts")
		apiEndpoint = "api/admin/users/" + firstUserAccount.Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)
		Expect(userAccount.Id).To(Equal(firstUserAccount.Id), "User account Id must equal user account Id from list")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
		Expect(userAccount.Names).To(Equal(firstUserAccount.Names), "User account names must equal user account names from list")
		Expect(userAccount.Names).ToNot(Equal(""), "User account names must not be empty")
		Expect(userAccount.Email).To(Equal(firstUserAccount.Email), "User account email must equal user account email from list")
		Expect(userAccount.Email).ToNot(Equal(""), "User account email must not be empty")
	})
})

func httpRequestWithHeader(verb string, url string, data interface{}) (resp *http.Response, err error) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer([]byte(fmt.Sprintf("%v", data))))
	req.Header.Set("Authorization", "bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	return resp, err
}

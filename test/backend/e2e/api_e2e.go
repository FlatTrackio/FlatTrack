package api_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/database"
	"gitlab.com/flattrack/flattrack/pkg/migrations"
	"gitlab.com/flattrack/flattrack/pkg/registration"
	"gitlab.com/flattrack/flattrack/pkg/routes"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

var jwtToken string

var _ = ginkgo.Describe("API e2e tests", func() {
	apiServer := common.GetEnvOrDefault("APP_HOST", "http://localhost:8080")
	apiServerAPIprefix := "api"
	jwtToken = os.Getenv("APP_TEST_JWT")

	cwd, _ := os.Getwd()
	os.Setenv("APP_DB_MIGRATIONS_PATH", fmt.Sprintf("%v/../../../migrations", cwd))

	regstrationForm := types.Registration{
		Timezone: "Pacific/Auckland",
		FlatName: "My flat",
		Language: "en_US",
		User: types.UserSpec{
			Names:    "Admin account",
			Email:    "adminaccount@example.com",
			Password: "Password123!",
			Groups:   []string{"flatmember", "admin"},
		},
	}

	// _ = godotenv.Load(".env")
	db, err := database.DB(common.GetDBusername(), common.GetDBpassword(), common.GetDBhost(), common.GetDBdatabase(), common.GetDBsslMode())
	if err != nil {
		log.Fatalln(err)
		return
	}

	ginkgo.BeforeSuite(func() {
		err = migrations.Reset(db)
		gomega.Expect(err).To(gomega.BeNil(), "failed to reset migrations")
		err = migrations.Migrate(db)
		gomega.Expect(err).To(gomega.BeNil(), "failed to migrate")

		registered, jwt, err := registration.Register(db, regstrationForm)

		gomega.Expect(err).To(gomega.BeNil(), "failed to register the instance")
		gomega.Expect(jwt).ToNot(gomega.Equal(""), "failed to register the instance")
		gomega.Expect(registered).To(gomega.Equal(true), "failed to register the instance")

		jwtToken = jwt
	})

	ginkgo.AfterSuite(func() {
		err = migrations.Reset(db)
		gomega.Expect(err).To(gomega.BeNil(), "failed to reset migrations")
		err = migrations.Migrate(db)
		gomega.Expect(err).To(gomega.BeNil(), "failed to migrate")
	})

	ginkgo.It("should reach root of API endpoint", func() {
		ginkgo.By("fetching from API's root")
		apiEndpoint := apiServerAPIprefix + ""
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		gomega.Expect(routes.GetHTTPresponseBodyContents(resp).Metadata.URL).To(gomega.Equal(fmt.Sprintf("/%v", apiEndpoint)))
	})

	ginkgo.It("should be initialized", func() {
		ginkgo.By("querying the API")
		apiEndpoint := apiServerAPIprefix + "/system/initialized"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		gomega.Expect(routes.GetHTTPresponseBodyContents(resp).Data.(bool)).To(gomega.Equal(true), "instance should be initialized")
	})

	ginkgo.It("should have a flat name", func() {
		ginkgo.By("querying the API")
		apiEndpoint := apiServerAPIprefix + "/system/flatName"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		gomega.Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).ToNot(gomega.Equal(""), "flatName should not be empty")
	})

	// TODO /api/admin/register - not yet possible in the same manner

	ginkgo.It("should have at least one user account", func() {
		ginkgo.By("querying the API")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		response := routes.GetHTTPresponseBodyContents(resp)
		gomega.Expect(len(response.List.([]interface{})) > 0).To(gomega.Equal(true), "should at least one user account")
	})

	ginkgo.It("should return properties of a single user account", func() {
		ginkgo.By("listing all user accounts")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		allUserAccountsResponse := routes.GetHTTPresponseBodyContents(resp).List
		allUserAccountsJSON, err := json.Marshal(allUserAccountsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var allUserAccounts []types.UserSpec
		json.Unmarshal(allUserAccountsJSON, &allUserAccounts)
		firstUserAccount := allUserAccounts[0]

		ginkgo.By("listing all user accounts")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + firstUserAccount.ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)
		gomega.Expect(userAccount.ID).To(gomega.Equal(firstUserAccount.ID), "User account ID must equal user account ID from list")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(firstUserAccount.Names), "User account names must equal user account names from list")
		gomega.Expect(userAccount.Names).ToNot(gomega.Equal(""), "User account names must not be empty")
		gomega.Expect(userAccount.Email).To(gomega.Equal(firstUserAccount.Email), "User account email must equal user account email from list")
		gomega.Expect(userAccount.Email).ToNot(gomega.Equal(""), "User account email must not be empty")
	})

	ginkgo.It("should allow the creation of valid user accounts", func() {
		accounts := []types.UserSpec{
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user2@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember", "admin"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user3@example.coop",
				Password: "Password123!",
				Groups:   []string{"flatmember", "admin"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "us.er1@example.coop",
				Password: "Password123!",
				Groups:   []string{"flatmember", "admin"},
			},
			{
				Names:       "Joe Bloggs",
				Email:       "us.er2@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "020 000 0000",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "us.er3@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "+64200000000",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "us.er4@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "64200000000",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "us.er5@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "64-20-000-000",
			},
			{
				Names:  "Joe Bloggs",
				Email:  "user4@example.com",
				Groups: []string{"flatmember"},
			},
		}
		for _, account := range accounts {
			accountBytes, err := json.Marshal(account)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			ginkgo.By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			response := routes.GetHTTPresponseBodyContents(resp)
			userAccountResponse := response.Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			ginkgo.By("checking the response")
			gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
			gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
			gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

			ginkgo.By("deleting the account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}
	})

	ginkgo.It("should disallow invalid fields for creating an account", func() {
		ginkgo.By("preparing accounts")
		dateNow := time.Now().Unix()
		accounts := []types.UserSpec{
			{
				Names:    "",
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "",
				Password: "Password123!",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "Password123!",
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "Password",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Groups:   []string{"flatmember"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{""},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember", "non existent group"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{"non existent group"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "test",
				Groups:   []string{"flatmember"},
			},
			{
				Names:       "Joe Bloggs",
				Email:       "user1@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
				PhoneNumber: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "user1@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
				PhoneNumber: "+64200000000x",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "user1@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
				PhoneNumber: "+64200000000",
				Birthday:    dateNow,
			},
		}

		for _, account := range accounts {
			accountBytes, err := json.Marshal(account)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			ginkgo.By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), apiServerAPIprefix+" should return error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusBadRequest), "api have return code of http.StatusBadRequest")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			ginkgo.By("checking the response")
			gomega.Expect(userAccount.ID).To(gomega.Equal(""), "User account ID must be empty")
			if userAccount.Names != "" {
				gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
			}
		}
	})

	ginkgo.It("should return user accounts if they exist", func() {
		accounts := []types.UserSpec{
			{
				Names:       "Joe Bloggs",
				Email:       "user1@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
				PhoneNumber: "020 000 0000",
			},
			{
				Names:       "Joe Bloggs 2",
				Email:       "user2@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
				PhoneNumber: "020 000 0000",
			},
			{
				Names:       "Joe Bloggs 3",
				Email:       "user3@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
				PhoneNumber: "020 000 0000",
			},
			{
				Names:       "Joe Bloggs 4",
				Email:       "user4@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
				PhoneNumber: "020 000 0000",
			},
		}
		for _, account := range accounts {
			accountBytes, err := json.Marshal(account)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			ginkgo.By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			ginkgo.By("checking the response")
			gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
			gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
			gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

			ginkgo.By("creating fetching user accounts")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			ginkgo.By("checking the response")
			gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
			gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
			gomega.Expect(userAccount.Email).To(gomega.Equal(account.Email), "User account email must match what was posted")
			gomega.Expect(userAccount.Groups).To(gomega.Equal(account.Groups), "User account email must match what was posted")
			gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

			ginkgo.By("deleting the account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}
	})

	ginkgo.It("should fail to get user account by id if account doesn't exist", func() {
		ids := []string{
			"a",
			"ab",
			"a-----c",
			"aasdsaddfgsdfgsaasd",
			"a-an32;8npwnsdvnsDNad",
			"albknlfknksln",
			"a=========",
			"a,,,,,,asdnasdasud",
		}
		for _, id := range ids {
			ginkgo.By("creating fetching user accounts")
			apiEndpoint := apiServerAPIprefix + "/admin/users/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	ginkgo.It("should allow updating of user accounts", func() {
		account := types.UserSpec{
			Names:  "Joe Bloggs",
			Email:  "user1@example.com",
			Groups: []string{"flatmember"},
		}

		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		response := routes.GetHTTPresponseBodyContents(resp)
		userAccountResponse := response.Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

		ginkgo.By("updating the resource")
		accountUpdates := []types.UserSpec{
			{
				Names:       "John Blogs",
				Email:       "user123@example.com",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "02000000",
				Birthday:    0,
				Password:    "Password1234!",
			},
			{
				Names:       "John 'golang' Smith",
				Email:       "user1237@example.com",
				Groups:      []string{"flatmember"},
				PhoneNumber: "0200000000",
				Birthday:    0420001,
				Password:    "Password123!",
			},
			{
				Names:       "E",
				Email:       "user1237@example.com",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "0200000001",
				Birthday:    0420002,
				Password:    "Password123!",
			},
		}
		for _, accountUpdate := range accountUpdates {
			ginkgo.By("updating account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			accountUpdateBytes, err := json.Marshal(accountUpdate)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountUpdateBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

			ginkgo.By("creating fetching the user account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			ginkgo.By("checking the account values")
			gomega.Expect(userAccount.Names).To(gomega.Equal(accountUpdate.Names), "user account names does not match update names")
			gomega.Expect(userAccount.Email).To(gomega.Equal(accountUpdate.Email), "user account email does not match update email")
			gomega.Expect(userAccount.Groups).To(gomega.Equal(accountUpdate.Groups), "user account groups does not match update groups")
			gomega.Expect(userAccount.PhoneNumber).To(gomega.Equal(accountUpdate.PhoneNumber), "user account phoneNumber does not match update phoneNumber")
			gomega.Expect(userAccount.Birthday).To(gomega.Equal(accountUpdate.Birthday), "user account birthday does not match update birthday")
		}

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should create account confirms when a password is not provided and allow it to be confirmed", func() {
		account := types.UserSpec{
			Names:  "Joe Bloggs",
			Email:  "user1@example.com",
			Groups: []string{"flatmember"},
		}

		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

		ginkgo.By("list user account confirms")
		apiEndpoint = apiServerAPIprefix + "/admin/useraccountconfirms"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v?userID=%v", apiServer, apiEndpoint, userAccount.ID), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		confirmsListResponse := routes.GetHTTPresponseBodyContents(resp).List
		confirmsListJSON, err := json.Marshal(confirmsListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var confirmsList []types.UserCreationSecretSpec
		json.Unmarshal(confirmsListJSON, &confirmsList)
		gomega.Expect(len(confirmsList) > 0).To(gomega.Equal(true), "must contain at least one confirm")

		ginkgo.By("fetching the user account confirm")
		apiEndpoint = apiServerAPIprefix + "/admin/useraccountconfirms/" + confirmsList[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		confirmResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		confirmJSON, err := json.Marshal(confirmResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var confirm types.UserCreationSecretSpec
		json.Unmarshal(confirmJSON, &confirm)
		gomega.Expect(confirm.ID).ToNot(gomega.Equal(""), "confirm id must not be empty")
		gomega.Expect(confirm.UserID).ToNot(gomega.Equal(""), "confirm userid must not be empty")
		gomega.Expect(confirm.Secret).ToNot(gomega.Equal(""), "confirm secret must not be empty")
		gomega.Expect(confirm.Valid).To(gomega.Equal(true), "confirm valid must be true")

		ginkgo.By("fetching the public route for user account confirm to check for it to be valid")
		apiEndpoint = apiServerAPIprefix + "/user/confirm/" + confirmsList[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		confirmValid := routes.GetHTTPresponseBodyContents(resp).Data
		gomega.Expect(confirmValid.(bool)).To(gomega.Equal(true), "confirm valid must be true")

		ginkgo.By("fetching the user account confirm")
		apiEndpoint = apiServerAPIprefix + "/user/confirm/" + confirm.ID + "?secret=" + confirm.Secret
		confirmUserAccount := types.UserSpec{
			Password: "Password123!",
		}
		confirmUserAccountJSON, err := json.Marshal(confirmUserAccount)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), confirmUserAccountJSON, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("fetching the user account confirm to check for it to be unavailable")
		apiEndpoint = apiServerAPIprefix + "/admin/useraccountconfirms/" + confirmsList[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound), "api have return code of http.StatusOK")

		ginkgo.By("fetching the user account to check if it's been registered")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err = json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccountRegistered types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccountRegistered)
		gomega.Expect(userAccountRegistered.ID).ToNot(gomega.Equal(""), "user account id must not be empty")
		gomega.Expect(userAccountRegistered.Registered).To(gomega.Equal(true), "account must be registered")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should disallow a deleted account to log in", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

		ginkgo.By("logging in")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginValid := routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		gomega.Expect(userAccountLoginValid).To(gomega.Equal(true), "JWT should be valid")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusUnauthorized), "api have return code of http.StatusUnauthorized")
		userAccountLoginValid = routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		gomega.Expect(userAccountLoginValid).To(gomega.Equal(false), "JWT should be valid")
	})

	ginkgo.It("should disallow a disabled account to log in", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

		ginkgo.By("logging in")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginValid := routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		gomega.Expect(userAccountLoginValid).To(gomega.Equal(true), "JWT should be valid")

		ginkgo.By("patching the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID + "/disabled"
		profilePatch := types.UserSpec{
			Disabled: true,
		}
		profilePatchData, err := json.Marshal(profilePatch)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		response := routes.GetHTTPresponseBodyContents(resp)
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse := response.Spec
		profileJSON, err := json.Marshal(profileResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		gomega.Expect(profile.Disabled).To(gomega.Equal(profilePatch.Disabled), "profile disabled does not match profilePatch disabled")

		ginkgo.By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusUnauthorized), "api have return code of http.StatusUnauthorized")
		userAccountLoginValid = routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		gomega.Expect(userAccountLoginValid).To(gomega.Equal(false), "JWT should be valid")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusUnauthorized), "api have return code of http.StatusUnauthorized")
		userAccountLoginValid = routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		gomega.Expect(userAccountLoginValid).To(gomega.Equal(false), "JWT should be valid")
	})

	ginkgo.It("should disallow login in after reset auth nonce", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

		ginkgo.By("logging in")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("resetting auth")
		apiEndpoint = apiServerAPIprefix + "/user/auth/reset"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusUnauthorized), "api have return code of http.StatusUnauthorized")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should disallow non-existent confirms", func() {
		confirmsIDs := []string{
			"a",
			"gggg",
			"239r2938rh",
			"fffffffffffp",
			"a48hrt894",
			"vsdvdvs",
			"0000000000000000",
		}
		ginkgo.By("fetching the user account confirm to check for it to be unavailable")
		for _, confirmID := range confirmsIDs {
			apiEndpoint := apiServerAPIprefix + "/admin/useraccountconfirms/" + confirmID
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	ginkgo.It("should return invalid for non-existent confirms on public route", func() {
		confirmsIDs := []string{
			"a",
			"gggg",
			"239r2938rh",
			"fffffffffffp",
			"a48hrt894",
			"vsdvdvs",
			"0000000000000000",
		}
		ginkgo.By("fetching the user account confirm to check for it to be unavailable")
		for _, confirmID := range confirmsIDs {
			apiEndpoint := apiServerAPIprefix + "/user/confirm/" + confirmID
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	ginkgo.It("should authenticate an existing user", func() {
		ginkgo.By("posting to auth")
		apiEndpoint := apiServerAPIprefix + "/user/auth"
		userAccountLogin := types.UserSpec{
			Email:    regstrationForm.User.Email,
			Password: regstrationForm.User.Password,
		}
		userAccountLoginData, err := json.Marshal(userAccountLogin)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), userAccountLoginData, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("checking validation of the token")
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), userAccountLoginData, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginValid := routes.GetHTTPresponseBodyContents(resp).Data
		gomega.Expect(userAccountLoginValid).To(gomega.Equal(true), "JWT should be valid")
	})

	ginkgo.It("should return the profile of the same account and allow patching of an account", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")
		account.ID = userAccount.ID

		ginkgo.By("getting a JWT")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("checking the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err := json.Marshal(profileResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		gomega.Expect(profile.ID).To(gomega.Equal(userAccount.ID), "profile id does not match account id")
		gomega.Expect(profile.Names).To(gomega.Equal(userAccount.Names), "profile names does not match account names")
		gomega.Expect(profile.Email).To(gomega.Equal(userAccount.Email), "profile email does not match account email")
		gomega.Expect(profile.PhoneNumber).To(gomega.Equal(userAccount.PhoneNumber), "profile phoneNumber does not match account phoneNumber")
		gomega.Expect(profile.Birthday).To(gomega.Equal(userAccount.Birthday), "profile birthday does not match account birthday")

		ginkgo.By("patching the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		profilePatch := types.UserSpec{
			ID:          "aaaaaaa",
			Names:       "Jonno bloggo",
			Email:       "user2@example.com",
			Password:    "Password1234!",
			PhoneNumber: "+64200000001",
			Birthday:    432001,
		}
		profilePatchData, err := json.Marshal(profilePatch)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err = json.Marshal(profileResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		profile = types.UserSpec{}
		json.Unmarshal(profileJSON, &profile)
		gomega.Expect(profile.ID).To(gomega.Equal(account.ID), "profile id does not match account id")
		gomega.Expect(profile.Names).To(gomega.Equal(profilePatch.Names), "profile names does not match profilePatch names")
		gomega.Expect(profile.Email).To(gomega.Equal(profilePatch.Email), "profile email does not match profilePatch email")
		gomega.Expect(profile.PhoneNumber).To(gomega.Equal(profilePatch.PhoneNumber), "profile phoneNumber does not match profilePatch phoneNumber")
		gomega.Expect(profile.Birthday).To(gomega.Equal(profilePatch.Birthday), "profile birthday does not match profilePatch birthday")

		ginkgo.By("getting a new JWT using new credentials")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should list all user accounts", func() {
		ginkgo.By("listing user accounts and checking the count")
		apiEndpoint := apiServerAPIprefix + "/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountsResponse := routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err := json.Marshal(userAccountsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccounts []types.UserSpec
		json.Unmarshal(userAccountsBytes, &userAccounts)

		ginkgo.By("checking the response")
		gomega.Expect(len(userAccounts)).To(gomega.Equal(1), "invalid amount of users")

		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account")
		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

		ginkgo.By("listing user accounts and checking the count")
		apiEndpoint = apiServerAPIprefix + "/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountsResponse = routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err = json.Marshal(userAccountsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		userAccounts = []types.UserSpec{}
		json.Unmarshal(userAccountsBytes, &userAccounts)

		ginkgo.By("checking the response")
		gomega.Expect(len(userAccounts)).To(gomega.Equal(2), "invalid amount of users")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("listing user accounts and checking the count")
		apiEndpoint = apiServerAPIprefix + "/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountsResponse = routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err = json.Marshal(userAccountsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		userAccounts = []types.UserSpec{}
		json.Unmarshal(userAccountsBytes, &userAccounts)

		ginkgo.By("checking the response")
		gomega.Expect(len(userAccounts)).To(gomega.Equal(1), "invalid amount of users")
	})

	ginkgo.It("should return a user by their id", func() {
		ginkgo.By("creating a user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")

		ginkgo.By("fetch the user account by the account's id")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err = json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		userAccount = types.UserSpec{}
		json.Unmarshal(userAccountBytes, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account Names must be equal to account Names")
		gomega.Expect(userAccount.Email).To(gomega.Equal(account.Email), "User account Email must be equal to account Email")
		gomega.Expect(userAccount.PhoneNumber).To(gomega.Equal(account.PhoneNumber), "User account PhoneNumber must be equal to account PhoneNumber")
		gomega.Expect(userAccount.Birthday).To(gomega.Equal(account.Birthday), "User account PhoneNumber must be equal to account PhoneNumber")
		gomega.Expect(userAccount.Groups).To(gomega.Equal(account.Groups), "User account Groups must be equal to account Groups")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should not allow two emails with the same email when patching a profile", func() {
		ginkgo.By("creating the first user account")
		account1 := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account1)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount1 types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount1)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount1.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")

		ginkgo.By("creating the first user account")
		account2 := types.UserSpec{
			Names:    "Bloggs Joe",
			Email:    "joeblogblog@example.com",
			Password: "Password123!",
			Groups:   []string{"flatmember"},
		}
		accountBytes, err = json.Marshal(account2)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err = json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount2 types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount2)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount2.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")

		ginkgo.By("logging in as the 2nd account")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("patching the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		profilePatch := types.UserSpec{
			Email: "user123@example.com",
		}
		profilePatchData, err := json.Marshal(profilePatch)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusBadRequest), "api have return code of http.StatusBadRequest")
		profileResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err := json.Marshal(profileResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		gomega.Expect(profile.Email).ToNot(gomega.Equal(profilePatch.Email), "profile email does not match profilePatch email")

		ginkgo.By("deleting the account1")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount1.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("deleting the account2")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount2.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should disallow own account from being deleted", func() {
		ginkgo.By("fetching user profile")
		apiEndpoint := apiServerAPIprefix + "/user/profile"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		userAccount := types.UserSpec{}
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusForbidden), "api have return code of http.StatusOK")
	})

	ginkgo.It("should not allow non admins to patch their groups", func() {
		ginkgo.By("creating the first user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")

		ginkgo.By("logging in as the account")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("patching the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		profilePatch := types.UserSpec{
			Groups: []string{"flatmember", "admin"},
		}
		profilePatchData, err := json.Marshal(profilePatch)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		response := routes.GetHTTPresponseBodyContents(resp)
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse := response.Spec
		profileJSON, err := json.Marshal(profileResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		gomega.Expect(profile.Groups).To(gomega.Equal(account.Groups), "profile email does not match profilePatch email")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should not allow non admins to update their groups", func() {
		ginkgo.By("creating the first user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		response := routes.GetHTTPresponseBodyContents(resp)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := response.Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")

		ginkgo.By("logging in as the account")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("updating the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		profilePatch := account
		profilePatch.Groups = []string{"flatmember", "admin"}
		profilePatchData, err := json.Marshal(profilePatch)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		response = routes.GetHTTPresponseBodyContents(resp)
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse := response.Spec
		profileJSON, err := json.Marshal(profileResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		gomega.Expect(profile.Groups).To(gomega.Equal(account.Groups), "profile email does not match profilePatch email")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should fail to list non-existent user accounts", func() {
		ids := []string{
			"aa",
			"234fasdsad",
			"----------asd",
			",,,,,,,,,,,,",
			"0-93jr9-q23nfiunq398nr3948n5q3c3",
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		}

		ginkgo.By("fetch the user account by the account's id")
		for _, id := range ids {
			apiEndpoint := apiServerAPIprefix + "/admin/users/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	ginkgo.It("should list groups and include the default groups", func() {
		apiEndpoint := apiServerAPIprefix + "/groups"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		groupsResponse := routes.GetHTTPresponseBodyContents(resp).List
		groupsBytes, err := json.Marshal(groupsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var groups []types.GroupSpec
		json.Unmarshal(groupsBytes, &groups)

		gomega.Expect(len(groups) >= 2).To(gomega.Equal(true), "There must be at least two groups")

		for _, groupItem := range groups {
			apiEndpoint := apiServerAPIprefix + "/groups/" + groupItem.ID
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			groupResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			groupBytes, err := json.Marshal(groupResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var group types.GroupSpec
			json.Unmarshal(groupBytes, &group)

			gomega.Expect(groupItem.ID).To(gomega.Equal(group.ID), "GroupItem ID must match Group ID")
			gomega.Expect(groupItem.Name).To(gomega.Equal(group.Name), "GroupItem Name must match Group Name")
			gomega.Expect(groupItem.DefaultGroup).To(gomega.Equal(group.DefaultGroup), "GroupItem DefaultGroup must match Group DefaultGroup")
			gomega.Expect(groupItem.Description).To(gomega.Equal(group.Description), "GroupItem Description must match Group Description")
		}
	})

	ginkgo.It("should fail to list groups which don't exist", func() {
		ids := []string{
			"aa",
			"234fasdsad",
			"----------asd",
			",,,,,,,,,,,,",
			"0-93jr9-q23nfiunq398nr3948n5q3c3",
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		}

		ginkgo.By("fetch the user account by the account's id")
		for _, id := range ids {
			apiEndpoint := apiServerAPIprefix + "/groups/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	ginkgo.It("should have the correct account groups for user account checking", func() {
		accounts := []types.UserSpec{
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember", "admin"},
			},
			{
				Names:    "Joe Bloggs 2",
				Email:    "user2@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember"},
			},
		}

		ginkgo.By("fetching all groups")
		apiEndpoint := apiServerAPIprefix + "/groups"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		groupsResponse := routes.GetHTTPresponseBodyContents(resp).List
		groupsBytes, err := json.Marshal(groupsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var groups []types.GroupSpec
		json.Unmarshal(groupsBytes, &groups)

		for _, account := range accounts {
			accountBytes, err := json.Marshal(account)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			ginkgo.By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			ginkgo.By("checking the response")
			gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")

			ginkgo.By("creating fetching user accounts")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			ginkgo.By("logging in as the new user account")
			apiEndpoint = apiServerAPIprefix + "/user/auth"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			jwt := routes.GetHTTPresponseBodyContents(resp).Data.(string)
			gomega.Expect(jwt).ToNot(gomega.Equal(""), "JWT in response must not be empty")

			defer func() {
				ginkgo.By("deleting the account")
				apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
				resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
				gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
				gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			}()

			for _, groupItem := range groups {
				expectGroup := false
				for _, accountGroup := range account.Groups {
					if groupItem.Name == accountGroup {
						expectGroup = true
					}
				}
				ginkgo.By("ensuring user account is or is not in a group")
				apiEndpoint := apiServerAPIprefix + "/user/can-i/group/" + groupItem.Name
				resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, jwt)
				gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
				gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
				canIgroupResponse := routes.GetHTTPresponseBodyContents(resp).Data.(bool)
				gomega.Expect(canIgroupResponse).To(gomega.Equal(expectGroup), "Group was expected for this user account", account.Names, account.Groups, groupItem.Name, expectGroup)
			}
		}
	})

	ginkgo.It("should disallow creating an account with an existing email", func() {
		ginkgo.By("creating a user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		ginkgo.By("creating the account again")
		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusBadRequest), "api have return code of http.StatusOK")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should create a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing all shopping lists")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListsBytes, err := json.Marshal(shoppingListsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingLists []types.ShoppingListSpec
		json.Unmarshal(shoppingListsBytes, &shoppingLists)

		gomega.Expect(len(shoppingLists)).To(gomega.Equal(1), "there must be one shopping list")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should patch a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		shoppingListPatches := []types.ShoppingListSpec{
			{
				Name: "Week 19",
			},
			{
				Name: "Week 19a",
			},
			{
				Name:  "My neat list",
				Notes: "Well, this list is neat.",
			},
		}
		for _, shoppingListPatch := range shoppingListPatches {
			shoppingListPatchBytes, err := json.Marshal(shoppingListPatch)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			ginkgo.By("creating a shopping list")
			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
			resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListPatchBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingListPatchedResponse := routes.GetHTTPresponseBodyContents(resp)
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			shoppingListPatchedResponseSpec := shoppingListPatchedResponse.Spec
			shoppingListPatchedBytes, err := json.Marshal(shoppingListPatchedResponseSpec)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingListPatched types.ShoppingListSpec
			json.Unmarshal(shoppingListPatchedBytes, &shoppingListPatched)

			gomega.Expect(shoppingListPatched.ID).To(gomega.Equal(shoppingListCreated.ID), "shopping list id must be equal to shopping list created id")
			gomega.Expect(shoppingListPatched.Name).To(gomega.Equal(shoppingListPatch.Name), "shopping list name does not match shopping list created name")
			gomega.Expect(shoppingListPatched.Notes).To(gomega.Equal(shoppingListPatch.Notes), "shopping list notes does not match shopping list created notes")
		}

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should update a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		shoppingListUpdates := []types.ShoppingListSpec{
			{
				Name: "Week 19",
			},
			{
				Name:      "Week 19a",
				Completed: true,
			},
			{
				Name:      "My neat list",
				Notes:     "Well, this list is neat.",
				Completed: false,
			},
		}
		for _, shoppingListUpdate := range shoppingListUpdates {
			shoppingListUpdateBytes, err := json.Marshal(shoppingListUpdate)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			ginkgo.By("creating a shopping list")
			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
			resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListUpdateBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingListUpdatedResponse := routes.GetHTTPresponseBodyContents(resp)
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			shoppingListUpdatedResponseSpec := shoppingListUpdatedResponse.Spec
			shoppingListUpdatedBytes, err := json.Marshal(shoppingListUpdatedResponseSpec)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingListUpdated types.ShoppingListSpec
			json.Unmarshal(shoppingListUpdatedBytes, &shoppingListUpdated)

			gomega.Expect(shoppingListUpdated.ID).To(gomega.Equal(shoppingListCreated.ID), "shopping list id must be equal to shopping list created id")
			gomega.Expect(shoppingListUpdated.Name).To(gomega.Equal(shoppingListUpdate.Name), "shopping list name does not match shopping list created name")
			gomega.Expect(shoppingListUpdated.Notes).To(gomega.Equal(shoppingListUpdate.Notes), "shopping list notes does not match shopping list created notes")
			gomega.Expect(shoppingListUpdated.Completed).To(gomega.Equal(shoppingListUpdate.Completed), "shopping list completed does not match shopping list created completed")
		}

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should not allow invalid shopping list properties", func() {
		shoppingLists := []types.ShoppingListSpec{
			{
				Name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			{
				Name: "",
			},
			{
				Name:  "My list",
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			{
				Name:       "My list",
				TemplateID: "x",
			},
		}

		for _, shoppingList := range shoppingLists {
			shoppingListBytes, err := json.Marshal(shoppingList)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			ginkgo.By("creating a shopping list")
			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusBadRequest), "api have return code of http.StatusOK")
			shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListBytes, err = json.Marshal(shoppingListResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingListCreated types.ShoppingListSpec
			json.Unmarshal(shoppingListBytes, &shoppingListCreated)

			gomega.Expect(shoppingListCreated.ID).To(gomega.Equal(""), "shopping list created id must not be empty")
		}
	})

	ginkgo.It("should allow adding items to a list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "Eggs",
				Quantity: 1,
			},
			{
				Name:     "Onions",
				Price:    2,
				Quantity: 1,
			},
			{
				Name:     "Pasta",
				Price:    0.8,
				Quantity: 3,
			},
			{
				Name:     "Bread",
				Price:    3.5,
				Quantity: 4,
				Notes:    "Sourdough",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges",
			},
		}

		ginkgo.By("listing shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		shoppingListItems = []types.ShoppingItemSpec{}
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(len(newShoppingListItems)), "There should be as many items added as on the shopping list")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should disallow adding items to a non-existent list", func() {
		ginkgo.By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "Eggs",
				Quantity: 1,
			},
			{
				Name:     "Onions",
				Price:    2,
				Quantity: 1,
			},
			{
				Name:     "Pasta",
				Price:    0.8,
				Quantity: 3,
			},
			{
				Name:     "Bread",
				Price:    3.5,
				Quantity: 4,
				Notes:    "Sourdough",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges",
			},
		}

		ginkgo.By("listing shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists/" + "xxxxxxxx" + "/items"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusNotFound), "api have return code of http.StatusBadRequest")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err := json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			gomega.Expect(shoppingItem.ID).To(gomega.Equal(""), "invalid shopping list items must not have ids")
		}
	})

	ginkgo.It("should not allow adding of invalid items to a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 0,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Tag:      "Fruits and veges",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
		}

		ginkgo.By("listing shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusBadRequest), "api must have return code of http.StatusBadRequest")
		}

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api must have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		shoppingListItems = []types.ShoppingItemSpec{}
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should allow updating of shopping list items", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating item on the list")
		newShoppingListItem := types.ShoppingItemSpec{
			Name:     "Lettuce",
			Price:    3,
			Quantity: 2,
			Notes:    "Not plastic bagged ones",
			Tag:      "Fruits and veges",
		}

		shoppingItemBytes, err := json.Marshal(newShoppingListItem)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingItem types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
		gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("updating item on the list")
		updatedShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "Iceberg lettuce",
				Price:    4,
				Quantity: 1,
				Notes:    "",
				Tag:      "Salad",
				Obtained: true,
			},
			{
				Name:     "Iceberg lettuce",
				Price:    4,
				Quantity: 1,
				Notes:    "",
				Tag:      "Salad",
				Obtained: false,
			},
		}

		for _, updatedShoppingListItem := range updatedShoppingListItems {
			shoppingItemBytes, err = json.Marshal(updatedShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items/" + shoppingItem.ID
			resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingItemResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemBytes, err := json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItemUpdated types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemBytes, &shoppingItemUpdated)
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

			gomega.Expect(shoppingItemUpdated.ID).ToNot(gomega.Equal(""), "shopping item id should not be nil")
			gomega.Expect(updatedShoppingListItem.Name).To(gomega.Equal(shoppingItemUpdated.Name), "shopping item name was not updated")
			gomega.Expect(updatedShoppingListItem.Price).To(gomega.Equal(shoppingItemUpdated.Price), "shopping item price was not updated")
			gomega.Expect(updatedShoppingListItem.Quantity).To(gomega.Equal(shoppingItemUpdated.Quantity), "shopping item quantity was not updated")
			gomega.Expect(updatedShoppingListItem.Notes).To(gomega.Equal(shoppingItemUpdated.Notes), "shopping item notes was not updated")
			gomega.Expect(updatedShoppingListItem.Tag).To(gomega.Equal(shoppingItemUpdated.Tag), "shopping item tag was not updated")
			gomega.Expect(updatedShoppingListItem.Obtained).To(gomega.Equal(shoppingItemUpdated.Obtained), "shopping item obtained was not updated")
		}

		ginkgo.By("deleting the shopping list item")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items/" + shoppingItem.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should allow patching of shopping list items", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating item on the list")
		newShoppingListItem := types.ShoppingItemSpec{
			Name:     "Lettuce",
			Price:    3,
			Quantity: 2,
			Notes:    "Not plastic bagged ones",
			Tag:      "Fruits and veges",
		}

		shoppingItemBytes, err := json.Marshal(newShoppingListItem)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingItem types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
		gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("updating item on the list")
		updatedShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "Iceberg lettuce",
				Price:    4,
				Quantity: 1,
				Notes:    "Just a note",
				Tag:      "Salad",
			},
			{
				Name:     "Iceberg lettuce",
				Price:    4,
				Quantity: 1,
				Notes:    "This note should have some useful meaning",
				Tag:      "Salad",
			},
		}

		for _, updatedShoppingListItem := range updatedShoppingListItems {
			shoppingItemBytes, err = json.Marshal(updatedShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items/" + shoppingItem.ID
			resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingItemResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemBytes, err := json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItemUpdated types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemBytes, &shoppingItemUpdated)
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

			gomega.Expect(updatedShoppingListItem.Name).To(gomega.Equal(shoppingItemUpdated.Name), "shopping item name was not updated")
			gomega.Expect(updatedShoppingListItem.Price).To(gomega.Equal(shoppingItemUpdated.Price), "shopping item price was not updated")
			gomega.Expect(updatedShoppingListItem.Quantity).To(gomega.Equal(shoppingItemUpdated.Quantity), "shopping item quantity was not updated")
			gomega.Expect(updatedShoppingListItem.Notes).To(gomega.Equal(shoppingItemUpdated.Notes), "shopping item notes was not updated")
			gomega.Expect(updatedShoppingListItem.Tag).To(gomega.Equal(shoppingItemUpdated.Tag), "shopping item tag was not updated")
		}

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should return a list of tags from a list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "Eggs",
				Quantity: 1,
				Tag:      "Dairy",
			},
			{
				Name:     "Onions",
				Price:    2,
				Quantity: 1,
				Tag:      "Fruits and veges",
			},
			{
				Name:     "Pasta",
				Price:    0.8,
				Quantity: 3,
				Tag:      "General",
			},
			{
				Name:     "Bread",
				Price:    3.5,
				Quantity: 4,
				Notes:    "Sourdough",
				Tag:      "General",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges",
			},
		}

		ginkgo.By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		ginkgo.By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err := json.Marshal(shoppingListTags)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var tags []string
		json.Unmarshal(shoppingListTagBytes, &tags)

		gomega.Expect(len(tags)).To(gomega.Equal(3), "invalid amount of tags")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should only return tags from a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating the first shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "Eggs",
				Quantity: 1,
				Tag:      "Dairy",
			},
			{
				Name:     "Onions",
				Price:    2,
				Quantity: 1,
				Tag:      "Fruits and veges",
			},
			{
				Name:     "Pasta",
				Price:    0.8,
				Quantity: 3,
				Tag:      "General",
			},
			{
				Name:     "Bread",
				Price:    3.5,
				Quantity: 4,
				Notes:    "Sourdough",
				Tag:      "General",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges",
			},
		}

		ginkgo.By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		ginkgo.By("creating the 2nd list")
		shoppingList2 := types.ShoppingListSpec{
			Name: "My list 2",
		}
		shoppingListBytes, err = json.Marshal(shoppingList2)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingList2Created types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingList2Created)

		gomega.Expect(shoppingList2Created.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingList2Created.Name).To(gomega.Equal(shoppingList2.Name), "shopping list name does not match shopping list created name2")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingList2Items []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingList2Items)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating items on the list")
		newShoppingList2Items := []types.ShoppingItemSpec{
			{
				Name:     "Eggs",
				Quantity: 1,
				Tag:      "Dairy1",
			},
			{
				Name:     "Oions",
				Price:    2,
				Quantity: 1,
				Tag:      "Fruits and veges1",
			},
			{
				Name:     "Pasta",
				Price:    0.8,
				Quantity: 3,
				Tag:      "General1",
			},
			{
				Name:     "Bread",
				Price:    3.5,
				Quantity: 4,
				Notes:    "Sourdough",
				Tag:      "General1",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges1",
			},
		}

		ginkgo.By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingList2Items {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingList2Created.ID), "shopping item must belong to a list")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		ginkgo.By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err := json.Marshal(shoppingListTags)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var listTags []string
		json.Unmarshal(shoppingListTagBytes, &listTags)

		gomega.Expect(len(listTags)).To(gomega.Equal(3), "invalid amount of tags")
		containsTagsFromOtherLists := false
		for _, tag := range listTags {
			for _, listItems := range newShoppingList2Items {
				if tag == listItems.Tag {
					containsTagsFromOtherLists = true
				}
			}
		}
		gomega.Expect(containsTagsFromOtherLists).To(gomega.Equal(false), "list of tags contains tags from other lists")

		ginkgo.By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err = json.Marshal(shoppingListTags)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var list2Tags []string
		json.Unmarshal(shoppingListTagBytes, &list2Tags)

		gomega.Expect(len(list2Tags)).To(gomega.Equal(3), "invalid amount of tags")
		containsTagsFromOtherLists = false
		for _, tag := range list2Tags {
			for _, listItems := range newShoppingListItems {
				if tag == listItems.Tag {
					containsTagsFromOtherLists = true
					break
				}
			}
		}
		gomega.Expect(containsTagsFromOtherLists).To(gomega.Equal(false), "list of tags contains tags from other lists")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("deleting the 2nd shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should allow updating of tags in a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating items on the list")
		newShoppingListItem := types.ShoppingItemSpec{
			Name:     "Lettuce",
			Price:    3,
			Quantity: 2,
			Notes:    "Not plastic bagged ones",
			Tag:      "Fruits and veges",
		}

		ginkgo.By("creating shopping list items")
		shoppingListItemBytes, err := json.Marshal(newShoppingListItem)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListItemBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingItem types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
		gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		shoppingItemTagUpdate := types.ShoppingTag{
			Name: "Veges",
		}
		ginkgo.By("patching the tag")
		shoppingItemTagBytes, err := json.Marshal(shoppingItemTagUpdate)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags/Fruits%20and%20veges"
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemTagBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err := json.Marshal(shoppingListTags)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var listTags []string
		json.Unmarshal(shoppingListTagBytes, &listTags)

		var foundUpdatedTag bool
		for _, tag := range listTags {
			if tag == shoppingItemTagUpdate.Name {
				foundUpdatedTag = true
			}
		}
		gomega.Expect(foundUpdatedTag).To(gomega.Equal(true), "Unable to find updated tag")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should allow management of shopping tags", func() {
		shoppingTags := []types.ShoppingTag{
			{
				Name: "Fruits and Veges",
			},
			{
				Name: "Meat",
			},
			{
				Name: "Poultry",
			},
		}

		ginkgo.By("creating shopping tags")
		for _, shoppingTag := range shoppingTags {
			shoppingTagBytes, err := json.Marshal(shoppingTag)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/tags"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingTagBytes, "")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingTagResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingTagBytes, err = json.Marshal(shoppingTagResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingTag types.ShoppingTag
			json.Unmarshal(shoppingTagBytes, &shoppingTag)
			gomega.Expect(shoppingTag.ID).ToNot(gomega.Equal(""), "shopping tag must have an ID")
		}

		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/tags"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingTagsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingTagBytes, err := json.Marshal(shoppingTagsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var tags []types.ShoppingTag
		json.Unmarshal(shoppingTagBytes, &tags)

		gomega.Expect(len(tags)).To(gomega.Equal(len(shoppingTags)), "failed to find the correct number (%v) of shopping tags in length of list (%v)", len(shoppingTags), len(tags))
		foundTags := 0
		for _, tag := range tags {
			for _, expectedTag := range shoppingTags {
				if expectedTag.Name == tag.Name {
					foundTags++
				}
			}
		}
		gomega.Expect(foundTags).To(gomega.Equal(len(shoppingTags)), "failed to find the correct (%v) number of shopping tags in list of tags from response (%v)", len(shoppingTags), foundTags)

		// update tag name
		tagUpdate := types.ShoppingTag{
			Name: "Fruits, Veges, and Fresh",
		}
		shoppingTagUpdateBytes, err := json.Marshal(tagUpdate)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/tags/" + tags[0].ID
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingTagUpdateBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		// get tag
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/tags/" + tags[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		shoppingTagUpdateGetResponse := routes.GetHTTPresponseBodyContents(resp)
		shoppingTagBytes, err = json.Marshal(shoppingTagUpdateGetResponse.Spec)
		var shoppingTagUpdated types.ShoppingTag
		json.Unmarshal(shoppingTagBytes, &shoppingTagUpdated)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK", shoppingTagUpdateGetResponse.Metadata.Response)
		gomega.Expect(shoppingTagUpdated.ID).To(gomega.Equal(tags[0].ID), "shopping tag must have an ID (%v) matching it's previous ID (%v)", shoppingTagUpdated.ID, tags[0].ID)
		gomega.Expect(shoppingTagUpdated.Name).To(gomega.Equal(tagUpdate.Name), "shopping tag must have the new tag name")

		for _, tag := range tags {
			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/tags/" + tag.ID
			resp, err := httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

	})

	ginkgo.It("should allow templating of a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		gomega.Expect(len(shoppingListItems)).To(gomega.Equal(0), "There should be no items on the shopping list")

		ginkgo.By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name:     "Eggs",
				Quantity: 1,
				Tag:      "Dairy",
			},
			{
				Name:     "Onions",
				Price:    2,
				Quantity: 1,
				Tag:      "Fruits and veges",
			},
			{
				Name:     "Pasta",
				Price:    0.8,
				Quantity: 3,
				Tag:      "General",
			},
			{
				Name:     "Bread",
				Price:    3.5,
				Quantity: 4,
				Notes:    "Sourdough",
				Tag:      "General",
			},
			{
				Name:     "Lettuce",
				Price:    3,
				Quantity: 2,
				Notes:    "Not plastic bagged ones",
				Tag:      "Fruits and veges",
			},
		}

		ginkgo.By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		shoppingListFromTemplate := types.ShoppingListSpec{
			Name:       "My list (from template)",
			Notes:      "This is a templated list",
			TemplateID: shoppingListCreated.ID,
		}
		shoppingListBytes, err = json.Marshal(shoppingListFromTemplate)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a templated shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		templatedShoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		templatedShoppingListBytes, err := json.Marshal(templatedShoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListTemplatedCreated types.ShoppingListSpec
		json.Unmarshal(templatedShoppingListBytes, &shoppingListTemplatedCreated)

		ginkgo.By("listing items of the templated shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListTemplatedCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		shoppingListItems = []types.ShoppingItemSpec{}
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)
		gomega.Expect(len(newShoppingListItems)).To(gomega.Equal(len(shoppingListItems)), "templated list must have the same amount of items as the orignal list")

		foundTotal := 0
		for _, item := range newShoppingListItems {
			for _, templatedItem := range shoppingListItems {
				if templatedItem.Name == item.Name {
					foundTotal++
					continue
				}
			}
		}

		gomega.Expect(foundTotal).To(gomega.Equal(len(newShoppingListItems)), "unable to find all items from original list in templated list")

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api must have return code of http.StatusOK")

		ginkgo.By("deleting the template shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListTemplatedCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api must have return code of http.StatusOK")
	})

	ginkgo.It("should retain template origin shopping list id in items and list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		gomega.Expect(shoppingListCreated.ID).ToNot(gomega.Equal(""), "shopping list created id must not be empty")
		gomega.Expect(shoppingListCreated.Name).To(gomega.Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		ginkgo.By("creating an item on the origin list")
		newShoppingListItem := types.ShoppingItemSpec{
			Name:     "Lettuce",
			Price:    3,
			Quantity: 2,
			Notes:    "Not plastic bagged ones",
			Tag:      "Fruits and veges",
		}

		shoppingListItemBytes, err := json.Marshal(newShoppingListItem)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListItemBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListItemsBytes, err := json.Marshal(shoppingItemResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingItem types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
		gomega.Expect(shoppingItem.ListID).To(gomega.Equal(shoppingListCreated.ID), "shopping item must belong to a list")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		shoppingListFromTemplate := types.ShoppingListSpec{
			Name:       "My list (from template)",
			Notes:      "This is a templated list",
			TemplateID: shoppingListCreated.ID,
		}
		shoppingListBytes, err = json.Marshal(shoppingListFromTemplate)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a templated shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		templatedShoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		templatedShoppingListBytes, err := json.Marshal(templatedShoppingListResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var shoppingListTemplatedCreated types.ShoppingListSpec
		json.Unmarshal(templatedShoppingListBytes, &shoppingListTemplatedCreated)

		gomega.Expect(shoppingListTemplatedCreated.TemplateID).To(gomega.Equal(shoppingListCreated.ID), "templated list must have templateID field matching origin ID")

		ginkgo.By("listing items of the templated shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListTemplatedCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		shoppingListItems := []types.ShoppingItemSpec{}
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		for _, item := range shoppingListItems {
			gomega.Expect(item.TemplateID).To(gomega.Equal(shoppingListCreated.ID), "templated list item must have templateID field matching origin ID")
		}

		ginkgo.By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api must have return code of http.StatusOK")

		ginkgo.By("deleting the template shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListTemplatedCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api must have return code of http.StatusOK")
	})

	ginkgo.It("should require authorization for protected routes", func() {
		apiEndpoint := apiServer + "/" + apiServerAPIprefix + "/user/profile"
		req, err := http.NewRequest("GET", apiEndpoint, nil)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusUnauthorized), "endpoint should be restricted")
		requestResp := routes.GetHTTPresponseBodyContents(resp)
		gomega.Expect(requestResp.Metadata.Response).To(gomega.Equal("Unauthorized"), "")
	})

	ginkgo.It("should require admin for admin protected routes", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		ginkgo.By("creating a user account with no admin access")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		ginkgo.By("checking the response")
		gomega.Expect(userAccount.ID).ToNot(gomega.Equal(""), "User account ID must not be empty")
		gomega.Expect(userAccount.Names).To(gomega.Equal(account.Names), "User account names must match what was posted")
		gomega.Expect(userAccount.Password).To(gomega.Equal(""), "User account password must return an empty string")

		ginkgo.By("logging in")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		gomega.Expect(userAccountLoginResponseData).ToNot(gomega.Equal(""), "JWT in response must not be empty")

		ginkgo.By("trying to use an admin route")
		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, userAccountLoginResponseData)
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusForbidden), "api have return code of http.StatusForbidden")

		ginkgo.By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should allow configuration of shopping list notes", func() {
		ginkgo.By("fetching the notes")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/settings/notes"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		gomega.Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).To(gomega.Equal(""), "notes should be empty")

		ginkgo.By("updating the notes")
		notesUpdate := types.ShoppingListNotes{
			Notes: "Our budget is $200. Please go to the closest supermarket",
		}
		notesUpdateBytes, err := json.Marshal(notesUpdate)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/admin/settings/shoppingListNotes"
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), notesUpdateBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")

		ginkgo.By("fetching the notes")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/settings/notes"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
		gomega.Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).To(gomega.Equal(notesUpdate.Notes), "notes should be empty")

		ginkgo.By("resetting the notes")
		notesUpdate = types.ShoppingListNotes{
			Notes: "",
		}
		notesUpdateBytes, err = json.Marshal(notesUpdate)
		gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/admin/settings/shoppingListNotes"
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), notesUpdateBytes, "")
		gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
		gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	ginkgo.It("should not allow invalid shopping list notes", func() {
		notesUpdates := []types.ShoppingListNotes{
			{
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
		}

		for _, notesUpdate := range notesUpdates {
			ginkgo.By("updating the notes to " + notesUpdate.Notes)
			notesUpdateBytes, err := json.Marshal(notesUpdate)
			gomega.Expect(err).To(gomega.BeNil(), "failed to marshal to JSON")

			apiEndpoint := apiServerAPIprefix + "/admin/settings/shoppingListNotes"
			resp, err := httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), notesUpdateBytes, "")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusBadRequest), "api have return code of http.StatusOK")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).To(gomega.Equal(""), "notes should be empty")

			ginkgo.By("fetching the notes")
			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/settings/notes"
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			gomega.Expect(err).To(gomega.BeNil(), "Request should not return an error")
			gomega.Expect(resp.StatusCode).To(gomega.Equal(http.StatusOK), "api have return code of http.StatusOK")
			gomega.Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).To(gomega.Equal(""), "notes should be empty")
		}
	})
})

func httpRequestWithHeader(verb string, url string, data []byte, jwt string) (resp *http.Response, err error) {
	if jwt == "" {
		jwt = jwtToken
	}
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	return resp, err
}

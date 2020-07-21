package api_e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/migrations"
	"gitlab.com/flattrack/flattrack/src/backend/registration"
	"gitlab.com/flattrack/flattrack/src/backend/routes"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

var jwtToken string

var _ = Describe("API e2e tests", func() {
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
	db, err := database.DB(common.GetDBusername(), common.GetDBpassword(), common.GetDBhost(), common.GetDBdatabase())
	if err != nil {
		log.Fatalln(err)
		return
	}

	BeforeSuite(func() {
		err = migrations.Reset(db)
		Expect(err).To(BeNil(), "failed to reset migrations")
		err = migrations.Migrate(db)
		Expect(err).To(BeNil(), "failed to migrate")

		registered, jwt, err := registration.Register(db, regstrationForm)

		Expect(err).To(BeNil(), "failed to register the instance")
		Expect(jwt).ToNot(Equal(""), "failed to register the instance")
		Expect(registered).To(Equal(true), "failed to register the instance")

		jwtToken = jwt
	})

	AfterSuite(func() {
		err = migrations.Reset(db)
		Expect(err).To(BeNil(), "failed to reset migrations")
		err = migrations.Migrate(db)
		Expect(err).To(BeNil(), "failed to migrate")
	})

	It("should reach root of API endpoint", func() {
		By("fetching from API's root")
		apiEndpoint := apiServerAPIprefix + ""
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		Expect(routes.GetHTTPresponseBodyContents(resp).Metadata.URL).To(Equal(fmt.Sprintf("/%v", apiEndpoint)))
	})

	It("should be initialized", func() {
		By("querying the API")
		apiEndpoint := apiServerAPIprefix + "/system/initialized"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		Expect(routes.GetHTTPresponseBodyContents(resp).Data.(bool)).To(Equal(true), "instance should be initialized")
	})

	It("should have a flat name", func() {
		By("querying the API")
		apiEndpoint := apiServerAPIprefix + "/system/flatName"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).ToNot(Equal(""), "flatName should not be empty")
	})

	// TODO /api/admin/register - not yet possible in the same manner

	It("should have at least one user account", func() {
		By("querying the API")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		response := routes.GetHTTPresponseBodyContents(resp)
		Expect(len(response.List.([]interface{})) > 0).To(Equal(true), "should at least one user account")
	})

	It("should return properties of a single user account", func() {
		By("listing all user accounts")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		allUserAccountsResponse := routes.GetHTTPresponseBodyContents(resp).List
		allUserAccountsJSON, err := json.Marshal(allUserAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var allUserAccounts []types.UserSpec
		json.Unmarshal(allUserAccountsJSON, &allUserAccounts)
		firstUserAccount := allUserAccounts[0]

		By("listing all user accounts")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + firstUserAccount.ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)
		Expect(userAccount.ID).To(Equal(firstUserAccount.ID), "User account ID must equal user account ID from list")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(firstUserAccount.Names), "User account names must equal user account names from list")
		Expect(userAccount.Names).ToNot(Equal(""), "User account names must not be empty")
		Expect(userAccount.Email).To(Equal(firstUserAccount.Email), "User account email must equal user account email from list")
		Expect(userAccount.Email).ToNot(Equal(""), "User account email must not be empty")
	})

	It("should allow the creation of valid user accounts", func() {
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
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			response := routes.GetHTTPresponseBodyContents(resp)
			userAccountResponse := response.Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("deleting the account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}
	})

	It("should disallow invalid fields for creating an account", func() {
		By("preparing accounts")
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
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), apiServerAPIprefix + " should return error")
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest), "api have return code of http.StatusBadRequest")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.ID).To(Equal(""), "User account ID must be empty")
			if userAccount.Names != "" {
				Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			}
		}
	})

	It("should return user accounts if they exist", func() {
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
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("creating fetching user accounts")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Email).To(Equal(account.Email), "User account email must match what was posted")
			Expect(userAccount.Groups).To(Equal(account.Groups), "User account email must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("deleting the account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}
	})

	It("should fail to get user account by id if account doesn't exist", func() {
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
			By("creating fetching user accounts")
			apiEndpoint := apiServerAPIprefix + "/admin/users/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	It("should allow updating of user accounts", func() {
		account := types.UserSpec{
			Names:  "Joe Bloggs",
			Email:  "user1@example.com",
			Groups: []string{"flatmember"},
		}

		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		response := routes.GetHTTPresponseBodyContents(resp)
		userAccountResponse := response.Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("updating the resource")
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
			By("updating account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			accountUpdateBytes, err := json.Marshal(accountUpdate)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountUpdateBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

			By("creating fetching the user account")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the account values")
			Expect(userAccount.Names).To(Equal(accountUpdate.Names), "user account names does not match update names")
			Expect(userAccount.Email).To(Equal(accountUpdate.Email), "user account email does not match update email")
			Expect(userAccount.Groups).To(Equal(accountUpdate.Groups), "user account groups does not match update groups")
			Expect(userAccount.PhoneNumber).To(Equal(accountUpdate.PhoneNumber), "user account phoneNumber does not match update phoneNumber")
			Expect(userAccount.Birthday).To(Equal(accountUpdate.Birthday), "user account birthday does not match update birthday")
		}

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should create account confirms when a password is not provided and allow it to be confirmed", func() {
		account := types.UserSpec{
			Names:  "Joe Bloggs",
			Email:  "user1@example.com",
			Groups: []string{"flatmember"},
		}

		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("list user account confirms")
		apiEndpoint = apiServerAPIprefix + "/admin/useraccountconfirms"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v?userID=%v", apiServer, apiEndpoint, userAccount.ID), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		confirmsListResponse := routes.GetHTTPresponseBodyContents(resp).List
		confirmsListJSON, err := json.Marshal(confirmsListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var confirmsList []types.UserCreationSecretSpec
		json.Unmarshal(confirmsListJSON, &confirmsList)
		Expect(len(confirmsList) > 0).To(Equal(true), "must contain at least one confirm")

		By("fetching the user account confirm")
		apiEndpoint = apiServerAPIprefix + "/admin/useraccountconfirms/" + confirmsList[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		confirmResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		confirmJSON, err := json.Marshal(confirmResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var confirm types.UserCreationSecretSpec
		json.Unmarshal(confirmJSON, &confirm)
		Expect(confirm.ID).ToNot(Equal(""), "confirm id must not be empty")
		Expect(confirm.UserID).ToNot(Equal(""), "confirm userid must not be empty")
		Expect(confirm.Secret).ToNot(Equal(""), "confirm secret must not be empty")
		Expect(confirm.Valid).To(Equal(true), "confirm valid must be true")

		By("fetching the public route for user account confirm to check for it to be valid")
		apiEndpoint = apiServerAPIprefix + "/user/confirm/" + confirmsList[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		confirmValid := routes.GetHTTPresponseBodyContents(resp).Data
		Expect(confirmValid.(bool)).To(Equal(true), "confirm valid must be true")

		By("fetching the user account confirm")
		apiEndpoint = apiServerAPIprefix + "/user/confirm/" + confirm.ID + "?secret=" + confirm.Secret
		confirmUserAccount := types.UserSpec{
			Password: "Password123!",
		}
		confirmUserAccountJSON, err := json.Marshal(confirmUserAccount)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), confirmUserAccountJSON, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("fetching the user account confirm to check for it to be unavailable")
		apiEndpoint = apiServerAPIprefix + "/admin/useraccountconfirms/" + confirmsList[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound), "api have return code of http.StatusOK")

		By("fetching the user account to check if it's been registered")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err = json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccountRegistered types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccountRegistered)
		Expect(userAccountRegistered.ID).ToNot(Equal(""), "user account id must not be empty")
		Expect(userAccountRegistered.Registered).To(Equal(true), "account must be registered")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should disallow a deleted account to log in", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("logging in")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginValid := routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		Expect(userAccountLoginValid).To(Equal(true), "JWT should be valid")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized), "api have return code of http.StatusUnauthorized")
		userAccountLoginValid = routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		Expect(userAccountLoginValid).To(Equal(false), "JWT should be valid")
	})

	It("should disallow login in after reset auth nonce", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("logging in")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("resetting auth")
		apiEndpoint = apiServerAPIprefix + "/user/auth/reset"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("checking validation of the token")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized), "api have return code of http.StatusUnauthorized")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should disallow non-existent confirms", func() {
		confirmsIDs := []string{
			"a",
			"gggg",
			"239r2938rh",
			"fffffffffffp",
			"a48hrt894",
			"vsdvdvs",
			"0000000000000000",
		}
		By("fetching the user account confirm to check for it to be unavailable")
		for _, confirmID := range confirmsIDs {
			apiEndpoint := apiServerAPIprefix + "/admin/useraccountconfirms/" + confirmID
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	It("should return invalid for non-existent confirms on public route", func() {
		confirmsIDs := []string{
			"a",
			"gggg",
			"239r2938rh",
			"fffffffffffp",
			"a48hrt894",
			"vsdvdvs",
			"0000000000000000",
		}
		By("fetching the user account confirm to check for it to be unavailable")
		for _, confirmID := range confirmsIDs {
			apiEndpoint := apiServerAPIprefix + "/user/confirm/" + confirmID
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	It("should authenticate an existing user", func() {
		By("posting to auth")
		apiEndpoint := apiServerAPIprefix + "/user/auth"
		userAccountLogin := types.UserSpec{
			Email:    regstrationForm.User.Email,
			Password: regstrationForm.User.Password,
		}
		userAccountLoginData, err := json.Marshal(userAccountLogin)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), userAccountLoginData, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("checking validation of the token")
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), userAccountLoginData, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginValid := routes.GetHTTPresponseBodyContents(resp).Data
		Expect(userAccountLoginValid).To(Equal(true), "JWT should be valid")
	})

	It("should return the profile of the same account and allow patching of an account", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")
		account.ID = userAccount.ID

		By("getting a JWT")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("checking the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err := json.Marshal(profileResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		Expect(profile.ID).To(Equal(userAccount.ID), "profile id does not match account id")
		Expect(profile.Names).To(Equal(userAccount.Names), "profile names does not match account names")
		Expect(profile.Email).To(Equal(userAccount.Email), "profile email does not match account email")
		Expect(profile.PhoneNumber).To(Equal(userAccount.PhoneNumber), "profile phoneNumber does not match account phoneNumber")
		Expect(profile.Birthday).To(Equal(userAccount.Birthday), "profile birthday does not match account birthday")

		By("patching the profile")
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
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err = json.Marshal(profileResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		profile = types.UserSpec{}
		json.Unmarshal(profileJSON, &profile)
		Expect(profile.ID).To(Equal(account.ID), "profile id does not match account id")
		Expect(profile.Names).To(Equal(profilePatch.Names), "profile names does not match profilePatch names")
		Expect(profile.Email).To(Equal(profilePatch.Email), "profile email does not match profilePatch email")
		Expect(profile.PhoneNumber).To(Equal(profilePatch.PhoneNumber), "profile phoneNumber does not match profilePatch phoneNumber")
		Expect(profile.Birthday).To(Equal(profilePatch.Birthday), "profile birthday does not match profilePatch birthday")

		By("getting a new JWT using new credentials")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should list all user accounts", func() {
		By("listing user accounts and checking the count")
		apiEndpoint := apiServerAPIprefix + "/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountsResponse := routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err := json.Marshal(userAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccounts []types.UserSpec
		json.Unmarshal(userAccountsBytes, &userAccounts)

		By("checking the response")
		Expect(len(userAccounts)).To(Equal(1), "invalid amount of users")

		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("listing user accounts and checking the count")
		apiEndpoint = apiServerAPIprefix + "/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountsResponse = routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err = json.Marshal(userAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		userAccounts = []types.UserSpec{}
		json.Unmarshal(userAccountsBytes, &userAccounts)

		By("checking the response")
		Expect(len(userAccounts)).To(Equal(2), "invalid amount of users")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("listing user accounts and checking the count")
		apiEndpoint = apiServerAPIprefix + "/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountsResponse = routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err = json.Marshal(userAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		userAccounts = []types.UserSpec{}
		json.Unmarshal(userAccountsBytes, &userAccounts)

		By("checking the response")
		Expect(len(userAccounts)).To(Equal(1), "invalid amount of users")
	})

	It("should return a user by their id", func() {
		By("creating a user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")

		By("fetch the user account by the account's id")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err = json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		userAccount = types.UserSpec{}
		json.Unmarshal(userAccountBytes, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account Names must be equal to account Names")
		Expect(userAccount.Email).To(Equal(account.Email), "User account Email must be equal to account Email")
		Expect(userAccount.PhoneNumber).To(Equal(account.PhoneNumber), "User account PhoneNumber must be equal to account PhoneNumber")
		Expect(userAccount.Birthday).To(Equal(account.Birthday), "User account PhoneNumber must be equal to account PhoneNumber")
		Expect(userAccount.Groups).To(Equal(account.Groups), "User account Groups must be equal to account Groups")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should not allow two emails with the same email when patching a profile", func() {
		By("creating the first user account")
		account1 := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account1)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount1 types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount1)

		By("checking the response")
		Expect(userAccount1.ID).ToNot(Equal(""), "User account ID must not be empty")

		By("creating the first user account")
		account2 := types.UserSpec{
			Names:       "Bloggs Joe",
			Email:       "joeblogblog@example.com",
			Password:    "Password123!",
			Groups:      []string{"flatmember"},
		}
		accountBytes, err = json.Marshal(account2)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err = json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount2 types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount2)

		By("checking the response")
		Expect(userAccount2.ID).ToNot(Equal(""), "User account ID must not be empty")

		By("logging in as the 2nd account")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("patching the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		profilePatch := types.UserSpec{
			Email:       "user123@example.com",
		}
		profilePatchData, err := json.Marshal(profilePatch)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest), "api have return code of http.StatusBadRequest")
		profileResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err := json.Marshal(profileResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		Expect(profile.Email).ToNot(Equal(profilePatch.Email), "profile email does not match profilePatch email")

		By("deleting the account1")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount1.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("deleting the account2")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount2.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should not allow non admins to patch their groups", func() {
		By("creating the first user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")

		By("logging in as the account")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("patching the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		profilePatch := types.UserSpec{
			Groups: []string{"flatmember", "admin"},
		}
		profilePatchData, err := json.Marshal(profilePatch)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		response := routes.GetHTTPresponseBodyContents(resp)
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse := response.Spec
		profileJSON, err := json.Marshal(profileResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		Expect(profile.Groups).To(Equal(account.Groups), "profile email does not match profilePatch email")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should not allow non admins to update their groups", func() {
		By("creating the first user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		response := routes.GetHTTPresponseBodyContents(resp)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := response.Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")

		By("logging in as the account")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("updating the profile")
		apiEndpoint = apiServerAPIprefix + "/user/profile"
		profilePatch := account
		profilePatch.Groups = []string{"flatmember", "admin"}
		profilePatchData, err := json.Marshal(profilePatch)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		response = routes.GetHTTPresponseBodyContents(resp)
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		profileResponse := response.Spec
		profileJSON, err := json.Marshal(profileResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		Expect(profile.Groups).To(Equal(account.Groups), "profile email does not match profilePatch email")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should fail to list non-existent user accounts", func() {
		ids := []string{
			"aa",
			"234fasdsad",
			"----------asd",
			",,,,,,,,,,,,",
			"0-93jr9-q23nfiunq398nr3948n5q3c3",
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		}

		By("fetch the user account by the account's id")
		for _, id := range ids {
			apiEndpoint := apiServerAPIprefix + "/admin/users/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	It("should list groups and include the default groups", func() {
		apiEndpoint := apiServerAPIprefix + "/groups"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		groupsResponse := routes.GetHTTPresponseBodyContents(resp).List
		groupsBytes, err := json.Marshal(groupsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var groups []types.GroupSpec
		json.Unmarshal(groupsBytes, &groups)

		Expect(len(groups) >= 2).To(Equal(true), "There must be at least two groups")

		for _, groupItem := range groups {
			apiEndpoint := apiServerAPIprefix + "/groups/" + groupItem.ID
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			groupResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			groupBytes, err := json.Marshal(groupResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var group types.GroupSpec
			json.Unmarshal(groupBytes, &group)

			Expect(groupItem.ID).To(Equal(group.ID), "GroupItem ID must match Group ID")
			Expect(groupItem.Name).To(Equal(group.Name), "GroupItem Name must match Group Name")
			Expect(groupItem.DefaultGroup).To(Equal(group.DefaultGroup), "GroupItem DefaultGroup must match Group DefaultGroup")
			Expect(groupItem.Description).To(Equal(group.Description), "GroupItem Description must match Group Description")
		}
	})

	It("should fail to list groups which don't exist", func() {
		ids := []string{
			"aa",
			"234fasdsad",
			"----------asd",
			",,,,,,,,,,,,",
			"0-93jr9-q23nfiunq398nr3948n5q3c3",
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		}

		By("fetch the user account by the account's id")
		for _, id := range ids {
			apiEndpoint := apiServerAPIprefix + "/groups/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound), "api have return code of http.StatusNotFound")
		}
	})

	It("should have the correct account groups for user account checking", func() {
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

		By("fetching all groups")
		apiEndpoint := apiServerAPIprefix + "/groups"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		groupsResponse := routes.GetHTTPresponseBodyContents(resp).List
		groupsBytes, err := json.Marshal(groupsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var groups []types.GroupSpec
		json.Unmarshal(groupsBytes, &groups)

		for _, account := range accounts {
			accountBytes, err := json.Marshal(account)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user account")
			apiEndpoint := apiServerAPIprefix + "/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")

			By("creating fetching user accounts")
			apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			By("logging in as the new user account")
			apiEndpoint = apiServerAPIprefix + "/user/auth"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			jwt := routes.GetHTTPresponseBodyContents(resp).Data.(string)
			Expect(jwt).ToNot(Equal(""), "JWT in response must not be empty")

			defer func() {
				By("deleting the account")
				apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
				resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
				Expect(err).To(BeNil(), "Request should not return an error")
				Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			}()

			for _, groupItem := range groups {
				expectGroup := false
				for _, accountGroup := range account.Groups {
					if groupItem.Name == accountGroup {
						expectGroup = true
					}
				}
				By("ensuring user account is or is not in a group")
				apiEndpoint := apiServerAPIprefix + "/user/can-i/group/" + groupItem.Name
				resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, jwt)
				Expect(err).To(BeNil(), "Request should not return an error")
				Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
				canIgroupResponse := routes.GetHTTPresponseBodyContents(resp).Data.(bool)
				Expect(canIgroupResponse).To(Equal(expectGroup), "Group was expected for this user account", account.Names, account.Groups, groupItem.Name, expectGroup)
			}
		}
	})

	It("should disallow creating an account with an existing email", func() {
		By("creating a user account")
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		By("creating the account again")
		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest), "api have return code of http.StatusOK")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should create a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing all shopping lists")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListsBytes, err := json.Marshal(shoppingListsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingLists []types.ShoppingListSpec
		json.Unmarshal(shoppingListsBytes, &shoppingLists)

		Expect(len(shoppingLists)).To(Equal(1), "there must be one shopping list")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should patch a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

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
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a shopping list")
			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
			resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListPatchBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingListPatchedResponse := routes.GetHTTPresponseBodyContents(resp)
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			shoppingListPatchedResponseSpec := shoppingListPatchedResponse.Spec
			shoppingListPatchedBytes, err := json.Marshal(shoppingListPatchedResponseSpec)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingListPatched types.ShoppingListSpec
			json.Unmarshal(shoppingListPatchedBytes, &shoppingListPatched)

			Expect(shoppingListPatched.ID).To(Equal(shoppingListCreated.ID), "shopping list id must be equal to shopping list created id")
			Expect(shoppingListPatched.Name).To(Equal(shoppingListPatch.Name), "shopping list name does not match shopping list created name")
			Expect(shoppingListPatched.Notes).To(Equal(shoppingListPatch.Notes), "shopping list notes does not match shopping list created notes")
		}

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should update a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

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
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a shopping list")
			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
			resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListUpdateBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingListUpdatedResponse := routes.GetHTTPresponseBodyContents(resp)
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			shoppingListUpdatedResponseSpec := shoppingListUpdatedResponse.Spec
			shoppingListUpdatedBytes, err := json.Marshal(shoppingListUpdatedResponseSpec)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingListUpdated types.ShoppingListSpec
			json.Unmarshal(shoppingListUpdatedBytes, &shoppingListUpdated)

			Expect(shoppingListUpdated.ID).To(Equal(shoppingListCreated.ID), "shopping list id must be equal to shopping list created id")
			Expect(shoppingListUpdated.Name).To(Equal(shoppingListUpdate.Name), "shopping list name does not match shopping list created name")
			Expect(shoppingListUpdated.Notes).To(Equal(shoppingListUpdate.Notes), "shopping list notes does not match shopping list created notes")
			Expect(shoppingListUpdated.Completed).To(Equal(shoppingListUpdate.Completed), "shopping list completed does not match shopping list created completed")
		}

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should not allow invalid shopping list properties", func() {
		shoppingLists := []types.ShoppingListSpec{
			{
				Name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			{
				Name: "",
			},
			{
				Name: "My list",
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
			{
				Name: "My list",
				TemplateID: "x",
			},
		}

		for _, shoppingList := range shoppingLists {
			shoppingListBytes, err := json.Marshal(shoppingList)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a shopping list")
			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest), "api have return code of http.StatusOK")
			shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListBytes, err = json.Marshal(shoppingListResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingListCreated types.ShoppingListSpec
			json.Unmarshal(shoppingListBytes, &shoppingListCreated)

			Expect(shoppingListCreated.ID).To(Equal(""), "shopping list created id must not be empty")
		}
	})

	It("should allow adding items to a list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "Eggs",
				Quantity: 1,
			},
			{
				Name: "Onions",
				Price: 2,
				Quantity: 1,
			},
			{
				Name: "Pasta",
				Price: 0.8,
				Quantity: 3,
			},
			{
				Name: "Bread",
				Price: 3.5,
				Quantity: 4,
				Notes: "Sourdough",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges",
			},
		}

		By("listing shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			Expect(shoppingItem.ListID).To(Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		shoppingListItems = []types.ShoppingItemSpec{}
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(len(newShoppingListItems)), "There should be as many items added as on the shopping list")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should disallow adding items to a non-existent list", func() {
		By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "Eggs",
				Quantity: 1,
			},
			{
				Name: "Onions",
				Price: 2,
				Quantity: 1,
			},
			{
				Name: "Pasta",
				Price: 0.8,
				Quantity: 3,
			},
			{
				Name: "Bread",
				Price: 3.5,
				Quantity: 4,
				Notes: "Sourdough",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges",
			},
		}

		By("listing shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists/" + "xxxxxxxx" + "/items"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound), "api have return code of http.StatusBadRequest")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err := json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			Expect(shoppingItem.ID).To(Equal(""), "invalid shopping list items must not have ids")
		}
	})

	It("should not allow adding of invalid items to a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 0,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				Tag: "Fruits and veges",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
			},
		}

		By("listing shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest), "api must have return code of http.StatusBadRequest")
		}

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api must have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		shoppingListItems = []types.ShoppingItemSpec{}
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should allow updating of shopping list items", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating item on the list")
		newShoppingListItem := types.ShoppingItemSpec{
			Name: "Lettuce",
			Price: 3,
			Quantity: 2,
			Notes: "Not plastic bagged ones",
			Tag: "Fruits and veges",
		}

		shoppingItemBytes, err := json.Marshal(newShoppingListItem)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingItem types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
		Expect(shoppingItem.ListID).To(Equal(shoppingListCreated.ID), "shopping item must belong to a list")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("updating item on the list")
		updatedShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "Iceberg lettuce",
				Price: 4,
				Quantity: 1,
				Notes: "",
				Tag: "Salad",
				Obtained: true,
			},
			{
				Name: "Iceberg lettuce",
				Price: 4,
				Quantity: 1,
				Notes: "",
				Tag: "Salad",
				Obtained: false,
			},
		}

		for _, updatedShoppingListItem := range updatedShoppingListItems {
			shoppingItemBytes, err = json.Marshal(updatedShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items/" + shoppingItem.ID
			resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingItemResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemBytes, err := json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItemUpdated types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemBytes, &shoppingItemUpdated)
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

			Expect(shoppingItemUpdated.ID).ToNot(Equal(""), "shopping item id should not be nil")
			Expect(updatedShoppingListItem.Name).To(Equal(shoppingItemUpdated.Name), "shopping item name was not updated")
			Expect(updatedShoppingListItem.Price).To(Equal(shoppingItemUpdated.Price), "shopping item price was not updated")
			Expect(updatedShoppingListItem.Quantity).To(Equal(shoppingItemUpdated.Quantity), "shopping item quantity was not updated")
			Expect(updatedShoppingListItem.Notes).To(Equal(shoppingItemUpdated.Notes), "shopping item notes was not updated")
			Expect(updatedShoppingListItem.Tag).To(Equal(shoppingItemUpdated.Tag), "shopping item tag was not updated")
			Expect(updatedShoppingListItem.Obtained).To(Equal(shoppingItemUpdated.Obtained), "shopping item obtained was not updated")
		}

		By("deleting the shopping list item")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items/" + shoppingItem.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should allow patching of shopping list items", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating item on the list")
		newShoppingListItem := types.ShoppingItemSpec{
			Name: "Lettuce",
			Price: 3,
			Quantity: 2,
			Notes: "Not plastic bagged ones",
			Tag: "Fruits and veges",
		}

		shoppingItemBytes, err := json.Marshal(newShoppingListItem)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingItem types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
		Expect(shoppingItem.ListID).To(Equal(shoppingListCreated.ID), "shopping item must belong to a list")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("updating item on the list")
		updatedShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "Iceberg lettuce",
				Price: 4,
				Quantity: 1,
				Notes: "Just a note",
				Tag: "Salad",
			},
			{
				Name: "Iceberg lettuce",
				Price: 4,
				Quantity: 1,
				Notes: "This note should have some useful meaning",
				Tag: "Salad",
			},
		}

		for _, updatedShoppingListItem := range updatedShoppingListItems {
			shoppingItemBytes, err = json.Marshal(updatedShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items/" + shoppingItem.ID
			resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingItemResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemBytes, err := json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItemUpdated types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemBytes, &shoppingItemUpdated)
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

			Expect(updatedShoppingListItem.Name).To(Equal(shoppingItemUpdated.Name), "shopping item name was not updated")
			Expect(updatedShoppingListItem.Price).To(Equal(shoppingItemUpdated.Price), "shopping item price was not updated")
			Expect(updatedShoppingListItem.Quantity).To(Equal(shoppingItemUpdated.Quantity), "shopping item quantity was not updated")
			Expect(updatedShoppingListItem.Notes).To(Equal(shoppingItemUpdated.Notes), "shopping item notes was not updated")
			Expect(updatedShoppingListItem.Tag).To(Equal(shoppingItemUpdated.Tag), "shopping item tag was not updated")
		}

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should return a list of tags from a list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "Eggs",
				Quantity: 1,
				Tag: "Dairy",
			},
			{
				Name: "Onions",
				Price: 2,
				Quantity: 1,
				Tag: "Fruits and veges",
			},
			{
				Name: "Pasta",
				Price: 0.8,
				Quantity: 3,
				Tag: "General",
			},
			{
				Name: "Bread",
				Price: 3.5,
				Quantity: 4,
				Notes: "Sourdough",
				Tag: "General",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges",
			},
		}

		By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			Expect(shoppingItem.ListID).To(Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err := json.Marshal(shoppingListTags)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var tags []string
		json.Unmarshal(shoppingListTagBytes, &tags)

		Expect(len(tags)).To(Equal(3), "invalid amount of tags")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should only return tags from a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating the first shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "Eggs",
				Quantity: 1,
				Tag: "Dairy",
			},
			{
				Name: "Onions",
				Price: 2,
				Quantity: 1,
				Tag: "Fruits and veges",
			},
			{
				Name: "Pasta",
				Price: 0.8,
				Quantity: 3,
				Tag: "General",
			},
			{
				Name: "Bread",
				Price: 3.5,
				Quantity: 4,
				Notes: "Sourdough",
				Tag: "General",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges",
			},
		}

		By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			Expect(shoppingItem.ListID).To(Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		By("creating the 2nd list")
		shoppingList2 := types.ShoppingListSpec{
			Name: "My list 2",
		}
		shoppingListBytes, err = json.Marshal(shoppingList2)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingList2Created types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingList2Created)

		Expect(shoppingList2Created.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingList2Created.Name).To(Equal(shoppingList2.Name), "shopping list name does not match shopping list created name2")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingList2Items []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingList2Items)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating items on the list")
		newShoppingList2Items := []types.ShoppingItemSpec{
			{
				Name: "Eggs",
				Quantity: 1,
				Tag: "Dairy1",
			},
			{
				Name: "Oions",
				Price: 2,
				Quantity: 1,
				Tag: "Fruits and veges1",
			},
			{
				Name: "Pasta",
				Price: 0.8,
				Quantity: 3,
				Tag: "General1",
			},
			{
				Name: "Bread",
				Price: 3.5,
				Quantity: 4,
				Notes: "Sourdough",
				Tag: "General1",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges1",
			},
		}

		By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingList2Items {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			Expect(shoppingItem.ListID).To(Equal(shoppingList2Created.ID), "shopping item must belong to a list")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err := json.Marshal(shoppingListTags)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var listTags []string
		json.Unmarshal(shoppingListTagBytes, &listTags)

		Expect(len(listTags)).To(Equal(3), "invalid amount of tags")
		containsTagsFromOtherLists := false
		for _, tag := range listTags {
			for _, listItems := range newShoppingList2Items {
				if tag == listItems.Tag {
					containsTagsFromOtherLists = true
				}
			}
		}
		Expect(containsTagsFromOtherLists).To(Equal(false), "list of tags contains tags from other lists")

		By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err = json.Marshal(shoppingListTags)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var list2Tags []string
		json.Unmarshal(shoppingListTagBytes, &list2Tags)

		Expect(len(list2Tags)).To(Equal(3), "invalid amount of tags")
		containsTagsFromOtherLists = false
		for _, tag := range list2Tags {
			for _, listItems := range newShoppingListItems {
				if tag == listItems.Tag {
					containsTagsFromOtherLists = true
					break
				}
			}
		}
		Expect(containsTagsFromOtherLists).To(Equal(false), "list of tags contains tags from other lists")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("deleting the 2nd shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingList2Created.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should allow updating of tags in a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating items on the list")
		newShoppingListItem := types.ShoppingItemSpec{
			Name: "Lettuce",
			Price: 3,
			Quantity: 2,
			Notes: "Not plastic bagged ones",
			Tag: "Fruits and veges",
		}

		By("creating shopping list items")
		shoppingListItemBytes, err := json.Marshal(newShoppingListItem)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListItemBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingItem types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
		Expect(shoppingItem.ListID).To(Equal(shoppingListCreated.ID), "shopping item must belong to a list")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		shoppingItemTagUpdate := types.ShoppingTag{
			Name: "Veges",
		}
		By("patching the tag")
		shoppingItemTagBytes, err := json.Marshal(shoppingItemTagUpdate)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags/Fruits%20and%20veges"
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingItemTagBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		By("fetching the shopping list tags")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/tags"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListTags := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListTagBytes, err := json.Marshal(shoppingListTags)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var listTags []string
		json.Unmarshal(shoppingListTagBytes, &listTags)

		var foundUpdatedTag bool
		for _, tag := range listTags {
			if tag == shoppingItemTagUpdate.Name {
				foundUpdatedTag = true
			}
		}
		Expect(foundUpdatedTag).To(Equal(true), "Unable to find updated tag")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
	})

	It("should allow management of shopping tags", func() {
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

		By("creating shopping tags")
		for _, shoppingTag := range shoppingTags {
			shoppingTagBytes, err := json.Marshal(shoppingTag)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/tags"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingTagBytes, "")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingTagResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingTagBytes, err = json.Marshal(shoppingTagResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingTag types.ShoppingTag
			json.Unmarshal(shoppingTagBytes, &shoppingTag)
			Expect(shoppingTag.ID).ToNot(Equal(""), "shopping tag must have an ID")
		}

		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/tags"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingTagsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingTagBytes, err := json.Marshal(shoppingTagsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var tags []types.ShoppingTag
		json.Unmarshal(shoppingTagBytes, &tags)

		Expect(len(tags)).To(Equal(len(shoppingTags)), "failed to find the correct number (%v) of shopping tags in length of list (%v)", len(shoppingTags), len(tags))
		foundTags := 0
		for _, tag := range tags {
			for _, expectedTag := range shoppingTags {
				if expectedTag.Name == tag.Name {
					foundTags++
				}
			}
		}
		Expect(foundTags).To(Equal(len(shoppingTags)), "failed to find the correct (%v) number of shopping tags in list of tags from response (%v)", len(shoppingTags), foundTags)

		// update tag name
		tagUpdate := types.ShoppingTag{
			Name: "Fruits, Veges, and Fresh",
		}
		shoppingTagUpdateBytes, err := json.Marshal(tagUpdate)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/tags/" + tags[0].ID
		resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingTagUpdateBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")

		// get tag
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/tags/" + tags[0].ID
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		shoppingTagUpdateGetResponse := routes.GetHTTPresponseBodyContents(resp)
		shoppingTagBytes, err = json.Marshal(shoppingTagUpdateGetResponse.Spec)
		var shoppingTagUpdated types.ShoppingTag
		json.Unmarshal(shoppingTagBytes, &shoppingTagUpdated)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK", shoppingTagUpdateGetResponse.Metadata.Response)
		Expect(shoppingTagUpdated.ID).To(Equal(tags[0].ID), "shopping tag must have an ID (%v) matching it's previous ID (%v)", shoppingTagUpdated.ID, tags[0].ID)
		Expect(shoppingTagUpdated.Name).To(Equal(tagUpdate.Name), "shopping tag must have the new tag name")

		for _, tag := range tags {
			apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/tags/" + tag.ID
			resp, err := httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

	})

	It("should allow templating of a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.ID).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing shopping list items")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err := json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListItems []types.ShoppingItemSpec
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)

		Expect(len(shoppingListItems)).To(Equal(0), "There should be no items on the shopping list")

		By("creating items on the list")
		newShoppingListItems := []types.ShoppingItemSpec{
			{
				Name: "Eggs",
				Quantity: 1,
				Tag: "Dairy",
			},
			{
				Name: "Onions",
				Price: 2,
				Quantity: 1,
				Tag: "Fruits and veges",
			},
			{
				Name: "Pasta",
				Price: 0.8,
				Quantity: 3,
				Tag: "General",
			},
			{
				Name: "Bread",
				Price: 3.5,
				Quantity: 4,
				Notes: "Sourdough",
				Tag: "General",
			},
			{
				Name: "Lettuce",
				Price: 3,
				Quantity: 2,
				Notes: "Not plastic bagged ones",
				Tag: "Fruits and veges",
			},
		}

		By("creating shopping list items")
		for _, newShoppingListItem := range newShoppingListItems {
			shoppingListBytes, err := json.Marshal(newShoppingListItem)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID + "/items"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingItemResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			shoppingListItemsBytes, err = json.Marshal(shoppingItemResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingItem types.ShoppingItemSpec
			json.Unmarshal(shoppingListItemsBytes, &shoppingItem)
			Expect(shoppingItem.ListID).To(Equal(shoppingListCreated.ID), "shopping item must belong to a list")
			Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		}

		shoppingListFromTemplate := types.ShoppingListSpec{
			Name: "My list (from template)",
			Notes: "This is a templated list",
			TemplateID: shoppingListCreated.ID,
		}
		shoppingListBytes, err = json.Marshal(shoppingListFromTemplate)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a templated shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		templatedShoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		templatedShoppingListBytes, err := json.Marshal(templatedShoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListTemplatedCreated types.ShoppingListSpec
		json.Unmarshal(templatedShoppingListBytes, &shoppingListTemplatedCreated)

		By("listing items of the templated shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListTemplatedCreated.ID + "/items"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		shoppingListItemsResponse = routes.GetHTTPresponseBodyContents(resp).List
		shoppingListItemsBytes, err = json.Marshal(shoppingListItemsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		shoppingListItems = []types.ShoppingItemSpec{}
		json.Unmarshal(shoppingListItemsBytes, &shoppingListItems)
		Expect(len(newShoppingListItems)).To(Equal(len(shoppingListItems)), "templated list must have the same amount of items as the orignal list")

		foundTotal := 0
		for _, item := range newShoppingListItems {
			for _, templatedItem := range shoppingListItems {
				if templatedItem.Name == item.Name {
					foundTotal++
					continue
				}
			}
		}

		Expect(foundTotal).To(Equal(len(newShoppingListItems)), "unable to find all items from original list in templated list")

		By("deleting the shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api must have return code of http.StatusOK")

		By("deleting the template shopping list")
		apiEndpoint = apiServerAPIprefix + "/apps/shoppinglist/lists/" + shoppingListTemplatedCreated.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api must have return code of http.StatusOK")
	})

	It("should require authorization for protected routes", func() {
		apiEndpoint := apiServer + "/" + apiServerAPIprefix + "/user/profile"
		req, err := http.NewRequest("GET", apiEndpoint, nil)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized), "endpoint should be restricted")
		requestResp := routes.GetHTTPresponseBodyContents(resp)
		Expect(requestResp.Metadata.Response).To(Equal("Unauthorized"), "")
	})

	It("should require admin for admin protected routes", func() {
		account := types.UserSpec{
			Names:       "Joe Bloggs",
			Email:       "user123@example.com",
			Password:    "Password123!",
			PhoneNumber: "64200000000",
			Birthday:    43200,
			Groups:      []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account with no admin access")
		apiEndpoint := apiServerAPIprefix + "/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.ID).ToNot(Equal(""), "User account ID must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("logging in")
		apiEndpoint = apiServerAPIprefix + "/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("trying to use an admin route")
		apiEndpoint = apiServerAPIprefix + "/admin/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusForbidden), "api have return code of http.StatusForbidden")

		By("deleting the account")
		apiEndpoint = apiServerAPIprefix + "/admin/users/" + userAccount.ID
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(http.StatusOK), "api have return code of http.StatusOK")
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

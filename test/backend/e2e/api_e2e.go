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
		apiEndpoint := "api"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		Expect(routes.GetHTTPresponseBodyContents(resp).Metadata.URL).To(Equal(fmt.Sprintf("/%v", apiEndpoint)))
	})

	It("should be initialized", func() {
		By("querying the API")
		apiEndpoint := "api/system/initialized"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		Expect(routes.GetHTTPresponseBodyContents(resp).Data.(bool)).To(Equal(true), "instance should be initialized")
	})

	It("should have a flat name", func() {
		By("querying the API")
		apiEndpoint := "api/system/flatName"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		Expect(routes.GetHTTPresponseBodyContents(resp).Spec.(string)).ToNot(Equal(""), "flatName should not be empty")
	})

	// TODO /api/admin/register - not yet possible in the same manner

	It("should have at least one user account", func() {
		By("querying the API")
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		response := routes.GetHTTPresponseBodyContents(resp)
		Expect(len(response.List.([]interface{})) > 0).To(Equal(true), "should at least one user account")
	})

	It("should return properties of a single user account", func() {
		By("listing all user accounts")
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		allUserAccountsResponse := routes.GetHTTPresponseBodyContents(resp).List
		allUserAccountsJSON, err := json.Marshal(allUserAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var allUserAccounts []types.UserSpec
		json.Unmarshal(allUserAccountsJSON, &allUserAccounts)
		firstUserAccount := allUserAccounts[0]

		By("listing all user accounts")
		apiEndpoint = "api/admin/users/" + firstUserAccount.Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
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
				Email:    "user1@example.com",
				Password: "Password123!",
				Groups:   []string{"flatmember", "admin"},
			},
			{
				Names:    "Joe Bloggs",
				Email:    "user1@example.coop",
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
				Email:       "us.er1@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "020 000 0000",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "us.er1@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "+64200000000",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "us.er1@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "64200000000",
			},
			{
				Names:       "Joe Bloggs",
				Email:       "us.er1@example.coop",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
				PhoneNumber: "64-20-000-000",
			},
			{
				Names:  "Joe Bloggs",
				Email:  "user1@example.com",
				Groups: []string{"flatmember"},
			},
		}
		for _, account := range accounts {
			accountBytes, err := json.Marshal(account)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user account")
			apiEndpoint := "api/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("deleting the account")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
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
			apiEndpoint := "api/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "API should return error")
			Expect(resp.StatusCode).To(Equal(400), "API must have return code of 400")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.Id).To(Equal(""), "User account Id must be empty")
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
			apiEndpoint := "api/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("creating fetching user accounts")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Email).To(Equal(account.Email), "User account email must match what was posted")
			Expect(userAccount.Groups).To(Equal(account.Groups), "User account email must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("deleting the account")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
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
			apiEndpoint := "api/admin/users/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(404), "API must have return code of 404")
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
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("updating the resource")
		accountUpdates := []types.UserSpec{
			{
				Names: "John Blogs",
				Email: "user123@example.com",
				Groups: []string{"flatmember", "admin"},
				PhoneNumber: "02000000",
				Birthday: 0,
				Password: "Password1234!",
			},
			{
				Names: "John 'golang' Smith",
				Email: "user1237@example.com",
				Groups: []string{"flatmember"},
				PhoneNumber: "0200000000",
				Birthday: 0420001,
				Password: "Password123!",
			},
			{
				Names: "E",
				Email: "user1237@example.com",
				Groups: []string{"flatmember", "admin"},
				PhoneNumber: "0200000001",
				Birthday: 0420002,
				Password: "Password123!",
			},
		}
		for _, accountUpdate := range accountUpdates {
			By("updating account")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			accountUpdateBytes, err := json.Marshal(accountUpdate)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			resp, err = httpRequestWithHeader("PUT", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountUpdateBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")

			By("creating fetching the user account")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
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
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
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
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("list user account confirms")
		apiEndpoint = "api/admin/useraccountconfirms"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v?userId=%v", apiServer, apiEndpoint, userAccount.Id), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		confirmsListResponse := routes.GetHTTPresponseBodyContents(resp).List
		confirmsListJSON, err := json.Marshal(confirmsListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var confirmsList []types.UserCreationSecretSpec
		json.Unmarshal(confirmsListJSON, &confirmsList)
		Expect(len(confirmsList) > 0).To(Equal(true), "must contain at least one confirm")

		By("fetching the user account confirm")
		apiEndpoint = "api/admin/useraccountconfirms/" + confirmsList[0].Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		confirmResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		confirmJSON, err := json.Marshal(confirmResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var confirm types.UserCreationSecretSpec
		json.Unmarshal(confirmJSON, &confirm)
		Expect(confirm.Id).ToNot(Equal(""), "confirm id must not be empty")
		Expect(confirm.UserId).ToNot(Equal(""), "confirm userid must not be empty")
		Expect(confirm.Secret).ToNot(Equal(""), "confirm secret must not be empty")
		Expect(confirm.Valid).To(Equal(true), "confirm valid must be true")

		By("fetching the public route for user account confirm to check for it to be valid")
		apiEndpoint = "api/user/confirm/" + confirmsList[0].Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		confirmValid := routes.GetHTTPresponseBodyContents(resp).Data
		Expect(confirmValid.(bool)).To(Equal(true), "confirm valid must be true")

		By("fetching the user account confirm")
		apiEndpoint = "api/user/confirm/" + confirm.Id + "?secret=" + confirm.Secret
		confirmUserAccount := types.UserSpec{
			Password: "Password123!",
		}
		confirmUserAccountJSON, err := json.Marshal(confirmUserAccount)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), confirmUserAccountJSON, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")

		By("fetching the user account confirm to check for it to be unavailable")
		apiEndpoint = "api/admin/useraccountconfirms/" + confirmsList[0].Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(404), "API must have return code of 200")

		By("fetching the user account to check if it's been registered")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err = json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccountRegistered types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccountRegistered)
		Expect(userAccountRegistered.Id).ToNot(Equal(""), "user account id must not be empty")
		Expect(userAccountRegistered.Registered).To(Equal(true), "account must be registered")

		By("deleting the account")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
	})

	It("should disallow a deleted account to log in", func() {
		account := types.UserSpec{
			Names: "Joe Bloggs",
			Email: "user123@example.com",
			Password: "Password123!",
			PhoneNumber: "64200000000",
			Birthday: 43200,
			Groups: []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("logging in")
		apiEndpoint = "api/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("checking validation of the token")
		apiEndpoint = "api/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountLoginValid := routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		Expect(userAccountLoginValid).To(Equal(true), "JWT should be valid")

		By("deleting the account")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")

		By("checking validation of the token")
		apiEndpoint = "api/user/auth"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(401), "API must have return code of 401")
		userAccountLoginValid = routes.GetHTTPresponseBodyContents(resp).Data.(bool)
		Expect(userAccountLoginValid).To(Equal(false), "JWT should be valid")
	})

	It("should disallow non-existent confirms", func() {
		confirmsIds := []string{
			"a",
			"gggg",
			"239r2938rh",
			"fffffffffffp",
			"a48hrt894",
			"vsdvdvs",
			"0000000000000000",
		}
		By("fetching the user account confirm to check for it to be unavailable")
		for _, confirmId := range confirmsIds {
			apiEndpoint := "api/admin/useraccountconfirms/" + confirmId
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(404), "API must have return code of 404")
		}
	})

	It("should return invalid for non-existent confirms on public route", func() {
		confirmsIds := []string{
			"a",
			"gggg",
			"239r2938rh",
			"fffffffffffp",
			"a48hrt894",
			"vsdvdvs",
			"0000000000000000",
		}
		By("fetching the user account confirm to check for it to be unavailable")
		for _, confirmId := range confirmsIds {
			apiEndpoint := "api/user/confirm/" + confirmId
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(404), "API must have return code of 404")
		}
	})

	It("should authenticate an existing user", func() {
		By("posting to auth")
		apiEndpoint := "api/user/auth"
		userAccountLogin := types.UserSpec{
			Email: regstrationForm.User.Email,
			Password: regstrationForm.User.Password,
		}
		userAccountLoginData, err := json.Marshal(userAccountLogin)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), userAccountLoginData, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("checking validation of the token")
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), userAccountLoginData, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountLoginValid := routes.GetHTTPresponseBodyContents(resp).Data
		Expect(userAccountLoginValid).To(Equal(true), "JWT should be valid")
	})

	It("should return the profile of the same account and allow patching of an account", func() {
		account := types.UserSpec{
			Names: "Joe Bloggs",
			Email: "user123@example.com",
			Password: "Password123!",
			PhoneNumber: "64200000000",
			Birthday: 43200,
			Groups: []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")
		account.Id = userAccount.Id

		By("getting a JWT")
		apiEndpoint = "api/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountLoginResponseData := routes.GetHTTPresponseBodyContents(resp).Data.(string)
		Expect(userAccountLoginResponseData).ToNot(Equal(""), "JWT in response must not be empty")

		By("checking the profile")
		apiEndpoint = "api/user/profile"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		profileResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err := json.Marshal(profileResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var profile types.UserSpec
		json.Unmarshal(profileJSON, &profile)
		Expect(profile.Id).To(Equal(userAccount.Id), "profile id does not match account id")
		Expect(profile.Names).To(Equal(userAccount.Names), "profile names does not match account names")
		Expect(profile.Email).To(Equal(userAccount.Email), "profile email does not match account email")
		Expect(profile.PhoneNumber).To(Equal(userAccount.PhoneNumber), "profile phoneNumber does not match account phoneNumber")
		Expect(profile.Birthday).To(Equal(userAccount.Birthday), "profile birthday does not match account birthday")

		By("patching the profile")
		apiEndpoint = "api/user/profile"
		profilePatch := types.UserSpec{
			Id: "aaaaaaa",
			Names: "Jonno bloggo",
			Email: "user2@example.com",
			Password: "Password1234!",
			PhoneNumber: "+64200000001",
			Birthday: 432001,
		}
		profilePatchData, err := json.Marshal(profilePatch)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, userAccountLoginResponseData)
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		profileResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		profileJSON, err = json.Marshal(profileResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		profile = types.UserSpec{}
		json.Unmarshal(profileJSON, &profile)
		Expect(profile.Id).To(Equal(account.Id), "profile id does not match account id")
		Expect(profile.Names).To(Equal(profilePatch.Names), "profile names does not match profilePatch names")
		Expect(profile.Email).To(Equal(profilePatch.Email), "profile email does not match profilePatch email")
		Expect(profile.PhoneNumber).To(Equal(profilePatch.PhoneNumber), "profile phoneNumber does not match profilePatch phoneNumber")
		Expect(profile.Birthday).To(Equal(profilePatch.Birthday), "profile birthday does not match profilePatch birthday")

		By("getting a new JWT using new credentials")
		apiEndpoint = "api/user/auth"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), profilePatchData, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")

		By("deleting the account")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
	})

	It("should list all user accounts", func() {
		By("listing user accounts and checking the count")
		apiEndpoint := "api/users"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountsResponse := routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err := json.Marshal(userAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccounts []types.UserSpec
		json.Unmarshal(userAccountsBytes, &userAccounts)

		By("checking the response")
		Expect(len(userAccounts)).To(Equal(1), "invalid amount of users")

		account := types.UserSpec{
			Names: "Joe Bloggs",
			Email: "user123@example.com",
			Password: "Password123!",
			PhoneNumber: "64200000000",
			Birthday: 43200,
			Groups: []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user account")
		apiEndpoint = "api/admin/users"
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the response")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("listing user accounts and checking the count")
		apiEndpoint = "api/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountsResponse = routes.GetHTTPresponseBodyContents(resp).List
		userAccountsBytes, err = json.Marshal(userAccountsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		userAccounts = []types.UserSpec{}
		json.Unmarshal(userAccountsBytes, &userAccounts)

		By("checking the response")
		Expect(len(userAccounts)).To(Equal(2), "invalid amount of users")

		By("deleting the account")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")

		By("listing user accounts and checking the count")
		apiEndpoint = "api/users"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
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
			Names: "Joe Bloggs",
			Email: "user123@example.com",
			Password: "Password123!",
			PhoneNumber: "64200000000",
			Birthday: 43200,
			Groups: []string{"flatmember"},
		}
		accountBytes, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountBytes, &userAccount)

		By("checking the response")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")

		By("fetch the user account by the account's id")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountBytes, err = json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		userAccount = types.UserSpec{}
		json.Unmarshal(userAccountBytes, &userAccount)

		By("checking the response")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account Names must be equal to account Names")
		Expect(userAccount.Email).To(Equal(account.Email), "User account Email must be equal to account Email")
		Expect(userAccount.PhoneNumber).To(Equal(account.PhoneNumber), "User account PhoneNumber must be equal to account PhoneNumber")
		Expect(userAccount.Birthday).To(Equal(account.Birthday), "User account PhoneNumber must be equal to account PhoneNumber")
		Expect(userAccount.Groups).To(Equal(account.Groups), "User account Groups must be equal to account Groups")

		By("deleting the account")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
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
			apiEndpoint := "api/admin/users/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(404), "API must have return code of 404")
		}
	})

	It("should list groups and include the default groups", func() {
		apiEndpoint := "api/groups"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		groupsResponse := routes.GetHTTPresponseBodyContents(resp).List
		groupsBytes, err := json.Marshal(groupsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var groups []types.GroupSpec
		json.Unmarshal(groupsBytes, &groups)

		Expect(len(groups) >= 2).To(Equal(true), "There must be at least two groups")

		for _, groupItem := range groups {
			apiEndpoint := "api/groups/" + groupItem.Id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			groupResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			groupBytes, err := json.Marshal(groupResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var group types.GroupSpec
			json.Unmarshal(groupBytes, &group)

			Expect(groupItem.Id).To(Equal(group.Id), "GroupItem Id must match Group Id")
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
			apiEndpoint := "api/groups/" + id
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(404), "API must have return code of 404")
		}
	})

	It("should have the correct account groups for user account checking", func() {
		accounts := []types.UserSpec{
			{
				Names:       "Joe Bloggs",
				Email:       "user1@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember", "admin"},
			},
			{
				Names:       "Joe Bloggs 2",
				Email:       "user2@example.com",
				Password:    "Password123!",
				Groups:      []string{"flatmember"},
			},
		}

		By("fetching all groups")
		apiEndpoint := "api/groups"
		resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		groupsResponse := routes.GetHTTPresponseBodyContents(resp).List
		groupsBytes, err := json.Marshal(groupsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var groups []types.GroupSpec
		json.Unmarshal(groupsBytes, &groups)

		for _, account := range accounts {
			accountBytes, err := json.Marshal(account)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user account")
			apiEndpoint := "api/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the response")
			Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")

			By("creating fetching user accounts")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			By("logging in as the new user account")
			apiEndpoint = "api/user/auth"
			resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			jwt := routes.GetHTTPresponseBodyContents(resp).Data.(string)
			Expect(jwt).ToNot(Equal(""), "JWT in response must not be empty")

			defer func() {
				By("deleting the account")
				apiEndpoint = "api/admin/users/" + userAccount.Id
				resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
				Expect(err).To(BeNil(), "Request should not return an error")
				Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			}()

			for _, groupItem := range groups {
				expectGroup := false
				for _, accountGroup := range account.Groups {
					if groupItem.Name == accountGroup {
						expectGroup = true
					}
				}
				By("ensuring user account is or is not in a group")
				apiEndpoint := "api/user/can-i/group/" + groupItem.Name
				resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, jwt)
				Expect(err).To(BeNil(), "Request should not return an error")
				Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
				canIgroupResponse := routes.GetHTTPresponseBodyContents(resp).Data.(bool)
				Expect(canIgroupResponse).To(Equal(expectGroup), "Group was expected for this user account", account.Names, account.Groups, groupItem.Name, expectGroup)
			}
		}
	})

	It("should create a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := "api/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.Id).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		By("listing all shopping lists")
		apiEndpoint = "api/apps/shoppinglist/lists"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		shoppingListsResponse := routes.GetHTTPresponseBodyContents(resp).List
		shoppingListsBytes, err := json.Marshal(shoppingListsResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingLists []types.ShoppingListSpec
		json.Unmarshal(shoppingListsBytes, &shoppingLists)

		Expect(len(shoppingLists)).To(Equal(1), "there must be one shopping list")

		By("deleting the shopping list")
		apiEndpoint = "api/apps/shoppinglist/lists/" + shoppingListCreated.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
	})

	It("should patch a shopping list", func() {
		shoppingList := types.ShoppingListSpec{
			Name: "My list",
		}
		shoppingListBytes, err := json.Marshal(shoppingList)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a shopping list")
		apiEndpoint := "api/apps/shoppinglist/lists"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListBytes, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		shoppingListResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		shoppingListBytes, err = json.Marshal(shoppingListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var shoppingListCreated types.ShoppingListSpec
		json.Unmarshal(shoppingListBytes, &shoppingListCreated)

		Expect(shoppingListCreated.Id).ToNot(Equal(""), "shopping list created id must not be empty")
		Expect(shoppingListCreated.Name).To(Equal(shoppingList.Name), "shopping list name does not match shopping list created name")

		shoppingListPatches := []types.ShoppingListSpec{
			{
				Name: "Week 19",
			},
			{
				Name: "Week 19a",
			},
			{
				Name: "My neat list",
				Notes: "Well, this list is neat.",
			},
		}
		for _, shoppingListPatch := range shoppingListPatches {
			shoppingListPatchBytes, err := json.Marshal(shoppingListPatch)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a shopping list")
			apiEndpoint = "api/apps/shoppinglist/lists/" + shoppingListCreated.Id
			resp, err = httpRequestWithHeader("PATCH", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), shoppingListPatchBytes, "")
			Expect(err).To(BeNil(), "Request should not return an error")
			shoppingListPatchedResponse := routes.GetHTTPresponseBodyContents(resp)
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			shoppingListPatchedResponseSpec := shoppingListPatchedResponse.Spec
			shoppingListPatchedBytes, err := json.Marshal(shoppingListPatchedResponseSpec)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var shoppingListPatched types.ShoppingListSpec
			json.Unmarshal(shoppingListPatchedBytes, &shoppingListPatched)

			Expect(shoppingListPatched.Id).To(Equal(shoppingListCreated.Id), "shopping list id must be equal to shopping list created id")
			Expect(shoppingListPatched.Name).To(Equal(shoppingListPatch.Name), "shopping list name does not match shopping list created name")
			Expect(shoppingListPatched.Notes).To(Equal(shoppingListPatch.Notes), "shopping list notes does not match shopping list created notes")
		}

		By("deleting the shopping list")
		apiEndpoint = "api/apps/shoppinglist/lists/" + shoppingListCreated.Id
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil, "")
		Expect(err).To(BeNil(), "Request should not return an error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
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

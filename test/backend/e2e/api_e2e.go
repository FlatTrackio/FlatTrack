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
	})

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
				Names:    "Joe Bloggs",
				Email:    "user1@example.com",
				Groups:   []string{"flatmember"},
			},
		}
		for _, account := range accounts {
			accountData, err := json.Marshal(account)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user accounts")
			apiEndpoint := "api/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
			Expect(err).To(BeNil(), "API should not return error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the results")
			Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("deleting the account")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
			Expect(err).To(BeNil(), "API should not return error")
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
			accountData, err := json.Marshal(account)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user accounts")
			apiEndpoint := "api/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
			Expect(err).To(BeNil(), "API should return error")
			Expect(resp.StatusCode).To(Equal(400), "API must have return code of 400")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the results")
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
			accountData, err := json.Marshal(account)
			Expect(err).To(BeNil(), "failed to marshal to JSON")

			By("creating a user accounts")
			apiEndpoint := "api/admin/users"
			resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
			Expect(err).To(BeNil(), "API should not return error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err := json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			var userAccount types.UserSpec
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the results")
			Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("creating fetching user accounts")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
			Expect(err).To(BeNil(), "API should not return error")
			Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
			userAccountResponse = routes.GetHTTPresponseBodyContents(resp).Spec
			userAccountJSON, err = json.Marshal(userAccountResponse)
			Expect(err).To(BeNil(), "failed to marshal to JSON")
			userAccount = types.UserSpec{}
			json.Unmarshal(userAccountJSON, &userAccount)

			By("checking the results")
			Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
			Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
			Expect(userAccount.Email).To(Equal(account.Email), "User account email must match what was posted")
			Expect(userAccount.Groups).To(Equal(account.Groups), "User account email must match what was posted")
			Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

			By("deleting the account")
			apiEndpoint = "api/admin/users/" + userAccount.Id
			resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
			Expect(err).To(BeNil(), "API should not return error")
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
			resp, err := httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
			Expect(err).To(BeNil(), "API should not return error")
			Expect(resp.StatusCode).To(Equal(404), "API must have return code of 404")
		}
	})

	It("should create account confirms when a password is not provided and allow it to be confirmed", func() {
		account := types.UserSpec{
			Names: "Joe Bloggs",
			Email: "user1@example.com",
			Groups: []string{"flatmember"},
		}

		accountData, err := json.Marshal(account)
		Expect(err).To(BeNil(), "failed to marshal to JSON")

		By("creating a user accounts")
		apiEndpoint := "api/admin/users"
		resp, err := httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		userAccountResponse := routes.GetHTTPresponseBodyContents(resp).Spec
		userAccountJSON, err := json.Marshal(userAccountResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var userAccount types.UserSpec
		json.Unmarshal(userAccountJSON, &userAccount)

		By("checking the results")
		Expect(userAccount.Id).ToNot(Equal(""), "User account Id must not be empty")
		Expect(userAccount.Names).To(Equal(account.Names), "User account names must match what was posted")
		Expect(userAccount.Password).To(Equal(""), "User account password must return an empty string")

		By("list user account confirms")
		apiEndpoint = "api/admin/useraccountconfirms"
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v?userId=%v", apiServer, apiEndpoint, userAccount.Id), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")
		confirmsListResponse := routes.GetHTTPresponseBodyContents(resp).List
		confirmsListJSON, err := json.Marshal(confirmsListResponse)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		var confirmsList []types.UserCreationSecretSpec
		json.Unmarshal(confirmsListJSON, &confirmsList)
		Expect(len(confirmsList) > 0).To(Equal(true), "must contain at least one confirm")

		By("fetching the user account confirm")
		apiEndpoint = "api/admin/useraccountconfirms/" + confirmsList[0].Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
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

		By("fetching the user account confirm")
		apiEndpoint = "api/user/confirm/" + confirm.Id + "?secret=" + confirm.Secret
		confirmUserAccount := types.UserSpec{
			Password: "Password123!",
		}
		confirmUserAccountJSON, err := json.Marshal(confirmUserAccount)
		Expect(err).To(BeNil(), "failed to marshal to JSON")
		resp, err = httpRequestWithHeader("POST", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), confirmUserAccountJSON)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")

		By("fetching the user account confirm to check for it to be unavailable")
		apiEndpoint = "api/admin/useraccountconfirms/" + confirmsList[0].Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), nil)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(404), "API must have return code of 200")

		By("fetching the user account to check if it's been registered")
		apiEndpoint = "api/admin/users/" + userAccount.Id
		resp, err = httpRequestWithHeader("GET", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
		Expect(err).To(BeNil(), "API should not return error")
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
		resp, err = httpRequestWithHeader("DELETE", fmt.Sprintf("%v/%v", apiServer, apiEndpoint), accountData)
		Expect(err).To(BeNil(), "API should not return error")
		Expect(resp.StatusCode).To(Equal(200), "API must have return code of 200")

	})
})

func httpRequestWithHeader(verb string, url string, data []byte) (resp *http.Response, err error) {
	req, err := http.NewRequest(verb, url, bytes.NewBuffer(data))
	req.Header.Set("Authorization", "bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	return resp, err
}

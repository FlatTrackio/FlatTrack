/*
	route related
*/

package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/users"
)

// PostUser
// create a user
func PostUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO
		// - handle existing user
		// - handle empty fields

		code := 500
		response := "Failed to create user account"
		invalid := false

		accountName := r.FormValue("names")
		if common.RegexMatchName(accountName) == false {
			invalid = true
			code = 400
			response = "Unable to use that name"
		}
		accountEmail := r.FormValue("email")
		if common.RegexMatchEmail(accountEmail) == false {
			invalid = true
			code = 400
			response = "Unable to use that email"
		}
		// TODO add group validation - requires creating admin and flatmember in migrations
		accountGroups := r.FormValue("groups")
		fmt.Println(accountGroups)
		accountPassword := r.FormValue("password")
		if common.RegexMatchPassword(accountPassword) == false {
			invalid = true
			code = 400
			response = "Unable to use that password"
		}
		accountPhoneNumber := r.FormValue("phoneNumber")
		if accountPhoneNumber != "" && common.RegexMatchPhoneNumber(accountPhoneNumber) == false {
			invalid = true
			code = 400
			response = "Unable to use that phone number"
		}
		user := types.UserSpec{
			Names:       accountName,
			Email:       accountEmail,
			Password:    accountPassword,
			PhoneNumber: accountPhoneNumber,
		}

		if invalid == true {
			JSONresp := types.JSONMessageResponse{
				Metadata: types.JSONResponseMetadata{
					Response: response,
				},
				Spec: types.User{},
			}
			JSONResponse(r, w, code, JSONresp)
			return
		}
		userAccount, err := users.CreateUser(db, user)
		if err == nil && invalid == false {
			code = 200
			response = "Created user account"
		}
		JSONresp := types.JSONMessageResponse{
			Metadata: types.JSONResponseMetadata{
				Response: response,
			},
			Spec: userAccount,
		}
		JSONResponse(r, w, code, JSONresp)
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	// root of API
	JSONresp := types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "Hey! you're talking to the Flattrack API",
		},
	}
	JSONResponse(r, w, 200, JSONresp)
}

func UnknownEndpoint(w http.ResponseWriter, r *http.Request) {
	JSONResponse(r, w, 404, types.JSONMessageResponse{
		Metadata: types.JSONResponseMetadata{
			Response: "This endpoint doesn't seem to exist.",
		},
	})
}

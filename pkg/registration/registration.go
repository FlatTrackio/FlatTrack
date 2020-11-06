/*
  registration
    manage the registration of the FlatTrack instance
*/

package registration

import (
	"database/sql"
	"gitlab.com/flattrack/flattrack/pkg/settings"
	"gitlab.com/flattrack/flattrack/pkg/system"
	"gitlab.com/flattrack/flattrack/pkg/types"
	"gitlab.com/flattrack/flattrack/pkg/users"
)

// Default groups
var (
	defaultInitalizationGroups = []string{"flatmember", "admin"}
)

// Register ...
// perform initial FlatTrack instance setup
func Register(db *sql.DB, registration types.Registration) (successful bool, jwt string, err error) {
	// TODO add timezone validation
	err = settings.SetTimezone(db, registration.Timezone)
	if err != nil {
		return successful, jwt, err
	}
	// TODO add language validation
	err = settings.SetTimezone(db, registration.Language)
	if err != nil {
		return successful, jwt, err
	}
	err = settings.SetFlatName(db, registration.FlatName)
	if err != nil {
		return successful, jwt, err
	}
	registration.User.Groups = defaultInitalizationGroups
	registration.User.Registered = true
	user, err := users.CreateUser(db, registration.User, false)
	if err != nil || user.ID == "" {
		return successful, jwt, err
	}
	jwt, err = users.GenerateJWTauthToken(db, user.ID, user.AuthNonce, 0)
	if err != nil {
		return successful, "", err
	}
	err = system.SetHasInitialized(db)
	if err != nil {
		return successful, "", err
	}
	return true, jwt, err
}

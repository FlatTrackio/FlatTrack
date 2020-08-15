/*
  registration
    manage the registration of the FlatTrack instance
*/

package registration

import (
	"database/sql"
	"gitlab.com/flattrack/flattrack/src/backend/settings"
	"gitlab.com/flattrack/flattrack/src/backend/system"
	"gitlab.com/flattrack/flattrack/src/backend/types"
	"gitlab.com/flattrack/flattrack/src/backend/users"
)

// Default groups
var (
	defaultInitalizationGroups = []string{"flatmember", "admin"}
)

// Register ...
// perform initial FlatTrack instance setup
func Register(db *sql.DB, registration types.Registration) (successful bool, err error) {
	// TODO add timezone validation
	err = settings.SetTimezone(db, registration.Timezone)
	if err != nil {
		return successful, err
	}
	// TODO add language validation
	err = settings.SetTimezone(db, registration.Language)
	if err != nil {
		return successful, err
	}
	err = settings.SetFlatName(db, registration.FlatName)
	if err != nil {
		return successful, err
	}
	registration.User.Groups = defaultInitalizationGroups
	registration.User.Registered = true
	user, err := users.CreateUser(db, registration.User, true)
	if err != nil || user.ID == "" {
		return successful, err
	}
	err = system.SetHasInitialized(db)
	if err != nil {
		return successful, err
	}
	return true, err
}

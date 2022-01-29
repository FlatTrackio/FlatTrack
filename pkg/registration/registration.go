/*
  registration
    manage the registration of the FlatTrack instance
*/

// This program is free software: you can redistribute it and/or modify
// it under the terms of the Affero GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the Affero GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
	userManager := users.UserManager{DB: db}
	systemManager := system.SystemManager{DB: db}
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
	jwt, err = userManager.GenerateJWTauthToken(user.ID, user.AuthNonce, 0)
	if err != nil {
		return successful, "", err
	}
	err = systemManager.SetHasInitialized()
	if err != nil {
		return successful, "", err
	}
	return true, jwt, err
}

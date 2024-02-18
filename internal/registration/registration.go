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
	"fmt"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/groups"
	"gitlab.com/flattrack/flattrack/internal/settings"
	"gitlab.com/flattrack/flattrack/internal/system"
	"gitlab.com/flattrack/flattrack/internal/users"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

// Default groups
var (
	defaultInitalizationGroups = []string{groups.GroupFlatmember, groups.GroupAdmin}
)

type Manager struct {
	user     *users.Manager
	system   *system.Manager
	settings *settings.Manager
	secret   string
}

func NewManager(users *users.Manager, system *system.Manager, settings *settings.Manager) *Manager {
	return &Manager{
		user:     users,
		system:   system,
		settings: settings,
		secret:   common.GetRegistrationSecret(),
	}
}

// Register ...
// perform initial FlatTrack instance setup
func (m *Manager) Register(registration types.Registration) (successful bool, jwt string, err error) {
	if m.secret != "" && registration.Secret != m.secret {
		return false, "", fmt.Errorf("a matching setup secret must be passed to registration")
	}
	// TODO add timezone validation
	err = m.settings.SetTimezone(registration.Timezone)
	if err != nil {
		return false, "", err
	}
	// TODO add language validation
	err = m.settings.SetTimezone(registration.Language)
	if err != nil {
		return false, "", err
	}
	err = m.settings.SetFlatName(registration.FlatName)
	if err != nil {
		return false, "", err
	}
	registration.User.Groups = defaultInitalizationGroups
	registration.User.Registered = true
	user, err := m.user.Create(registration.User, false)
	if err != nil || user.ID == "" {
		return false, "", err
	}
	jwt, err = m.user.GenerateJWTauthToken(user.ID, user.AuthNonce, 0)
	if err != nil {
		return false, "", err
	}
	err = m.system.SetHasInitialized()
	if err != nil {
		return false, "", err
	}
	return true, jwt, err
}

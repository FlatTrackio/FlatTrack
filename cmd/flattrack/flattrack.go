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

// Package classification for FlatTrack API.
//
//     Schemes: http
//     Host: localhost
//     BasePath: /api
//     Version: 0.16.1
//     License: AGPL-3.0 http://www.gnu.org/licenses/agpl-3.0.html
//     Contact: Caleb Woodbine <calebwoodbine.public@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta

// Package flattrack ...
// backend cmd
package flattrack

import (
	"gitlab.com/flattrack/flattrack/internal/flattrack"
)

func Run() {
	flattrack.NewManager().Init().Run()
}

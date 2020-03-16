// manage user accounts
package users

import (
	"database/sql"

	//"gitlab.com/flattrack/flattrack/src/backend/database"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

func CreateUser(db *sql.DB, user types.UserSpec) (userInserted types.UserSpec, err error) {
	sqlStatement := `
  insert into users (name, email, groups, password, phonenumber, creationtimestamp)
  values ($1, $2, $3, $4, $5, $5)
  returning (name, email, groups, password, phonenumber, creationtimestamp)`
  err = db.QueryRow(sqlStatement, user.Name, user.Email, user.Groups, user.Password, user.PhoneNumber, user.CreationTimestamp).Scan(&user)
	return userInserted, err
}

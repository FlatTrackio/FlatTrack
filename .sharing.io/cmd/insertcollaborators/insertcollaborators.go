package main

import (
	"encoding/json"
	"log"
	"os"

	"gitlab.com/flattrack/flattrack/pkg/common"
	"gitlab.com/flattrack/flattrack/pkg/database"
	"gitlab.com/flattrack/flattrack/pkg/migrations"
	"gitlab.com/flattrack/flattrack/pkg/types"
	"gitlab.com/flattrack/flattrack/pkg/users"
)

var (
	defaultUserAccountGroups = []string{"flatmember", "admin"}
	defaultUserPassword      = "P@ssw0rd123!"
)

func main() {
	log.Println("Preparing to add new guests")
	userAccountsRaw := os.Getenv("SHARINGIO_PAIR_GUESTS")
	if userAccountsRaw == "" {
		log.Fatalln("No guests declared to create")
		return
	}
	var userAccounts []types.UserSpec
	err := json.Unmarshal([]byte(userAccountsRaw), &userAccounts)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("Connecting to local Postgres instance")
	dbUsername := common.GetDBusername()
	dbPassword := common.GetDBpassword()
	dbHostname := common.GetDBhost()
	dbDatabase := common.GetDBdatabase()
	db, err := database.DB(dbUsername, dbPassword, dbHostname, dbDatabase)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = migrations.Migrate(db)
	if err != nil {
		log.Fatalln("migrations:", err)
		return
	}

	if len(userAccounts) == 0 {
		log.Fatalln("No user accounts found in environment")
		return
	}
	for _, user := range userAccounts {
		user.Groups = defaultUserAccountGroups
		user.Password = defaultUserPassword
		log.Printf("Inserting '%v'/'%v'\n", user.Names, user.Email)
		userInDB, err := users.CreateUser(db, user, false)
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Printf("Inserted '%v'/'%v'\n", user.Names, user.Email)
		userAccountPatch := types.UserSpec{
			Names:       user.Names,
			Email:       user.Email,
			Password:    defaultUserPassword,
			Birthday:    user.Birthday,
			PhoneNumber: user.PhoneNumber,
			Registered:  true,
		}
		log.Println("Patching account to be registered")
		userConfirmed, err := users.PatchProfileAdmin(db, userInDB.ID, userAccountPatch)
		if err != nil {
			log.Fatalln(err)
			return
		}
		if userConfirmed.ID == "" || userConfirmed.Registered == false {
			log.Fatalln("Failed to patch profile")
			return
		}
		log.Println("Completed processing account")
	}

	log.Println("It's collaboration time!")
}

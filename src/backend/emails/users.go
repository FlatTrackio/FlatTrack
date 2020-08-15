package emails

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"

	"gitlab.com/flattrack/flattrack/src/backend/common"
	"gitlab.com/flattrack/flattrack/src/backend/settings"
	"gitlab.com/flattrack/flattrack/src/backend/smtp"
	"gitlab.com/flattrack/flattrack/src/backend/types"
)

type SmtpTemplateAccountSignup struct {
	Subject            string
	InstanceURL        string
	CssFiles           []string
	FlatName           string
	User               types.UserSpec
	UserCreationSecret types.UserCreationSecretSpec
}

func SendAccountSignup(db *sql.DB, user types.UserSpec, userCreationSecret types.UserCreationSecretSpec) (err error) {
	flatName, err := settings.GetFlatName(db)
	if err != nil {
		return err
	}
	cssFiles, err := common.GetFileNamesFromFolder("dist/css")
	if err != nil {
		return err
	}
	context := &SmtpTemplateAccountSignup{
		Subject:            fmt.Sprintf("FlatTrack (%v) - account signup", flatName),
		FlatName:           flatName,
		InstanceURL:        common.GetInstanceURL(),
		CssFiles:           cssFiles,
		User:               user,
		UserCreationSecret: userCreationSecret,
	}
	emailReportTemplate, err := template.ParseFiles("templates/accountSignup.html")
	if err != nil {
		return err
	}
	templatedEmailBuffer := new(bytes.Buffer)
	err = emailReportTemplate.Execute(templatedEmailBuffer, context)
	if err != nil {
		log.Println(err)
		return err
	}
	if err != nil {
		return err
	}
	templateEmailRendered := templatedEmailBuffer.String()
	err = smtp.SendEmail(templateEmailRendered, context.Subject, user.Email)
	return err
}

package emails

import (
	"bytes"
	"gitlab.com/flattrack/flattrack/pkg/smtp"
	"html/template"
	"log"
)

// SMTPTemplateData ...
// basic email template
type SMTPTemplateData struct {
	Subject string
}

// SendTestEmail ...
// sends a test email from a template
func SendTestEmail(recipient string) (err error) {
	context := &SMTPTemplateData{
		Subject: "FlatTrack SMTP test",
	}
	emailReportTemplate, err := template.ParseFiles("templates/test.html")
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
	err = smtp.SendEmail(templateEmailRendered, context.Subject, recipient)
	return err
}

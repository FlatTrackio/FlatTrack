/*
  emails
    manage and send email alerts
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

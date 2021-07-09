/*
  smtp
    smtp management
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

package smtp

import (
	"fmt"
	"gitlab.com/flattrack/flattrack/pkg/common"
	"log"
	"net/smtp"
)

// EmailUser ...
// fields for authenticating to an SMTP server
type EmailUser struct {
	Username string
	Password string
	Host     string
	Port     string
}

// SendEmail ...
// send a HTML email to a subject
func SendEmail(content string, subject string, recipient string) error {
	smtpConfiguration := &EmailUser{
		Username: common.GetSMTPUsername(),
		Password: common.GetSMTPPassword(),
		Host:     common.GetSMTPHost(),
		Port:     common.GetSMTPPort(),
	}

	var auth smtp.Auth = nil
	if smtpConfiguration.Username != "" && smtpConfiguration.Password != "" {
		auth = smtp.PlainAuth("",
			smtpConfiguration.Username,
			smtpConfiguration.Password,
			smtpConfiguration.Host,
		)
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subjectFmt := fmt.Sprintf("Subject: %v\n", subject)
	msg := []byte(fmt.Sprintf("%v%v\n%v", subjectFmt, mime, content))
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", smtpConfiguration.Host, smtpConfiguration.Port),
		auth,
		smtpConfiguration.Username,
		[]string{recipient},
		msg,
	)

	if err != nil {
		log.Printf("Error: failed to send email mail %v\n", err)
		return err
	}
	return err
}

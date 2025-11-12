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
	"log/slog"
	"net/smtp"

	"gitlab.com/flattrack/flattrack/internal/common"
)

type Manager struct {
	username string
	password string
	host     string
	port     string

	auth smtp.Auth
}

func NewManager() *Manager {
	username := common.GetSMTPUsername()
	password := common.GetSMTPPassword()
	host := common.GetSMTPHost()
	port := common.GetSMTPPort()
	auth := smtp.PlainAuth("", username, password, host)

	return &Manager{
		username: username,
		password: password,
		host:     host,
		port:     port,

		auth: auth,
	}
}

// SendEmail ...
// send a HTML email to a subject
func (m *Manager) SendEmail(content string, subject string, recipient string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subjectFmt := fmt.Sprintf("Subject: %v\n", subject)
	msg := []byte(fmt.Sprintf("%v%v\n%v", subjectFmt, mime, content))
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", m.host, m.port),
		m.auth,
		m.username,
		[]string{recipient},
		msg,
	)

	if err != nil {
		slog.Error("Failed to send email mail", "error", err)
		return err
	}
	return err
}

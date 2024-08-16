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
	"html/template"
	"log"
	"net/url"
	"path"
	"time"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/settings"
	"gitlab.com/flattrack/flattrack/internal/smtp"
	"gitlab.com/flattrack/flattrack/internal/users"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

type Manager struct {
	Enabled       bool
	smtpManager   *smtp.Manager
	templatesPath string
	instanceURL   *url.URL
	flatName      string
	usersManager  *users.Manager
}

func NewManager(settingsManager *settings.Manager, usersManager *users.Manager) *Manager {
	flatName, _ := settingsManager.GetFlatName()
	return &Manager{
		Enabled:       common.GetSMTPEnabled(),
		smtpManager:   smtp.NewManager(),
		templatesPath: common.GetEmailTemplatesPath(),
		instanceURL:   common.GetInstanceURL(),
		flatName:      flatName,
		usersManager:  usersManager,
	}
}

func (m *Manager) templatePath(name string) string {
	return path.Join(m.templatesPath, name+".html")
}
func (m *Manager) renderTemplate(templateName string, context any) (string, error) {
	emailReportTemplate, err := template.ParseFiles(m.templatePath(templateName))
	if err != nil {
		return "", err
	}
	templatedEmailBuffer := new(bytes.Buffer)
	if err := emailReportTemplate.Execute(templatedEmailBuffer, context); err != nil {
		log.Println(err)
		return "", err
	}
	templateEmailRendered := templatedEmailBuffer.String()
	return templateEmailRendered, nil
}

type emailUserAccountCreated struct {
	Account           *types.UserSpec
	Confirm           *types.UserCreationSecretSpec
	InstanceDomain    string
	AccountConfirmURL string
	Subject           string
}

func (m *Manager) SendUserAccountCreated(account *types.UserSpec, confirm *types.UserCreationSecretSpec) error {
	templateName := "userAccountCreated"
	confirmURL := &url.URL{}
	*confirmURL = *m.instanceURL
	confirmURL.Path = path.Join("/useraccountconfirm", confirm.ID)
	params := url.Values{"secret": {confirm.Secret}}
	confirmURL.RawQuery = params.Encode()
	context := &emailUserAccountCreated{
		Account:           account,
		Confirm:           confirm,
		AccountConfirmURL: confirmURL.String(),
		InstanceDomain:    m.instanceURL.Hostname(),
		Subject:           "Welcome to FlatTrack!",
	}
	r, err := m.renderTemplate(templateName, context)
	if err != nil {
		return err
	}
	if err := m.smtpManager.SendEmail(r, context.Subject, account.Email); err != nil {
		return err
	}
	return nil
}

func (m *Manager) ReconcileUserCreatedEmails() error {
	if !m.Enabled {
		return nil
	}
	list, err := m.usersManager.UserCreationSecrets().List(types.UserCreationSecretSelector{})
	if err != nil {
		return err
	}
	for _, item := range list {
		if item.EmailSentStatus == types.EmailSentStatusNoAction {
			user, err := m.usersManager.GetByID(item.UserID, false)
			if err != nil {
				return err
			}
			log.Printf("sending email for user account creation %v for uid %v", item.ID, item.UserID)
			if err := m.SendUserAccountCreated(&user, &item); err != nil {
				return err
			}
			item.EmailSentDate = time.Now().Unix()
			item.EmailSentStatus = types.EmailSentStatusSentOnce
			log.Println("sent email", item.ID)
			if _, err := m.usersManager.UserCreationSecrets().Update(item); err != nil {
				return err
			}
		}
	}
	return nil
}

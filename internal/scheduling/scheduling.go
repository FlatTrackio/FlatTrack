package scheduling

import (
	"database/sql"

	"gitlab.com/flattrack/flattrack/internal/scheduling/leaderelection"
)

type Manager struct {
	db             *sql.DB
	leaderelection *leaderelection.Lock
}

func NewManager(db *sql.DB) *Manager {
	return &Manager{
		db:             db,
		leaderelection: leaderelection.NewLock(db),
	}
}

func (m *Manager) Run() {
	m.leaderelection.Run(func() error {
		return nil
	})
}

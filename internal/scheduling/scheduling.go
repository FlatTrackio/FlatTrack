package scheduling

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/scheduling/leaderelection"
)

type Manager struct {
	db              *sql.DB
	leaderelection  *leaderelection.Lock
	fns             []func() error
	endpointEnabled bool
	secret          string
}

func NewManager(db *sql.DB) *Manager {
	m := &Manager{
		db:              db,
		leaderelection:  leaderelection.NewLock(db),
		endpointEnabled: common.GetSchedulerDisableUseEndpoint(),
		secret:          common.GetSchedulerEndpointSecret(),
	}
	if m.endpointEnabled && m.secret == "" {
		log.Panicln("warning: APP_SCHEDULER_ENDPOINT_SECRET must be set when scheduler is disabled to ensure that only expected authorities call the scheduler endpoint")
	}
	return m
}

func (m *Manager) GetEndpointEnabled() bool {
	return m.endpointEnabled
}

func (m *Manager) GetEndpointSecret() string {
	return m.secret
}

func (m *Manager) RegisterFunc(fns ...func() error) *Manager {
	m.fns = append(m.fns, fns...)
	return m
}

func (m *Manager) PerformWork() error {
	var eg []error
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, fn := range m.fns {
		wg.Add(1)
		go func(fn func() error) {
			if err := fn(); err != nil {
				mu.Lock()
				eg = append(eg, err)
				mu.Unlock()
			}
			wg.Done()
		}(fn)
	}
	wg.Wait()
	if len(eg) > 0 {
		log.Printf("%+v\n", eg)
		return fmt.Errorf("scheduling errors: %v", eg)
	}
	return nil

}

func (m *Manager) Run() {
	if m.endpointEnabled {
		return
	}
	m.leaderelection.Run(m.PerformWork)
}

package scheduling

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/go-co-op/gocron/v2"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/scheduling/leaderelection"
)

type Manager struct {
	db              *sql.DB
	leaderelection  *leaderelection.Lock
	fns             []func() error
	endpointEnabled bool
	secret          string
	cronScheduler   gocron.Scheduler
}

func NewManager(db *sql.DB) *Manager {
	leaderelection := leaderelection.NewLock(db)
	cronScheduler, err := gocron.NewScheduler(
		gocron.WithLogger(
			gocron.NewLogger(gocron.LogLevelError),
		),
		gocron.WithDistributedElector(leaderelection),
	)
	if err != nil {
		log.Printf("Error creating scheduler: %v", err)
	}
	m := &Manager{
		db:              db,
		leaderelection:  leaderelection,
		endpointEnabled: common.GetSchedulerDisableUseEndpoint(),
		secret:          common.GetSchedulerEndpointSecret(),
		cronScheduler:   cronScheduler,
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

func (m *Manager) RegisterCronFunc(crontab string, fn func() error) *Manager {
	if !m.endpointEnabled {
		m.fns = append(m.fns, fn)
		return m
	}
	if _, err := m.cronScheduler.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(fn),
	); err != nil {
		log.Printf("Error: creating cron func: %v", err)
	}
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
	m.cronScheduler.Start()
	defer func() {
		if err := m.cronScheduler.Shutdown(); err != nil {
			log.Println(err)
		}
	}()
	m.leaderelection.Run(m.PerformWork)
}

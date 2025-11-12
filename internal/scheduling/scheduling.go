package scheduling

import (
	"database/sql"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"

	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/scheduling/leaderelection"
	"gitlab.com/flattrack/flattrack/internal/system"
	"gitlab.com/flattrack/flattrack/pkg/types"
)

type Manager struct {
	db              *sql.DB
	system          *system.Manager
	leaderelection  *leaderelection.Lock
	fns             []func() error
	endpointEnabled bool
	secret          string
	cronScheduler   gocron.Scheduler
}

func NewManager(db *sql.DB, system *system.Manager) *Manager {
	leaderelection := leaderelection.NewLock(db)
	cronScheduler, err := gocron.NewScheduler(
		gocron.WithLogger(
			gocron.NewLogger(gocron.LogLevelError),
		),
		gocron.WithDistributedElector(leaderelection),
	)
	if err != nil {
		slog.Info("Error creating scheduler", "error", err)
	}
	m := &Manager{
		db:              db,
		system:          system,
		leaderelection:  leaderelection,
		endpointEnabled: common.GetSchedulerDisableUseEndpoint(),
		secret:          common.GetSchedulerEndpointSecret(),
		cronScheduler:   cronScheduler,
	}
	if m.endpointEnabled && m.secret == "" {
		slog.Warn("APP_SCHEDULER_ENDPOINT_SECRET must be set when scheduler is disabled to ensure that only expected authorities call the scheduler endpoint")
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
	if m.endpointEnabled {
		m.fns = append(m.fns, fn)
		return m
	}
	if _, err := m.cronScheduler.NewJob(
		gocron.CronJob(crontab, false),
		gocron.NewTask(fn),
	); err != nil {
		slog.Info("Error: creating cron func", "error", err)
	}
	return m
}

func (m *Manager) PerformWork() error {
	slog.Debug("scheduler", "error", "Work running")
	now := time.Now()
	if err := m.system.SetSchedulerLastRun(types.SchedulerLastRun{
		Time:  now.Unix(),
		State: types.SchedulerRunStateRunning,
	}); err != nil {
		return err
	}
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
		slog.Debug("%+v\n", "errorCount", eg)
		if err := m.system.SetSchedulerLastRun(types.SchedulerLastRun{
			Time:  now.Unix(),
			State: types.SchedulerRunStateFailure,
		}); err != nil {
			return err
		}
		slog.Debug("scheduler", "error", "Work failed")
		return fmt.Errorf("scheduling errors: %v", eg)
	}
	if err := m.system.SetSchedulerLastRun(types.SchedulerLastRun{
		Time:  now.Unix(),
		State: types.SchedulerRunStateComplete,
	}); err != nil {
		return err
	}
	slog.Debug("scheduler", "error", "Work complete")
	return nil

}

func (m *Manager) Run() {
	if m.endpointEnabled {
		return
	}
	m.cronScheduler.Start()
	defer func() {
		if err := m.cronScheduler.Shutdown(); err != nil {
			slog.Error("Failed to shutdown cron scheduler", "error", err)
		}
	}()
	m.leaderelection.Run(m.PerformWork)
}

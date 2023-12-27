package flattrack

import (
	"log"

	"github.com/joho/godotenv"
	"gitlab.com/flattrack/flattrack/internal/common"
	"gitlab.com/flattrack/flattrack/internal/database"
	"gitlab.com/flattrack/flattrack/internal/emails"
	"gitlab.com/flattrack/flattrack/internal/groups"
	"gitlab.com/flattrack/flattrack/internal/health"
	"gitlab.com/flattrack/flattrack/internal/httpserver"
	"gitlab.com/flattrack/flattrack/internal/metrics"
	"gitlab.com/flattrack/flattrack/internal/migrations"
	"gitlab.com/flattrack/flattrack/internal/registration"
	"gitlab.com/flattrack/flattrack/internal/settings"
	"gitlab.com/flattrack/flattrack/internal/shoppinglist"
	"gitlab.com/flattrack/flattrack/internal/system"
	"gitlab.com/flattrack/flattrack/internal/users"
)

type Manager struct {
	httpserver   *httpserver.HTTPServer
	metrics      *metrics.Manager
	emails       *emails.Manager
	groups       *groups.Manager
	health       *health.Manager
	migrations   *migrations.Manager
	registration *registration.Manager
	settings     *settings.Manager
	system       *system.Manager
}

func NewManager() *Manager {
	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)
	db, err := database.Open()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		return nil
	}
	users := users.NewManager(db)
	shoppinglist := shoppinglist.NewManager(db)
	emails := emails.NewManager()
	groups := groups.NewManager(db)
	health := health.NewManager(db)
	migrations := migrations.NewManager(db)
	settings := settings.NewManager(db)
	system := system.NewManager(db)
	registration := registration.NewManager(users, system, settings)
	metrics := metrics.NewManager()
	httpserver := httpserver.NewHTTPServer(db, users, shoppinglist, emails, groups, health, migrations, registration, settings, system)
	return &Manager{
		httpserver:   httpserver,
		metrics:      metrics,
		emails:       emails,
		groups:       groups,
		health:       health,
		migrations:   migrations,
		registration: registration,
		settings:     settings,
		system:       system,
	}
}

type ManagerInit struct {
	*Manager
}

func (m *Manager) Init() *ManagerInit {
	if err := m.migrations.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return &ManagerInit{
		Manager: m,
	}
}

func (mi *ManagerInit) Run() {
	log.Printf("launching FlatTrack (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())
	go mi.metrics.Listen()
	go mi.health.Listen()
	mi.httpserver.Listen()
}

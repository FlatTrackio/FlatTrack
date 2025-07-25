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
	"gitlab.com/flattrack/flattrack/internal/scheduling"
	"gitlab.com/flattrack/flattrack/internal/settings"
	"gitlab.com/flattrack/flattrack/internal/shoppinglist"
	"gitlab.com/flattrack/flattrack/internal/system"
	"gitlab.com/flattrack/flattrack/internal/users"
)

type manager struct {
	httpserver   *httpserver.HTTPServer
	metrics      *metrics.Manager
	emails       *emails.Manager
	groups       *groups.Manager
	health       *health.Manager
	migrations   *migrations.Manager
	registration *registration.Manager
	settings     *settings.Manager
	system       *system.Manager
	scheduling   *scheduling.Manager

	maintenanceMode bool
}

func NewManager() *manager {
	log.Printf("launching FlatTrack (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())
	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)
	maintenanceMode := common.GetMaintenanceMode()
	db, err := database.Open()
	if err != nil && !maintenanceMode {
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
	scheduling := scheduling.NewManager(db)
	httpserver := httpserver.NewHTTPServer(db, users, shoppinglist, emails, groups, health, migrations, registration, settings, system, scheduling, maintenanceMode)
	return &manager{
		httpserver:      httpserver,
		metrics:         metrics,
		emails:          emails,
		groups:          groups,
		health:          health,
		migrations:      migrations,
		registration:    registration,
		settings:        settings,
		system:          system,
		scheduling:      scheduling,
		maintenanceMode: maintenanceMode,
	}
}

type managerInit struct {
	httpserver   *httpserver.HTTPServer
	metrics      *metrics.Manager
	health       *health.Manager
	registration *registration.Manager
	scheduling   *scheduling.Manager

	maintenanceMode bool
}

func (m *manager) Init() *managerInit {
	if err := m.migrations.Migrate(); err != nil && !m.maintenanceMode {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return &managerInit{
		httpserver:      m.httpserver,
		metrics:         m.metrics,
		health:          m.health,
		registration:    m.registration,
		scheduling:      m.scheduling,
		maintenanceMode: m.maintenanceMode,
	}
}

func (mi *managerInit) Run() {
	go mi.metrics.Listen()
	go mi.health.Listen()
	if !mi.maintenanceMode {
		go mi.scheduling.Run()
	} else {
		log.Println("Instance in maintenance mode. Will only serve message stating as such.")
	}
	mi.httpserver.Listen()
}

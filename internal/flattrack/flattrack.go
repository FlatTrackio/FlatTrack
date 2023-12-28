package flattrack

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	"github.com/lib/pq"
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
	db           *sql.DB
	listener     *pq.Listener
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
}

func NewManager() *manager {
	log.Printf("launching FlatTrack (%v, %v, %v, %v)\n", common.GetAppBuildVersion(), common.GetAppBuildHash(), common.GetAppBuildDate(), common.GetAppBuildMode())
	envFile := common.GetAppEnvFile()
	_ = godotenv.Load(envFile)
	db, err := database.Open()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
		return nil
	}
	listener := database.OpenListener()
	if err := listener.Listen("events"); err != nil {
		log.Fatalf("error listening to database events: %v\n", err)
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
	httpserver := httpserver.NewHTTPServer(db, listener, users, shoppinglist, emails, groups, health, migrations, registration, settings, system, scheduling)
	return &manager{
		db:           db,
		listener:     listener,
		httpserver:   httpserver,
		metrics:      metrics,
		emails:       emails,
		groups:       groups,
		health:       health,
		migrations:   migrations,
		registration: registration,
		settings:     settings,
		system:       system,
		scheduling:   scheduling,
	}
}

type managerInit struct {
	db           *sql.DB
	listener     *pq.Listener
	httpserver   *httpserver.HTTPServer
	metrics      *metrics.Manager
	health       *health.Manager
	registration *registration.Manager
	scheduling   *scheduling.Manager
}

func (m *manager) Init() *managerInit {
	if err := m.migrations.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	return &managerInit{
		db:           m.db,
		listener:     m.listener,
		httpserver:   m.httpserver,
		metrics:      m.metrics,
		health:       m.health,
		registration: m.registration,
		scheduling:   m.scheduling,
	}
}

func (mi *managerInit) Run() {
	go mi.metrics.Listen()
	go mi.health.Listen()
	go mi.scheduling.Run()
	defer func() {
		log.Println("closing database connection...")
		if err := database.Close(mi.db); err != nil {
			log.Fatalf("error closing database connection: %v", err)
		}
	}()
	defer func() {
		log.Println("closing pg listener connection...")
		if err := mi.listener.Close(); err != nil {
			log.Fatalf("error closing pg listener: %v", err)
		}
	}()
	mi.httpserver.Listen()
}

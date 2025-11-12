package flattrack

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"

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
	log          *slog.Logger
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

type FlatTrack struct {
	DB struct {
		ConnectionString string `mapstructure:"connection_string"`
		Database         string `mapstructure:"database"`
		Username         string `mapstructure:"username"`
		Host             string `mapstructure:"host"`
		Port             string `mapstructure:"port"`
		Password         string `mapstructure:"password"`
		sslMode          string `mapstructure:"ssl_mode"`
		migrationsPath   string `mapstructure:"migrationsPath"`
	} `mapstructure:"db"`
	URL  string `mapstructure:"url"`
	SMTP struct {
		Enabled  string `mapstructure:"enabled"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"Host"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"smtp"`
	Port    string `mapstructure:"port"`
	Metrics struct {
		Enabled string `mapstructure:"enabled"`
		Port    string `mapstructure:"port"`
	} `mapstructure:"metrics"`
	Health struct {
		Enabled string `mapstructure:"enabled"`
		Port    string `mapstructure:"port"`
	} `mapstructure:"health"`
	HTTPRealIPHeader string `mapstructure:"http_real_ip_header"`
	SetupMessage     string `mapstructure:"setup_message"`
	LoginMessage     string `mapstructure:"login_message"`
	EmbeddedHTML     string `mapstructure:"embedded_html"`
	Minio            struct {
		AccessKey string `mapstructure:"access_key"`
		SecretKey string `mapstructure:"secret_key"`
		Bucket    string `mapstructure:"bucket"`
		Host      string `mapstructure:"host"`
		UseSSL    bool   `mapstructure:"use_ssl"`
	} `mapstructure:"minio"`
	Scheduler struct {
		EndpointSecret     string `mapstructure:"endpoint_secret"`
		DisableUseEndpoint string `mapstructure:"disable_use_endpoint"`
	} `mapstructure:"scheduler"`
	RegistrationSecret string `mapstructure:"registration_secret"`
	MaintenanceMode    string `mapstructure:"maintenance_mode"`
	WebFolder          string `mapstructure:"web_folder"`
}

func NewManager() *manager {
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{AddSource: true, ReplaceAttr: common.SLogReplaceSource},
		)),
	)
	slog.Info("launching FlatTrack",
		slog.String("buildVersion", common.GetAppBuildVersion()),
		slog.String("buildHash", common.GetAppBuildHash()),
		slog.String("buildDate", common.GetAppBuildDate()),
		slog.String("buildMode", common.GetAppBuildMode()),
	)
	viper.SetEnvPrefix("app")
	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.BindEnv("port", "PORT"); err != nil {
		slog.Error("Unable bind config flag", "error", err)
		return nil
	}
	if err := viper.BindEnv("web_folder", "KO_DATA_PATH"); err != nil {
		slog.Error("Unable bind config flag", "error", err)
		return nil
	}
	viper.SetDefault("port", ":8080")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			slog.Error("Unable read in config", "error", err)
			return nil
		}
	}
	flattrack := &FlatTrack{}
	if err := viper.Unmarshal(flattrack); err != nil {
		slog.Error("Unable to decode FlatTrack config into struct", "error", err)
	}
	slog.Info("config", "flattrack", flattrack, "APP_PORT", viper.GetString("port"), "keys", viper.AllKeys())
	maintenanceMode := common.GetMaintenanceMode()
	db, err := database.Open()
	if err != nil && !maintenanceMode {
		slog.Error("failed to connect to database", "error", err)
		return nil
	}
	users := users.NewManager(db)
	settings := settings.NewManager(db)
	shoppinglist := shoppinglist.NewManager(db, settings)
	emails := emails.NewManager()
	groups := groups.NewManager(db)
	health := health.NewManager(db)
	migrations := migrations.NewManager(db)
	system := system.NewManager(db)
	registration := registration.NewManager(users, system, settings)
	metrics := metrics.NewManager()
	scheduling := scheduling.NewManager(db, system).
		RegisterCronFunc(shoppinglist.ShoppingList().DeleteCleanup())
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
	log          *slog.Logger
	httpserver   *httpserver.HTTPServer
	metrics      *metrics.Manager
	health       *health.Manager
	registration *registration.Manager
	scheduling   *scheduling.Manager

	maintenanceMode bool
}

func (m *manager) Init() *managerInit {
	if err := m.migrations.Migrate(); err != nil && !m.maintenanceMode {
		m.log.Error("failed to migrate database", "error", err)
	}
	return &managerInit{
		log:             m.log,
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
		mi.log.Info("Instance in maintenance mode. Will only serve message stating as such.")
	}
	mi.httpserver.Listen()
}

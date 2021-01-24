package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/sharmarajdaksh/yorpoll-api/internal/log"
)

type env string

const (
	// Prod env
	Prod env = "prod"
	// Stage env
	Stage env = "stage"
	// Test env
	Test env = "test"
	// Dev env
	Dev env = "dev"
	// Trace enabled env
	Trace env = "trace"
)

type database struct {
	Type        string
	Hostname    string
	Port        int
	Name        string
	Username    string
	password    string
	NetworkType string
}

type server struct {
	Host string
	Port int
}

type global struct {
	Env     env
	Logfile string
}

type subconfig interface {
	loadConfig() (subconfig, error)
}

// Config represents a shared global configuration object
type Config struct {
	Global   global
	Database database
	Server   server
}

// Password returns the private database password
func (d database) Password() string {
	return d.password
}

// DbType returns the configured database type for the application
func (c Config) DbType() string {
	return c.Database.Type
}

// LogFileName returns the name of the file to writes logs to
func LogFileName() string {
	return envOrDefault(logFile, defaults.Global.Logfile)
}

// ServerAddress returns the address to listen on
func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

var defaults = Config{
	Global: global{
		Env:     Dev,
		Logfile: "logs/server.log",
	},
	Database: database{
		Type:        "mysql",
		Hostname:    "localhost",
		Port:        3306,
		Name:        "mysql",
		Username:    "root",
		password:    "",
		NetworkType: "tcp",
	},
	Server: server{
		Host: "localhost",
		Port: 8910,
	},
}

var c Config

func envOrDefault(envVariable string, defaultValue string) string {
	e := os.Getenv(envVariable)
	if e == "" {
		return defaultValue
	}
	return e
}

func getEnvForString(e string) env {
	envs := []env{Dev, Trace, Test, Stage, Prod}
	for _, env := range envs {
		if e == string(env) {
			return env
		}
	}
	log.Logger.Warn().Msgf("provided project environment type \"%s\" is invalid: defaulting to %s", e, defaults.Global.Env)
	return defaults.Global.Env
}

func setGlobalLogLevel(e env) {
	envLogLevel := map[env]zerolog.Level{
		Trace: zerolog.TraceLevel,
		Dev:   zerolog.DebugLevel,
		Test:  zerolog.DebugLevel,
		Stage: zerolog.InfoLevel,
		Prod:  zerolog.WarnLevel,
	}

	zerolog.SetGlobalLevel(envLogLevel[e])
}

func (global) loadConfig() (global, error) {
	env := getEnvForString(envOrDefault(envEnv, string(defaults.Global.Env)))
	logf := envOrDefault(logFile, defaults.Global.Logfile)

	return global{Env: env, Logfile: logf}, nil
}

func (database) loadConfig() (database, error) {
	dbType := envOrDefault(envDatabaseType, defaults.Database.Type)
	dbHostname := envOrDefault(envDatabaseHost, defaults.Database.Hostname)
	defaultDbPort := strconv.Itoa(defaults.Database.Port)
	dbPort, err := strconv.Atoi(envOrDefault(envDatabasePort, defaultDbPort))
	if err != nil {
		return database{}, fmt.Errorf("loaded invalid configuration: non-numeric database port")
	}
	dbName := envOrDefault(envDatabaseName, defaults.Database.Name)
	dbUsername := envOrDefault(envDatabaseUser, defaults.Database.Username)
	dbPassword := envOrDefault(envDatabasePassword, defaults.Database.password)
	sHost := envOrDefault(envServerHost, defaults.Server.Host)
	if sHost == "" {
		sHost = defaults.Server.Host
	}
	dbNetType := envOrDefault(envDatabaseNetworkType, defaults.Database.NetworkType)

	return database{
		Type:        dbType,
		Hostname:    dbHostname,
		Port:        dbPort,
		Name:        dbName,
		Username:    dbUsername,
		password:    dbPassword,
		NetworkType: dbNetType,
	}, nil
}

func (server) loadConfig() (server, error) {
	sHost := envOrDefault(envServerHost, defaults.Server.Host)
	if sHost == "" {
		sHost = defaults.Server.Host
	}
	defaultServerPort := strconv.Itoa(defaults.Server.Port)
	sPort, err := strconv.Atoi(envOrDefault(envServerPort, defaultServerPort))
	if err != nil {
		return server{}, fmt.Errorf("loaded invalid configuration")
	}

	return server{
		Host: sHost,
		Port: sPort,
	}, nil
}

// Init Initializes application configuration
func Init() (*Config, error) {
	glblc, err := global{}.loadConfig()
	if err != nil {
		return &Config{}, err
	}

	setGlobalLogLevel(glblc.Env)

	dbc, err := database{}.loadConfig()
	if err != nil {
		return &Config{}, err
	}

	srvrc, err := server{}.loadConfig()
	if err != nil {
		return &Config{}, err
	}

	c = Config{
		Global:   glblc,
		Database: dbc,
		Server:   srvrc,
	}

	return &c, nil
}

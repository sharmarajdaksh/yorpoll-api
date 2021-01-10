package db

import (
	"github.com/sharmarajdaksh/yorpoll-api/config"
	"github.com/sharmarajdaksh/yorpoll-api/internal/db/mongodb"
	"github.com/sharmarajdaksh/yorpoll-api/internal/db/mysql"
	"github.com/sharmarajdaksh/yorpoll-api/internal/poll"
)

// Connection is an interface to connect and interact with a database
type Connection interface {
	GetConnectionString(*config.Config) string
	Connect(*config.Config) error
	Close() error
	Ping() error
	poll.Repository
}

// Init returns a Connection object based on current application configuration
func Init(c *config.Config) Connection {
	if c.DbType() == mysqldb {
		return &mysql.Connection{}
	}
	if c.DbType() == mongo {
		return &mongodb.Connection{}
	}

	return nil
}

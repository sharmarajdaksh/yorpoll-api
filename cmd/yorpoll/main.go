package main

import (
	"os"
	"os/signal"

	"github.com/sharmarajdaksh/yorpoll-api/config"
	"github.com/sharmarajdaksh/yorpoll-api/internal/db"
	"github.com/sharmarajdaksh/yorpoll-api/internal/log"
	"github.com/sharmarajdaksh/yorpoll-api/internal/server"
)

func main() {
	log.Logger.Debug().Msg("Initializing yorpoll server")

	err := log.Configure(config.LogFileName())
	checkErr(err, "failed to configure logger")
	log.Logger.Debug().Str("event", "logger configured").Msg("succesfully configured logger")

	conf, err := config.Init()
	checkErr(err, "failed to load application config")
	log.Logger.Debug().Str("event", "loaded config").Interface("config", conf).Msg("loaded config")

	// Establish database connectivity
	dbc := db.Init(conf)
	err = dbc.Connect(conf)
	checkErr(err, "failed to establish database connection")
	defer func() {
		if err = dbc.Close(); err != nil {
			checkErr(err, "could not close established database connection")
		}
	}()

	log.Logger.Debug().Str("event", "connected to db").Msg("database connection established")

	srv := server.Init(conf, dbc)
	addr := conf.ServerAddress()

	// Listen for SIGINT to shut down gracefully
	sgnl := make(chan os.Signal, 1)
	signal.Notify(sgnl, os.Interrupt)
	go func() {
		_ = <-sgnl
		log.Logger.Debug().Msg("gracefully shutting down server")
		_ = srv.Shutdown()
	}()

	log.Logger.Debug().Msgf("starting server at %s", addr)
	err = srv.Listen(addr)
	checkErr(err, "could not bind to address "+addr)
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Logger.Panic().Err(err).Msg(msg)
	}
}

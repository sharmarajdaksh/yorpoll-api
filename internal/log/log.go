package log

import (
	"fmt"
	"os"

	zerolog "github.com/rs/zerolog"
	// "github.com/sharmarajdaksh/yorpoll-api/config"
)

// Logger is the global shared logger instance
var Logger zerolog.Logger

// Configure initializes the global logger with
func Configure(logf string) error {

	cnsl := zerolog.ConsoleWriter{Out: os.Stdout}

	f, err := os.OpenFile(logf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	multi := zerolog.MultiLevelWriter(cnsl, f)

	Logger = zerolog.New(multi).With().Timestamp().Logger()

	return nil
}

// Error is a helper function to log errors
func Error(err error) {
	Logger.Error().Err(err).Send()
}

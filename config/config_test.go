package config

import (
	"os"
	"testing"
)

var globalConfigValidTestCases = []struct {
	Env        env
	Logfile    string
	DefaultEnv bool
}{
	{
		Env:        Prod,
		Logfile:    "logs/log.log",
		DefaultEnv: false,
	},
	{
		Env:        Trace,
		Logfile:    "logs/trace.log",
		DefaultEnv: false,
	},
	{
		Env:        Dev,
		Logfile:    "dev.log",
		DefaultEnv: false,
	},
	{
		Env:        "non-existent",
		Logfile:    "dev.log",
		DefaultEnv: true,
	},
	{
		Env:        "non existent",
		Logfile:    "log/logger.log",
		DefaultEnv: true,
	},
}

func TestGlobalLoadConfig(t *testing.T) {
	for _, tglbl := range globalConfigValidTestCases {
		t.Run(string(tglbl.Env), func(t *testing.T) {
			os.Setenv(envEnv, string(tglbl.Env))
			os.Setenv(envLogFile, string(tglbl.Logfile))

			lglbl, err := global{}.loadConfig()
			if err != nil {
				t.Errorf("Unexpected error in loadConfig %s", err.Error())
			}

			if lglbl.Logfile != tglbl.Logfile {
				t.Errorf("Loaded logfile config not as given via environment variable")
			}

			if tglbl.DefaultEnv && lglbl.Env != defaults.Global.Env {
				t.Errorf("Did not default to default in case of invalid environment specified")
			}

			if !tglbl.DefaultEnv && lglbl.Env != tglbl.Env {
				t.Errorf("Loaded env config not as given via environment variable")
			}
		})
	}
}

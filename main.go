package main

import (
	"github.com/alecthomas/kong"
	"github.com/gentoomaniac/logging"
	"github.com/rs/zerolog/log"
)

var (
	version = "0.0.1"
)

var cli struct {
	logging.LoggingConfig

	Version kong.VersionFlag `short:"v" help:"Display version."`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Vars{
		"version": version,
	})
	logging.Setup(&cli.LoggingConfig)

	log.Info().Msg("Starting game")
	run()

	ctx.Exit(0)
}

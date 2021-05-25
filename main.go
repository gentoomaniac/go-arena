package main

import (
	"os"
	"runtime/pprof"

	"github.com/alecthomas/kong"
	"github.com/gentoomaniac/logging"
	"github.com/rs/zerolog/log"
)

var (
	version = "0.0.1"
)

var cli struct {
	logging.LoggingConfig

	Bot           []string `short:"b" help:"add another bot with this filename to the arena" required:""`
	ProfileMemory string   `help:"write a memory profile"`
	ProfileCPU    string   `help:"write a cpu profile"`

	Version kong.VersionFlag `short:"v" help:"Display version."`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Vars{
		"version": version,
	})
	logging.Setup(&cli.LoggingConfig)

	log.Info().Msg("Starting game")

	if cli.ProfileCPU != "" {
		f, err := os.Create(cli.ProfileCPU)
		if err != nil {
			log.Error().Err(err).Msg("could not create cpu profile")
			ctx.Exit(1)
		}

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Error().Err(err).Msg("could not start cpu profile")
			ctx.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	run(cli.Bot)

	if cli.ProfileMemory != "" {
		f, err := os.Create(cli.ProfileMemory)
		if err != nil {
			log.Error().Err(err).Msg("could not create memory profile")
			ctx.Exit(1)
		}
		defer f.Close()

		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Error().Err(err).Msg("could not start memory profile")
			ctx.Exit(1)
		}
	}

	pprof.StopCPUProfile()
	ctx.Exit(0)
}

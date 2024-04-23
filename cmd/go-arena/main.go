package main

import (
	"os"
	"runtime/pprof"

	"github.com/alecthomas/kong"
	"github.com/gentoomaniac/logging"
	"github.com/rs/zerolog/log"

	go_arena "github.com/gentoomaniac/go-arena/pkg/go-arena"
)

var (
	version = "0.0.1"
)

var cli struct {
	logging.LoggingConfig

	Bot      []string `short:"b" help:"add another bot with this filename to the arena" required:""`
	Respawns int      `short:"r" help:"Number of respawns"`

	MapPath string `short:"m" help:"Path to the map file" default:"maps\test.tmx"`

	ProfileMemory string `help:"write a memory profile"`
	ProfileCPU    string `help:"write a cpu profile"`
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

	go_arena.Run(go_arena.RunOpts{Bots: cli.Bot, Respawns: cli.Respawns, MapPath: cli.MapPath})

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

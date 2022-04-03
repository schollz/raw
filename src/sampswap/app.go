package sampswap

import (
	"os"

	"github.com/urfave/cli/v2"
)

func App() (err error) {
	ss := Init()
	app := &cli.App{
		Name:    "raw",
		Version: "v0.0.1",
		Usage:   "random audio workstation",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Required:    true,
				Name:        "file-in",
				Usage:       "input file",
				Destination: &ss.FileIn,
			},
			&cli.StringFlag{
				Required:    true,
				Name:        "file-out",
				Usage:       "output file",
				Destination: &ss.FileOut,
			},
			&cli.Float64Flag{
				Name:        "tempo-in",
				Usage:       "tempo of input (will be estimated if 0)",
				Value:       0,
				Destination: &ss.TempoIn,
			},
			&cli.Float64Flag{
				Name:        "tempo-out",
				Usage:       "final tempo of output",
				Value:       0,
				Destination: &ss.TempoOut,
			},
			&cli.Float64Flag{
				Name:        "beats",
				Usage:       "number of beats to render",
				Value:       16,
				Destination: &ss.BeatsOut,
			},
			&cli.Float64Flag{
				Name:        "jump",
				Usage:       "probability for jump",
				Value:       0,
				Destination: &ss.ProbJump,
			},
			&cli.Float64Flag{
				Name:        "stutter",
				Usage:       "probability for stutter",
				Value:       0,
				Destination: &ss.ProbStutter,
			},
			&cli.Float64Flag{
				Name:        "reverse",
				Usage:       "probability for reverse",
				Value:       0,
				Destination: &ss.ProbReverse,
			},
			&cli.Float64Flag{
				Name:        "rereverb",
				Usage:       "probability for reversed reverb",
				Value:       0,
				Destination: &ss.ProbRereverb,
			},
			&cli.Float64Flag{
				Name:        "filter-in",
				Usage:       "beats for filter ramp up at start",
				Value:       0,
				Destination: &ss.FilterIn,
			},
			&cli.Float64Flag{
				Name:        "filter-out",
				Usage:       "beats for filter ramp down at end",
				Value:       0,
				Destination: &ss.FilterOut,
			},
			&cli.Float64Flag{
				Name:        "sidechain",
				Usage:       "add sidechain every X beats",
				Value:       0,
				Destination: &ss.Sidechain,
			},
			&cli.BoolFlag{
				Name:        "tapedeck",
				Usage:       "process final output with tape emulator",
				Destination: &ss.Tapedeck,
			},
			&cli.BoolFlag{
				Name:        "tempo-ignore-pitch",
				Usage:       "ignores pitch when re-tempoing",
				Destination: &ss.ReTempoSpeed,
			},
			&cli.BoolFlag{
				Name:        "tempo-ignore-all",
				Usage:       "ignores re-tempoing",
				Destination: &ss.ReTempoNone,
			},
		},
		Action: func(c *cli.Context) error {
			return ss.Run()
		},
	}

	err = app.Run(os.Args)
	return
}

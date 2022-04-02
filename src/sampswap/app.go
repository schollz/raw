package sampswap

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func App() (err error) {
	ss := Init()
	app := &cli.App{
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
				Destination: &ss.BeatsIn,
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
		},
		Action: func(c *cli.Context) error {
			fmt.Printf("%+v", ss)
			return nil
		},
	}

	err = app.Run(os.Args)
	return
}
